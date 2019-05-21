## RECORDING RULES
---
用于记录规则的一致的命名方案使得更容易一目了然地解释规则的含义。它也避免错误，使错误或无意义的计算脱颖而出。

本页面介绍了如何正确地进行聚合，并提出了一个命名约定。

### 命名和聚合
记录规则应为一般形式`level：metric：operations`。 `level`表示规则输出的聚合级别和标签。 `metric`是度量名称，并且在使用`rate()`或`irate()`时除了剥离`_total` off计数器之外应该保持不变。操作是应用于度量的操作的列表，首先是最新的操作。

保持度量标准名称不变，可以轻松地了解代码库中的指标，并且易于查找。

为了保持操作清洁，如果有其他操作，则`_sum`被省略为`sum()`。关联操作可以合并（例如`min_min`与`min`相同）。

如果没有明显的操作使用，使用`sum`。通过进行分割取一个比率时，使用`_per_`分离度量，并调用`ratio`。

当汇总比率时，分别分解分子和分母，然后除以。 不要采用平均值的平均值，而不是统计学上的平均值。

当汇总总结的`_count`和`_sum`并将其划分以计算平均观察大小时，将其视为比例将是笨重的。 相反，保留度量名称不带`_count`或`_sum`后缀，并用平均值替换操作中的速率。 这表示该时间段内的平均观察尺寸。

总是使用您要聚合的标签指定一个无条款。 这是为了保留所有其他标签，如作业，这将避免冲突，并为您提供更有用的指标和警报。

### 例子
聚合每秒请求的label标签：
```
instance_path:requests:rate5m =
  rate(requests_total{job="myjob"}[5m])

path:requests:rate5m =
  sum without (instance)(instance_path:requests:rate5m{job="myjob"})
```

计算请求失败率并聚合到作业级失败率：
```
instance_path:request_failures:rate5m =
  rate(request_failures_total{job="myjob"}[5m])

instance_path:request_failures_per_requests:ratio_rate5m =
    instance_path:request_failures:rate5m{job="myjob"}
  /
    instance_path:requests:rate5m{job="myjob"}

// Aggregate up numerator and denominator, then divide to get path-level ratio.
path:request_failures_per_requests:ratio_rate5m =
    sum without (instance)(instance_path:request_failures:rate5m{job="myjob"})
  /
    sum without (instance)(instance_path:requests:rate5m{job="myjob"})

// No labels left from instrumentation or distinguishing instances,
// so we use 'job' as the level.
job:request_failures_per_requests:ratio_rate5m =
    sum without (instance, path)(instance_path:request_failures:rate5m{job="myjob"})
  /
    sum without (instance, path)(instance_path:requests:rate5m{job="myjob"})
```

从一个摘要计算一段时间内的平均延迟：
```
instance_path:request_latency_seconds_count:rate5m =
  rate(request_latency_seconds_count{job="myjob"}[5m])

instance_path:request_latency_seconds_sum:rate5m =
  rate(request_latency_seconds_sum{job="myjob"}[5m])

instance_path:request_latency_seconds:mean5m =
    instance_path:request_latency_seconds_sum:rate5m{job="myjob"}
  /
    instance_path:request_latency_seconds_count:rate5m{job="myjob"}

// Aggregate up numerator and denominator, then divide.
path:request_latency_seconds:mean5m =
    sum without (instance)(instance_path:request_latency_seconds_sum:rate5m{job="myjob"})
  /
    sum without (instance)(instance_path:request_latency_seconds_count:rate5m{job="myjob"})
```

在实例和路径之间计算平均查询速率是使用avg（）函数完成的：
```
job:request_latency_seconds_count:avg_rate5m =
  avg without (instance, path)(instance:request_latency_seconds_count:rate5m{job="myjob"})
```

请注意，在聚合时，与输入度量名称相比，将从输出度量名称的级别中删除without子句中的标签。 当没有聚合时，级别总是匹配的。 如果不是这样，规则中可能会出现错误。

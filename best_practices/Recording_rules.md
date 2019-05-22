## RECORDING RULES
---
用于[记录规则](https://prometheus.io/docs/prometheus/latest/configuration/recording_rules/)的一致命名方案使得更容易一目了然地解释规则的含义。 它还通过使错误或无意义的计算脱颖而出来避免错误。

本页介绍了如何正确进行聚合并建议命名约定。

### 命名和聚合
录制规则应具有一般形式`level:metric:operations`。 `level`表示聚合级别和规则输出的标签。` metric`是度量标准名称，除了在使用`rate()`或`irate()`时剥离`_total` off计数器之外，应该保持不变。 operations是首先应用于度量标准，最新操作的操作列表。

保持度量标准名称不变可以很容易地知道度量标准是什么，并且很容易在代码库中找到。

为了保持操作干净，如果有其他操作，则省略`_sum`，如`sum()`。可以合并关联操作（例如，`min_min`与`min`相同）。

如果没有明显的操作要使用，请使用`sum`。通过划分获取比率时，使用`_per_`分隔指标并调用操作`rate`。

在汇总比率时，分别汇总分子和分母然后除。不要取平均值的平均值或平均值，因为这在统计上无效。

当汇总摘要的`_count`和`_sum`并除以计算平均观察大小时，将其视为比率将是笨重的。而是保持度量标准名称不带`_count`或`_sum`后缀，并将操作中的`rate`替换为`mean`。这表示该时间段内的平均观察大小。

始终使用要聚合的标签指定`without`子句。这是为了保留所有其他标签，例如`job`，这将避免冲突并为您提供更有用的指标和警报。

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

在实例和路径之间计算平均查询速率是使用`avg()`函数完成的：
```
job:request_latency_seconds_count:avg_rate5m =
  avg without (instance, path)(instance:request_latency_seconds_count:rate5m{job="myjob"})
```

请注意，在聚合时，与输入度量标准名称相比，将从输出度量标准名称的级别中删除`without`子句中的标签。 如果没有聚合，则级别始终匹配。 如果不是这种情况，则规则中可能存在错误。

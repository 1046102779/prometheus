## 查询例子
---
### 简单的时间序列选择
返回度量指标名称是`http_requests_total`的所有时间序列样本数据：
> http_requests_total

返回度量指标名称是`http_requests_total`, 标签分别是`job="apiserver`, `handler="/api/comments`的所有时间序列样本数据：
> http_requests_total{job="apiserver", hanlder="/api/comments"}

返回度量指标名称是`http_requests_total`, 标签分别是`job="apiserver`, `handler="/api/comments`，且是5分钟内的所有时间序列样本数据：
> http_requests_total{job="apiserver", handler="/api/comments"}[5m]

注意：一个范围向量表达式结果不能直接在Graph图表中，但是可以在"console"视图中展示。

使用正则表达式，你可以通过特定模式匹配标签为job的特定任务名，获取这些任务的时间序列。在下面这个例子中, 所有任务名称以`server`结尾。
> http_requests_total{job=~"server$"}

返回度量指标名称是`http_requests_total`， 且http返回码不以4开头的所有时间序列数据：
> http_requests_total{status!~"^4..$"}

### 使用函数，操作符等
返回度量指标名称`http_requests_total`，且过去5分钟的所有时间序列数据值速率。
> rate(http_requests_total[5m])

假设度量名称是`http_requests_total`，且过去5分钟的所有时间序列数据的速率和，并保留输出时间序列标签名称`job`
> sum(rate(http_requests_total)[5m]) by (job)

如果我们有相同维度标签，但是不同的度量指标名称，我们可以使用二元操作符。具有相同标签集合的元素将会输出。例如，下面这个表达式返回每一个实例剩余内存，单位是M, 如果不同，则需要使用`ignoring(label_lists)`，如果多对一，则采用group_left, 如果是一对多，则采用group_right。
> (instance_memory_limit_byte - instant_memory_usage_bytes) / 1024 / 1024

相同表达式，求和可以采用下面表达式：
> sum( instance_memory_limit_bytes - instance_memory_usage_bytes) by (app, proc) / 1024 / 1024

如果相同集群调度器任务，显示CPU使用率度量指标的话，如下所示：
> instance_cpu_time_ns{app="lion", pro="web", rev="34d0f99", env="prod", job="cluster-manager"}
> instance_cpu_time_ns{app="elephant", proc="worker", rev="34d0f99", env="prod", job="cluster-manager"}
> instance_cpu_time_ns{app="turtle", proc="api", rev="4d3a513", env="prod", job="cluster-manager"}
> ...

我们可以获取最高的3个CPU使用率，且结果中带有`app`和`proc`的标签时间序列。

假设一个服务实例只有一个时间序列数据，那么我们通过下面表达式，可以统计出每个应用的实例数量：
> count(instance_cpu_time_ns) by (app)

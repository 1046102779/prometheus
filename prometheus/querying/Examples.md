##### 一、简单的时间序列选择
使用度量标准`http_requests_total`返回所有时间序列：
> http_requests_total

使用度量标准`http_requests_total`以及给定的`job`和`handler`标签返回所有时间系列：
> http_requests_total{job="apiserver", hanlder="/api/comments"}

返回相同向量的整个时间范围（在本例中为5分钟），使其成为范围向量：
> http_requests_total{job="apiserver", handler="/api/comments"}[5m]

请注意，导致范围向量的表达式不能直接绘制，而是在表达式浏览器的表格（"Console"）视图中查看。

使用正则表达式，您只能为名称与特定模式匹配的作业选择时间序列，在本例中为所有以`server`结尾的作业。 请注意，这会进行子字符串匹配，而不是完整的字符串匹配：
> http_requests_total{job=~"server$"}

Prometheus中的所有正则表达式都使用[RE2语法](https://github.com/google/re2/wiki/Syntax)。

要选择除4xx之外的所有HTTP状态代码，您可以运行：
> http_requests_total{status!~"^4..$"}

##### 二、子查询
此查询返回过去30分钟的5分钟`http_requests_total`指标率，分辨率为1分钟。
> rate(http_requests_total[5m])[30m:1m]

这是嵌套子查询的示例。 `deri`函数的子查询使用默认分辨率。 请注意，不必要地使用子查询是不明智的。
> max_over_time(deriv(rate(distance_covered_total[5s])[30s:5s])[10m:])

##### 三、使用函数，操作符等
使用`http_requests_total`指标名称返回所有时间序列的每秒速率，在过去5分钟内测量：
> rate(http_requests_total[5m])

假设`http_requests_total`时间序列都有标签`job`（按作业名称扇出）和`instance`（按作业实例扇出），我们可能想要总结所有实例的速率，因此我们得到的输出时间序列更少，但仍然 保留`job`维度：
> sum(rate(http_requests_total)[5m]) by (job)

如果我们有两个具有相同维度标签的不同指标，我们可以对它们应用二元运算符，并且两侧具有相同标签集的元素将匹配并传播到输出。 例如，此表达式为每个实例返回MiB中未使用的内存（在虚构的群集调度程序上公开它运行的实例的这些度量标准）：
> (instance_memory_limit_byte - instant_memory_usage_bytes) / 1024 / 1024

相同的表达式，但由应用程序总结，可以这样写：
> sum( instance_memory_limit_bytes - instance_memory_usage_bytes) by (app, proc) / 1024 / 1024

如果相同的虚构集群调度程序为每个实例公开了如下所示的CPU使用率指标：
> instance_cpu_time_ns{app="lion", pro="web", rev="34d0f99", env="prod", job="cluster-manager"}
> instance_cpu_time_ns{app="elephant", proc="worker", rev="34d0f99", env="prod", job="cluster-manager"}
> instance_cpu_time_ns{app="turtle", proc="api", rev="4d3a513", env="prod", job="cluster-manager"}
> ...

...我们可以按应用程序（`app`）和进程类型（`proc`）分组排名前3位的CPU用户：
> topk(3, sum(rate(instance_cpu_time_ns[5m])) by(app, proc))

假设此度量标准包含每个运行实例的一个时间系列，您可以计算每个应用程序运行实例的数量，如下所示：
> count(instance_cpu_time_ns) by (app)

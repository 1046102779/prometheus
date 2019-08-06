Prometheus从根本上存储的所有数据都是[时间序列](http://en.wikipedia.org/wiki/Time_series): 具有时间戳的数据流只属于单个度量指标和该度量指标下的多个标签维度。除了存储时间序列数据外，Prometheus还可以生成临时派生的时间序列作为查询的结果。

##### 一、metrics和labels(度量指标名称和标签)
每一个时间序列数据由`metric`度量指标名称和它的标签`labels`键值对集合唯一确定。

这个`metric`度量指标名称指定监控目标系统的测量特征（如：`http_requests_total`- 接收http请求的总计数）。 `metric`度量指标可能包含ASCII字母、数字、下划线和冒号，他必须配正则表达式`[a-zA-Z_:][a-zA-Z0-9_:]*`
。

注意：冒号保留用于用户定义的录制规则。 它们不应被`exporter`或直接仪表使用。

`labels`开启了Prometheus的多维数据模型：对于相同的度量名称，通过不同标签列表的结合, 会形成特定的度量维度实例。(例如：所有包含度量名称为`/api/tracks`的http请求，打上`method=POST`的标签，则形成了具体的http请求)。这个查询语言在这些度量和标签列表的基础上进行过滤和聚合。改变任何度量上的任何标签值，则会形成新的时间序列图。

标签名称可以包含ASCII字母、数字和下划线。它们必须匹配正则表达式`[a-zA-Z_][a-zA-Z0-9_]*`。带有`_`下划线的标签名称被保留内部使用。

标签值包含任意的Unicode码。

具体详见[metrics和labels命名最佳实践](https://prometheus.io/docs/practices/naming/)。

##### 二、样本
样本形成了实际的时间序列数据列表。每个采样值包括：
 - 一个`float64`类型的值
 - 一个精确到毫秒级的时间戳

##### 三、符号
表示一个度量指标和一组键值对标签，需要使用以下符号：
>  [metric name]{[label name]=[label value], ...}

例如，度量指标名称是`api_http_requests_total`， 标签为`method="POST"`, `handler="/messages"` 的示例如下所示：
> api_http_requests_total{method="POST", handler="/messages"}

这些命名和OpenTSDB使用方法是一样的。

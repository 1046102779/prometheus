## 数据模型
---
Prometheus从根本上存储的所有数据都是[时间序列](http://en.wikipedia.org/wiki/Time_series): 具有时间戳的数据流只属于单个度量指标和该度量指标下的多个标签维度。除了存储时间序列数据外，Prometheus也可以利用查询表达式存储5分钟的返回结果中的时间序列数据
### metrics和labels(度量指标名称和标签)
每一个时间序列数据由metric度量指标名称和它的标签labels键值对集合唯一确定。

这个metric度量指标名称指定监控目标系统的测量特征（如：`http_requests_total`- 接收http请求的总计数）. metric度量指标命名ASCII字母、数字、下划线和冒号，他必须配正则表达式`[a-zA-Z_:][a-zA-Z0-9_:]*`。

标签labels是指Prometheus多维度数据模型：任何相同的度量指标的多个维度标签结合标识了一个度量指标实例(例如：所有包含method=`POST`， URL=`/api/tracks`的HTTP请求)。在这些标签维度上，查询语言可以做过滤和聚合操作。添加或者移除一个标签会创建新的时间序列数据。

标签label名称可以包含ASCII字母、数字和下划线。它们必须匹配正则表达式`[a-zA-Z_][a-zA-Z0-9_]*`。带有`_`下划线的标签名称被保留内部使用。

标签labels值包含任意的Unicode码。

具体详见[metrics和labels命名最佳实践](https://prometheus.io/docs/practices/naming/)。

### Samples(样本)
样本形成了时间序列数据。每一个样本包括：
>  [metric name]{[label name]=[label value], ...}

例如，度量指标名称是`api_http_requests_total`， 标签为`method="POST"`, `handler="/messages"` 的示例如下所示：
> api_http_requests_total{method="POST", handler="/messages"}

这些命名和OpenTSDB使用方法是一样的


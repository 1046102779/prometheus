Prometheus客户端库提供了四个核心的metrics类型。这四种类型目前仅在客户端库和wire协议中区分。Prometheus服务还没有充分利用这些类型。不久的将来就会发生改变。

##### 一、Counter(计数器)
*counter* 是表示单个单调递增计数器的累积度量，其值只能在重启时增加或重置为零。 例如，您可以使用计数器来表示所服务的请求数，已完成的任务或错误。

不要使用计数器来暴露可能减少的值。例如，不要使用计数器来处理当前正在运行的进程数; 而是使用仪表。

客户端使用计数器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Counter.java)
 - [Python](https://github.com/prometheus/client_python#counter)
 - [Ruby](https://github.com/prometheus/client_ruby#counter)

##### 二、Gauge(测量器)
*gauge*是一个度量指标，它表示一个既可以递增, 又可以递减的值。

测量器主要测量类似于温度、当前内存使用量等，也可以统计当前服务运行随时增加或者减少的Goroutines数量

客户端使用计量器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Gauge.java)
 - [Python](https://github.com/prometheus/client_python#gauge)
 - [Ruby](https://github.com/prometheus/client_ruby#gauge)

##### 三、Histogram(柱状图)
*histogram*，直方图对观察结果进行采样（通常是请求持续时间或响应大小等），并将其计入可配置存储桶中。它还提供所有观察值的总和。

基本度量标准名称为<basename>的直方图在scrape期间显示多个时间序列：

 - 暴露的观察桶的累积计数器：`<basename>_bucket{le="<upper inclusive bound>"}`
 - 所有观测值的总和：`<basename>_sum`
 - 已观察到的事件数：`<basename>_count`，和`<basename>_bucket{le="+Inf"}`相同

使用[histogram_quantile](https://prometheus.io/docs/querying/functions/#histogram_quantile)函数, 计算直方图或者是直方图聚合计算的分位数阈值。 一个直方图计算[Apdex值](http://en.wikipedia.org/wiki/Apdex)也是合适的, 当在buckets上操作时，记住直方图是累计的。详见[直方图和总结](https://prometheus.io/docs/practices/histograms)

客户库的直方图使用文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Histogram.java)
 - [Python](https://github.com/prometheus/client_python#histogram)
 - [Ruby](https://github.com/prometheus/client_ruby#histogram)

##### 四、[Summary]总结
类似*histogram*柱状图，*summary*是采样点分位图统计(通常是请求持续时间和响应大小等)。虽然它还提供观察的总数和所有观测值的总和，但它在滑动时间窗口上计算可配置的分位数。

基本度量标准名称`<basename>`的`summary`在scrape期间公开了多个时间序列：
 - 流φ-quantiles (0 ≤ φ ≤ 1), 显示为`<basename>{quantiles="[φ]"}`
 - `<basename>_sum`， 是指所有观察值的总和
 - `<basename>_count`, 是指已观察到的事件计数值

有关φ-分位数，Summary用法和histogram图差异的详细说明，详见[histogram和summaries](https://prometheus.io/docs/practices/histograms)

有关`summaries`的客户端使用文档：

 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Summary)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Summary.java)
 - [Python](https://github.com/prometheus/client_python#summary)
 - [Ruby](https://github.com/prometheus/client_ruby#summary)

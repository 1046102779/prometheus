Prometheus客户端库提供了四个核心的`metrics`类型。这四种类型目前仅在客户端库(启用针对特定类型使用量身定制的API)和`wire`协议中区分。Prometheus服务还没有充分利用这些类型。不久的将来就会发生改变。

##### 一、Counter
*counter* 是表示单个[单调递增计数器](https://en.wikipedia.org/wiki/Monotonic_function)的累积度量，其值只能在重启时增加或重置为零。 例如，您可以使用`counter`来表示所服务的请求数，已完成的任务或错误。

不要使用`counter`来暴露可能减少的值。例如，不要使用`counter`来处理当前正在运行的进程数; 而是使用`gauge`。

客户端使用`counter`的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Counter.java)
 - [Python](https://github.com/prometheus/client_python#counter)
 - [Ruby](https://github.com/prometheus/client_ruby#counter)

##### 二、Gauge
*gauge*是一个度量指标，它表示一个既可以递增, 又可以递减的值。

`Gauge`通常用于测量值，如温度或当前内存使用情况，但也可用于可以上下的"计数"，例如并发请求的数量。

客户端使用`gauge`的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Gauge.java)
 - [Python](https://github.com/prometheus/client_python#gauge)
 - [Ruby](https://github.com/prometheus/client_ruby#gauge)

##### 三、Histogram(柱状图)
*histogram*，对观察结果进行采样（通常是请求持续时间或响应大小等），并将其计入可配置存储桶中。它还提供所有观察值的总和。

基本度量标准名称为`<basename>`的直方图在scrape期间显示多个时间序列：

 - 暴露的观察桶的累积计数器：`<basename>_bucket{le="<upper inclusive bound>"}`
 - 所有观测值的总和：`<basename>_sum`
 - 已观察到的事件数：`<basename>_count`，和`<basename>_bucket{le="+Inf"}`相同

使用[`histogram_quantile()`](https://prometheus.io/docs/querying/functions/#histogram_quantile), 计算直方图或者是直方图聚合计算的分位数阈值。 `histogram`适合计算计算[Apdex值](http://en.wikipedia.org/wiki/Apdex), 当在`buckets`上操作时，记住`histogram`是累计的。详见[ histograms和summaries](https://prometheus.io/docs/practices/histograms)

客户库的`histogram`使用文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Histogram.java)
 - [Python](https://github.com/prometheus/client_python#histogram)
 - [Ruby](https://github.com/prometheus/client_ruby#histogram)

##### 四、Summary
类似*histogram*，*summary*是采样点分位图统计(通常是请求持续时间和响应大小等)。虽然它还提供观察的总数和所有观测值的总和，但它在滑动时间窗口上计算可配置的分位数。

基本度量标准名称`<basename>`的`summary`在scrape期间公开了多个时间序列：
 - 流φ-quantiles (0 ≤ φ ≤ 1), 显示为`<basename>{quantiles="[φ]"}`
 - `<basename>_sum`， 是指所有观察值的总和
 - `<basename>_count`, 是指已观察到的事件计数值

有关`φ`-分位数，`Summary`用法和`histogram`图差异的详细说明，详见[histogram和summaries](https://prometheus.io/docs/practices/histograms)

有关`summaries`的客户端使用文档：

 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Summary)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Summary.java)
 - [Python](https://github.com/prometheus/client_python#summary)
 - [Ruby](https://github.com/prometheus/client_ruby#summary)

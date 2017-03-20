## metrics类型
---
Prometheus客户库提供了四个核心的metrics类型。这四种类型目前仅在客户库和wire协议中区分。Prometheus服务还没有充分利用这些类型。不久的将来就会发生改变。

### Counter(计数器)
*counter* 是一个累计度量指标，它仅仅是永远递增的数值。计数器主要用于统计服务的请求数、任务完成数和错误出现的次数等等。计数器是一个递增的值。反例：统计goroutines的数量。计数器的使用方式在下面的各个客户端例子中：

客户端使用计数器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Counter.java)
 - [Python](https://github.com/prometheus/client_python#counter)
 - [Ruby](https://github.com/prometheus/client_ruby#counter)

### Gauge(计量器)
*gauge*是一个度量指标，它表示一个既可以递增或者递减的值。

计量器主要用在类似于温度、当前内存使用量等，也可以统计当前服务运行随时增加或者减少的Goroutines数量

客户端使用计量器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Gauge.java)
 - [Python](https://github.com/prometheus/client_python#gauge)
 - [Ruby](https://github.com/prometheus/client_ruby#gauge)

### Histogram[直方图]
*histogram*对观察结果结果(通常是请求持续时间或者响应大小)进行采样，并在可配置的桶中对其进行技术。它还提供所有观察值的总和

带有度量指标名称为`[basename]`的直方图会展示Prometheus服务抓取的时间序列数据
 - [base]_bucket{le="<="}, 是指观察buckets的累计计数器
 - [basename]_sum, 是指观察值总和
 - [basename]_count,是指已经观察到的事件总计数

使用[histogram_quantile()](https://prometheus.io/docs/querying/functions/#histogram_quantile)函数, 计算直方图或者是直方图聚合计算的分位数阈值。 一个直方图计算[Apdex值](http://en.wikipedia.org/wiki/Apdex)也是合适的, 当在buckets上操作时，记住直方图是累计的。详见[直方图和总结](https://prometheus.io/docs/practices/histograms)

客户库的直方图使用文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Histogram.java)
 - [Python](https://github.com/prometheus/client_python#histogram)
 - [Ruby](https://github.com/prometheus/client_ruby#histogram)

### [Summary]总结
类似*histogram*，*summary*观察样本值(使用场景类似：请求持续时间和响应大小)。它也提供观察的总计数和所有观察值的总和。同时它可以在滑动时间窗口上计算可配置的分位数。

带有度量指标的`[basename]`的`summary` 在抓取时间序列数据展示。
 - 观察时间的φ-quantiles (0 ≤ φ ≤ 1), 显示为`[basename]{分位数="[φ]"}`
 - `[basename]_sum`， 是指所有观察值的总和
 - `[basename]_count`, 是指已观察到的事件计数值

详见[histogram和summaries](https://prometheus.io/docs/practices/histograms)

有关`summaries`的客户端使用文档：

 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Summary)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Summary.java)
 - [Python](https://github.com/prometheus/client_python#summary)
 - [Ruby](https://github.com/prometheus/client_ruby#summary)

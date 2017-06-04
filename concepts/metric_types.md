## metrics类型
---
Prometheus客户库提供了四个核心的metrics类型。这四种类型目前仅在客户库和wire协议中区分。Prometheus服务还没有充分利用这些类型。不久的将来就会发生改变。

### Counter(计数器)
*counter* 是一个累计度量指标，它是一个只能递增的数值。计数器主要用于统计服务的请求数、任务完成数和错误出现的次数等等。计数器是一个递增的值。反例：统计goroutines的数量。计数器的使用方式在下面的各个客户端例子中：

客户端使用计数器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Counter.java)
 - [Python](https://github.com/prometheus/client_python#counter)
 - [Ruby](https://github.com/prometheus/client_ruby#counter)

### Gauge(计量器)
*gauge*是一个度量指标，它表示一个既可以递增, 又可以递减的值。

计量器主要用在类似于温度、当前内存使用量等，也可以统计当前服务运行随时增加或者减少的Goroutines数量

客户端使用计量器的文档：
 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Gauge.java)
 - [Python](https://github.com/prometheus/client_python#gauge)
 - [Ruby](https://github.com/prometheus/client_ruby#gauge)

### Histogram
*histogram*对观察结果(通常是请求持续时间或者响应大小)进行采样，并在可配置的桶中对其进行统计。它还提供所有观察值的总和

带有度量指标名称为`[basename]`的histogram会展示Prometheus服务抓取的时间序列数据
 - [basename]_bucket{le="<upper inclusive bound"}, 是指观察buckets的累计计数器
 - [basename]_sum, 是指观察值总和
 - [basename]_count,是指已经观察到的事件总计数

`histogram理解：`对每个度量指标进行histogram统计，会生成三个度量指标数据，分别是<basename>_bucket, <basename>_sum, 和<basename>_count三个度量指标，对于`<basename>_bucket`: 会有三个数据输入，一个是基准值，一个是每次增长的步长，一个是横坐标的长度。 
例子：统计`pond_temperature_celsius`的histogram, 输入<20, 5, 5>, 采样两千次:30 + math.Floor(120*math.Sin(float64(i)*0.1))/10), 结果：
```
histogram: <
  sample_count: 2000
  sample_sum: 59968.60000000001
  bucket: <
    cumulative_count: 383
    upper_bound: 20
  >
  bucket: <
    cumulative_count: 729
    upper_bound: 25
  >
  bucket: <
    cumulative_count: 1002
    upper_bound: 30
  >
  bucket: <
    cumulative_count: 1278
    upper_bound: 35
  >
  bucket: <
    cumulative_count: 1632
    upper_bound: 40
  >
>
```
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

`summaries理解`:  分位数分别是0.5,  0.9,  0.99
 ```
summary: <
  sample_count: 1000
  sample_sum: 29969.50000000001
  quantile: <
    quantile: 0.5
    value: 31.1
  >
  quantile: <
    quantile: 0.9
    value: 41.3
  >
  quantile: <
    quantile: 0.99
    value: 41.9
  >
>
 ```

详见[histogram和summaries](https://prometheus.io/docs/practices/histograms)

有关`summaries`的客户端使用文档：

 - [Go](http://godoc.org/github.com/prometheus/client_golang/prometheus#Summary)
 - [Java](https://github.com/prometheus/client_java/blob/master/simpleclient/src/main/java/io/prometheus/client/Summary.java)
 - [Python](https://github.com/prometheus/client_python#summary)
 - [Ruby](https://github.com/prometheus/client_ruby#summary)

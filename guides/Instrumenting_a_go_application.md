## INSTRUMENTING A GO APPLICATION FOR PROMETHEUS
---
Prometheus有一个官方[Go客户端库](https://github.com/prometheus/client_golang)，您可以使用它来检测Go应用程序。 在本指南中，我们将创建一个简单的Go应用程序，通过HTTP公开Prometheus指标。

### 安装
您可以使用go get安装指南所需的`prometheus`，`promauto`和`promhttp`库：
```
go get github.com/prometheus/client_golang/prometheus
go get github.com/prometheus/client_golang/prometheus/promauto
go get github.com/prometheus/client_golang/prometheus/promhttp
```

### Go exposition如何运作
要在Go应用程序中公开Prometheus指标，您需要提供`/metrics` HTTP端点。 您可以使用`prometheus/promhttp`库的HTTP` Handler`作为处理函数。

例如，这个最小的应用程序将通过`http://localhost:2112/metrics`公开Go应用程序的默认指标：
```
package main

import (
        "net/http"

        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}
```
启动这个应用：
```
go run main.go
```
获取metrics：
```
curl http://localhost:2112/metrics
```

### 添加你自己的指标
上面的应用程序仅公开默认的Go指标。 您还可以注册自己的自定义应用程序特定度量标准。 此示例应用程序公开`myapp_processed_ops_total`计数器，该计数器计算到目前为止已处理的操作数。 每2秒钟，计数器加1。
```
package main

import (
        "net/http"
        "time"

        "github.com/prometheus/client_golang/prometheus"
        "github.com/prometheus/client_golang/prometheus/promauto"
        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func recordMetrics() {
        go func() {
                for {
                        opsProcessed.Inc()
                        time.Sleep(2 * time.Second)
                }
        }()
}

var (
        opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
                Name: "myapp_processed_ops_total",
                Help: "The total number of processed events",
        })
)

func main() {
        recordMetrics()

        http.Handle("/metrics", promhttp.Handler())
        http.ListenAndServe(":2112", nil)
}
```
启动这个应用：
```
go run main.go
```
获取metrics：
```
curl http://localhost:2112/metrics
```
在metrics输出中，您将看到`myapp_processed_ops_total`计数器的帮助文本，类型信息和当前值：
```
# HELP myapp_processed_ops_total The total number of processed events
# TYPE myapp_processed_ops_total counter
myapp_processed_ops_total 5
```
您可以[配置](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cscrape_config)本地运行的Prometheus实例以从应用程序中删除指标。 这是`prometheus.yml`配置示例：
```
scrape_configs:
- job_name: myapp
  scrape_interval: 10s
  static_configs:
  - targets:
    - localhost:2112
```

### 其他Go客户端特征
在本指南中，我们仅介绍了Prometheus Go客户端库中的一小部分功能。 您还可以公开其他度量标准类型，例如[gauges](https://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)和[直方图](https://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram)，[非全局注册表](https://godoc.org/github.com/prometheus/client_golang/prometheus#Registry)，将[指标推送](https://godoc.org/github.com/prometheus/client_golang/prometheus/push)到Prometheus [PushGateways](https://prometheus.io/docs/instrumenting/pushing/)的功能，桥接Prometheus和[Graphite](https://godoc.org/github.com/prometheus/client_golang/prometheus/graphite)等等。

### 总结
在本指南中，您创建了两个示例Go应用程序，这些应用程序向Prometheus公开指标 - 一个仅公开默认Go指标，另一个公开自定义Prometheus计数器 - 并配置Prometheus实例以从这些应用程序中提取指标。
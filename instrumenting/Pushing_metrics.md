 Pushgateway](https://github.com/prometheus/pushgateway)允许您将时间序列从[短期服务级批处理作业](https://prometheus.io/docs/practices/pushing/)推送到Prometheus可以抓取的中间作业。 结合Prometheus简单的基于文本的展示格式，这使得即使没有客户端库的shell脚本也很容易。

有关使用Pushgateway和从Unix shell使用的更多信息，请参阅项目的[README.md](https://github.com/prometheus/pushgateway/blob/master/README.md)。

- 要从Java中使用，请参阅[PushGateway](https://prometheus.github.io/client_java/io/prometheus/client/exporter/PushGateway.html)类。

- 要在Go中使用，请参阅[Push](https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Push)和[Add](https://godoc.org/github.com/prometheus/client_golang/prometheus/push#Pusher.Add)方法。

- 要在Python中使用，请参阅导出到[Pushgateway](https://github.com/prometheus/client_python#exporting-to-a-pushgateway)。

- 要从Ruby使用，请参阅[Pushgateway文档](https://github.com/prometheus/client_ruby#pushgateway)。


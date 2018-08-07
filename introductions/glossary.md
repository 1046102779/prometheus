## 词汇表
### Alert(警告)
警告是Prometheus服务正在激活警报规则的结果。警报将数据从Prometheus服务发送到警告管理器
### (Alertmanager)警告管理器
警告管理器接收警告，并把它们聚合成组、去重复数据、应用静默和节流，然后发送通知到邮件、Pageduty或者Slack等系统中
### (Bridge)网桥
网桥是一个从客户端库提取样本，然后将其暴露给非Prometheus监控系统的组件。例如：Python客户端可以将度量指标数据导出到Graphite。
### (Client library)客户库
客户库是使用某种语言（Go、Java、Python、Ruby等），可以轻松直接调试代码，编写样本收集器去拉取来自其他系统的数据，并将这些度量指标数据输送给Prometheus服务。
### (Collector) 收集器
收集器是表示一组度量指标导出器的一部分。它可以是单个度量指标，也可以是从另一个系统提取的多维度度量指标。
### (Direct instrumentation)直接测量
直接测量是将测量在线添加到程序的代码中
### (Exporter)导出器
导出器是暴露Prometheus度量指标的二进制文件，通常将非Prometheus数据格式转化为Prometheus支持的数据处理格式
### (Notification)通知
通知表示一组或者多组的警告，通过警告管理器将通知发送到邮件，Pagerduty或者Slack等系统中
### (PromDash) 面板
[PromDash](https://prometheus.io/docs/visualization/promdash/)是Prometheus的Ruby-on-rails主控面板构建器。它和Grafana有高度的相似之处，但是它只能为Prometheus服务
### Prometheus
Prometheus经常称作Prometheus系统的核心二进制文件。它也可以作为一个整体，被称作Prometheus监控系统
### (PromQL) Prometheus查询语言
[PromQL](https://prometheus.io/docs/querying/basics/)是Prometheus查询语言。它支持聚合、分片、切割、预测和连接操作
### Pushgateway
Pushgateway会保留最近从批处理作业中推送的度量指标。这允许服务中断后Prometheus能够抓取它们的度量指标数据
### Silence
在AlertManager中的静默可以阻止符合标签的警告通知
### Target
在Prometheus服务中，一个应用程序、服务、端点的度量指标数据

##### 一、(Alert)警告
警告是Prometheus服务正在激活警报规则的结果。警报从Prometheus服务发送到警告管理器。

##### 二、(Alertmanager)警告管理器
警告管理器接收警告，并把它们聚合成组、去重复数据、应用静默和节流，然后发送通知到邮件、Pageduty或者Slack等系统中。

##### 三、(Bridge)网桥
网桥是一个从客户端库提取样本，然后将其暴露给非Prometheus监控系统的组件。例如：Python、Go和Java客户端可以将度量指标数据导出到Graphite。

##### 四、(Client library)客户库
客户端库是使用某种语言（Go、Java、Python、Ruby等），可以轻松直接调试代码，编写样本收集器去拉取来自其他系统的数据，并将这些度量指标数据输送给Prometheus服务。

##### 五、(Collector) 收集器
收集器是表示一组度量指标导出器的一部分。它可以是单个度量指标，也可以是从另一个系统提取的多维度度量指标。

##### 六、(Direct instrumentation)直接测量
直接检测是使用[客户端库](https://prometheus.io/docs/introduction/glossary/#client-library)内联作为程序源代码的一部分内联添加的检测。

##### 七、(Endpoint)端点
可以抓取的度量标准源，通常对应于单个进程。

##### 八、(Exporter)导出器
导出器是暴露Prometheus度量指标的二进制文件，通常将非Prometheus数据格式转化为Prometheus支持的数据处理格式

##### 九、(Instance)实例
实例是唯一标识作业中目标的标签。

##### 十、(Job)作业
具有相同目的的目标集合（例如，监视为可伸缩性或可靠性而复制的一组类似进程）被称为作业。

##### 十一、(Notification)通知
通知表示一组或者多组的警告，通过警告管理器将通知发送到邮件，Pagerduty或者Slack等系统中

##### 十二、(PromDash) 面板
[PromDash](https://prometheus.io/docs/visualization/promdash/)是Prometheus的Ruby-on-rails主控面板构建器。它和Grafana有高度的相似之处，但是它只能为Prometheus服务

##### 十三、Prometheus
Prometheus经常称作Prometheus系统的核心二进制文件。它也可以作为一个整体，被称作Prometheus监控系统

##### 十四、(PromQL) Prometheus查询语言
[PromQL](https://prometheus.io/docs/querying/basics/)是Prometheus查询语言。它支持聚合、分片、切割、预测和连接操作

##### 十五、Pushgateway
[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/)会保留最近从批处理作业中推送的度量指标。这允许服务中断后Prometheus能够抓取它们的度量指标数据

##### 十六、(Remote read)远程读取
远程读取是Prometheus功能，允许从其他系统（例如长期存储）透明读取时间序列作为查询的一部分。

##### 十七、(Remote Read Adapter)远程读取适配器
并非所有系统都直接支持远程读取。远程读取适配器位于Prometheus和另一个系统之间，用于转换时间序列请求和它们之间的响应。

##### 十八、(Remote Read Endpoint)远程读取端点
远程读取端点是Prometheus在进行远程读取时所说的。

##### 十九、(Remote Write)远程写入
远程写入是Prometheus功能，允许动态地将采集的样本发送到其他系统，例如长期存储。

##### 二十、(Remote Write Adapter)远程写入适配器
并非所有系统都直接支持远程写入。 远程写入适配器位于Prometheus和另一个系统之间，将远程写入中的样本转换为其他系统可以理解的格式。

##### 二十一、(Remote Write Endpoint)远程写入端点
远程写入端点是Prometheus在进行远程写入时所说的。

##### 二十二、(Sample)样本
样本是时间序列中某个时间点的单个值。
在Prometheus中，每个样本都包含一个float64值和一个毫秒精度的时间戳。

##### 二十三、(Silence)静默
在AlertManager中的静默可以阻止符合标签的警告通知

##### 二十四、(Target)目标
在Prometheus服务中，一个应用程序、服务、端点的度量指标数据

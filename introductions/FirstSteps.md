欢迎来到Prometheus！Prometheus是一个监控平台，通过在监控目标目标上的HTTP端点来收集受监控目标的指标。本指南将向您展示如何使用Prometheus安装，配置和监控我们的第一个资源。 您将下载，安装并运行Prometheus。您还将下载并安装exporter，这些工具可在主机和服务上公开时间序列数据。我们的第一个exporter将是Prometheus本身，它提供了有关内存使用，垃圾收集等的各种主机级指标。

##### 一、下载Prometheus
根据你的平台[https://prometheus.io/download/](https://note.youdao.com/)，然后解压它:
```
tar xvfz prometheus-*.tar.gz
cd prometheus-*
```
Prometheus服务器是一个名为prometheus的二进制文件（或Microsoft Windows上的prometheus.exe）。 我们可以通过传递--help标志来运行二进制文件并查看其选项的帮助。
```
./prometheus --help
usage: prometheus [<flags>]

The Prometheus monitoring server
. . .
```
在使用欢迎来到Prometheus之前，让我们配置它。

##### 二、配置Prometheus
Prometheus配置是YAML。Prometheus下载附带一个名为prometheus.yml的文件中的示例配置，这是一个很好的入门之处。

我们删除了示例文件中的大部分注释，使其更简洁（注释是以＃为前缀的行）。
```
global:
  scrape_interval:     15s
  evaluation_interval: 15s

rule_files:
  # - "first.rules"
  # - "second.rules"

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets: ['localhost:9090']
```
示例配置文件中有三个配置块：global，rule_files和scrape_configs。

全局块控制Prometheus服务器的全局配置。 我们有两种选择。 第一个是scrape_interval，它控制Prometheus抓取目标的频率。 您可以为单个目标重写此值。 在这种例子下，全局设置是每15抓取一次。 evaluation_interval选项控制Prometheus评估规则的频率。 Prometheus使用规则创建新的时间序列并生成警报。

rule_files块指定我们希望Prometheus服务器加载的任何规则的位置。 现在我们没有规则。

最后一个块scrape_configs控制Prometheus监视的资源。 由于Prometheus还将自己的数据公开为HTTP端点，因此它可以抓取并监控自身的健康状况。 在默认配置中，有一个名为prometheus的作业，它会抓取Prometheus服务器公开的时间序列数据。 该作业包含一个静态配置的目标，即端口9090上的localhost。Prometheus希望指标在/metrics路径上的目标上可用。 所以这个默认的工作是通过URL抓取：[http//localhost:9090/metrics](http//localhost:9090/metrics)。

返回的时间序列数据将详细说明Prometheus服务器的状态和性能。

有关配置选项的完整规范，请参阅[配置文档](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)。

##### 三、启动Prometheus
要使用我们新创建的配置文件启动Prometheus，请切换到包含Prometheus二进制文件的目录并运行：
```
./prometheus --config.file=prometheus.yml
```
prometheus应该启动。您还应该能够在[http//localhost:9090](http//localhost:9090)浏览到自己的状态页面。给它大约30秒的时间从自己的HTTP指标端点收集有关自己的数据。

您还可以通过导航到其自己的指标端点来验证Prometheus是否正在提供有关自身的指标：[http//localhost:9090/metrics](http//localhost:9090/metrics)。

##### 四、使用表达式浏览器
让我们试着看一下Prometheus收集的关于自己的一些数据。 要使用Prometheus的内置表达式浏览器，请导航到[http//localhost:9090/graph](http//localhost:9090/graph)并在“Graph”选项卡中选择“Console”视图。

正如您可以从[http//localhost：9090/metrics](http//localhost：9090/metrics)收集的那样，Prometheus导出的一个度量标准称为`promhttp_metric_handler_requests_total`（Prometheus服务器已服务的/ metrics请求的总数）。 继续并将其输入表达式控制台：
```
promhttp_metric_handler_requests_total
```
这应该返回许多不同的时间序列（以及为每个记录的最新值），所有时间序列都使用度量标准名称`promhttp_metric_handler_requests_total`，但具有不同的标签。 这些标签指定不同的请求状态。

如果我们只对导致HTTP代码200的请求感兴趣，我们可以使用此查询来检索该信息：
```
promhttp_metric_handler_requests_total{code="200"}
```
要计算返回的时间序列总数，您可以写：
```
count(promhttp_metric_handler_requests_total)
```
有关表达式语言的更多信息，请参阅[表达式语言文档](https://prometheus.io/docs/prometheus/latest/querying/basics/)。

##### 五、适用图表接口
要绘制表达式图表，请导航到[http//localhost:9090/graph](http//localhost:9090/mgraph) graph并使用“图表”选项卡。

例如，输入以下表达式来绘制在自我抓取的Prometheus中发生的返回状态代码200的每秒HTTP请求率：
```
rate(promhttp_metric_handler_requests_total{code="200"}[1m])
```
您可以尝试图形范围参数和其他设置。

##### 六、监控其他目标
仅从Prometheus那里收集指标并不能很好地反映Prometheus的能力。 为了更好地了解Prometheus可以做什么，我们建议您浏览有关其他exporter的文档。 使用node exporter指南[监控Linux或macOS主机指标](https://prometheus.io/docs/guides/node-exporter/)是一个很好的起点。

##### 七、总结
在本指南中，您安装了Prometheus，配置了Prometheus实例来监视资源，并学习了在Prometheus表达式浏览器中处理时间序列数据的一些基础知识。 要继续了解Prometheus，请查看[概述](https://prometheus.io/docs/introduction/overview/)，了解接下来要探索的内容。


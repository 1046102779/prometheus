##### 一、一般问题
###### 1. Prometheus是什么？
> Prometheus是一款高活跃生态系统的开源系统监控和预警工具包。详见[概览](https://prometheus.io/docs/introduction/overview/)

###### 2. Prometheus与其他的监控系统比较
> 详见[比较](https://prometheus.io/docs/introduction/comparison/)

###### 3. Prometheus有什么依赖？
> Prometheus服务独立运行，没有其他依赖

###### 4. Prometheus有高可用的保证吗？
> 是的，在多台服务器上运行相同的Prometheus服务，相同的预警会由[预警管理器](https://github.com/prometheus/alertmanager)删除
> 为了[提高Alertmanager的可用性](https://github.com/prometheus/alertmanager#high-availability)，您可以在[Mesh群集](https://github.com/weaveworks/mesh)中运行多个实例，并将Prometheus服务器配置为向每个实例发送通知。

###### 5. 我被告知Prometheus"不能水平扩展"
> 事实上，有许多方式可以扩展Prometheus。 阅读Robust Percetion的博客关于Prometheus的[扩展](https://www.robustperception.io/scaling-and-federating-prometheus/)

###### 6. Prometheus是什么语言写的？
> 大多数Prometheus组件是由Go语言写的。还有一些是由Java，Python和Ruby写的

###### 7. Prometheus的特性、存储格式和APIs有多稳定？
> Prometheus从v1.0.0版本开始就非常稳定了，我们现在有一些版本功能规划,详见[路线图](https://prometheus.io/docs/introduction/roadmap/)。重大更改以主要版本的增量表示。 实验组件可能会出现例外，声明中会明确标出此类例外。
> 通常，即使尚未达到1.0.0版的存储库也相当稳定。 我们的目标是为每个存储库制定适当的发布流程并最终发布1.0.0。 在任何情况下，重大更改都将在发行说明中指出（由[`CHANGE`]标记），或者对于尚未正式发行的组件进行清楚地传达。

###### 8. 为什么是使用的是pull而不是push？
基于Http方式的`pull`模型提供了一下优点：
 - 当开发变化时，你可以在笔记本上运行你的监控
 - 如果目标实例挂掉，你可以很容易地知道
 - 你可以手动指定一个目标，并通过浏览器检查该目标实例的监控状况

总体而言，我们认为`pull`比`push`略好，但在考虑使用监视系统时，不应将其视为重点。

如果你必须要用`Push`模式，我们提供[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/)

###### 9. 怎么样把日志推送到Prometheus系统中？
> 简单地回答：千万别这样做，你可以使用[ELK栈](https://www.elastic.co/cn/products/)去实现
> 比较详细的回答：Prometheus是一款收集和处理度量指标的系统，并非事件日志系统。Raintank的博客有关[日志、度量指标和图表](https://blog.raintank.io/logs-and-metrics-and-graphs-oh-my/)在日志和度量指>
标之间，进行了详尽地阐述。

> 如果你想要从应用日志中提取Prometheus度量指标中。 谷歌的[mtail](https://github.com/google/mtail)可能会更有帮助

###### 10. 谁写的Prometheus？
> Prometheus项目发起人是Matt T. Proud和Julius Volz。 一开始大部分的开发是由[SoundCloud](https://soundcloud.com/)赞助的
> 现在它由许多公司和个人维护和扩展

###### 11. 当前Prometheus的许可证是用的哪个？
> Prometheus是根据[Apache 2.0](https://github.com/prometheus/prometheus/blob/master/LICENSE)许可发布的

###### 12. Prometheus单词的复数是什么？
> 经过[广泛研究](https://www.youtube.com/watch?v=B_CDeYrqxjQ&feature=youtu.be)，已确定`Prometheus`的正确复数是`Prometheis`。

###### 13. 我能够动态地加载Prometheus的配置吗？
> 是的，将`SIGHUP`发送到Prometheus进程或将HTTP POST请求发送到`/-/reload`端点将重新加载并应用配置文件。 各种组件尝试妥善处理失败的更改。

###### 14. 我能发送告警吗？
> 是的，通过[预警管理器](https://github.com/prometheus/alertmanager)
当前，下面列表的外部系统都是被支持的
 - Email
 - General Webhooks
 - HipChat(https://www.hipchat.com/)
 - OpsGenie(https://www.atlassian.com/software/opsgenie)
 - PagerDuty(http://www.pagerduty.com/)
 - Pushover(https://pushover.net/)
 - Slack(https://slack.com/)
 
###### 15. 我能创建Dashboard吗？
> 是的，但是在生产使用中，我们推荐用[Grafana](https://prometheus.io/docs/visualization/grafana/)。[PromDash](https://prometheus.io/docs/visualization/promdash/)和[Console templates](https://prom
etheus.io/docs/visualization/consoles/)也可以

###### 16. 我能改变timezone和UTC吗？
> 为避免任何时区混乱，特别是在涉及所谓的夏时制时，我们决定在Prometheus的所有组件中内部专门使用Unix time和UTC进行显示。 可以将精心选择的时区选择引入UI。 欢迎捐款。 有关此工作的当前状态，请参阅[issue#500](https://github.com/prometheus/prometheus/issues/500)。

##### 二、仪表

###### 1. 哪些语言有工具库？
> 这里有很多客户端库，用Prometheus的度量指标度量你的服务。详见[client库](https://prometheus.io/docs/instrumenting/clientlibs/)
> 如果你对功能工具库非常感兴趣，详见[exposition formats](https://prometheus.io/docs/instrumenting/exposition_formats/)

###### 2. 我能监控机器吗？
> 是的。[Node Exporter](https://github.com/prometheus/node_exporter)暴露了很多机器度量指标，包括CPU使用率、内存使用率和磁盘利用率、文件系统的余量和网络带宽等数据。

###### 3. 我能监控网络数据吗？
> 是的。[SNMP Exporter](https://github.com/prometheus/snmp_exporter)允许监控网络设备。

###### 4. 我能监控批量任务吗？
> 是的，通过[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/). 详见[最佳实践](https://prometheus.io/docs/practices/instrumentation/#batch-jobs)

###### 5. Prometheus开箱即用的监控应用程序是什么？
> 详见[exporters和integrations列表](https://prometheus.io/docs/instrumenting/exporters/)

###### 6. 我能通过JMX监控JVM应用程序吗？
> 是的。不能直接使用Java客户端进行测试的应用程序，你可以将[JMX Exporter](https://github.com/prometheus/jmx_exporter)单独使用或者Java代理使用

###### 7. 工具对性能的影响是什么？
> 客户端库和语言之间的性能可能会有所不同。 对于`Java`，[基准测试](https://github.com/prometheus/client_java/blob/master/benchmark/README.md)表明，根据争用，使用`Java`客户端增加计数器/表将花费12-17ns。 除了对延迟最关键的代码之外，所有其他代码都可以忽略不计。

##### 三、故障排除

###### 1. 我的Prometheus 1.x服务器需要很长时间才能启动并使用有关崩溃恢复的大量信息来保存日志。。
> 你的服务可能遭到了不干净的关闭。Prometheus必须在SIGTERM后彻底关闭，特别地对于一些重量级服务可能需要比较长的时间去。如果服务器崩溃或者强制杀死（如：在等待Prometheus关闭时，内核的OOM杀死你的Promethe
us服务），必须执行崩溃恢复，这在正常情况下需要不到一分钟。详见[崩溃恢复](https://prometheus.io/docs/operating/storage/#crash-recovery)

###### 2. 我的Prometheus 1.x服务器内存不足。
> 请参阅有关[内存使用](https://prometheus.io/docs/prometheus/1.8/storage/#memory-usage)情况的部分，以配置Prometheus可用的内存量。

###### 3. 我的Prometheus 1.x服务器报告处于“匆忙模式”或“存储需要限制”。
> 您的存储空间很重。阅读有关[配置本地存储](https://prometheus.io/docs/prometheus/1.8/storage/)的部分，了解如何调整设置以获得更好的性能。

##### 四、实现

###### 1. 为什么所有样品值都是float64数据类型？我想要integer数据类型。
> 我们限制了float64以简化设计[,IEEE 754双精度二进制浮点格式](https://en.wikipedia.org/wiki/Double-precision_floating-point_format)支持高达253的值的整数精度。如果您需要高于253但低于263的整数精度，支持本地64位整数将有帮助。原则上，支持不同的样本值类型 （包括某种大整数
，支持甚至超过64位）可以实现，但它现在不是一个优先级。 注意，一个计数器，即使每秒增加100万次，只有在超过285年后才会出现精度问题。

###### 2. 为什么Prometheus服务器组件不支持TLS或身份验证？ 我可以添加这些吗？
> 注意：Prometheus团队在2018年8月11日的开发峰会期间已经改变了对此的立场，现在正在项目的[路线图](https://prometheus.io/docs/introduction/roadmap/#tls-and-authentication-in-http-serving-endpoints)中支持TLS和服务端点的身份验证。 代码更改后，将更新此文档。

> 虽然TLS和身份验证是经常被请求的功能，但我们故意没有在Prometheus的任何服务器端组件中实现它们。 我们已经决定专注于构建最佳监控系统，而不是在每个服务器组件中支持完全通用的TLS和身份验证解决方案，因此有两个不同的选项和参数（仅TLS的10多个选项）。

> 如果您需要TLS或身份验证，我们建议将反向代理放在Prometheus前面。 参见例如[使用Nginx添加对Prometheus的基本认证](https://www.robustperception.io/adding-basic-auth-to-prometheus-with-nginx)。

> 这仅适用于入站连接。 Prometheus确实支持抓取[TLS-和auth启用的目标](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cscrape_config%3E)，以及其他创建出站连接的Prometheus组件具有类似的支持。

##### 一、一般问题
###### 1. Prometheus是什么？
> Prometheus是一款高活跃生态系统的开源系统监控和警告工具包。详见[概览](https://prometheus.io/docs/introduction/overview/)

###### 2. Prometheus与其他的监控系统比较
> 详见[比较](https://prometheus.io/docs/introduction/comparison/)

###### 3. Prometheus有什么依赖？
> Prometheus服务独立运行，没有其他依赖

###### 4. Prometheus有高可用的保证吗？
> 是的，在多台服务器上运行相同的Prometheus服务，相同的报警会由警告管理器删除
> 警告管理器当前不能保证高可用，但高可用是目标

###### 5. 我被告知Prometheus"不能水平扩展"
> 事实上，有许多方式可以扩展Prometheus。 阅读Robust Percetion的博客关于Prometheus的[扩展](https://www.robustperception.io/scaling-and-federating-prometheus/)

###### 6. Prometheus是什么语言写的？
> 大多数Prometheus组件是由Go语言写的。还有一些是由Java，Python和Ruby写的

###### 7. Prometheus的特性、存储格式和APIs有多稳定？
> Prometheus从v1.0.0版本开始就非常稳定了，我们现在有一些版本功能规划,详见[路线图](https://prometheus.io/docs/introduction/roadmap/)

###### 8. 为什么是使用的是pull而不是push？
基于Http方式的拉模型提供了一下优点：
 - 当开发变化时，你可以在笔记本上运行你的监控
 - 如果目标实例挂掉，你可以很容易地知道
 - 你可以手动指定一个目标，并通过浏览器检查该目标实例的监控状况

总体来说，我们相信拉模式比推模式要好一地啊你，但是当考虑一个监控系统时，它不是主要的考虑点
[Push vs. Pull](http://www.boxever.com/push-vs-pull-for-monitoring)监控在Brian Brazil的博客中被详细的描述

如果你必须要用Push模式，我们提供[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/)

###### 9. 怎么样把日志推送到Prometheus系统中？
> 简单地回答：千万别这样做，你可以使用ELK栈去实现
> 比较详细的回答：Prometheus是一款收集和处理度量指标的系统，并非事件日志系统。Raintank的博客有关[日志、度量指标和图表](https://blog.raintank.io/logs-and-metrics-and-graphs-oh-my/)在日志和度量指>
标之间，进行了详尽地阐述。

> 如果你想要从应用日志中提取Prometheus度量指标中。 谷歌的[mtail](https://github.com/google/mtail)可能会更有帮助

###### 10. 谁写的Prometheus？
> Prometheus项目发起人是Matt T. Proud和Julius Volz。 一开始大部分的开发是由SoundCloud赞助的
> 现在它由许多公司和个人维护和扩展

###### 11. 当前Prometheus的许可证是用的哪个？
> Apache 2.0

###### 12. Prometheus单词的复数是什么？
> Prometheis

###### 13. 我能够动态地加载Prometheus的配置吗？
> 是的，通过发送SIGHUP信号量给Prometheus进行，将会重载配置文件。不同的组件会优雅地处理失败的更改

###### 14. 我能发送告警吗？
> 是的，通过警告管理器
当前，下面列表的外部系统都是被支持的
 - Email
 - General Webhooks
 - PagerDuty(http://www.pagerduty.com/)
 - HipChat(https://www.hipchat.com/)
 - Slack(https://slack.com/)
 - Pushover(https://pushover.net/)
 - Flowdock(https://www.flowdock.com/)
 
###### 15. 我能创建Dashboard吗？
> 是的，但是在生产使用中，我们推荐用[Grafana](https://prometheus.io/docs/visualization/grafana/)。[PromDash](https://prometheus.io/docs/visualization/promdash/)和[Console templates](https://prom
etheus.io/docs/visualization/consoles/)也可以

###### 16. 我能改变timezone和UTC吗？
> 不行。为了避免任何时区的困惑和混乱，我们用了UTC这个通用单位

##### 二、仪表

###### 1. 哪些语言有工具库？
> 这里有很多客户端库，用Prometheus的度量指标度量你的服务。详见[客户库](https://prometheus.io/docs/instrumenting/clientlibs/)
> 如果你对功能工具库非常感兴趣，详见[exposition formats](https://prometheus.io/docs/instrumenting/exposition_formats/)

###### 2. 我能监控机器吗？
> 是的。[Node Exporter](https://github.com/prometheus/node_exporter)暴露了很多机器度量指标，包括CPU使用率、内存使用率和磁盘利用率、文件系统的余量和网络带宽等数据。

###### 3. 我能监控网络数据吗？
> 是的。[SNMP Exporter](https://github.com/prometheus/snmp_exporter)允许监控网络设备。

###### 4. 我能监控批量任务吗？
> 是的，通过[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/). 详见[最佳实践](https://prometheus.io/docs/practices/instrumentation/#batch-jobs)

###### 5. Prometheus开箱即用的监控应用程序是什么？
> 详见[the list of exporters and integrations.](https://prometheus.io/docs/instrumenting/exporters/)

###### 6. 我能通过JMX监控JVM应用程序吗？
> 是的。不能直接使用Java客户端进行测试的应用程序，你可以将[JMX Exporter](https://github.com/prometheus/jmx_exporter)单独使用或者Java代理使用

###### 7. 工具对性能的影响是什么？
> 客户端和语言的性能可能不同。对于Java，基准表明使用Java客户端递增计数器需要12~17ns，具体依赖于竞争。最关键的延迟关键代码之外的所有代码都是可以忽略的。

##### 三、故障排除

###### 1. 我的Prometheus 1.x服务器需要很长时间才能启动并使用有关崩溃恢复的大量信息来保存日志。。
> 你的服务可能遭到了不干净的关闭。Prometheus必须在SIGTERM后彻底关闭，特别地对于一些重量级服务可能需要比较长的时间去。如果服务器崩溃或者司机（如：在等待Prometheus关闭时，内核的OOM杀死你的Promethe
us服务），必须执行崩溃恢复，这在正常情况下需要不到一分钟。详见[崩溃恢复](https://prometheus.io/docs/operating/storage/#crash-recovery)

###### 2. 我的Prometheus 1.x服务器内存不足。
> 请参阅有关[内存使用](https://prometheus.io/docs/prometheus/1.8/storage/#memory-usage)情况的部分，以配置Prometheus可用的内存量。

###### 3. 我的Prometheus 1.x服务器报告处于“匆忙模式”或“存储需要限制”。
> 您的存储空间很重。阅读有关[配置本地存储](https://prometheus.io/docs/prometheus/1.8/storage/)的部分，了解如何调整设置以获得更好的性能。

##### 四、实现

###### 1. 为什么所有样品值都是float64数据类型？我想要integer数据类型。
> 我们限制了float64以简化设计,IEEE 754双精度二进制浮点格式支持高达253的值的整数精度。如果您需要高于253但低于263的整数精度，支持本地64位整数将有帮助。原则上，支持不同的样本值类型 （包括某种大整数
，支持甚至超过64位）可以实现，但它现在不是一个优先级。 注意，一个计数器，即使每秒增加100万次，只有在超过285年后才会出现精度问题。

###### 2. 为什么Prometheus服务器组件不支持TLS或身份验证？ 我可以添加这些吗？
> 注意：Prometheus团队在2018年8月11日的开发峰会期间已经改变了对此的立场，现在正在项目的路线图中支持TLS和服务端点的身份验证。 代码更改后，将更新此文档。

> 虽然TLS和身份验证是经常被请求的功能，但我们故意没有在Prometheus的任何服务器端组件中实现它们。 我们已经决定专注于构建最佳监控系统，而不是在每个服务器组件中支持完全通用的TLS和身份验证解决方案，因此有两个不同的选项和参数（仅TLS的10多个选项）。

> 如果您需要TLS或身份验证，我们建议将反向代理放在Prometheus前面。 参见例如使用Nginx添加对Prometheus的基本认证。

> 这仅适用于入站连接。 Prometheus确实支持删除TLS-和auth启用的目标，以及其他创建出站连接的Prometheus组件具有类似的支持。

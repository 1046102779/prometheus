## 常见问题
### 一般问题
#### 1. Prometheus是什么？
> Prometheus是一款高活跃生态系统的开源系统监控和警告工具包。详见[概览](https://prometheus.io/docs/introduction/overview/)   
#### 2. Prometheus与其他的监控系统比较
> 详见[比较](https://prometheus.io/docs/introduction/comparison/)
#### 3. Prometheus有什么依赖？
> Prometheus服务独立运行，没有其他依赖
#### 4. Prometheus有高可用的保证吗？
> 是的，在多台服务器上运行相同的Prometheus服务，相同的报警会由警告管理器删除
> 警告管理器当前不能保证高可用，但高可用是目标
#### 5. 我被告知Prometheus"不能水平扩展"
> 事实上，有许多方式可以扩展Prometheus。 阅读Robust Percetion的博客关于Prometheus的[扩展](https://www.robustperception.io/scaling-and-federating-prometheus/)
#### 5. Prometheus是什么语言写的？
> 大多数Prometheus组件是由Go语言写的。还有一些是由Java，Python和Ruby写的
#### 6. Prometheus的特性、存储格式和APIs有多稳定？
> Prometheus从v1.0.0版本开始就非常稳定了，我们现在有一些版本功能规划,详见[路线图](https://prometheus.io/docs/introduction/roadmap/)
#### 7. 为什么是使用的是pull而不是push？
基于Http方式的拉模型提供了一下优点：
 - 当开发变化时，你可以在笔记本上运行你的监控
 - 如果目标实例挂掉，你可以很容易地知道
 - 你可以手动指定一个目标，并通过浏览器检查该目标实例的监控状况

总体来说，我们相信拉模式比推模式要好一地啊你，但是当考虑一个监控系统时，它不是主要的考虑点
[Push vs. Pull](http://www.boxever.com/push-vs-pull-for-monitoring)监控在Brian Brazil的博客中被详细的描述

如果你必须要用Push模式，我们提供[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/)
#### 8. 怎么样把日志推送到Prometheus系统中？
> 简单地回答：千万别这样做，你可以使用ELK栈去实现
> 比较详细的回答：Prometheus是一款收集和处理度量指标的系统，并非事件日志系统。Raintank的博客有关[日志、度量指标和图表](https://blog.raintank.io/logs-and-metrics-and-graphs-oh-my/)在日志和度量指标之间，进行了详尽地阐述。

如果你想要从应用日志中提取Prometheus度量指标中。 谷歌的[mtail](https://github.com/google/mtail)可能会更有帮助
#### 9. 谁写的Prometheus？
> Prometheus项目发起人是Matt T. Proud和Julius Volz。 一开始大部分的开发是由SoundCloud赞助的
> 现在它由许多公司和个人维护和扩展
#### 10. 当前Prometheus的许可证是用的哪个？
> Apache 2.0
#### 11. Prometheus单词的复数是什么？
> Prometheis
#### 12. 我能够动态地加载Prometheus的配置吗？
> 是的，通过发送SIGHUP信号量给Prometheus进行，将会重载配置文件。不同的组件会优雅地处理失败的更改
#### 13. 我能发送告警吗？
> 是的，通过警告管理器
当前，下面列表的外部系统都是被支持的
 - Email
 - General Webhooks
 - PagerDuty(http://www.pagerduty.com/)
 - HipChat(https://www.hipchat.com/)
 - Slack(https://slack.com/)
 - Pushover(https://pushover.net/)
 - Flowdock(https://www.flowdock.com/)
#### 14. 我能创建Dashboard吗？
> 是的，但是在生产使用中，我们推荐用[Grafana](https://prometheus.io/docs/visualization/grafana/)。[PromDash](https://prometheus.io/docs/visualization/promdash/)和[Console templates](https://prometheus.io/docs/visualization/consoles/)也可以
#### 15. 我能改变timezone和UTC吗？
> 不行。为了避免任何时区的困惑和混乱，我们用了UTC这个通用单位
### 工具库
#### 1. 哪些语言有工具库？
> 这里有很多客户端库，用Prometheus的度量指标度量你的服务。详见[客户库](https://prometheus.io/docs/instrumenting/clientlibs/)
> 如果你对功能工具库非常感兴趣，详见[exposition formats](https://prometheus.io/docs/instrumenting/exposition_formats/)
#### 2. 我能监控机器吗？
> 是的。[Node Exporter](https://github.com/prometheus/node_exporter)暴露了很多机器度量指标，包括CPU使用率、内存使用率和磁盘利用率、文件系统的余量和网络带宽等数据
#### 3. 我能监控网络数据吗？
> 是的。[SNMP Exporter](https://github.com/prometheus/snmp_exporter)允许监控网络设备
#### 4. 我能监控批量任务吗？
> 是的，通过[Pushgateway](https://prometheus.io/docs/instrumenting/pushing/). 详见[最佳实践](https://prometheus.io/docs/practices/instrumentation/#batch-jobs)
#### 5. Prometheus的第三方工具有哪些？
> 详见[exporters for third-party systems](https://prometheus.io/docs/instrumenting/exporters/)
#### 6. 我能通过JMX监控JVM应用程序吗？
> 是的。不能直接使用Java客户端进行测试的应用程序，你可以将[JMX Exporter](https://github.com/prometheus/jmx_exporter)单独使用或者Java代理使用
#### 7. 工具对性能的影响是什么？
> 客户端和语言的性能可能不同。对于Java，基准表明使用Java客户端递增计数器需要12~17ns，具体依赖于竞争。最关键的延迟关键代码之外的所有代码都是可以忽略的。
### 故障排除
#### 1. 当服务崩溃恢复后，我的服务需要很多时间启动和清理垃圾日志。
> 你的服务可能遭到了不干净的关闭。Prometheus必须在SIGTERM后彻底关闭，特别地对于一些重量级服务可能需要比较长的时间去。如果服务器崩溃或者司机（如：在等待Prometheus关闭时，内核的OOM杀死你的Prometheus服务），必须执行崩溃恢复，这在正常情况下需要不到一分钟。详见[崩溃恢复](https://prometheus.io/docs/operating/storage/#crash-recovery)
#### 2. 我在Linux上使用ZFS，单元测试TestPersistLoadDropChunks失败。尽管测试失败，我运行Prometheus服务，奇怪的事情会发生。
你在Linux上有bug的ZFS文件系统运行Prometheus服务。详见[Issue #484](https://github.com/prometheus/prometheus/issues/484), 在linux v.0.6.4上升级ZFS应该可以解决该问题。

### 实现
#### 1. 为什么所有样品值都是float64数据类型？我想要integer数据类型。
> 我们限制了float64以简化设计,IEEE 754双精度二进制浮点格式支持高达253的值的整数精度。如果您需要高于253但低于263的整数精度，支持本地64位整数将有帮助。原则上，支持不同的样本值类型 （包括某种大整数，支持甚至超过64位）可以实现，但它现在不是一个优先级。 注意，一个计数器，即使每秒增加100万次，只有在超过285年后才会出现精度问题。
#### 2. 为什么Prometheus使用自定义的存储后端，而不是使用其他的存储方法？是不是“一个时间序列一个文件”会大大地伤害性能？
> 一开始，Prometheus是在LevelDB上存储事件序列数据，但不能达到比较好的性能，我们必须改变大量时间序列的存储方式。我们评估了当时可用的许多存储系统，但是没有得到满意的结果。所以我们实现了我们需要的部分。同时保持LevelDB的索引和大量使用文件系统功能。我们最重要的要求是对于常见查询的可接受查询速度，以及每秒数千个样本的可持续速率。后者取决于样本数据的可压缩性和样本所属的时间序列数，但是给你一个想法，这里有一些基准的结果：
 - 在具有Intel Core i7 CPU，8GiB RAM和两个旋转磁盘（三星HD753LJ）的老式8核机器上，Prometheus在每个RAID-1设置中的吞吐速率为34k样本，属于170k时间序列， 600个目标。
 - 在具有64GiB RAM，32个CPU内核和SSD的现代服务器上，Prometheus的每秒吞吐率为525k样本，属于1.4M时间序列，从1650个目标中剔除。

在这两种情况下，没有明显的瓶颈。在相同的流入速度下，各个阶段的处理管道或多或少都会达到他们的限度。

在通常的设置中，不可能使用inode。 有一个可能的缺点：如果你想删除Prometheus的存储目录，你会注意到，一些文件系统在删除文件时非常慢。

#### 3. 为什么Prometheus服务器组件不支持TLS或身份验证？ 我可以添加这些吗？
> 虽然TLS和身份验证是经常请求的功能，但我们有意不在Prometheus的任何服务器端组件中实现它们。 有这么多不同的选项和参数（仅限TLS的10多个选项），我们决定专注于建立最佳的监控系统，而不是在每个服务器组件中支持完全通用的TLS和身份验证解决方案。

如果您需要TLS或身份验证，我们建议将反向代理放在Prometheus前面。 参见例如使用Nginx添加对Prometheus的基本认证。

请注意，这仅适用于入站连接。 Prometheus确实支持删除TLS-和auth启用的目标，以及其他创建出站连接的Prometheus组件具有类似的支持。

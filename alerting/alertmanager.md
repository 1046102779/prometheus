## Alertmanager 警报管理器
---
[Alertmanager](https://github.com/prometheus/alertmanager)处理由客户端应用程序发送的警告，例如：Prometheus服务发送的警告。它关心deduplication, grouping和routing到正确的接收器，如：email，PageDuty或者OpsGenie。它还负责保护和抑制警报。

下面描述有关实现Alertmanager的核心概念。参考[配置文件](https://prometheus.io/docs/alerting/configuration)，它会教你怎样使用Alertmanager。

### Grouping 
Grouping分组将性质类似的警告分组成一个通知类。当许多系统同时出现故障时，这种情况尤其有用，而数百到数千个警报可能同时触发。

例如: 当出现网络分区时，十个到数百个服务实例正在集群中运行。这时多半服务实例暂时无法访问数据库。如果服务实例不能和数据库通信，则对于已经配置好警报规则的Prometheus服务将会对每个服务实例发送一个警报。这样数百个警报会发送到Alertmanager。

如果一个用户仅仅想看到一个页面，这个页面上的数据是精确地表示哪个服务实例受影响了。Alertmanager通过它们的集群和警报名称来分组标签, 这样它可以发送一个单独受影响的通知。

警报分组，分组通知的时间，和通知的接受者是在配置文件中由一个路由树配置的。

### Inhibition 抑制
如果某些其他警报已经触发了，则对于某些警报，Inhibition是一个抑制通知的概念。

例如：一个警报已经触发，它正在通知整个集群是不可达的时，Alertmanager则可以配置成关心这个集群的其他警报无效。这可以防止与实际问题无关的数百或数千个触发警报的通知。

通过Alertmanager的配置文件配置Inhibition。

### Sliences 静默
静默是一个非常简单的方法，可以在给定时间内简单地忽略所有警报。slience基于matchers配置，类似路由树。来到的警告将会被检查，判断它们是否和活跃的slience相等或者正则表达式匹配。如果匹配成功，则不会将这些警报发送给接收者。

Silences在Alertmanager的web接口中配置。

### Client behavior 客户行为
对于客户行为，Alertmanager有[特别要求](https://prometheus.io/docs/alerting/clients)。这些仅仅适用于Prometheus服务不用于发送警报的高级用例。

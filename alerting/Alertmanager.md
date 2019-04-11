Alertmanager处理客户端应用程序（如Prometheus服务器）发送的警报。 它负责对它们进行重复数据删除，分组和路由，以及正确的接收器集成，例如电子邮件，PagerDuty或OpsGenie。 它还负责警报的静音和抑制。

以下描述了Alertmanager实现的核心概念。 请[参阅配置](https://prometheus.io/docs/alerting/configuration/)文档以了解如何更详细地使用它们。

##### 一、分组
分组将类似性质的警报分类为单个通知。 在许多系统一次性失败并且数百到数千个警报可能同时发生的较大中断期间，这尤其有用。

示例：发生网络分区时，群集中正在运行数十或数百个服务实例。 一半的服务实例无法再访问数据库。 Prometheus中的警报规则配置为在每个服务实例无法与数据库通信时发送警报。 结果，数百个警报被发送到Alertmanager。

作为用户，人们只想获得单个页面，同时仍能够确切地看到哪些服务实例受到影响。 因此，可以将Alertmanager配置为按群集和alertname对警报进行分组，以便发送单个紧凑通知。

通过配置文件中的路由树配置警报的分组，分组通知的定时以及这些通知的接收器。

##### 二、抑制
如果某些其他警报已经触发，则抑制是抑制某些警报的通知的概念。

示例：正在触发警报，通知无法访问整个集群。 Alertmanager可以配置为在该特定警报触发时将与该集群有关的所有其他警报静音。 这可以防止发送与实际问题无关的数百或数千个触发警报。

通过Alertmanager的配置文件配置禁止。

##### 三、静默
沉默是在给定时间内简单地静音警报的简单方法。 基于匹配器配置静默，就像路由树一样。 检查传入警报它们是否匹配活动静默的所有相等或正则表达式匹配器。 如果他们这样做，则不会发送该警报的通知。

在Alertmanager的Web界面中配置了静音。

##### 四、客户端行为
Alertmanager对其客户的行为有[特殊要求](https://prometheus.io/docs/alerting/clients/)。 这些仅适用于不使用Prometheus发送警报的高级用例。

##### 五、高可用
Alertmanager支持配置以创建用于高可用性的集群。 这可以使用--[cluster-*](https://github.com/prometheus/alertmanager#high-availability)标志进行配置。

重要的是不要在Prometheus和它的Alertmanagers之间加载平衡流量，而是将Prometheus指向所有Alertmanagers的列表。

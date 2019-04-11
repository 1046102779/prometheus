使用普罗米修斯进行警报分为两部分。 Prometheus服务器中的警报规则会向Alertmanager发送警报。 然后，[Alertmanager](https://prometheus.io/docs/alerting/alertmanager/)管理这些警报，包括静音，禁止，聚合以及通过电子邮件，PagerDuty和HipChat等方法发送通知。

设置警报和通知的主要步骤如下：

- 设置并[配置](https://prometheus.io/docs/alerting/configuration/)Alertmanager
- [配置Prometheus](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Calertmanager_config)与Alertmanager交谈
- 在Prometheus中创建[警报规则](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)

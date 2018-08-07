## 警告概览 alerting overview
---
Pormetheus的警告由独立的两部分组成。Prometheus服务中的警告规则发送警告到Alertmanager。然后这个[Alertmanager](https://prometheus.io/docs/alerting/alertmanager)管理这些警告。包括silencing, inhibition, aggregation，以及通过一些方法发送通知，例如：email，PagerDuty和HipChat。

建立警告和通知的主要步骤：
 - 创建和配置Alertmanager
 - 启动Prometheus服务时，通过`-alertmanager.url`标志配置Alermanager地址，以便Prometheus服务能和Alertmanager建立连接。
 - 在Prometheus服务中创建[警告规则](https://prometheus.io/docs/alerting/rules)

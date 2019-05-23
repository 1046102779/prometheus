## 集成
---
除了客户端库和导出器以及相关库之外，Prometheus还有许多其他通用集成点。 此页面列出了与这些集成的一些集成。

由于功能重叠或仍处于开发阶段，并非所有集成都列在此处。 [导出器默认端口维基页面](https://github.com/prometheus/prometheus/wiki/Default-port-allocations)也恰好包含一些适合这些类别的非导出器集成。

### 文件服务发现
对于Prometheus本身不支持的服务发现机制，[基于文件的服务发现](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cfile_sd_config%3E)提供了集成接口。

- [Docker Swarm](https://github.com/ContainerSolutions/prometheus-swarm-discovery)
- [Scaleway](https://github.com/scaleway/prometheus-scw-sd)

### 远端端点和存储
Prometheus的[远程写入](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cremote_write%3E)和[远程读取](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cremote_read%3E)功能允许透明地发送和接收样本。 这主要用于长期存储。 建议您仔细评估此空间中的任何解决方案，以确认它可以处理您的数据量。

- [AppOptics](https://github.com/solarwinds/prometheus2appoptics):写
- [Chronix](https://github.com/ChronixDB/chronix.ingester):写
- [Cortex](https://github.com/cortexproject/cortex):读和写
- [CrateDB](https://github.com/crate/crate_adapter):读和写
- [Elasticsearch](https://www.elastic.co/guide/en/beats/metricbeat/current/metricbeat-module-prometheus.html):写
- [Gnocchi](https://gnocchi.xyz/prometheus.html):写
- [Graphite]https://github.com/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter):写
- [InfluxDB](https://docs.influxdata.com/influxdb/v1.7/supported_protocols/prometheus):读和写
- [IRONdb](https://github.com/circonus-labs/irondb-prometheus-adapter):读和写
- [Kafka](https://github.com/Telefonica/prometheus-kafka-adapter):写
- [M3DB](https://m3db.github.io/m3/integrations/prometheus/):读和写
- [OpenTSDB](https://github.com/prometheus/prometheus/tree/master/documentation/examples/remote_storage/remote_storage_adapter):写
- [PostgreSQL/TimescaleDB](https://github.com/timescale/prometheus-postgresql-adapter):读和写
- [SignalFx](https://github.com/signalfx/gateway#prometheus):写
- [Splunk](https://github.com/lukemonahan/splunk_modinput_prometheus#prometheus-remote-write):写
- [TiKV](https://github.com/bragfoo/TiPrometheus):读和写
- [VictoriaMetrics](https://github.com/VictoriaMetrics/VictoriaMetrics):写
- [Wavefront](https://github.com/wavefrontHQ/prometheus-storage-adapter):写

### Alertmanager Webhook接收器
对于Alertmanager本身不支持的通知机制，[webhook接收器](https://prometheus.io/docs/alerting/configuration/#webhook_config)允许集成。

- [Alertsnitch](https://gitlab.com/yakshaving.art/alertsnitch):将报警存入MySQL数据库
- [AWS SNS](https://github.com/DataReply/alertmanager-sns-forwarder)
- [DingTalk](https://github.com/timonwong/prometheus-webhook-dingtalk)
- [IRC Bot](https://github.com/multimfi/bot)
- [JIRAlert](https://github.com/free/jiralert)
- [Phabricator / Maniphest](https://github.com/knyar/phalerts)
- [prom2teams](https://github.com/idealista/prom2teams):将通知转发给Microsoft Teams
- [SMS](https://github.com/messagebird/sachet):支持[多个提供商](https://github.com/messagebird/sachet/blob/master/examples/config.yaml)
- [SNMP traps](https://github.com/maxwo/snmp_notifier)
- [Telegram bot](https://github.com/inCaller/prometheus_bot)
- [XMPP Bot](https://github.com/jelmer/prometheus-xmpp-alerts)

### 管理
Prometheus不包含配置管理功能，允许您将其与现有系统集成或构建在其上。

- [Prometheus Operator](https://github.com/coreos/prometheus-operator):在Kubernetes上管理Prometheus
- [Promgen](https://github.com/line/promgen):Prometheus和Alertmanager的Web UI和配置生成器

### 其他
- [karma](https://github.com/prymitive/karma):报警看板
- [PushProx](https://github.com/RobustPerception/PushProx):代理到横向NAT和类似的网络设置
- [Promregator](https://github.com/promregator/promregator):发现和抓取Cloud Foundry应用程序
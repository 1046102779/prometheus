## 导出和集成
---
有很多库和服务能够从帮助第三方系统的度量指标导出为Prometheus的度量指标。这对于（例如：HAProxy或者Linux系统统计信息）不能直接使用Prometheus度量指标是非常有用的。

### 第三方导出器
有一些exporter主要有Prometheus Github组织维护的,详见[地址](https://github.com/prometheus)，这些项目被标记为**official**， 其他都是由外部贡献和维护的。

我们鼓励更多exporters的出现，但无法对所有人进行推广，通常这些exporters被托管在Prometheus Github组织之外。

[JMX exporter](https://github.com/prometheus/jmx_exporter)能够从大多数JVM应用中导出数据，例如：kafka和cassandra。

#### 数据库Databases
 - [Aerospike exporter](https://github.com/alicebob/asprom)
 - [ClickHouse exporter](https://github.com/f1yegor/clickhouse_exporter)
 - [Consul exporter](https://github.com/prometheus/consul_exporter) **官方**
 - [CouchDB exporter](https://github.com/gesellix/couchdb-exporter)
 - [ElasticSearch exporter](https://github.com/justwatchcom/elasticsearch_exporter)
 - [Memcached exporter](https://github.com/prometheus/memcached_exporter) **官方**
 - [MongoDB exporter](https://github.com/dcu/mongodb_exporter)
 - [MySQL server exporter](https://github.com/prometheus/mysqld_exporter) **官方**
 - [PgBouncer exporter](http://git.cbaines.net/prometheus-pgbouncer-exporter/about)
 - [PostgreSQL exporter](https://github.com/wrouesnel/postgres_exporter)
 - [ProxySQL exporter](https://github.com/percona/proxysql_exporter)
 - [Redis exporter](https://github.com/oliver006/redis_exporter)
 - [RethinkDB exporter](https://github.com/oliver006/rethinkdb_exporter)
 - [SQL query result set metrics exporter](https://github.com/chop-dbhi/prometheus-sql)
 - [Tarantool metric library](https://github.com/tarantool/prometheus)

#### 硬件相关Hardware related
 - [apcupsd exporter](https://github.com/mdlayher/apcupsd_exporter)
 - [IoT Edison exporter](https://github.com/lovoo/ipmi_exporter)
 - [knxd exporter](https://github.com/RichiH/knxd_exporter)
 - [Node/System metrics exporter](https://github.com/prometheus/node_exporter) **官方**
 - [Ubiquiti UniFi exporter](https://github.com/mdlayher/unifi_exporter)

#### 消息系统
 - [NATS exporter](https://github.com/lovoo/nats_exporter)
 - [NSQ exporter](https://github.com/lovoo/nsq_exporter)
 - [RabbitMQ exporter](https://github.com/kbudde/rabbitmq_exporter)
 - [RabbitMQ Management Plugin exporter](https://github.com/deadtrickster/prometheus_rabbitmq_exporter)
 - [Mirth Connect exporter](https://github.com/vynca/mirth_exporter)

#### 存储Storage
 - [Ceph exporter](https://github.com/digitalocean/ceph_exporter)
 - [ScaleIO exporter](https://github.com/syepes/sio2prom)
 - [Gluster exporter](https://github.com/ofesseler/gluster_exporter)

#### HTTP
 - [Apache exporter](https://github.com/neezgee/apache_exporter)
 - [HAProxy exporter](https://github.com/prometheus/haproxy_exporter) **官方**
 - [Nginx metric library](https://github.com/knyar/nginx-lua-prometheus)
 - [Nginx VTS exporter](https://github.com/hnlq715/nginx-vts-exporter)
 - [Passenger exporter](https://github.com/stuartnelson3/passenger_exporter)
 - [Varnish exporter](https://github.com/jonnenauha/prometheus_varnish_exporter)
 - [WebDriver exporter](https://github.com/mattbostock/webdriver_exporter)

#### APIs
 - [AWS ECS exporter](https://github.com/slok/ecs-exporter)
 - [Cloudfare exporter](https://github.com/wehkamp/docker-prometheus-cloudflare-exporter)
 - [DigitalOcean exporter](https://github.com/metalmatze/digitalocean_exporter)
 - [Docker Cloud exporter](https://github.com/infinityworksltd/docker-cloud-exporter)
 - [Docker Hub exporter](https://github.com/infinityworksltd/docker-hub-exporter)
 - [Github exporter](https://github.com/infinityworksltd/github-exporter)
 - [Mozilla Observatory exporter](https://github.com/Jimdo/observatory-exporter)
 - [OpenWeatherMap exporter](https://github.com/RichiH/openweathermap_exporter)
 - [Rancher exporter](https://github.com/infinityworksltd/prometheus-rancher-exporter)
 - [Speedtest.net.exporter](https://github.com/RichiH/speedtest_exporter)

#### Logging
 - [Google's mtail log data extractor](https://github.com/google/mtail)
 - [Grok exporter](https://github.com/fstab/grok_exporter)

#### 其他的监控系统
 - [Akamai Colud monitor exporter](https://github.com/ExpressenAB/cloudmonitor_exporter)
 - [AWS CloudWatch exporter](https://github.com/prometheus/cloudwatch_exporter) **官方**
 - [Cloud Foundry Firehose exporter](https://github.com/cloudfoundry-community/firehose_exporter)
 - [Collectd exporter](https://github.com/prometheus/collectd_exporter) **官方**
 - [Graphite exporter](https://github.com/prometheus/graphite_exporter) **官方**
 - [Heka dashboard exporter](https://github.com/docker-infra/heka_exporter)
 - [Heka exporter](https://github.com/imgix/heka_exporter)
 - [InfluxDB exporter](https://github.com/prometheus/influxdb_exporter) **官方**
 - [JMX exporter](https://github.com/prometheus/jmx_exporter) **官方**
 - [Munin exporter](https://github.com/pvdh/munin_exporter)
 - [New Relic exporter](https://github.com/jfindley/newrelic_exporter)
 - [Pingdom exporter](https://github.com/giantswarm/prometheus-pingdom-exporter)
 - [scollector exporter](https://github.com/tgulacsi/prometheus_scollector)
 - [SNMP exporter](https://github.com/prometheus/snmp_exporter) **官方**
 - [StatsD exporter](https://github.com/prometheus/statsd_exporter)

#### 其他杂项
 - [BIG-IP exporter](https://github.com/ExpressenAB/bigip_exporter)
 - [BIND exporter](https://github.com/digitalocean/bind_exporter)
 - [BlackBox exporter](https://github.com/prometheus/blackbox_exporter) **官方**
 - [BOSH exporter](https://github.com/cloudfoundry-community/bosh_exporter)
 - [Dovecot exporter](https://github.com/kumina/dovecot_exporter)
 - [Jenkins exporter](https://github.com/lovoo/jenkins_exporter)
 - [Kemp LoadBalancer exporter](https://github.com/giantswarm/prometheus-kemp-exporter)
 - [Meteor JS web framework exporter](https://atmospherejs.com/sevki/prometheus-exporter)
 - [Minecraft exporter module](https://github.com/Baughn/PrometheusIntegration)
 - [PowerDNS exporter](https://github.com/janeczku/powerdns_exporter)
 - [Process exporter](https://github.com/ncabatoff/process-exporter)
 - [rTorrent exporter](https://github.com/mdlayher/rtorrent_exporter)
 - [Script exporter](https://github.com/adhocteam/script_exporter)
 - [SMTP/Maildir MDA blackbox prober](https://github.com/cherti/mailexporter)
 - [Transmission exporter](https://github.com/metalmatze/transmission-exporter)
 - [Unbound exporter](https://github.com/kumina/unbound_exporter)
 - [Xen exporter](https://github.com/lovoo/xenstats_exporter)

当实现一个新的Prometheus导出器时，请遵循[writing exporter指南]（https://prometheus.io/docs/instrumenting/writing_exporters）,也请参考[邮件列表](https://groups.google.com/forum/#!forum/prometheus-developers)。我们乐意给出建议告诉你怎样写一个尽可能有用地，一致性地导出器

### 可直接使用的软件
一些第三方软件本身就已经暴露了Prometheus度量指标, 因此不需要exporter：
 - [cAdvisor](https://github.com/google/cadvisor)
 - [Doorman](https://github.com/youtube/doorman)
 - [Etcd](https://github.com/coreos/etcd)
 - [Kubernetes-Mesos](https://github.com/mesosphere/kubernetes-mesos)
 - [Kubernetes](https://github.com/kubernetes/kubernetes)
 - [RobustIRC](http://robustirc.net/)
 - [Quobyte](https://www.quobyte.com/)
 - [SkyDNS](https://github.com/skynetservices/skydns)
 - [Weave Flux](http://weaveworks.github.io/flux/)

### 其他第三方实用工具
本节列出了帮助您以特定语言编写代码的库和其他实用程序。它们不是Prometheus客户端库，而是使用客户库。对于所有独立维护的软件，我们无法对其进行优化和改进。
 - Coljure: [prometheus-clj](https://github.com/soundcloud/prometheus-clj)
 - Go: [go-metrics instrumentation library](https://github.com/armon/go-metrics)
 - Go: [gokit](https://github.com/peterbourgon/gokit)
 - Java/JVM: [Hystrix metrics publisher](https://github.com/soundcloud/prometheus-hystrix)
 - Python-Django: [django-prometheus](https://github.com/korfuri/django-prometheus)

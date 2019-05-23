## 使用NODE EXPORTER监控LINUX主机指标
---
Prometheus [Node Exporter](https://github.com/prometheus/node_exporter)公开了各种与硬件和内核相关的指标。

在本指南中，您将：

- 在`localhost`上启动node exporter
- 在`localhost`上启动一个Prometheus实例，该实例配置为从正在运行的node exporter中获取指标

### 安装并运行Node Exporter
Prometheus节点导出器是一个可以[通过tarball](https://prometheus.io/docs/guides/node-exporter/#tarball-installation)安装的单个静态二进制文件。 从Prometheus[下载页面](https://prometheus.io/download/#node_exporter)下载后，将其解压缩并运行：
```
wget https://github.com/prometheus/node_exporter/releases/download/v*/node_exporter-*.*-amd64.tar.gz
tar xvfz node_exporter-*.*-amd64.tar.gz
cd node_exporter-*.*-amd64
./node_exporter
```
您应该看到这样的输出，表明Node Exporter现在正在运行并在端口9100上公开指标：
```
INFO[0000] Starting node_exporter (version=0.16.0, branch=HEAD, revision=d42bd70f4363dced6b77d8fc311ea57b63387e4f)  source="node_exporter.go:82"
INFO[0000] Build context (go=go1.9.6, user=root@a67a9bc13a69, date=20180515-15:53:28)  source="node_exporter.go:83"
INFO[0000] Enabled collectors:                           source="node_exporter.go:90"
INFO[0000]  - boottime                                   source="node_exporter.go:97"
...
INFO[0000] Listening on :9100                            source="node_exporter.go:111"
```

### Node Exporter指标
安装并运行节点导出程序后，您可以通过cURLing `/metrics`端点验证是否正在导出度量标准：
```
curl http://localhost:9100/metrics
```
你应该看到像这样的输出：
```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 3.8996e-05
go_gc_duration_seconds{quantile="0.25"} 4.5926e-05
go_gc_duration_seconds{quantile="0.5"} 5.846e-05
# etc.
```
成功！ Node Exporter现在公开了Prometheus可以抓取的指标，包括输出中更多的系统指标（以`node_`为前缀）。 要查看这些指标（以及帮助和类型信息）：
```
curl http://localhost:9100/metrics | grep "node_"
```

### 配置你的Prometheus实例
需要正确配置本地运行的Prometheus实例才能访问节点导出器指标。 以下`scrape_config`块（在`prometheus.yml`配置文件中）将告诉Prometheus实例通过`localhost：9100`从节点导出器中删除：
```
scrape_configs:
- job_name: 'node'
  static_configs:
  - targets: ['localhost:9100']
```
要安装Prometheus，请下载适用于您平台的最新版本并解压缩它：
```
wget https://github.com/prometheus/prometheus/releases/download/v*/prometheus-*.*-amd64.tar.gz
tar xvf prometheus-*.*-amd64.tar.gz
cd prometheus-*.*
```
安装Prometheus后，您可以启动它，使用`--config.file`标志指向您在上面创建的Prometheus配置：
```
./prometheus --config.file=./prometheus.yml
```

### 通过Prometheus表达式浏览器探索节点导出器指标
现在Prometheus正在从正在运行的Node Exporter实例中抓取指标，您可以使用Prometheus UI（也就是表达式浏览器）来探索这些指标。 在浏览器中导航到`localhost：9090 / graph`，然后使用页面顶部的主表达式栏输入表达式。 表达式栏看起来像这样：

![prometheus-expression-bar](https://prometheus.io/assets/prometheus-expression-bar.png)

特定于节点导出器的度量标准以`node_`为前缀，并包含`node_cpu_seconds_total`和`node_exporter_build_info`等度量标准。

点击下面的链接查看一些示例指标：

| 指标 | 含义 | 
|---|---|
| `rate(node_cpu_seconds_total{mode="system"}[1m])` | 系统模式下每秒平均花费的CPU时间（以秒为单位） | 
| `node_filesystem_avail_bytes` | 非root用户可用的文件系统空间（以字节为单位） |
| `非root用户可用的文件系统空间（以字节为单位）` | 每秒接收的平均网络流量（以字节为单位） | 
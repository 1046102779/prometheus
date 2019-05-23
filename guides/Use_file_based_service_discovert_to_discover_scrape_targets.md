## 使用基于文件的服务发现来抓取目标
---
Prometheus提供各种[服务发现](https://github.com/prometheus/prometheus/tree/master/discovery)选项，用于发现抓取目标，包括[Kubernetes](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Ckubernetes_sd_config)，[Consul](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cconsul_sd_config)等等。 如果您需要使用当前不支持的服务发现系统，Prometheus[基于文件的服务发现机制](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#%3Cfile_sd_config)可以最好地为您的用例提供服务，该机制允许您列出JSON文件中的scrape目标（以及有关这些目标的元数据））。

在本指南中，我们将：

- 在本地安装并运行Prometheus [Node Exporter](https://prometheus.io/docs/guides/node-exporter/)
- 创建一个`targets.json`文件，指定Node Exporter的主机和端口信息
- 安装并运行配置为使用`targets.json`文件发现Node Exporter的Prometheus实例

### 安装并运行Node Exporter
请参阅[使用节点导出器指南监视Linux主机指标](https://prometheus.io/docs/guides/node-exporter/)的此部分。 节点导出程序在端口9100上运行。要确保节点导出程序正在公开指标：
```
curl http://localhost:9100/metrics
```
指标输出应该像这样：
```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
...
```

### 安装，配置并运行Prometheus
与Node Exporter一样，Prometheus是一个可以通过tarball安装的静态二进制文件。 [下载适用于您平台的最新版本](https://prometheus.io/download/#prometheus)并解压缩它：
```
wget https://github.com/prometheus/prometheus/releases/download/v*/prometheus-*.*-amd64.tar.gz
tar xvf prometheus-*.*-amd64.tar.gz
cd prometheus-*.*
```
untarred目录包含`prometheus.yml`配置文件。 用以下内容替换该文件的当前内容：
```
scrape_configs:
- job_name: 'node'
  file_sd_configs:
  - files:
    - 'targets.json'
```
此配置指定存在一个名为`node`（用于节点导出器）的作业，该作业从`targets.json`文件中检索节点导出器实例的主机和端口信息。

现在创建`targets.json`文件并将其添加到其中：
```
[
  {
    "labels": {
      "job": "node"
    },
    "targets": [
      "localhost:9100"
    ]
  }
]
```
此配置指定存在具有一个目标的`node`作业：`localhost:9100`。

现在你可以启动Prometheus了：
```
./prometheus
```
如果Prometheus已成功启动，您应该在日志中看到如下所示的行：
```
level=info ts=2018-08-13T20:39:24.905651509Z caller=main.go:500 msg="Server is ready to receive web requests."
```

### 探索已发现服务的指标
启动并运行Prometheus后，您可以使用Prometheus表达式浏览器探索`node`服务公开的指标。 例如，如果您浏览`up{job="node"}`指标，则可以看到正在正确发现节点导出器。

### 动态更改目标列表
使用Prometheus基于文件的服务发现机制时，Prometheus实例将侦听对文件的更改并自动更新scrape目标列表，而无需重新启动实例。 为了演示这一点，请在端口9200上启动第二个Node Exporter实例。首先导航到包含Node Exporter二进制文件的目录，然后在新的终端窗口中运行此命令：
```
./node_exporter --web.listen-address=":9200"
```
现在通过为新的Node Exporter添加一个条目来修改`targets.json`中的配置：
```
[
  {
    "targets": [
      "localhost:9100"
    ],
    "labels": {
      "job": "node"
    }
  },
  {
    "targets": [
      "localhost:9200"
    ],
    "labels": {
      "job": "node"
    }
  }
]
```
保存更改后，Prometheus将自动收到新目标列表的通知。 `up{jo ="node"}`指标应显示两个`instance`标签为`localhost：9100`和`localhost：9200`。

### 总结
在本指南中，您安装并运行了Prometheus节点导出程序并配置了Prometheus，以使用基于文件的服务发现从节点导出程序中发现和搜索指标。
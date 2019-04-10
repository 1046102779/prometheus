根据我们的[稳定承诺](https://prometheus.io/blog/2016/07/18/prometheus-1-0-released/#fine-print)，Prometheus 2.0版本包含许多向后不兼容的更改。 本文档提供了从Prometheus 1.8迁移到Prometheus 2.0的指导。

##### 一、标志
Prometheus命令行标志的格式已更改。现在所有标志都使用双破折号而不是单个破折号。公共标志（`--config.file`， `--web.listen-address`和`--web.external-url`）仍然相同，但除此之外，几乎所有与存储相关的标志都已被删除。

一些值得注意的标志已删除：

- `-alertmanager.url`在Prometheus 2.0中，已删除用于配置静态Alertmanager URL的命令行标志。现在必须通过服务发现来发现Alertmanager，请参阅Alertmanager服务发现。

- `-log.format`在Prometheus 2.0中，日志只能流式传输到标准错误。

- `-query.staleness-delta`已重命名为`--query.lookback-delta`; Prometheus 2.0引入了一种处理陈旧性的新机制，请参见陈旧性。

- `-storage.local.*` Prometheus 2.0引入了一个新的存储引擎，因此删除了与旧引擎相关的所有标志。有关新引擎的信息，请参阅存储。

- `-storage.remote.*` Prometheus 2.0删除了已经弃用的远程存储标志，如果提供它们将无法启动。要写入InfluxDB，Graphite或OpenTSDB，请使用相关的存储适配器。

##### 二、Alertmanager服务发现
在Prometheus 1.4中引入了Alertmanager服务发现，允许Prometheus使用与刮擦目标相同的机制动态发现Alertmanager复制品。 在Prometheus 2.0中，已删除静态Alertmanager配置的命令行标志，因此以下命令行标志：
> ./prometheus -alertmanager.url=http://alertmanager:9093/

将在`prometheus.yml`配置文件中替换为以下内容：
```
alerting:
  alertmanagers:
  - static_configs:
    - targets:
      - alertmanager:9093
```
您还可以在Alertmanager配置中使用所有常用的Prometheus服务发现集成和重新标记。 此代码段指示Prometheus使用`name: alertmanger`：alertmanager和非空端口在`default`命名空间中搜索Kubernetes pod。
```
alerting:
  alertmanagers:
  - kubernetes_sd_configs:
      - role: pod
    tls_config:
      ca_file: /var/run/secrets/kubernetes.io/serviceaccount/ca.crt
    bearer_token_file: /var/run/secrets/kubernetes.io/serviceaccount/token
    relabel_configs:
    - source_labels: [__meta_kubernetes_pod_label_name]
      regex: alertmanager
      action: keep
    - source_labels: [__meta_kubernetes_namespace]
      regex: default
      action: keep
    - source_labels: [__meta_kubernetes_pod_container_port_number]
      regex:
      action: drop
```
##### 三、记录规则和报警
配置警报和录制规则的格式已更改为YAML。 旧格式的录制规则和警报示例：
```
job:request_duration_seconds:histogram_quantile99 =
  histogram_quantile(0.99, sum(rate(request_duration_seconds_bucket[1m])) by (le, job))

ALERT FrontendRequestLatency
  IF job:request_duration_seconds:histogram_quantile99{job="frontend"} > 0.1
  FOR 5m
  ANNOTATIONS {
    summary = "High frontend request latency",
  }
```
看起来像这样：
```
groups:
- name: example.rules
  rules:
  - record: job:request_duration_seconds:histogram_quantile99
    expr: histogram_quantile(0.99, sum(rate(request_duration_seconds_bucket[1m]))
      BY (le, job))
  - alert: FrontendRequestLatency
    expr: job:request_duration_seconds:histogram_quantile99{job="frontend"} > 0.1
    for: 5m
    annotations:
      summary: High frontend request latency
```
为了帮助进行更改，`promtool`工具具有自动执行规则转换的模式。 给定`.rules`文件，它将以新格式输出`.rules.yml`文件。 例如：
```
$ promtool update rules example.rules
```
请注意，您需要使用2.0中的promtool，而不是1.8。
##### 四、存储
Prometheus 2.0中的数据格式已完全改变，并且不向后兼容1.8。 为了保持对历史监控数据的访问，我们建议您运行至少与Prometheus 2.0实例并行运行至少版本1.8.1的非刮擦Prometheus实例，并让新服务器通过远程读取协议从旧服务器读取现有数据。

您的Prometheus 1.8实例应该使用以下标志和仅包含`external_labels`设置（如果有）的配置文件启动：
> $ ./prometheus-1.8.1.linux-amd64/prometheus -web.listen-address ":9094" -config.file old.yml

然后可以使用以下标志启动Prometheus 2.0（在同一台机器上）：
> $ ./prometheus-2.0.0.linux-amd64/prometheus --config.file prometheus.yml

除了完整的现有配置之外，prometheus.yml还包含哪些节：
```
remote_read:
  - url: "http://localhost:9094/api/v1/read"
```
##### 五、PromQL
从PromQL中删除了以下功能：

- `drop_common_labels`函数 - 应该使用不使用聚合修饰符。
- `keep_common`聚合修饰符 - 应该使用修饰符。
- `count_scalar`函数 - `absent()`或在操作中正确传播标签可以更好地处理用例。
- 
有关详细信息，请参阅[issue＃3060](https://github.com/prometheus/prometheus/issues/3060)。

##### 六、杂项
###### 6.1 普罗米修斯非root用户
Prometheus Docker镜像现在可以作为[非root用户运行Prometheus](https://github.com/prometheus/prometheus/pull/2859)。 如果您希望Prometheus UI / API侦听低端口号（例如，端口80），则需要覆盖它。 对于Kubernetes，您将使用以下YAML：
```
apiVersion: v1
kind: Pod
metadata:
  name: security-context-demo-2
spec:
  securityContext:
    runAsUser: 0
...
```
有关更多详细信息，请参阅[为Pod或容器配置安全上下文](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)。

如果您使用的是Docker，则会使用以下代码段：
```
docker run -u root -p 80:80 prom/prometheus:v2.0.0-rc.2  --web.listen-address :80
```
###### 6.2 普罗米修斯生命周期
如果您使用Prometheus `/-/reload`加载HTTP端点在更改时[自动重新加载Prometheus配置](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)，则出于安全原因，默认情况下会禁用这些端点，这是Prometheus 2.0中的。 要启用它们，请设置`--web.enable-lifecycle`标志。

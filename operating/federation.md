## FEDERATION(联合)
---
Federation允许一个Prometheus服务从另一个Prometheus服务中获取选中的时间序列数据。

### Use cases 用例
对于federation，有几种不同的用例。它常常用于实现Prometheus监控的扩展，或者从另一个Prometheus服务中拉取相关度量指标数据。

#### Hierarchical federation分层扩展
分层扩展允许Prometheus扩展到数十个数据中心和上百万的节点数量。在这个用例中，federation拓扑结构类似一颗树，较高层级的Prometheus服务从大量较低层级的Prometheus服务中检索和聚集时间序列数据。

例如：环境中可能包含许多以数据中心为基础的Prometheus服务，可以从较高层级收集数据，还有一组全局的Prometheus服务，仅仅从本地服务器收集和存储聚合的数据。这提供了一个聚合的全局视图和本地视图。

#### 跨服务的federation 
在跨服务的federation中，一个Prometheus服务配置成从另一个Prometheus服务中获取选中的数据，这个Prometheus服务能够对单个服务中的两个数据集启用警告和查询。

例如，一个运行多个服务的集群调度器可能暴露了集群资源使用信息（例如：CPU和内存使用量）。另一方面，运行在集群上的一个服务仅仅暴露应用程序指定的服务度量指标。经常地，独立的Prometheus服务抓取这两个度量指标CPU和内存。使用federation，这个Prometheus服务包含服务级别的度量指标，这个指标可以从集群Prometheus服务中获取有关其指定的服务集群资源使用量，以便在该服务中使用两组度量指标。

### 配置federation
任何Prometheus服务，这个`/federation`允许从服务中选中的时间序列检索当前值。至少一个`match[]`URL参数必须指定要选择的暴露时间序列。每一个`match[]`参数需要指定一个[即时向量选择器](https://prometheus.io/docs/querying/basics/#instant-vector-selectors)，例如：`up`或者`{job="api-server"}`。如果提供了多个`match`参数，将会选取所有匹配的时间序列数据的并集。

为了一个Prometheus服务从另一个Prometheus服务中federate度量指标,从一个源服务的`/federate`端点，配置你的目标Prometheus服务。当然，也可以使用`honor_labels`获取选项和输入想要的`match[]`的参数。例如，下面`scrape_config`federates任何带有`job="prometheus`标签的所有时间序列，或者一个以`job`开头的度量指标名称：从这个Prometheus服务上端口服务`prometheus-{1,2,3}:9090`上获取其他Prometheus服务度量指标数据。
```
- job_name: 'federate'
  scrape_interval: 15s

  honor_labels: true
  metrics_path: '/federate'

  params:
    'match[]':
      - '{job="prometheus"}'
      - '{__name__=~"job:.*"}'

  static_configs:
    - targets:
      - 'source-prometheus-1:9090'
      - 'source-prometheus-2:9090'
      - 'source-prometheus-3:9090'
``` 

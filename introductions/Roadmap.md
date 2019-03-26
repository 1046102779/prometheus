下面的一些功能使我们即将要做的事情。如果你想查看整个计划和当前工作的完整概述，请查看github上Prometheus项目的issue，如：[Prometheus server](https://github.com/prometheus/prometheus/issues)

##### 一、服务器端指标元数据支持
此时，度量标准类型和其他元数据仅在客户端库和展示中使用，但不会在Prometheus服务器中保留或使用。 我们计划将来使用这个元数据。第一步是在Prometheus中将这些数据聚合在内存中，并通过实验性API端点提供。

##### 二、采用OpenMetrics
OpenMetrics工作组正在为度量标准开发新标准。我们计划在我们的客户端库和Prometheus本身支持这种格式。

##### 三、回填时间序列
回填将允许过去的大量数据。这将允许追溯规则评估，并从其他监控系统传输旧数据。

##### 四、HTTP服务端点中的TLS和身份验证
Prometheus，Alertmanager和官方exporter中的HTTP服务端点尚未内置对TLS和身份验证的支持。 添加此支持将使人们更容易安全地部署Prometheus组件，而无需反向代理从外部添加这些功能。

##### 五、支持生态
Prometheus拥有一系列客户端库和exporter。总是可以支持更多语言，或者从中导出指标有用的系统。 我们将支持生态系统创建和扩展这些系统。

## 路线路
---
下面的一些功能使我们即将要做的事情。如果你想查看整个计划和当前工作的完整概述，请查看github上Prometheus项目的issue，如：[Prometheus服务](https://github.com/prometheus/prometheus/issues)

### 新的存储引擎
现在一个新的，更加高效的存储引擎正在开发中。它将减少资源的使用率，更好的倒排索引，并且可以使Prometheus更好地扩展

### 长期存储
当前Prometheus支持本地存储样本数据，同时也有一些实验性地支持：通过一个通用机制，发送数据到远程系统。例如：TSDBs
我们计划通过Prometheus的通用机制（如：Cortex）添加来自其他TSDB的回读支持。Github issue：[#10](https://github.com/prometheus/prometheus/issues/10)
### 改进陈旧性处理
当前对于一个表达式的查询结果时间超过5分钟后，Prometheus会丢弃结果中的时间序列数据，言外之意是说Prometheus当前只保存5分钟内的查询结果。目前禁止使用Pushgateway和CloudWatch导出时间序列数据，因为它可能表示超过过去5分钟的时间序列，这是不准确的查询结果。如果近期不发生时间序列数据的抓取操作，我们计划仅仅考虑时间序列的无效性。github issue：[#398](https://github.com/prometheus/prometheus/issues/398)
### 服务端度量指标元数据支持
现在度量指标类型和其他元数据仅仅在客户库和展示格式中使用，并不会在Prometheus服务中持久保留或者利用。将来我们计划充分利用这些元数据。第一步是在Prometheus服务的内存中聚合这些数据，并开放一些实验性的API来提供服务
### Prometheus度量指标格式作为一个标准
我们打算提交一个标准化的干净版本给IETF等组织。
### 回填时间序列
回填时间序列的含义是将过去大量的时间序列数据，根据一定的回溯规则，传输到其他的监控系统中。
### 支持生态
Prometheus有大量的客户库和导出数据器。也有大量的语言被支持，或者有一些可以从Prometheus服务中导出的时间序列系统。我们将会为这个生态做更多的创建和扩展，丰富这个生态。

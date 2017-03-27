## 展示格式 exposition formats
---
Prometheus实现了两种不同的wire格式，客户端可以使用它们将度量指标公开到Prometheus服务：简单的基于文本的格式和更有效和更强大的协议缓冲格式。Prometheus服务和客户端使用内容协商来建立使用的实际格式。如果客户端不支持前者，服务器将更愿意接收协议缓冲区格式，并且将返回到基于文本的格式。

大多数用户应该使用已经实现了相关导出格式的客户端库。

### Format v0.0.4
这是当前度量指标展示格式版本

到这个版本，Prometheus有两种替代格式：基于协议缓冲的格式和文本格式。客户端必须支持这两种备用格式中的至少一种。


另外，客户端可选择性地暴露不被Prometheus服务理解的其他文本格式。他们仅仅是为了方便调试。强烈建议客户端库至少一直一种可读的格式。如果客户端库不能立即HTTP Content-Type的头部，那么这种可读的格式应该被回退。v0.0.4的文本格式通常被认为是可读的，所以它是一个很好的候选者（并且也被Prometheus理解）。

#### 格式变体比较
| | 协议缓冲格式 | 文本格式 |
|--|----------:|:--------|
|开始时间| 2014.4月|2014.4月|
|传输协议| HTTP | HTTP|
| 编码| - 32-bit varint-encoded record length-delimited | utf-9, \n行结尾|
|     | - Protocol Buffer messages of type  | |
|     | - io.prometheus.client.MetricFamily | |
| HTTP Content-Type| application/vnd.google.protobuf; proto=io.prometheus.client.MetricFamily; encoding=delimited | text/plain; version=0.0.4 (A missing version value will lead to a fall-back to the most recent text format version.)|
| Optional HTTP Content-Encoding| gzip | gzip|
| 优点 | - 跨平台 | -可读性好
|      |- Size    | - 易于组合，特别适用于简单情况（无需嵌套）|
|      |- 编解码代价小| - 逐行阅读（处理类型提示和文本字符串外）|
|      |- 严格的范式 ||
|      |- 支持链接和理论流（仅服务端行为需要更改）||
| 限制 | - 不具有可读性| - 信息不全 |
|      |          | - 类型和文档格式不是语法的组成部分，意味着很少到不存在的度量契约验证
| | | - 解析代价大|
| 支持度量指标原子性 | - Counter | - Counter |
| | - Gauge | - Gauge|
| | - Histogram |  - Histogram|
| | - Summary | - Summary |
| | - Untyped | - Untyped|
| 兼容性| - 低于v.0.0.3无效 | 无 |
| | - v0.0.4有效 | |

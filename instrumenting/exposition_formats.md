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

### Protocol buffer format details 协议缓冲区格式详细信息
重复展示中协议缓冲去的可重复排序是优选的，但不是必需的，即如果计算成本过高，则不排序

同一个展示内的相同`MetricFamily`的必须有一个唯一的名称。相同`MetricFamily`中的每个度量指标都必须有一组唯一的`LabelPair`字段。否则，这种嵌入行为是未定义的。

#### Text format details 文本格式细节

协议是面向行的。换行字符(\n)分隔行。最后一行必须以换行字符结束。空行被忽略。

在一行内，令牌可以被任何数量的空格和/或制表符分开（如果它们不想和以前的token合并，则必须至少分开一个）。头尾部的空白被忽略。

具有"#"作为第一个费空格字符的行是注释。它们被忽略，除非“#”之后的一个令牌是HELP或TYPE。这些行被视为如下：如果令牌是HELP，则至少需要一个令牌，这是度量指标名称。所有剩余的令牌都被认为是该度量指标名称的docstring。HELP行可以包含任何UTF-8字符序列（度量名称后）。但反斜杠和换行字符必须分别转移为\\和\n。相同的度量名称只能有一个HELP行。

如果令牌是TYPE，则预期只有两个令牌。第一个是度量名称，第二个是计数器，规格，直方图，摘要或无类型，定义该名称的度量标准的类型。相同的度量名称只能有一个TYPE行。用于度量名称的TYPE行必须出现在为该度量名称报告的第一个样本之前。如果度量名称没有TYPE行，则该类型将设置为无类型。剩余行用以下语法描述样本，每行一个(EBNF)：

```
metric_name [
  "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } [ "," ] "}"
] value [ timestamp ]
```

`metric_name`和`label_name`具有普通的Prometheus表达式语言限制。 `label_value`可以是UTF-8字符的任何序列，但反斜杠，双引号和换行字符必须分别转义为\\，\“和\ n，value为浮点数和时间戳 一个int64（从时代以来的毫秒，即1970-01-01 00:00:00 UTC，不包括闰秒），由Go strconv软件包（见函数ParseInt和ParseFloat）表示，特别是Nan，+ Inf和 -Inf是有效值。

给定度量的所有行必须作为一个不间断的组提供，可选的HELP和TYPE行首先（不是特定的顺序）。 除此之外，重复展示的可重复排序是优选的，但不是必需的，即不计算计算成本是否可以排序。

每行必须具有度量名称和标签的唯一组合。 否则，嵌入行为是未定义的。

`histogram`和`summary`类型很难以文本格式表示。 适用以下约定：
 - 名为x的摘要或直方图的样本总和作为名为x_sum的单独样本给出。
 - 名为x的摘要或直方图的样本计数作为名为x_count的单独样本给出。
 - 名为x的摘要的每个分位数作为具有相同名称x和标签{quantile =“y”}的单独样本行给出。
 - 名为x的直方图的每个桶计数作为单独的样本行（名称为x_bucket）和一个标签{le =“y”}（其中y是存储桶的上限）给出。
 - 直方图必须带有{le =“+ Inf”}的存储桶。 其值必须与x_count的值相同。
 - 直方图的桶和总结的分位数必须以其标签值的增加数字顺序（分别为le或分位数标签）出现。
另见下面的例子。

```
# HELP http_requests_total The total number of HTTP requests.
# TYPE http_requests_total counter
http_requests_total{method="post",code="200"} 1027 1395066363000
http_requests_total{method="post",code="400"}    3 1395066363000

# Escaping in label values:
msdos_file_access_time_seconds{path="C:\\DIR\\FILE.TXT",error="Cannot find file:\n\"FILE.TXT\""} 1.458255915e9

# Minimalistic line:
metric_without_timestamp_and_labels 12.47

# A weird metric from before the epoch:
something_weird{problem="division by zero"} +Inf -3982045

# A histogram, which has a pretty complex representation in the text format:
# HELP http_request_duration_seconds A histogram of the request duration.
# TYPE http_request_duration_seconds histogram
http_request_duration_seconds_bucket{le="0.05"} 24054
http_request_duration_seconds_bucket{le="0.1"} 33444
http_request_duration_seconds_bucket{le="0.2"} 100392
http_request_duration_seconds_bucket{le="0.5"} 129389
http_request_duration_seconds_bucket{le="1"} 133988
http_request_duration_seconds_bucket{le="+Inf"} 144320
http_request_duration_seconds_sum 53423
http_request_duration_seconds_count 144320

# Finally a summary, which has a complex representation, too:
# HELP rpc_duration_seconds A summary of the RPC duration in seconds.
# TYPE rpc_duration_seconds summary
rpc_duration_seconds{quantile="0.01"} 3102
rpc_duration_seconds{quantile="0.05"} 3272
rpc_duration_seconds{quantile="0.5"} 4773
rpc_duration_seconds{quantile="0.9"} 9001
rpc_duration_seconds{quantile="0.99"} 76656
rpc_duration_seconds_sum 1.7560473e+07
rpc_duration_seconds_count 2693
```

#### 可选文本表示 Optional Text Representation
以下三种可选文本格式仅供使用，且不被Prometheus理解。因此，他们的定义可能会有些随意。客户端库可能支持或可能不支持这些格式。工具不应该依赖这些格式。

 1. HTML：此格式由HTTP Content-Type头部请求，其值为text/html。在浏览器中查看的指标是一个“pretty”渲染。虽然生成客户端在技术上完全免费组装HTML，但客户端库之间的一致性应该是针对的。
 2. 协议缓冲区文本格式：与协议缓冲区格式相同，但以文本形式。它由以文本格式（也称为“调试字符串”）连接的协议消息组成，由另外的新行字符分隔（即在协议消息之间有空行）。请求格式为协议缓冲区格式，但HTTP Content-Type头文件中的编码设置为文本。
 3. 协议缓冲区紧凑文本格式：与（2）相同，但使用紧凑文本格式而不是普通文本格式。紧凑文本格式将整个协议消息放在一行上。协议消息仍然以新的行字符分隔，但是不需要“空行”来进行分离。 （每行只需一个协议消息）。格式被请求为协议缓冲区格式，但HTTP Content-Type头中的编码设置为compact-text。
### 历史版本
有关历史格式版本的详细信息，请参阅旧版客户端数据展示格式[文档](https://docs.google.com/document/d/1ZjyKiKxZV83VI9ZKAXRGKaUKK2BIWCT7oiGBKDBpjEY/edit?usp=sharing)。

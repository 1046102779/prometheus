可以使用简单的[基于文本](https://prometheus.io/docs/instrumenting/exposition_formats/#text-based-format)的展示格式向Memetheus公开指标。 有各种各样的客户端库可以为您实现此格式。 如果您的首选语言没有客户端库，则可以[创建自己的语言库](https://prometheus.io/docs/instrumenting/writing_clientlibs/)。

#### 一、基于文本格式
从Prometheus 2.0版开始，向Prometheus公开指标的所有流程都需要使用基于文本的格式。 在本节中，您可以找到有关此格式的一些[基本信息](https://prometheus.io/docs/instrumenting/exposition_formats/#basic-info)以及格式的更详细分类。

##### 1.1 基本信息
| 条目 | 描述 |
|--|----------:|
|成立| 2014.4月|
|支持|Prometheus 版本 `>=0.4.0`
|传输协议| HTTP |
| 编码|  utf-8, `\n`行结尾|
| HTTP `Content-Type`| `text/plain; version=0.0.4`（缺少版本值将导致回退到最新的文本格式版本。）
| Optional HTTP `Content-Encoding`| `gzip` |
| 优点 | -可读性好
|      | -易于组合，特别适用于简单情况（无需嵌套）
|      |- 逐行阅读（处理类型提示和文本字符串外）
| 限制 | - 繁琐
|      | - 类型和文档字符串不是语法的组成部分，意味着很少 - 不存在的度量标准合同验证
|      | - 解析耗时
| 支持度量指标原子性 | - Counter 
| | - Gauge 
| | - Histogram 
| | - Summary 
| | - Untyped 

##### 1.2 文本格式详解

Prometheus基于文本的格式是面向行的。 行由换行符（`\n`）分隔。 最后一行必须以换行符结尾。 空行被忽略。

###### 1.2.1 行格式
在一行内，令牌可以由任意数量的空格和/或制表符分隔（如果它们否则将与前一个令牌合并，则必须至少分开一个）。 前导空格和尾随空格被忽略。

###### 1.2.2 注释，帮助文本和类型信息
带`#`作为第一个非空白字符的行是注释。除`#`之后的第一个标记是`HELP`或`TYPE`，否则它们将被忽略。这些行被视为如下：如果令牌是`HELP`，则预期至少还有一个令牌，即度量标准名称。所有剩余的标记都被视为该度量标准名称的文档字符串。 `HELP`行可以包含任何UTF-8字符序列（在度量标准名称之后），但反斜杠和换行符必须分别转义为`\\`和`\n`。对于任何给定的度量标准名称，只能存在一个HELP行。

如果令牌是`TYPE`，则预计还会有两个令牌。第一个是度量标准名称，第二个是`counter`, `gauge`, `histogram`, `summary`, or `untyped`,，定义该名称的度量标准的类型。给定的度量标准名称只能存在一个`TYPE`行。在为该度量标准名称报告第一个样本之前，必须显示度量标准名称的TYPE行。如果度量标准名称没有`TYPE`行，则类型将设置为`untyped`。

其余行使用以下语法（EBNF）描述样本（每行一个）：
```
metric_name [
  "{" label_name "=" `"` label_value `"` { "," label_name "=" `"` label_value `"` } [ "," ] "}"
] value [ timestamp ]
```
在示例语法中：

- `metric_name`和`label_name`带有通常的Prometheus表达式语言限制。
- `label_value`可以是任何UTF-8字符序列，但反斜杠（`\`，双引号(`"`)和换行符（`\n`））必须分别转义为`\\`，`\"`和`\n`。
- `value`是Go的`ParseFloat()`函数所需的浮点数。 除标准数值外，`Nan`，`+Inf`和`-Inf`分别是表示数字，正无穷大和负无穷大的有效值。
- `timestamp`是`int64`（自纪元以来的毫秒，即1970-01-01 00:00:00 UTC，不包括闰秒），表示为Go的`ParseInt()`函数所需。

###### 1.2.3 分组和排序
给定度量的所有行必须作为一个单独的组提供，首先是可选的`HELP`和`TYPE`行（没有特定的顺序）。 除此之外，重复展示中的可重复排序是优选的，但不是必需的，即如果计算成本过高则不进行排序。

每行必须具有度量标准名称和标签的唯一组合。 否则，摄取行为未定义。

###### 1.2.4 Histograms and summaries
`histogram `and `summary`类型很难以文本格式表示。 以下约定适用：

- 名为`x`的摘要或直方图的样本总和作为名为`x_sum`的单独样本给出。
- 名为`x`的摘要或直方图的样本计数作为名为`x_count`的单独样本给出。
- 名为`x`的摘要的每个分位数作为单独的样本行给出，其具有相同的名称x和标签`{quantile="y"}`。
- 名为`x`的直方图的每个桶计数作为单独的样本行给出，名称为`x_bucket`，标签为`{le="y"}`（其中`y`为桶的上限）。
- 直方图必须有一个`{l ="+Inf"}`的存储桶。 其值必须与`x_count`的值相同。
- 直方图的桶和摘要的分位数必须以其标签值的递增数字顺序出现（分别对于`le`或`quantile`标签）。

##### 1.3 文本格式例子
下面是一个完整的Prometheus度量标准展示示例，包括注释，`HELP`和`TYPE`表达式，直方图，摘要，字符转义示例等。
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
#### 二、历史版本
有关历史格式版本的详细信息，请参阅旧版[客户端数据展示格式](https://docs.google.com/document/d/1ZjyKiKxZV83VI9ZKAXRGKaUKK2BIWCT7oiGBKDBpjEY/edit)文档。

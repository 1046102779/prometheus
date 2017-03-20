## Prometheus查询
---
Prometheus提供一个基本的表达式语言，可以使用户实时地查找和聚合时间序列数据。表达式查询结果可以在图表中展示，查询结果在Prometheus表达式浏览器中展示为表格数据，或者通过HTTP API提供给外部系统使用。

### examples
这个文档仅供参考。对于学习，可能更容易从几个例子开始。

### 表达式语言数据类型

在Prometheus的表达式语言中，一个表达式有四种表达式类型：
 - `instant vector` 即时向量 - 时间序列集合包含了一个度量指标的时间序列集合，所有时间序列共享时间戳
 - `range vector` 范围向量 - 时间序列集合包含一些时间序列数据范围的一组时间序列。
 - `scalar` 标量 - 一个简单的浮点值
 - `string` 字符串  - 一个当前没有被使用的简单字符串

### Literals
#### 字符串
字符串可以用单引号、双引号或者反引号表示

PromQL遵循与Go相同的转义规则。在单引号，双引号中，反斜杠开始一个专一序列，后面可以跟着a, b, f, n, r, t, v或者\。 可以使用八进制(\nnn)或者十六进制(\xnn, \unnnn和\Unnnnnnnn)提供特定字符。

在反引号内不处理转义。与Go不同，Prom不会丢弃反引号中的换行符。例如：
> "this is a string"
> 'these are unescaped: \n \\ \t'
> `these are not unescaped: \n ' " \t"'`

### 浮点数
标量浮点值可以直接写成形式[-](digits)[.(digits)]。
> -2.43

### 时间序列选择器
#### 即时向量选择器
即时向量选择器允许在一个给定的时间戳，选择一组时间序列和对应的样本值：最简单的形式是，仅仅指定度量指标名称，那么含有该度量指标名称的时间序列数据都属于即时向量。

下面这个例子选择所有时间序列度量名称为`http_requests_total`的时间序列数据：
> http_requests_total

通过在度量指标后面增加{}一组标签可以进一步地过滤这些时间序列数据。

下面这个例子在度量指标名称为`http_requests_total`，且一组标签为`job=prometheus`, `group=canary`:
> http_requests_total{job="prometheus",group="canary"}

可以采用不匹配的标签值也是可以的，或者用正则表达式不匹配标签。标签匹配操作如下所示：
  - `=`: 精确地匹配标签给定的值
  - `!=`: 不等于给定的标签值
  - `=~`: 正则表达匹配给定的标签值
  - `!=`: 给定的标签值不符合正则表达式

例如：度量指标名称为`http_requests_total`，正则表达式匹配标签`environment`为`staging, testing, development`的值，且http请求方法不等于`GET`。
> http_requests_total{environment=~"staging|testing|development", method!="GET"}

匹配空标签值得标签匹配器还会选择没有设置任何标签的所有时间序列数据。正则表达式完全匹配。

向量选择器必须指定一个度量指标名称或者至少不能为空字符串的标签匹配器。以下表达式是非法的:
>  {job=~".*"} #Bad!

上面这个例子既没有度量指标名称，也可以匹配空字符串，所以不符合向量选择器的条件

相反地，下面这些表达式是有效地，第一个一定有一个字符。第二个有一个有用的标签method
> {job=~".+"}  # Good!
> {job=~".*", method="get"} # Good!

标签匹配器能够被应用到度量指标名称自己，使用`__name__`标签筛选。例如：表达式`http_requests_total`等价于`{__name__="http_requests_total"}`。 其他的匹配器，如：`= ( !=, =~, !~)`都可以使用。下面的表达式选择了度量指标名称以`job:`开头的时间序列数据：
> {__name__=~"^job:.*"} #
 
#### 范围向量选择器
范围向量类似即时向量, 增加的特性是，它们从当前时刻选择样本范围。在语法上，返回持续时间被追加在向量选择器尾部的方括号[]中，用以指定对于每个生成的范围向量元素应该取回多少时间值。

持续时间的大小有一个数值决定，后面可以跟下面单位中的一个：
 - `s` - seconds
 - `m` - minutes
 - `h` - hours
 - `d` - days
 - `w` - weeks
 - `y` - years

在下面这个例子中，我们选择5分钟内，选择度量指标名称为`http_requests_total`， 标签为`job="prometheus`的时间序列数据:
> http_requests_total{job="prometheus"}[5m]

#### 偏移修饰符
这个`offset`偏移修饰符允许改变向量选择器和标量选择器中的时间偏移

例如，下面的表达式返回度量指标名称为`http_requests_total`，从当前时间到过去5分钟时间内的时间序列值：
> http_requests_total offset 5m

注意：`offset`偏移修饰符总是需要直接跟着选择器，例如：
> sum(http_requests_total{method="GET"} offset 5m} // GOOD.

然而，下面这种情况是不正确的
>  sum(http_requests_total{method="GET"}) offset 5m // INVALID.

offset偏移修饰符在范围向量上和即时向量用法一样的。下面这个返回了度量指标名称为`http_requests_total`, 范围是5分钟粒度，过去一周的速率表达式：
> rate(http_requests_total[5m] offset 1w)

### 操作符
Prometheus支持二元和聚合操作符。详见[表达式语言操作符](https://prometheus.io/docs/querying/operators/)

### 函数
Prometheus提供了一些函数列表操作时间序列数据。详见[表达式语言函数](https://prometheus.io/docs/querying/functions/)

### 陷阱
#### 插值和陈旧
当运行查询后，独立于当前时间序列数据来选择采用数据的时间戳。这主要是为了支持聚合(sum, avg等)情况，其中多喝聚合时间在时间上不完全一致。由于它们的独立性，Prometheus需要在每个相关时间序列的时间戳上分配一个值。他通过简单地在这个时间戳之前采取最新的样本。

如果在采样时间戳之前的5分钟内都没有找到存储的样本，则在此时间点之前没有任何时间序列值。只意味着5分钟之前的图是没有的，其中最新收的样本需要大于过去的5分钟。
```
注意：差值和陈旧处理可能会发生变化。详见[https://github.com/prometheus/prometheus/issues/398](https://github.com/prometheus/prometheus/issues/398)和[https://github.com/prometheus/prometheus/issues/581](https://github.com/prometheus/prometheus/issues/581)

#### 避免慢查询和高负载
如果一个查询需要操作非常大的数据量，图表绘制很可能会曹氏，或者服务器负载过高。因此，在对未知数据构建查询时，始终在Prometheus表达式浏览器的表格视图中开始构建查询，直到结果是看起来合理的（最多为数百个，而不是数千个）。只有当您已经充分过滤或者聚合数据时，才切换到图表模式。如果表达式的查询结果仍然需要很长时间才能绘制出来，则需要通过记录规则重新清洗数据。

像`api_http_requests_total`这样简单的度量指标名称选择器，可以扩展到具有不同标签的数千个时间序列中，这对于Prometheus的查询语言是非常重要的。还要记住，聚合操作即使输出的结果集非常少，但是它会在服务器上产生负载。这类似于关系型数据库查询可一个字段的总和，总是非常缓慢。
```

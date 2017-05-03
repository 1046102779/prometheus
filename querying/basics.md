## Prometheus查询
---
Prometheus提供一个函数式的表达式语言，可以使用户实时地查找和聚合时间序列数据。表达式计算结果可以在图表中展示，也可以在Prometheus表达式浏览器中以表格形式展示，或者作为数据源, 以HTTP API的方式提供给外部系统使用。

### examples
这个文档仅供参考, 这里先举几个容易上手的例子。

### 表达式语言数据类型

在Prometheus的表达式语言中，任何表达式或者子表达式都可以归为四种类型：
 - `instant vector` 瞬时向量 - 它是指在同一时刻，抓取的所有度量指标数据。这些度量指标数据的key都是相同的，也即相同的时间戳。
 - `range vector` 范围向量 - 它是指在任何一个时间范围内，抓取的所有度量指标数据。
 - `scalar` 标量 - 一个简单的浮点值
 - `string` 字符串  - 一个当前没有被使用的简单字符串

依赖于使用场景（例如：图表 vs. 表格），根据用户所写的表达式，仅仅只有一部分类型才适用这种表达式。例如：瞬时向量类型是唯一可以直接在图表中使用的。

### Literals
#### 字符串
字符串可以用单引号、双引号或者反引号表示

PromQL遵循与Go相同的转义规则。在单引号，双引号中，反斜杠成为了转义字符，后面可以跟着a, b, f, n, r, t, v或者\。 可以使用八进制(\nnn)或者十六进制(\xnn, \unnnn和\Unnnnnnnn)提供特定字符。

在反引号内不处理转义字符。与Go不同，Prom不会丢弃反引号中的换行符。例如：
> "this is a string"
> 'these are unescaped: \n \\ \t'
> `these are not unescaped: \n ' " \t"'`

### 浮点数
标量浮点值可以直接写成形式[-](digits)[.(digits)]。
> -2.43

### 时间序列选择器
#### 即时向量选择器
瞬时向量选择器可以对一组时间序列数据进行筛选，并给出结果中的每个结果键值对（时间戳-样本值）: 最简单的形式是，只有一个度量名称被指定。在一个瞬时向量中这个结果包含有这个度量指标名称的所有样本数据键值对。

下面这个例子选择所有时间序列度量名称为`http_requests_total`的样本数据：
> http_requests_total

通过在度量指标后面增加{}一组标签可以进一步地过滤这些时间序列数据。

下面这个例子选择了度量指标名称为`http_requests_total`，且一组标签为`job=prometheus`, `group=canary`:
> http_requests_total{job="prometheus",group="canary"}

可以采用不匹配的标签值也是可以的，或者用正则表达式不匹配标签。标签匹配操作如下所示：
  - `=`: 精确地匹配标签给定的值
  - `!=`: 不等于给定的标签值
  - `=~`: 正则表达匹配给定的标签值
  - `!=`: 给定的标签值不符合正则表达式

例如：度量指标名称为`http_requests_total`，正则表达式匹配标签`environment`为`staging, testing, development`的值，且http请求方法不等于`GET`。
> http_requests_total{environment=~"staging|testing|development", method!="GET"}

匹配空标签值的标签匹配器也可以选择没有设置任何标签的所有时间序列数据。正则表达式完全匹配。

向量选择器必须指定一个度量指标名称或者至少不能为空字符串的标签值。以下表达式是非法的:
>  {job=~".*"} #Bad!

上面这个例子既没有度量指标名称，标签选择器也可以正则匹配空标签值，所以不符合向量选择器的条件

相反地，下面这些表达式是有效的，第一个一定有一个字符。第二个有一个有用的标签method
> {job=~".+"}  # Good!
> {job=~".*", method="get"} # Good!

标签匹配器能够被应用到度量指标名称，使用`__name__`标签筛选度量指标名称。例如：表达式`http_requests_total`等价于`{__name__="http_requests_total"}`。 其他的匹配器，如：`= ( !=, =~, !~)`都可以使用。下面的表达式选择了度量指标名称以`job:`开头的时间序列数据：
> {__name__=~"^job:.*"} #
 
#### 范围向量选择器
范围向量类似瞬时向量, 不同在于，它们从当前实例选择样本范围区间。在语法上，时间长度被追加在向量选择器尾部的方括号[]中，用以指定对于每个样本范围区间中的每个元素应该抓取的时间范围样本区间。

时间长度有一个数值决定，后面可以跟下面的单位：
 - `s` - seconds
 - `m` - minutes
 - `h` - hours
 - `d` - days
 - `w` - weeks
 - `y` - years

在下面这个例子中, 选择过去5分钟内，度量指标名称为`http_requests_total`， 标签为`job="prometheus"`的时间序列数据:
> http_requests_total{job="prometheus"}[5m]

#### 偏移修饰符
这个`offset`偏移修饰符允许在查询中改变单个瞬时向量和范围向量中的时间偏移

例如，下面的表达式返回相对于当前时间的前5分钟时的时刻, 度量指标名称为`http_requests_total`的时间序列数据：
> http_requests_total offset 5m

注意：`offset`偏移修饰符必须直接跟在选择器后面，例如：
> sum(http_requests_total{method="GET"} offset 5m} // GOOD.

然而，下面这种情况是不正确的
>  sum(http_requests_total{method="GET"}) offset 5m // INVALID.

offset偏移修饰符在范围向量上和瞬时向量用法一样的。下面这个返回了相对于当前时间的前一周时，过去5分钟的度量指标名称为`http_requests_total`的速率：
> rate(http_requests_total[5m] offset 1w)

### 操作符
Prometheus支持二元和聚合操作符。详见[表达式语言操作符](https://prometheus.io/docs/querying/operators/)

### 函数
Prometheus提供了一些函数列表操作时间序列数据。详见[表达式语言函数](https://prometheus.io/docs/querying/functions/)

### 陷阱
#### 插值和陈旧
当运行查询后，独立于当前时刻被选中的时间序列数据所对应的时间戳，这个时间戳主要用来进行聚合操作，包括`sum`, `avg`等，大多数聚合的时间序列数据所对应的时间戳没有对齐。由于它们的独立性，我们需要在这些时间戳中选择一个时间戳，并已这个时间戳为基准，获取小于且最接近这个时间戳的时间序列数据。

如果5分钟内，没有获取到任何的时间序列数据，则这个时间戳不会存在。那么在图表中看到的数据都是在当前时刻5分钟前的数据。
```
注意：差值和陈旧处理可能会发生变化。详见[https://github.com/prometheus/prometheus/issues/398](https://github.com/prometheus/prometheus/issues/398)和[https://github.com/prometheus/prometheus/issues/581](https://github.com/prometheus/prometheus/issues/581)

#### 避免慢查询和高负载
如果一个查询需要操作非常大的数据量，图表绘制很可能会超时，或者服务器负载过高。因此，在对未知数据构建查询时，始终需要在Prometheus表达式浏览器的表格视图中构建查询，直到结果是看起来合理的（最多为数百个，而不是数千个）。只有当你已经充分过滤或者聚合数据时，才切换到图表模式。如果表达式的查询结果仍然需要很长时间才能绘制出来，则需要通过记录规则重新清洗数据。

像`api_http_requests_total`这样简单的度量指标名称选择器，可以扩展到具有不同标签的数千个时间序列中，这对于Prometheus的查询语言是非常重要的。还要记住，聚合操作即使输出的结果集非常少，但是它会在服务器上产生负载。这类似于关系型数据库查询可一个字段的总和，总是非常缓慢。
```

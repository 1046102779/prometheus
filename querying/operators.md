## 操作符
---
### 二元操作符
Prometheus的查询语言支持基本的逻辑运算和算术运算。对于两个(instant vector)即时向量。[匹配行为](https://prometheus.io/docs/querying/operators/#vector-matching)可以被改变。

#### 算术二元运算符
在Prometheus系统中支持下面的二元算术操作符：
 - `+` 加法
 - `-` 减法
 - `*` 乘法
 - `/` 除法
 - `%` 模
 - `^` 幂等

二元运算操作符支持`scalar/scalar(标量/标量)`、`vector/scalar(向量/标量)`、和`vector/vector(向量/向量)`之间的操作。

在两个标量之间进行操作符运算，得到的结果也是标量

在向量和标量之间，这个操作符会每个向量的样本值上进行运算。例如：如果一个时间序列向量除以2，操作结果也是一个新的向量，它是每一个向量的每个样本值除以2.

在两个向量之间，一个二元算术操作符作用在左边向量的每一个样本值，且该元素与操作符右边的[向量匹配](https://prometheus.io/docs/querying/operators/#vector-matching)。这个结果是另一个向量。没有匹配到的右边向量entry，不会在结果集中

#### 比较二元操作符
在Prometheus系统中，比较二元操作符有：
 - `==` 等于
 - `!=` 不等于
 - `>`  大于
 - `<`  小于
 - `>=` 大于等于
 - `<=` 小于等于

比较二元操作符被应用于`scalar/scalar（标量/标量）`、`vector/scalar(向量/标量)`，和`vector/vector（向量/向量）`。比较操作符得到的结果是`bool`布尔类型值，返回1或者0值。

在两个标量之间的比较运算，bool结果写入到另一个结果标量中。

即时向量和标量之间的比较运算，是向量中的米一个数据样本值和标量进行比较操作，如果提供了修饰符，则会把运算结果值：true或者false写入到这个修饰符中。否则直接丢弃为false的向量，保留结果为1的向量。所以它们之间的比较运算，相当于一个过滤器，过滤为false的向量

在两个即时向量之间的比较运算，也相当于是一个过滤器，两个向量完全相同，则保留向量，否则，丢弃向量。如果提供了修饰符，则把结果值为true或者false，写入到该修饰符中。

#### 逻辑/集合二元操作符
逻辑/集合二元操作符只能作用在即时向量， 包括：
  - `and` 交集
  - `or`  并集
  - `unless` 补集

`vector1 and vector2` 的逻辑/集合二元操作符，规则：`vector1`向量中的每一个样本数据与`vector2`向量中的所有样本数据进行标签匹配，不匹配的，全部丢弃。运算结果是保留左边的向量元素。

`vector1 or vector2`的逻辑/集合二元操作符，规则: 保留`vector1`向量中的每一个元素，对于`vector2`向量元素，则不匹配`vector1`向量的任何元素，则追加到结果元素中。

`vector1 unless vector2`的逻辑/集合二元操作符，又称差积。规则：包含在`vector1`中的元素，但是该元素不在`vector2`向量所有元素列表中，则写入到结果集中。
### 向量匹配
向量之间的匹配是指右边向量中的每一个元素，在左边向量中也存在。这里有两种基本匹配行为特征：
 - 一对一，找到这个操作符的两边向量元素的相同元素。默认情况下，操作符的格式是`vector1 [operate] vector2`。如果它们有相同的标签和值，则表示相匹配。`ingoring`关键字是指，向量匹配时，可以忽略指定标签。`on`关键字是指，在指定标签上进行匹配。格式如下所示：
> [vector expr] [bin-op] ignoring([label list]) [vector expr]

> [vector expr] [bin-op] on([lable list]) [vector expr]

例如样本数据：
> method_code:http_errors:rate5m{method="get", code="500"} 24
> method_code:http_errors:rate5m{method="get", code="404"} 30
> method_code:http_errors:rate5m{method="put", code="501"} 3
> method_code:http_errors:rate5m{method="post", code="404"} 21

> method:http_requests:rate5m{method="get"} 600

> method:http_requests:rate5m{method="delete"} 34

> method:http_requests:rate5m{method="post"} 120

查询例子：
> method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m

两个向量之间的除法操作运算的向量结果是，每一个向量样本http请求方法标签值是500，且在过去5分钟的运算值。如果没有忽略`code="500"`的标签，这里不能匹配到向量样本数据。两个向量的请求方法是`put`和`delete`的样本数据不会出现在结果列表中

> {method="get"} 0.04     // 24 / 600
> {method="post"} 0.05    //  6 / 120

多对一和一对多的匹配，是指向量元素中的一个样本数据匹配标签到了多个样本数据标签。这里必须直接指定两个修饰符`group_left`或者`group_right`， 左或者右决定了哪边的向量具有较高的子集。
> \<vector expr\> \<bin-op\> ignoring(\<label list\>) group_left(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> ignoring(\<label list\>) group_right(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> on(\<label list\>) group_left(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> on(\<label list\>) group_right(\<label list\>) \<vector expr\>

这个group带标签的修饰符标签列表包含了“一对多”中的“一”一侧的额外标签。对于`on`标签只能是这些列表中的一个。结果向量中的每一个时间序列数据都是唯一的。

`group`修饰符只能被用在比较操作符和算术运算符。

查询例子：
> method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m

在这个例子中，左向量的标签数量多于左边向量的标签数量，所以我们使用`group_left`。右边向量的时间序列元素匹配左边的所有相同`method`标签:
> {method="get", code="500"} 0.04   // 24 /600

> {method="get", code="404"} 0.05   // 30 /600

> {method="post", code="500"} 0.05  //  6 /600

> {method="post", code="404"} 0.175     // 21 /600


多对一和一对多匹配应该更多地被谨慎使用。经常使用`ignoring(\<labels\>)`输出想要的结果。

### 聚合操作符
Prometheus支持下面的内置聚合操作符。这些聚合操作符被用于聚合单个即时向量的所有时间序列列表，把聚合的结果值存入到新的向量中。

 - `sum` (在维度上求和)
 - `max` (在维度上求最大值)
 - `min` (在维度上求最小值)
 - `avg` (在维度上求平均值)
 - `stddev` (求标准差)
 - `stdvar` (求方差)
 - `count` (统计向量元素的个数)
 - `count_values` (统计相同数据值的元素数量)
 - `bottomk` (样本值第k个最小值)
 - `topk` (样本值第k个最大值)
 - `quantile` (统计分位数)

这些操作符被用于聚合所有标签维度，或者通过`without`或者`by`子句来保留不同的维度。
> \<aggr-op\>([parameter,] \<vector expr\>) [without | by (\<label list\>)] [keep_common]

`parameter`只能用于`count_values`, `quantile`, `topk`和`bottomk`。`without`移除结果向量中的标签集合，其他标签被保留输出。`by`关键字的作用正好相反，即使它们的标签值在向量的所有元素之间。`keep_common`子句允许保留额外的标签（在元素之间相同，但不在by子句中的标签）

`count_values`对每个唯一的样本值输出一个时间序列。每个时间序列都附加一个标签。这个标签的名字有聚合参数指定，同时这个标签值是唯一的样本值。每一个时间序列值是结果样本值出现的次数。

`topk`和`bottomk`与其他输入样本子集聚合不同，返回的结果中包括原始标签。`by`和`without`仅仅用在输入向量的桶中

例如：
如果度量指标名称`http_requests_total`包含由`group`, `application`, `instance`的标签组成的时间序列数据，我们可以通过以下方式计算去除`instance`标签的http请求总数：
> sum(http_requests_total) without (instance)

如果我们对所有应用程序的http请求总数，我们可以简单地写下：
> sum(http_requests_total)

统计每个编译版本的二进制文件数量，我们可以如下写：
> count_values("version", build_version)

通过所有实例，获取http请求第5个最大值，我们可以简单地写下：
> topk(5, http_requests_total)

### 二元运算符优先级
在Prometheus系统中，二元运算符优先级从高到低：
 1. ^
 2. *, /, %
 3. +, -
 4. ==, !=, <=, <, >=, >
 5. and, unless
 6. or

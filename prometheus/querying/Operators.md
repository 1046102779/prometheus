##### 一、二元操作符
Prometheus的查询语言支持基本的逻辑运算和算术运算。对于两个瞬时向量, [匹配行为](https://prometheus.io/docs/querying/operators/#vector-matching)可以被改变。

###### 1.1 算术二元运算符
在Prometheus系统中支持下面的二元算术操作符：
 - `+` 加法
 - `-` 减法
 - `*` 乘法
 - `/` 除法
 - `%` 模
 - `^` 幂等

二元运算操作符定义在`scalar/scalar(标量/标量)`、`vector/scalar(向量/标量)`、和`vector/vector(向量/向量)`之间。

在两个标量之间，行为是显而易见的：它们评估另一个标量，这是运算符应用于两个标量操作数的结果。

在即时向量和标量之间，将运算符应用于向量中的每个数据样本的值。 例如。 如果时间序列即时向量乘以2，则结果是另一个向量，其中原始向量的每个样本值乘以2。

在两个即时向量之间，二进制算术运算符应用于左侧向量中的每个条目及其右侧向量中的[匹配元素](https://prometheus.io/docs/prometheus/latest/querying/operators/#vector-matching)。 结果将传播到结果向量中，并删除度量标准名称。 可以找到右侧向量中没有匹配条目的条目不是结果的一部分。

###### 1.2 比较二元操作符
在Prometheus系统中，比较二元操作符有：
 - `==` 等于
 - `!=` 不等于
 - `>`  大于
 - `<`  小于
 - `>=` 大于等于
 - `<=` 小于等于

比较二元操作符定义在`scalar/scalar（标量/标量）`、`vector/scalar(向量/标量)`，和`vector/vector（向量/向量）`。默认情况下他们过滤。 可以通过在运算符之后提供`bool`来修改它们的行为，这将为值返回`0`或`1`而不是过滤。

在两个标量之间，必须提供`bool`修饰符，并且这些运算符会产生另一个标量，即`0`（假）或`1`（真），具体取决于比较结果。

在即时向量和标量之间，将这些运算符应用于向量中的每个数据样本的值，并且从结果向量中删除比较结果为假的向量元素。 如果提供了`bool`修饰符，则将被删除的向量元素的值为`0`，而将保留的向量元素的值为1。

在两个即时向量之间，这些运算符默认表现为过滤器，应用于匹配条目。 表达式不正确或在表达式的另一侧找不到匹配项的向量元素将从结果中删除，而其他元素将传播到具有其原始（左侧）度量标准名称的结果向量中 标签值。 如果提供了`bool`修饰符，则已经删除的向量元素的值为`0`，而保留的向量元素的值为`1`，左侧标签值为`1`。

###### 1.3 逻辑/集合二元操作符
逻辑/集合二元操作符只能作用在即时向量， 包括：
  - `and` 交集
  - `or`  并集
  - `unless` 补集

`vector1 and vector2`得到一个由`vector1`元素组成的向量，其中`vector2`中的元素具有完全匹配的标签集。 其他元素被删除。 度量标准名称和值从左侧向量转移

`vector1 or vector2`得到包含`vector1`的所有原始元素（标签集+值）的向量以及`vector2`中`vector1`中没有匹配标签集的所有元素。。

`vector1 unless vector2`得到一个由`vector1`元素组成的向量，其中`vector2`中没有元素，具有完全匹配的标签集。 两个向量中的所有匹配元素都被删除。
##### 二、向量匹配
向量之间的操作尝试在左侧的每个条目的右侧向量中找到匹配元素。 匹配行为有两种基本类型：一对一和多对一/一对多。
###### 2.1 一对一向量匹配
一对一从操作的每一侧找到一对唯一条目。 在默认情况下，这是格式为`vector1<operator>vector2`之后的操作。 如果两个条目具有完全相同的标签集和相应的值，则它们匹配。 忽略关键字允许在匹配时忽略某些标签，而`on`关键字允许将所考虑的标签集减少到提供的列表：
> [vector expr] [bin-op] ignoring([label list]) [vector expr]

> [vector expr] [bin-op] on([lable list]) [vector expr]

例如样本数据：
```
 method_code:http_errors:rate5m{method="get", code="500"} 24
 method_code:http_errors:rate5m{method="get", code="404"} 30
 method_code:http_errors:rate5m{method="put", code="501"} 3
 method_code:http_errors:rate5m{method="post", code="404"} 21

 method:http_requests:rate5m{method="get"} 600
 method:http_requests:rate5m{method="delete"} 34
 method:http_requests:rate5m{method="post"} 120
```

查询例子：
> method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m

这将返回一个结果向量，其中包含每个方法的状态代码为500的HTTP请求部分，在过去的5分钟内进行测量。 没有`ignoring(code)`就没有匹配，因为度量标准不共享同一组标签。 方法`put`和`del`的条目没有匹配，并且不会显示在结果中：

> {method="get"} 0.04     // 24 / 600

> {method="post"} 0.05    //  6 / 120

###### 2.2 多对一和一对多向量匹配
多对一和一对多匹配指的是“一”侧的每个向量元素可以与“多”侧的多个元素匹配的情况。 必须使用`group_left`或`group_right`修饰符明确请求，其中`left/right`确定哪个向量具有更高的基数。
> \<vector expr\> \<bin-op\> ignoring(\<label list\>) group_left(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> ignoring(\<label list\>) group_right(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> on(\<label list\>) group_left(\<label list\>) \<vector expr\>

> \<vector expr\> \<bin-op\> on(\<label list\>) group_right(\<label list\>) \<vector expr\>

随组修饰符提供的标签列表包含来自“一”侧的其他标签，以包含在结果度量标准中。 对于标签，只能出现在其中一个列表中。 每次结果向量的序列必须是唯一可识别的。

分组修饰符只能用于比较和算术。 默认情况下，操作as和除非和或操作与右向量中的所有可能条目匹配。

示例查询：
> method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m

在这种情况下，左向量每个`method`标签值包含多个条目。 因此，我们使用`group_left`表明这一点。 右侧的元素现在与多个元素匹配，左侧具有相同的`method`标签：
> {method="get", code="500"} 0.04   // 24 /600
> {method="get", code="404"} 0.05   // 30 /600

> {method="post", code="500"} 0.05  //  6 /600

> {method="post", code="404"} 0.175     // 21 /600


多对一和一对多匹配是高级用例，应该仔细考虑。 通常正确使用忽略`ignoring(<labels>) `可提供所需的结果。

##### 三、聚合操作符
Prometheus支持以下内置聚合运算符，这些运算符可用于聚合单个即时向量的元素，从而生成具有聚合值的较少元素的新向量：

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

这些运算符可以用于聚合所有标签维度，也可以通过包含`without`或`by`子句来保留不同的维度。
> \<aggr-op\>([parameter,] \<vector expr\>) [without | by (\<label list\>)] [keep_common]

`parameter`仅用于`count_values`，`quantile`，`topk`和`bottomk`。不从结果向量中删除列出的标签，而所有其他标签都保留输出。 `by`相反并删除未在`by`子句中列出的标签，即使它们的标签值在向量的所有元素之间是相同的。

`count_values`输出每个唯一样本值的一个时间序列。每个系列都有一个额外的标签。该标签的名称由聚合参数给出，标签值是唯一的样本值。每个时间序列的值是样本值存在的次数。

`topk`和`bottomk`与其他聚合器的不同之处在于，输入样本的子集（包括原始标签）在结果向量中返回。 `by`和`without`仅用于存储输入向量。

例：

如果度量标准`http_requests_total`具有按应用程序，实例和组标签扇出的时间序列，我们可以通过以下方式计算每个应用程序和组在所有实例上看到的HTTP请求总数：
> sum(http_requests_total) without (instance)

等价于：
> sum(http_requests_total)

要计算运行每个构建版本的二进制文件的数量，我们可以编写：
> count_values("version", build_version)

要在所有实例中获取5个最大的HTTP请求计数，我们可以编写：
> topk(5, http_requests_total)

##### 四、二元运算符优先级
以下列表显示了Prometheus中二进制运算符的优先级，从最高到最低。
 1. ^
 2. *, /, %
 3. +, -
 4. ==, !=, <=, <, >=, >
 5. and, unless
 6. or

具有相同优先级的运算符是左关联的。 例如，`2 * 3％2`相当于`（2 * 3）％2`。但是`^`是右关联的，因此`2 ^ 3 ^ 2`相当于`2 ^（3 ^ 2）`。

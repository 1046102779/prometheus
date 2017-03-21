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
> \<vectr\>


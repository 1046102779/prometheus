## 函数列表
---
一些函数有默认的参数，例如：`year(v=vector(time()) instant-vector)`。v是参数值，instant-vector是参数类型。vector(time())是默认值。

### abs()
`abs(v instant-vector)`返回输入向量的所有样本的绝对值。

### absent()
`absent(v instant-vector)`，如果赋值给它的向量具有任何元素，则返回空向量；如果传递给该函数的向量没有元素，则返回这个传入向量的所有标签，同时向量值等于1。

> absent(nonexistent{job="myjob"})  # => {job="myjob"} 1

> absent(nonexistent{job="myjob", instance=~".*"}) # => {job="myjob"} 1
> so smart !

> absent(sum(nonexistent{job="myjob"})) # => {} 0

对于给定的度量指标表达式，如果没有样本数据结果，我们使用absent()函数打印输出是非常有用的。

### ceil()
`ceil(v instant-vector)` 是一个四舍五入的函数。它将v中的所有元素的样本值四舍五入到最接近的整数。

### changes()
对于每一个输入的时间序列数据，`changes(v range-vector)` 返回一个向量表达式在给定时间内向量样本值变化的次数。

### clamp_max()
`clamp_max(v instant-vector, max scalar)`函数，把所有向量元素的样本值不能超过最大范围max限制。

### clamp_min()
`clamp_min(v instant-vector)`函数，把所有向量的元素样本值设置在不能低于最小范围min值。

### count_saclar()
`count_scalar(v instant-vector)` 函数, 返回值是一个scalar标量，表示v向量的元素个数。而`count()`聚合函数正好相反，它总是返回一个向量和向量值为元素个数，并允许通过`by`条件分组。

### day_of_month()
`day_of_month(v=vector(time()) instant-vector)`函数，返回被给定UTC时间所在月的第几天。返回值范围：1~31。

### day_of_week()
`day_of_week(v=vector(time()) instant-vector)`函数，返回被给定UTC时间所在月的天数。返回值范围：28~31。

### delta()
`delta(v range-vector)`函数，计算一个范围向量v的第一个元素和最后一个元素之间的差值。返回相同向量，且向量值为这个差值。这个值以差值的形式插入时间序列各个点中。

下面这个表达式例子，返回过去两小时的CPU温度差：
> delta(cpu_temp_celsius{host="zeus"}[2h])

`delta`差值仅仅用在gauges上。

### deriv()
`deriv(v range-vector)`函数，计算一个范围向量v中各个时间序列二阶导数，使用[简单线性回归](https://en.wikipedia.org/wiki/Simple_linear_regression)

`deriv`二阶导数仅仅用在gauges。

### drop_common_labels()
`drop_common_labels(instant-vector)`函数，返回去掉形参向量的所有标签，且向量值不变。

### exp()
`exp(v instant-vector)`函数计算向量v值的指数函数。特殊情况如下所示：
> Exp(+inf) = +Inf
> Exp(NaN) = NaN

### floor()
`floor(v instant-vector)`函数，是表示向量样本值小于等于该值的最接近整数。

### histogram_quantile()
`histogram_quatile(φ float, b instant-vector)` 函数计算b向量的φ-直方图 (0 ≤ φ ≤ 1) 。参考中文文献[https://www.howtoing.com/how-to-query-prometheus-on-ubuntu-14-04-part-2/]

### holt_winters()
`holt_winters(v range-vector, sf scalar, tf scalar)`函数基于范围向量v，生成事件序列数据平滑值。平滑因子`sf`越低, 对老数据越重要。趋势因子`tf`越高，越多的数据趋势应该被重视。0< sf, tf <=1。 `holt_winters`仅用于gauges

### hour()
`hour(v=vector(time()) instant-vector)`函数返回被给定UTC时间的当前第几个小时，时间范围：0~23。

### idelta()
`idelta(v range-vector)`函数，返回范围向量v的最后两个样本值的差值和相同向量标签。

### increase()
`increase(v range-vector)`计算范围向量中时间序列数据的增长，自动调整单调性，如：服务实例重启，则计数器重置。

下面的表达式例子，返回过去5分钟，连续两个时间序列数据样本值的http请求差值。
> increase(http_requests_total{job="api-server"}[5m])

`increase`仅仅用于统计，主要作用是增加图表和数据的可读性，使用`rate`记录规则的使用率，以便持续跟踪数据样本值的变化。

### irate
`irate(v range-vector)`函数, 计算范围向量中时间序列的每秒增长速率。它是基于最后两个数据点，自动调整单调性， 如：服务实例重启，则计数器重置。

下面表达式针对范围向量中的每个时间序列数据，返回两个最新数据点过去5分钟的HTTP请求速率。
> irate(http_requests_total{job="api-server"}[5m])

`irate`只能用于绘制快速移动的计数器。因为速率的简单更改可以重置FOR子句，利用警报和缓慢移动的计数器，完全由罕见的尖峰组成的图形很难阅读。

### label_replace()
对于v中的每个时间序列，`label_replace(v instant-vector, dst_label string, replacement string, src_label string, regex string)` 将正则表达式与标签src_label匹配。如果匹配，则返回时间序列，标签dst_label被替换的扩展替换。$1替换为第一个匹配子组，$2替换为第二个等。如果正则表达式不匹配，则时间序列不会更改。

下面这个例子返回一个向量值a带有`foo`标签：
`label_replace(up{job="api-server", serice="a:c"}, "foo", "$1", "service", "(.*):.*")`

### ln()
`ln(v instance-vector)`计算向量v中所有元素的自然对数。特殊例子：
 > ln(+Inf) = +Inf
 > ln(0) = -Inf
 > ln(x<0) = NaN
 > ln(NaN) = NaN

### log2()
`log2(v instant-vector)`函数计算向量v中所有元素的二进制对数。

### log10()
`log10(v instant-vector)`函数计算向量v中所有元素的10进制对数。相当于ln()

### minute()
`minute(v=vector(time()) instant-vector)`函数返回给定UTC时间当前小时的第多少分钟。结果范围：0~59。

### month()
`month(v=vector(time()) instant-vector)`函数返回给定UTC时间当前属于第几个月，结果范围：0~12。

### predict_linear()
`predict_linear(v range-vector, t scalar)`函数，是指基于一个范围向量v，预测从现在起t秒内的时间序列数据值。

### rate()
`rate(v range-vector)`函数计算范围向量v的时间序列数据样本值的增长率。

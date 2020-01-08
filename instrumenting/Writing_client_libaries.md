本文档介绍了Prometheus客户端库应提供的功能和API，旨在实现库之间的一致性，简化易用用例，避免提供可能导致用户走错路的功能。

在撰写本文时已经支持了[10种语言](https://prometheus.io/docs/instrumenting/clientlibs/)，因此我们现在已经很好地理解了如何编写客户端。 这些指南旨在帮助新客户端库的作者生成良好的库。

##### 一、Conventions约定
MUST/MUST NOT/SHOULD/SHOULD NOT/MAY具有给出的含义在[https://www.ietf.org/rfc/rfc2119.txt](https://www.ietf.org/rfc/rfc2119.txt)

此外，ENCOURAGED意味着某个功能对于库来说是理想的，但如果它不存在则可以。 换句话说，一个很好的。

记住下面的几点：
 - 利用每种语言的功能。
 - 常用用例应该很简单
 - 做事情正确方式是简单的方法
 - 更复杂的例子应该是可能的

常用用例（有序）：
 - 没有标签的Counters在库/应用程序之间传播
 - Summaries/Histograms的时序功能/代码块
 - Gauges跟踪事情的当前状态
 - 批量任务监控

##### 二、总体结构
必须将客户端编写为内部回调。客户通常应该遵循这里描述的结构。

关键类是`Collector`。有一个方法（通常称为`collect`），返回零个或多个指标及其样本。`Collector`在`CollectorRegistry`注册。通过将`CollectorRegistry`传递给`class/method/function``bridge`来公开数据，该类以Prometheus支持的格式返回指标。每次抓取`CollectorRegistry`时，它都必须回调每个`Collector`的`collect`方法。

大多数用户与之交互的界面是`Counter`，`Gauge`，`Summary`和`Histogram` Collectors。这些代表一个度量标准，应涵盖用户正在使用自己的代码的绝大多数用例。

更高级的用例（例如从另一个监视/检测系统代理）需要编写自定义`Collector`。有人可能还想编写一个`bridge`，它采用`CollectorRegistry`并以不同监控/仪表系统理解的格式生成数据，从而允许用户只需考虑一个仪器系统。

`CollectorRegistry`应该提供`register()/unregister()`函数，并且应该允许收集器注册到多个`CollectorRegistrys`。

客户端库必须是线程安全的。

对于诸如C的非OO语言，客户端库应该尽可能地遵循这种结构的精神。

###### 2.1 命名
客户端库应该遵循本文档中提到的`function/method/class`，记住它们所使用的语言的命名约定。例如，`set_to_current_time()`适用于方法名称Python，但在Go中`SetToCurrentTime()`更好，`setToCurrentTime()`是Java中的约定。 如果名称因技术原因而不同（例如，不允许函数重载），文档/帮助字符串应该将用户指向其他名称。

库不得提供与此处给出的名称相同或相似的函数/方法/类，但具有不同的语义。

##### 三、Metrics
`Counter`、`Gauge`、`Summary`和`Histogram`[度量指标类型](https://prometheus.io/docs/concepts/metric_types/)是最主要的接口。

`Counter`和`Gauge`必须是客户库的一部分。`Summary`和`Histogram`至少被提供一个。

这些应该主要用作文件静态变量，即在与它们正在检测的代码相同的文件中定义的全局变量。客户端库应该启用它。常见的用例是整体编写一段代码，而不是在一个对象实例的上下文中编写代码。用户不必担心在他们的代码中管理他们的指标，客户端库应该为他们做这些（如果没有，用户将在库周围编写一个包装器以使其“更容易” - 这很少倾向于好吧）。

必须有一个默认的`CollectorRegistry`，默认情况下，标准指标必须隐式注册到它中，而不需要用户进行任何特殊工作。必须有一种方法可以将指标注册到默认的`CollectorRegistry`，以便在批处理作业和单元测试中使用。定制收藏家也应该遵循这一点。

究竟应该如何创建指标因语言而异。对于某些人（Java，Go），构建器方法是最好的，而对于其他人（Python），函数参数足够丰富，可以在一次调用中完成。

例如，在Java Simpleclient中，我们有：
```Java
class YourClass {
  static final Counter requests = Counter.build()
      .name("requests_total")
      .help("Requests.").register();
}
```

这将使用默认的`CollectorRegistry`注册请求。 通过调用`build()`而不是`register()`，度量标准将不会被注册（方便单元测试），您还可以将`CollectorRegistry`传递给`register()`（便于批处理作业）。

###### 3.1 Counter
[`Counter`](https://prometheus.io/docs/concepts/metric_types/#counter)是一个单调递增的计数器。它不允许counter值下降，但是它可以被重置为0（例如：客户端服务重启）。

一个counter必须有以下方法：
 - `inc()`: 增量为1.
 - `inc(double v)`: 增加给定值v。必须检查v>=0。

一个`Counter`鼓励有：

一种计算在给定代码段中抛出/引发异常的方法，以及可选的仅某些类型的异常。 这是Python中的count_exceptions。

计数器必须从0开始。

###### 3.2 Gauge
[`Gauge`](https://prometheus.io/docs/concepts/metric_types/#gauge)表示一个可以上下波动的值。

gauge必须有以下的方法：
 - `inc()`: 每次增加1
 - `inc(double v)`: 每次增加给定值v
 - `dec()`: 每次减少1
 - `dec(double v)`: 每次减少给定值v
 - `set(double v)`: 设置gauge值成v

Gauges值必须从0开始，你可以为给定的量表提供一种方法，以不同的数字开始。

gauge应该有以下方法：
 - `set_to_current_time()`: 将gauge设置为当前的unix时间（以秒为单位）。

gauge被建议有：
一种跟踪某些代码/功能中正在进行的请求的方法。 这是Python中的`track_inprogress`。

一种为一段代码计时并将仪表设置为其持续时间的方法，以秒为单位。 这对批处理作业很有用。 这是Java中的`startTimer/setDuration`和Python中的`time()`装饰器/上下文管理器。 这应该与`Summary/Histogram`中的模式匹配（尽管是`set()`而不是`observe()`）。

###### 3.3 Summary
[`Summary`](https://prometheus.io/docs/concepts/metric_types/#summary)通过时间滑动窗口抽样观察（通常是要求持续时间），并提供对其分布、频率和总和的即时观察。

`Summary`绝不允许用户将“quantile”设置为标签名称，因为这在内部用于指定摘要分位数。 一个`Summary`是ENCOURAGED提供分位数作为出口，虽然这些不能汇总，往往很慢。 总结必须允许没有分位数，因为`_count/_sum`非常有用，这必须是默认值。

`Summary`必须具有以下方法：

- `observe(double v)`：观察给定量

`Summary`应该有以下方法：

一些方法可以在几秒钟内为用户计时。 在Python中，这是`time()`装饰器/上下文管理器。 在Java中，这是`startTimer/observeDuration`。 绝不能提供秒以外的单位（如果用户想要其他东西，他们可以手工完成）。 这应该遵循与`Gauge/Histogram`相同的模式。

`Summary``_count/_sum`必须从0开始。

###### 3.4 Histogram
[`Histogram`](https://prometheus.io/docs/concepts/metric_types/#histogram)允许可聚合的事件分布，例如请求延迟。 这是每个桶的核心。

`Histogram`绝不允许`le`作为用户设置标签，因为`le`在内部用于指定存储桶。

`Histogram`必须提供一种手动选择存储桶的方法。应该提供以`linear(start, width, count)`和`exponential(start, factor, count)`方式设置桶的方法。计数必须排除`+Inf`桶。

`Histogram`应该与其他客户端库具有相同的默认存储桶。创建度量标准后，不得更改存储桶。

`Histogram`必须有以下方法：

- `observe(double v)`：观察给定量

`Histogram`应该有以下方法：

一些方法可以在几秒钟内为用户计时。在Python中，这是`time()`装饰器/上下文管理器。在Java中，这是`startTimer/observeDuration`。绝不能提供秒以外的单位（如果用户想要其他东西，他们可以手工完成）。这应该遵循与`Gauge/Summary`相同的模式。

`Histogram``_count/_sum`和桶必须从0开始。

进一步的指标考虑

除了上面记录的对于给定语言有意义的指标之外，还提供额外的功能，这是ENCOURAGED。

如果有一个常见的用例，你可以做得更简单然后去做，只要它不会鼓励不良行为（例如次优的度量/标签布局，或在客户端进行计算）。

###### 3.5 标签
标签是普罗米修斯[最强大的方面](https://prometheus.io/docs/practices/instrumentation/#use-labels)之一，但很容易[被滥用](https://prometheus.io/docs/practices/instrumentation/#do-not-overuse-labels)。因此，客户端库必须非常小心地向用户提供标签。

在任何情况下，客户端库都不允许用户为`Gauge/Counter/Summary/Histogram`或库提供的任何其他`Collector`的相同度量标准指定不同的标签名称。

自定义收集器中的度量标准几乎总是具有一致的标签名称。由于仍然存在罕见但有效的用例，但事实并非如此，客户端库不应对此进行验证。

虽然标签功能强大，但大多数指标都没有标签。因此，API应该允许标签但不支配它。

客户端库必须允许在`Gauge/Counter/Summary/Histogram`创建时指定标签名称列表。客户端库应该支持任意数量的标签名称。客户端库必须验证标签名称是否符合记录的要求。

提供对度量标注维度的访问的一般方法是使用`labels()`方法，该方法获取标签值列表或从标签名称到标签值的映射并返回“Child”。然后可以在Child上调用通常的`.inc()/.dec()/.observe()`等方法。

`labels()`返回的子项应该由用户缓存，以避免再次查找 - 这在延迟关键代码中很重要。

带标签的度量标准应该支持一个`remove()`方法，该方法具有与`labels()`相同的签名，它将从不再导出它的度量中删除Child，以及一个从度量中删除所有Children的`clear()`方法。这些无效的缓存儿童。

应该是一种使用默认值初始化给定Child的方法，通常只是调用`labels()`。必须始终初始化没有标签的度量标准以避免缺少度量标准的问题。

###### 3.6 度量指标名称
度量标准名称必须遵循规范。 与标签名称一样，必须满足使用`Gauge/Counter/Summary/Histogram`以及随库提供的任何其他Collector。

许多客户端库提供了三个部分的名称设置：`namespace_subsystem_name`，其中只有`name`是必需的。

除非自定义收集器从其他检测/监视系统进行代理，否则不得禁止动态/生成的度量标准名称或度量标准名称的子部分。 生成/动态度量标准名称是您应该使用标签的标志。

###### 3.7 度量指标描述和帮助
`Gauge/Counter/Summary/Histogram`必须要求提供度量标准描述/帮助。

随客户端库提供的任何自定义收集器必须具有其指标的描述/帮助。

建议将其作为强制性参数，但不要检查它是否具有一定的长度，好像有人真的不想写文档，否则我们不会说服它们。 图书馆提供的收藏家（实际上我们可以在生态系统中的任何地方）应该有很好的度量描述，以身作则。

##### 四、导出

客户必须实现博览会格式文档中概述的基于文本的[导出格式](https://prometheus.io/docs/instrumenting/exposition_formats/)。

如果可以在没有显着资源成本的情况下实现暴露度量的可重现顺序是ENCOURAGED（特别是对于人类可读格式）。

##### 五、标准化和运行时收集器

客户端库应该提供标准导出的功能，如下所示。

这些应该作为自定义收集器实现，并默认注册在默认的CollectorRegistry上。 应该有一种方法来禁用它们，因为有一些非常小的用例会妨碍它们。

###### 5.1 处理度量指标
这些导出应该有前缀`process_`。 如果语言或运行时没有公开其中一个变量，那么它就不会导出它。 所有内存值，以字节为单位，所有时间均为`unixtime/seconds`。

| 度量指标名称          |        含义                |    单位   |
| --------------------------| :-------------------------:| ---------:|
| process_cpu_seconds_total | 用户和系统CPU花费的时间    |  秒       |
| process_open_fds          | 打开的文件描述符数量       | 文件描述符|
| process_max_fds           | 打开描述符最大值           | 文件描述符|
| process_virtual_memory_bytes| 虚拟内存大小             | 字节|
| process_virtual_memory_max_bytes| 最大可用虚拟内存量（以字节为单位）             | 字节|
| process_resident_memory_bytes| 驻留内存大小|字节|
| process_heap_bytes | 进程head堆大小| 字节|
| process_start_time_seconds| unix时间 | 秒|

###### 5.2 运行时度量指标
此外，还鼓励客户端库提供其语言运行时（例如垃圾收集统计信息）的度量标准，并提供适当的前缀，如`go_`，`hostspot_`等。
##### 六、单元测试
客户端库应该有单元测试，涵盖核心工具库和博览会。

客户端库鼓励提供方便用户对其使用仪器代码进行单元测试的方法。 例如，Python中的`CollectorRegistry.get_sample_value`。

##### 七、包和依赖
理想情况下，客户端库可以包含在任何应用程序中，以便在不破坏应用程序的情况下添加一些检测。

因此，在向客户端库添加依赖项时，建议谨慎。 例如，如果添加使用Prometheus客户端的库，该客户端需要x.y版本的库但应用程序在其他地方使用x.z，那么这会对应用程序产生负面影响吗？

建议在可能出现这种情况时，将核心工具与给定格式的度量的桥梁/展示分开。 例如，Java simpleclient `simpleclient`模块没有依赖关系，`simpleclient_servlet`具有HTTP位。

##### 八、性能考虑
由于客户端库必须是线程安全的，因此需要某种形式的并发控制，并且必须考虑多核机器和应用程序的性能。

根据我们的经验，效果最差的是互斥体。

处理器原子指令往往位于中间，并且通常是可接受的。

避免不同CPU改变相同RAM的方法最有效，例如Java的simpleclient中的DoubleAdder。但是有记忆费用。

如上所述，`labels()`的结果应该是可缓存的。倾向于使用标签返回度量标准的并发映射往往相对较慢。没有标签的特殊套管指标可以避免`labels()`- 像查找一样可以提供很多帮助。

度量标准应当在递增/递减/设置等时避免阻塞，因为在刮擦正在进行时整个应用程序被阻止是不可取的。

主要仪器操作（包括标签）的基准测试是鼓励的。

在进行博览会时，应牢记资源消耗，特别是RAM。考虑通过流式传输结果减少内存占用量，并可能限制并发擦除次数。

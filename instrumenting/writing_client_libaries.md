## 编写客户端库
---
这篇文档包括Prometheus客户端API应该提供的基础功能，目的是在客户端库之间保持一致性，轻松上手并避免提供导致用户出错的功能。

已经有10种客户端语言支持Prometheus客户端了，因此我们知道怎么写好一个客户端。这个指南旨在帮助写Prometheus客户端其他语言的作者写一个好的库。

### Conventions约定
MUST/MUST NOT/SHOULD/SHOULD NOT/MAY在[https://www.ietf.org/rfc/rfc2119.txt](https://www.ietf.org/rfc/rfc2119.txt)

另一个ENCOURAGE的含义是，一个特性对于一个库是非常好的，但是如果关闭这个特性的话，不会影响库的使用。

记住下面的几点：
 - 记住每个特性的好处。
 - 常用用例应该很简单
 - 做事情正确方式是简单的方法
 - 更复杂的例子应该是可能的

常用用例（有序）：
 - 没有标签的Counters在库或者应用程序之间传播
 - Summaries/Histograms的时序功能/代码块
 - Gauges跟踪事情的当前状态
 - 批量任务监控

### 总体结构
客户端`必须`在内部写入回调。客户通常`应该`遵循下面描述的结构。

这个关键类是`Collector`。这个有一个典型的方法`collect`, 返回0~N个度量指标和这些指标的样本数据。`Collector`用`CollectorRegistry`进行注册。通过传递`CollectorRegistry`给称之为`bridge`的class/method/function来暴露数据。该`bridge`返回Prometheus支持的数据格式数据。每次这个`CollectorRegistry`被收集时，都必须回调`Collector`的collect方法。

和用户交互最多的接口是`Counter`, `Gauge`, `Summary`和`Histogram Collectors`。这些表示单个度量指标，写的代码覆盖绝大多数的用例。

更高级的用例（例如来自其他监控/检测系统的代理）需要编写一个自定义`Collector`收集器。有人也可能像写一个带有`CollectorRegistry`的"bridge"，以不同的监控/测量系统理解的格式生产数据, 允许用户只需要考虑一个测量系统。

`CollectorRegistry`应该提供`register()/unregister()`方法，以及一个`Collector`应该注册多个`CollectorRegistrys`

客户库必须是线程安全的。

对于非面向对象的客户端，如：C语言，客户库编写在实践中应该遵循这种结构的理念。

#### 命名
客户库应该遵循`function/method/class`在这个文档中提及的命名规则，记住他们正在使用的语言命名规范。例如：`set_to_current_time()`对于python而言是非常好的方法名称，`SetToCurrentTime`对于Go语言是更好的，`setToCurrentTime()`对于Java是更好的。由于技术原因（例如：不允许功能重载），名称不能，文档/帮助文档应该将用户指向其他名称。

库禁止提供与此处给出的相同或者相似`functions/methods/classes`，但具有不同的语义。

### Metrics
`Counter`、`Gauge`、`Summary`和`Histogram`度量指标类型是最主要的接口。

`Counter`和`Gauge`必须是客户库的一部分。`Summary`和`Histogram`至少被提供一个。

这些主要用作文件静态变量，也就是说，全局变量与他们正在调试的代码在同一个文件中定义。客户端库应该启用此功能。常见的用例是整体测试一段代码，而不是在对象的一个实例上下文中的一段代码。用户不必担心在他们的代码中管理他们的指标，客户端库应该为他们做到这一点（如果不这样做，用户将会围绕库写一个`wrapper`, 使其更容易，少即是多）。

必须有一个默认的`CollectorRegistry`， 标准的度量指标必须默认被注册，不需要用户干预。必须有一种方法，使度量指标不默认注册到`CollectorRegistry`中，用于批处理作业和单元测试。自定义的`Collectors`也应该遵循这点。

究竟应该如何创建度量指标因语言而异。对于某些语言（Go，Java），构建器是最好的，对于其他（Python）函数参数足够丰富，可以在一个调用中执行。

例如，一个简单的Java客户端，我们可以这样写：
```Java
class YourClass {
  static final Counter requests = Counter.build()
      .name("requests_total")
      .help("Requests.").register();
}
```

使用默认的`CollectorRegistry`进行注册。通过调用build()而不是register(), 度量指标将不会被注册（对于单元测试来说很方便），你还可以将`CollectorRegistry`传递给register()(方便批作业处理)。

#### Counter
`Counter`[https://prometheus.io/docs/concepts/metric_types/#counter]是一个单调递增的计数器。它不允许counter值下降，但是它可以被重置为0（例如：客户端服务重启）。

一个counter必须有以下方法：
 - `inc()`: 增量为1.
 - `inc(double v)`: 增加给定值v。必须检查v>=0。

Counter在给定代码段抛出/引发异常的方式，也可以只选择某些类型的一场，这是Python中的count_exceptions。

Counters必须从0开始。

#### Gauge
[Gauge](https://prometheus.io/docs/concepts/metric_types/#gauge)表示一个可以上下波动的值。

gauge必须有以下的方法：
 - `inc()`: 每次增加1
 - `inc(double v)`: 每次增加给定值v
 - `dec()`: 每次减少1
 - `dec(double v)`: 每次减少给定值v
 - `set(double v)`: 设置gauge值成v

Gauges值必须从0开始，你可以提供一个从不等于0的值开始。

gauge应该有以下方法：
 - `set_to_current_time()`: 将gauge设置为当前的unix时间（以秒为单位）。

gauge被建议有：
 - 一种在某些代码/方法中跟踪正在进行的请求方法。这是python种的`track_inprogress`。

执行代码块的时间，并将测量仪设置为其持续时间（秒），这对于批量任务是非常有用的。在Java中是`startTimer/setDuration`， 在python中是`time()` decorator/上下文管理器。这应该符合在`Summary`和`Histogram`中的pattern(通过`set()`而不是`observe()`)。

#### Summary
[summary](https://prometheus.io/docs/concepts/metric_types/#summary)通过时间滑动窗口抽样观察（通常是要求持续时间），并提供对其分布、频率和总和的即时观察。

summary不允许用户设置"quantile"作为一个标签，因为这个名称已在内部使用，用来指定分位数。summary鼓励提供“quantile”导出，虽然这些不能被汇总，而且需要大量时间。summary必须允许没有quantiles，因为只有`_count/_sum`是飞铲更拥有的，这必须是默认值。

summary必须有以下方法：
 - `observe(double v)`: 观察被给定值

summary应该有以下方法：
 - 统计用户执行代码的时间，以秒为单位。在python中，这是`time()`decorateor/context管理器。在Java中这是`startTimer/observeDuration`。 不能提供秒意外的单位（如果用户想要别的，自己手动做）。这应该遵循Gauge/Histogram相同的模式。

Summary `_count/_sum`必须从0开始。

#### Histogram
[Histogram](https://prometheus.io/docs/concepts/metric_types/#histogram)允许时间的可聚合分布，如：请求延迟。这是counter/bucket的核心。

一个histogram直方图不允许使用`le`作为一个用户集合标签，该标签内部用于指定buckets。

直方图必须提供一个方法来手动选择buckets。应该提供一现行（start，factor和count）和指数（start，factor和count）方式设置buckets的方法。counter必须排序+Inf bucket

直方图应该具有与其他客户端库相同的默认值，创建度量指标后bucket不能再更改。

一个直方图必须有下面的方法：
 - `observe(double v)`: 观察给定值

直方图应该有以下的方法：
统计代码执行时间的一些方法，以秒为单位。在Python中是`time()`decorator/context管理器。在Java中是`startTimer/observeDuration`。不提供秒以外的单位（如果用户需要别的，可以手动做）。这应该遵循与Gauge/Summary相同的模式。

直方图`_count/_sum`和buckets必须从0开始。

进一步的度量指标考量
提供额外的功能，超出以上记录的指标，对于给定的语言是有意义的

如果有一个常见用例，例如：次优度量指标/标签布局或者在客户端进行计算，可以使其更简单。

#### 标签
标签Labels是Prometheus系统最强大的一个方面，但是很容易被滥用。因此，客户端库必须非常小心地如何向用户提供labels。

客户库在任何情况下不得允许用户对于"Gauge/counter/summary/histogram"或者由库提供的其他Collector相同度量指标名称，有相同不同的labels名称。

如果你的客户库在收集时刻对其进行了度量指标的验证，那么它也可以为自定义Collector进行验证。

虽然标签功能很强大，但大多数度量指标不会有标签。因此，API应该允许有标签，但不支持配标签。

客户库必须允许在Gauge/Counter/Summary/Histogram创建时间可选地指定标签名称列表。客户端库应该支持任意数量的标签名称。客户端库必须验证标签名称符合已记录的要求。

提供对度量指标的标记维度访问的一般方式是通过`labels()`方法，该方法可以使用标签纸的列表或者从标签名称到标签纸的映射，并返回“child”，然后在Child上调用通常的`.inc()/.desc()/.observe()`等方法。

`label()`返回Child应该由用户缓存，以避免再次查找，这在延迟至关重要的代码中很重要。
带有标签的度量指标应该支持一个具有与`labels()`相同签名的`remove()`方法，它将从不再导出它的度量标准中删除一个Child，另一个clear()方法可以从度量指标中删除所有的`Child`。

应该有一种使用默认初始化给定Child的方法，通常只需要调用labels()。没有标签的度量指标必须被初始化，已避免缺少度量指标的问题。

#### 度量指标名称
度量指标名称补习遵循规范。与标签名称一样，必须满足使用`Counter/Gauge/Summary/Histogram`和库中提供的任意其他`Collector`的使用。

许多客户库提供三个部分的名称：`namespace_subsystem_name`, 其中只有该名称是强制性的。

必须不鼓励使用动态/生成的度量指标名称或者其子部分，除非自定义"Collector"是从其他工具/监控系统代理的。生成/动态度量指标名称是你应该使用标签的标志。

#### 度量指标描述和帮助
`Gauge/Counter/Summary/Histogram`要求必须提供度量指标的描述和帮助。

客户端中任何自定义的Collectors必须在度量指标名称上有一个描述和帮助。

建议将其作为强制性参数，但不要检查它是否具有一定长度，就好像有人真的不想写文档，否则我们不会说服他们。库中的Collector（以及我们再生态系统内部的任何地方）应该以良好的度量指标为例。

### 阐述

客户端必须实现一个文档阐述格式。

客户端可以实现多种格式。应该是可读性非常好的格式。

如果有疑问，请去文本格式。它不具有依赖性（protobuf），往往易于生成，是可读取的，并且protobuf的性能优势对于大多数用例来说并不重要。

如果可以在没有显著的资源成本情况下实现，可以重现可用的度量指标顺序（特别是对于人类可读格式）。

### 标准化和运行时收集器

客户端库应该提供标准导出的内容，如下所述：
这些应该作为自定义收集器实现，默认情况下在默认的CollectorRegistry上注册。应该有一种方法来禁用这些，因为有一些非常适用于他们的使用方式。

#### 处理度量指标
这些导出应该有前缀process_。如果一种语言或者运行时没有公开其中一个变量，它不会被导出它。所有内存值以字节为单位，以时间戳/秒为单位。

| 度量指标名称          |        含义                |    单位   |
| --------------------------| :-------------------------:| ---------:|
| process_cpu_seconds_total | 用户和系统CPU花费的时间    |  秒       |
| process_open_fds          | 打开的文件描述符数量       | 文件描述符|
| process_max_fds           | 打开描述符最大值           | 文件描述符|
| process_virtual_memory_bytes| 虚拟内存大小             | 字节|
| process_resident_memory_bytes| 驻留内存大小|字节|
| process_heap_bytes | 进程head堆大小| 字节|
| process_start_time_seconds| unix时间 | 秒|

### 运行时度量指标
另外，客户端库也被提供给他们的语言运行时（如：垃圾回收统计信息）的指标方面，提供了一些合适的前缀，比如： go_, hostspot_等。

### 单元测试
客户端库应该包含核心工具库和展示的单元测试。

客户端库被鼓励提供方便用户单元测试其使用的工具代码。例如，python中的CollectorRegistry.get_sample_value。

### 包和依赖
理想情况下，客户端库可以包含在任何应用程序中以添加一些工具，而无需担心它会破坏应用程序。

因此，当向客户端添加依赖关系时，建议谨慎。例如：如果用户添加使用添加版本1.4的Protobuf的Prometheus客户端库，但是应用程序在其他地方使用1.2，会发生什么？

建议在可能出现的情况下，将核心工具和给定格式的度量指标/展示分开。例如：Java简单客户端模块没有依赖关系，而simpleclient_servlet具有Http比特位。

### 性能考虑
由于客户端库必须是线程安全的，因此需要进行某种形式的并发控制，并且必须考虑多核机器和应用程序的性能。

在我们的经验中，性能最差的是互斥体。

处理器原子指令往往处于中间，并且通常可以接受。

避免不同CPU突然使用RAM的方法效果最好，例如：Java简单客户端中的DoubleAdder。有内存成本。

如上所述，labels()的结果应该是可缓存的。趋向于使用标签返回度量的并发映射往往相对较慢。没有标签的特殊套管指标，已避免labels()，像查找可以帮助很多。

度量指标应该避免阻塞，当它们递增/递减/设置等时，因为整个应用程序在持续获取时不会被组织。

主要工具操作的基准（包括labels）得到了鼓励。

应该牢记资源消耗，特别是RAM。考虑通过stream传输结果来减少内存占用，并且可能对并发获取的数量有限制。

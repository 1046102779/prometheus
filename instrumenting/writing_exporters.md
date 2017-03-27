## 编写导出器
---
当直接调试你自己的代码时，如何使用Prometheus客户端编写代码的通用规则能够直接进行。当从另一个监控系统中获取度量指标数据时，往往不会那么难。

这篇文档包含了当你写一个exporter或者自定义的collector时，你需要注意的一些点。所涉及的理论也很有意思。

如果你正在写一个exporter，且有任何不清楚的，可以随时联系我们, 详见IRC和[邮件](https://prometheus.io/community)

### 可维护性和Purity
当你写一个exporter时，你需要做出的最大决定是你会投入多大的代价，输出完美的度量指标

如果所讨论系统的度量指标很少变化，那么获取一切都是轻而易举的选择（例如：HAProxy exporter）。

另一方面, 如果系统具有数百种不断变化的度量指标，如果您尝试使其完美，那么你已经签署了许多正在进行的工作。一个例子：mysql exporter。

这个[node exporter](https://github.com/prometheus/node_exporter)是混合的，不同的模块，度量指标会有很大不同。对于mdadm，我们必须手动解析文件，并提出自己的度量指标，所以我们可以在需要的时候获取度量指标。对于meminfo，另一方面，结果在内核版本之间有所不同，所以我们最终只做了足够的变换来床架有效的度量。

### 配置
使用应用程序时，你应该针对不需要用户自定义配置的exporter，而不是告诉应用程序在哪里。你可能还需要提供过滤特定度量指标的功能，如果它们需要很高的代价才能获取度量指标数据（例如：HAProxy exporter允许过滤每个服务器的统计信息）。类似地，默认情况下可能会禁用代价比较高的度量指标

使用监控系统时，框架和协议不是那么简单的。

在最好的情况下，系统具有与Prometheus有相似的，且足够多的数据模型，你可以自己确定如何转换指标数据。Cloudwatch, SNMP和Collectd就是这样。最多我们需要能够让用户选择要获取的度量指标。

更多的用例是完全非标准的，依赖于用户如何使用它，以及底层的应用程序是什么。在这样的情况下，用户必须告诉我们如何转换数据。JMX exporter是最糟糕的，graphite和statsd exporter也要求从配置提取labels。

labelsexporter编写配置时，应牢记本文档。

YAML是Prometheus的标准配置格式

### Metrics度量指标
####  Naming命名
遵循[度量指标最佳实践](https://prometheus.io/docs/practices/naming)

通常度量指标名称应该允许一些人对Prometheus非常熟悉，而不是一个度量指标名称含义需要进行猜测的特殊系统。一个度量指标命名为`http_requests_total`并不是非常有用的，这些指标是在进入某些过滤器还是进入用户代码时进行测量？而`requests_total`甚至更糟，这是什么类型的请求？

为了以另一种方式使用工具，给定的度量指标应该存在于一个文件中。因此，在exporter和Collector中，度量指标适应于一个子系统，并相应地命名。

除非在编写自定义Collector或者exporter时，度量指标名称不应该被程序自动生成。

应用程序的度量指标名称应该以exporter名称为前缀，如haproxy_up。

度量指标名称必须使用基本单位（例如：秒，字节），并将其转换为图形软件更易读的内容。无论你最终使用什么样的计量单位，指标名称中的单位都必须与正在使用的单位相匹配。类似地展示比率，而不是百分比（尽管两个组件的计数器比率更好）。

度量指标名称不应该包括他们导出的labels（如：`by_type`）如果标签聚在一起，就没有意义。

这里有一个特例，当你通过多个度量指标导出不同标签的相同数据时，通常是区分他们最好的方式。这只有在导出具有所有标签的单个度量，具有太高的基数时才会出现。

Prometheus度量指标和标签名称被写在`snake_case`中。转换`camelCase`成`snake_case`是可取的，尽管这样做并不总是为`myTCPExample`或`isNaN`这样的事情产生不错的结果，所以有时最好将他们保留。

暴露的度量指标不应包含冒号，这些用于聚合时可供用户使用。度量指标名称只有符合正则表达式[a-zA-Z0-9:_]是有效的，任何其他字符串都需要被改成"_"下划线

`_sum`, `_count`, `_bucket`, 和`_total`用于Summaries，Histogram和Counters的后缀。除非生产其中之一，否则避免使用这些后缀。

`_total`是计数器的约定，如果使用COUNTER类型，则应该使用它。

这个`process_`和`scrape_`的前缀被保留。如果它们遵循匹配的语义，可以将自己的前缀添加到这些上。例如：Prometheus的持续获取时间为`scrape_duration_seconds`，这是很好的做法。

`jmx_scrape_duration_seconds`表示JMX Collector做到这一点花了多长时间。对于您可以访问pid的进程统计信息，Go和Python都会为你提供处理此选项的collector（请参阅HAProxy exporter示例）

当你统计请求的成功数量和失败数量时，最佳方法是暴露两个度量指标，一个是总的请求数，另一个度量指标是失败的请求度量指标。这使得计算失败比率变得很容易。不要使用带有failed/success的标签。类似于缓存中的`hit/miss`，有一个总的度量指标和另一个hits的度量指标是更好的。

考虑使用监控的人可能会对指标名称执行代码或者网络搜索。如果这些名称是非常完善的，并且不太可能在用于这些名称的人的领域之外（例如：SNMP和网络工程师）使用，那么按原样离开它们可能是一个好主意。该逻辑不适用于（例如：作为非数据库管理员的MySQL），可以预期会围绕这些指标。具有原始名称的帮助文档可以提供与使用以前的名称大致相同的好处。

#### labels
阅读一个有关labels标签的建议,详见[advice](https://prometheus.io/docs/practices/instrumentation/#things-to-watch-out-for)

避免把`type`做为标签名称，它太通用化，而且无意义。你也应该试着尽可能地避免与目标标签冲突，例如：`region`, `zone`, `cluster`, `availability`, `az`, `datacenter`, `dc`, `owner`, `customer`, `stage`, `environment`和`env`， 尽管这是应用程序调用的东西，但是最好不要通过重命名造成混乱。

避免这个，将东西放在一个度量指标中，因为它们共享一个前缀。除非你确定某些度量是由意义的，否则多个度量指标更加安全。

标签`le`对于Histograms有特殊的含义。同样对于Summaries来说，`quantile`也是如此。避免使用这些标签。

Read/write和send/receive作为独立的度量指标是非常好的，而不是作为一个标签。因为你关心的是某个时刻它们中的一个，度量指标命名并使用它们是容易的。

实践经验是，当求总和或平均时，一个度量指标应该是有意义的。 还有另一种情况出现了exporter，数据基本是用表格的方式呈现，否则将要求用户对度量标准名称进行正则表达式可用。 考虑您主板上的电压传感器，而对它们进行数学计算是无意义的，将它们置于一个度量标准中而不是每个传感器有一个度量是有意义的。 度量中的所有值（几乎）总是具有相同的单位（如果风扇速度与电压混合，则无法自动分离）。

不要做这些：
```
my_metric{label=a} 1
my_metric{label=b} 6
**my_metric{label=total} 7**
```

或者
```
my_metric{label=a} 1
my_metric{label=b} 6
**my_metric{} 7**
```

前者打破了对你的度量指标做`sum()`的人，后者打破了总和，也很难与之合作。某些客户端库（如：Go）会唧唧地尝试组织你在已定义收集器中执行后者，并且所有客户端都应组织你使用工具进行前者。不要这样做，依赖Prometheus的聚合操作，可以轻易地达到这一点。

如果你的监控暴露了一个像这样的total，请删除它。如果你必须因为一些原因（例如：这个total并不是统计计数）保留它，请使用其他的度量指标名称。

#### Target labels, not static scraped labels(目标标签，非静态获取标签)
如果你发现自己想要对所有度量指标应用相同的标签，请放弃这样做。

通常有两种情况出现。

第一个是一些标签对有关度量指标（例如软件版本号）是非常有用的。请使用文档[地址](https://www.robustperception.io/how-to-have-labels-for-machine-roles/)中所述的方法。

另一种情况是真正的目标标签。这些是区域，集群名称等，它们来自你的基础架构设置而不是应用程序本身。应用程序不应该在标签分类中说明，哪些地方适合Prometheus服务配置人员使用，而监控同一个应用程序的不同人员可能会给它不同的名称。

因此，这些标签通过你使用的任何服务发现，都属于Prometheus的获取配置。还可以在这里应用机器角色的概念，因为至少有一些人获取它，可能是非常有用的信息。

#### Types类型
你应该尝试将你的度量指标类型与Prometheus类型相匹配。这通常意味着Counter和Gauge。总结`_count`和`_sum`也是比较常见的，有时你会看到分位数。quantile柱状图是罕见的，如果你碰到了，记得展示格式并暴露累计值。

通常，衡量度量指标的类型是不明显的（特别是如果你自动处理一组指标），那么在这种情况下使用UNTYPED，一般来说，UNTYPED是一个安全默认值。

Counter不能下降，所以如果你有一个来自另一个测量系统的计数器类型，有一种方法来减少它（例如：Dropwizard指标），这不是一个计数器Counter而是一个计量表Gauge, UNTYPED可能实在那里使用的最好的类型，因为如果它被用作计数器，则GAUGE将会被误导。

### 帮助文档
当你转换度量指标时，用户能够跟踪原始内容记忆导致该转换的规则是有用的。以收件人/导出者的身份，应用程序的任何规则ID和原始度量的名称/详细信息记录到帮助文档，会极大地帮助用户。

Prometheus不喜欢一个具有不同帮助文档的度量指标。如果你从许多其他公司制定一个度量指标，请选择其中一个来放置帮助文档。

例如：SNMP导出器使用OID，JMX导出器放入一个样例mBean名称中。HAProxy exporter有手写字符串。node exporter有大量的例子可以使用。

#### 放弃无用的统计数据
某些测量系统暴露1m/5m/15m速率，从应用程序启动以来的平均速率（例如：在dropwizard指标中统计为平均值），最小值，最大值和标准偏差。

这些都应该被抛弃，因为它们不是非常有用, 并且增加了混乱。Prometheus可以自己计算费率，通常更准确（这些通常是指数衰减平均值）。你不知道什么时候计算最小/最大值，而stddev在统计上是无用的（如果你需要计算，则显示平方和，`_sum`和`_count`）

Quantiles有相关问题，你可以选择丢弃或者将其放在摘要中

#### .字符串(Dotted strings)
许多监控系统没有标签，而是做成像`my.class.path.mymetric.labelvalue1.labelvalue2.labelvalue3`这样。

graphite和statsd exporters分享一种使用小型配置语言执行此操作的方法。其他exporters也应该这样做。它目前仅在Go中实现，并将受益于将其分解为单独的库。

### Collectors
在为你的exporter实现一个collector时，你绝对不要使用通常最直接的测量方法，然后更新每个获取的度量指标。

相反，每次都会创建新的度量指标。在Go中，使用Update()方法中的MustNewConstMetric完成此操作。对于Python，请参与[https://github.com/prometheus/client_python#custom-collectors](https://github.com/prometheus/client_python#custom-collectors), 对于Java，在Collector方法中生成列表\<MetricFamilySamples\>, 请参与StardardExports.java作为示例。

这样做得原因首先是两次获取可能出现同时现象，直接使用有效的（文件级别）的全局变量，你会获取竞争条件。第二个原因是如果标签值消失，它仍然可以被导出。

通过直接测量来调整你的exporter是非常好的，例如：传输的总字节数或者exporter获取的所有度量指标数据中的操作。对于exporters而言，例如：黑盒exporters和snmp exporter，这个并不是单一的目标，所以这些只能在一个vanilla/metrics调用上公开，而不是在特定目标上。

#### 关于获取本身的度量指标
有时，你可以导出获取的度量指标数据，例如：你花费了多长时间或处理了多少记录。

这些应该被显示为Gauges（因为它们是关于event，scrape）和以exporter名字为前缀的度量指标名称，例如：`jmx_scrape_duration_seconds`。通常`_exporter`被排除（如果exporter仅仅作为collector使用，绝对排除它）。

### Machine & Process metrics（机器，进程度量指标）
许多系统（如：elasticsearch）暴露一些度量指标，如：cpu，内存和文件系统信息。当node exporter在Prometheus生态系统中提供了这些时，应该删除这些度量指标

在Java世界中，许多测量框架暴露了进程级别和JVM级别的统计信息，例如：CPU和GC。Java客户端和JMX exporter已经通过DefaultExports.java以首选方式包含了这些，因此这些应该被删除。

与其他语言类似。

### Deployment部署
每个exporter应该准确地监控一个实例应用程序，最好是在同一台机器上。这意味着你运行的每个HAProxy，你运行一个haproxy_exporter进程。对于具有mesos从节点的每台机器，都可以在其上运行mesos exporter（如果一台机器有两台机器，则运行另一台主机）。

这背后的理论是，对于测量系统来说，这是你正在做的，我们试图尽可能接近于其他布局。这意味着所有的服务发现都是在Prometheus而不是exporter完成的。这有利于Prometheus具有所需的目标信息，允许用户使用黑盒exporter探测你的服务。

有两个例外：

第一个实在应用程序旁边运行的见识是完全无意义的。SNMP，黑盒和IPMI是这一方的主要例子。IPMI和SNMP作为设备是有效的黑盒子，这是不可能运行代码（虽然如果你可以运行一个node exporter）和黑盒子，就像你正在监视像DNS名称这样没什么可以跑的在这种情况下，prometheus仍然应该做服务发现，并传递获取的目标。有关示例，请参阅blackbox和SNMP exporter）

请注意，目前只有使用python和Java客户端库编写这种类型的exporter（在Go中编写的blacxbox exporter手工执行文本格式，请不要这样做）。

另一个是你从系统的随机实例中提取一些统计信息的地方，而不在乎您正在谈论哪一个。考虑一组MySQL从机，您希望对数据进行一些业务查询，然后导出。有一个exporter使用你通常的负载平衡方法来谈话

### ::TODO 有时间再翻译，想吐

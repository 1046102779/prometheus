## 监控系统产品比较
### Prometheus vs. Graphite
#### 范围
Graphite专注于查询语言和图表特征的时间序列数据库。其他需求都需要依赖外部组件实现。

Prometheus是一个基于时间序列数据的完整监控系统和趋势系统，包括内置和主动抓取、存储、查询、图表和警告。它懂得监控系统和趋势系统应该是什么样的（哪些目标机应该存在，哪些时间序列模型存在问题等等），并积极地试着找出故障

#### 数据模型
Graphite和Prometheus一样，存储时间序列数值采样点。然而，Prometheus的元数据模型更加丰富：Graphite的度量指标名称是由隐式地编码多维度的点分隔命名的，而Prometheus的度量指标是可以自定义名称的，并以key-value键值对的标签，赋给度量指标。通过这些标签，我们使用Prometheus查询语言，可以方便地进行时间序列数据的过滤、分组、匹配操作。

进一步地，当Graphite与StatsD结合使用时，Graphite就只是对一个聚合数据的存储系统了，而不是把目标实例作为一个维度，并深入分析目标实例出现的各种问题。

例如：我们用Graphite/StatsD监控系统存储api-server服务Http请求的数量，前置条件包括：返回码是`500`，请求方法是`POST`，访问URL为`/tracks` , 度量指标自动隐式地编码，如下所示：
> stats.api-server.tracks.post.500 -> 93

但是在Prometheus中同样的数据存储可能像下面一样(假设有三个api-server)：
> api_server_http_requests_total{method="POST",handler="/tracks",status="500",instance="<sample1>"} -> 34
> api_server_http_requests_total{method="POST",handler="/tracks",status="500",instance="<sample2>"} -> 28
> api_server_http_requests_total{method="POST",handler="/tracks",status="500",instance="<sample3>"} -> 31

由上可以看到，三个api-server各自的度量指标数据，Prometheus把api-server也作为了一个维度，便于分析api-server服务出现的各种问题

#### 存储
Graphite存储以[Whisper](http://graphite.readthedocs.org/en/latest/whisper.html)把时间序列数据存储到本地磁盘，这种数据存储格式是RRD风格的数据库，它期望采样点定期地到达。 任何时间序列在一个单独的文件中存储，一段时间后新采集的样本会覆盖老数据

Prometheus也为每一个时间序列创建了一个本地文件，但是它允许时间序列以任意时间到达。新采集的样本被简单地追加到文件尾部，老数据可以任意长的时间保留。Prometheus对于短生命周期、且经常变化的时间序列集也可以表现得很好

#### 总结

Prometheus提供了一个丰富的数据模型和查询语言，而且更加容易地运行和集成到你的环境中。如果你想要一个可以长期保留历史数据的集群解决方案，Graphite可能是一个更好的选择。

### Prometheus vs. InfluxDB
InfluxDB是一个开源的时间序列数据库，它的商业版本具有可扩展和集群化的特性。在Prometheus刚刚开始开发时，InfluxDB项目已经发布了近一年时间。但是这两款产品还是有很大的不同之处，这两个系统也有一些略有不同的应用小场景。

#### 范围
公平起见，我们必须把InfluxDB和Kapacitor结合起来，与Prometheus和Prometheus的报警管理工具比较。

Graphite与Prometheus的范围差异，同样适用于InfluxDB本身。此外InfluxDB提供了连续查询，和Prometheus的记录规则一样。

Kapacitor的作用范围相当于，Prometheus的记录规则、告警规则和警告通知功能的结合。Prometheus提供了一个更加丰富地用于图表化和警告的查询语言，Prometheus告警器还提供了分组、重复数据删除和静默功能(silencing functionality)。

#### 数据模型/存储
和Prometheus一样，InfluxDB数据模型采用的标签也是键值对形式，被称为tags。而且InfluxDB有第二级标签，被称为fields，它被更多地限制使用。InfluxDB支持高达纳秒级的时间戳，以及float64、int64、bool和string的数据类型。相反地，Prometheus仅仅支持float64的数据类型，strings和毫秒只能小范围地支持

InfluxDB使用变种的日志结构合并树结构来存储预写日志，并按时间分片。这比Prometheus的文件追加更适合事件记录
[Logs and Metrics and Graphs, Oh My](https://blog.raintank.io/logs-and-metrics-and-graphs-oh-my)描述了事件日志和度量指标记录的不同

#### 框架
Prometheus服务独立运行，没有集群架构，它仅仅依赖于本地存储。Prometheus有四个核心的功能：抓取、规则处理和警告。InfluxDB的开源版本也是类似的。
InfluxDB的商业版本具有存储和查询的分布式版本, 存储和查询由集群中的节点同时处理。
这意味着商业版本的InfluxDB更加容易的水平扩展，同时也表示你必须从一开始就要管理分布式存储系统的复杂性。而Prometheus运行非常简单，而且在某些时候，你需要在可扩展性边界（如产品，服务，数据中心或者类似方面）明确分片服务器。单Prometheus服务也可以为您提供更好的可靠性和故障隔离。

Kapacitor对规则、警告和通知当前还没有内置/冗余选项。相反地，通过运行Prometheus的冗余副本和使用警告管理器的高可用模式提供了冗余选项。Kapacitor通过用户手动水平切分能够被缩放，这点类似于Prometheus本身

#### 总结
在两个系统之间有许多相似点。1. 利用标签（tags/labels）有效地支持多维度量指标。2. 使用相同的压缩算法。3.都可扩展集成。4.允许使用第三方进行监控系统的扩展，如：统计分析工具、自动化操作

InfluxDB更好之处：
 - 使用事件日志
 - 商业版本提供的集群方案，对于长期的时间序列存储是非常不错的
 - 复制的数据最终一致性

Prometheus更好之处：
 - 主要做度量指标监控
 - 更强大的查询语言，警告和通知功能
 - 图表和警告的高可靠性和稳定性

InfluxDB是有一家商业公司按照开放核心模式运营，提供高级功能，如：集群是闭源的，托管和支持。Prometheus是一个完全开放和独立的项目，有许多公司和个人维护，其中也提供一些商业服务和支持。

### Prometheus vs. OpenTSDB
[OpenTSDB](http://opentsdb.net/)是一个基于hadoop和Hbase的分布式时间序列数据库

#### 范围
和Graphite vs. Prometheus的范围一样

#### 数据模型
OpenTSDB的数据模型几乎和Prometheus一样：时间序列由任意的tags键值对集合表示。所有的度量指标存放在一起，并限制度量指标的总数量大小。Prometheus和OpenTSDB有一些细微的差别，例如：Prometheus允许任意的标签字符，而OpenTSDB的tags命名有一定的限制.OpenTSDB缺乏灵活的查询语言支持，通过它提供的API只能简单地进行聚合和数学计算

#### 存储
OpenTSDB的存储由Hadoop和HBase实现的。这意味着水平扩展OpenTSDB是非常容易的，但是你必须接受集群的总体复杂性

Prometheus初始运行非常简单，但是一旦超过单个节点的容量，就需要进行水平切分服务操作

#### 总结
Prometheus提供了一个非常灵活且丰富的查询语言，能够支持更多的度量指标数量，组成整个监控系统的一部分。如果你对hadoop非常熟悉，并且对时间序列数据有长期的存储要求，OpenTSDB是一个不错的选择

### Prometheus vs. Nagios
[Nagios](https://www.nagios.org/)是一款产生于90s年代的监控系统

#### 范围
Nagios是基于脚本运行结果的警告系统，又称"运行结果检查"。有警告通知，但是没有分组、路由和重复数据删除功能。

Nagios有大量的插件。例如：perfData插件抓取数据后写入到时间序列数据库（Graphite）或者使用NRPE在远程计算机上运行检查

#### 数据模型
Nagios是基于主机的，每一台主机有一个或者多个服务。其中一个是check运行检查，但是没有标签和查询语言的概念

Nagios除了检查脚本运行状态，没有任何存储功能。有第三方插件可以存储数据，并可视化数据

#### 架构
Nagios是单实例服务，所有的检查配置项统一由一个文件配置。

#### 总结
Nagios对于小型监控或者黑盒测试时非常有效的。如果你想要做白盒监控，或者动态地，基于云环境的数据监控，Prometheus是一个不错的选择

### Prometheus vs. Sensu
广义上说，[Sensu](https://sensuapp.org/)是一个更加现代的Nagios。

#### 范围
主要不同点在于Sensu客户端注册自己，并确定从本地还是其他地方获取配置检查。Senus对perfData的数量没有限制。 还有一个客户端socket允许把任意检查结果推送到Senus

#### 数据模型
和Nagios一样

#### 存储
Sensu在Redis中存储数据，存储被称作stashes。主要是静默存储，同时它也存储在Senus上注册的所有客户端

#### 架构
Sensu有很多组件。它使用Rabbit消息队列进行数据传输，使用Redis存储当前状态，独立的服务处理数据

RabbitMQ和Redis都可以是集群的，运行多个服务器副本可是实现副本和冗余

#### 总结
如果已经有了Nagios服务，你希望扩展它，同时希望使用Senus的注册特性，那么Senus是一个不错的选择

如果你想要使用白盒、或者有一个动态的云环境，那么Prometheus是一个很好的选择。

## STORAGE 存储
---
Prometheus有一个复杂的本地存储子系统。对于索引，它使用levelDB。对于批量的样本数据，它由自己的自定义存储层，并以固定大小（1024个字节有效负载）的块组织样本数据。然后将这些块存储在每个时间序列的一个文件中的磁盘上。

### Memory usage 内存使用量
Prometheus将所有当前使用的块保留在内存中。此外，它将最新使用的块保留在内存中，最大内存可以通过`storage.local.memory-chunks`标志配置。如果你有较多的可用内存，你可能希望将其增加到默认值1048576字节以上（反之亦然，如果遇到RAM问题，可以尝试减少内存值）。请注意，服务器的实际RAM使用率将高于将`storage.local.memory-chunks`*1024字节所期望的RAM使用率。管理存储层中的样本数据是不可避免的开销。此外，服务器正在做更多的事情，而不仅仅存储样本数据。实际开销取决于你的使用模式。在极端情况下，Prometheus必须保持更多的内存块，而不是配置，因为所有这些块都在同一时间使用。你必须试一下。Prometheus导出导出的度量指标`prometheus_local_storage_memory_chunks`和`process_resident_memory_bytes`将派上用场。作为经验法则，你应该至少拥有内存块所需三倍以上。

设计到大量时间序列的PromQL查询大量使用LevelDB支持的索引。如果需要运行这种查询，则可能需要调整索引缓存大小。以下标志是相关的：
 - `-storage.local.index-cache-size.label-name-to-label-values`: 正则表达式匹配
 - `-storage.local.index-cache-size.label-pair-to-fingerprints`: 如果大量的时间序列共享相同的标签，增加内存大小
 - `-storage.local.index-cache-size.fingerprint-to-metric` and `-storage.local.index-cache-size.fingerprint-to-timerange`: 如果你有大量的目标时间序列，例如：一段时间还没有被接收的样本数据时间序列，但是数据又还没有失效。这时也需要增加内存. 

你必须尝试使用flag，才能找出有用的。如果一个查询触及到100000多个时间序列数据，几百M内存使用可能是合理的。如果你有足够的内存可用，对于LevelDB使用更多的内存不会有任何伤害。

### 磁盘使用量 disk usage
Prometheus存储时间序列在磁盘上，目录由flag `storage.local.path`指定。默认path是`./data`（关联到工作目录），这是很好的快速尝试，但很可能不是你想要的实际操作。这个flag`stroage.local.retention`允许你配置这个保留的样本数据。根据你的需求和你的可用磁盘空间做出合适的调整。

### Chunking encoding
Prometheus当前提供三种不同类型的块编码（chunk encodings）。对于新创建块的编码由flag `-storage.local.chunk-encoding-version`决定。 有效值分别是0，1和2.

对于Prometheus的第一块存储存，类型值为0实现了delta编码。类型值为1是当前默认编码, 这是有更好的压缩算法的双delta编码，比类型值为0的delta编码要好。这两种编码在整个块中都具有固定的每个样本字节宽度，这允许快速随机访问。然而类型值为0 的delta编码是最快的编码，与类型值为1的编码相比，编码成本的差异很小。由于具有更好的压缩算法的编码1，除了兼容Prometheus更老的版本，一般建议使用编码1。

类型2是可变宽度的编码，例如：在块中的每个样本能够使用一个不同数量的bit位数。时间戳也是双delta编码。但是算法稍微有点不同。一些不同编码范式对于样本值都是可用的。根据样本值类型来决定使用哪种编码范式，样本值类型有：constant，int型，递增，随机等

编码2的主要部分的灵感来源于Facebook工程师发表的一篇论文：[Gorilla: A Fast, Scalable, In-Memory Time Series Database](http://www.vldb.org/pvldb/vol8/p1816-teller.pdf)

编码2必须顺序地访问块，并且编解码的代价比较高。总体来看，对比编码1，编码2造成了更高的CPU使用量和增加了查询延时，但是它提供一个改进的压缩比。准确值非常依赖于数据集和查询类型。下面的结果来自典型的生产环境的服务器中：
| 块编码类型 | 每个样本数据占用的比特位数  |  CPU核数 | 规则评估时间 | 
| ---------- | -------------------------:  |  -------:| :----------: |
|     1      |           3.3               |     1.6  |      2.9s    |
|     2      |           1.3               |     2.4  |      4.9s    |


每次启动Prometheus服务时，你可以改变块的编码类型，因此在实验中测试不同编码类型是我们非常鼓励的。但是考虑到，仅仅是新创建的块会使用新选择块编码，因此你将需要一段时间才能看到效果。

### 设置大量时间序列数据
Prometheus能够处理百万级别的时间序列数据。然而，你必须调整存储设置到处理多余100000活跃的时间序列。基本上，对于每个时间序列要存储到内存中，你想要允许这几个确定数量的块。对于`storage.local.memory-chunks`flag标志的默认值是1048567。高达300000个时间序列时，平均来看，每个时间序列仍然有三个可用的块。对于更多的时间序列，应该增加`storage.local.memory-chunks`值。三倍于时间序列的数量是一个非常好的近似值。但请注意内存使用的含义（见上文）。

如果你比配置的内存块有更多的时间序列数据，Prometheus不可避免地遇到一种情况，它必须保持比配置更多的内存块。如果使用块数量超过配置限制的10%， Prometheus将会减少获取的样本数据量（通过skip scrape和rule evaluation）直到减少到超过配置的5%。减少获取样本数量是非常糟糕的情况，这是你我都不愿意看到的。

同样重要地，特别是如果写入磁盘，会增长`storage.local.max-chunks-to-persist`flag值。根据经验，保持它是`storage.local.memory-chunks`值的50%是比较好的。`storage.local.max-chunks-to-persist`控制了多少块等待写入到你的存储设备，它既可以是spinning磁盘，也可以是SSD。如果等待块过多，这Prometheus将会减少获取样本数量，知道等待写入的样本数量下降到配置值的95%以下。在发生这种情况之前，Prometheus试着加速写入块。详见[文档](https://prometheus.io/docs/operating/storage/#persistence-pressure-and-rushed-mode)

每个时间序列可以保留更多的内存块，你就可以批量编写更多的写操作。对spinning磁盘是相当重要的。注意每个活跃的时间序列将会有个不完整的头块，目前还不能被持久化。它是在内存的块，不是磁盘块数据。如果你有1M的活跃时间序列数据，你需要3M`storage.local.memory-chunks`块，为每个时间序列提供可用的3块内存。仅仅有2M可持久化，因此设置`storage.local.max-to-persist`值大于2M，可以很容易地让内存超过3M块。尽管存储`storage.local.memory-chunks`的设置，这再次导致可怕的减少样本数量（Prometheus服务将尽快再此出现之前加速消费）。

等待持久性块的高价值的另一个缺点是检查点较大。

如果将大量时间序列与非常快速和/或较大的scrapes相结合，则预先分配的时间序列互斥所可能不会很奏效。如果你在Prometheus服务正在编写检查点或者处理代价大的查询时，看到获取较慢，请尝试增加`storage.local.num-fingerprint-mutexes`flag值。有时需要数万甚至更多。

### 持续压力和“冲动模式” Persist pressure and “rushed mode”
本质上，Prometheus服务将尽可能快递将完成的块持久化到磁盘上。这样的策略可能会导致许多小的写入操作，会占用更多的I/O带宽并保持服务器的繁忙。spinning磁盘在这里更加敏感，但是即使是SSD也不会喜欢这样。Prometheus试图尽可能的批量编写写操作，如果允许使用更多的内存，这样做法更好。因此，将上述flag设置为导致充分利用可用内存的值对于高性能非常重要。

Prometheus还将在每次写入后同步时间序列文件（使用`storage.local.series-sync-strategy = adaptive`, 这是默认值）， 并将磁盘带宽用于更频繁的检查点（根据“脏的时间序列”的计数，见下文），都试图在崩溃的情况下最小化数据丢失。

但是，如果等待写入的块数量增长太多，怎么办？Prometheus计算一个持久块的紧急度分数，这取决于等待与`storage.local.max-chunks-to-persist`值相关的持久性的快数量，以及内存中的快数量超过存储空间的数量。`local.memory-chunks`值（如果有的话，只有等待持久性的块的最小数量，以便更快的帮助）。分数在0~1.其中1是指对应于最高的紧急程度。根据得分，Prometheus将更频繁地写入磁盘。如果得分超过0.8的门槛，Prometheus将进入“冲动模式”（你可以在日志中国看到）。在冲动模式下，采用以下策略来加速持久化块：
  - 时间序列文件不再在写操作之后同步（更好地利用操作系统的页面缓存，在服务器崩溃的情况下，丢失数据的风险会增加）， 这个行为通过`storage.local.series-sync-strategy`flag。
 - 检查点仅仅通过`storage.local.checkpoint-interval`flag启动配置时创建（对于持久化块，以崩溃的情况下更多丢失数据的代码和运行随后崩溃恢复的时间增加，来释放更多的磁盘带宽）
 - 对于持久化块的写操作不再被限制，并且尽可能快地执行。

一段得分下降到0.7以下，Prometheus将退出冲动模式。

### 设置更长的保留时间 setting for very long retention time
如果你有通过`storage.local.retention`flag(超过一个月), 设置一个更长的留存时间，你可能想要增加`storage.local.series-file-shrink-ratio`flag值。

每当Prometheus需要从一系列文件的开头切断一些块时，它将简单地重写整个文件。（某些文件系统支持“头截断”，Prometheus目前由于几个原因目前不使用）。为了不重写一个非常大的系列文件来摆脱很少的块，重写只会发生在至少10％的块中 系列文件被删除。 该值可以通过上述的`storage.local.series-file-shrink-ratio`flag来更改。 如果您有很多磁盘空间，但希望最小化重写（以浪费磁盘空间为代价），请将标志值增加到更高的值，例如。 30％所需的块删除为0.3。

### 有用的度量指标
在Prometheus暴露自己的度量指标之外，以下内容对于调整上述flag特别有用：
 - prometheus_local_storage_memory_series: 时间序列持有的内存当前块数量
 - prometheus_local_storage_memory_chunks: 在内存中持久块的当前数量
 - `prometheus_local_storage_chunks_to_persist`: 当前仍然需要持久化到磁盘的的内存块数量
 - `prometheus_local_storage_persistence_urgency_score`: 上述讨论的紧急程度分数
 - 如果Prometheus处于冲动模式下，`prometheus_local_storage_rushed_mode`值等于1; 否则等于0.

### Crash恢复 Carsh Recovery
Prometheus在完成后尽快将块保存到磁盘。常规检查点中不完整的块保存到磁盘。您可以使用`storage.local.checkpoint-interval`flag配置检查点间隔。如果太多的时间序列处于“脏”状态，那么Prometheus更频繁地创建检查点，即它们当前的不完整的头部块不是包含在最近检查点中的。此限制可通过`storage.local.checkpoint-dirty-series-limit`flag进行配置。

然而，如果您的服务器崩溃，您可能仍然丢失数据，并且您的存储空间可能处于不一致的状态。因此，Prometheus在意外关机后执行崩溃恢复，类似于文件系统的fsck运行。将记录关于崩溃恢复的详细信息，因此如果需要，您可以将其用于取证。无法恢复的数据被移动到名为孤立的目录（位于storage.local.path下）。如果不再需要，请记住删除该数据。

崩溃恢复通常需要不到一分钟。如果需要更长时间，请咨询日志，以了解出现的问题。

### Data corrution 数据损坏
如果您怀疑数据库中的损坏引起的问题，则可以通过使用`storage.local.dirty`flag启动服务器来强制执行崩溃恢复。

如果没有帮助，或者如果您只想删除现有的数据库，可以通过删除存储目录的内容轻松地启动：
 1. stop prometheus.
 2. `rm -r <storage path>/*`
 3. start prometheus

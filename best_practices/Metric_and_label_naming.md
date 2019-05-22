## 度量指标和标签命名
---
使用Prometheus不需要本文档中提供的度量标准和标签约定，但可以同时用作样式指南和最佳实践集合。 个别组织可能希望采用其中一些做法，例如： 命名约定，不同。

### 指标名字
指标名称......

- ... 必须符合有效字符的[数据模型](https://prometheus.io/docs/concepts/data_model/#metric-names-and-labels)。 
- ... 应该具有与度量所属的域相关的（单字）应用程序前缀。前缀有时被客户端库称为命名空间。对于特定于应用程序的度量标准，前缀通常是应用程序名称本身。但是，有时候，指标更通用，就像客户端库导出的标准化指标一样。例子：
     - `prometheus_notifications_total`(特定于Prometheus服务)
     - `process_cpu_seconds_total`(由许多客户库导出)
     - `http_request_duration_seconds`(适用于所有HTTP请求)
- ...必须有一个单位（即不要与毫秒混合秒，或与字节混合秒）.
- ...应该有一个基本单位（例如：秒，字节，米，不是毫秒，兆字节，公里）。请参阅下面的基本单位列表。
- ...应该有一个以复数形式描述单位的后缀。请注意，累积计数除了单位（如果适用）外，还有总数作为后缀。
     - `http_request_request_seconds`
     - `node_memory_usage_bytes`
     - `http_requests_total`(对于无单位累计计数)
     - `process_cpu_seconds_total`(用于累计计数单位)
- ...应该代表所有标签尺寸上相同的逻辑事物
     - 请求持续时间
     - 数据传输的字节数
     - 瞬时资源使用率百分比

根据经验，给定度量的所有维度上的`sum()`或`avg()`应该是有意义的（尽管不一定有用）。如果没有意义，请将数据拆分为多个指标。例如，在一个度量中具有各种队列的容量是好的，而将队列的容量与队列中的当前数量的元素混合则不是。

### 标签
使用标签来区分正在测量的事物特征：

- `api_http_requests_total`- 区分请求类型： `type="created| update | delete`.
- `api_request_duration_seconds`- 区分请求阶段：`stage="extract | transform | load"` 

不要将标签名称放在度量标准名称中，因为这会引入冗余，如果聚合了相应的标签，则会引起混淆。

### 基本单位
Prometheus没有硬编码的任何单位。 为了更好的兼容性，应使用基本单元。 以下列出了一些具有基本单位的度量标准系列。 该清单并非详尽无遗。

|Family|基本单位 | 备注|
|---|---|---|
|Time | seconds |  | 
| Temperature | celsius | 摄氏度是实践中遇到的最常见的摄氏度 | 
| Length | meters | |
| Bytes | bytes | |
| Bits | bytes | | 
| Percent | ratio(*) | 值为0-1。`*`）通常'ratio'不用作后缀，而是用作A_per_B。 例外是例如disk_usage_ratio | 
| Voltage | volts | |
| Electric current | amperes | |
| Energy | joules | | 
| Weight | grams | | 

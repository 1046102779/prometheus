## Jobs和Instances(任务和实例)
---
就Prometheus而言，任何抓取的目标都被称作*instance*，通常是一个服务进程。相同类型的实例集合被称为一个任务*job*。

例如, 一个被称作*api-server*的任务有四个相同的实例。
 - 任务: `api-server`
     1. 实例1: `1.2.3.4:5670`
     2. 实例2：`1.2.3.4:5671`
     3. 实例3：`5.6.7.8:5670`
     4. 实例4：`5.6.7.8:5671`

### 自动化生成的标签和时间序列
当Prometheus抓取一个目标数据时，它会把一些标签自动化地赋予给被抓取的时间序列数据：
  - `job`: 目标所属于的配置任务名称。
  - `instance`: 被抓取的目标服务`host:port`

判断任何一个标签是否在抓取的时间序列数据中，取决于`honor_labels`配置选项。详见[文档](https://prometheus.io/docs/operating/configuration/#%3Cscrape_config%3E)

对于每个实例，在下面的时间序列数据中，Prometheus存储了一个样本：
 - up{job="[job-name]", instance="instance-id"}: 如果实例是健康的，则up值等于1，否则，up值等于0，表示目标服务不可达。
 - scrape_duration_seconds{job="[job-name]", instance="[instance-id]"}: 抓取目标数据的持久性

`up`度量指标对目标的健康监控是非常有用的。

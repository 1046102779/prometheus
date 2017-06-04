## Jobs和Instances(任务和实例)
---
就Prometheus而言，任何抓取的进程都被称作*instance*。相同的多实例进程集合被称为一个任务*job*。

例如, 一个被称作*api-server*的任务有四个相同的实例。
 - 任务: `api-server`
     - 实例1：`1.2.3.4:5670`
     - 实例2：`1.2.3.4:5671`
     - 实例3：`5.6.7.8:5670`
     - 实例4：`5.6.7.8:5671`

### 自动化生成的标签和时间序列
当Prometheus抓取一个进程的度量指标数据时，默认会有一些度量指标存在。
  - `job`: 目标所属于的配置任务名称。
  - `instance`: 被抓取的目标服务`host:port`

判断任何一个标签是否在抓取的时间序列数据中，取决于`honor_labels`配置选项。详见[文档](https://prometheus.io/docs/operating/configuration/#%3Cscrape_config%3E)

对于每个进程，Prometheus都会默认为它创建一些度量指标：
 - up{job="[job-name]", instance="instance-id"}: 如果进程是健康的，则up值等于1，否则，up值等于0，表示进程不可用。
 - scrape_duration_seconds{job="[job-name]", instance="[instance-id]"}: 表示抓取一次度量指标数据花费的时间。
 - scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}: 表示度量指标的标签变化后，标签没有变化的度量指标数量。
 - scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}: 进程的所有度量指标总数

`up`度量指标对进程健康的监控是非常有用的。

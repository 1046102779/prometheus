就Prometheus而言，pull拉取采样点的端点服务称之为**instance**，通常对应一个过程。具有相同目的的**instance**，例如，为可伸缩性或可靠性而复制的流程称为作业。, 则构成了一个**job**

例如, 一个被称作**api-server**的任务有四个相同的实例。
 - job: `api-server`
     - instance 1：`1.2.3.4:5670`
     - instance 2：`1.2.3.4:5671`
     - instance 3：`5.6.7.8:5670`
     - instance 4：`5.6.7.8:5671`

##### 一、自动化生成的标签和时间序列
当Prometheus拉取一个目标，会自动地把两个标签添加到度量名称的标签列表中，分别是：
  - **job**: 目标所属的配置任务名称。
  - **instance**: 被抓取的目标网址的一部分务: `host:port`

如果以上两个标签二者之一存在于采样点中，这个取决于`honor_labels`配置选项。详见[文档](https://prometheus.io/docs/operating/configuration/#%3Cscrape_config%3E)

对于每个采样点所在服务instance，Prometheus都会存储以下的度量指标采样点：
 - `up{job="[job-name]", instance="instance-id"}`:`1`，表示采样点所在服务健康;`0`，标识抓取失败
 - `scrape_duration_seconds{job="[job-name]", instance="[instance-id]"}`: 抓取的持续时间
 - `scrape_samples_post_metric_relabeling{job="<job-name>", instance="<instance-id>"}`: 应用度量标准重新标记后剩余的样本数。
 - `scrape_samples_scraped{job="<job-name>", instance="<instance-id>"}`: 目标暴露的样本数量。

`up`度量指标对服务健康的监控是非常有用的。

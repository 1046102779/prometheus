本教程是类似"hello,world"的教程，展示怎样在一个简单地例子中安装、配置和使用Prometheus。你将下载和本地化运行Prometheus服务，并写一个配置文件，监控Prometheus服务本身和一个简单的应用，然后配合使用query、rules和graphs展示收集的时间序列数据。

##### 一、下载和运行Prometheus
[现在Prometheus最新的发布版本](https://prometheus.io/download),然后提取和运行它：
```shell
tar zxvf prometheus-*.tar.gz
cd prometheus-*
```
在开始启动Prometheus之前，我们要配置它

##### 二、配置Prometheus监控自身
Prometheus从监控的目标上通过http方式拉取指标数据,它也可以拉取自身服务数据并监控自身的健康状况。

当然Prometheus服务拉取自身服务数据，并没有多大的用处，但是它是一个好的开始例子。保存下面的基本Prometheus配置，并命名为：`prometheus.yml`:
```
global:
  scrape_interval:     15s # 默认情况下，每15s拉取一次目标采样点数据。

  # 我们可以附加一些指定标签到采样点度量标签列表中, 用于和第三方系统进行通信, 包括：federation, remote storage, Alertmanager
  external_labels:
    monitor: 'codelab-monitor'

# 下面就是拉取自身服务数据配置
scrape_configs:
  # job名称会增加到拉取到的所有采样点上，同时还有一个instance目标服务的host：port标签也会增加到采样点上
  - job_name: 'prometheus'

    # 覆盖global的采样点，拉取时间间隔5s
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']
```

对于一个完整的配置选项，请见[配置文档](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)

##### 三、启动Prometheus
要使用新创建的配置文件启动Prometheus，请切换到包含Prometheus二进制文件的目录并运行：
```shell
./prometheus --config.file=prometheus.yml
```
Prometheus服务应该启动了。你可以在浏览器上输入：`http://localhost:9090`, 给它几秒钟从自己的HTTP指标端点收集有关自身的数据。

您还可以通过导航到其指标端点来验证Prometheus是否正在提供有关自身的指标：`http://localhost:9090/metrics`

##### 四、使用expression browser
让我们试着看一下Prometheus收集的关于自己的一些数据。 使用Prometheus的内置表达式浏览器，导航到`http://localhost:9090/graph`，并选择带有"Graph"的"Console".

在`http://localhost:9090/gmetrics`中收集中，有一个metric叫`prometheus_target_interval_length_seconds`(从目标收集数据的实际时间量)，在表达式的console中输入:
```shell
prometheus_target_interval_length_seconds
```
这个应该会返回很多不同的时间序列数据(以及每个记录的最新值)，这些度量名称都是`prometheus_target_interval_length_seconds`，但是带有不同的标签列表值，这些标签列表值指定了不同的延迟百分比和目标组间隔。

如果我们仅仅对99%的延迟感兴趣，则我们可以使用下面的查询去清洗信息：
```shell
prometheus_target_interval_length_seconds{quantile="0.99"}
```
为了统计返回时间序列数据个数，你可以写：
```shell
count(prometheus_target_interval_length_seconds)
```

有关更多的表达式语言，请见[表达式语言文档](https://prometheus.io/docs/prometheus/latest/querying/basics/)

##### 五、使用graph interface
见图表表达式，导航到`http://localhost:9090/graph`， 然后使用"Graph" tab

例如，输入以下表达式来绘制在自我抓取的Prometheus中创建的每秒块速率：
```shell
rate(prometheus_tsdb_head_chunks_created_total[1m])
```
试验graph范围参数和其他设置。

##### 六、启动其他一些采样目标
让我们让这个更有趣，并开始一些示例目标，让Prometheus抓取。

Go客户端库包含一个示例，该示例为具有不同延迟分布的三个服务导出虚构的RPC延迟。

确保已安装Go编译器并设置了正常工作的Go构建环境（具有正确的GOPATH）。

下载Prometheus的Go客户端，运行三个服务:
```shell
git clone https://github.com/prometheus/client_golang.git
cd client_golang/examples/random
go get -d
go build

## 启动三个服务
./random -listen-address=:8080
./random -listen-address=:8081
./random -listen-address=:8082
```
现在你在浏览器输入:`http://localhost:8080/metrics`, `http://localhost:8081/metrics`, `http://localhost:8082/metrics`, 能看到所有采集到的采样点数据。

##### 七、配置Prometheus去监控这三个目标服务
现在我们将会配置Prometheus，拉取三个目标服务的采样点。我们把这三个目标服务组成一个job, 叫`example-radom`。 然而，想象成，前两个服务是生产环境服务，后者是测试环境服务。我们可以通过group标签分组，要在Prometheus中对此进行建模，我们可以将多组端点添加到单个作业中，为每组目标添加额外的标签。在此示例中，我们将`group ="production"`标签添加到第一组目标，同时将`group ="canary"`添加到第二组。

要实现此目的，请将以下作业定义添加到prometheus.yml中的scrape_configs部分，然后重新启动Prometheus实例：
```shell
scrape_configs:
  - job_name:       'example-random'

    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:8080', 'localhost:8081']
        labels:
          group: 'production'

      - targets: ['localhost:8082']
        labels:
          group: 'test'
```

转到表达式浏览器并验证Prometheus现在是否有关于这些示例端点公开的时间序列的信息，例如`rpc_durations_seconds`指标。

##### 八、为抓取的数据聚合配置规则
虽然在我们的示例中不是问题，但是在计算ad-hoc时，聚合了数千个时间序列的查询会变慢。 为了提高效率，Prometheus允许您通过配置的录制规则将表达式预先记录到全新的持久时间序列中。 假设我们感兴趣的是记录在5分钟窗口内测量的所有实例（但保留作业和服务维度）的平均示例RPC（`rpc_durations_seconds_count`）的每秒速率。 我们可以这样写：
```shell
avg(rate(rpc_durations_seconds_count[5m])) by (job, service)
```
要将此表达式生成的时间序列记录到名为`job_service：rpc_durations_seconds_count：avg_rate5m`的新度量标准中，请使用以下记录规则创建一个文件并将其另存为`prometheus.rules.yml`：
```shell
groups:
- name: example
  rules:
  - record: job_service:rpc_durations_seconds_count:avg_rate5m
    expr: avg(rate(rpc_durations_seconds_count[5m])) by (job, service)
```

要使Prometheus选择此新规则，请在`prometheus.yml`中添加`rule_files`语句。 配置现在应该如下所示：
```shell
global:
  scrape_interval:     15s # By default, scrape targets every 15 seconds.
  evaluation_interval: 15s # Evaluate rules every 15 seconds.

  # Attach these extra labels to all timeseries collected by this Prometheus instance.
  external_labels:
    monitor: 'codelab-monitor'

rule_files:
  - 'prometheus.rules.yml'

scrape_configs:
  - job_name: 'prometheus'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']

  - job_name:       'example-random'

    # Override the global default and scrape targets from this job every 5 seconds.
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:8080', 'localhost:8081']
        labels:
          group: 'production'

      - targets: ['localhost:8082']
        labels:
          group: 'canary'
```
使用新配置重新启动Prometheus，并通过表达式浏览器查询或绘制图表，验证带有度量标准名称`job_service：rpc_durations_seconds_count：avg_rate5m`的新时间序列现在可用。


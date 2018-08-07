## 启动
---
这是个类似"hello,world"的试验，教大家怎样快速安装、配置和简单地搭建一个DEMO。你会下载和本地化运行Prometheus服务，并写一个配置文件，监控Prometheus服务本身和一个简单的应用，然后配合使用query、rules和图表展示采样点数据

### 下载和运行Prometheus
[最新下载页](https://prometheus.io/download), 然后提取和运行它，so easy：
```shell
tar zxvf prometheus-*.tar.gz
cd prometheus-*
```
在开始启动Prometheus之前，我们要配置它

### 配置Prometheus监控自身
Prometheus从目标机上通过http方式拉取采样点数据, 它也可以拉取自身服务数据并监控自身的健康状况

当然Prometheus服务拉取自身服务采样数据，并没有多大的用处，但是它是一个好的DEMO。保存下面的Prometheus配置，并命名为：`prometheus.yml`:
```shell
global:
  scrape_interval:     15s # 默认情况下，每15s拉取一次目标采样点数据。

  # 我们可以附加一些指定标签到采样点度量标签列表中, 用于和第三方系统进行通信, 包括：federation, remote storage, Alertmanager
  external_labels:
    monitor: 'codelab-monitor'

# 下面就是拉取自身服务采样点数据配置
scrape_configs:
  # job名称会增加到拉取到的所有采样点上，同时还有一个instance目标服务的host：port标签也会增加到采样点上
  - job_name: 'prometheus'

    # 覆盖global的采样点，拉取时间间隔5s
    scrape_interval: 5s

    static_configs:
      - targets: ['localhost:9090']
```

对于一个完整的配置选项，请见[配置文档](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)

### 启动Prometheus
指定启动Prometheus的配置文件，然后运行
```shell
./prometheus --config.file=prometheus.yml
```

这样Prometheus服务应该起来了。你可以在浏览器上输入：`http://localhost:9090`, 就可以看到Prometheus的监控界面

你也可以通过输入`http://localhost:9090/metrics`，直接拉取到所有最新的采样点数据集

### 使用expression browser(暂翻译：浏览器上输入表达式)
为了使用Prometheus内置浏览器表达式，导航到`http://localhost:9090/graph`，并选择带有"Graph"的"Console".

在拉取到的度量采样点数据中， 有一个metric叫`prometheus_target_interval_length_seconds`, 两次拉取实际的时间间隔，在表达式的console中输入:
```shell
prometheus_target_interval_length_seconds
```

这个应该会返回很多不同的倒排时间序列数据，这些度量名称都是`prometheus_target_interval_length_seconds`, 但是带有不同的标签列表值，这些标签列表值指定了不同的延迟百分比和目标组间隔

如果我们仅仅对99%的延迟感兴趣，则我们可以使用下面的查询去清洗信息：
```shell
prometheus_target_interval_length_seconds{quantile="0.99"}
```

为了统计返回时间序列数据个数，你可以写：
```shell
count(prometheus_target_interval_length_seconds)
```

有关更多的表达式语言，请见[表达式语言文档](https://prometheus.io/docs/prometheus/latest/querying/basics/)

### 使用graph interface
见图表表达式，导航到`http://localhost:9090/graph`， 然后使用"Graph" tab

例如，进入下面表达式，绘图最近1分钟产生chunks的速率：
```shell
rate(prometheus_tsdb_head_chunks_created_total[1m])
```

### 启动其他一些采样目标
Go客户端包括了一个例子，三个服务只见的RPC调用延迟

首先你必须有Go的开发环境，然后才能跑下面的DEMO, 下载Prometheus的Go客户端，运行三个服务:
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
现在你在浏览器输入:`http://localhost:8080/metrics`, `http://localhost:8081/metrics`, `http://localhost:8082/metrics`, 能看到所有采集到的采样点数据

### 配置Prometheus去监控这三个目标服务
现在我们将会配置Prometheus，拉取三个目标服务的采样点。我们把这三个目标服务组成一个job, 叫`example-radom`. 然而，想象成，前两个服务是生产环境服务，后者是测试环境服务。我们可以通过group标签分组，在这个例子中，我们通过`group="production"`标签和`group="test"`来区分生产和测试
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

进入浏览器，输入`rpc_duration_seconds`, 验证Prometheus所拉取到的采样点中每个点都有group标签，且这个标签只有两个值`production`, `test`

### 聚集到的采样点数据配置规则
上面的例子没有什么问题， 但是当采样点海量时，计算成了瓶颈。查询、聚合成千上万的采样点变得越来越慢。为了提高性能，Prometheus允许你通过配置文件设置规则，对表达式预先记录为全新的持续时间序列。让我们继续看RPCs的延迟速率(`rpc_durations_seconds_count`),  如果存在很多实例，我们只需要对特定的`job`和`service`进行时间窗口为5分钟的速率计算，我们可以写成这样：
```shell
avg(rate(rpc_durations_seconds_count[5m])) by (job, service)
```
为了记录这个计算结果，我们命名一个新的度量：`job_service:rpc_durations_seconds_count:avg_rate5m`, 创建一个记录规则文件，并保存为`prometheus.rules.yml`:
```shell
groups:
- name: example
  rules:
  - record: job_service:rpc_durations_seconds_count:avg_rate5m
    expr: avg(rate(rpc_durations_seconds_count[5m])) by (job, service)
```

然后再在Prometheus配置文件中，添加`rule_files`语句到`global`配置区域， 最后配置文件应该看起来是这样的：
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
          group: 'test'
```

然后重启Prometheus服务，并指定最新的配置文件，查询并验证`job_service:rpc_durations_seconds_count:avg_rate5m`度量指标

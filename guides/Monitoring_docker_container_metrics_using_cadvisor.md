## 使用CADVISOR监控DOCKER容器的指标
---
[cAdvisor](https://github.com/google/cadvisor)（容器顾问的简称）分析并公开运行容器的资源使用和性能数据。 cAdvisor开箱即用暴露Prometheus指标。 在本指南中，我们将：

- 创建一个本地多容器[Docker Compose](https://docs.docker.com/compose/)安装，其中包括分别运行Prometheus，cAdvisor和[Redis](https://redis.io/)服务器的容器
- 检查由Redis容器生成的一些容器指标，由cAdvisor收集并由Prometheus抓取

### Prometheus配置
首先，您需要配置Prometheus以从cAdvisor中获取指标。 创建一个prometheus.yml文件并使用以下配置填充它：
```
scrape_configs:
- job_name: cadvisor
  scrape_interval: 5s
  static_configs:
  - targets:
    - cadvisor:8080
```
### Docker Compose配置
现在我们需要创建一个Docker Compose[配置](https://docs.docker.com/compose/compose-file/)，指定哪些容器是我们安装的一部分，以及每个容器暴露的端口，使用哪些卷等等。

在您创建`prometheus.yml`文件的同一文件夹中，创建一个`docker-compose.yml`文件并使用此Docker Compose配置填充它：
```
version: '3.2'
services:
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    ports:
    - 9090:9090
    command:
    - --config.file=/etc/prometheus/prometheus.yml
    volumes:
    - ./prometheus.yml:/etc/prometheus/prometheus.yml:ro
    depends_on:
    - cadvisor
  cadvisor:
    image: google/cadvisor:latest
    container_name: cadvisor
    ports:
    - 8080:8080
    volumes:
    - /:/rootfs:ro
    - /var/run:/var/run:rw
    - /sys:/sys:ro
    - /var/lib/docker/:/var/lib/docker:ro
    depends_on:
    - redis
  redis:
    image: redis:latest
    container_name: redis
    ports:
    - 6379:6379
```

此配置指示Docker Compose运行三个服务，每个服务对应一个Docker容器：

1. `prometheus`服务使用本地`prometheus.yml`配置文件（由`volumes`参数导入容器）。
2. `cadvisor`服务公开端口8080（cAdvisor度量标准的默认端口）并依赖于各种本地卷（`/`，`/var/run`等）。
3. `redis`服务是标准的Redis服务器。 cAdvisor将自动从该容器收集容器指标，即无需进一步配置。

要运行安装：
```
docker-compose up
```
如果Docker Compose成功启动所有三个容器，您应该看到如下输出：
```
prometheus  | level=info ts=2018-07-12T22:02:40.5195272Z caller=main.go:500 msg="Server is ready to receive web requests."
```
您可以使用ps命令验证所有三个容器是否都在运行：
```
docker-compose ps
```
您的输出将如下所示：
```
   Name                 Command               State           Ports
----------------------------------------------------------------------------
cadvisor     /usr/bin/cadvisor -logtostderr   Up      8080/tcp
prometheus   /bin/prometheus --config.f ...   Up      0.0.0.0:9090->9090/tcp
redis        docker-entrypoint.sh redis ...   Up      0.0.0.0:6379->6379/tcp
```
### 探索cAdvisor Web UI
您可以访问位于`http://localhost:8080`的cAdvisor Web UI。 您可以在我们的安装位置`http://localhost:8080/docker/<container>`中浏览特定Docker容器的统计信息和图形。 例如，可以在`http://localhost:8080/docker/redis`，Prometheus在`http://localhost:8080/docker/prometheus`上访问Redis容器的度量标准，依此类推。

### 在表达式浏览器中探索指标
cAdvisor的Web UI是一个有用的界面，用于探索cAdvisor监控的各种事物，但它没有提供用于探索容器度量的界面。 为此，我们需要Prometheus表达式浏览器，它可以在`http://localhost:9090/graph`中找到。 您可以在表达式栏中输入Prometheus表达式，如下所示：

![prometheus-expression-bar](https://prometheus.io/assets/prometheus-expression-bar.png)

让我们首先探索`container_start_time_seconds`指标，该指标记录容器的开始时间（以秒为单位）。 您可以使用`name="<container_name>"`表达式按名称选择特定容器。 容器名称对应于Docker Compose配置中的`container_name`参数。 例如，`container_start_time_seconds{name="redis"}`表达式显示了`redis`容器的开始时间。

### 其他表达式
下表列出了一些其他示例表达式：

|表达式|描述|目的|
|---|---|---|
| `rate(container_cpu_usage_seconds_total{name="redis"}[1m])` | cgroup在最后一分钟的CPU使用率（按核心划分） | `redis`容器 | 
| `container_memory_usage_bytes{name="redis"}` | cgroup的总内存使用量（以字节为单位） | `redis`容器 |
| `rate(container_network_transmit_bytes_total[1m])` | 在最后一分钟，容器每秒通过网络传输的字节数 | 所有容器 | 
| `rate(container_network_receive_bytes_total[1m])` | 在最后一分钟，容器每秒接收网络传输的字节数 | 所有容器 | 

### 总结
在本指南中，我们使用Docker Compose在一个安装中运行了三个独立的容器：一个Prometheus容器从cAdvisor容器中抓取指标，该容器反过来收集由Redis容器生成的指标。 然后，我们使用Prometheus表达式浏览器探索了一些cAdvisor容器指标。
## 推送度量指标
---
偶尔你需要监控不能被获取的实例。它们可能被防火墙保护，或者它们生命周期太短而不能通过拉模式获取数据。Prometheus的Pushgateway允许你将这些实例的时间序列数据推送到Prometheus的代理任务中。结合Prometheus简单的文本导出格式，这使得即使没有客户库，也能使用shell脚本获取数据。
 - shell实现用例，查看[Readme](https://github.com/prometheus/pushgateway/blob/master/README.md)
 - Java, 详见[PushGateway](https://prometheus.io/client_java/io/prometheus/client/exporter/PushGateway.html)类
 - Go，详见[Push](http://godoc.org/github.com/prometheus/client_golang/prometheus#Push)和[PushAdd](http://godoc.org/github.com/prometheus/client_golang/prometheus#PushAdd)
 - Python, 详见[Pushgateway](https://github.com/prometheus/client_python#exporting-to-a-pushgateway)
 - Ruby, 详见[Pushgateway](https://github.com/prometheus/client_ruby#pushgateway)

### Java批量任务例子
这个例子主要说明, 如何执行一个批处理任务，并且在没有执行成功时报警

如果使用Maven，添加下面的代码到`pom.xml`文件中：
```Java
        <dependency>
            <groupId>io.prometheus</groupId>
            <artifactId>simpleclient</artifactId>
            <version>0.0.10</version>
        </dependency>
        <dependency>
            <groupId>io.prometheus</groupId>
            <artifactId>simpleclient_pushgateway</artifactId>
            <version>0.0.10</version>
        </dependency>
```

执行批量作业的代码：
```Java
import io.prometheus.client.CollectorRegistry;
import io.prometheus.client.Gauge;
import io.prometheus.client.exporter.PushGateway;

void executeBatchJob() throws Exception {
 CollectorRegistry registry = new CollectorRegistry();
 Gauge duration = Gauge.build()
     .name("my_batch_job_duration_seconds")
     .help("Duration of my batch job in seconds.")
     .register(registry);
 Gauge.Timer durationTimer = duration.startTimer();
 try {
   // Your code here.

   // This is only added to the registry after success,
   // so that a previous success in the Pushgateway is not overwritten on failure.
   Gauge lastSuccess = Gauge.build()
       .name("my_batch_job_last_success_unixtime")
       .help("Last time my batch job succeeded, in unixtime.")
       .register(registry);
   lastSuccess.setToCurrentTime();
 } finally {
   durationTimer.setDuration();
   PushGateway pg = new PushGateway("127.0.0.1:9091");
   pg.pushAdd(registry, "my_batch_job");
 }
}
```

警报一个Pushgateway，如果需要的话，修改host和port

如果任务最近没有运行，请创建一个警报到Alertmanager。将以下内容添加到Pushgateway的Prometheus服务的记录规则中：
```record rules
ALERT MyBatchJobNotCompleted
  IF min(time() - my_batch_job_last_success_unixtime{job="my_batch_job"}) > 60 * 60
  FOR 5m
  WITH { severity="page" }
  SUMMARY "MyBatchJob has not completed successfully in over an hour"
```

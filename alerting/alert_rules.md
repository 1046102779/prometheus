## Alerting rules 警报规则
---
警报规则允许你基于Prometheus表达式语言的表达式定义报警报条件，并在触发警报时发送通知给外部的接收者。每当警报表达式在给定时间点产生一个或者多个向量元素，这个警报统计活跃的这些元素标签集。

警报规则在Prometheus系统中用同样的record rules方式进行配置

### 定义警报规则
警报规则的定义遵循下面的风格：
```
ALERT <alert name>
  IF <expression>
    [ FOR <duration> ]
      [ LABELS <label set> ]
        [ ANNOTATIONS <label set> ]
```
`FOR`选项语句会使Prometheus服务等待指定的时间, 在第一次遇到新的表达式输出向量元素（如：具有高HTTP错误率的实例）之间，并将该警报统计为该元素的触发。如果该元素的活跃的，且尚未触发，表示正在挂起状态。

`LABELS`选项语句允许指定额外的标签列表，把它们附加在警告上。任何已存在的冲突标签会被重写。这个标签值能够被模板化。

`ANNOTATIONS`选项语句指定了另一组标签，它们不被当做警告实例的身份标识。它们经常用于存储额外的信息，例如：警告描述，后者runbook链接。这个注释值能够被模板化。

### Templating 模板
```
# Alert for any instance that is unreachable for >5 minutes.
ALERT InstanceDown
  IF up == 0
    FOR 5m
      LABELS { severity = "page" }
        ANNOTATIONS {
                summary = "Instance {{ $labels.instance }} down",
                    description = "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes.",
                      }

# Alert for any instance that have a median request latency >1s.
ALERT APIHighRequestLatency
  IF api_http_request_latencies_second{quantile="0.5"} > 1
    FOR 1m
      ANNOTATIONS {
              summary = "High request latency on {{ $labels.instance }}",
                  description = "{{ $labels.instance }} has a median request latency above 1s (current value: {{ $value }}s)",
                    }
```

### 运行时检查警告
为了能够手动检查哪个警告是活跃的（挂起或者触发），导航到你的Prometheus服务实例的"Alerts"tab页面。这个会显示精确的标签集合，它们每一个定义的警告都是当前活跃的。

对于挂起和触发警告，Prometheus也存储形如`ALERTS{alertname="<alert name>", alertstat=s"pending|firing", <additional alert labels>}`. 只要警告是在指定的活跃（挂起或者触发）状态上，这个样本值设置为1。当一个警告从活跃状态变成不活跃状态时，这个样本值被设置为0。一旦不活跃，这个时间序列将不会再更新。

### 发送警告通知
Prometheus的警告规则擅长确定当前哪个实例有问题。但它们并不是一个完整的通知解决方案。在简单的警报定义上，需要另一个层来添加总结，通知速率限制，silencing，警报依赖。在Prometheus的生态系统中，Alertmanager发挥了这一作用。因此，Prometheus可能被配置为定期向Alertmanager实例发送有关警报信息，该实例然后负责调用正确的通知，可以通过`-alertmanager.url`命令行标志配置Alertmanager实例。

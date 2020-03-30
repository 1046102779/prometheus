警报规则允许您基于Prometheus表达式语言表达式定义警报条件，并将有关触发警报的通知发送到外部服务。 每当警报表达式在给定时间点生成一个或多个向量元素时，警报将计为这些元素的标签集的活动状态。

##### 一、定义报警规则
警报规则在Prometheus中以与记录规则相同的方式配置。

带警报的示例规则文件将是：
```
groups:
- name: example
  rules:
  - alert: HighErrorRate
    expr: job:request_latency_seconds:mean5m{job="myjob"} > 0.5
    for: 10m
    labels:
      severity: page
    annotations:
      summary: High request latency
```
可选的`for`子句使Prometheus在第一次遇到新的表达式输出向量元素和将此警告作为此元素的触发计数之间等待一段时间。 在这种情况下，Prometheus将在每次评估期间检查警报是否继续处于活动状态10分钟，然后再触发警报。 处于活动状态但尚未触发的元素处于暂挂状态。

`labels`子句允许指定要附加到警报的一组附加标签。 任何现有的冲突标签都将被覆盖。 标签值可以是模板化的。

`annotations`子句指定一组信息标签，可用于存储更长的附加信息，例如警报描述或Runbook链接。 注释值可以是模板化的。

##### 二、模板
可以使用控制台模板模板化标签和注释值。 `$labels`变量保存警报实例的标签键/值对，`$value`保存警报实例的评估值。
```
# 要插入触发元素的标签值：
{{ $labels.<labelname> }}
# 要插入触发元素的数值表达式值：
{{ $value }}
```
例子：
```
groups:
- name: example
  rules:

  # 对于任何无法访问> 5分钟的实例的警报。
  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
      severity: page
    annotations:
      summary: "Instance {{ $labels.instance }} down"
      description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  # 对中值请求延迟> 1s的任何实例发出警报。
  - alert: APIHighRequestLatency
    expr: api_http_request_latencies_second{quantile="0.5"} > 1
    for: 10m
    annotations:
      summary: "High request latency on {{ $labels.instance }}"
      description: "{{ $labels.instance }} has a median request latency above 1s (current value: {{ $value }}s)"
```

##### 三、在运行时检查警报
要手动检查哪些警报处于活动状态（待处理或触发），请导航至Prometheus实例的"警报"选项卡。 这将显示每个定义的警报当前处于活动状态的确切标签集。

对于待处理和触发警报，Prometheus还存储`ALERTS{alertname="<alert name>",alertstate ="pending|firing",<additional alert labels>}`形式的合成时间序列。 只要警报处于指示的活动（挂起或触发）状态，样本值就会设置为`1`，并且当不再是这种情况时，系列会标记为过时。

##### 四、发送提醒通知
普罗米修斯的警报规则很好地解决了现在的问题，但它们并不是一个完全成熟的通知解决方案。 需要另一层来在简单警报定义之上添加摘要，通知速率限制，静默和警报依赖性。 在普罗米修斯的生态系统中，Alertmanager承担了这一角色。 因此，Prometheus可以被配置为周期性地向Alertmanager实例发送关于警报状态的信息，然后该实例负责调度正确的通知。
Prometheus可以配置为通过其服务发现集成自动发现可用的Alertmanager实例。

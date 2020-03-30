您可以使用`promtool`来测试您的规则。
```
# 单个测试文件
./promtool test rules test.yml

# 多个测试文件
./promtool test rules test1.yml test2.yml test3.yml
```
##### 一、测试文件格式
```
# 这是要考虑进行测试的规则文件列表。
rule_files:
  [ - <file_name> ]

# 可选的, 默认 = 1m
evaluation_interval: <duration>

# 下面列出组名称的顺序将是规则组的评估顺序（在给定的评估时间）。 订单仅适用于下面提到的组。
# 下面不需要提到所有组。
group_eval_order:
  [ - <group_name> ]

# 所有测试都列在这里。
tests:
  [ - <test_group> ]
```
###### 1.1 `<test_group>`
```
# 系列数据
interval: <duration>
input_series:
  [ - <series> ]

# 上述数据的单元测试

# 警报规则的单元测试。 我们从输入文件中考虑警报规则。
alert_rule_test:
  [ - <alert_test_case> ]

# 单元测试PromQL表达式。
promql_expr_test:
  [ - <promql_test_case> ]
```
###### 1.2 `<series>`
```
# 这遵循通常的系列符号 '<metric name>{<label name>=<label value>, ...}'
# 例子: 
#      series_name{label1="value1", label2="value2"}
#      go_goroutines{job="prometheus", instance="localhost:9090"}
series: <string>

# 这使用扩展表示法。
# 扩展符号：
#     'a+bxc' becomes 'a a+b a+(2*b) a+(3*b) … a+(c*b)'
#     'a-bxc' becomes 'a a-b a-(2*b) a-(3*b) … a-(c*b)'
# 例子:
#     1. '-2+4x3' becomes '-2 2 6 10'
#     2. ' 1-2x4' becomes '1 -1 -3 -5 -7'
values: <string>
```
###### 1.3 `<alert_test_case>`
Prometheus允许您为不同的警报规则使用相同的警报名称。 因此，在此单元测试中，您必须在单个`<alert_test_case>`下列出`alertname`的所有触发警报的并集。
```
# 这是从必须检查警报的时间= 0开始经过的时间。
eval_time: <duration>

# 要测试的警报的名称。
alertname: <string>

# 在给定评估时间在给定警报名称下触发的预期警报列表。 如果您想测试是否不应该触发警报规则，那么您可以提及上述字段并将“exp_alerts”留空。
exp_alerts:
  [ - <alert> ]
```
###### 1.4 `<alert>`
```
# 这些是预期警报的扩展标签和注释。 
# 注意：标签还包括与警报关联的样本标签（与您在`/alerts`中看到的标签相同，没有系列`__name__`和`alertname`）
exp_labels:
  [ <labelname>: <string> ]
exp_annotations:
  [ <labelname>: <string> ]
```
###### 1.5 `<promql_test_case>`
```
# 表达评估
expr: <string>

# 这是从必须检查警报的时间= 0开始经过的时间。
eval_time: <duration>

# 在给定评估时间的预期样品。
exp_samples:
  [ - <sample> ]
```
###### 1.6 `<sample>`
```
# 通常系列表示法中的样本标签 '<metric name>{<label name>=<label value>, ...}'
# 例子: 
#      series_name{label1="value1", label2="value2"}
#      go_goroutines{job="prometheus", instance="localhost:9090"}
labels: <string>

# promql表达式的期望值。
value: <number>
```
##### 二、例子
这是用于通过测试的单元测试的示例输入文件。 `test.yml`是遵循上述语法的测试文件，`alerts.yml`包含警报规则。

将`alerts.yml`放在同一目录中，运行`./promtool test rules test.ym`l。
###### 2.1 `<test.yml>`
```
# This is the main input for unit testing. 
# Only this file is passed as command line argument.

rule_files:
    - alerts.yml

evaluation_interval: 1m

tests:
    # Test 1.
    - interval: 1m
      # Series data.
      input_series:
          - series: 'up{job="prometheus", instance="localhost:9090"}'
            values: '0 0 0 0 0 0 0 0 0 0 0 0 0 0 0'
          - series: 'up{job="node_exporter", instance="localhost:9100"}'
            values: '1+0x6 0 0 0 0 0 0 0 0' # 1 1 1 1 1 1 1 0 0 0 0 0 0 0 0
          - series: 'go_goroutines{job="prometheus", instance="localhost:9090"}'
            values: '10+10x2 30+20x5' # 10 20 30 30 50 70 90 110 130
          - series: 'go_goroutines{job="node_exporter", instance="localhost:9100"}'
            values: '10+10x7 10+30x4' # 10 20 30 40 50 60 70 80 10 40 70 100 130

      # Unit test for alerting rules.
      alert_rule_test:
          # Unit test 1.
          - eval_time: 10m
            alertname: InstanceDown
            exp_alerts:
                # Alert 1.
                - exp_labels:
                      severity: page
                      instance: localhost:9090
                      job: prometheus
                  exp_annotations:
                      summary: "Instance localhost:9090 down"
                      description: "localhost:9090 of job prometheus has been down for more than 5 minutes."
      # Unit tests for promql expressions.
      promql_expr_test:
          # Unit test 1.
          - expr: go_goroutines > 5
            eval_time: 4m
            exp_samples:
                # Sample 1.
                - labels: 'go_goroutines{job="prometheus",instance="localhost:9090"}'
                  value: 50
                # Sample 2.
                - labels: 'go_goroutines{job="node_exporter",instance="localhost:9100"}'
                  value: 50
```
###### 2.2 `<alerts.yml>`
```
# This is the rules file.

groups:
- name: example
  rules:

  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
        severity: page
    annotations:
        summary: "Instance {{ $labels.instance }} down"
        description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."

  - alert: AnotherInstanceDown
    expr: up == 0
    for: 10m
    labels:
        severity: page
    annotations:
        summary: "Instance {{ $labels.instance }} down"
        description: "{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes."
```

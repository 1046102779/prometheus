##### 一、配置规则
Prometheus支持两种类型的规则，这些规则可以定期配置，然后定期评估：记录规则和[警报规则](https://prometheus.io/docs/prometheus/latest/configuration/alerting_rules/)。 要在Prometheus中包含规则，请创建包含必要规则语句的文件，并让Prometheus通过Prometheus配置中的`rule_files`字段加载文件。 规则文件使用YAML。

通过将`SIGHUP`发送到Prometheus进程，可以在运行时重新加载规则文件。 仅当所有规则文件格式正确时才会应用更改。
##### 二、语法检查规则
要在不启动Prometheus服务器的情况下快速检查规则文件在语法上是否正确，请安装并运行Prometheus的`promtool`命令行实用工具：
```
go get github.com/prometheus/prometheus/cmd/promtool
promtool check rules /path/to/example.rules.yml
```
当文件在语法上有效时，检查器将已解析规则的文本表示打印到标准输出，然后以`0`返回状态退出。

如果存在任何语法错误或无效的输入参数，则会向标准错误输出错误消息，并以`1`返回状态退出。

##### 三、录制规则
录制规则允许您预先计算经常需要或计算上昂贵的表达式，并将其结果保存为一组新的时间序列。 因此，查询预先计算的结果通常比每次需要时执行原始表达式快得多。 这对于仪表板尤其有用，仪表板需要在每次刷新时重复查询相同的表达式。

记录和警报规则存在于规则组中。 组内的规则以固定间隔顺序运行。

规则文件的语法是：
```
groups:
  [ - <rule_group> ]
```
一个简单的示例规则文件将是：
```
groups:
  - name: example
    rules:
    - record: job:http_inprogress_requests:sum
      expr: sum(http_inprogress_requests) by (job)
```
###### 3.1 `<rule_group>`
```
# 组的名称。 在文件中必须是唯一的。
name: <string>

# 评估组中的规则的频率。
[ interval: <duration> | default = global.evaluation_interval ]

rules:
  [ - <rule> ... ]
```
###### 3.2 `<rule>`
记录规则的语法是：
```
# 要输出的时间序列的名称。 必须是有效的度量标准名称。
record: <string>

# 要评估的PromQL表达式。 每个评估周期都会在当前时间进行评估，并将结果记录为一组新的时间序列，其中度量标准名称由“记录”给出。
expr: <string>

# 在存储结果之前添加或覆盖的标签。
labels:
  [ <labelname>: <labelvalue> ]
```
警报规则的语法是：
```
# 警报的名称。 必须是有效的度量标准名称。
alert: <string>

# 要评估的PromQL表达式。 每个评估周期都会在当前时间进行评估，并且所有结果时间序列都会成为待处理/触发警报。
expr: <string>

# 警报一旦被退回这段时间就会被视为开启。
# 尚未解雇的警报被认为是未决的。
[ for: <duration> | default = 0s ]

# 为每个警报添加或覆盖的标签。
labels:
  [ <labelname>: <tmpl_string> ]

# 要添加到每个警报的注释。
annotations:
  [ <labelname>: <tmpl_string> ]
```

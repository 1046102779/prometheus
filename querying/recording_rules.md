## 定义recording rules
### 配置规则
Prometheus支持可以配置，然后定期执行的两种规则: recording rules(记录规则)和alerting rules[警告规则](https://prometheus.io/docs/alerting/rules)。为了在Prometheus系统中包括规则，我们需要创建一个包含规则语句的文件，并通过在[Prometheus配置](https://prometheus.io/docs/operating/configuration)的`rule_fields`字段加载这个记录规则文件。

这些规则文件可以通过像Prometheus服务发送`SIGNUP`信号量，实时重载记录规则。如果所有的记录规则有正确的格式和语法，则这些变化能够生效。

### 语法检查规则
在没有启动Prometheus服务之前，想快速知道一个规则文件是否正确，可以通过安装和运行Prometheus的`promtool`命令行工具检验:
> go get github.com/prometheus/prometheus/cmd/promtool
> promtool check-rules /path/to/examples.rules

当记录规则文件是有效的，则这个检查会打印出解析到规则的文本表示，并以返回值0退出程序。

如果有任何语法错误的话，则这个命令行会打印出一个错误信息到标准输出，并以返回值1退出程序。无效的输入参数，则以返回值2退出程序。

### 记录规则
记录规则允许你预先计算经常需要的，或者计算复杂度高的表达式，并将结果保存为一组新的时间序列数据。查询预计算结果通常比需要时进行计算表达式快得多。对于dashboard是非常有用的，因为dashboard需要实时刷新查询表达式的结果。

为了增加一个新记录规则，增加下面的记录规则到你的规则文件中：
> \<new time series name\>[{\<label overrides\>}] = \<expression to record\>

一些例子：
> # 计算每个进程http请求总数，保存到新的度量指标中
> job:http_inprogree_requests:sum = sum(http_inprogress_requests) by (job)
> 
> # 放弃老标签，写入新标签的结果时间序列数据：
> new_time_series{label_to_change="new_value", label_to_drop=""} = old_time_series

记录规则将以Prometheus配置中的`evaluate_interval`字段指定的间隔进行评估。在么个评估周期中，规则语句的右侧表达式将在当前时刻进行评估，并将生成的样本向量作为一组新的时间序列存储，这个时间序列带有当前时间戳和新的指标名称

Prometheus支持在警报的注释和标签以及服务的控制台页面中进行模板化。 模板能够针对本地数据库运行查询，迭代数据，使用条件，格式化数据等.Prometheus模板语言基于[Go模板系统](https://golang.org/pkg/text/template/)。
##### 一、简单的警报字段模板
```
alert: InstanceDown
expr: up == 0
for: 5m
labels:
  severity: page
annotations:
  summary: "Instance {{$labels.instance}} down"
  description: "{{$labels.instance}} of job {{$labels.job}} has been down for more than 5 minutes."
```
对于每个触发的警报，将在每次规则迭代期间执行警报字段模板，因此请保持所有查询和模板的轻量级。 如果您需要更复杂的警报模板，建议您改为链接到控制台。
##### 二、简单的迭代
这将显示实例列表以及它们是否已启动：
```
{{ range query "up" }}
  {{ .Labels.instance }} {{ .Value }}
{{ end }}
```
特别的`.`变量包含每个循环迭代的当前样本的值。
##### 三、展示一个值
```
{{ with query "some_metric{instance='someinstance'}" }}
  {{ . | first | value | humanize }}
{{ end }}
```
Go和Go的模板语言都是强类型的，因此必须检查是否返回了样本以避免执行错误。 例如，如果抓取或规则评估尚未运行，或者主机已关闭，则可能会发生这种情况。

包含的`prom_query_drilldown`模板处理此问题，允许格式化结果，并链接到表达式浏览器。
##### 四、使用命令行URL参数
```
{{ with printf "node_memory_MemTotal{job='node',instance='%s'}" .Params.instance | query }}
  {{ . | first | value | humanize1024}}B
{{ end }}
```
如果以`console.html?instance = hostname`的身份访问，`.Params.instance`将评估为`hostname`。
##### 五、高级迭代
```
<table>
{{ range printf "node_network_receive_bytes{job='node',instance='%s',device!='lo'}" .Params.instance | query | sortByLabel "device"}}
  <tr><th colspan=2>{{ .Labels.device }}</th></tr>
  <tr>
    <td>Received</td>
    <td>{{ with printf "rate(node_network_receive_bytes{job='node',instance='%s',device='%s'}[5m])" .Labels.instance .Labels.device | query }}{{ . | first | value | humanize }}B/s{{end}}</td>
  </tr>
  <tr>
    <td>Transmitted</td>
    <td>{{ with printf "rate(node_network_transmit_bytes{job='node',instance='%s',device='%s'}[5m])" .Labels.instance .Labels.device | query }}{{ . | first | value | humanize }}B/s{{end}}</td>
  </tr>{{ end }}
<table>
```
在这里，我们遍历所有网络设备并显示每个网络设备的网络流量。

由于`range`操作未指定变量，因此。循环内部`.Params.instance`不可用。 现在是循环变量。
##### 六、定义可重用模板
Prometheus支持定义可重用的模板。 当与控制台库支持结合使用时，这一功能特别强大，允许跨控制台共享模板。
```
{{/* Define the template */}}
{{define "myTemplate"}}
  do something
{{end}}

{{/* Use the template */}}
{{template "myTemplate"}}
```
模板仅限于一个参数。 `args`函数可用于包装多个参数。
```
{{define "myMultiArgTemplate"}}
  First argument: {{.arg0}}
  Second argument: {{.arg1}}
{{end}}
{{template "myMultiArgTemplate" (args 1 2)}}
```

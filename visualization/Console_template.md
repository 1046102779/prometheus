
[Golang模板](http://golang.org/pkg/text/template/)创建任意的console。这些由Prometheus服务提供。

控制台模板是创建可在源代码管理中轻松管理的模板的最强大方法。 虽然有一个学习曲线，所以对这种监控方式不熟悉的用户应首先尝试Grafana。

##### 一、开始
Prometheus附带一套示例，让您学习。 这些可以在运行的Prometheus上的`/consoles/index.html.example`中找到，如果Prometheus正在使用`job="node"`标签来抓取节点导出器，则会显示节点导出器控制台。

这个例子控制台包括5部分：
 1. 在顶部的导航栏
 2. 左边的一个菜单
 3. 底部的时间控制
 4. 在中心的主内容，通常是图表
 5. 右边的表格

导航栏用于指向其他系统的链接，例如其他Prometheis，文档以及其他任何对您有意义的内容。 该菜单用于在同一个Prometheus服务器内导航，这对于能够在另一个选项卡中快速打开控制台以关联信息非常有用。 两者都在`console_libraries/menu.lib`中配置。

时间控制允许更改图形的持续时间和范围。 控制台URL可以共享，并为其他人显示相同的图表。

主要内容通常是图表。 提供了一个可配置的JavaScript图形库，可以处理来自Prometheus的请求数据，并通过[Rickshaw](https://tech.shutterstock.com/rickshaw/)进行渲染。

最后，右侧的表格可用于以比图形更紧凑的形式显示统计数据。

##### 二、控制台例子
这是一个基本的控制台。 它显示了右侧表中的任务数，其中有多少，平均CPU使用率和平均内存使用量。 主要内容具有每秒查询图。
```template
{{template "head" .}}

{{template "prom_right_table_head"}}
<tr>
  <th>MyJob</th>
  <th>{{ template "prom_query_drilldown" (args "sum(up{job='myjob'})") }}
      / {{ template "prom_query_drilldown" (args "count(up{job='myjob'})") }}
  </th>
</tr>
<tr>
  <td>CPU</td>
  <td>{{ template "prom_query_drilldown" (args
      "avg by(job)(rate(process_cpu_seconds_total{job='myjob'}[5m]))"
      "s/s" "humanizeNoSmallPrefix") }}
  </td>
</tr>
<tr>
  <td>Memory</td>
  <td>{{ template "prom_query_drilldown" (args
       "avg by(job)(process_resident_memory_bytes{job='myjob'})"
       "B" "humanize1024") }}
  </td>
</tr>
{{template "prom_right_table_tail"}}


{{template "prom_content_head" .}}
<h1>MyJob</h1>

<h3>Queries</h3>
<div id="queryGraph"></div>
<script>
new PromConsole.Graph({
  node: document.querySelector("#queryGraph"),
  expr: "sum(rate(http_query_count{job='myjob'}[5m]))",
  name: "Queries",
  yAxisFormatter: PromConsole.NumberFormatter.humanizeNoSmallPrefix,
  yHoverFormatter: PromConsole.NumberFormatter.humanizeNoSmallPrefix,
  yUnits: "/s",
  yTitle: "Queries"
})
</script>

{{template "prom_content_tail" .}}

{{template "tail"}}
```
`prom_right_table_head`和`prom_right_table_tail`模板包含右侧表。这是可选的。

`prom_query_drilldown`是一个模板，它将评估传递给它的表达式，格式化它，并链接到表达式浏览器中的表达式。第一个参数是表达式。第二个参数是要使用的单位。第三个参数是如何格式化输出。只需要第一个参数。

`prom_query_drilldown`的第三个参数的有效输出格式：

- 未指定：默认转到显示输出。
- `humanize`：使用指标[前缀显示](https://en.wikipedia.org/wiki/Metric_prefix)结果。
- `humanizeNoSmallPrefix`：对于大于1的绝对值，使用度量标准[前缀显示](https://en.wikipedia.org/wiki/Metric_prefix)结果。对于小于1的绝对值，显示3位有效数字。这对于避免可以通过人性化生成的诸如每秒毫微秒的单位是有用的。
- `humanize1024`：使用1024而不是1000的基数显示人性化结果。这通常与`B`一起用作生成`KiB`和`MiB`等单位的第二个参数。
- `printf.3g`：显示3位有效数字。

可以定义自定义格式。有关示例，请参阅[prom.lib](https://github.com/prometheus/prometheus/blob/master/console_libraries/prom.lib)。

##### 三、图库
图库被调用为：
```
<div id="queryGraph"></div>
<script>
new PromConsole.Graph({
  node: document.querySelector("#queryGraph"),
  expr: "sum(rate(http_query_count{job='myjob'}[5m]))"
})
</script>
```
`head`模板加载所需的Javascript和CSS。

图库的参数：
名字|	描述
---|---
expr|	必选. 表达式到图表。 可以是一个清单。
node|	必选. 要渲染的DOM节点。
duration|	可选. 图表的持续时间。 默认为1小时。
endTime	|可选. 图表结束时的Unixtime。 默认为现在。
width	|可选. 图表的宽度，不包括标题。 默认为自动检测。
height	|可选. 图表的高度，不包括标题和图例。 默认为200像素。
min	| 可选. 最小x轴值。 默认为最低数据值。
max	 | 可选. 最小y轴值。 默认为最高数据值。
renderer|	可选. 图表类型。 选项`line`和`area`（堆叠图）。 默认为行。
name	| 可选. 图例和悬停细节中的图表标题。 如果传递了一个字符串，`[[label]]`将被替换为标签值。 如果传递了一个函数，它将传递一个标签映射，并应该将该名称作为字符串返回。 可以是一个清单。
xTitle	| 可选. x轴的标题。 默认为`Time`。
yUnits	| 可选. y轴的单位。 默认为空。
yTitle	| 可选. y轴的标题。 默认为空。
yAxisFormatter	| 可选. y轴的数字格式化程序。 默认为`PromConsole.NumberFormatter.humanize`。
yHoverFormatter|	可选. 悬停细节的数字格式化程序。 默认为`PromConsole.NumberFormatter.humanizeExact`。
colorScheme| 	可选. 图表使用的配色方案。 可以是十六进制[颜色代码列表](https://github.com/shutterstock/rickshaw/blob/master/src/js/Rickshaw.Fixtures.Color.js)，也可以是人力车支持的颜色方案名称之一。 默认为`colorwheel`。

如果`expr`和`name`都是列表，则它们的长度必须相同。 该名称将应用于相应表达式的图。

`yAxisFormatter`和`yHoverFormatter`的有效选项：

- `PromConsole.NumberFormatter.humanize`：使用度[量标准前缀](https://en.wikipedia.org/wiki/Metric_prefix)的格式。
- `PromConsole.NumberFormatter.humanizeNoSmallPrefix`：对于大于1的绝对值，使用度量标准前缀进行格式化。 对于小于1的绝对值，请使用3位有效数字格式。 这对于避免`PromConsole.NumberFormatter.humanize`可以生成的每秒毫秒数等单位很有用。
- `PromConsole.NumberFormatter.humanize1024`：使用1024而不是1000的基数格式化人性化结果。

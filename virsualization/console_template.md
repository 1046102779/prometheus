## console template（控制模板）
---
控制模板允许使用[Go语言模板](http://golang.org/pkg/text/template/)创建任意的console。这些由Prometheus服务提供

console模板是最强有力的方法，它可以在源码控制中创建容易管理的模板，这里有一个学习曲线，所以用户在使用这种新的风格时，应该首先尝试Grafana。

### Getting started
Prometheus提供了一系列的控制模板来帮助您。这些可以在Prometheus服务上的`console/index.html.example`中找到，如果Prometheus服务正在删除带有标签`job="node"`的Node Exporter, 则会显示NodeExporter控制台

这个例子控制台包括5部分：
 1. 在顶部的导航栏
 2. 左边的一个菜单
 3. 底部的时间控制
 4. 在中心的主内容，通常是图表
 5. 右边的表格

这个导航栏是链接到其他系统，例如Prometheus其他方面的文档，以及其他任何使你明白的。该菜单用于在同一个Prometheus服务中导航，它可以快速在另一个tar中打开一个控制台。这些都是在`console_libraries/menu.lib`中配置。

时间控制台允许持久性和图表范围的改变。控制台URLs能够被分享，并且在其他的控制台中显示相同的图表。

主要内容通常是图表。这里有一个可配置的JavaScript图表库，它可以处理来自Prometheus服务的请求，并通过[Rickshaw](http://code.shutterstock.com/rickshaw/)来渲染

最后，在右边的表格可以用笔图表更紧凑的形式显示统计信息。

### 例子控制台 
这是一个最基本的控制台。它显示任务的数量，其中CPU平均使用率、以及右侧表中的平均内存使用率。主要内容具有每秒查询数据。
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

### 模板部分不翻译了，建议大家用Grafana，不喜欢后台服务渲染模板，还是让前端的童鞋去做数据呈现工作吧

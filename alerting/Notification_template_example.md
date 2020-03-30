以下是警报和相应的Alertmanager配置文件设置（alertmanager.yml）的所有不同示例。 每个都使用Go模板系统。

##### 一、自定义Slack通知
在这个例子中，我们定制了Slack通知，以便向我们组织的wiki发送一个URL，告知如何处理已发送的特定警报。
```
global:
  slack_api_url: '<slack_webhook_url>'

route:
  receiver: 'slack-notifications'
  group_by: [alertname, datacenter, app]

receivers:
- name: 'slack-notifications'
  slack_configs:
  - channel: '#alerts'
    text: 'https://internal.myorg.net/wiki/alerts/{{ .GroupLabels.app }}/{{ .GroupLabels.alertname }}'
```

##### 二、访问CommonAnnotations中的注释
在这个例子中，我们再次定制发送给Slack接收器的文本，访问存储在Alertmanager发送的数据的`CommonAnnotations`中的摘要和描述。

警报
```
groups:
- name: Instances
  rules:
  - alert: InstanceDown
    expr: up == 0
    for: 5m
    labels:
      severity: page
    # Prometheus templates apply here in the annotation and label fields of the alert.
    annotations:
      description: '{{ $labels.instance }} of job {{ $labels.job }} has been down for more than 5 minutes.'
      summary: 'Instance {{ $labels.instance }} down'
```
接收器
```
- name: 'team-x'
  slack_configs:
  - channel: '#alerts'
    # Alertmanager templates apply here.
    text: "<!channel> \nsummary: {{ .CommonAnnotations.summary }}\ndescription: {{ .CommonAnnotations.description }}"
```

##### 三、范围内所有收到的警报
最后，假设与前一个示例相同的警报，我们定制我们的接收器以覆盖从Alertmanager接收的所有警报，在新线路上打印它们各自的注释摘要和描述。

接收器
```
- name: 'default-receiver'
  slack_configs:
  - channel: '#alerts'
    title: "{{ range .Alerts }}{{ .Annotations.summary }}\n{{ end }}"
    text: "{{ range .Alerts }}{{ .Annotations.description }}\n{{ end }}"
```

##### 四、定义可重用模板
回到我们的第一个例子，我们还可以提供一个包含命名模板的文件，然后由Alertmanager加载，以避免跨越多行的复杂模板。 在`/alertmanager/template/myorg.tmpl`下创建一个文件，并在其中创建一个名为“slack.myorg.txt”的模板：
```
{{ define "slack.myorg.text" }}https://internal.myorg.net/wiki/alerts/{{ .GroupLabels.app }}/{{ .GroupLabels.alertname }}{{ end}}
```
配置现在加载具有“text”字段的给定名称的模板，并提供自定义模板文件的路径：
```
global:
  slack_api_url: '<slack_webhook_url>'

route:
  receiver: 'slack-notifications'
  group_by: [alertname, datacenter, app]

receivers:
- name: 'slack-notifications'
  slack_configs:
  - channel: '#alerts'
    text: '{{ template "slack.myorg.text" . }}'

templates:
- '/etc/alertmanager/templates/myorg.tmpl'
```
此[博客文章](https://prometheus.io/blog/2016/03/03/custom-alertmanager-templates/)中进一步详细说明了此示例。


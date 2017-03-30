## configuration 配置
---
[Alertmanager](https://github.com/prometheus/alertmanager)通过命令行和一个配置文件配置。命令行配置不可变的系统参数，而配置文件定义inhibiton规则，通知路由和通知接收者。

[可视化编辑器](https://prometheus.io/webtools/alerting/routing-tree-editor)可以帮助构建路由树。

如果想要查看所有命令，请使用命令`alertmanager -h`。

`Alertmanager`能够在运行时动态加载配置文件。如果新的配置有错误，则配置中的变化不会生效，同时错误日志被输出到终端。通过发送`SIGHUP`信号量给这个进程，或者通过HTTP POST请求`/-/reload`，Alertmanager配置动态加载到内存。

### 配置文件
使用`-config.file`指定要加载的配置文件
> ./alertmanager -config.file=simple.yml

这个配置文件使用`yaml`格式编写的，括号表示参数是可选的，对于非列表参数，该值将设置为指定的默认值。
 - `<duration>`: 与正则表达式匹配的持续时间`[0-9]+(ms|[smhdwy])`
 - `<labeltime>`: 与正则表达式匹配的字符串`[a-zA-Z_][a-zA-Z0-9_]*`
 - `<filepath>`: 当前工作目录下的有效路径
 - `<boolean>`: 布尔值： `false` 或者 `true`。
 - `<string>`: 常规字符串
 - `<tmpl_string>`: 一个在使用前被模板扩展的字符串

其他占位符被分开指定， 一个有效的示例，点击[这里](https://github.com/prometheus/alertmanager/blob/master/doc/examples/simple.yml)

全局配置指定的参数在所有其他上下文配置中是有效的。它们也作为其他区域的默认值。
```
global:
  # ResolveTimeout is the time after which an alert is declared resolved
  # if it has not been updated.
  [ resolve_timeout: <duration> | default = 5m ]

  # The default SMTP From header field.
  [ smtp_from: <tmpl_string> ]
  # The default SMTP smarthost used for sending emails.
  [ smtp_smarthost: <string> ]
  # SMTP authentication information.
  [ smtp_auth_username: <string> ]
  [ smtp_auth_password: <string> ]
  [ smtp_auth_secret: <string> ]
  # The default SMTP TLS requirement.
  [ smtp_require_tls: <bool> | default = true ]

  # The API URL to use for Slack notifications.
  [ slack_api_url: <string> ]

  [ pagerduty_url: <string> | default = "https://events.pagerduty.com/generic/2010-04-15/create_event.json" ]
  [ opsgenie_api_host: <string> | default = "https://api.opsgenie.com/" ]
  [ hipchat_url: <string> | default = "https://api.hipchat.com/" ]
  [ hipchat_auth_token: <string> ]

# Files from which custom notification template definitions are read.
# The last component may use a wildcard matcher, e.g. 'templates/*.tmpl'.
templates:
  [ - <filepath> ... ]

# The root node of the routing tree.
route: <route>

# A list of notification receivers.
receivers:
  - <receiver> ...

# A list of inhibition rules.
inhibit_rules:
  [ - <inhibit_rule> ... ]
```

### <route>
一个路由块在路由树和它的孩子中定义了一个节点。如果不设置，它的可选配置参数从父节点中继承其值。

每个警报在已配置路由树的顶部节点，这个节点必须匹配所有警报。然后遍历所有的子节点。如果`continue`设置成`false`, 当匹配到第一个孩子时，它会停止下来；如果`continue`设置成`true`, 则警报将继续匹配后续的兄弟姐妹节点。如果一个警报不匹配一个节点的任何孩子，这个警报将会基于当前节点的配置参数来处理警报。
```
[ receiver: <string> ]
[ group_by: '[' <labelname>, ... ']' ]

# Whether an alert should continue matching subsequent sibling nodes.
[ continue: <boolean> | default = false ]

# A set of equality matchers an alert has to fulfill to match the node.
match:
  [ <labelname>: <labelvalue>, ... ]

# A set of regex-matchers an alert has to fulfill to match the node.
match_re:
  [ <labelname>: <regex>, ... ]

# How long to initially wait to send a notification for a group
# of alerts. Allows to wait for an inhibiting alert to arrive or collect
# more initial alerts for the same group. (Usually ~0s to few minutes.)
[ group_wait: <duration> ]

# How long to wait before sending notification about new alerts that are
# in are added to a group of alerts for which an initial notification
# has already been sent. (Usually ~5min or more.)
[ group_interval: <duration> ]

# How long to wait before sending a notification again if it has already
# been sent successfully for an alert. (Usually ~3h or more).
[ repeat_interval: <duration> ]

# Zero or more child routes.
routes:
  [ - <route> ... ]
```

#### 示例
```
# The root route with all parameters, which are inherited by the child
# routes if they are not overwritten.
route:
  receiver: 'default-receiver'
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h
  group_by: [cluster, alertname]
  # All alerts that do not match the following child routes
  # will remain at the root node and be dispatched to 'default-receiver'.
  routes:
  # All alerts with service=mysql or service=cassandra
  # are dispatched to the database pager.
  - receiver: 'database-pager'
    group_wait: 10s
    match_re:
      service: mysql|cassandra
  # All alerts with the team=frontend label match this sub-route.
  # They are grouped by product and environment rather than cluster
  # and alertname.
  - receiver: 'frontend-pager'
    group_by: [product, environment]
    match:
      team: frontend
```

### <inhibit_rule>
一个inhibition规则是在与另一组匹配器匹配的警报存在的条件下，使匹配一组匹配器的警报失效的规则。两个警报必须具有一组相同的标签。
```
# Matchers that have to be fulfilled in the alerts to be muted.
target_match:
  [ <labelname>: <labelvalue>, ... ]
target_match_re:
  [ <labelname>: <regex>, ... ]

# Matchers for which one or more alerts have to exist for the
# inhibition to take effect.
source_match:
  [ <labelname>: <labelvalue>, ... ]
source_match_re:
  [ <labelname>: <regex>, ... ]

# Labels that must have an equal value in the source and target
# alert for the inhibition to take effect.
[ equal: '[' <labelname>, ... ']' ]
```

### <receiver>
接收者是一个或者多个通知集成的命名配置

**Alertmanager在v0.0.4中可用的其他接收器尚未实现。我们乐意接受任何贡献，并将其添加到新的实现中**
```
# The unique name of the receiver.
name: <string>

# Configurations for several notification integrations.
email_configs:
  [ - <email_config>, ... ]
hipchat_configs:
  [ - <hipchat_config>, ... ]
pagerduty_configs:
  [ - <pagerduty_config>, ... ]
pushover_configs:
  [ - <pushover_config>, ... ]
slack_configs:
  [ - <slack_config>, ... ]
opsgenie_configs:
  [ - <opsgenie_config>, ... ]
webhook_configs:
  [ - <webhook_config>, ... ]
```
### <email_config>
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The email address to send notifications to.
to: <tmpl_string>
# The sender address.
[ from: <tmpl_string> | default = global.smtp_from ]
# The SMTP host through which emails are sent.
[ smarthost: <string> | default = global.smtp_smarthost ]
# SMTP authentication information.
[ auth_username: <string> ]
[ auth_password: <string> ]
[ auth_secret: <string> ]
[ auth_identity: <string> ]

[ require_tls: <bool> | default = global.smtp_require_tls ]

# The HTML body of the email notification.
[ html: <tmpl_string> | default = '{{ template "email.default.html" . }}' ]

# Further headers email header key/value pairs. Overrides any headers
# previously set by the notification implementation.
[ headers: { <string>: <tmpl_string>, ... } ]
```

### <hipchat_config>
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The HipChat Room ID.
room_id: <tmpl_string>
# The auth token.
[ auth_token: <string> | default = global.hipchat_auth_token ]
# The URL to send API requests to.
[ url: <string> | default = global.hipchat_url ]

# See https://www.hipchat.com/docs/apiv2/method/send_room_notification
# A label to be shown in addition to the sender's name.
[ from:  <tmpl_string> | default = '{{ template "hipchat.default.from" . }}' ]
# The message body.
[ message:  <tmpl_string> | default = '{{ template "hipchat.default.message" . }}' ]
# Whether this message should trigger a user notification.
[ notify:  <boolean> | default = false ]
# Determines how the message is treated by the alertmanager and rendered inside HipChat. Valid values are 'text' and 'html'.
[ message_format:  <string> | default = 'text' ]
# Background color for message.
[ color:  <tmpl_string> | default = '{{ if eq .Status "firing" }}red{{ else }}green{{ end }}' ]
```

### <pagerduty_config>
通过PagerDuty ApI发送通知：
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The PagerDuty service key.
service_key: <tmpl_string>
# The URL to send API requests to
[ url: <string> | default = global.pagerduty_url ]

# The client identification of the Alertmanager.
[ client:  <tmpl_string> | default = '{{ template "pagerduty.default.client" . }}' ]
# A backlink to the sender of the notification.
[ client_url:  <tmpl_string> | default = '{{ template "pagerduty.default.clientURL" . }}' ]

# A description of the incident.
[ description: <tmpl_string> | default = '{{ template "pagerduty.default.description" .}}' ]

# A set of arbitrary key/value pairs that provide further detail
# about the incident.
[ details: { <string>: <tmpl_string>, ... } | default = {
  firing:       '{{ template "pagerduty.default.instances" .Alerts.Firing }}'
  resolved:     '{{ template "pagerduty.default.instances" .Alerts.Resolved }}'
  num_firing:   '{{ .Alerts.Firing | len }}'
  num_resolved: '{{ .Alerts.Resolved | len }}'
} ]
```

### <pushover_config>
通过PUSHover API发送通知：
```
# The recipient user’s user key.
user_key: <string>

# Your registered application’s API token, see https://pushover.net/apps
token: <string>

# Notification title.
[ title: <tmpl_string> | default = '{{ template "pushover.default.title" . }}' ]

# Notification message.
[ message: <tmpl_string> | default = '{{ template "pushover.default.message" . }}' ]

# A supplementary URL shown alongside the message.
[ url: <tmpl_string> | default = '{{ template "pushover.default.url" . }}' ]

# Priority, see https://pushover.net/api#priority
[ priority: <tmpl_string> | default = '{{ if eq .Status "firing" }}2{{ else }}0{{ end }}' ]

# How often the Pushover servers will send the same notification to the user.
# Must be at least 30 seconds.
[ retry: <duration> | default = 1m ]

# How long your notification will continue to be retried for, unless the user
# acknowledges the notification.
[ expire: <duration> | default = 1h ]
```

### <slack_config>
通过Slack webhooks发送通知：
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The Slack webhook URL.
[ api_url: <string> | default = global.slack_api_url ]

# The channel or user to send notifications to.
channel: <tmpl_string>

# API request data as defined by the Slack webhook API.
[ color: <tmpl_string> | default = '{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}' ]
[ username: <tmpl_string> | default = '{{ template "slack.default.username" . }}'
[ title: <tmpl_string> | default = '{{ template "slack.default.title" . }}' ]
[ title_link: <tmpl_string> | default = '{{ template "slack.default.titlelink" . }}' ]
[ icon_emoji: <tmpl_string> ]
[ icon_url: <tmpl_string> ]
[ pretext: <tmpl_string> | default = '{{ template "slack.default.pretext" . }}' ]
[ text: <tmpl_string> | default = '{{ template "slack.default.text" . }}' ]
[ fallback: <tmpl_string> | default = '{{ template "slack.default.fallback" . }}' ]
```

### <opsgenie_config>
通过OpsGenie API发送通知:
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The API key to use when talking to the OpsGenie API.
api_key: <string>

# The host to send OpsGenie API requests to.
[ api_host: <string> | default = global.opsgenie_api_host ]

# A description of the incident.
[ description: <tmpl_string> | default = '{{ template "opsgenie.default.description" . }}' ]
# A backlink to the sender of the notification.
[ source: <tmpl_string> | default = '{{ template "opsgenie.default.source" . }}' ]

# A set of arbitrary key/value pairs that provide further detail
# about the incident.
[ details: { <string>: <tmpl_string>, ... } ]

# Comma separated list of team responsible for notifications.
[ teams: <tmpl_string> ]
# Comma separated list of tags attached to the notifications.
[ tags: <tmpl_string> ]
```

### <webhook_config>
webhook接收者允许配置一个通用的接收者
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The endpoint to send HTTP POST requests to.
url: <string>
```

Alertmanager通过HTTP POST请求发送json格式的数据到配置端点：
```
{
  "version": "3",
  "groupKey": <number>     // key identifying the group of alerts (e.g. to deduplicate)
  "status": "<resolved|firing>",
  "receiver": <string>,
  "groupLabels": <object>,
  "commonLabels": <object>,
  "commonAnnotations": <object>,
  "externalURL": <string>,  // backling to the Alertmanager.
  "alerts": [
    {
      "labels": <object>,
      "annotations": <object>,
      "startsAt": "<rfc3339>",
      "endsAt": "<rfc3339>"
    },
    ...
  ]
}
```

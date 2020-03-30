Alertmanager通过命令行标志和配置文件进行配置。 虽然命令行标志配置了不可变的系统参数，但配置文件定义了禁止规则，通知路由和通知接收器。

[可视化编辑器](https://prometheus.io/webtools/alerting/routing-tree-editor/)可以帮助构建路由树。

要查看所有可用的命令行标志，请运行`alertmanager -h`。

Alertmanager可以在运行时重新加载其配置。 如果新配置格式不正确，则不会应用更改并记录错误。 通过向进程发送`SIGHUP`或向`/-/reload`端点发送HTTP POST请求来触发配置重新加载。

##### 一、配置文件
指定要加载的配置文件，使用`--config.file`标志
```
./alertmanager --config.file=simple.yml
```
该文件以YAML格式写入，由下面描述的方案定义。括号表示参数是可选的。对于非列表参数，该值设置为指定的默认值。

通用占位符定义如下：

- `<duration>`：与正则表达式匹配的持续时间`[0-9]+(ms|[smhdwy])`
- `<labelname>`：与正则表达式匹配的字符串`[a-zA-Z _][a-zA-Z0-9 _]*`
- `<labelvalue>`：一串unicode字符
- `<filepath>`：当前工作目录中的有效路径
- `<boolean>`：一个可以取值为`true`或`false`的布尔值
- `<string>`：常规字符串
- `<secret>`：一个秘密的常规字符串，例如密码
- `<tmpl_string>`：在使用前进行模板扩展的字符串
- `<tmpl_secret>`：在使用之前进行模板扩展的字符串，它是一个秘密
其他占位符是单独指定的。

可以在[此处](https://github.com/prometheus/alertmanager/blob/master/doc/examples/simple.yml)找到有效的示例文件。

全局配置指定在所有其他配置上下文中有效的参数。它们还可用作其他配置节的默认值。
```
global:
  # ResolveTimeout is the time after which an alert is declared resolved
  # if it has not been updated.
  [ resolve_timeout: <duration> | default = 5m ]

  # The default SMTP From header field.
  [ smtp_from: <tmpl_string> ]
  # The default SMTP smarthost used for sending emails, including port number.
  # Port number usually is 25, or 587 for SMTP over TLS (sometimes referred to as STARTTLS).
  # Example: smtp.example.org:587
  [ smtp_smarthost: <string> ]
  # The default hostname to identify to the SMTP server.
  [ smtp_hello: <string> | default = "localhost" ]
  [ smtp_auth_username: <string> ]
  # SMTP Auth using LOGIN and PLAIN.
  [ smtp_auth_password: <secret> ]
  # SMTP Auth using PLAIN.
  [ smtp_auth_identity: <string> ]
  # SMTP Auth using CRAM-MD5. 
  [ smtp_auth_secret: <secret> ]
  # The default SMTP TLS requirement.
  [ smtp_require_tls: <bool> | default = true ]

  # The API URL to use for Slack notifications.
  [ slack_api_url: <secret> ]
  [ victorops_api_key: <secret> ]
  [ victorops_api_url: <string> | default = "https://alert.victorops.com/integrations/generic/20131114/alert/" ]
  [ pagerduty_url: <string> | default = "https://events.pagerduty.com/v2/enqueue" ]
  [ opsgenie_api_key: <secret> ]
  [ opsgenie_api_url: <string> | default = "https://api.opsgenie.com/" ]
  [ hipchat_api_url: <string> | default = "https://api.hipchat.com/" ]
  [ hipchat_auth_token: <secret> ]
  [ wechat_api_url: <string> | default = "https://qyapi.weixin.qq.com/cgi-bin/" ]
  [ wechat_api_secret: <secret> ]
  [ wechat_api_corp_id: <string> ]

  # The default HTTP client configuration
  [ http_config: <http_config> ]

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

##### 二、`<route>`
路由块定义路由树中的节点及其子节点。 如果未设置，则其可选配置参数将从其父节点继承。

每个警报都在配置的顶级路由中进入路由树，该路由必须匹配所有警报（即没有任何已配置的匹配器）。 然后它遍历子节点。 如果将`continue`设置为false，则在第一个匹配的子项后停止。 如果匹配节点上的`continue`为true，则警报将继续与后续兄弟节点匹配。 如果警报与节点的任何子节点都不匹配（没有匹配的子节点，或者不存在），则根据当前节点的配置参数处理警报。
```
[ receiver: <string> ]
# The labels by which incoming alerts are grouped together. For example,
# multiple alerts coming in for cluster=A and alertname=LatencyHigh would
# be batched into a single group.
#
# To aggregate by all possible labels use the special value '...' as the sole label name, for example:
# group_by: ['...'] 
# This effectively disables aggregation entirely, passing through all 
# alerts as-is. This is unlikely to be what you want, unless you have 
# a very low alert volume or your upstream notification system performs 
# its own grouping.
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
[ group_wait: <duration> | default = 30s ]

# How long to wait before sending a notification about new alerts that
# are added to a group of alerts for which an initial notification has
# already been sent. (Usually ~5m or more.)
[ group_interval: <duration> | default = 5m ]

# How long to wait before sending a notification again if it has already
# been sent successfully for an alert. (Usually ~3h or more).
[ repeat_interval: <duration> | default = 4h ]

# Zero or more child routes.
routes:
  [ - <route> ... ]
```
例子：
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

##### 三、`<inhibit_rule>`
当存在与另一组匹配器匹配的警报（源）时，禁止规则将匹配一组匹配器的警报（目标）静音。 目标和源警报必须具有相同列表中标签名称的`equal`标签值。

从语义上讲，缺少标签和具有空值的标签是`equal`的。 因此，如果源和目标警报中都缺少所有相同的标签名称，则禁用规则将适用。

为了防止警报抑制自身，禁止规则将永远不会禁止与规则的目标和源侧匹配的警报。 但是，我们建议以警报永远不会匹配双方的方式选择目标和源匹配器。 理由更容易，并且不会触发这种特殊情况。
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

##### 四、`<http_config>`
`http_config`允许配置接收器用于与基于HTTP的API服务通信的HTTP客户端。
```
# Note that `basic_auth`, `bearer_token` and `bearer_token_file` options are
# mutually exclusive.

# Sets the `Authorization` header with the configured username and password.
# password and password_file are mutually exclusive.
basic_auth:
  [ username: <string> ]
  [ password: <secret> ]
  [ password_file: <string> ]

# Sets the `Authorization` header with the configured bearer token.
[ bearer_token: <secret> ]

# Sets the `Authorization` header with the bearer token read from the configured file.
[ bearer_token_file: <filepath> ]

# Configures the TLS settings.
tls_config:
  [ <tls_config> ]

# Optional proxy URL.
[ proxy_url: <string> ]
```

##### 五、`<tls_config>`
`tls_config`允许配置TLS连接。
```
# CA certificate to validate the server certificate with.
[ ca_file: <filepath> ]

# Certificate and key files for client cert authentication to the server.
[ cert_file: <filepath> ]
[ key_file: <filepath> ]

# ServerName extension to indicate the name of the server.
# http://tools.ietf.org/html/rfc4366#section-3.1
[ server_name: <string> ]

# Disable validation of the server certificate.
[ insecure_skip_verify: <boolean> | default = false]
```

##### 六、`<receiver>`
Receiver是一个或多个通知集成的命名配置。

我们没有主动添加新的接收器，我们建议通过webhook接收器实现自定义通知集成。
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
victorops_configs:
  [ - <victorops_config>, ... ]
wechat_configs:
  [ - <wechat_config>, ... ]
```

##### 七、`<email_config>`
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The email address to send notifications to.
to: <tmpl_string>

# The sender address.
[ from: <tmpl_string> | default = global.smtp_from ]

# The SMTP host through which emails are sent.
[ smarthost: <string> | default = global.smtp_smarthost ]

# The hostname to identify to the SMTP server.
[ hello: <string> | default = global.smtp_hello ]

# SMTP authentication information.
[ auth_username: <string> | default = global.smtp_auth_username ]
[ auth_password: <secret> | default = global.smtp_auth_password ]
[ auth_secret: <secret> | default = global.smtp_auth_secret ]
[ auth_identity: <string> | default = global.smtp_auth_identity ]

# The SMTP TLS requirement.
[ require_tls: <bool> | default = global.smtp_require_tls ]

# TLS configuration.
tls_config:
  [ <tls_config> ]

# The HTML body of the email notification.
[ html: <tmpl_string> | default = '{{ template "email.default.html" . }}' ]
# The text body of the email notification.
[ text: <tmpl_string> ]

# Further headers email header key/value pairs. Overrides any headers
# previously set by the notification implementation.
[ headers: { <string>: <tmpl_string>, ... } ]
```

##### 八、`<hipchat_config>`
HipChat通知使用[Build Your Own](https://confluence.atlassian.com/hc/integrations-with-hipchat-server-683508267.html)集成。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The HipChat Room ID.
room_id: <tmpl_string>
# The auth token.
[ auth_token: <secret> | default = global.hipchat_auth_token ]
# The URL to send API requests to.
[ api_url: <string> | default = global.hipchat_api_url ]

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

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```

##### 九、`<pagerduty_config>`
PagerDuty通知通过[PagerDuty API](https://v2.developer.pagerduty.com/)发送。 PagerDuty提供了有关如何[在此](https://www.pagerduty.com/docs/guides/prometheus-integration-guide/)集成的文档。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The following two options are mutually exclusive.
# The PagerDuty integration key (when using PagerDuty integration type `Events API v2`).
routing_key: <tmpl_secret>
# The PagerDuty integration key (when using PagerDuty integration type `Prometheus`).
service_key: <tmpl_secret>

# The URL to send API requests to
[ url: <string> | default = global.pagerduty_url ]

# The client identification of the Alertmanager.
[ client:  <tmpl_string> | default = '{{ template "pagerduty.default.client" . }}' ]
# A backlink to the sender of the notification.
[ client_url:  <tmpl_string> | default = '{{ template "pagerduty.default.clientURL" . }}' ]

# A description of the incident.
[ description: <tmpl_string> | default = '{{ template "pagerduty.default.description" .}}' ]

# Severity of the incident.
[ severity: <tmpl_string> | default = 'error' ]

# A set of arbitrary key/value pairs that provide further detail
# about the incident.
[ details: { <string>: <tmpl_string>, ... } | default = {
  firing:       '{{ template "pagerduty.default.instances" .Alerts.Firing }}'
  resolved:     '{{ template "pagerduty.default.instances" .Alerts.Resolved }}'
  num_firing:   '{{ .Alerts.Firing | len }}'
  num_resolved: '{{ .Alerts.Resolved | len }}'
} ]

# Images to attach to the incident.
images:
  [ <image_config> ... ]

# Links to attach to the incident.
links:
  [ <link_config> ... ]

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```

###### 9.1 `<image_config>`
这些字段记录在[PagerDuty API文档](https://v2.developer.pagerduty.com/docs/send-an-event-events-api-v2#section-the-images-property)中。
```
source: <tmpl_string>
alt: <tmpl_string>
text: <tmpl_string>
```
###### 9.2 `<link_config>`
这些字段记录在[PagerDuty API文档](https://v2.developer.pagerduty.com/docs/send-an-event-events-api-v2#section-the-images-property)中。
```
href: <tmpl_string>
text: <tmpl_string>
```

##### 十、`<pushover_config>`
推送通知通过[Pushover API](https://pushover.net/api)发送。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The recipient user’s user key.
user_key: <secret>

# Your registered application’s API token, see https://pushover.net/apps
token: <secret>

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

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```

##### 十一、`<slack_config>`
Slack通知通过[Slack webhooks](https://api.slack.com/incoming-webhooks)发送。 通知包含[附件](https://api.slack.com/docs/message-attachments)。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The Slack webhook URL.
[ api_url: <secret> | default = global.slack_api_url ]

# The channel or user to send notifications to.
channel: <tmpl_string>

# API request data as defined by the Slack webhook API.
[ icon_emoji: <tmpl_string> ]
[ icon_url: <tmpl_string> ]
[ link_names: <boolean> | default = false ]
[ username: <tmpl_string> | default = '{{ template "slack.default.username" . }}' ]
# The following parameters define the attachment.
actions:
  [ <action_config> ... ]
[ callback_id: <tmpl_string> | default = '{{ template "slack.default.callbackid" . }}' ]
[ color: <tmpl_string> | default = '{{ if eq .Status "firing" }}danger{{ else }}good{{ end }}' ]
[ fallback: <tmpl_string> | default = '{{ template "slack.default.fallback" . }}' ]
fields:
  [ <field_config> ... ]
[ footer: <tmpl_string> | default = '{{ template "slack.default.footer" . }}' ]
[ pretext: <tmpl_string> | default = '{{ template "slack.default.pretext" . }}' ]
[ short_fields: <boolean> | default = false ]
[ text: <tmpl_string> | default = '{{ template "slack.default.text" . }}' ]
[ title: <tmpl_string> | default = '{{ template "slack.default.title" . }}' ]
[ title_link: <tmpl_string> | default = '{{ template "slack.default.titlelink" . }}' ]
[ image_url: <tmpl_string> ]
[ thumb_url: <tmpl_string> ]

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```
###### 11.1 `<action_config>`
这些字段记录在[Slack API文档](https://api.slack.com/docs/message-attachments#action_fields)中。
```
type: <tmpl_string>
text: <tmpl_string>
url: <tmpl_string>
[ style: <tmpl_string> [ default = '' ]
```
###### 11.2 `<field_config>`
这些字段记录在[Slack API文档](https://api.slack.com/docs/message-attachments#action_fields)中。
```
title: <tmpl_string>
value: <tmpl_string>
[ short: <boolean> | default = slack_config.short_fields ]
```

##### 十二、`<opsgenie_config>`
OpsGenie通知通过[OpsGenie API](https://docs.opsgenie.com/docs/alert-api)发送。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The API key to use when talking to the OpsGenie API.
[ api_key: <secret> | default = global.opsgenie_api_key ]

# The host to send OpsGenie API requests to.
[ api_url: <string> | default = global.opsgenie_api_url ]

# Alert text limited to 130 characters.
[ message: <tmpl_string> ]

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

# Additional alert note.
[ note: <tmpl_string> ]

# Priority level of alert. Possible values are P1, P2, P3, P4, and P5.
[ priority: <tmpl_string> ]

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```

##### 十三、`<victorcops_config>`
VictorOps通知通过[VictorOps API](https://help.victorops.com/knowledge-base/victorops-restendpoint-integration/)发送出去
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The API key to use when talking to the VictorOps API.
[ api_key: <secret> | default = global.victorops_api_key ]

# The VictorOps API URL.
[ api_url: <string> | default = global.victorops_api_url ]

# A key used to map the alert to a team.
routing_key: <tmpl_string>

# Describes the behavior of the alert (CRITICAL, WARNING, INFO).
[ message_type: <tmpl_string> | default = 'CRITICAL' ]

# Contains summary of the alerted problem.
[ entity_display_name: <tmpl_string> | default = '{{ template "victorops.default.entity_display_name" . }}' ]

# Contains long explanation of the alerted problem.
[ state_message: <tmpl_string> | default = '{{ template "victorops.default.state_message" . }}' ]

# The monitoring tool the state message is from.
[ monitoring_tool: <tmpl_string> | default = '{{ template "victorops.default.monitoring_tool" . }}' ]

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```

##### 十四、`<webhook_config>`
webhook接收器允许配置通用接收器。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = true ]

# The endpoint to send HTTP POST requests to.
url: <string>

# The HTTP client's configuration.
[ http_config: <http_config> | default = global.http_config ]
```
Alertmanager将以下列JSON格式将HTTP POST请求发送到配置的端点：
```
{
  "version": "4",
  "groupKey": <string>,    // key identifying the group of alerts (e.g. to deduplicate)
  "status": "<resolved|firing>",
  "receiver": <string>,
  "groupLabels": <object>,
  "commonLabels": <object>,
  "commonAnnotations": <object>,
  "externalURL": <string>,  // backlink to the Alertmanager.
  "alerts": [
    {
      "status": "<resolved|firing>",
      "labels": <object>,
      "annotations": <object>,
      "startsAt": "<rfc3339>",
      "endsAt": "<rfc3339>",
      "generatorURL": <string> // identifies the entity that caused the alert
    },
    ...
  ]
}
```
有一个与此功能[集成](https://prometheus.io/docs/operating/integrations/#alertmanager-webhook-receiver)的列表。


##### 十五、`<wechat_config>`
微信通知通过[微信API](https://mp.weixin.qq.com/)发送。
```
# Whether or not to notify about resolved alerts.
[ send_resolved: <boolean> | default = false ]

# The API key to use when talking to the WeChat API.
[ api_secret: <secret> | default = global.wechat_api_secret ]

# The WeChat API URL.
[ api_url: <string> | default = global.wechat_api_url ]

# The corp id for authentication.
[ corp_id: <string> | default = global.wechat_api_corp_id ]

# API request data as defined by the WeChat API.
[ message: <tmpl_string> | default = '{{ template "wechat.default.message" . }}' ]
[ agent_id: <string> | default = '{{ template "wechat.default.agent_id" . }}' ]
[ to_user: <string> | default = '{{ template "wechat.default.to_user" . }}' ]
[ to_party: <string> | default = '{{ template "wechat.default.to_party" . }}' ]
[ to_tag: <string> | default = '{{ template "wechat.default.to_tag" . }}' ]``

```

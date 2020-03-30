免责声明：Prometheus会自动处理由其配置的警报规则生成的警报。 强烈建议根据时间序列数据在Prometheus中配置警报规则，而不是实现直接客户端。

Alertmanager在`/api/v1/alerts`上侦听API端点上的警报。 只要客户仍处于活动状态（通常为30秒至3分钟），客户就会不断重新发送警报。 客户端可以通过以下格式的POST请求将警报列表推送到该端点：
```
[
  {
    "labels": {
      "alertname": "<requiredAlertName>",
      "<labelname>": "<labelvalue>",
      ...
    },
    "annotations": {
      "<labelname>": "<labelvalue>",
    },
    "startsAt": "<rfc3339>",
    "endsAt": "<rfc3339>",
    "generatorURL": "<generator_url>"
  },
  ...
]
```
标签用于标识警报的相同实例并执行重复数据删除。 注释始终设置为最近收到的注释，而不是识别警报。

两个时间戳都是可选的。 如果省略`startsAt`，则当前时间由Alertmanager分配。 endsAt仅在已知警报结束时间时设置。 否则，它将设置为自上次收到警报以来的可配置超时时间。

`generatorURL`字段是唯一的反向链接，用于标识客户端中此警报的生成实体。

## 发送警报
---
免责声明：Prometheus自动处理发送由其配置的警报规则生成的警报。强烈建议你根据时间序列数据配置Prometheus中的警报规则，而不是直接使用客户端。

Alertmanager用http的API`/api/v1/alerts`监听警报。只要Alertmanager仍然活跃（经常使用30s到3min时间），客户端期望持续地重发警报。客户端通过下面的POST请求，能够推送警报列表到指定端点：
```
[
  {
    "labels": {
      "<labelname>": "<labelvalue>",
      ...
    },
    "annotations": {
      "<labelname>": "<labelvalue>",
    },
    "startsAt": "<rfc3339>",
    "endsAt": "<rfc3339>"
    "generatorURL": "<generator_url>"
  },
  ...
]
```

这个标签用于识别一个警告的唯一实例和执行去重数据操作。这个注释总是设置给最近经常被接收的警告实例。

timestamps是可选的。如果`startsAt`省略，这个当前时间被赋值给Alertmanager。如果一个警报的结束时间是已知的，则只有`endsAt`被设置。如果这个警报是最后被接收的，它将会设置一个可配置的超时时间。

`generatorURL`字段是唯一的后端链接，用于标识客户端中此警报的引发实体。

Alertmanager还支持`/api/alerts`上的传统端点。与Prometheus的v0.16.2级更低版本兼容。


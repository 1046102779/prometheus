## 工具
---
### 客户端库
在你能够监控你的服务器之前，你需要通过Prometheus客户端库把监控的代码放在被监控的服务代码中。下面实现了Prometheus的度量指标类型metric types。

选择你需要的客户端语言，在你的服务实例上通过HTTP端口提供内部度量指标。

 - [Go](https://github.com/prometheus/client_golang)
 - [Java or Scala](https://github.com/prometheus/client_java)
 - [Python](https://github.com/prometheus/client_python)
 - [Ruby](https://github.com/prometheus/client_ruby)

非正式的第三方客户端库
 - [Bash](https://github.com/aecolley/client_bash)
 - [C++](https://github.com/jupp0r/prometheus-cpp)
 - [Common Lisp](https://github.com/deadtrickster/prometheus.cl)
 - [Elixir](https://github.com/deadtrickster/prometheus.ex)
 - [Erlang](https://github.com/deadtrickster/prometheus.erl)
 - [Haskell](https://github.com/fimad/prometheus-haskell)
 - [Lua for Nginx](https://github.com/knyar/nginx-lua-prometheus)
 - [Lua for Tarantool](https://github.com/tarantool/prometheus)
 - [.Net/C#](https://github.com/andrasm/prometheus-net)
 - [Node.js](https://github.com/siimon/prom-client)
 - [PHP](https://github.com/Jimdo/prometheus_client_php)
 - [Rust](https://github.com/pingcap/rust-prometheus)

当Prometheus获取实例的HTTP端点时，客户库发送所有跟踪的度量指标数据到服务器上。

如果没有可用的客户端语言版本，或者你想要避免依赖，你也可以实现一个支持的导入格式到度量指标数据中。

在实现一个新的Prometheus客户端库时，请遵循[客户端指南](https://prometheus.io/docs/instrumenting/writing_clientlibs)。注意，这个文档在仍然在更新中。同时也请关注[开发邮件列表](https://groups.google.com/forum/#!forum/prometheus-developers)。我们非常乐意地给出合适的意见或者建议。

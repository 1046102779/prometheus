客户端库向其代码添加检测。 这些实现了Prometheus度量标准类型。

选择与您的应用程序编写语言相匹配的Prometheus客户端库。 这允许您通过应用程序实例上的HTTP端点定义和公开内部指标：

 - [Go](https://github.com/prometheus/client_golang)
 - [Java or Scala](https://github.com/prometheus/client_java)
 - [Python](https://github.com/prometheus/client_python)
 - [Ruby](https://github.com/prometheus/client_ruby)

非正式的第三方客户端库
 - [Bash](https://github.com/aecolley/client_bash)
 - [C](https://github.com/digitalocean/prometheus-client-c)
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

当Prometheus擦除实例的HTTP端点时，客户端库会将所有跟踪的度量标准的当前状态发送到服务器。

如果您的语言没有可用的客户端库，或者您希望避免依赖性，那么您也可以自己实现一种受支持的[展示格式](https://prometheus.io/docs/instrumenting/exposition_formats/)以公开指标。

在实施新的Prometheus客户端库时，请[遵循编写客户端库的指导原则](https://prometheus.io/docs/instrumenting/writing_clientlibs/)。 请注意，此文档仍在进行中。 另请考虑咨询[开发邮件列表](https://groups.google.com/forum/#!forum/prometheus-developers)。 我们很乐意就如何使您的库尽可能有用和一致提供建议。实现一个新的Prometheus客户端库时，请遵循[客户端指南](https://prometheus.io/docs/instrumenting/writing_clientlibs)。注意，这个文档在仍然在更新中。同时也请关注[开发邮件列表](https://groups.google.com/forum/#!forum/prometheus-developers)。我们非常乐意地给出合适的意见或者建议。


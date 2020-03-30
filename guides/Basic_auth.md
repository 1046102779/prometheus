## 使用基本知识保护PROMETHEUS API和UI端点
---
Prometheus不直接支持与Prometheus表达式浏览器和HTTP API连接的[基本身份验证](https://en.wikipedia.org/wiki/Basic_access_authentication)（也称为“基本身份验证”）。 如果您要为这些连接强制执行基本身份验证，我们建议将Prometheus与反向代理结合使用，并在代理层应用身份验证。 您可以使用Prometheus的任何反向代理，但在本指南中我们将提供一个nginx示例。

### Nginx例子
假设你想在运行在`localhost：12321`上的nginx服务器后面运行一个Prometheus实例，并且所有Prometheus端点都可以通过`/prometheus`端点运行。 因此，Prometheus'`/metrics`端点的完整URL将是：
```
http://localhost:12321/prometheus/metrics
```
我们还要假设，您需要访问Prometheus实例的所有用户的用户名和密码。 对于此示例，请使用`admin`作为用户名并选择您想要的任何密码。

首先，使用`htpasswd`工具创建一个`.htpasswd`文件来存储用户名/密码，并将其存储在`/etc/nginx`目录中：
```
mkdir -p /etc/nginx
htpasswd -c /etc/nginx/.htpasswd admin
```

### Nginx配置
下面是一个示例`nginx.conf`配置文件（存储在`/etc/nginx/.htpasswd中`）。 使用此配置，nginx将对`/prometheus`端点（代理Prometheus）的所有连接强制执行基本身份验证：
```
http {
    server {
        listen 12321;

        location /prometheus {
            auth_basic           "Prometheus";
            auth_basic_user_file /etc/nginx/.htpasswd;

            proxy_pass           http://localhost:9090/;
        }
    }
}

events {}
```
使用上面的配置启动nginx：
```
nginx -c /etc/nginx/nginx.conf
```

### Prometheus配置
在nginx代理后面运行Prometheus时，您需要将外部URL设置为`http：//localhost：12321/prometheus`，并将路由前缀设置为`/`：
```
prometheus \
  --config.file=/path/to/prometheus.yml \
  --web.external-url=http://localhost:12321/prometheus \
  --web.route-prefix="/"
```

### 测试
您可以使用cURL与本地nginx / Prometheus设置进行交互。 试试这个请求：
```
curl --head http://localhost:12321/prometheus/graph
```
这将返回`401 Unauthorized`响应，因为您未能提供有效的用户名和密码。 响应还将包含由nginx提供的`WWW-Authenticate：Basic realm="Prometheus"`标头，表示强制执行由nginx的auth_basic参数指定的`Prometheus``auth_basic`领域。

要使用基本身份验证（例如`/metrics`端点）成功访问Prometheus端点，请使用`-u`标志提供正确的用户名，并在提示时提供密码：
```
curl -u admin http://localhost:12321/prometheus/metrics
Enter host password for user 'admin':
```
那应该返回Prometheus指标输出，看起来应该是这样的：
```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.0001343
go_gc_duration_seconds{quantile="0.25"} 0.0002032
go_gc_duration_seconds{quantile="0.5"} 0.0004485
...
```

### 总结
在本指南中，您将用户名和密码存储在`.htpasswd`文件中，配置nginx以使用该文件中的凭据来验证访问Prometheus的HTTP端点，启动nginx以及配置Prometheus以进行反向代理的用户。
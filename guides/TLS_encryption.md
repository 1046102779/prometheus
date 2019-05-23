## 使用TLS加密保护PROMETHEUS API和UI端点
---
Prometheus不直接支持与Prometheus实例（即表达式浏览器或HTTP API）连接的[传输层安全性](https://en.wikipedia.org/wiki/Transport_Layer_Security)（TLS）加密。 如果您想为这些连接强制执行TLS，我们建议将Prometheus与[反向代理](https://www.nginx.com/resources/glossary/reverse-proxy-server/)结合使用，并在代理层应用TLS。 您可以使用Prometheus的任何反向代理，但在本指南中我们将提供一个[nginx示例](https://prometheus.io/docs/guides/tls-encryption/#nginx-example)。

### nginx例子
假设您想要在`example.co`m域（您拥有）的nginx服务器后面运行Prometheus实例，并且可以通过`/prometheus`端点获得所有Prometheus端点。 因此，Prometheus'`/metrics`端点的完整URL将是：
```
https://example.com/prometheus/metrics
```
我们还假设你使用[OpenSSL](https://www.digitalocean.com/community/tutorials/openssl-essentials-working-with-ssl-certificates-private-keys-and-csrs)或类似工具生成了以下内容：

- `/root/certs/example.com/example.com.crt`上的SSL证书
- `/root/certs/example.com/example.com.key`上的SSL密钥

您可以使用以下命令生成自签名证书和私钥：
```
mkdir -p /root/certs/example.com && cd /root/certs/example.com
openssl req \
  -x509 \
  -newkey rsa:4096 \
  -nodes \
  -keyout example.com.key \
  -out example.com.crt
```
根据提示填写相应的信息，并确保在`Common Name`提示符下输入`example.com`。

### nginx配置
下面是一个示例`nginx.conf`配置文件。 使用此配置，nginx将：

- 使用您提供的证书和密钥强制执行TLS加密
- 将与`/prometheus`端点的所有连接代理到在同一主机上运行的Prometheus服务器（同时从URL中删除`/prometheus`）

```
http {
    server {
        listen              443 ssl;
        server_name         example.com;
        ssl_certificate     /root/certs/example.com/example.com.crt;
        ssl_certificate_key /root/certs/example.com/example.com.key;

        location /prometheus {
            proxy_pass http://localhost:9090/;
        }
    }
}

events {}
```
以root身份启动nginx（因为nginx需要绑定到端口443）：
```
sudo nginx -c /usr/local/etc/nginx/nginx.conf
```

### Prometheus配置
在nginx代理后面运行Prometheus时，您需要将外部URL设置为`http://example.com/prometheus`，并将路由前缀设置为`/`：
```
prometheus \
  --config.file=/path/to/prometheus.yml \
  --web.external-url=http://example.com/prometheus \
  --web.route-prefix="/"
```

### 测试
如果您想使用`example.com`域在本地测试nginx代理，可以在`/etc/hosts`文件中添加一个条目，将`example.com`重新路由到`localhost`：
```
127.0.0.1     example.com
```
然后，您可以使用cURL与本地nginx / Prometheus设置进行交互：
```
curl --cacert /root/certs/example.com/example.com.crt \
  https://example.com/prometheus/api/v1/label/job/values
```
您可以使用`--insecure`或`-k`标志连接到nginx服务器而不指定证书：
```
curl -k https://example.com/prometheus/api/v1/label/job/values
```
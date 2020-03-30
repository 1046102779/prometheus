Prometheus提供了一组管理API，以简化自动化和集成。

##### 一、健康检查
```
GET /-/healthy
```
这个端点始终返回200，应用于检查Prometheus的健康状况。

##### 二、准备检查
```
GET /-/ready
```
当Prometheus准备服务流量（即响应查询）时，此端点返回200。

##### 三、重新加载
```
PUT  /-/reload
POST /-/reload
```
该端点触发Prometheus配置和规则文件的重新加载。 默认情况下它是禁用的，可以通过`--web.enable-lifecycle`标志启用。

或者，可以通过将`SIGHUP`发送到Prometheus进程来触发配置重载。

##### 四、退出
```
PUT  /-/quit
POST /-/quit
```
该端点触发Prometheus的正常关闭。 默认情况下它是禁用的，可以通过`--web.enable-lifecycle`标志启用。

或者，可以通过将`SIGTERM`发送到Prometheus进程来触发正常关闭。
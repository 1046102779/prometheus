## 安装
---
### 使用预编译二进制文件
我们为大多数Prometheus的官方组件，提供了预编译二进制文件。可用版本[下载](https://prometheus.io/download)列表

### 源码安装
对于源码安装Prometheus组件，可以查看`Makefile`目标文件。
> 注意点：在web上的文档指向最新的稳定版(不包括预发布版)。[下一个版本](https://github.com/prometheus/docs/compare/next-release)指向master分支还没有发布的版本

### Docker安装
所有的Prometheus服务在Docker镜像[prom](https://hub.docker.com/u/prom/)仓库中，都是可用的

在Docker上运行Prometheus服务，只需要简单地执行`docker run -p 9090:9090 prom/prometheus`即可。这条命令会启动Prometheus服务，使用的是默认配置文件，并暴露出web可以访问的9090端口

Prometheus镜像使用Volumn存储实际度量指标。在生产环境上使用[数据卷容器](https://docs.docker.com/engine/userguide/containers/dockervolumes/#creating-and-mounting-a-data-volume-container)模式达到轻松管理Prometheus升级数据的目的，它是被强烈推荐的 

为了提供你自己的配置，这儿有几个选项。下面有两个例子。

#### 卷&绑定挂载
在运行Prometheus服务的主机上，绑定挂载你的prometheus.yml配置文件:
> docker run -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

或者为这个配置文件使用另外的卷

> docker run -p 9090:9090 -v /prometheus-data \
       prom/prometheus -config.file=/prometheus-data/prometheus.yml

#### 自定义镜像
为了避免在主机绑定和挂载配置文件，配置直接打包到镜像中。如果配置是静态的，并且在所有环境中配置都是一样的话，这种把配置文件直接打包到镜像中的方式是非常直接推荐的。

例如：在Prometheus配置和Dockerfile中，创建一个新目录：
> FROM prom/prometheus
> ADD prometheus.yml /etc/prometheus/

现在编译和运行它：
> docker build -t my-prometheus .
> docker run -p 9090:9090 my-prometheus

一个更高级的选项是可以通过一些工具动态地渲染配置，甚至后台定期地更新配置

### 使用配置管理系统
如果你喜欢使用配置管理系统，你可能对下面地第三方库感兴趣：

Ansible：
 - [griggheo/ansible-prometheus](https://github.com/griggheo/ansible-prometheus)
 - [William-Yeh/ansible-prometheus](https://github.com/William-Yeh/ansible-prometheus)

Chef:
 - [rayrod2030/chef-prometheus](https://github.com/rayrod2030/chef-prometheus)

SaltStack:
 - [bechtoldt/saltstack-prometheus-formula](https://github.com/bechtoldt/saltstack-prometheus-formula)

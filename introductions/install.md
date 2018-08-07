## 安装
---
### 使用预编译二进制文件
我们为Prometheus大多数的官方组件，提供了预编译二进制文件。可用版本[下载](https://prometheus.io/download)列表

### 源码安装
如果要从源码安装Prometheus的官方组件，可以查看各个项目源码目录下的`Makefile`
> 注意点：在web上的文档指向最新的稳定版(不包括预发布版)。[下一个版本](https://github.com/prometheus/docs/compare/next-release)指向master分支还没有发布的版本

### Docker安装
所有Prometheus服务的Docker镜像在官方组织[prom](https://hub.docker.com/u/prom/)下，都是可用的

在Docker上运行Prometheus服务，只需要简单地执行`docker run -p 9090:9090 prom/prometheus`命令行即可。这条命令会启动Prometheus服务，使用的是默认配置文件，并对外界暴露9090端口

Prometheus镜像使用docker中的volumn卷存储实际度量指标。在生产环境上使用[容器卷](https://docs.docker.com/engine/userguide/containers/dockervolumes/#creating-and-mounting-a-data-volume-container)模式, 可以在Prometheus更新和升级时轻松管理Prometheus数据， 这种使用docker volumn卷方式存储数据，是被docker官方强烈推荐的.

通过几个选项，可以达到使用自己的配置的目的。下面有两个例子。

#### 卷&绑定挂载
在运行Prometheus服务的主机上，做一个本地到docker容器的配置文件关系映射。
> docker run -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

或者为这个配置文件使用一个独立的volumn
> docker run -p 9090:9090 -v /prometheus-data \
       prom/prometheus -config.file=/prometheus-data/prometheus.yml

#### 自定义镜像
为了避免在主机上与docker映射配置文件，我们可以直接将配置文件拷贝到docker镜像中。如果Prometheus配置是静态的，并且在所有服务器上的配置相同，这种把配置文件直接拷贝到镜像中的方式是非常好的。

例如：利用Dockerfile创建一个Prometheus配置目录， Dockerfile应该这样写：
```
FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
```

然后编译和运行它：
```
docker build -t my-prometheus .
docker run -p 9090:9090 my-prometheus
```

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

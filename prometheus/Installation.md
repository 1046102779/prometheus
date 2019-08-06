##### 一、使用预编译二进制文件
我们为Prometheus大多数的官方组件，提供了预编译二进制文件。可用版本[下载](https://prometheus.io/download)列表

##### 二、源码安装
如果要从源码安装Prometheus的官方组件，可以查看各个项目源码目录下的`Makefile`

##### 三、Docker安装
所有Prometheus服务的Docker镜像在官方组织[Quay.io](https://quay.io/repository/prometheus/prometheus)或者[Docker Hub](https://hub.docker.com/u/prom/)下，都是可用的。

在`Docker`上运行Prometheus服务，只需要简单地执行`docker run -p 9090:9090 prom/prometheus`命令行即可。这条命令会启动Prometheus服务，使用的是默认配置文件，并对外界暴露`9090`端口。

Prometheus镜像使用`docker`中的`volumn`卷存储实际度量指标。在生产环境上使用[容器卷](https://docs.docker.com/engine/userguide/containers/dockervolumes/#creating-and-mounting-a-data-volume-container)模式,可以在Prometheus更新和升级时轻松管理Prometheus数据，这种使用`docker volumn`卷方式存储数据，是被`docker`官方强烈推荐的。

通过几个选项，可以达到使用自己的配置的目的。下面有两个例子。

###### 3.1 卷&绑定挂载
通过运行以下命令从主机绑定您的prometheus.yml：
> docker run -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

或者为这个配置文件使用一个独立的volumn
> docker run -p 9090:9090 -v /prometheus-data \
       prom/prometheus -config.file=/prometheus-data/prometheus.yml

###### 3.2 自定义镜像
为避免管理主机上的文件并对其进行绑定安装，可以将配置烘焙到映像中。 如果配置本身相当静态并且在所有环境中都相同，则此方法很有效。

为此，使用Prometheus配置和`Dockerfile`创建一个新目录，如下所示：
```
FROM prom/prometheus
ADD prometheus.yml /etc/prometheus/
```

然后编译和运行它：
```
docker build -t my-prometheus .
docker run -p 9090:9090 my-prometheus
```

一个更高级的选项是可以通过一些工具动态地渲染配置，甚至后台定期地更新配置。

##### 四、使用配置管理系统
如果你喜欢使用配置管理系统，你可能对下面地第三方库感兴趣：

Ansible：
 - [griggheo/ansible-prometheus](https://github.com/griggheo/ansible-prometheus)
 - [William-Yeh/ansible-prometheus](https://github.com/William-Yeh/ansible-prometheus)

Chef:
 - [rayrod2030/chef-prometheus](https://github.com/rayrod2030/chef-prometheus)

Puppet：
 - [puppet/prometheus](https://forge.puppet.com/puppet/prometheus)

SaltStack:
 - [bechtoldt/saltstack-prometheus-formula](https://github.com/bechtoldt/saltstack-prometheus-formula)

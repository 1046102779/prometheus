## Grafana支持Prometheus可视化
---
[Grafana](http://grafana.org/)支持Prometheus查询。从Grafana 2.5.0 (2015-10-28)开始Prometheus可以作为它的数据源。

下面的例子：Prometheus查询在Grafana Dashboard界面的图表展示
![Deme Dashboard](https://prometheus.io/assets/grafana_prometheus-cbb943f0bb3.png)

### Grafana安装
如果要Grafana的完整安装教程，详见[Grafana官方文档](http://docs.grafana.org/installation/)

在Linux安装Grafana，如下所示：
```Grafana install
# Download and unpack Grafana from binary tar (adjust version as appropriate).
curl -L -O https://grafanarel.s3.amazonaws.com/builds/grafana-2.5.0.linux-x64.tar.gz
tar zxf grafana-2.5.0.linux-x64.tar.gz

# Start Grafana.
cd grafana-2.5.0/
./bin/grafana-server web
```

## 使用方法
默认情况下，Grafana服务端口[http://localhost:3000](http://localhost:3000)。默认登录用户名和密码“admin/admin”。

#### 创建一个Prometheus数据源

为了创建一个Prometheus数据源Data source：
 1. 点击Grafana的logo，打开工具栏。
 2. 在工具栏中，点击"Data Source"菜单。
 3. 点击"Add New"。
 4. 数据源Type选择“Prometheus”。
 5. 设置Prometheus服务访问地址（例如：`http://localhost:9090`）。
 6. 调整其他想要的设置（例如：关闭代理访问）。
 7. 点击“Add”按钮，保存这个新数据源。

下面显示了一个Prometheus数据源配置例子：
![Prometheus configuration in Grafana](https://prometheus.io/assets/grafana_configuring_datasource-cb0e78b7cfa.png)

#### 创建一个Prometheus Graph图表
下面是添加一个新的Grafana的标准方法：
 1. 点击图表Graph的title，它在图表上方中间。然后点击“Edit”。
 2. 在“Metrics”tab下面，选择你的Prometheus数据源（下面右边）。
 3. 在“Query”字段中输入你想查询的Prometheus表达式，同时使用“Metrics”字段通过自动补全查找度量指标。
 4. 为了格式化时间序列的图例名称，使用“Legend format”图例格式输入。例如，为了仅仅显示这个标签为`method`和`status`的查询结果，你可以使用图例格式`{{method{} - {{status}}`。
 5. 调节其他的Graph设置，知道你有一个工作图表。

下面显示了一个Prometheus图表配置：
![Prometheus图表](https://prometheus.io/assets/grafana_qps_graph-cb702994700.png)

#### 从Grafana.net导入预构建的dashboard
Grafana.net维护一个共享dashboard的收集，它们能够被下载，并在Grafana服务中使用。使用Grafana.net的“Filter”选项去浏览来自Prometheus数据源的dashboards

你当前必须手动编辑下载下来的JSON文件和更改`datasource`: 选择Prometheus服务作为Grafana的数据源，使用“Dashboard”->"Home"->"Import"选项去导入编辑好的dashboard文件到你的Grafana中。

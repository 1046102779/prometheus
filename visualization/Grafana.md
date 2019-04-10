[Grafana](http://grafana.org/)支持查询Prometheus。从Grafana 2.5.0 (2015-10-28)开始Prometheus可以作为它的数据源。

以下显示了一个示例Grafana仪表板，它向Prometheus查询数据：
![Deme Dashboard](https://prometheus.io/assets/grafana_prometheus.png)

##### 一、Grafan安装
Grafana的完整安装教程，详见[Grafana官方文档](http://docs.grafana.org/installation/)

例如，在Linux上，安装Grafana可能如下所示：
```Grafana install
# Download and unpack Grafana from binary tar (adjust version as appropriate).
curl -L -O https://grafanarel.s3.amazonaws.com/builds/grafana-2.5.0.linux-x64.tar.gz
tar zxf grafana-2.5.0.linux-x64.tar.gz

# Start Grafana.
cd grafana-2.5.0/
./bin/grafana-server web
```

##### 二、使用方法
默认情况下，Grafana将监听[http://localhost:3000](http://localhost:3000)。默认登录用户名和密码“admin/admin”。

###### 2.1 创建一个Prometheus数据源

创建一个Prometheus数据源Data source：
 1. 点击Grafana的logo，打开工具栏。
 2. 在工具栏中，点击"Data Source"菜单。
 3. 点击"Add New"。
 4. 数据源Type选择“Prometheus”。
 5. 设置Prometheus服务访问地址（例如：`http://localhost:9090`）。
 6. 调整其他想要的设置（例如：关闭代理访问）。
 7. 点击“Add”按钮，保存这个新数据源。

下面显示了一个Prometheus数据源配置例子：
![Prometheus configuration in Grafana](https://prometheus.io/assets/grafana_configuring_datasource.png)

###### 2.2 创建一个Prometheus Graph图表
按照添加新Grafana图的标准方式。 然后：
 1. 点击图表Graph的title，它在图表上方中间。然后点击“Edit”。
 2. 在“Metrics”tab下面，选择你的Prometheus数据源（下面右边）。
 3. 在“Query”字段中输入你想查询的Prometheus表达式，同时使用“Metrics”字段通过自动补全查找度量指标。
 4. 为了格式化时间序列的图例名称，使用“Legend format”图例格式输入。例如，为了仅仅显示这个标签为`method`和`status`的查询结果，你可以使用图例格式`{{method{} - {{status}}`。
 5. 调节其他的Graph设置，知道你有一个工作图表。

以下显示了Prometheus图配置示例：
![Prometheus图表](https://prometheus.io/assets/grafana_qps_graph.png)

###### 2.3 从Grafana.net导入预构建的dashboard
Grafana.com维护着一组[共享仪表板](https://grafana.com/dashboards)，可以下载并与Grafana的独立实例一起使用。 使用Grafana.com“过滤器”选项仅浏览“Prometheus”数据源的仪表板。

您当前必须手动编辑下载的JSON文件并更正`datasource`：条目以反映您为Prometheus服务器选择的Grafana数据源名称。 使用“仪表板”→“主页”→“导入”选项将已编辑的仪表板文件导入Grafana安装。

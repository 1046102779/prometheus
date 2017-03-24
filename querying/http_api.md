## HTTP API
---
在Prometheus服务上`/api/v1`版本api是稳定版。

### 格式概述
这个API返回是JSON格式。每个请求成功的返回值都是以`2xx`开头的编码。

到达API处理的无效请求，返回一个JSON错误对象，并返回下面的错误码：
  - `400 Bad Request`。当参数错误或者丢失时。
  - `422 Unprocessable Entity`。当一个表达式不能被执行时。
  - `503 Service Unavailable`。当查询超时或者中断时。

在请求到达API之前，其他非`2xx`的错误码可能会被返回。

JSON返回格式如下所示：
```JSON
{
"status": "success" | "error",
"data": <data>,

// 如果status是"error", 这个数据字段还会包括下面的数据
"errorType": "<string>",
"error": "<string>"
}
```

输入时间戳可以被RFC3339格式或者Unix时间戳提供。输出时间戳以Unix时间戳的方式呈现。

查询参数名称可以用`[]`中括号重复次数。
`\<series_selector\>`占位符提供像`http_requests_total`或者`http_requests_total{method=~"^GET|POST$"}`的Prometheus时间序列选择器，并需要在URL中编码传输。

`\<duration\>`占位符涉及到`[0-9]-[smhdwy]`。例如：`5m`表示5分钟的持续时间。

### 表达式查询
查询语言表达式可以是表示一个瞬时向量值，或者一个范围向量值。

#### Instant queries(即时查询)
下面这个在该时刻的即时查询方式：
> GET /api/v1/query

URL查询参数：
 - `query=\<string\>`: Prometheus表达式查询字符串。
 - `time=\<rfc3339 | uninx_timestamp\>`: 评估时间戳，可选项。

如果`time`缺省，则用当前服务器时间表示即时时刻。

这个查询结果的`data`部分有下面格式：
```JSON
{
 "resultType": "matrix" | "vector" | "scalar" | "string",
 "result": <value>
}
```

\<value\>提供一个查询结果数据，依赖于这个`resultType`有很多格式。见[表达式查询结果格式](https://prometheus.io/docs/querying/api/#expression-query-result-formats)。

下面例子评估了在时刻是`2015-07-01T20:10:51.781Z`的`up`表达式：
```JSON
$ curl 'http://localhost:9090/api/v1/query?query=up&time=2015-07-01T20:10:51.781Z'
{
 "status": "success",
 "data":{
    "resultType": "vector",
    "result" : [
         {
            "metric" : {
               "__name__" : "up",
               "job" : "prometheus",
               "instance" : "localhost:9090"
            },
            "value": [ 1435781451.781, "1" ]
         },
         {
            "metric" : {
               "__name__" : "up",
               "job" : "node",
               "instance" : "localhost:9100"
            },
            "value" : [ 1435781451.781, "0" ]
         }
    ]
 }
}
```

#### 范围查询
下面评估了一个范围时间的查询表达式：
> GET /api/v1/query_range

URL查询参数
 - `query=\<string\>`: Prometheus表达式查询字符串。
 - `start=\<rfc3339 | unix_timestamp\>`: 开始时间戳。
 - `end=\<rfc3339 | unix_timestamp\>`: 结束时间戳。
 - `step=\<duration\>`: 查询步长。

下面查询结果格式的`data`部分：
```json
{
    "resultType": "matrix",
    "result": \<value\>
}
```

对于`\<value\>`占位符的格式，详见[范围向量结果格式](https://prometheus.io/docs/querying/api/#range-vectors)。

下面例子评估的查询条件`up`，且30s范围的查询，步长是15s。
```JSON
$ curl 'http://localhost:9090/api/v1/query_range?query=up&start=2015-07-01T20:10:30.781Z&end=2015-07-01T20:11:00.781Z&step=15s'
{
   "status" : "success",
   "data" : {
      "resultType" : "matrix",
      "result" : [
         {
            "metric" : {
               "__name__" : "up",
               "job" : "prometheus",
               "instance" : "localhost:9090"
            },
            "values" : [
               [ 1435781430.781, "1" ],
               [ 1435781445.781, "1" ],
               [ 1435781460.781, "1" ]
            ]
         },
         {
            "metric" : {
               "__name__" : "up",
               "job" : "node",
               "instance" : "localhost:9091"
            },
            "values" : [
           [ 1435781430.781, "0" ],
               [ 1435781445.781, "0" ],
               [ 1435781460.781, "1" ]
            ]
         }
      ]
   }
}
```

### 查询元数据
#### 通过标签匹配器找到时间序列
下面例子返回了匹配时间序列数据, 且不返回时间序列数据值。
> GET /api/v1/series

URL查询参数：
 - `match[]=\<series_selector\>`: 选择器是series_selector。这个参数个数必须大于等于1.
 - `start=\<rfc3339 | unix_timestamp\>`: 开始时间戳。
 - `end=\<rfc3339 | unix_timestamp\>`: 结束时间戳。

返回结果的`data`部分，是由key-value键值对的对象列表组成的。

下面这个例子返回时间序列数据, 选择器是`up`或者`process_start_time_seconds{job="prometheus"}`
```JSON
$ curl -g 'http://localhost:9090/api/v1/series?match[]=up&match[]=process_start_time_seconds{job="prometheus"}'
{
   "status" : "success",
   "data" : [
      {
         "__name__" : "up",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      },
      {
         "__name__" : "up",
         "job" : "node",
         "instance" : "localhost:9091"
      },
      {
         "__name__" : "process_start_time_seconds",
         "job" : "prometheus",
         "instance" : "localhost:9090"
      }
   ]
}
```

#### 查询标签值
下面这个例子，返回了带有指定标签和标签值
> GET /api/v1/label/\<label_name\>/values

这个返回JSON结果的`data`部分是带有label_name=job的值列表：
```JSON
$ curl http://localhost:9090/api/v1/label/job/values
{
   "status" : "success",
   "data" : [
      "node",
      "prometheus"
   ]
}
```

### 删除时间序列
下面的例子，是从Prometheus服务中删除所有的时间序列数据：
> DELETE /api/v1/series

URL查询参数
 - `match[]=\<series_selector\>`: 删除符合series_selector匹配器的时间序列数据。参数个数必须大于等于1.

返回JSON数据中的`data`部分有以下的格式
> {
>    "numDeleted": \<number of deleted series\>
> }

下面的例子删除符合度量指标名称是`up`或者时间序列为`process_start_time_seconds{job="prometheus"}`：
```JSON
$ curl -XDELETE -g 'http://localhost:9090/api/v1/series?match[]=up&match[]=process_start_time_seconds{job="prometheus"}'
{
   "status" : "success",
   "data" : {
      "numDeleted" : 3
   }
}
```

### 表达式查询结果格式
表达式查询结果，在`data`部分的`result`部分中，返回下面的数据。`\<sample_value\>`占位符有数值样本值。JSON不支持特殊浮点值，例如：`NaN`, `Inf`和`-Inf`。因此样本值返回结果是字符串，不是原生的数值。

#### 范围向量
范围向量返回的result类型是一个`matrix`矩阵。下面返回的结果是`result`部分的数据格式：
```JSON
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "values": [ [ <unix_time>, "<sample_value>" ], ... ]
  },
  ...
]
```

#### 瞬时向量
瞬时向量的`result`类型是`vector`。下面是`result`部分的数据格式
```JSON
[
  {
    "metric": { "<label_name>": "<label_value>", ... },
    "value": [ <unix_time>, "<sample_value>" ]
  },
  ...
]
```

#### Scalars标量
标量查询返回`result`类型是`scalar`。下面是`result`部分的数据格式：
> [ \<unix_time\>, "\<scalar_value\>" ]

#### 字符串
字符串的`result`类型是`string`。下面是`result`部分的数据格式：
> [ \<unix_time\>, "\<string_value\>" ]

### Targets目标
这个API是实验性的，暂不翻译。

### Alertmanagers
这个API也是实验性的，暂不翻译

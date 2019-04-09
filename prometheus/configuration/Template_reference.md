Prometheus支持在警报的注释和标签以及服务的控制台页面中进行模板化。 模板能够针对本地数据库运行查询，迭代数据，使用条件，格式化数据等.Prometheus模板语言基于Go模板系统。
##### 一、数据结构题
处理时间序列数据的主要数据结构是样本，定义如下：
```
type sample struct {
        Labels map[string]string
        Value  float64
}
```
样本的度量标准名称在`Labels`映射中的特殊`__name__`标签中进行编码。

`[]sample`表示样本列表。

Go中的`interface{}`类似于C中的void指针。
##### 二、函数
除了Go模板提供的[默认功能](https://golang.org/pkg/text/template/#hdr-Functions)外，Prometheus还提供了更轻松处理模板中查询结果的功能。

如果在管道中使用函数，则管道值将作为最后一个参数传递。
###### 2.1 查询
名字 |	参数 |	返回值 |解析
---|---|---|---
query|	query string	|[]sample	|查询数据库，不支持返回范围向量。
first|	[]sample|	sample	|索引等于0
label|	label, sample	|string	|相当于`index sample.Labels标签`
value|	sample	|float64|	相当于 `sample.Value`
sortByLabel|	label, []samples|	[]sample|	按给定标签对样品进行排序。 是稳定排序。
`first`，`label`和`value`旨在使查询结果易于在管道中使用。
###### 2.2 数字
名字|	参数|	返回|	解析
---|---|---|---
humanize	|number|	string	|使用度量标准前缀将数字转换为更易读的格式。
humanize1024|	number	|string	|像`humanize`一样，但使用1024作为基础而不是1000。
humanizeDuration|	number	|string|	将持续时间（以秒为单位）转换为更易读的格式。
humanizeTimestamp|	number	|string	|将Unix时间戳以秒为单位转换为更易读的格式。
`Humanizing `功能旨在为人类消费产生合理的输出，并且不保证在Prometheus版本之间返回相同的结果。
###### 2.3 字符串
名字|	参数|	返回|	解析
---|---|---|---
title|	string	|string	|`strings.Title`, 大写每个单词的第一个字符。
toUpper|	string	|string	|`strings.ToUpper`, 将所有字符转换为大写。
toLower	|string	|string	|`strings.ToLower`, 将所有字符转换为小写。
match|	pattern, text	|boolean	|`regexp.MatchString` 测试未锚定的正则表达式匹配。
reReplaceAll|	pattern, replacement, text|	string	|`Regexp.ReplaceAllString` Regexp替换，未经修复。
graphLink|	expr	|string	|返回表达式的表达式浏览器中图表视图的路径。
tableLink|	expr	|string|	返回表达式的表达式浏览器中表格（“Console”）视图的路径。
###### 2.4 其他
名字|	参数|	返回|	解析
---|---|---|---
args|	[]interface{}|	map[string]interface{}|	这会将对象列表转换为具有键`arg0`，`arg1`等的映射。这旨在允许将多个参数传递给模板。
tmpl|	string, []interface{}	|nothing	|与内置模板一样，但允许非文字作为模板名称。 请注意，结果被认为是安全的，不会自动转义。 仅适用于游戏机。
safeHtml|	string	|string|	将字符串标记为不需要自动转义的HTML。

##### 三、模板类型的区别
每种类型的模板都提供可用于参数化模板的不同信息，并具有一些其他差异。
###### 3.1 报警字段模板
`.Value`和`.Labels`包含警报值和标签。 为方便起见，它们也作为`$value`和`$labels`变量公开。
###### 3.2 控制台模板
控制台暴露在`/consoles/`上，并且来自`-web.console.templates`标志指向的目录。

控制台模板使用[html/template](https://golang.org/pkg/html/template/)呈现，提供自动转义功能。 要绕过自动转义，请使用`safe*`功能。，

URL参数在`.Params`中以地图形式提供。 要使用相同的名称访问多个URL参数，`.RawParams`是每个参数的列表值的映射。 URL路径在`.Path`中可用，不包括`/consoles/`前缀。

控制台还可以访问在`-web.console.libraries`标志指向的目录中的* `.lib`文件中找到的`{{define"templateName"}}...{{end}}`定义的所有模板。 由于这是一个共享命名空间，请注意避免与其他用户发生冲突。 以`prom`，`_prom`和`__`开头的模板名称保留供Prometheus使用，上面列出的函数也是如此。

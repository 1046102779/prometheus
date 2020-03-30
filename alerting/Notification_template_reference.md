Prometheus创建并向Alertmanager发送警报，然后Alertmanager根据标签向不同的接收者发送通知。 接收器可以是众多集成中的一种，包括：Slack，PagerDuty，电子邮件或通过通用webhook接口的自定义集成。

发送给接收者的通知是通过模板构建的。 Alertmanager附带默认模板，但也可以自定义。 为避免混淆，请务必注意Alertmanager模板与Prometheus中的模板不同，但Prometheus模板还包括警报规则标签/注释中的模板。

Alertmanager的通知模板基于Go模板系统。 请注意，某些字段将作为文本进行评估，而其他字段将作为HTML进行评估，这将影响转义。

##### 一、数据
`Data`是传递给通知模板和webhook推送的结构。

Name |	Type|	Notes
---|---|---
Receiver|	string	|定义通知将被发送到的接收者名称（松弛，电子邮件等）。
Status	|string| 如果至少有一个警报被触发，则定义为触发，否则解析。	
Alerts	|Alert|	警报对象列表（见下文）。
GroupLabels	|KV	|这些警报的标签按分组。
CommonLabels|	KV	|所有警报共有的标签。
CommonAnnotations|	KV	|所有警报的常用注释集。 用于有关警报的更长的其他信息串。
ExternalURL	|string|反向链接到发送通知的Alertmanager。

##### 二、警报
`Alert`为通知模板保留一个警报。

Name|	Type|	Notes
---|---|---
Status	| string |	定义警报是已解决还是当前正在触发。
Labels	| KV	|要附加到警报的一组标签。
Annotations	|KV|	警报的一组注释。
StartsAt|	time.Time|	警报开始发射的时间。 如果省略，则当前时间由Alertmanager分配。
EndsAt	|time.Time	|仅在已知警报结束时间时设置。 否则设置为自上次收到警报以来的可配置超时时间。
GeneratorURL|	string	|一个反向链接，用于标识此警报的生成实体。

##### 三、KV
`KV`是一组用于表示标签和注释的键/值字符串对。
```
type KV map[string]string
```
包含两个注释的注释示例：
```
{
  summary: "alert summary",
  description: "alert description",
}
```
除了直接访问存储为KV的数据（标签和注释）之外，还有用于排序，删除和查看LabelSet的方法：

KV方法

Name|	Arguments|	Returns|	Notes
---|---|---|---
SortedPairs	|-	|Pairs (list of key/value string pairs.)|返回键/值对的排序列表。	
Remove|	[]string|	KV|	返回没有给定键的键/值映射的副本。
Names|	-	|[]string|	返回LabelSet中标签名称的名称。
Values|	-	|[]string|	返回LabelSet中的值列表。

##### 四、函数
请注意Go模板也提供的[默认函数](https://golang.org/pkg/text/template/#hdr-Functions)。

Name	|Arguments	|Returns	
---|---|---
title	|string	|strings.Title, 大写每个单词的第一个字符。
toUpper	|string|	strings.ToUpper, 将所有字符转换为大写。
toLower	|string|	strings.ToLower, 将所有字符转换为小写。
match	|pattern, string	|Regexp.MatchString. 使用Regex匹配字符串。
reReplaceAll|	pattern, replacement, text|	Regexp.ReplaceAllString Regexp替换，未经修复。	
join	|sep string, s []string	|strings.Join, 连接s的元素以创建单个字符串。 分隔符字符串sep放在结果字符串中的元素之间。 （注意：参数顺序已反转，以便在模板中更容易管道化。）
safeHtml|	text string|	html/template.HTML, 将字符串标记为不需要自动转义的HTML。

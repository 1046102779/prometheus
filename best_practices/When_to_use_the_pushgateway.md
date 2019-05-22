## WHEN TO USE THE PUSHGATEWAY
---
Pushgateway是一个中介服务，允许您从不能被刮除的作业推出指标。 有关详细信息，请参阅[推送指标](https://prometheus.io/docs/instrumenting/pushing/)。

### 我应该用Pushgateway吗？
我们只建议在某些有限的情况下使用Pushgateway。 盲目使用Pushgateway代替Prometheus通常的拉动模型进行一般指标收集时，有几个陷阱：
 - 通过单个Pushgateway监控多个实例时，Pushgateway成为单点故障和潜在瓶颈。
 - 您通过up metric（在每次抓取时生成）丢失了Prometheus的自动实例运行状况监视。
 - Pushgateway永远不会忘记推送到它的系列，并将永远暴露给Prometheus，除非这些系列是通过Pushgateway的API手动删除的。
当作业的多个实例通过实例标签或类似物区分其在Pushgateway中的度量时，后一点尤其重要。即使重命名或删除了原始实例，实例的度量标准仍将保留在Pushgateway中。这是因为Pushgateway作为指标缓存的生命周期基本上与将指标推送到它的流程的生命周期分开。将此与普罗米修斯通常的拉式监控进行对比：当实例消失（有意或无意）时，其指标会随之自动消失。使用Pushgateway时，情况并非如此，您现在必须手动删除任何过时的指标或自行自动执行此生命周期同步。

通常，Pushgateway唯一有效的用例是捕获服务级批处理作业的结果。 “服务级”批处理作业是与特定计算机或作业实例在语义上无关的作业（例如，删除整个服务的多个用户的批处理作业）。此类作业的度量标准不应包含计算机或实例标签，以将特定计算机或实例的生命周期与推送的度量标准分离。这减少了管理Pushgateway中过时指标的负担。另请参阅[监视批处理作业的最佳实践](https://prometheus.io/docs/practices/instrumentation/#batch-jobs)

### 替代策略
如果入站防火墙或NAT阻止您从目标中提取指标，请考虑将Prometheus服务器也移到网络屏障后面。 我们通常建议在与受监控实例相同的网络上运行Prometheus服务器。

对于与计算机相关的批处理作业（例如自动安全更新cronjobs或配置管理客户端运行），使用[Node exporter](https://github.com/prometheus/node_exporter)的文本文件模块而不是Pushgateway公开生成的指标。

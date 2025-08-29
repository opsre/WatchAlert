use watchalert;
INSERT ignore INTO watchalert.rule_template_groups (name,`number`,description,`type`) VALUES
	 ('APISIX',0,'APISIX指标监控','Metrics'),
	 ('服务日志监控',0,'服务日志监控','Logs'),
	 ('Docker',0,'Docker容器监控','Metrics'),
	 ('ElasticSearch',0,'ElasticSearch资源监控','Metrics'),
	 ('ETCD',0,'ETCD','Metrics'),
	 ('Jaeger',0,'Jaeger链路监控','Traces'),
	 ('Kafka',0,'Kafka监控','Metrics'),
	 ('Kubernetes',0,'Kubernetes事件监控','Events'),
	 ('KubernetesMetric',0,'Kubernetes指标监控','Metrics'),
	 ('MongoDB',0,'MongoDB监控','Metrics');
INSERT ignore INTO watchalert.rule_template_groups (name,`number`,description,`type`) VALUES
	 ('MySQL',0,'MySQL资源监控','Metrics'),
	 ('Node节点监控',0,'Node服务器监控','Metrics'),
	 ('Redis',0,'Redis资源监控','Metrics'),
	 ('RocketMQ',0,'RocketMQ监控','Metrics');

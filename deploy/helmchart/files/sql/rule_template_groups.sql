use watchalert;
INSERT ignore INTO watchalert.rule_template_groups (name,`number`,`type`) VALUES
	 ('APISIX',0,'Metrics'),
	 ('服务日志监控',0,'Logs'),
	 ('Docker',0,'Metrics'),
	 ('ElasticSearch',0,'Metrics'),
	 ('ETCD',0,'Metrics'),
	 ('Jaeger',0,'Traces'),
	 ('Kafka',0,'Metrics'),
	 ('Kubernetes',0,'Events'),
	 ('KubernetesMetric',0,'Metrics'),
	 ('MongoDB',0,'Metrics');
INSERT ignore INTO watchalert.rule_template_groups (name,`number`,`type`) VALUES
	 ('MySQL',0,'Metrics'),
	 ('Node节点监控',0,'Metrics'),
	 ('Redis',0,'Metrics'),
	 ('RocketMQ',0,'Metrics');

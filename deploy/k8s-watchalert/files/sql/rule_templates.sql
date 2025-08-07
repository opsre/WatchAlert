use watchalert;
INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'APISIX',
        'APISIX ETCD DB Total Size过大',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(etcd_mvcc_db_total_size_in_bytes / etcd_server_quota_backend_bytes) * 100',
            'annotations', '实例：${instance}，ETCD 数据库过大可能导致查询和写入操作的延迟增加，影响 APISIX 的性能。 数据库过大还可能增加 ETCD 集群的负载，导致资源消耗增加，甚至可能触发 OOM（内存溢出）等问题。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'APISIX',
        'APISIX ETCD leader变更频繁',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'changes(etcd_server_leader_changes_seen_total{}[1h])',
            'annotations', '实例：${instance}，频繁的 leader 变更会导致配置更新延迟和性能下降。 可能引发 ETCD 的可用性问题，从而影响 APISIX 的配置同步和运行稳定性。 可能是由于网络抖动、节点故障或资源瓶颈等问题导致的，需要进一步排查。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>5')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'APISIX',
        'APISIX ETCD 缺失leader',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(etcd_server_is_leader{})',
            'annotations', '实例：${instance}，当缺少 leader 时，ETCD 集群将无法处理写操作，这会导致 APISIX 无法进行配置更新。 配置无法更新可能导致流量调度、插件配置等关键功能失效，直接影响业务的运行。 这种情况通常是由于 ETCD 集群的严重故障或节点失联导致的，需及时修复。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '!=1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'APISIX',
        'APISIX ETCD连接状态',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'apisix_etcd_reachable',
            'annotations', '实例：${instance}，APISIX 无法从 Etcd 读取或写入配置信息，可能导致路由规则更新失败、服务发现异常等问题，进而影响整个 API 网关的正常运行。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '!=1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'APISIX',
        'APISIX 每秒请求量过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'increase(apisix_http_requests_total{}[1m]) / 60',
            'annotations', '实例：${instance}，生产环境APISix服务正在经历高请求速率。每秒超过1000个请求。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>1000')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'APISIX',
        'APISIX 活跃连接接近峰值',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(apisix_nginx_http_current_connections{state="active"})',
            'annotations', '实例：${instance}，连接数超过apisix工作节点的worker connections总数后，网关会停止工作，新的请求会被拒绝，导致业务受影响。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>135313')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Docker',
        'Docker ContainerKilled',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'time() - container_last_seen',
            'annotations', '容器：${instance}，被意外杀死，请注意排查！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>60')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Docker',
        'Docker容器CPU使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum(rate(container_cpu_usage_seconds_total{container!=""}[5m])) by (pod, container) / sum(container_spec_cpu_quota{container!=""}/container_spec_cpu_period{container!=""}) by (pod, container) * 100)',
            'annotations', '容器：${instance}，CPU使用率过高，当前：${value}%，请注意排查！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Docker',
        'Docker容器内存使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum(container_memory_working_set_bytes{name!=""}) BY (instance, name) / sum(container_spec_memory_limit_bytes > 0) BY (instance, name) * 100)',
            'annotations', '容器：${instance}，内存使用率过高，当前：${value}%，请注意排查！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Docker',
        'Docker容器存储使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(1 - (sum(container_fs_inodes_free{name!=""}) BY (instance) / sum(container_fs_inodes_total) BY (instance))) * 100',
            'annotations', '容器：${instance}，存储使用率过高，当前：${value}%，请注意排查！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'ElasticSearch',
        'ES慢查询告警',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'rate(elasticsearch_indices_search_query_time_seconds[5m])',
            'annotations', '实例：${instance}，当前：${value}，慢查询占用较多计算和存储资源，可能导致其他查询任务响应变慢，影响系统的整体性能。 频繁的慢查询会增加 Elasticsearch 集群的负载，导致 CPU、内存和磁盘 I/O 的压力增大，甚至导致集群不稳定。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>0.2')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES节点CPU使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(elasticsearch_process_cpu_percent{})',
            'annotations', '实例：${instance}，当前：${value}，高 CPU 使用率可能是由于高查询或写入负载、复杂的查询操作、数据压缩或解压引起的，可能导致节点性能下降、响应时间变长，甚至引发故障。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES节点堆内存使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(((elasticsearch_jvm_memory_used_bytes{area="heap"} / elasticsearch_jvm_memory_max_bytes{area="heap"}) * 100 ) * 100) / 100',
            'annotations', '实例：${instance}，当前：${value}，内存使用率过高可能导致频繁的垃圾回收 (GC)，增加延迟，甚至导致“out of memory” (OOM) 错误，使节点宕机。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES节点状态',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'elasticsearch_node_stats_up',
            'annotations', '实例：${instance}，当前：${value}，节点状态不等于1可能是节点出现down机的情况，会影响到整个集群的稳定性。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '!=1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES节点磁盘使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round((100 - (elasticsearch_filesystem_data_available_bytes{}/elasticsearch_filesystem_data_size_bytes{} * 100)) * 100 ) / 100',
            'annotations', '实例：${instance}，当前：${value}，磁盘空间不足会导致分片无法写入或重新分配，影响集群的正常运行，可能导致集群进入黄色或红色状态。超过 90% 时，集群可能会停止写入。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES集群写入QPS大于50',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'rate(elasticsearch_indices_indexing_index_total[5m])',
            'annotations', '实例：${instance}，当前：${value}，写入 QPS 过高可能导致写入延迟增加、集群性能下降，甚至引发写入失败。大量写入请求可能导致分片分配不均、节点负载过高等问题。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>50')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES集群查询QPS大于50',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'rate(elasticsearch_indices_search_query_total[5m])',
            'annotations', '实例：${instance}，当前：${value}，查询 QPS 过高可能导致查询延迟增加、响应时间变长，影响查询性能。集群的查询负载增加，可能造成节点资源紧张。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>50')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES集群状态为红色',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'elasticsearch_cluster_health_status{color="red"}',
            'annotations', '实例：${instance}，红色状态表示集群中有主分片和副本分片不可用，数据丢失或无法访问。此时，部分或全部数据不可用，可能影响业务的正常运行。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ElasticSearch',
        'ES集群状态为黄色',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'elasticsearch_cluster_health_status{color="yellow"}',
            'annotations', '实例：${instance}，黄色状态表示集群中的主分片可用，但副本分片有问题，通常是由于资源不足、分片分配失败等原因。虽然数据可以读取，但副本的不可用会降低数据冗余性和容灾能力。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ETCD',
        'Etcd大量失败的 GRPC 请求',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(rate(grpc_server_handled_total{grpc_code!="OK"}[1m])) BY (grpc_service, grpc_method) / sum(rate(grpc_server_handled_total[1m])) BY (grpc_service, grpc_method)',
            'annotations', 'Etcd 中检测到超过 1% 的 GRPC 请求失败，当前：${value}%',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>0.01')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'ETCD',
        'Etcd大量失败的 HTTP 请求',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(rate(etcd_http_failed_total[1m])) BY (method) / sum(rate(etcd_http_received_total[1m])) BY (method)',
            'annotations', 'Etcd 中检测到超过 1% 的 HTTP 故障，当前：${value}%',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>0.01')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ETCD',
        'Etcd成员不足',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'count(etcd_server_id) % 2',
            'annotations', 'Etcd 集群应具有奇数个成员，当前成员个数：${value}',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ETCD',
        'Etcd无Leader',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'etcd_server_has_leader',
            'annotations', 'Etcd 集群没有领导者',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'ETCD',
        'Etcd频繁发生领导者变更',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'increase(etcd_server_leader_changes_seen_total[10m])',
            'annotations', 'Etcd 领导者在 10 分钟内变更超过 2 次。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>2')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kafka',
        'KafkaTopic偏移量过低',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'rate(kafka_topic_partition_current_offset{topic!=""}[2m])',
            'annotations', 'Kafka集群偏移量异常，当前：${value}。可能存在数据异常问题，请观察!',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '< 1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kafka',
        'KafkaTopic副本丢失',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'count(kafka_topic_partition_under_replicated_partition{topic!="",partition="0"}==0) by(topic,job)',
            'annotations', 'Kafka集群节点缺失，当前：${value}个副本。当一个分区的某个副本所在的 broker 宕机、网络不可达、数据损坏或由于其他原因导致副本不可用时，就会发生副本丢失。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '<3')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kafka',
        'Kafka节点离线告警',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'kafka_brokers',
            'annotations', 'Kafka集群节点缺失，当前：${value}个节点。Kafka 集群中的某个节点下线，可能会导致部分分区数据无法读取或写入，并可能导致消费者或生产者请求失败。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '<3')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'KubernetesAPI 服务器延迟',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'histogram_quantile(0.99, sum(rate(apiserver_request_duration_seconds_bucket{verb!~"(?:CONNECT|WATCHLIST|WATCH|PROXY)"} [10m])) WITHOUT (subresource))',
            'annotations', 'Kubernetes API 服务器对于 ${verb} ${resource} 的第 99 个百分位延迟为 ${value} 秒。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesHPA获取外部指标失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'HPA', 'reason', 'FailedGetExternalMetric', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesHPA获取资源指标失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'HPA', 'reason', 'FailedGetResourceMetric', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'Kubernetes',
        'KubernetesHPA调整副本数失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'HPA', 'reason', 'FailedRescale', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesNode内存压力过大',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Nodes', 'reason', 'NodeUnderMemoryPressure', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesNode磁盘压力过大',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Nodes', 'reason', 'NodeUnderDiskPressure', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesNode节点不可用',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Nodes', 'reason', 'NodeNotReady', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod不健康',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'Unhealthy', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod反覆崩溃',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'CrashLoopBackOff', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod启动失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'Failed', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod挂载卷失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'FailedMount', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod调度失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'FailedScheduling', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPod超出运行期限',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'DeadlineExceeded', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'Kubernetes',
        'KubernetesPod附加卷到节点失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'Pods', 'reason', 'FailedAttachVolume', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Kubernetes',
        'KubernetesPV/PVC绑定失败',
        'KubernetesEvent',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Events',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', 'PVC/PV', 'reason', 'FailedBinding', 'value', 1, 'filter', NULL, 'scope', 10),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'Kubernetes客户端证书下周到期',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'apiserver_client_certificate_expiration_seconds_count{job="apiserver"} > 0 and histogram_quantile(0.01, sum by (job, le) (rate(apiserver_client_certificate_expiration_seconds_bucket{job="apiserver"}[5m])))',
            'annotations', 'K8s集群证书下周到期，请立即处理。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '<604800')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'Kubernetes节点网络不可用',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'kube_node_status_condition{condition="NetworkUnavailable",status="true"}',
            'annotations', '节点 ${instance}，存在 NetworkUnavailable 情况',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'Kubernetes节点超出 pod 容量',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum by (node) ((kube_pod_status_phase{phase="Running"} == 1) + on(uid, instance) group_left(node) (0 * kube_pod_info{pod_template_hash=""})) / sum by (node) (kube_node_status_allocatable{resource="pods"}) * 100',
            'annotations', '节点 ${instance}，超出 pod 容量',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>110')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MongoDB',
        'MongoDB副本成员复制存在延迟',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'mongodb_mongod_replset_member_optime_date{state="PRIMARY"} - on (set) group_right mongodb_mongod_replset_member_optime_date{state="SECONDARY"}',
            'annotations', 'MongoDB实例：${instance} ，复制延迟过高，当前：${value}。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>10')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MongoDB',
        'MongoDB副本成员状态异常',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'mongodb_mongod_replset_member_health',
            'annotations', '实例：${instance} 不健康，如果副本集中有节点健康状态异常，可能会导致数据同步延迟、读写操作失败，甚至可能引发服务不可用的问题。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '!=1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MongoDB',
        'MongoDB节点宕机',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'mongodb_up',
            'annotations', 'MongoDB实例：${instance} ，节点宕机。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MongoDB',
        'MongoDB连接数过多',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round((max(mongodb_connections{state="current"}) by (instance,job) / sum(mongodb_connections{state=~"current|available"}) by (instance,job) * 100) * 100) / 100',
            'annotations', '实例：${instance} 连接数过多，当连接数达到上限时，新的客户端连接请求将被拒绝，进而影响应用程序的正常运行，可能引发严重的系统不可用问题。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>50')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MySQL',
        'MySQL存在慢查询',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'increase(mysql_global_status_slow_queries[1m])',
            'annotations', 'MySQL实例：${instance}，存在慢查询，当前条数：${value}。慢查询通常表明数据库中存在未优化的 SQL 操作，可能导致查询性能下降，影响整体响应时间。需要通过优化查询语句、添加索引、调整表结构或优化数据库配置来解决慢查询问题。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'MySQL',
        'MySQL宕机',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'mysql_up',
            'annotations', 'MySQL实例：${instance}，意外宕机。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MySQL',
        'MySQL实例发生重启',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'mysql_global_status_uptime',
            'annotations', 'MySQL实例：${instance}，发生重启。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '<60')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MySQL',
        'MySQL连接数过多',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'max_over_time(mysql_global_status_threads_connected[1m]) / mysql_global_variables_max_connections * 100',
            'annotations', 'MySQL实例：${instance}，连接数过多，当前：${value}。连接数接近上限意味着数据库负载很高，可能导致连接请求被拒绝，影响业务操作。 连接数耗尽时，新请求会被拒绝，导致数据库不可用。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'MySQL',
        'MySQL高线程运行',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'max_over_time(mysql_global_status_threads_running[1m]) / mysql_global_variables_max_connections * 100',
            'annotations', '实例：${instance}，当前：${value}，高线程数可能导致数据库负载增加，需要检查数据库的并发操作、慢查询、事务等待等问题，优化线程池配置和查询效率。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>60')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'Pod CPU使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round((sum(rate(container_cpu_usage_seconds_total{namespace="kube-system",container!~"^POD$"}[5m])) by (namespace,pod,container) / sum(kube_pod_container_resource_limits{namespace="kube-system",resource="cpu"}) by (container,namespace,pod) * 100) * 100) / 100',
            'annotations', 'Namespace：${namespace}\nPod：${pod}\nCPU使用率过高，当前：${value}%',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'KubernetesMetric',
        'Pod 内存使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(((avg by (namespace,pod,container) (container_memory_usage_bytes{namespace="kube-system",container!~"^POD$"}) / avg by (namespace,pod,container) (container_spec_memory_limit_bytes{namespace="kube-system",container!=""})) * 100) * 100) / 100',
            'annotations', 'Namespace：${namespace}\nPod：${pod}\n内存使用率过高，当前：${value}%',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Redis',
        'Redis主节点过多',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'count(redis_instance_info{role="master"})',
            'annotations', 'Redis实例：${instance}，集群中有太多节点标记为主节点。！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Redis',
        'Redis实例宕机',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'redis_up',
            'annotations', 'Redis实例：${instance}，发生宕机，请跟进排查！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Redis',
        'Redis拒绝连接',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'increase(redis_rejected_connections_total[1m])',
            'annotations', 'Redis实例：${instance}，一些与 Redis 的连接已被拒绝。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Redis',
        'Redis系统内存不足',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'redis_memory_used_bytes / redis_total_system_memory_bytes * 100',
            'annotations', 'Redis实例：${instance}，系统内存不足（> 80％），当前：${value}%。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'Redis',
        'Redis缺少主节点',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(count(redis_instance_info{role="master"}) or vector(0))',
            'annotations', 'Redis实例：${instance}，集群没有标记为主节点！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '<1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Redis',
        'Redis连接数过多',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'redis_connected_clients / redis_config_maxclients * 100',
            'annotations', 'Redis实例：${instance}，连接数不足（使用率 > 90%），当前：${value}%。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'RocketMQ',
        'RocketMQ Broker5分钟内QPS增长量过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(increase(rocketmq_broker_qps[5m]) * 100) / 100',
            'annotations', '实例：${instance}，QPS 的突然增长通常意味着突发的业务流量或异常的请求量。 如果增长过快且超出系统处理能力，可能会导致请求积压，响应时间增加，甚至服务不可用。 持续的高 QPS 可能导致资源耗尽，如 CPU、内存和 I/O 负载增高。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'RocketMQ',
        'RocketMQ Broker5分钟内TPS增长量过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(increase(rocketmq_broker_tps[5m]) * 100) / 100',
            'annotations', '实例：${instance}，TPS 的快速增长通常意味着大量的消息发送或消费请求。这可能导致 Broker 的压力增加。 当 TPS 增长过快时，Broker 可能面临资源瓶颈（例如内存、CPU、磁盘 I/O），导致处理效率下降。 突增的 TPS 可能由于异常业务流量或系统故障，需进一步排查原因。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'RocketMQ',
        'RocketMQ磁盘使用率过高',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(rocketmq_brokeruntime_commitlog_disk_ratio) by (brokerIP) * 100',
            'annotations', '实例：${instance}，磁盘空间不足可能导致消息无法写入，导致消息积压，甚至服务中断。 当磁盘使用率过高时，文件系统性能也可能下降，影响读写效率。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Jaeger',
        'Trace链路存在异常状态码',
        'Jaeger',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Traces',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 10, 'tags', '%7B\"http.status_code\"%3A\"5.%2A%3F\"%7D'),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        '服务日志监控',
        '接口耗时大于300ms',
        'Loki',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '{app="apisix"} !~ `socket.io` | json | upstream_response_time > 0.3 | upstream_response_time != ""',
            'logScope', 10,
            'evalCondition', JSON_OBJECT('type', 'count', 'operator', '>=', 'queryValue', 0, 'value', 1)
        ),
        10,
        0,
        'Logs',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        '服务日志监控',
        '服务ERROR错误日志「占比超过1%」',
        'Loki',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', 'sum(count_over_time({namespace="prod"} |~ `ERROR`[10m])) / sum(count_over_time({namespace="prod"}[10m])) * 100 > 1',
            'logScope', 10,
            'evalCondition', JSON_OBJECT('type', 'count', 'operator', '>=', 'queryValue', 0, 'value', 1)
        ),
        10,
        0,
        'Logs',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        '服务日志监控',
        '服务ERROR错误日志「增长率超过2%」',
        'Loki',
        JSON_OBJECT(
            'promQL', '',
            'annotations', '',
            'forDuration', 0,
            'rules', NULL
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', 'sum(rate({namespace="prod"} |~ `level` |~ `ERROR`[5m])) by (app) > 2',
            'logScope', 10,
            'evalCondition', JSON_OBJECT('type', 'count', 'operator', '>=', 'queryValue', 0, 'value', 1)
        ),
        10,
        0,
        'Logs',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器CPU使用率监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(100 - (avg(irate(node_cpu_seconds_total{mode="idle",}[5m])) by (instance) * 100))',
            'annotations', '节点：${instance}，CPU使用率过高，当前：${value}%，高 CPU 使用率意味着服务器可能接近处理能力上限，影响性能，导致应用程序响应变慢或服务中断！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80'),
                JSON_OBJECT('severity', 'P1', 'expr', '>75'),
                JSON_OBJECT('severity', 'P2', 'expr', '>70')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
        '',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );INSERT ignore INTO watchalert.rule_templates (
    rule_group_name,
    rule_name,
    datasource_type,
    prometheus_config,
    ali_cloud_sls_config,
    loki_config,
    eval_interval,
    for_duration,
    type,
    repeat_notice_interval,
    description,
    effective_time,
    jaeger_config,
    kubernetes_config,
    elastic_search_config
) VALUES
    (
        'Node节点监控',
        '服务器TCP连接数监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'sum(label_replace(max(node_sockstat_TCP_alloc{}) by (instance),"host_ip","$1","instance","(.*):.*"))by(instance)',
            'annotations', '节点：${instance}，TCP连接数过高，当前：${value}%，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>500')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器入网络吞吐量异常监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum by (instance) (rate(node_network_transmit_bytes_total[2m])) / 1024 / 1024 > 100) * on(instance) group_left (nodename) node_uname_info{nodename=~\".+\"}\n',
            'annotations', '节点：${instance}，主机网络接口可能接收了太多数据（> 100 MB/s），当前：${value}%，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>100')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器内存使用率监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(100 - (node_memory_MemAvailable_bytes{} / node_memory_MemTotal_bytes{} * 100))',
            'annotations', '节点：${instance}，内存使用率过高，当前：${value}%，内存紧张可能导致系统频繁进行内存交换（swap），进而引起性能下降。如果内存耗尽，应用程序可能崩溃或无法分配更多资源。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>80'),
                JSON_OBJECT('severity', 'P1', 'expr', '>75'),
                JSON_OBJECT('severity', 'P2', 'expr', '>70')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器出网络吞吐量异常监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum by (instance) (rate(node_network_receive_bytes_total[2m])) / 1024 / 1024 > 100) * on(instance) group_left (nodename) node_uname_info{nodename=~\".+\"}',
            'annotations', '节点：${instance}，主机网络接口可能发送了太多数据（> 100 MB/s），当前：${value}%，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>100')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器文件系统设备错误',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'node_filesystem_device_error',
            'annotations', '节点：${instance}，${mountpoint} 文件系统出现设备错误，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==1')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器时钟偏差',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '((node_timex_offset_seconds > 0.05 and deriv(node_timex_offset_seconds[5m]) >= 0) or (node_timex_offset_seconds < -0.05 and deriv(node_timex_offset_seconds[5m]) <= 0)) * on(instance) group_left (nodename) node_uname_info{nodename=~\".+\"}\n',
            'annotations', '节点：${instance}，检测到时钟偏差，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>=0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器状态异常',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'up{}',
            'annotations', '实例：${instance}，节点状态可能是exporter运行状态异常，也能是服务器宕机，当此告警发生后，需要进一步排查来确保节点的状态。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '==0')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器磁盘IO写流量过大',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum by (instance) (rate(node_disk_written_bytes_total[2m])) / 1024 / 1024 > 100) * on(instance) group_left (nodename) node_uname_info{nodename=~\".+\"}\n',
            'annotations', '节点：${instance}，磁盘可能写入了太多数据 (> 100 MB/s)，当前：${value}%，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>100')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器磁盘IO读流量过大',
        'Prometheus',
        JSON_OBJECT(
            'promQL', '(sum by (instance) (rate(node_disk_read_bytes_total[2m])) / 1024 / 1024 > 100) * on(instance) group_left (nodename) node_uname_info{nodename=~\".+\"}',
            'annotations', '节点：${instance}，磁盘可能读取了太多数据 (> 100 MB/s)，当前：${value}%，请尽快处理！',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>100')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    ),
    (
        'Node节点监控',
        '服务器磁盘使用率监控',
        'Prometheus',
        JSON_OBJECT(
            'promQL', 'round(max((node_filesystem_size_bytes{fstype=~\"ext.?|xfs\",}-node_filesystem_free_bytes{fstype=~\"ext.?|xfs\",}) *100/(node_filesystem_avail_bytes {fstype=~\"ext.?|xfs\",}+(node_filesystem_size_bytes{fstype=~\"ext.?|xfs\",}-node_filesystem_free_bytes{fstype=~\"ext.?|xfs\",})))by(ecs_cname,instance))',
            'annotations', '节点：${instance}，磁盘使用率过高，当前：${value}%，磁盘空间不足会导致文件无法写入、新日志无法记录，甚至可能使服务无法正常运行。尤其在数据增长迅速的系统中，这需要重点监控。',
            'forDuration', 60,
            'rules', JSON_ARRAY(
                JSON_OBJECT('severity', 'P0', 'expr', '>85'),
                JSON_OBJECT('severity', 'P1', 'expr', '>80'),
                JSON_OBJECT('severity', 'P2', 'expr', '>75')
            )
        ),
        JSON_OBJECT(
            'project', '',
            'logstore', '',
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        JSON_OBJECT(
            'logQL', '',
            'logScope', 0,
            'evalCondition', JSON_OBJECT('type', '', 'operator', '', 'queryValue', 0, 'value', 0)
        ),
        10,
        0,
        'Metrics',
        0,
		'',
        JSON_OBJECT('week', NULL, 'startTime', 0, 'endTime', 0),
        JSON_OBJECT('service', '', 'scope', 0, 'tags', ''),
        JSON_OBJECT('resource', '', 'reason', '', 'value', 0, 'filter', NULL, 'scope', 0),
        JSON_OBJECT('index', '', 'scope', 0, 'filter', NULL)
    );
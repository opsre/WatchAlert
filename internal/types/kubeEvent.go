package types

var EventResourceTypeList = []string{
	"Pods",
	"Nodes",
	"PVC/PV",
	"HPA",
}

type reason struct {
	Type   string `json:"type"`
	TypeCN string `json:"typeCN"`
}

var EventReasonLMapping = map[string][]reason{
	"Pods": {
		{
			"Failed",
			"容器启动失败(Failed)",
		},
		{
			"Unhealthy",
			"容器健康状况不佳(Unhealthy)",
		},
		{
			"Killing",
			"容器被终止(Killing)",
		},
		{
			"BackOff",
			"容器崩溃(BackOff)",
		},
		{
			"FailedMount",
			"挂载卷失败(FailedMount)",
		},
		{
			"FailedAttachVolume",
			"附加卷到节点失败(FailedAttachVolume)",
		},
		{
			"DeadlineExceeded",
			"Pod超过其运行期限(DeadlineExceeded)",
		},
		{
			"FailedScheduling",
			"Pod调度失败(FailedScheduling)",
		},
	},
	"Nodes": {
		{
			"NodeNotReady",
			"节点处于不可用状态(NodeNotReady)",
		},
		{
			"NodeUnderMemoryPressure",
			"节点处于内存压力下(NodeUnderMemoryPressure)",
		},
		{
			"NodeUnderDiskPressure",
			"节点处于磁盘压力下(NodeUnderDiskPressure)",
		},
	},
	"PVC/PV": {
		{
			"FailedBinding",
			"PVC/PV绑定失败(FailedBinding)",
		},
	},
	"HPA": {
		{
			"FailedRescale",
			"调整副本数失败(FailedRescale)",
		},
		{
			"FailedGetResourceMetric",
			"获取资源指标失败(FailedGetResourceMetric)",
		},
		{
			"FailedGetExternalMetric",
			"获取外部指标失败(FailedGetExternalMetric)",
		},
	},
}

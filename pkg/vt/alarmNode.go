package vt

import "watchAlert/internal/models"

const (
	Firing  string = "Firing"
	Recover        = "Recover"
)

// AlarmTreeNode 告警事件节点
type AlarmTreeNode struct {
	Value    string
	Alerts   map[string]models.AlertCurEvent
	Children map[string]*AlarmTreeNode
}

func NewTreeNode(value string) *AlarmTreeNode {
	return &AlarmTreeNode{
		Value:    value,
		Alerts:   make(map[string]models.AlertCurEvent),
		Children: make(map[string]*AlarmTreeNode),
	}
}

func (an *AlarmTreeNode) Set(TName string, Alerts map[string]models.AlertCurEvent) error {
	if len(Alerts) == 0 {
		return nil
	}

	exitsAlerts := an.Gets(TName)
	for fingerprint, alert := range Alerts {
		exitsAlerts[fingerprint] = alert
	}

	an.Children[TName] = &AlarmTreeNode{
		Alerts: exitsAlerts,
	}
	return nil
}

func (an *AlarmTreeNode) Gets(TName string) map[string]models.AlertCurEvent {
	if child, ok := an.Children[TName]; ok {
		return child.Alerts
	}

	return map[string]models.AlertCurEvent{}
}

func (an *AlarmTreeNode) List() *AlarmTreeNode {
	return an
}

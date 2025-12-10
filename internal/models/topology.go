package models

import (
	"fmt"
)

const (
	TopologyPrefix = "topology"
)

type Topology struct {
	TenantId  string `json:"tenantId"`
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Nodes     []Node `json:"nodes" gorm:"serializer:json"`
	Edges     []Edge `json:"edges" gorm:"serializer:json"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedAt int64  `json:"updatedAt"`
}

// TopologyList 用于列表查询，不包含nodes和edges字段以减少数据传输量
type TopologyList struct {
	TenantId  string `json:"tenantId"`
	ID        string `json:"id"`
	Name      string `json:"name"`
	UpdatedBy string `json:"updatedBy"`
	UpdatedAt int64  `json:"updatedAt"`
}

func (t *Topology) TableName() string {
	return "w8t_topology"
}

type Node struct {
	ID               string   `json:"id"`
	Type             string   `json:"type"`
	Data             NodeData `json:"data" gorm:"serializer:json"`
	PositionAbsolute Position `json:"positionAbsolute" gorm:"serializer:json"`
	Position         Position `json:"position" gorm:"serializer:json"`
	Dragging         bool     `json:"dragging"`
	Draggable        bool     `json:"draggable"`
	Selected         bool     `json:"selected"`
	Style            Style    `json:"style" gorm:"serializer:json"`
	Width            int      `json:"width"`
	Height           int      `json:"height"`
}

type NodeData struct {
	Label            string `json:"label"`
	SubLabel         string `json:"subLabel"`
	Type             string `json:"type"`
	EnablePrometheus bool   `json:"enablePrometheus"`
	MetricsLabel     string `json:"metricsLabel"`
	PrometheusQuery  string `json:"prometheusQuery"`
	Operator         string `json:"operator"`
	Threshold        string `json:"threshold"`
	DatasourceId     string `json:"datasourceId"`
}

type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Edge struct {
	ID           string    `json:"id"`
	Source       string    `json:"source"`
	Target       string    `json:"target"`
	SourceHandle string    `json:"sourceHandle"`
	TargetHandle string    `json:"targetHandle"`
	Type         string    `json:"type"`
	MarkerEnd    MarkerEnd `json:"markerEnd" gorm:"serializer:json"`
	Style        Style     `json:"style" gorm:"serializer:json"`
	Label        string    `json:"label"`
}

type MarkerEnd struct {
	Type  string `json:"type"`
	Color string `json:"color"`
}
type Style struct {
	Width       int    `json:"width"`
	Height      int    `json:"height"`
	Stroke      string `json:"stroke"`
	StrokeWidth int    `json:"strokeWidth"`
}

type TopologyCacheKey string

func BuildTopologyCacheKey(tenantId, topologyId string) TopologyCacheKey {
	return TopologyCacheKey(fmt.Sprintf("w8t:%s:%s:%s.info", tenantId, TopologyPrefix, topologyId))
}

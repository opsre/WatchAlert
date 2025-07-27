package models

type Page struct {
	Total int64 `json:"total" form:"total"`
	Index int64 `json:"index" form:"index"`
	Size  int64 `json:"size" form:"size"`
}

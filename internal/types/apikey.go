package types

type RequestApiKeyCreate struct {
	UserId      string `json:"userId" form:"userId"`
	Name        string `json:"name" form:"name" binding:"required"`
	Description string `json:"description" form:"description"`
}

type RequestApiKeyUpdate struct {
	ID          int    `json:"id" form:"id" binding:"required"`
	UserId      string `json:"userId" form:"userId"`
	Name        string `json:"name" form:"name"`
	Description string `json:"description" form:"description"`
}

type RequestApiKeyQuery struct {
	ID     int    `json:"id" form:"id"`
	UserId string `json:"userId" form:"userId"`
	Name   string `json:"name" form:"name"`
	Query  string `json:"query" form:"query"`
}

type ResponseApiKeyInfo struct {
	ID          int    `json:"id"`
	UserId      string `json:"userId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Key         string `json:"key"`
	CreatedAt   int64  `json:"createdAt"`
}

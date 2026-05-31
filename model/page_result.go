package model

type PageResult struct {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
type PageQuery struct {
	PageNum  int    `json:"pageNum"`
	PageSize int    `json:"pageSize"`
	Type     string `json:"type"`
}

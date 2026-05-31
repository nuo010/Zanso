package model

import (
	"time"
)

type TblFq struct {
	ID        string    `json:"id" gorm:"primarykey"`
	CreatTime time.Time `json:"creatTime"`
	KeyName   string    `json:"keyName"`
	Url       string    `json:"url"`
}

// 设置 User 模型对应的表名为 "tbl_unpw"
func (TblFq) TableName() string {
	return "tbl_fq"
}

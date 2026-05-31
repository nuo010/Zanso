package model

import (
	"time"
)

type TblEmail struct {
	ID        string    `json:"id" gorm:"primarykey"`
	CreatTime time.Time `json:"creatTime"`
	EmailName string    `json:"emailName"`
	SendEmail string    `json:"sendEmail"`
}

// 设置 User 模型对应的表名为 "tbl_unpw"
func (TblEmail) TableName() string {
	return "tbl_email"
}

package model

import (
	"time"
)

type TblMemo struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	U         string    `json:"u"`
	P         string    `json:"p"`
	BZ        string    `json:"bz"`
	TokenKey  string    `json:"tokenKey"`
}

// 设置 User 模型对应的表名为 "tbl_unpw"
func (TblMemo) TableName() string {
	return "tbl_memo"
}

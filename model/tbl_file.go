package model

import "time"

type TblFile struct {
	ID        string    `json:"id" gorm:"primarykey"`
	CreatTime time.Time `json:"creatTime" gorm:"creat_time"`
	Fsize     int64     `json:"fsize" gorm:"fsize"`
	Key       string    `json:"key" gorm:"key"`
	Type      string    `json:"type" gorm:"type"`
}

func (TblFile) TableName() string {
	return "tbl_file"
}

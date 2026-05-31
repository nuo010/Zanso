package model

import "time"

type TblTask struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	Text       string    `json:"text"`
	Complete   string    `json:"complete"`
	UpdateTime time.Time `json:"updateTime"`
	CreatTime  time.Time `json:"creatTime"`
	Del        string    `gorm:"-"`
}

func (TblTask) TableName() string {
	return "tbl_task"
}

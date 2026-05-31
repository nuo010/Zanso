package model

import "time"

type TblWeightRecord struct {
	ID        string    `json:"id" gorm:"primarykey"`
	Num       float64   `json:"num"`
	UserId    string    `json:"userId"`
	CreatTime time.Time `json:"creatTime"`
	FilePath  string    `json:"filePath"`
	NickName  string    `json:"nickName"`
	Avatar    string    `json:"avatar"`
}

func (TblWeightRecord) TableName() string {
	return "tbl_weight_record"
}

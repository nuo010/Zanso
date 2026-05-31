package model

import "time"

type TblQuotation struct {
	ID        string    `json:"id" gorm:"primarykey"`
	Quotation string    `json:"quotation"`
	CreatTime time.Time `json:"creatTime"`
}

func (TblQuotation) TableName() string {
	return "tbl_quotation"
}

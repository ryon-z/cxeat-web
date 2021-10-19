package models

// SubsItemMst : 구독 상품 관리
type SubsItemMst struct {
	SubsItemNo int    `json:"SubsItemNo" db:"SubsItemNo"`
	ItemName   string `json:"ItemName" db:"ItemName"`
	DpName     string `json:"DpName" db:"DpName"`
	UnitType   string `json:"UnitType" db:"UnitType"`
	UnitAmt    string `json:"UnitAmt" db:"UnitAmt"`
	StatusCode string `json:"StatusCode" db:"StatusCode"`
	RegDate    string `json:"RegDate" db:"RegDate"`
}

// TableName : 구독 상품 관리 테이블명
func (SubsItemMst) TableName() string {
	return "SubsItemMst"
}

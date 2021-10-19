package models

// FaqMst : FAQ관리
type FaqMst struct {
	FaqNo      int    `json:"FaqNo" db:"FaqNo"`
	FaqType    string `json:"FaqType" db:"FaqType"`
	FaqTitle   string `json:"FaqTitle" db:"FaqTitle"`
	FaqDesc    string `json:"FaqDesc" db:"FaqDesc"`
	StatusCode string `json:"StatusCode" db:"StatusCode"`
	RegDate    string `json:"RegDate" db:"RegDate"`
}

// TableName : FAQ관리 테이블명
func (FaqMst) TableName() string {
	return "FaqMst"
}

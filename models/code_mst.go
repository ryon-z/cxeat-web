package models

// CodeMst : 코드
type CodeMst struct {
	CodeNo     int     `json:"CodeNo" db:"CodeNo"`
	CodeType   string  `json:"CodeType" db:"CodeType"`
	CodeKey    string  `json:"CodeKey" db:"CodeKey"`
	CodeValue  *string `json:"CodeValue" db:"CodeValue"`
	CodeValue2 *string `json:"CodeValue2" db:"CodeValue2"`
	CodeLabel  string  `json:"CodeLabel" db:"CodeLabel"`
	CodeOrder  int     `json:"CodeOrder" db:"CodeOrder"`
	IsUse      int     `json:"IsUse" db:"IsUse"`
	RegDate    string  `json:"RegDate" db:"RegDate"`
	RegUser    string  `json:"RegUser" db:"RegUser"`
	UpdDate    *string `json:"UpdDate" db:"UpdDate"`
	UpdUser    *string `json:"UpdUser" db:"UpdUser"`
}

// TableName : 코드 테이블명
func (CodeMst) TableName() string {
	return "CodeMst"
}

package models

// CodeType : 코드 타입
type CodeType struct {
	CodeTypeNo   int    `json:"CodeTypeNo" db:"CodeTypeNo"`
	CodeType     string `json:"CodeType" db:"CodeType"`
	CodeTypeDesc string `json:"CodeTypeDesc" db:"CodeTypeDesc"`
	IsUse        int    `json:"IsUse" db:"IsUse"`
	RegDate      string `json:"RegDate" db:"RegDate"`
	RegUser      string `json:"RegUser" db:"RegUser"`
}

// TableName : 코드 타입 테이블명
func (CodeType) TableName() string {
	return "CodeType"
}

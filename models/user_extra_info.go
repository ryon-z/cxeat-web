package models

// UserExtraInfo : 회원 부가 정보
type UserExtraInfo struct {
	UserNo       int    `json:"UserNo" db:"UserNo"`
	SubsNo       int    `json:"SubsNo" db:"SubsNo"`
	InfoType     string `json:"InfoType" db:"InfoType"`
	InfoTypeDesc string `json:"InfoTypeDesc" db:"InfoTypeDesc"`
	InfoData     string `json:"InfoData" db:"InfoData"`
	ExtraDesc    string `json:"ExtraDesc" db:"ExtraDesc"`
	RegDate      string `json:"RegDate" db:"RegDate"`
}

// TableName : 회원 부가 정보 테이블명
func (UserExtraInfo) TableName() string {
	return "UserExtraInfo"
}

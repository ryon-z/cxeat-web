package models

// AgreementMst : 약관관리
type AgreementMst struct {
	AgreementNo    int    `json:"AgreementNo" db:"AgreementNo"`
	AgreementType  string `json:"AgreementType" db:"AgreementType"`
	AgreementTitle string `json:"AgreementTitle" db:"AgreementTitle"`
	AgreementDesc  string `json:"AgreementDesc" db:"AgreementDesc"`
	AttachFileURL  string `json:"AttachFileURL" db:"AttachFileURL"`
	StatusCode     string `json:"StatusCode" db:"StatusCode"`
	RegDate        string `json:"RegDate" db:"RegDate"`
}

// TableName : 약관관리 테이블명
func (AgreementMst) TableName() string {
	return "AgreementMst"
}

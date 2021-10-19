package models

// SubsPlan : 구독 금액
type SubsPlan struct {
	SubsPlanCode  string `json:"SubsPlanCode" db:"SubsPlanCode"`
	SubsPlanName  string `json:"SubsPlanName" db:"SubsPlanName"`
	SubsPlanDesc  string `json:"SubsPlanDesc" db:"SubsPlanDesc"`
	SubsPlanPrice int    `json:"SubsPlanPrice" db:"SubsPlanPrice"`
	StatusCode    string `json:"StatusCode" db:"StatusCode"`
	RegDate       string `json:"RegDate" db:"RegDate"`
}

// TableName : 구독 금액 테이블명
func (SubsPlan) TableName() string {
	return "SubsPlan"
}

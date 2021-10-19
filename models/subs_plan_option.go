package models

// SubsPlanOption : 구독 구상 옵션
type SubsPlanOption struct {
	SubsPlanCode   string `json:"SubsPlanCode" db:"SubsPlanCode"`
	OptionCode     string `json:"OptionCode" db:"OptionCode"`
	OptionName     string `json:"OptionName" db:"OptionName"`
	OptionAddPrice int    `json:"OptionAddPrice" db:"OptionAddPrice"`
	StatusCode     string `json:"StatusCode" db:"StatusCode"`
	RegDate        int    `json:"RegDate" db:"RegDate"`
}

// TableName : 구독 구상 옵션 테이블명
func (SubsPlanOption) TableName() string {
	return "SubsPlanOption"
}

package models

// UserAddress : 회원 주소록
type UserAddress struct {
	AddressNo    int     `json:"AddressNo" db:"AddressNo"`
	UserNo       int     `json:"UserNo" db:"UserNo"`
	AddressLabel *string `json:"AddressLabel" db:"AddressLabel"`
	RcvName      *string `json:"RcvName" db:"RcvName"`
	RoadAddress  *string `json:"RoadAddress" db:"RoadAddress"`
	LotAddress   *string `json:"LotAddress" db:"LotAddress"`
	SubAddress   *string `json:"SubAddress" db:"SubAddress"`
	PostNo       *string `json:"PostNo" db:"PostNo"`
	ContactNo    *string `json:"ContactNo" db:"ContactNo"`
	ReqMsg       *string `json:"ReqMsg" db:"ReqMsg"`
	IsBasic      int     `json:"IsBasic" db:"IsBasic"`
	StatusCode   string  `json:"StatusCode" db:"StatusCode"`
	RegDate      string  `json:"RegDate" db:"RegDate"`
}

// TableName : 회원 주소록 테이블명
func (UserAddress) TableName() string {
	return "UserAddress"
}

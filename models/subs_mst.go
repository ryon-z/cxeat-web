package models

// SubsMst : 회원 구독
type SubsMst struct {
	SubsNo      int     `json:"SubsNo" db:"SubsNo"`
	UserNo      int     `json:"UserNo" db:"UserNo"`
	CardRegNo   int     `json:"CardRegNo" db:"CardRegNo"`
	TagGroupNo  *int    `json:"TagGroupNo" db:"TagGroupNo"`
	SubsType    string  `json:"SubsType" db:"SubsType"`
	CateType    string  `json:"CateType" db:"CateType"`
	BoxType     string  `json:"BoxType" db:"BoxType"`
	SubsPrice   int     `json:"SubsPrice" db:"SubsPrice"`
	PeriodDay   int     `json:"PeriodDay" db:"PeriodDay"`
	FirstDate   string  `json:"FirstDate" db:"FirstDate"`
	NextDate    *string `json:"NextDate" db:"NextDate"`
	RcvName     *string `json:"RcvName" db:"RcvName"`
	MainAddress *string `json:"MainAddress" db:"MainAddress"`
	SubAddress  *string `json:"SubAddress" db:"SubAddress"`
	ContactNo   *string `json:"ContactNo" db:"ContactNo"`
	PostNo      *string `json:"PostNo" db:"PostNo"`
	ReqMsg      *string `json:"ReqMsg" db:"ReqMsg"`
	AddnlDesc   *string `json:"AddnlDesc" db:"AddnlDesc"`
	StatusCode  string  `json:"StatusCode" db:"StatusCode"`
	RegDate     string  `json:"RegDate" db:"RegDate"`
	UpdDate     *string `json:"UpdDate" db:"UpdDate"`
	CnlDate     *string `json:"CnlDate" db:"CnlDate"`
	CnlReason   *string `json:"CnlReason" db:"CnlReason"`
	SubsReason  *string `json:"SubsReason" db:"SubsReason"`
	RefCode     *string `json:"RefCode" db:"RefCode"`
}

// TableName : 회원 구독 테이블명
func (SubsMst) TableName() string {
	return "SubsMst"
}

package models

// ReviewMst : 리뷰 정보
type ReviewMst struct {
	ReviewNo   int     `json:"ReviewNo" db:"ReviewNo"`
	OrderNo    int     `json:"OrderNo" db:"OrderNo"`
	ReviewDesc *string `json:"ReviewDesc" db:"ReviewDesc"`
	RegDate    string  `json:"RegDate" db:"RegDate"`
	UpdDate    *string `json:"UpdDate" db:"UpdDate"`
}

// TableName : 리뷰 정보 테이블명
func (ReviewMst) TableName() string {
	return "ReviewMst"
}

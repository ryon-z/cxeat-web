package models

// ItemReview : 아이템 리뷰 정보
type ItemReview struct {
	ItemReviewNo int     `json:"ItemReviewNo" db:"ItemReviewNo"`
	ReviewNo     int     `json:"ReviewNo" db:"ReviewNo"`
	OrderNo      int     `json:"OrderNo" db:"OrderNo"`
	ItemNo       int     `json:"ItemNo" db:"ItemNo"`
	ReviewScore  float64 `json:"ReviewScore" db:"ReviewScore"`
	ReviewDesc   *string `json:"ReviewDesc" db:"ReviewDesc"`
	RegDate      string  `json:"RegDate" db:"RegDate"`
	UpdDate      *string `json:"UpdDate" db:"UpdDate"`
}

// TableName : 아이템 리뷰 정보 테이블명
func (ItemReview) TableName() string {
	return "ItemReview"
}

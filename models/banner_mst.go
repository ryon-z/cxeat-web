package models

// BannerMst : 배너관리
type BannerMst struct {
	BannerNo     int    `json:"BannerNo" db:"BannerNo"`
	BannerType   string `json:"BannerType" db:"BannerType"`
	BannerCode   string `json:"BannerCode" db:"BannerCode"`
	BannerTitle  string `json:"BannerTitle" db:"BannerTitle"`
	BannerDesc   string `json:"BannerDesc" db:"BannerDesc;"`
	BannerImgURL string `json:"BannerImgURL" db:"BannerImgURL;"`
	StatusCode   string `json:"StatusCode" db:"StatusCode"`
	RegDate      string `json:"RegDate" db:"RegDate"`
}

// TableName : 배너관리 테이블명
func (BannerMst) TableName() string {
	return "BannerMst"
}

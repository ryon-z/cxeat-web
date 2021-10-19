package models

// SubsBundleMst : 구독 상품 구성 관리
type SubsBundleMst struct {
	BundleNo     int    `json:"BundleNo" db:"BundleNo"`
	SubsPlanType string `json:"SubsPlanType" db:"SubsPlanType"`
	BundleName   string `json:"BundleName" db:"BundleName"`
	StatusCode   string `json:"StatusCode" db:"StatusCode"`
	RegDate      string `json:"RegDate" db:"RegDate"`
}

// TableName : 구독 상품 구성 관리 테이블명
func (SubsBundleMst) TableName() string {
	return "SubsBundleMst"
}

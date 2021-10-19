package models

// SubsBundleItem : 구독 구성 상품
type SubsBundleItem struct {
	BundleNo   int `json:"BundleNo" db:"BundleNo"`
	SubsItemNo int `json:"SubsItemNo" db:"SubsItemNo"`
	ItemCnt    int `json:"ItemCnt" db:"ItemCnt"`
}

// TableName : 구독 구성 상품 테이블명
func (SubsBundleItem) TableName() string {
	return "SubsBundleItem"
}

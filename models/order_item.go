package models

// OrderItem : 주문 상품
type OrderItem struct {
	OrderNo int    `json:"OrderNo" db:"OrderNo"`
	ItemNo  string `json:"ItemNo" db:"ItemNo"`
	ItemCnt int    `json:"ItemCnt" db:"ItemCnt"`
	RegDate string `json:"RegDate" db:"RegDate"`
}

// TableName : 주문 상품 테이블명
func (OrderItem) TableName() string {
	return "OrderItem"
}

package models

// OrderPayment : 주문 결제 정보
type OrderPayment struct {
	OrderNo      int    `json:"OrderNo" db:"OrderNo"`
	PaymentType  string `json:"PaymentType" db:"PaymentType"`
	HostData     string `json:"HostData" db:"HostData"`
	TID          string `json:"TID" db:"TID"`
	PaymentPrice int    `json:"PaymentPrice" db:"PaymentPrice"`
	StatusCode   string `json:"StatusCode" db:"StatusCode"`
	RegDate      string `json:"RegDate" db:"RegDate"`
}

// TableName : 주문 결제 정보 테이블명
func (OrderPayment) TableName() string {
	return "OrderPayment"
}

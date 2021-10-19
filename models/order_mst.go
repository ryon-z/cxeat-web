package models

// OrderMst : 주문
type OrderMst struct {
	OrderNo        int     `json:"OrderNo" db:"OrderNo"`
	UserNo         int     `json:"UserNo" db:"UserNo"`
	CardRegNo      int     `json:"CardRegNo" db:"CardRegNo"`
	TagGroupNo     *int    `json:"TagGroupNo" db:"TagGroupNo"`
	SubsNo         *int    `json:"SubsNo" db:"SubsNo"`
	OrderType      string  `json:"OrderType" db:"OrderType"`
	CateType       string  `json:"CateType" db:"CateType"`
	BoxType        string  `json:"BoxType" db:"BoxType"`
	OrderRound     *int    `json:"OrderRound" db:"OrderRound"`
	OrderPrice     int     `json:"OrderPrice" db:"OrderPrice"`
	DiscountPrice  int     `json:"DiscountPrice" db:"DiscountPrice"`
	RcvName        *string `json:"RcvName" db:"RcvName"`
	MainAddress    *string `json:"MainAddress" db:"MainAddress"`
	SubAddress     *string `json:"SubAddress" db:"SubAddress"`
	ContactNo      *string `json:"ContactNo" db:"ContactNo"`
	PostNo         *string `json:"PostNo" db:"PostNo"`
	ReqDelivDate   string  `json:"ReqDelivDate" db:"ReqDelivDate"`
	ReqMsg         *string `json:"ReqMsg" db:"ReqMsg"`
	DelivCo        *string `json:"DelivCo" db:"DelivCo"`
	DelivInvoiceNo *string `json:"DelivInvoiceNo" db:"DelivInvoiceNo"`
	AddnlDesc      *string `json:"AddnlDesc" db:"AddnlDesc"`
	StatusCode     string  `json:"StatusCode" db:"StatusCode"`
	RegDate        string  `json:"RegDate" db:"RegDate"`
	UpdDate        *string `json:"UpdDate" db:"UpdDate"`
	CnlDate        *string `json:"CnlDate" db:"CnlDate"`
}

// TableName : 주문 테이블명
func (OrderMst) TableName() string {
	return "OrderMst"
}

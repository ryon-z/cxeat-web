package models

// OrderHist : FAQ관리
type OrderHist struct {
	OrderNo    int     `json:"OrderNo" db:"OrderNo"`
	StatusCode *string `json:"StatusCode" db:"StatusCode"`
	HistDesc   *string `json:"HistDesc" db:"HistDesc"`
	ExecUser   string  `json:"ExecUser" db:"ExecUser"`
	ExecDate   string  `json:"ExecDate" db:"ExecDate"`
}

// TableName : FAQ관리 테이블명
func (OrderHist) TableName() string {
	return "OrderHist"
}

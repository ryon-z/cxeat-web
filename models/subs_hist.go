package models

// SubsHist : FAQ관리
type SubsHist struct {
	SubsNo     int     `json:"SubsNo" db:"SubsNo"`
	StatusCode *string `json:"StatusCode" db:"StatusCode"`
	HistDesc   *string `json:"HistDesc" db:"HistDesc"`
	ExecUser   string  `json:"ExecUser" db:"ExecUser"`
	ExecDate   string  `json:"ExecDate" db:"ExecDate"`
}

// TableName : FAQ관리 테이블명
func (SubsHist) TableName() string {
	return "SubsHist"
}

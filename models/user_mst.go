package models

// UserMst : 회원
type UserMst struct {
	UserNo        int     `json:"UserNo" db:"UserNo"`
	UserType      string  `json:"UserType" db:"UserType"`
	UserID        string  `json:"UserID" db:"UserID"`
	UserSecretKey string  `json:"UserSecretKey" db:"UserSecretKey"`
	UserName      string  `json:"UserName" db:"UserName"`
	UserEmail     *string `json:"UserEmail" db:"UserEmail"`
	UserPhone     string  `json:"UserPhone" db:"UserPhone"`
	UserGender    *string `json:"UserGender" db:"UserGender"`
	BirthDay      *string `json:"BirthDay" db:"BirthDay"`
	IsMktAgree    int     `json:"IsMktAgree" db:"IsMktAgree"`
	Funnel        *string `json:"Funnel" db:"Funnel"`
	StatusCode    string  `json:"StatusCode" db:"StatusCode"`
	RegDate       string  `json:"RegDate" db:"RegDate"`
	LeaveDate     *string `json:"LeaveDate" db:"LeaveDate"`
	LastLoginDate *string `json:"LastLoginDate" db:"LastLoginDate"`
}

// TableName : 회원 테이블명
func (UserMst) TableName() string {
	return "UserMst"
}

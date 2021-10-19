package models

// UserCard : 회원 카드 등록 정보
type UserCard struct {
	CardRegNo    int     `json:"CardRegNo" db:"CardRegNo"`
	UserNo       int     `json:"UserNo" db:"UserNo"`
	CardNickName *string `json:"CardNickName" db:"CardNickName"`
	CardCode     *string `json:"CardCode" db:"CardCode"`
	CardName     *string `json:"CardName" db:"CardName"`
	CardNumber   *string `json:"CardNumber" db:"CardNumber"`
	CardKey      string  `json:"CardKey" db:"CardKey"`
	IsBasic      int     `json:"IsBasic" db:"IsBasic"`
	StatusCode   string  `json:"StatusCode" db:"StatusCode"`
	RegDate      string  `json:"RegDate" db:"RegDate"`
}

// TableName : 회원 카드 등록 정보 테이블명
func (UserCard) TableName() string {
	return "UserCard"
}

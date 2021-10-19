package models

// Tag : 태그 그룹
type Tag struct {
	TagNo      int     `json:"TagNo" db:"TagNo"`
	TagGroupNo int     `json:"TagGroupNo" db:"TagGroupNo"`
	TagType    string  `json:"TagType" db:"TagType"`
	TagLabel   string  `json:"TagLabel" db:"TagLabel"`
	TagValue   string  `json:"TagValue" db:"TagValue"`
	IsUse      int     `json:"IsUse" db:"IsUse"`
	RegDate    string  `json:"RegDate" db:"RegDate"`
	RegUser    string  `json:"RegUser" db:"RegUser"`
	UpdDate    *string `json:"UpdDate" db:"UpdDate"`
	UpdUser    *string `json:"UpdUser" db:"UpdUser"`
}

// TableName : 태그 그룹 테이블명
func (Tag) TableName() string {
	return "Tag"
}

package models

// TagGroup : 태그 그룹
type TagGroup struct {
	TagGroupNo   int    `json:"TagGroupNo" db:"TagGroupNo"`
	UserNo       int    `json:"UserNo" db:"UserNo"`
	TagGroupType string `json:"TagGroupType" db:"TagGroupType"`
	RegDate      string `json:"RegDate" db:"RegDate"`
}

// TableName : 태그 그룹 테이블명
func (TagGroup) TableName() string {
	return "TagGroup"
}

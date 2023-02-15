package model

type School struct {
	ID    int
	Name  string `gorm:"column:wxname"`
	Token string `gorm:"column:token"`
	Hurl  string `gorm:"column:hurl"`
}

func (School) TableName() string {
	return "school"
}

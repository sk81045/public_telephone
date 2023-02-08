package model

type Device struct {
	ID         int
	Sid        int     `gorm:"column:sid"`
	Tag        string  `gorm:"column:tag"`
	Key        string  `gorm:"column:key"`
	Status     bool    `gorm:"column:status"`
	Fee        float32 `gorm:"column:fee"`
	Category   string  `gorm:"column:category"`
	Created_at string  `gorm:"column:created_at"`
}

func (Device) TableName() string {
	return "device"
}

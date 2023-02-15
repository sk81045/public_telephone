package model

type Payorder struct {
	ID         int     `gorm:"column:id,primaryKey"`
	Sid        int     `gorm:"column:sid"`
	Pid        int     `gorm:"column:pid"`
	Lid        int     `gorm:"column:lid"`
	Studentid  string  `gorm:"column:student_id"`
	Ic         string  `gorm:"column:ic"`
	Orderid    string  `gorm:"column:orderid"`
	Price      float32 `gorm:"column:price"`
	Type       int     `gorm:"column:type"`
	From       string  `gorm:"column:from"`
	Paystatus  bool    `gorm:"column:paystatus"`
	Category   string  `gorm:"column:category"`
	Sync       bool    `gorm:"column:sync"`
	Created_at int64   `gorm:"column:created_at"`
}

func (Payorder) TableName() string {
	return "payorder"
}

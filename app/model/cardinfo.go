package model

type CardinfoResult struct {
	Code   int        `json:"code"`
	Result []Cardinfo `json:"result"`
}

type Cardinfo struct {
	Name       string `gorm:"column:Name"`
	UserNO     string `gorm:"column:UserNO"`
	AfterPay   string `gorm:"column:AfterPay"`
	MacID      string `gorm:"column:MacID"`
	MacType    string `gorm:"column:MacType"`
	Cardid     string `gorm:"column:Cardid"`
	MerchantID string `gorm:"column:MerchantID"`
	CardState  string `gorm:"column:CardState"`
}

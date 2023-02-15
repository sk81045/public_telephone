package model

type CardinfoRes struct {
	Code   int        `json:"code"`
	Result []Cardinfo `json:"result"`
}

type Cardinfo struct {
	FID          int    `gorm:"column:FID"`
	UserNO       string `gorm:"column:userNO"`
	MacID        string `gorm:"column:macID"`
	MacType      string `gorm:"column:macType"`
	PayMoney     string `gorm:"column:payMoney"`
	AfterPay     string `gorm:"column:afterPay"`
	PayTime      string `gorm:"column:payTime"`
	CardPayCount string `gorm:"column:cardPayCount"`
	PayTimeFrame string `gorm:"column:payTimeFrame"`
	PayKind      string `gorm:"column:payKind"`
	AddMode      string `gorm:"column:addMode"`
	UpdateState  string `gorm:"column:updateState"`
	TimeBucket   string `gorm:"column:timeBucket"`
	Sky          string `gorm:"column:sky"`
	PayListNO    string `gorm:"column:PayListNO"`
}

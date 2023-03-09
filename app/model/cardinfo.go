package model

type CardinfoResult struct {
	Code   int        `json:"code"`
	Result []Cardinfo `json:"result"`
}

type Cardinfo struct {
	FID        string `gorm:"column:FID"`
	UserNO     string `gorm:"column:UserNO"`
	AfterPay   string `gorm:"column:AfterPay"`
	MacID      string `gorm:"column:MacID"`
	MacType    string `gorm:"column:MacType"`
	Cardid     string `gorm:"column:Cardid"`
	MerchantID string `gorm:"column:MerchantID"`
}

// type Cardinfo struct {
// 	FID          int    `gorm:"column:FID"`
// 	UserNO       string `gorm:"column:userNO"`
// 	Cardid       string `gorm:"column:Cardid"`
// 	MacID        string `gorm:"column:macID"`
// 	MacType      string `gorm:"column:macType"`
// 	PayMoney     string `gorm:"column:payMoney"`
// 	AfterPay     string `gorm:"column:afterPay"`
// 	PayTime      string `gorm:"column:payTime"`
// 	CardPayCount string `gorm:"column:cardPayCount"`
// 	PayTimeFrame string `gorm:"column:payTimeFrame"`
// 	PayKind      string `gorm:"column:payKind"`
// 	AddMode      string `gorm:"column:addMode"`
// 	UpdateState  string `gorm:"column:updateState"`
// 	TimeBucket   string `gorm:"column:timeBucket"`
// 	Sky          string `gorm:"column:sky"`
// 	PayListNO    string `gorm:"column:PayListNO"`
// }

package main

import (
	// "fmt"
	// "Hwgen/app"
	"Hwgen/app/service"
	"Hwgen/core"
	"Hwgen/global"
)

func main() {
	global.H_VIPER = core.Viper() // 初始化Viper
	global.H_DB = core.Gorm()     // gorm连接数据库
	db, _ := global.H_DB.DB()
	defer db.Close()
	service.Run()
}

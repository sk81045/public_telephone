package controller

import (
	model "Hwgen/app/model"
	"Hwgen/global"
	helpers "Hwgen/utils"
	"fmt"
)

var (
	piece_1     string
	piece_2     string
	piece_3     string
	piece_4     string
	piece_5     string
	piece_6     string
	instruction string
)

type Phonecode struct {
	Id   int `json:"id"`
	Code int `json:"key"`
}

// @Description  获取亲情号码
// @param_1 初始报文
func Operation_01(origin string) (string, error) {
	piece_1 = origin[0:4]
	piece_2 = origin[4:6]
	piece_3 = origin[6:10]
	piece_4 = origin[10:18] //话机code

	device, err := GetDevice(piece_4)
	if err != nil {
		return "错误->话机已停用", err
	} else {
		fmt.Println("话机状态正常", device.Key)
	}

	piece_5 = helpers.Hex2Dec(origin[18:26]) //IC
	student, err := GetStudent(piece_5)
	if err != nil {
		fmt.Println("no record", err)
		return "错误->查无此人", err
	}

	if student.Balance < 5 {
		err := fmt.Errorf("err:余额不足")
		return "错误->余额不足", err
	}
	if len(student.Parents) == 0 {
		err := fmt.Errorf("err:没有查询到绑定的号码")
		return "错误->没有绑定家长", err
	}
	var Phones string
	var Relation string
	for _, v := range student.Parents {
		if len(v.Phone) == 11 {
			Phones += helpers.JoiningString(v.Phone, " ", 4) //拼接字符
		}

		switch v.Guanxi {
		case "01":
			Relation += "爸爸"

		case "02":
			Relation += "妈妈"

		case "03":
			Relation += "爷爷"

		case "04":
			Relation += "奶奶"

		case "05":
			Relation += "姥姥"

		case "06":
			Relation += "姥爷"

		case "07":
			Relation += "亲戚"

		default:
			Relation += "家长"
		}
	}
	piece_6 = origin[26:len(origin)] //time
	nums := fmt.Sprintf("%d", len(student.Parents))
	fmt.Println("last", piece_1+piece_2+piece_3+"10"+nums+Phones+Relation+"0000"+origin[26:len(origin)])
	//instruction = piece1 + piece2 + piece3 + "10315033414330    13785182848    18032453047    爸爸妈妈爷爷000020221105222159"
	instruction = piece_1 + piece_2 + piece_3 + "10" + nums + Phones + Relation + "0000" + origin[26:len(origin)]
	return instruction, err
}

// 获取设备信息
func GetDevice(key string) (model.Device, error) {
	Db := global.H_DB
	var device = model.Device{}
	err := Db.Model(&model.Device{}).Where("`key` = ? AND `category` = ? AND `status` = ?", key, 1, 1).Find(&device).Error
	return device, err
}

// 获取学生信息
func GetStudent(ic string) (model.Students, error) {
	Db := global.H_DB
	var student = model.Students{}
	err := Db.Model(&model.Students{}).Where("`cardid` = ?", ic).Preload("Parents").Find(&student).Error
	return student, err
}

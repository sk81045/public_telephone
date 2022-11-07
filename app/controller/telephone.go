package controller

import (
	model "Hwgen/app/model"
	"Hwgen/global"
	helpers "Hwgen/utils"
	"fmt"
)

var (
	instruction string
)

type Origin struct {
	piece_1 string
	piece_2 string
	piece_3 string
	piece_4 string
	piece_5 string
	piece_6 string
}

// @Description  公话认证
// @param_1 初始报文
func (o *Origin) Operation_10(origin string) (string, error) {
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	o.piece_4 = origin[10:18] //phone_code
	instruction := o.piece_1 + o.piece_2
	device, err := GetDevice(o.piece_4)
	if err != nil {
		instruction += o.piece_3 + "0"
		fmt.Println(err)
		return instruction, err
	} else {
		instruction += o.piece_3 + "1"
		fmt.Println("话机状态正常", device.Key)
		return instruction, err
	}
}

// @Description  获取亲情号码
// @param_1 初始报文
func (o *Origin) Operation_01(origin string) (string, error) {
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	o.piece_5 = helpers.Hex2Dec(origin[18:26]) //IC
	o.piece_6 = origin[26:len(origin)]         //time

	instruction := o.piece_1 + o.piece_2 + o.piece_3
	valid := "00"
	student, err := GetStudent(o.piece_5)
	if err != nil {
		fmt.Println(err)
		return instruction + valid, err
	}
	if student.Balance < 5 {
		err := fmt.Errorf("IC卡余额不足")
		return instruction + valid, err
	}
	if len(student.Parents) == 0 {
		err := fmt.Errorf("err:没有查询到绑定的号码")
		return instruction + valid, err
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
	fmt.Println("IC卡状态正常", student.Name)
	valid = "10"
	nums := fmt.Sprintf("%d", len(student.Parents))
	lastinstruction := instruction + valid + nums + Phones + Relation + "0000" + o.piece_6
	return lastinstruction, err
}

// 获取设备信息
func GetDevice(key string) (model.Device, error) {
	Db := global.H_DB
	var device = model.Device{}
	err := Db.Model(&model.Device{}).Where("`key` = ? AND `category` = ? AND `status` = ?", key, 1, 1).Find(&device).Error
	if device.ID == 0 {
		err = fmt.Errorf("func GetDevice():没有获取到此设备信息")
	}
	return device, err
}

// 获取学生信息
func GetStudent(ic string) (model.Students, error) {
	Db := global.H_DB
	var student = model.Students{}
	err := Db.Model(&model.Students{}).Where("`cardid` = ?", ic).Preload("Parents").Find(&student).Error
	if student.ID == 0 {
		err = fmt.Errorf("func GetStudent():没有获取到学生信息")
	}
	return student, err
}

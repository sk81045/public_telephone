package controller

import (
	model "Hwgen/app/model"
	"Hwgen/global"
	helpers "Hwgen/utils"
	"fmt"
	"strconv"
	"time"
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	o.piece_4 = origin[10:18] //phone_code
	head := helpers.PackageHead(o.piece_1 + o.piece_2 + o.piece_3 + "0")
	instruction = head + o.piece_2
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
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	o.piece_5 = helpers.Hex2Dec(origin[18:26]) //IC
	o.piece_6 = origin[26:len(origin)]         //time

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
			Relation += "教师"
		default:
			Relation += "家长"
		}
	}
	fmt.Println("IC卡状态正常", helpers.ConvertStr2GBK(Relation))
	valid = "10"
	nums := fmt.Sprintf("%d", len(student.Parents))
	Relation = helpers.ConvertStr2GBK(Relation)
	instruction = o.piece_1 + o.piece_2 + o.piece_3 + valid + nums + Phones + Relation + "0000" + o.piece_6
	head := helpers.PackageHead(instruction)
	lastinstruction := head + o.piece_2 + o.piece_3 + valid + nums + Phones + Relation + "0000" + o.piece_6
	return lastinstruction, err
}

// @Description  处理通话订单
// @param_1 初始报文
func (o *Origin) Operation_03(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	// basics := origin[0:4] + origin[4:6] + origin[6:10]
	KEY := origin[10:18]
	IC := helpers.Hex2Dec(origin[18:26]) //IC
	// ORDER := origin[44:46]               //ORDER
	STIME := origin[46:60] //STIME
	DURATION, _ := strconv.ParseFloat(origin[60:66], 32)
	NUMBER := origin[66:77] //NUMBER

	Student, _ := GetStudent(IC)

	Calling := CallingLog(model.Calllog{
		Pid:         Student.ID,
		Sid:         Student.Sid,
		Key:         KEY,
		Ic:          IC,
		Describe:    "Calling",
		PhoneNumber: NUMBER,
		CallTime:    int(DURATION),
		Cost:        float32(DURATION) / 60 * 0.2, //计费
		Stime:       STIME,
		Created_at:  time.Now().Unix(),
	})

	_ = CallingOrder(model.Payorder{
		Pid:      Student.ID,
		Sid:      Student.Sid,
		Lid:      Calling.ID,
		Orderid:  "tp" + fmt.Sprintf("%d", time.Now().UnixNano()) + helpers.RandStr(6),
		Price:    float32(DURATION) / 60 * 0.2, //计费
		From:     "telephone:" + KEY,
		Category: "1",
	})

	instruction = helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + "1")
	lastinstruction := instruction + origin[4:6] + origin[6:10] + "1"
	return lastinstruction, nil
}

// @Description  公话状态告警
// @param_1 初始报文
func (o *Origin) Operation_81(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	instruction = helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + "1")
	lastinstruction := instruction + origin[4:6] + origin[6:10] + "1"
	return lastinstruction, nil
}

// @Description  公话状态告警
// @param_1 初始报文
func (o *Origin) Operation_82(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	params := make(map[string]string)
	params["key"] = origin[10:18]
	params["version"] = origin[18:28]
	params["power"] = origin[28:29]
	params["receiver"] = origin[29:30]
	params["door"] = origin[30:31]

	fmt.Println("params", params)
	return instruction, nil
}

// @Description  获取公话状态
// @param_1 初始报文
func (o *Origin) TelephoneState() (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	t := time.Now().Unix()
	t2 := time.Unix(t, 0).Format("20060102150405")
	instruction = "0024820001" + t2
	return instruction, nil
}

//处理通话订单
func CallingOrder(data model.Payorder) model.Payorder {
	result := global.H_DB.Create(&data)
	if result.Error != nil {
		panic("func CallingOrder():处理通话订单失败")
	}
	return data
}

//处理通话记录
func CallingLog(data model.Calllog) model.Calllog {
	result := global.H_DB.Create(&data)
	if result.Error != nil {
		panic("func CallingLog():处理通话记录失败")
	}
	return data
}

// 获取设备信息
func GetDevice(key string) (model.Device, error) {
	var device = model.Device{}
	err := global.H_DB.Model(&model.Device{}).Where("`key` = ? AND `category` = ? AND `status` = ?", key, 1, 1).Find(&device).Error
	if device.ID == 0 {
		panic("func GetDevice():没有获取到此设备信息")
	}
	return device, err
}

// 获取学生信息
func GetStudent(ic string) (model.Students, error) {
	var student = model.Students{}
	err := global.H_DB.Model(&model.Students{}).Where("`cardid` = ?", ic).Preload("Parents").Find(&student).Error
	if student.ID == 0 {
		panic("func GetStudent():没有获取到学生信息")
	}
	return student, err
}

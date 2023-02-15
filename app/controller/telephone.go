package controller

import (
	model "Hwgen/app/model"
	"Hwgen/global"
	helpers "Hwgen/utils"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
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
			global.H_LOG.Warn("func Operation_03()", zap.String("公话认证:", err.(string)))
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
		fmt.Println("话机状态正常", "key:", device.Key, "费率:", device.Fee, "归属:", device.School.Name)
		return instruction, err
	}
}

// @Description  获取亲情号码
// @param_1 初始报文
func (o *Origin) Operation_01(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			global.H_LOG.Warn("func Operation_01()", zap.String("获取亲情号码:", err.(string)))
			fmt.Println(err)
		}
	}()
	o.piece_1 = origin[0:4]
	o.piece_2 = origin[4:6]
	o.piece_3 = origin[6:10]
	o.piece_5 = helpers.Hex2Dec(origin[18:26]) //IC
	o.piece_6 = origin[26:len(origin)]         //time

	valid := "00"
	student := GetStudent(o.piece_5)
	if student.ID == 0 {
		return instruction + valid, nil
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
	valid = "10"
	nums := fmt.Sprintf("%d", len(student.Parents))
	Relation = helpers.ConvertStr2GBK(Relation)
	instruction = o.piece_1 + o.piece_2 + o.piece_3 + valid + nums + Phones + Relation + "0000" + o.piece_6
	head := helpers.PackageHead(instruction)
	lastinstruction := head + o.piece_2 + o.piece_3 + valid + nums + Phones + Relation + "0000" + o.piece_6
	return lastinstruction, nil
}

// @Description  处理通话订单
// @param_1 初始报文
func (o *Origin) Operation_03(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			global.H_LOG.Warn("func Operation_03()", zap.String("处理通话订单:", err.(string)))
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

	Student := GetStudent(IC)
	device, _ := GetDevice(KEY)
	fee := float32(DURATION) / 60 * device.Fee //计费
	Calling := CallingLog(model.Calllog{
		Pid:         Student.ID,
		Sid:         Student.Sid,
		Key:         KEY,
		Ic:          IC,
		Describe:    "Calling",
		PhoneNumber: NUMBER,
		CallTime:    int(DURATION),
		Cost:        fee,
		Stime:       STIME,
		Created_at:  time.Now().Unix(),
	})

	_ = CallingOrder(model.Payorder{
		Pid:        Student.ID,
		Sid:        Student.Sid,
		Lid:        Calling.ID,
		Ic:         IC,
		Orderid:    "tp" + fmt.Sprintf("%d", time.Now().UnixNano()) + helpers.RandStr(6),
		Price:      fee,
		Type:       2,
		From:       "telephone:" + KEY,
		Category:   "1",
		Created_at: time.Now().Unix(),
	})

	instruction = helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + "1")
	lastinstruction := instruction + origin[4:6] + origin[6:10] + "1"
	return lastinstruction, nil
}

// @Description  计费通话话单结算
// @param_1 初始报文
func (o *Origin) Operation_13(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	device_id := origin[10:18]
	ic := helpers.Hex2Dec(origin[18:26])
	len_time, _ := strconv.ParseFloat(origin[60:66], 32) //通话时长
	start_time := origin[46 : 46+14]
	phone := origin[66 : 66+11]

	Student := GetStudent(ic)
	device, _ := GetDevice(device_id)
	fee := float32(len_time) / 60 * device.Fee //计费
	Calling := CallingLog(model.Calllog{
		Pid:         Student.ID,
		Sid:         Student.Sid,
		Key:         device_id,
		Ic:          ic,
		Describe:    "Calling",
		PhoneNumber: phone,
		CallTime:    int(len_time),
		Cost:        fee, //计费
		Stime:       start_time,
		Created_at:  time.Now().Unix(),
	})

	order := CallingOrder(model.Payorder{
		Pid:        Student.ID,
		Sid:        Student.Sid,
		Studentid:  Student.Studentid,
		Lid:        Calling.ID,
		Ic:         ic,
		Orderid:    "tp" + fmt.Sprintf("%d", time.Now().UnixNano()) + helpers.RandStr(6),
		Price:      fee, //计费
		Type:       2,
		From:       "telephone:" + device_id,
		Category:   "1",
		Created_at: time.Now().Unix(),
	})
	valid := "1"
	if order.ID == 0 {
		global.H_LOG.Info("func Operation_13()", zap.String("处理通话订单失败:", fmt.Sprintf("%d", Calling.ID)))
		valid = "0"
	}
	head := helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + valid)
	lastinstruction := head + origin[4:6] + origin[6:10] + valid
	global.H_LOG.Info("func Operation_13()", zap.String("处理通话订单(计费电话):", lastinstruction))
	return lastinstruction, nil
}

// @Description  获取通话费率
// @param_1 初始报文
func (o *Origin) Operation_75(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()

	valid := "1"

	device, err := GetDevice(origin[10:18])
	if err != nil {
		fmt.Println("没有获取设备信息")
		valid = "4"
	}

	f1 := fmt.Sprintf("%.0f", device.Fee*100)
	fee := helpers.JoiningString2(f1, "0", 4-len(f1)) //拼接字符得到费率

	// student := GetStudent(helpers.Hex2Dec(origin[18:26]))
	// if student.ID == 0 {
	// 	fmt.Println("没有获取学生信息")
	// 	valid = "0"
	// }

	ic, err := Geticcard(device.Sid, helpers.Hex2Dec(origin[18:26]))
	if err != nil {
		fmt.Println("无效的IC卡")
		valid = "0"
	}

	score, _ := strconv.ParseFloat(ic.AfterPay, 64)
	fmt.Println("IC卡号:", helpers.Hex2Dec(origin[18:26]))
	fmt.Println("人员编号:", ic.UserNO)
	fmt.Println("卡余额", score)

	f1 = fmt.Sprintf("%.0f", score*100)
	balance := helpers.JoiningString2(f1, "0", 10-len(f1)) //拼接字符得到余额
	// phone_code := origin[26:37]
	instruction := origin[4:6] + origin[6:10] + valid + balance + fee + "999999"
	head := helpers.PackageHead(origin[0:4] + instruction)
	lastinstruction := head + instruction
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
	global.H_LOG.Info("func Operation_81()", zap.String("记录话机状态:", instruction))
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
	err := global.H_DB.Model(&model.Device{}).Preload("School").Where("`key` = ? AND `category` = ? AND `status` = ?", key, 1, 1).Find(&device).Error
	if device.ID == 0 {
		panic("func GetDevice():没有获取到此设备信息")
	}
	return device, err
}

// 获取学生信息
func GetStudent(ic string) model.Students {
	var student = model.Students{}
	err := global.H_DB.Model(&model.Students{}).Where("`cardid` = ?", ic).Preload("Parents").Find(&student).Error
	// if student.ID == 0 {
	// 	panic("func GetStudent():没有获取到学生信息")
	// }
	fmt.Println(err)
	return student
}

// 获取学校信息
func GetSchool(id int) model.School {
	var sc = model.School{}
	global.H_DB.Model(&model.School{}).Where("`id` = ? ", id).Find(&sc)
	fmt.Println("学校:", sc.Name)
	return sc
}

// 获取最新的IC卡信息
func Geticcard(sid int, ic string) (model.Cardinfo, error) {
	sc := GetSchool(sid)
	body := helpers.HttpGet(sc.Hurl + "/work/cardinfo?ic=" + ic)
	var res model.CardinfoRes
	err := json.Unmarshal([]byte(body), &res)
	if nil != err {
		fmt.Println("Geticcard() json.Unmarshal err:", err)
		return model.Cardinfo{}, err
	}

	return res.Result[0], nil
}

package controller

import (
	model "Hwgen/app/model"
	"Hwgen/global"
	helpers "Hwgen/utils"
	"encoding/json"
	// "errors"
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
	o.piece_4 = origin[10:18]                  //key
	o.piece_5 = helpers.Hex2Dec(origin[18:26]) //IC
	o.piece_6 = origin[26:len(origin)]         //time

	device, err := GetDevice(o.piece_4) //利用设备信息获取所属单位
	if err != nil {
		fmt.Println(err)
	}

	// icDetails, err := GetIcDetails(device.Sid, o.piece_5) //从学校售饭取得最新的IC卡信息
	// if err != nil {
	// 	fmt.Println(err)
	// }

	valid := "10"
	student, err := GetStudentIc(o.piece_5, device.Sid) //从平台获取人员
	if err != nil {
		valid = "00"
		fmt.Println(err)
	}

	fmt.Println("IC卡号:", o.piece_5)
	// if student.Balance < 0.1 {
	// 	err := fmt.Errorf("IC卡余额不足")
	// 	valid = "00"
	// 	fmt.Println(err)
	// }
	if len(student.Parents) == 0 {
		err := fmt.Errorf("err:没有查询到绑定的号码")
		valid = "00"
		fmt.Println(err)
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
			global.H_LOG.Warn("func Operation_03()", zap.String("处理亲情通话订单失败:", err.(string)))
			fmt.Println(err)
		}
	}()

	KEY := origin[10:18]
	IC := helpers.Hex2Dec(origin[18:26]) //IC
	STIME := origin[46:60]               //STIME
	DURATION, _ := strconv.ParseFloat(origin[60:66], 32)
	NUMBER := origin[66:77] //NUMBER

	device, err := GetDevice(KEY) //利用设备信息获取所属单位
	if err != nil {
		fmt.Println(err)
	}
	// icDetails, err := GetIcDetails(device.Sid, IC) //从学校售饭取得最新的IC卡信息
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// Student, err := GetStudent(icDetails.UserNO, device.Sid)
	// if err != nil {
	// 	fmt.Println(err)
	// 	return instruction, err
	// }

	Student, err := GetStudentIc(IC, device.Sid) //从平台获取人员
	if err != nil {
		fmt.Println(err)
		return instruction, err
	}

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

	order := CallingOrder(model.Payorder{
		Pid:        Student.ID,
		Sid:        Student.Sid,
		Lid:        Calling.ID,
		Ic:         IC,
		Orderid:    "tp" + helpers.RandStr(3) + fmt.Sprintf("%d", time.Now().UnixNano()),
		Price:      fee,
		Type:       2,
		From:       "telephone:" + KEY,
		Category:   "1",
		Created_at: time.Now().Unix(),
	})

	if order.ID == 0 {
		return instruction, fmt.Errorf(fmt.Sprintf("%d", Calling.ID))
	}
	valid := "1"
	instruction = helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + valid)
	lastinstruction := instruction + origin[4:6] + origin[6:10] + valid
	return lastinstruction, nil
}

// @Description  计费通话话单结算
// @param_1 初始报文
func (o *Origin) Operation_13(origin string) (string, error) {
	defer func() {
		if err := recover(); err != nil {
			global.H_LOG.Warn("func Operation_13()", zap.String("处理计费通话订单失败:", err.(string)))
			fmt.Println(err)
		}
	}()

	key := origin[10:18]
	ic := helpers.Hex2Dec(origin[18:26])
	len_time, _ := strconv.ParseFloat(origin[60:66], 32) //通话时长
	start_time := origin[46 : 46+14]
	phone := origin[66 : 66+11]

	device, err := GetDevice(key) //利用设备信息获取所属单位
	if err != nil {
		fmt.Println(err)
	}
	icDetails, err := GetIcDetails(device.Sid, ic) //从学校售饭取得最新的IC卡信息
	if err != nil {
		fmt.Println(err)
	}
	Student, err := GetStudent(icDetails.UserNO, device.Sid)
	if err != nil {
		fmt.Println(err)
		return instruction, err
	}

	fee := float32(len_time) / 60 * device.Fee //计费
	Calling := CallingLog(model.Calllog{
		Pid:         Student.ID,
		Sid:         Student.Sid,
		Key:         key,
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
		Orderid:    "tp" + helpers.RandStr(4) + fmt.Sprintf("%d", time.Now().UnixNano()),
		Price:      fee, //计费
		Type:       2,
		From:       "telephone:" + key,
		Category:   "1",
		Created_at: time.Now().Unix(),
	})

	if order.ID == 0 {
		return instruction, fmt.Errorf(fmt.Sprintf("%d", Calling.ID))
	}
	fmt.Printf("订单id:%d 处理成功\n", order.ID)
	valid := "1"
	head := helpers.PackageHead(origin[0:4] + origin[4:6] + origin[6:10] + valid)
	lastinstruction := head + origin[4:6] + origin[6:10] + valid
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

	phone_code := origin[26:37]
	if phone_code == "110" {
		fmt.Println("禁止拨打此号码", phone_code)
		valid = "8"
	}

	device, err := GetDevice(origin[10:18])
	if err != nil {
		fmt.Println("错误 没有获取此设备信息")
		valid = "4"
	}
	fee := helpers.JoiningString2(fmt.Sprintf("%.0f", device.Fee*100), "0", 4-len(fmt.Sprintf("%.0f", device.Fee*100))) //拼接字符得到费率

	icDetails, err := GetIcDetails(device.Sid, helpers.Hex2Dec(origin[18:26]))
	if err != nil {
		fmt.Println(err)
		valid = "0"
	}
	score, _ := strconv.ParseFloat(icDetails.AfterPay, 64) //余额

	_, err = GetStudent(icDetails.UserNO, device.Sid) //平台验证IC卡
	if err != nil {
		fmt.Printf("人员编号 %s 平台没有录入\n", icDetails.UserNO)
		valid = "0"
	}

	fmt.Println("IC卡号:", icDetails.Cardid)
	fmt.Println("人员编号:", icDetails.UserNO)
	fmt.Println("卡余额", score)

	balance := helpers.JoiningString2(fmt.Sprintf("%.0f", score*100), "0", 10-len(fmt.Sprintf("%.0f", device.Fee*100))) //拼接字符得到余额

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
		return device, fmt.Errorf("没有获取到此设备信息")
	}
	return device, err
}

// 获取学生信息
func GetStudent(stuid string, sid int) (model.Students, error) {
	var student = model.Students{}
	global.H_DB.Model(&model.Students{}).Where("`studentid` = ? AND `sid` = ?", stuid, sid).Preload("Parents").Find(&student)
	if student.ID == 0 {
		return student, fmt.Errorf("没有获取到人员信息-根据学号")
	}
	return student, nil
}

// 获取学生信息
func GetStudentIc(ic string, sid int) (model.Students, error) {
	var student = model.Students{}
	global.H_DB.Model(&model.Students{}).Where("`cardid` = ? AND `sid` = ?", ic, sid).Preload("Parents").Find(&student)
	if student.ID == 0 {
		return student, fmt.Errorf("没有获取到人员信息-根据IC卡号")
	}
	return student, nil
}

// 获取学校信息
func GetSchool(id int) model.School {
	var sc = model.School{}
	global.H_DB.Model(&model.School{}).Where("`id` = ? ", id).Find(&sc)
	fmt.Println("学校:", sc.Name)
	return sc
}

// 获取最新的IC卡信息-旧版
func GetIcDetails(sid int, ic string) (model.Cardinfo, error) {
	sc := GetSchool(sid)
	body := helpers.HttpGet(sc.Hurl + "/work/cardinfo?ic=" + ic)
	fmt.Println("Hurl:", sc.Hurl+"/work/cardinfo?ic="+ic)
	var res model.CardinfoRes
	err := json.Unmarshal([]byte(body), &res)
	if nil != err {
		panic("FRP接口错误-请求失败")
	}

	if len(res.Result) == 0 {
		panic("FRP接口错误-没有获取到IC卡信息")
	}

	return res.Result[0], nil
}

// 获取最新的IC卡信息-新版
func GetIcDetailsNew(sid int, stuid string) (model.CardinfoNew, error) {
	sc := GetSchool(sid)
	body := helpers.HttpGet(sc.Hurl + "/work/record?stime=2023-02-09 00:00&etime=2023-02-25 00:00&factor=" + stuid)
	fmt.Println("Hurl:", sc.Hurl+"/work/record?stime=2023-02-09 00:00&etime=2023-02-25 00:00&factor="+stuid)
	var res model.CardinfoResNew
	err := json.Unmarshal([]byte(body), &res)
	if nil != err {
		panic("FRP接口错误-请求失败")
	}

	if len(res.Result) == 0 {
		panic("FRP接口错误-没有获取到IC卡信息")
	}

	return res.Result[0], nil
}

package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

// @Description  16进制转10进制
// @param_1 16进制字符
func Hex2Dec(val string) string {
	n, err := strconv.ParseUint(val, 16, 32)
	int64Str := strconv.FormatUint(n, 10)
	if err != nil {
		fmt.Println(err)
	}
	return int64Str
}

// @Description  拼接字符串达到规定长度
func JoiningString(s1 string, s2 string, le int) string {
	var p string
	for i := 0; i < le; i++ {
		p += s2
	}
	var build strings.Builder
	build.WriteString(s1)
	build.WriteString(p)
	return build.String()
}

// @Description  生成随机数
const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func RandStr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

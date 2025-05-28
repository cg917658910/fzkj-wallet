package THB

import (
	"strconv"
	"strings"
)

// 工具函数：泰文月份转数字
func thaiMonthToNum(month string) (res string) {
	var m int = 0
	arr := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	for i, v := range arr {
		if month == v {
			m = i
			break
		}
	}
	if m == 0 {
		res = ""
		return
	}
	prefix := ""
	if m < 10 {
		prefix = "0"
	}
	res = prefix + strconv.Itoa(m)

	return
}

// 工具函数：英文月份转数字
func engMonthToNum(month string) (res string) {
	var m int = 0
	arr := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sept", "Oct", "Nov", "Dec"}
	for i, v := range arr {
		if month == v {
			m = i + 1
			break
		}
	}
	if m == 0 {
		res = ""
		return
	}
	prefix := ""
	if m < 10 {
		prefix = "0"
	}
	res = prefix + strconv.Itoa(m)
	return
}

func amountStrToFloat64(aStr string) (a float64) {
	aStr = strings.ReplaceAll(aStr, ",", "")
	aStr = strings.ReplaceAll(aStr, "บ", "")
	aStr = strings.ReplaceAll(aStr, "-", "")
	aStr = strings.ReplaceAll(aStr, "+", "")
	a, _ = strconv.ParseFloat(aStr, 64)

	return
}

// 工具
func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

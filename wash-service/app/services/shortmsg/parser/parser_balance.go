package parser

import (
	"regexp"
	"strconv"
	"strings"
)

type Extracted struct {
	PayTime int64
	Balance float64
}

// 匹配余额 ใช้ได้36,447.43บ 或 ใช้ได้ 36,447.43บ
func MatchBalanceUsed(s string) (bal float64, b bool) {
	s = strings.ReplaceAll(s, ",", "")
	re := regexp.MustCompile(`ใช้ได้\s*(\d+(\.\d{1,5})?)\s*บ`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		bal, _ := strconv.ParseFloat(m[1], 64)
		return bal, true
	}
	return
}

// 匹配 เหลือ 7,420.57 或 คงเหลือ 8,236.90 บ.
func MatchBalanceRemain(s string) (bal float64, b bool) {
	s = strings.ReplaceAll(s, ",", "")
	re := regexp.MustCompile(`(?:เหลือ|คงเหลือ)\s*(\d+(\.\d{1,5})?)(\s*บ)?`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		bal, _ := strconv.ParseFloat(m[1], 64)
		return bal, true
	}
	return
}

// 匹配 ชใช้ได้29,760.21บ or ชใช้ได้ 29,760.21บ
func MatchBalanceUsed2(s string) (bal float64, b bool) {
	s = strings.ReplaceAll(s, ",", "")
	re := regexp.MustCompile(`ชใช้ได้\s*(\d+(\.\d{1,5})?)\s*บ`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		bal, _ := strconv.ParseFloat(m[1], 64)
		return bal, true
	}
	return
}

// 匹配  all
func MatchBalanceAll(s string) (bal float64, b bool) {
	s = strings.ReplaceAll(s, ",", "")
	re := regexp.MustCompile(`(?:เหลือ|คงเหลือ|ชใช้ได้|ใช้ได้)\s*(\d+(\.\d{1,5})?)(\s*บ)?`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 2 {
		bal, _ := strconv.ParseFloat(m[1], 64)
		return bal, true
	}
	return
}

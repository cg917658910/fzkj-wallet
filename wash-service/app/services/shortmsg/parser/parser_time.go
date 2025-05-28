package parser

import (
	"regexp"
	"strconv"
	"time"
)

// 1. 匹配 18/05@14:44 或 27-05@16:25
func MatchDayMonthAtTime(s string) (*time.Time, bool) {
	// 支持 / 或 - 分隔
	re := regexp.MustCompile(`(\d{1,2})[/-](\d{1,2})@(\d{1,2}):(\d{2})`)
	m := re.FindStringSubmatch(s)
	if len(m) == 5 {
		year := time.Now().Year()
		month, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[1])
		hour, _ := strconv.Atoi(m[3])
		min, _ := strconv.Atoi(m[4])
		t := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local)
		return &t, true
	}
	return nil, false
}

// 2. 匹配 (17/5/68,16:27)  //泰历年份处理
func MatchThaiDateParen(s string) (*time.Time, bool) {
	re := regexp.MustCompile(`\((\d{1,2})/(\d{1,2})/(\d{2,4}),(\d{1,2}):(\d{2})\)`)
	m := re.FindStringSubmatch(s)
	if len(m) == 6 {
		day, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		year, _ := strconv.Atoi(m[3])
		if year < 100 { // 68 -> 2568
			year += 2500
		}
		// 转公历
		year -= 543
		hour, _ := strconv.Atoi(m[4])
		min, _ := strconv.Atoi(m[5])
		t := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local)
		return &t, true
	}
	return nil, false
}

// 3. 匹配 เมื่อ 17 พ.ค. 2568 - 16:28
func MatchThaiTextDate(s string) (*time.Time, bool) {
	re := regexp.MustCompile(`เมื่อ (\d{1,2}) ([^ ]+) (\d{4}) - (\d{1,2}):(\d{2})`)
	m := re.FindStringSubmatch(s)
	if len(m) == 6 {
		day, _ := strconv.Atoi(m[1])
		month := thaiMonthToNum(m[2])
		year, _ := strconv.Atoi(m[3])
		year -= 543
		hour, _ := strconv.Atoi(m[4])
		min, _ := strconv.Atoi(m[5])
		t := time.Date(year, time.Month(month), day, hour, min, 0, 0, time.Local)
		return &t, true
	}
	return nil, false
}

// 4. 匹配 20/05/68 17:11
func MatchDayMonthYearTime(s string) (*time.Time, bool) {
	re := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{2,4}) (\d{1,2}):(\d{1,2})(:(\d{1,2}))?`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 6 {
		var sec int
		day, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		year, _ := strconv.Atoi(m[3])
		if year < 100 {
			year += 2500
		}
		year -= 543
		hour, _ := strconv.Atoi(m[4])
		min, _ := strconv.Atoi(m[5])
		if len(m) >= 8 {
			sec, _ = strconv.Atoi(m[7])
		}
		t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
		return &t, true
	}
	return nil, false
}

// 5. 匹配 20/05/25@17:11
func MatchDayMonthYearTimeUse2(s string) (*time.Time, bool) {
	re := regexp.MustCompile(`(\d{1,2})/(\d{1,2})/(\d{2,4})@(\d{1,2}):(\d{1,2})(:(\d{1,2}))?`)
	m := re.FindStringSubmatch(s)
	if len(m) >= 6 {
		var sec int
		day, _ := strconv.Atoi(m[1])
		month, _ := strconv.Atoi(m[2])
		year, _ := strconv.Atoi(m[3])
		if year < 100 {
			year += 2000
		}
		hour, _ := strconv.Atoi(m[4])
		min, _ := strconv.Atoi(m[5])
		if len(m) >= 8 {
			sec, _ = strconv.Atoi(m[7])
		}
		t := time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)
		return &t, true
	}
	return nil, false
}

// 泰文月份转数字
func thaiMonthToNum(month string) int {
	arr := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	for i, v := range arr {
		if month == v {
			return i
		}
	}
	return 0
}

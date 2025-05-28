package THB

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/parser"
	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/types"
)

// SCB
func ExtractSCB(ramk string) (res *types.Extracted) {

	res = &types.Extracted{
		Code: codes.Success,
	}
	// clean [,] 方便计算金额
	ramk = strings.ReplaceAll(ramk, ",", "")
	remarkArr := strings.Fields(ramk)
	var payCoin float64
	if len(remarkArr) < 2 {
		res.Code = codes.ErrUnsupportedMessageFormat
		return
	}
	if strings.Contains(ramk, "True") || strings.Contains(ramk, "โอนเงิน") || strings.Contains(ramk, "ถอน") {
		res.Code = codes.ErrUnsupportedMessageFormat.New("当前数据不是入款数据。")
		return
	}
	if remarkArr[0] == "เงิน" {
		//เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ
		payCoin, _ = strconv.ParseFloat(strings.TrimSuffix(remarkArr[1], "บ"), 64)
	} else if remarkArr[0] == "Transfer" && len(remarkArr) > 5 {
		payCoin, _ = strconv.ParseFloat(remarkArr[5], 64)
	} else if remarkArr[0] == "Cash/transfer" && len(remarkArr) > 4 {
		payCoin, _ = strconv.ParseFloat(remarkArr[4], 64)
	} else {
		//18/05@14:44 29.92 จากKTB/x129439เข้าx997010 ใช้ได้10,121.23บ
		payCoin, _ = strconv.ParseFloat(remarkArr[1], 64)
	}
	// 时间
	t, ok := parser.MatchDayMonthAtTime(ramk)
	if ok {
		res.PayTime = t.Unix()
	}
	// 余额
	b, ok := parser.MatchBalanceAll(ramk)
	if ok {
		res.Balance = b
	}
	res.PayCoin = payCoin
	return
}

// SCB读取
func ExtractSCBRead(ramk string) *types.Extracted {
	return ExtractSCB(ramk)
}

// SCB通知
func ExtractSCBNotify(ramk string) (res *types.Extracted) {
	//คุณได้รับเงิน 69.89 บาท ผ่านรายการพร้อมเพย์จาก KTB / xxxxxx6048 เข้าบัญชี xxxxxx1083 เมื่อ 17 พ.ค. 2568 - 16:28
	res = &types.Extracted{}
	remarkArr := strings.Fields(ramk)

	if len(remarkArr) == 15 {
		res.PayCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[1], ",", ""), 64)
		payTime, ok := parser.MatchThaiTextDate(ramk)
		if ok {
			res.PayTime = payTime.Unix()
		}
	}
	return
}

// SCB流水
func ExtractSCBWater(msgTime string) (res *types.Extracted) {
	res.Code = codes.ErrUnimplemented
	// 无处理
	return nil
}

// TM流水
func ExtractTMWater(ramk string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	// 无数据，暂不处理
	/* var obj struct {
		Time  string `json:"time"`
		Money string `json:"money"`
	}
	_ = json.Unmarshal([]byte(ramk), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	timeArr := strings.Fields(obj.Time)
	if len(timeArr) < 4 {
		return
	}
	// 泰历转公历
	year, _ := strconv.Atoi(timeArr[2])
	year = year - 543
	month := thaiMonthToNum(timeArr[1])
	payTimeStr := strconv.Itoa(year) + "-" + month + "-" + timeArr[0] + " " + timeArr[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local) */
	return
}

// KTB
func ExtractKTB(ramk string) (res *types.Extracted) {
	//17-05@16:27 บชX58589X:เงินเข้า 50.34บ ใช้ได้ 11,911.89บ
	// TODO:
	res = &types.Extracted{}
	remark := strings.Fields(ramk)
	if strings.Contains(ramk, "Deposit") {
		res.Code = codes.ErrUnimplemented.New("Deposit unimplemented")
		return
	}
	var payCoin float64
	if len(remark) > 2 {
		payCoin = amountStrToFloat64(remark[2])
	}
	payTime, ok := parser.MatchDayMonthAtTime(ramk)
	if ok {
		res.PayTime = payTime.Unix()
	}
	balances, ok := parser.MatchBalanceAll(ramk)
	if ok {
		res.Balance = balances
	}
	res.PayCoin = payCoin
	return
}

// KTBLine
func ExtractKTBLine(ramk string) (res *types.Extracted) {
	//เงินออก: -11,000.00 บาท จากบัญชี XX5568 เมื่อ 27/0
	//เงินเข้า: 1,000.20 บาท เข้าบัญชี XX5568 เมื่อ 27/0
	res = &types.Extracted{}
	remark := strings.Fields(ramk)
	var payCoin float64
	if len(remark) > 1 {
		payCoin = amountStrToFloat64(remark[1])
	}
	res.PayCoin = payCoin
	return
}

// KTB通知
func ExtractKTBNotice(ramk string) (res *types.Extracted) {
	//ได้รับ +59.84 บาท เข้าพร้อมเพย์ จากบัญชี ไทยพาณิชย์ XXX-X-XX805-8
	res = &types.Extracted{}
	remark := strings.Fields(ramk)
	if strings.Contains(ramk, "Deposit") {
		res.Code = codes.ErrUnimplemented.New("Deposit unimplemented")
		return
	}
	var payCoin float64
	if len(remark) >= 2 {
		payCoin = amountStrToFloat64(remark[1])
	}
	res.PayCoin = payCoin
	return
}

// KTB流水
func ExtractKTBWater(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	var obj struct {
		MsgTime string  `json:"msg_time"`
		Coin    float64 `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 2 {
		return
	}
	dateParts := strings.Split(parts[0], "-")
	if len(dateParts) != 3 {
		return
	}
	payTimeStr := dateParts[2] + "-" + dateParts[1] + "-" + dateParts[0] + " " + parts[1]
	time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	return
}

// KBANK通知
func ExtractKBANKNotify(ramk string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* remarkArr := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	//arr := []string{"", "ม.ค.", "ก.พ.", "มี.ค.", "เม.ย.", "พ.ค.", "มิ.ย.", "ก.ค.", "ส.ค.", "ก.ย.", "ต.ค.", "พ.ย.", "ธ.ค."}
	if len(remarkArr) == 14 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[4], ",", ""), 64)
		year, _ := strconv.Atoi(remarkArr[10])
		year = 2500 + year - 543
		month := thaiMonthToNum(remarkArr[9])
		day := remarkArr[8]
		payTimeStr := strconv.Itoa(year) + "-" + month + "-" + day + " " + remarkArr[12]
		t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
		payTime = t.Unix()
	} */
	return
}

// KBANK读取
func ExtractKBANKRead(ramk string) (res *types.Extracted) {
	//17/05/68 16:27 บช X-9426 เงินเข้า 999.63 คงเหลือ 8,236.90 บ.
	res = &types.Extracted{}
	remark := strings.ReplaceAll(ramk, "KBank:", "")
	remarkArr := strings.Fields(remark)
	var payCoin float64
	var balance float64
	if len(remarkArr) == 9 || len(remarkArr) == 14 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[5], ",", ""), 64)
	} else if len(remarkArr) == 10 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[6], ",", ""), 64)
		if !isNumeric(remarkArr[6]) {
			payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[5], ",", ""), 64)
		}
	} else if len(remarkArr) == 4 || len(remarkArr) == 5 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[3], ",", ""), 64)
	} else if len(remarkArr) == 6 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remarkArr[4], ",", ""), 64)
	}
	// 余额
	if len(remarkArr) > 2 {
		balanceStr := remarkArr[len(remarkArr)-2]
		balance, _ = strconv.ParseFloat(strings.ReplaceAll(balanceStr, ",", ""), 64)
	}
	payTime, ok := parser.MatchDayMonthYearTime(ramk)
	if ok {
		res.PayTime = payTime.Unix()
	}
	res.PayCoin = payCoin
	res.Balance = balance
	return
}

// KBANK流水
func ExtractKBANKWater(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* var obj struct {
		MsgTime string      `json:"msg_time"`
		Coin    interface{} `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 4 {
		return
	}
	year, _ := strconv.Atoi(parts[2])
	year = 2500 + year - 543
	month := thaiMonthToNum(parts[1])
	day := parts[0]
	payTimeStr := strconv.Itoa(year) + "-" + month + "-" + day + " " + parts[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local)
	coinStr := ""
	switch v := obj.Coin.(type) {
	case string:
		coinStr = v
	case float64:
		coinStr = strconv.FormatFloat(v, 'f', 2, 64)
	}
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(coinStr, ",", ""), 64) */
	return
}

// BBL流水
func ExtractBBLWater(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* var obj struct {
		MsgTime string  `json:"msg_time"`
		Coin    float64 `json:"coin"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	parts := strings.Fields(obj.MsgTime)
	if len(parts) < 4 {
		return
	}
	month := engMonthToNum(parts[1])
	payTimeStr := parts[2] + "-" + month + "-" + parts[0] + " " + parts[3]
	t, _ := time.ParseInLocation("2006-1-2 15:04", payTimeStr, time.Local) */
	return
}

// BBL
func ExtractBBL(remark string) (res *types.Extracted) {
	//เงินเข้าบ/ชX0272 จากพร้อมเพย์ผ่านMB 12,000.00บ เงินในบ/ชใช้ได้46,769.66บ
	res = &types.Extracted{}
	if strings.Contains(remark, "ถอน/โอน") {
		res.Code = codes.ErrUnimplemented.New("暂未实现ถอน/โอน")
		return
	}
	remarkArr := strings.Fields(remark)
	if len(remarkArr) == 0 {
		res.Code = codes.ErrInvalidArgument
		return
	}
	var payCoin float64

	if remarkArr[0] == "PromptPay" && len(remarkArr) > 8 {
		payCoin = amountStrToFloat64(remarkArr[8])
	} else if len(remarkArr) > 2 {
		payCoin = amountStrToFloat64(remarkArr[2])
	}
	balance, ok := parser.MatchBalanceUsed2(remark)
	if ok {
		res.Balance = balance
	}
	res.PayCoin = payCoin
	return
}

// BAAC
func ExtractBAAC(ramk string) (res *types.Extracted) {
	//20/05/68 17:52:48 บัญชี x9576 รับโอนพร้อมเพย์ จำนวน 10.00 บาท ยอดเงินคงเหลือ 10.39 บาท
	res = &types.Extracted{}
	// รับโอนพร้อมเพย์ 付款金额不处理
	if strings.Contains(ramk, "รับโอนพร้อมเพย์") {
		res.Code = codes.ErrUnsupportedMessageFormat.New("付款金额不处理")
		return
	}
	// TODO: 检查 数据合法格式
	remark := strings.Fields(ramk)
	// coin
	if len(remark) > 6 {
		res.PayCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[6], ",", ""), 64)
	}
	// 时间
	t, ok := parser.MatchDayMonthYearTime(ramk)
	if ok {
		res.PayTime = t.Unix()
	}
	// 余额
	b, ok := parser.MatchBalanceAll(ramk)
	if ok {
		res.Balance = b
	}
	return
}

// TTB
func ExtractTTB(remark string) (res *types.Extracted) {
	//โอนเงิน13,000.00บ.ไปยังเบอร์มือถือXX3781 เหลือ680.31บ.17/05/25@16:34
	//มีเงิน50.14บ.โอนเข้าบ/ชxx5440เหลือ17,081.71บ.24/05/25@11:42
	res = &types.Extracted{}
	var payCoin float64
	// TODO: 英文短信
	reCoin := regexp.MustCompile(`([\d,]+\.\d{2})บ\.`)
	if m := reCoin.FindStringSubmatch(remark); len(m) > 1 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(m[1], ",", ""), 64)
	}
	payTime, ok := parser.MatchDayMonthYearTimeUse2(remark)
	if ok {
		res.PayTime = payTime.Unix()
	}
	balance, ok := parser.MatchBalanceAll(remark)
	if ok {
		res.Balance = balance
	}
	res.PayCoin = payCoin
	return
}

// TTB读取/TTB通知
func ExtractTTBRead(remark string) *types.Extracted {
	return ExtractTTB(remark)
}
func ExtractTTBNotify(remark string) *types.Extracted {
	return ExtractTTB(remark)
}

// KKRLSCLI
func ExtractKKRLSCLI(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	var obj struct {
		Coin float64 `json:"coin"`
		Time int64   `json:"time"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	return
}

// BAY
func ExtractBAY(remark string) (res *types.Extracted) {
	//โอนเข้า xxx170203x  150.62 เหลือ 7,420.57 (17/5/68,16:27)
	res = &types.Extracted{}
	if remark == "" {
		res.Code = codes.ErrInvalidArgument
		return
	}
	if !strings.HasPrefix(remark, "โอนเข้า") && !strings.Contains(remark, "Money Deposit") {
		res.Code = codes.ErrUnsupportedMessageFormat.New("不是指定数据类型。")
		return
	}
	arr := strings.Fields(remark)
	if strings.Contains(remark, "Money Deposit") {
		res.Code = codes.ErrUnimplemented.New("Money Deposit unimplemented")
		return
	}
	if arr[0] == "โอนเข้า" && len(arr) >= 6 {
		res.PayCoin, _ = strconv.ParseFloat(strings.ReplaceAll(arr[2], ",", ""), 64)
		// 时间
		t, ok := parser.MatchThaiDateParen(remark)
		if ok {
			res.PayTime = t.Unix()
		}
		// 余额
		b, ok := parser.MatchBalanceAll(remark)
		if ok {
			res.Balance = b
		}
	}
	return
}

// TM流水ios
func ExtractTMIosWater(ramk string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	// 暂不处理，无数据
	/* var obj struct {
		Time  string `json:"time"`
		Money string `json:"money"`
	}
	_ = json.Unmarshal([]byte(ramk), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	timeArr := strings.Fields(obj.Time)
	if len(timeArr) < 4 {
		return
	}
	year, _ := strconv.Atoi(timeArr[2])
	year = year - 543
	month := thaiMonthToNum(timeArr[1])
	payTimeStr := strconv.Itoa(year) + "-" + month + "-" + timeArr[0] + " " + timeArr[3]
	t, _ := time.ParseInLocation("2006-01-02 15:04", payTimeStr, time.Local) */
	return
}

// SwooleTM
func ExtractSwooleTM(coin string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* coin = strings.ReplaceAll(coin, " ", "")
	coin = strings.ReplaceAll(coin, "฿", "")
	coin = strings.ReplaceAll(coin, ",", "")
	payCoin, _ := strconv.ParseFloat(coin, 64) */
	return
}

// GSB
func ExtractGSB(ramk string) (res *types.Extracted) {
	res = &types.Extracted{}
	remark := strings.Fields(ramk)
	var payCoin float64
	var payTime int64
	var balance float64
	if len(remark) == 0 {
		res.Code = codes.ErrInvalidArgument
		return
	}
	if remark[0] != "คุณได้รับเงิน" && remark[0] != "เงินเข้า:" && remark[0] != "Deposit" {
		res.Code = codes.ErrUnsupportedMessageFormat
		return
	}
	// 过滤
	if remark[0] == "คุณได้รับเงิน" {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(remark[1], ",", ""), 64)
	} else if len(remark) > 2 {
		payCoin, _ = strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(remark[2], ",", ""), "฿", ""), 64)
	}
	re := regexp.MustCompile(`วันที่ (\d+) (.*?) (\d+) เวลา (\d+):(\d+) น`)
	m := re.FindStringSubmatch(ramk)
	if len(m) == 6 {
		month := thaiMonthToNum(m[2])
		year := time.Now().Year()
		payTimeStr := strconv.Itoa(year) + "-" + month + "-" + m[1] + " " + m[4] + ":" + m[5] + ":00"
		t, err := time.Parse(time.DateTime, payTimeStr)
		if err != nil {
			logger.Warnf("Parse time failed|error=%v", err)
		}
		if err == nil {
			payTime = t.Unix()
		}
	}
	if len(remark) > 18 {
		balance, _ = strconv.ParseFloat(strings.ReplaceAll(remark[18], ",", ""), 64)
	}
	res.Code = codes.Success
	res.PayCoin = payCoin
	res.PayTime = payTime
	res.Balance = balance
	return
}

// GSB读取/GSB通知/GSBLine
func ExtractGSBRead(ramk string) *types.Extracted   { return ExtractGSB(ramk) }
func ExtractGSBNotify(ramk string) *types.Extracted { return ExtractGSB(ramk) }
func ExtractGSBLine(ramk string) *types.Extracted   { return ExtractGSB(ramk) }

// python-SCBGH
func ExtractSCBGH(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* var obj struct {
		Money string `json:"money"`
		Time  string `json:"time"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Money, ",", ""), 64)
	t, _ := time.ParseInLocation("2006-01-02 15:04", obj.Time, time.Local) */
	return
}

// python-KTBGH
func ExtractKTBGH(postMsg string) (res *types.Extracted) {
	return ExtractSCBGH(postMsg)
}

// python-Kbankgh
func ExtractKbankGH(postMsg string) (res *types.Extracted) {
	return ExtractSCBGH(postMsg)
}

// python-Ttbgh
func ExtractTtbGH(postMsg string) (res *types.Extracted) {
	return ExtractSCBGH(postMsg)
}

// python-BAYGH
func ExtractBayGH(postMsg string) (res *types.Extracted) {
	return ExtractSCBGH(postMsg)
}

// tm_protocol_water
func ExtractTMProtocolWater(postMsg string) (res *types.Extracted) {
	res = &types.Extracted{
		Code: codes.ErrUnimplemented,
	}
	/* var obj struct {
		DateTime string `json:"date_time"`
		Amount   string `json:"amount"`
	}
	_ = json.Unmarshal([]byte(postMsg), &obj)
	payCoin, _ := strconv.ParseFloat(strings.ReplaceAll(obj.Amount, ",", ""), 64)
	t, _ := time.ParseInLocation("2006-01-02 15:04", obj.DateTime, time.Local) */
	return
}

// scb_protocol_water
func ExtractSCBProtocolWater(postMsg string) (res *types.Extracted) {
	return ExtractTMProtocolWater(postMsg)
}

// ttb_protocol_water
func ExtractTTBProtocolWater(postMsg string) (res *types.Extracted) {
	return ExtractTMProtocolWater(postMsg)
}

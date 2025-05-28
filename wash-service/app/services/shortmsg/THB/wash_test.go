package THB

import (
	"testing"

	"github.com/cg917658910/fzkj-wallet/wash-service/app/services/shortmsg/codes"
)

func TestWashGSB(t *testing.T) {
	//msg1
	msg1 := `เงินเข้า: มีการฝาก/โอนเงิน 74.76 บาท จากบัญชี KBNK 0013XXXX1945 เข้าบัญชี GSBA 0204XXXX0743 วันที่ 24 พ.ค. 2568 เวลา 15:04 น. คงเหลือ 2,702.93 บาท`
	result := ExtractGSB(msg1)
	msg1ExpectedCoin := 74.76
	msg1ExpectedBalance := 2702.93
	if result.PayCoin != msg1ExpectedCoin {
		t.Errorf("wash GSB Coin expected %v, but got %v", msg1ExpectedCoin, result.PayCoin)
	}
	if result.Balance != msg1ExpectedBalance {
		t.Errorf("wash GSB Balance expected %v, but got %v", msg1ExpectedBalance, result.Balance)
	}
	// msg2
	msg2 := `เงินเข้า: มีการฝาก/โอนเงิน 10.00 บาท จากบัญชี KBNK`
	result2 := ExtractGSB(msg2)
	msg2ExpectedCoin := 10.00
	if result2.PayCoin != msg2ExpectedCoin {
		t.Errorf("wash GSB Coin expected %v, but got %v", msg2ExpectedCoin, result.PayCoin)
	}

}

type testMsg struct {
	con           string
	expectCoin    float64
	expectBalance float64
	expectTime    int64
	expectCode    codes.Code
}

func TestWashSCB(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `18/05@14:44 29.92 จากKTB/x129439เข้าx997010 ใช้ได้10,121.23บ`,
			expectCoin:    29.92,
			expectBalance: 10121.23,
			expectTime:    1747550640,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    399.69,
			expectBalance: 6302.39,
			expectTime:    1747550700,
		},
		{
			con:           `ถอน/โอนเงิน 10,000.00บ บ/ชx833473 18/05@12:29 ใช้ได้657.53บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			expectCode:    codes.ErrUnsupportedMessageFormat,
		},
	}

	for _, msg := range msgs {
		res := ExtractSCB(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("wash SCB code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("wash SCB time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("wash SCB coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("wash SCB balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractSCBNotify(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `คุณได้รับเงิน 69.89 บาท ผ่านรายการพร้อมเพย์จาก KTB / xxxxxx6048 เข้าบัญชี xxxxxx1083 เมื่อ 17 พ.ค. 2568 - 16:28`,
			expectCoin:    69.89,
			expectBalance: 0,
			expectTime:    1747470480,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
		},
	}

	for _, msg := range msgs {
		res := ExtractSCBNotify(msg.con)
		if res.Code.Code != codes.Success.Code {
			t.Errorf("wash SCBNotify code expected %v, but got %v", codes.Success.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("wash SCBNotify time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("wash SCBNotify coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("wash SCBNotify balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractBAAC(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `20/05/68 17:52:48 บัญชี x9576 รับโอนพร้อมเพย์ จำนวน 10.00 บาท ยอดเงินคงเหลือ 10.39 บาท`,
			expectCoin:    10.00,
			expectBalance: 10.39,
			expectTime:    1747734768,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
		},
	}

	for _, msg := range msgs {
		res := ExtractBAAC(msg.con)
		if res.Code.Code != codes.Success.Code {
			t.Errorf("wash SCBNotify code expected %v, but got %v", codes.Success.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("wash SCBNotify time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("wash SCBNotify coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("wash SCBNotify balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractBAY(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `โอนเข้า xxx170203x  150.62 เหลือ 7,420.57 (17/5/68,16:27)`,
			expectCoin:    150.62,
			expectBalance: 7420.57,
			expectTime:    1747470420,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			expectCode:    codes.ErrUnsupportedMessageFormat,
		},
	}

	for _, msg := range msgs {
		res := ExtractBAY(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractBBL(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `เงินเข้าบ/ชX0272 จากพร้อมเพย์ผ่านMB 12,000.00บ เงินในบ/ชใช้ได้46,769.66บ`,
			expectCoin:    12000.00,
			expectBalance: 46769.66,
			expectTime:    0,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			//expectCode:    0,
		},
	}

	for _, msg := range msgs {
		res := ExtractBBL(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractKBANKRead(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `17/05/68 16:27 บช X-9426 เงินเข้า 999.63 คงเหลือ 8,236.90 บ.`,
			expectCoin:    999.63,
			expectBalance: 8236.90,
			expectTime:    1747470420,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			//expectCode:    0,
		},
	}
	for _, msg := range msgs {
		res := ExtractKBANKRead(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractKTB(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `17-05@16:27 บชX58589X:เงินเข้า 100.28บ ใช้ได้ 11,861.55บ`,
			expectCoin:    100.28,
			expectBalance: 11861.55,
			expectTime:    1747470420,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 6302.39,
			expectTime:    1747550700,
			//expectCode:    0,
		},
	}
	for _, msg := range msgs {
		res := ExtractKTB(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractKTBLine(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `เงินเข้า: 1,000.20 บาท เข้าบัญชี XX5568 เมื่อ 27/0`,
			expectCoin:    1000.20,
			expectBalance: 0,
			expectTime:    0,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน s399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			//expectCode:    0,
		},
	}
	for _, msg := range msgs {
		res := ExtractKTBLine(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractKTBNotice(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `ได้รับ +59.84 บาท เข้าพร้อมเพย์ จากบัญชี ไทยพาณิชย์ XXX-X-XX805-8`,
			expectCoin:    59.84,
			expectBalance: 0,
			expectTime:    0,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน s399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 0,
			expectTime:    0,
			//expectCode:    0,
		},
	}
	for _, msg := range msgs {
		res := ExtractKTBLine(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}
func TestExtractTTB(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `มีเงิน50.14บ.โอนเข้าบ/ชxx5440เหลือ17,081.71บ.24/05/25@11:42`,
			expectCoin:    50.14,
			expectBalance: 17081.71,
			expectTime:    1748058120,
			expectCode:    codes.Success,
		},
		{
			con:           `เงิน s399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302.39บ`,
			expectCoin:    0,
			expectBalance: 6302.39,
			expectTime:    0,
			//expectCode:    0,
		},
	}
	for _, msg := range msgs {
		res := ExtractTTB(msg.con)
		if res.Code.Code != msg.expectCode.Code {
			t.Errorf("code expected %v, but got %v", msg.expectCode.Code, res.Code.Code)
		}
		if res.PayTime != msg.expectTime {
			t.Errorf("time expected %v, but got %v", msg.expectTime, res.PayTime)
		}
		if res.PayCoin != msg.expectCoin {
			t.Errorf("coin expected %v, but got %v", msg.expectCoin, res.PayCoin)
		}
		if res.Balance != msg.expectBalance {
			t.Errorf("balance expected %v, but got %v", msg.expectBalance, res.Balance)
		}
	}

}

package parser

import (
	"fmt"
	"testing"
)

func TestParserTimeMatchDayMonthAtTime(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `18/05@14:44 29.92 จากKTB/x129439เข้าx997010 ใช้ได้10,121.23บ`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchDayMonthAtTime(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("time: ", r1.String())
		}
		t.Error("done")
	}
}
func TestParserTimeMatchThaiDateParen(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `โอนเข้า xxx170203x  50.28 เหลือ 7,666.30 (17/5/68,16:30)`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchThaiDateParen(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("time: ", r1.String())
		}
		t.Error("done")
	}
}
func TestParserTimeMatchThaiTextDate(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `เงินเข้าบ/ชX1929 จากพร้อมเพย์ผ่านMB 10,000.00บ เงินในบ/ชใช้ได้75,509.90บ`,
			expectBalance: 10,
		},
		{
			con:           `คุณได้รับเงิน 49.74 บาท ผ่านรายการพร้อมเพย์จาก KTB / xxxxxx9424 เข้าบัญชี xxxxxx1083 เมื่อ 17 พ.ค. 2568 - 16:28`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchThaiTextDate(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("time: ", r1.String())
		}
		t.Error("done")
	}
}
func TestParserTimeMatchDayMonthYearTime(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `17/05/68 16:28 บช X-9426 เงินเข้า 100.27 คงเหลือ 8,437.27 บ.`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
		{
			con:           `18/05@14:44 29.92 จากKTB/x129439เข้าx997010 ใช้ได้10,121.23บ`,
			expectBalance: 10,
		},
		{
			con:           `20/05/68 17:52:48 บัญชี x9576 รับโอนพร้อมเพย์ จำนวน 10.00 บาท ยอดเงินคงเหลือ 10.39 บาท`,
			expectBalance: 10,
		},
		{
			con:           `20/05/23 17:52:48 บัญชี x9576 รับโอนพร้อมเพย์ จำนวน 10.00 บาท ยอดเงินคงเหลือ 10.39 บาท`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchDayMonthYearTime(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("time: ", r1.String())
		}
		t.Error("done")
	}
}

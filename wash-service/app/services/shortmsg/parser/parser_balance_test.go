package parser

import (
	"fmt"
	"testing"
)

type testMsg struct {
	con           string
	expectCoin    float64
	expectBalance float64
	expectTime    int64
}

func TestParserBalanceMatchBalanceUsed(t *testing.T) {
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

		r1, b1 := MatchBalanceUsed(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("balance: ", r1)
		}
		t.Error("done")
	}
}
func TestParserBalanceMatchBalanceRemain(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `17/05/68 16:28 บช X-9426 เงินเข้า 100.27 คงเหลือ 8,437.27 บ.`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchBalanceRemain(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("balance: ", r1)
		}
		t.Error("done")
	}
}
func TestParserBalanceMatchBalanceUsed2(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `เงินเข้าบ/ชX1929 จากพร้อมเพย์ผ่านMB 10,000.00บ เงินในบ/ชใช้ได้75,509.90บ`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchBalanceUsed2(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("balance: ", r1)
		}
		t.Error("done")
	}
}
func TestParserBalanceMatchAll(t *testing.T) {
	var msgs = []testMsg{
		{
			con:           `เงินเข้าบ/ชX1929 จากพร้อมเพย์ผ่านMB 10,000.00บ เงินในบ/ชใช้ได้75,509.90บ`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
		{
			con:           `เงินเข้าบ/ชX1929 จากพร้อมเพย์ผ่านMB 10,000.00บ เงินในบ/ชใช้ได้75,509.90บ`,
			expectBalance: 10,
		},
		{
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
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
			con:           `เงิน 399.69บ เข้าบ/ชx963071 18/05@14:45 ใช้ได้ 6,302บ`,
			expectBalance: 10,
		},
	}

	for _, msg := range msgs {

		r1, b1 := MatchBalanceAll(msg.con)
		fmt.Println("bool", b1)
		if b1 {
			fmt.Println("balance: ", r1)
		}
		t.Error("done")
	}
}

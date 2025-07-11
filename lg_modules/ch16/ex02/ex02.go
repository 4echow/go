package main

import (
	"fmt"
	"unsafe"
)

type OrderInfo struct {
	OrderCode   rune
	Amount      int
	OrderNumber uint16
	Items       []string
	IsReady     bool
}

// OrderCode, OrderNumber and IsReady fit in an 8-byte block with only 1-byte buffer!
// Amount and Items are multiples of 8-bytes, so are aligned
type SmallOrderInfo struct {
	OrderCode   rune
	OrderNumber uint16
	IsReady     bool
	Amount      int
	Items       []string
}

func main() {
	fmt.Println("OrderInfo:",
		unsafe.Sizeof(OrderInfo{}),
		unsafe.Offsetof(OrderInfo{}.OrderCode),
		unsafe.Offsetof(OrderInfo{}.Amount),
		unsafe.Offsetof(OrderInfo{}.OrderNumber),
		unsafe.Offsetof(OrderInfo{}.Items),
		unsafe.Offsetof(OrderInfo{}.IsReady),
	)
	fmt.Println("SmallOrderInfo:",
		unsafe.Sizeof(SmallOrderInfo{}),
		unsafe.Offsetof(SmallOrderInfo{}.OrderCode),
		unsafe.Offsetof(SmallOrderInfo{}.Amount),
		unsafe.Offsetof(SmallOrderInfo{}.OrderNumber),
		unsafe.Offsetof(SmallOrderInfo{}.Items),
		unsafe.Offsetof(SmallOrderInfo{}.IsReady),
	)
}

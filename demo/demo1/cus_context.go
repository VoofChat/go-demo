package demo1

import (
	"encoding/json"
	"fmt"
)

type Context interface {
	Println(obj interface{})
}

type GinContext struct {
}

func (c *GinContext) Println(obj interface{}) {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		fmt.Println(obj)
	} else {
		fmt.Println(string(jsonStr))
	}
}

func withContext(ctx Context) {
	ctx.Println("withContext")
}

func dbQuery1(ctx *GinContext) {
	withContext(ctx)
}

//func dbQuery2(ctx *Context) { // 不要写成这种方式 golang 中interface本身就是指针类型
func dbQuery2(ctx Context) {
	withContext(ctx)
}

// 以下示例说明了interface本身是指针的概念

type Rect struct {
	Width  int
	Height int
}

func ex1() {
	var a interface{}
	var r = Rect{50, 50}
	a = &r // 指向了结构体指针

	var rx = a.(*Rect) // 转换成指针类型
	r.Width = 100
	r.Height = 100
	fmt.Println("r:", r)                // r: {100 100}
	fmt.Println("rx:", rx)              // rx: &{100 100}
	fmt.Printf("rx:%p, r:%p\n", rx, &r) // rx:0xc0000160c0, r:0xc0000160c0
}

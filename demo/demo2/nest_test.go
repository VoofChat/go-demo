package demo2

import "testing"

// 接口嵌套接口
func Test_InterInInter(t *testing.T) {
	InterInInterTest()
}

// 结构体嵌套接口
func Test_StructInInter(t *testing.T) {
	StructInInterTest1()
	StructInInterTest2()
}

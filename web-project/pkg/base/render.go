package base

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type Render interface {
	SetReturnCode(int)
	SetReturnMsg(string)
	SetReturnData(interface{})
	GetReturnCode() int
	GetReturnMsg() string
}

var newRender func() Render

func RegisterRender(s func() Render) {
	newRender = s
}

func newJsonRender() Render {
	if newRender == nil {
		newRender = defaultNew
	}
	return newRender()
}

// default render

var defaultNew = func() Render {
	return &DefaultRender{}
}

type DefaultRender struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func (r *DefaultRender) GetReturnCode() int {
	return r.ErrNo
}
func (r *DefaultRender) SetReturnCode(code int) {
	r.ErrNo = code
}
func (r *DefaultRender) GetReturnMsg() string {
	return r.ErrMsg
}
func (r *DefaultRender) SetReturnMsg(msg string) {
	r.ErrMsg = msg
}
func (r *DefaultRender) GetReturnData() interface{} {
	return r.Data
}
func (r *DefaultRender) SetReturnData(data interface{}) {
	r.Data = data
}

func RenderJsonSucc(ctx *gin.Context, data interface{}) {
	r := newJsonRender()
	r.SetReturnCode(0)
	r.SetReturnMsg("succ")
	r.SetReturnData(data)
	ctx.JSON(http.StatusOK, r)
	return
}

func RenderJsonFail(ctx *gin.Context, err error) {
	r := newJsonRender()

	code, msg := -1, errors.Cause(err).Error()
	switch errors.Cause(err).(type) {
	case Error:
		code = errors.Cause(err).(Error).ErrNo
		msg = errors.Cause(err).(Error).ErrMsg
	default:
	}

	r.SetReturnCode(code)
	r.SetReturnMsg(msg)
	r.SetReturnData(gin.H{})

	ctx.JSON(http.StatusOK, r)

	// 打印错误栈
	StackLogger(ctx, err)
	return
}

// 打印错误栈
func StackLogger(ctx *gin.Context, err error) {
	if !strings.Contains(fmt.Sprintf("%+v", err), "\n") {
		return
	}

	// todo 错误栈输出格式优化
	fmt.Printf("-------------------stack-start-------------------\n%+v\n-------------------stack-end-------------------\n", err)
}

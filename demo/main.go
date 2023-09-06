package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	//"net/http"
)

type DefaultRender struct {
	ErrNo  int         `json:"errNo"`
	ErrMsg string      `json:"errMsg"`
	Data   interface{} `json:"data"`
}

func defaultRenderJson() {

}

func main() {
	obj := DefaultRender{
		ErrNo:  -1,
		ErrMsg: "no cuid",
		Data:   gin.H{},
	}

	data, err := json.Marshal(&obj)
	if err != nil {
		//panic(err)
		println(err.Error())
		return
	}
	fmt.Println(string(data))
	//r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//go get github.com/json-iterator/go

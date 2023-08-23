package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
{
	"code" : 10001, // 程序中的错误码
	"msg" : xx， // 提示信息
	"data" : {}, // 数据
}
*/

type ResponseData struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func ResponseError(c *gin.Context, code ResCode) {

	//gin.H{
	//	"code": "xx",
	//	"msg":  "xx",
	//	"date": "xx",
	//}

	rd := &ResponseData{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)

	//c.JSON(http.StatusOK, &ResponseData{
	//	Code: code,
	//	Msg:  code.Msg(),
	//	Data: nil,
	//})
}

func ResponseErrorWithMsg(c *gin.Context, code ResCode, msg interface{}) {
	rd := &ResponseData{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
	c.JSON(http.StatusOK, rd)
}

func ResponseSuccess(c *gin.Context, data interface{}) {

	//gin.H{
	//	"code": "xx",
	//	"msg":  "xx",
	//	"date": "xx",
	//}

	rd := &ResponseData{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	}
	c.JSON(http.StatusOK, rd)

	//c.JSON(http.StatusOK, &ResponseData{
	//	Code: CodeSuccess,
	//	Msg:  CodeSuccess.Msg(),
	//	Data: data,
	//})
}

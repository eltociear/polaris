package handler

import "github.com/gin-gonic/gin"

type HttpResponse struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func Response(c *gin.Context, httpCode, errCode int, data interface{}) {
	c.JSON(httpCode, HttpResponse{
		Code: errCode,
		Msg:  GetMsg(errCode),
		Data: data,
	})
	return
}

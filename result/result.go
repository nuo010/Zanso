package result

import (
	"fmt"
	"github.com/dromara/carbon/v2"
	"github.com/gin-gonic/gin"
)

// Response 基础序列化器
type Response struct {
	Code   int         `json:"code"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
	Status bool        `json:"status"`
	Time   string      `json:"time"`
}

const (
	SUCCESS = 200
	FAIL    = 500
)

func Ok(c *gin.Context) {
	res := Response{
		Code:   SUCCESS,
		Msg:    "操作成功",
		Status: true,
		Time:   carbon.Now().ToDateTimeString(),
	}
	c.JSON(SUCCESS, res)
}
func OkSetData(c *gin.Context, data interface{}) {
	res := Response{
		Code:   SUCCESS,
		Data:   data,
		Msg:    "操作成功",
		Status: true,
		Time:   carbon.Now().ToDateTimeString(),
	}
	c.JSON(SUCCESS, res)
}
func ErrSetMsg(c *gin.Context, msg string) {
	res := Response{
		Code:   FAIL,
		Msg:    msg,
		Status: false,
		Time:   carbon.Now().ToDateTimeString(),
	}
	c.JSON(SUCCESS, res)
}

// Err 通用错误处理
func Err(c *gin.Context, msg string, err error) {
	res := Response{
		Code:   FAIL,
		Msg:    msg,
		Status: false,
		Time:   carbon.Now().ToDateTimeString(),
	}
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	c.JSON(SUCCESS, res)
}

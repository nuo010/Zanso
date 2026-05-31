package service

import (
	"log"
	"time"
	"zanso/db"
	"zanso/model"
	"zanso/result"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func GetUrl(c *gin.Context) {
	var tblfq model.TblFq
	key := c.Query("key")
	log.Println("开始获取订阅", key)
	if key == "" {
		db.DB.Where("mr = ?", "1").Take(&tblfq)
	} else {
		db.DB.Where("key_name = ?", key).Take(&tblfq)
	}
	if tblfq.ID == "" {

		result.ErrSetMsg(c, "key错误")
		return
	}

	client := resty.New().
		SetTimeout(15*time.Second).
		SetHeader("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36").
		SetHeader("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7").
		SetHeader("Accept-Language", "zh-CN,zh;q=0.9").
		SetHeader("Upgrade-Insecure-Requests", "1").
		SetHeader("Priority", "u=0, i").
		SetHeader("Sec-CH-UA", `"Not(A:Brand";v="8", "Chromium";v="144", "Google Chrome";v="144"`).
		SetHeader("Sec-CH-UA-Mobile", "?0").
		SetHeader("Sec-CH-UA-Platform", `"macOS"`).
		SetHeader("Sec-Fetch-Dest", "document").
		SetHeader("Sec-Fetch-Mode", "navigate").
		SetHeader("Sec-Fetch-Site", "none").
		SetHeader("Sec-Fetch-User", "?1")

	resp, err := client.R().Get(tblfq.Url)
	if err != nil {

		result.ErrSetMsg(c, "dy错误"+err.Error())
		return
	}

	// 原封不动返回
	c.Status(resp.StatusCode())
	for k, v := range resp.Header() {
		c.Writer.Header()[k] = v
	}
	_, _ = c.Writer.Write(resp.Body())

}

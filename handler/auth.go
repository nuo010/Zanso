package handler

import (
	"errors"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func UserAuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := util.ExtractUserToken(c)
		if token == "" {
			result.ErrSetMsg(c, "未登录或登录令牌缺失")
			c.Abort()
			return
		}

		userID, err := db.RedisClient.Get(c.Request.Context(), "auth:user:token:"+token).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				result.ErrSetMsg(c, "登录状态已失效，请重新登录")
			} else {
				result.ErrSetMsg(c, "登录状态校验失败")
			}
			c.Abort()
			return
		}

		var user model.User
		if err = db.DB.Where("id = ? AND status = ?", userID, model.UserStatusActive).Take(&user).Error; err != nil {
			result.ErrSetMsg(c, "用户状态异常，请重新登录")
			c.Abort()
			return
		}

		c.Set(util.ContextUserKey, user)
		c.Next()
	}
}

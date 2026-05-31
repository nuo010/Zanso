package service

import (
	"context"
	"strings"
	"time"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

const (
	userTokenPrefix = "auth:user:token:"
)

func LoginUser(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "登录参数错误")
		return
	}

	loginName := strings.TrimSpace(req.LoginName)
	password := strings.TrimSpace(req.Password)
	if loginName == "" || password == "" {
		result.ErrSetMsg(c, "登录账号和密码不能为空")
		return
	}

	var user model.User
	if err := db.DB.Where("login_name = ? AND status = ?", loginName, model.UserStatusActive).Take(&user).Error; err != nil {
		result.ErrSetMsg(c, "账号或密码错误")
		return
	}
	if !util.CheckPassword(password, user.PasswordHash) {
		result.ErrSetMsg(c, "账号或密码错误")
		return
	}

	token, err := createUserToken(c.Request.Context(), user.ID)
	if err != nil {
		result.ErrSetMsg(c, "创建登录态失败")
		return
	}

	result.OkSetData(c, model.UserAuthResponse{
		User:  user,
		Token: token,
	})
}

func LogoutUser(c *gin.Context) {
	token := util.ExtractUserToken(c)
	if token == "" {
		result.ErrSetMsg(c, "未获取到登录令牌")
		return
	}
	if err := deleteUserToken(c.Request.Context(), token); err != nil {
		result.ErrSetMsg(c, "退出登录失败")
		return
	}
	result.Ok(c)
}

func GetCurrentUserProfile(c *gin.Context) {
	user, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	result.OkSetData(c, user)
}

func createUserToken(ctx context.Context, userID string) (string, error) {
	token := util.GetUuid() + util.GetUuid()[:16]
	ttl := getUserTokenTTL()
	key := userTokenPrefix + token
	if err := db.RedisClient.Set(ctx, key, userID, ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func deleteUserToken(ctx context.Context, token string) error {
	return db.RedisClient.Del(ctx, userTokenPrefix+token).Err()
}

func getUserTokenTTL() time.Duration {
	hours := viper.GetInt("auth.session_ttl_hours")
	if hours <= 0 {
		hours = 168
	}
	return time.Duration(hours) * time.Hour
}

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
	merchantTokenPrefix = "auth:merchant:token:"
)

func LoginMerchant(c *gin.Context) {
	var req model.MerchantLoginRequest
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

	var merchant model.Merchant
	if err := db.DB.Where("login_name = ? AND status = ?", loginName, model.MerchantStatusActive).Take(&merchant).Error; err != nil {
		result.ErrSetMsg(c, "账号或密码错误")
		return
	}
	if !util.CheckPassword(password, merchant.PasswordHash) {
		result.ErrSetMsg(c, "账号或密码错误")
		return
	}

	token, err := createMerchantToken(c.Request.Context(), merchant.ID)
	if err != nil {
		result.ErrSetMsg(c, "创建登录态失败")
		return
	}

	result.OkSetData(c, model.MerchantAuthResponse{
		Merchant: merchant,
		Token:    token,
	})
}

func LogoutMerchant(c *gin.Context) {
	token := util.ExtractMerchantToken(c)
	if token == "" {
		result.ErrSetMsg(c, "未获取到登录令牌")
		return
	}
	if err := deleteMerchantToken(c.Request.Context(), token); err != nil {
		result.ErrSetMsg(c, "退出登录失败")
		return
	}
	result.Ok(c)
}

func GetCurrentMerchantProfile(c *gin.Context) {
	merchant, ok := util.GetCurrentMerchant(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	result.OkSetData(c, merchant)
}

func createMerchantToken(ctx context.Context, merchantID string) (string, error) {
	token := util.GetUuid() + util.GetUuid()[:16]
	ttl := getMerchantTokenTTL()
	key := merchantTokenPrefix + token
	if err := db.RedisClient.Set(ctx, key, merchantID, ttl).Err(); err != nil {
		return "", err
	}
	return token, nil
}

func deleteMerchantToken(ctx context.Context, token string) error {
	return db.RedisClient.Del(ctx, merchantTokenPrefix+token).Err()
}

func GetMerchantIDByToken(ctx context.Context, token string) (string, error) {
	merchantID, err := db.RedisClient.Get(ctx, merchantTokenPrefix+token).Result()
	if err != nil {
		return "", err
	}
	ttl := getMerchantTokenTTL()
	if ttl > 0 {
		_ = db.RedisClient.Expire(ctx, merchantTokenPrefix+token, ttl).Err()
	}
	return merchantID, nil
}

func getMerchantTokenTTL() time.Duration {
	hours := viper.GetInt("auth.session_ttl_hours")
	if hours <= 0 {
		hours = 168
	}
	return time.Duration(hours) * time.Hour
}

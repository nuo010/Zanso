package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"zanso/db"
	"zanso/model"
	"zanso/result"
	"zanso/util"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func CreateMerchant(c *gin.Context) {
	var req model.CreateMerchantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "商家参数错误")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "商家名称不能为空")
		return
	}
	if strings.TrimSpace(req.LoginName) == "" || strings.TrimSpace(req.Password) == "" {
		result.ErrSetMsg(c, "登录账号和密码不能为空")
		return
	}

	passwordHash, err := util.HashPassword(strings.TrimSpace(req.Password))
	if err != nil {
		result.ErrSetMsg(c, "密码加密失败")
		return
	}

	merchant := model.Merchant{
		ID:           util.GetUuid(),
		Name:         strings.TrimSpace(req.Name),
		LoginName:    strings.TrimSpace(req.LoginName),
		PasswordHash: passwordHash,
		ContactName:  strings.TrimSpace(req.ContactName),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Status:       model.MerchantStatusActive,
		CreatedAt:    util.GetTime(),
		UpdatedAt:    util.GetTime(),
	}
	if err = db.DB.Create(&merchant).Error; err != nil {
		result.ErrSetMsg(c, "创建商家失败")
		return
	}
	result.OkSetData(c, merchant)
}

func GetMerchantList(c *gin.Context) {
	currentMerchant, ok := util.GetCurrentMerchant(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var merchantList []model.Merchant
	if err := db.DB.Where("id = ?", currentMerchant.ID).Order("created_at desc").Find(&merchantList).Error; err != nil {
		result.ErrSetMsg(c, "查询商家失败")
		return
	}
	result.OkSetData(c, merchantList)
}

func CreateProduct(c *gin.Context) {
	var req model.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "商品参数错误")
		return
	}
	currentMerchant, ok := util.GetCurrentMerchant(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		result.ErrSetMsg(c, "商品标题不能为空")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.ProductStatusDraft
	}

	product := model.Product{
		ID:          util.GetUuid(),
		MerchantID:  currentMerchant.ID,
		Title:       strings.TrimSpace(req.Title),
		Description: strings.TrimSpace(req.Description),
		CoverURL:    strings.TrimSpace(req.CoverURL),
		Status:      status,
		CreatedAt:   util.GetTime(),
		UpdatedAt:   util.GetTime(),
	}
	if err := db.DB.Create(&product).Error; err != nil {
		result.ErrSetMsg(c, "创建商品失败")
		return
	}
	result.OkSetData(c, product)
}

func GetMerchantProductList(c *gin.Context) {
	merchantID := c.Param("merchantId")
	currentMerchant, ok := util.GetCurrentMerchant(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if merchantID == "" || merchantID != currentMerchant.ID {
		result.ErrSetMsg(c, "无权查看其他商家的商品")
		return
	}

	var productList []model.Product
	if err := db.DB.Where("merchant_id = ?", merchantID).Order("created_at desc").Find(&productList).Error; err != nil {
		result.ErrSetMsg(c, "查询商品失败")
		return
	}
	result.OkSetData(c, productList)
}

func GetProductDetail(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		result.ErrSetMsg(c, "商品 ID 不能为空")
		return
	}
	currentMerchantID := util.CurrentMerchantID(c)
	if currentMerchantID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var product model.Product
	if err := db.DB.Where("id = ? AND merchant_id = ?", productID, currentMerchantID).Take(&product).Error; err != nil {
		result.ErrSetMsg(c, "商品不存在")
		return
	}

	var merchant model.Merchant
	var mediaList []model.MediaAsset
	var shareLinks []model.ShareLink
	db.DB.Where("id = ?", product.MerchantID).Take(&merchant)
	db.DB.Where("product_id = ?", productID).Order("sort asc, created_at asc").Find(&mediaList)
	db.DB.Where("product_id = ?", productID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, model.ProductDetail{
		Product:    product,
		Merchant:   merchant,
		MediaList:  mediaList,
		ShareLinks: shareLinks,
	})
}

func UploadProductMedia(c *gin.Context) {
	productID := c.Param("id")
	if productID == "" {
		result.ErrSetMsg(c, "商品 ID 不能为空")
		return
	}
	currentMerchantID := util.CurrentMerchantID(c)
	if currentMerchantID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var product model.Product
	if err := db.DB.Where("id = ? AND merchant_id = ?", productID, currentMerchantID).Take(&product).Error; err != nil {
		result.ErrSetMsg(c, "商品不存在")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		result.ErrSetMsg(c, "请上传文件")
		return
	}

	mediaType := strings.TrimSpace(c.PostForm("mediaType"))
	if mediaType == "" {
		mediaType = inferMediaType(file.Header.Get("Content-Type"), file.Filename)
	}

	sortValue, _ := strconv.Atoi(c.DefaultPostForm("sort", "0"))
	now := util.GetTime()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	relativePath := filepath.Join(
		"products",
		now.Format("2006"),
		now.Format("01"),
		productID+"_"+util.GetUuid()[:8]+ext,
	)
	localDir := getUploadLocalDir()
	fullPath := filepath.Join(localDir, relativePath)

	if err = os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		result.ErrSetMsg(c, "创建上传目录失败")
		return
	}
	if err = c.SaveUploadedFile(file, fullPath); err != nil {
		result.ErrSetMsg(c, "保存文件失败")
		return
	}

	publicURL := buildPublicAssetURL(c, relativePath)
	fileRecord := model.TblFile{
		ID:        util.GetUuid(),
		CreatTime: now,
		Fsize:     file.Size,
		Key:       publicURL,
		Type:      mediaType,
	}
	asset := model.MediaAsset{
		ID:         util.GetUuid(),
		MerchantID: product.MerchantID,
		ProductID:  productID,
		MediaType:  mediaType,
		FileName:   file.Filename,
		FileSize:   file.Size,
		MimeType:   file.Header.Get("Content-Type"),
		URL:        publicURL,
		Sort:       sortValue,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if errTx := tx.Create(&fileRecord).Error; errTx != nil {
			return errTx
		}
		if errTx := tx.Create(&asset).Error; errTx != nil {
			return errTx
		}
		if product.CoverURL == "" && mediaType == "image" {
			if errTx := tx.Model(&model.Product{}).Where("id = ?", productID).Updates(map[string]interface{}{
				"cover_url":  publicURL,
				"updated_at": now,
			}).Error; errTx != nil {
				return errTx
			}
		}
		return nil
	}); err != nil {
		result.ErrSetMsg(c, "保存媒体资源失败")
		return
	}

	result.OkSetData(c, asset)
}

func CreateShareLink(c *gin.Context) {
	var req model.CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分享参数错误")
		return
	}
	currentMerchantID := util.CurrentMerchantID(c)
	if currentMerchantID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if req.ProductID == "" {
		result.ErrSetMsg(c, "商品 ID 不能为空")
		return
	}

	var merchant model.Merchant
	if err := db.DB.Where("id = ?", currentMerchantID).Take(&merchant).Error; err != nil {
		result.ErrSetMsg(c, "商家不存在")
		return
	}

	var product model.Product
	if err := db.DB.Where("id = ? AND merchant_id = ?", req.ProductID, currentMerchantID).Take(&product).Error; err != nil {
		result.ErrSetMsg(c, "商品不存在或不属于当前商家")
		return
	}

	now := util.GetTime()
	shareLink := model.ShareLink{
		ID:          util.GetUuid(),
		MerchantID:  currentMerchantID,
		ProductID:   req.ProductID,
		ShareCode:   util.GetUuid()[:10],
		Title:       fallbackString(strings.TrimSpace(req.Title), product.Title),
		Description: fallbackString(strings.TrimSpace(req.Description), product.Description),
		Status:      model.ShareStatusActive,
		ExpiresAt:   req.ExpiresAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := db.DB.Create(&shareLink).Error; err != nil {
		result.ErrSetMsg(c, "创建分享链接失败")
		return
	}

	result.OkSetData(c, gin.H{
		"shareLink": shareLink,
		"shareUrl":  buildShareURL(c, shareLink.ShareCode),
	})
}

func GetShareLinkDetail(c *gin.Context) {
	shareCode := c.Param("code")
	if shareCode == "" {
		result.ErrSetMsg(c, "分享码不能为空")
		return
	}

	var shareLink model.ShareLink
	if err := db.DB.Where("share_code = ? AND status = ?", shareCode, model.ShareStatusActive).Take(&shareLink).Error; err != nil {
		result.ErrSetMsg(c, "分享链接不存在")
		return
	}
	if shareLink.ExpiresAt != nil && shareLink.ExpiresAt.Before(util.GetTime()) {
		result.ErrSetMsg(c, "分享链接已过期")
		return
	}

	var product model.Product
	var merchant model.Merchant
	var mediaList []model.MediaAsset
	db.DB.Where("id = ?", shareLink.ProductID).Take(&product)
	db.DB.Where("id = ?", shareLink.MerchantID).Take(&merchant)
	db.DB.Where("product_id = ?", shareLink.ProductID).Order("sort asc, created_at asc").Find(&mediaList)
	db.DB.Model(&model.ShareLink{}).Where("id = ?", shareLink.ID).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	shareLink.ViewCount++

	result.OkSetData(c, model.ShareView{
		ShareLink: shareLink,
		Product:   product,
		Merchant:  merchant,
		MediaList: mediaList,
		ShareURL:  buildShareURL(c, shareCode),
	})
}

func inferMediaType(contentType string, filename string) string {
	switch {
	case strings.HasPrefix(contentType, "video/"):
		return "video"
	case strings.HasPrefix(contentType, "image/"):
		return "image"
	}

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4", ".mov", ".avi", ".mkv", ".webm":
		return "video"
	default:
		return "image"
	}
}

func getUploadLocalDir() string {
	dir := strings.TrimSpace(viper.GetString("storage.local_dir"))
	if dir == "" {
		return "./uploads"
	}
	return dir
}

func buildPublicAssetURL(c *gin.Context, relativePath string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("storage.public_base_url")), "/")
	cleanPath := strings.ReplaceAll(relativePath, "\\", "/")
	if baseURL != "" {
		return baseURL + "/" + cleanPath
	}
	return fmt.Sprintf("%s/uploads/%s", requestBaseURL(c), cleanPath)
}

func buildShareURL(c *gin.Context, shareCode string) string {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("share.base_url")), "/")
	if baseURL == "" {
		baseURL = requestBaseURL(c)
	}
	return baseURL + "/share/" + shareCode
}

func requestBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func fallbackString(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

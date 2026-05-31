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

func CreateUser(c *gin.Context) {
	var req model.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "用户参数错误")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "用户名称不能为空")
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

	user := model.User{
		ID:           util.GetUuid(),
		Name:         strings.TrimSpace(req.Name),
		LoginName:    strings.TrimSpace(req.LoginName),
		PasswordHash: passwordHash,
		ContactName:  strings.TrimSpace(req.ContactName),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Status:       model.UserStatusActive,
		CreatedAt:    util.GetTime(),
		UpdatedAt:    util.GetTime(),
	}
	if err = db.DB.Create(&user).Error; err != nil {
		result.ErrSetMsg(c, "创建用户失败")
		return
	}
	var role model.Role
	err = db.DB.Where("code = ?", model.RoleCodeUser).Take(&role).Error
	if err == nil {
		_ = db.DB.Create(&model.UserRole{
			ID:        util.GetUuid(),
			UserID:    user.ID,
			RoleID:    role.ID,
			CreatedAt: util.GetTime(),
		}).Error
	}
	result.OkSetData(c, user)
}

func GetUserList(c *gin.Context) {
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var userList []model.User
	if err := db.DB.Where("id = ?", currentUser.ID).Order("created_at desc").Find(&userList).Error; err != nil {
		result.ErrSetMsg(c, "查询用户失败")
		return
	}
	result.OkSetData(c, userList)
}

func CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分类参数错误")
		return
	}
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "分类名称不能为空")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.CategoryStatusDraft
	}

	category := model.Category{
		ID:          util.GetUuid(),
		UserID:      currentUser.ID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		CoverURL:    strings.TrimSpace(req.CoverURL),
		Status:      status,
		CreatedAt:   util.GetTime(),
		UpdatedAt:   util.GetTime(),
	}
	if err := db.DB.Create(&category).Error; err != nil {
		result.ErrSetMsg(c, "创建分类失败")
		return
	}
	result.OkSetData(c, category)
}

func GetUserCategoryList(c *gin.Context) {
	userID := c.Param("userId")
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if userID == "" || userID != currentUser.ID {
		result.ErrSetMsg(c, "无权查看其他用户的分类")
		return
	}

	var categoryList []model.Category
	if err := db.DB.Where("user_id = ?", userID).Order("created_at desc").Find(&categoryList).Error; err != nil {
		result.ErrSetMsg(c, "查询分类失败")
		return
	}
	result.OkSetData(c, categoryList)
}

func GetCategoryDetail(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}

	var user model.User
	var resourceList []model.CategoryResourceRelation
	var shareLinks []model.ShareLink
	db.DB.Where("id = ?", category.UserID).Take(&user)
	db.DB.Where("category_id = ?", categoryID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Where("category_id = ?", categoryID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, model.CategoryDetail{
		Category:     category,
		User:         user,
		ResourceList: resourceList,
		ShareLinks:   shareLinks,
	})
}

func UploadCategoryResource(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		result.ErrSetMsg(c, "请上传文件")
		return
	}

	resourceType := strings.TrimSpace(c.PostForm("resourceType"))
	if resourceType == "" {
		resourceType = inferResourceType(file.Header.Get("Content-Type"), file.Filename)
	}

	sortValue, _ := strconv.Atoi(c.DefaultPostForm("sort", "0"))
	now := util.GetTime()
	ext := strings.ToLower(filepath.Ext(file.Filename))
	relativePath := filepath.Join(
		"categories",
		now.Format("2006"),
		now.Format("01"),
		categoryID+"_"+util.GetUuid()[:8]+ext,
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
	resource := model.Resource{
		ID:           util.GetUuid(),
		UserID:       category.UserID,
		ResourceType: resourceType,
		FileName:     file.Filename,
		FileExt:      ext,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		StoragePath:  strings.ReplaceAll(relativePath, "\\", "/"),
		URL:          publicURL,
		Status:       model.ResourceStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	relation := model.CategoryResourceRelation{
		ID:           util.GetUuid(),
		UserID:       category.UserID,
		CategoryID:   categoryID,
		ResourceID:   resource.ID,
		ResourceType: resourceType,
		FileName:     file.Filename,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		URL:          publicURL,
		Sort:         sortValue,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if errTx := tx.Create(&resource).Error; errTx != nil {
			return errTx
		}
		if errTx := tx.Create(&relation).Error; errTx != nil {
			return errTx
		}
		if category.CoverURL == "" && resourceType == "image" {
			if errTx := tx.Model(&model.Category{}).Where("id = ?", categoryID).Updates(map[string]interface{}{
				"cover_url":  publicURL,
				"updated_at": now,
			}).Error; errTx != nil {
				return errTx
			}
		}
		return nil
	}); err != nil {
		result.ErrSetMsg(c, "保存分类资源失败")
		return
	}

	result.OkSetData(c, relation)
}

func CreateShareLink(c *gin.Context) {
	var req model.CreateShareLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分享参数错误")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if req.CategoryID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}

	var user model.User
	if err := db.DB.Where("id = ?", currentUserID).Take(&user).Error; err != nil {
		result.ErrSetMsg(c, "用户不存在")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", req.CategoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在或不属于当前用户")
		return
	}

	now := util.GetTime()
	shareLink := model.ShareLink{
		ID:          util.GetUuid(),
		UserID:      currentUserID,
		CategoryID:  req.CategoryID,
		ShareCode:   util.GetUuid()[:10],
		Title:       fallbackString(strings.TrimSpace(req.Title), category.Name),
		Description: fallbackString(strings.TrimSpace(req.Description), category.Description),
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

	var category model.Category
	var user model.User
	var resourceList []model.CategoryResourceRelation
	db.DB.Where("id = ?", shareLink.CategoryID).Take(&category)
	db.DB.Where("id = ?", shareLink.UserID).Take(&user)
	db.DB.Where("category_id = ?", shareLink.CategoryID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Model(&model.ShareLink{}).Where("id = ?", shareLink.ID).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	_ = db.DB.Create(&model.ShareViewLog{
		ID:          util.GetUuid(),
		ShareLinkID: shareLink.ID,
		CategoryID:  shareLink.CategoryID,
		UserID:      shareLink.UserID,
		ViewerIP:    c.ClientIP(),
		UserAgent:   c.GetHeader("User-Agent"),
		Referer:     c.GetHeader("Referer"),
		CreatedAt:   util.GetTime(),
	}).Error
	shareLink.ViewCount++

	result.OkSetData(c, model.ShareView{
		ShareLink:    shareLink,
		Category:     category,
		User:         user,
		ResourceList: resourceList,
		ShareURL:     buildShareURL(c, shareCode),
	})
}

func inferResourceType(contentType string, filename string) string {
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

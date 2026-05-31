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
	if hasCategoryName(currentUser.ID, req.Name, "") {
		result.ErrSetMsg(c, "分类名称不能重复")
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

func CreateCategoryItem(c *gin.Context) {
	var req model.CreateCategoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分类项参数错误")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.CategoryID) == "" || strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "分类 ID 和分类项名称不能为空")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", req.CategoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}
	if hasCategoryItemName(currentUserID, req.CategoryID, req.Name, "") {
		result.ErrSetMsg(c, "分类项名称不能重复")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.CategoryItemStatusDraft
	}

	item := model.CategoryItem{
		ID:          util.GetUuid(),
		UserID:      currentUserID,
		CategoryID:  req.CategoryID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		CoverURL:    strings.TrimSpace(req.CoverURL),
		Status:      status,
		CreatedAt:   util.GetTime(),
		UpdatedAt:   util.GetTime(),
	}
	if err := db.DB.Create(&item).Error; err != nil {
		result.ErrSetMsg(c, "创建分类项失败")
		return
	}
	result.OkSetData(c, item)
}

func UpdateCategoryItem(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		result.ErrSetMsg(c, "分类项 ID 不能为空")
		return
	}

	var req model.UpdateCategoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分类项参数错误")
		return
	}

	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "分类项名称不能为空")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类项不存在")
		return
	}
	if hasCategoryItemName(currentUserID, item.CategoryID, req.Name, itemID) {
		result.ErrSetMsg(c, "分类项名称不能重复")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = item.Status
	}
	now := util.GetTime()
	if err := db.DB.Model(&model.CategoryItem{}).Where("id = ?", itemID).Updates(map[string]interface{}{
		"name":        strings.TrimSpace(req.Name),
		"description": strings.TrimSpace(req.Description),
		"status":      status,
		"updated_at":  now,
	}).Error; err != nil {
		result.ErrSetMsg(c, "更新分类项失败")
		return
	}

	item.Name = strings.TrimSpace(req.Name)
	item.Description = strings.TrimSpace(req.Description)
	item.Status = status
	item.UpdatedAt = now
	result.OkSetData(c, item)
}

func DeleteCategory(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Param("id"))
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

	resourcePaths, _ := collectCategoryResourcePaths(categoryID, "")
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		var itemIDs []string
		if err := tx.Model(&model.CategoryItem{}).Where("category_id = ? AND user_id = ?", categoryID, currentUserID).Pluck("id", &itemIDs).Error; err != nil {
			return err
		}
		if len(itemIDs) > 0 {
			if err := tx.Where("category_item_id IN ?", itemIDs).Delete(&model.ShareLink{}).Error; err != nil {
				return err
			}
			if err := tx.Where("category_item_id IN ?", itemIDs).Delete(&model.ShareViewLog{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("category_id = ?", categoryID).Delete(&model.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ?", categoryID).Delete(&model.ShareViewLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ?", categoryID).Delete(&model.CategoryResourceRelation{}).Error; err != nil {
			return err
		}
		if len(itemIDs) > 0 {
			if err := tx.Where("id IN ?", itemIDs).Delete(&model.CategoryItem{}).Error; err != nil {
				return err
			}
		}
		return tx.Where("id = ?", categoryID).Delete(&model.Category{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除分类失败")
		return
	}

	removeLocalFiles(resourcePaths)
	result.OkSetData(c, gin.H{"id": categoryID})
}

func DeleteCategoryItem(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		result.ErrSetMsg(c, "分类项 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类项不存在")
		return
	}

	resourcePaths, _ := collectCategoryResourcePaths("", itemID)
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_item_id = ?", itemID).Delete(&model.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_item_id = ?", itemID).Delete(&model.ShareViewLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_item_id = ?", itemID).Delete(&model.CategoryResourceRelation{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", itemID).Delete(&model.CategoryItem{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除分类项失败")
		return
	}

	removeLocalFiles(resourcePaths)
	result.OkSetData(c, gin.H{"id": itemID})
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
	var itemList []model.CategoryItem
	var resourceList []model.CategoryResourceRelation
	var shareLinks []model.ShareLink
	db.DB.Where("id = ?", category.UserID).Take(&user)
	db.DB.Where("category_id = ?", categoryID).Order("created_at desc").Find(&itemList)
	db.DB.Where("category_id = ? AND (category_item_id = '' OR category_item_id IS NULL)", categoryID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Where("category_id = ?", categoryID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, model.CategoryDetail{
		Category:      category,
		User:          user,
		CategoryItems: itemList,
		ResourceList:  resourceList,
		ShareLinks:    shareLinks,
	})
}

func GetCategoryItemDetail(c *gin.Context) {
	itemID := c.Param("id")
	if itemID == "" {
		result.ErrSetMsg(c, "分类项 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类项不存在")
		return
	}

	var resourceList []model.CategoryResourceRelation
	var shareLinks []model.ShareLink
	db.DB.Where("category_item_id = ?", itemID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Where("category_item_id = ?", itemID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, gin.H{
		"categoryItem": item,
		"resourceList": resourceList,
		"shareLinks":   shareLinks,
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

	categoryItemID := strings.TrimSpace(c.PostForm("categoryItemId"))
	if categoryItemID != "" {
		var item model.CategoryItem
		if err := db.DB.Where("id = ? AND category_id = ? AND user_id = ?", categoryItemID, categoryID, currentUserID).Take(&item).Error; err != nil {
			result.ErrSetMsg(c, "分类项不存在")
			return
		}
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
	pathID := categoryID
	if categoryItemID != "" {
		pathID = categoryItemID
	}
	relativePath := filepath.Join(
		"categories",
		now.Format("2006"),
		now.Format("01"),
		pathID+"_"+util.GetUuid()[:8]+ext,
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
		ID:             util.GetUuid(),
		UserID:         category.UserID,
		CategoryID:     categoryID,
		CategoryItemID: categoryItemID,
		ResourceID:     resource.ID,
		ResourceType:   resourceType,
		FileName:       file.Filename,
		FileSize:       file.Size,
		MimeType:       file.Header.Get("Content-Type"),
		URL:            publicURL,
		Sort:           sortValue,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if errTx := tx.Create(&resource).Error; errTx != nil {
			return errTx
		}
		if errTx := tx.Create(&relation).Error; errTx != nil {
			return errTx
		}
		if categoryItemID != "" {
			return tx.Model(&model.CategoryItem{}).Where("id = ? AND cover_url = ''", categoryItemID).Updates(map[string]interface{}{
				"cover_url":  publicURL,
				"updated_at": now,
			}).Error
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

func DeleteResource(c *gin.Context) {
	resourceID := strings.TrimSpace(c.Param("id"))
	if resourceID == "" {
		result.ErrSetMsg(c, "资源 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var resource model.Resource
	if err := db.DB.Where("id = ? AND user_id = ?", resourceID, currentUserID).Take(&resource).Error; err != nil {
		result.ErrSetMsg(c, "资源不存在")
		return
	}

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("resource_id = ?", resourceID).Delete(&model.CategoryResourceRelation{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ? AND user_id = ?", resourceID, currentUserID).Delete(&model.Resource{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除资源失败")
		return
	}

	removeLocalFiles([]string{resource.StoragePath})

	result.OkSetData(c, gin.H{"id": resourceID})
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

	targetType := strings.TrimSpace(req.TargetType)
	if targetType == "" {
		targetType = model.ShareTargetCategory
	}
	if targetType != model.ShareTargetCategory && targetType != model.ShareTargetItem {
		result.ErrSetMsg(c, "分享目标类型不正确")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", req.CategoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在或不属于当前用户")
		return
	}

	var item *model.CategoryItem
	if targetType == model.ShareTargetItem {
		if req.CategoryItemID == "" {
			result.ErrSetMsg(c, "分享分类项时必须传分类项 ID")
			return
		}
		var itemRecord model.CategoryItem
		if err := db.DB.Where("id = ? AND category_id = ? AND user_id = ?", req.CategoryItemID, req.CategoryID, currentUserID).Take(&itemRecord).Error; err != nil {
			result.ErrSetMsg(c, "分类项不存在或不属于当前分类")
			return
		}
		item = &itemRecord
	}

	now := util.GetTime()
	shareLink := model.ShareLink{
		ID:             util.GetUuid(),
		UserID:         currentUserID,
		CategoryID:     req.CategoryID,
		CategoryItemID: req.CategoryItemID,
		TargetType:     targetType,
		ShareCode:      util.GetUuid()[:10],
		Title:          shareTitle(req, category, item),
		Description:    shareDescription(req, category, item),
		Status:         model.ShareStatusActive,
		ExpiresAt:      req.ExpiresAt,
		CreatedAt:      now,
		UpdatedAt:      now,
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
	shareView, ok := loadShareView(c)
	if !ok {
		return
	}
	result.OkSetData(c, shareView)
}

func RenderSharePage(c *gin.Context) {
	shareCode := strings.TrimSpace(c.Param("code"))
	if shareCode == "" {
		result.ErrSetMsg(c, "分享码不能为空")
		return
	}
	c.Redirect(302, buildShareURL(c, shareCode))
}

func loadShareView(c *gin.Context) (model.ShareView, bool) {
	shareCode := c.Param("code")
	if shareCode == "" {
		result.ErrSetMsg(c, "分享码不能为空")
		return model.ShareView{}, false
	}

	var shareLink model.ShareLink
	if err := db.DB.Where("share_code = ? AND status = ?", shareCode, model.ShareStatusActive).Take(&shareLink).Error; err != nil {
		result.ErrSetMsg(c, "分享链接不存在")
		return model.ShareView{}, false
	}
	if shareLink.ExpiresAt != nil && shareLink.ExpiresAt.Before(util.GetTime()) {
		result.ErrSetMsg(c, "分享链接已过期")
		return model.ShareView{}, false
	}

	var category model.Category
	var user model.User
	var item *model.CategoryItem
	var itemList []model.CategoryItem
	var resourceList []model.CategoryResourceRelation
	db.DB.Where("id = ?", shareLink.CategoryID).Take(&category)
	db.DB.Where("id = ?", shareLink.UserID).Take(&user)
	if shareLink.TargetType == model.ShareTargetItem && shareLink.CategoryItemID != "" {
		var itemRecord model.CategoryItem
		if err := db.DB.Where("id = ?", shareLink.CategoryItemID).Take(&itemRecord).Error; err == nil {
			item = &itemRecord
		}
		db.DB.Where("category_item_id = ?", shareLink.CategoryItemID).Order("sort asc, created_at asc").Find(&resourceList)
	} else {
		db.DB.Where("category_id = ?", shareLink.CategoryID).Order("created_at asc").Find(&itemList)
		db.DB.Where("category_id = ?", shareLink.CategoryID).Order("sort asc, created_at asc").Find(&resourceList)
	}
	db.DB.Model(&model.ShareLink{}).Where("id = ?", shareLink.ID).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	_ = db.DB.Create(&model.ShareViewLog{
		ID:             util.GetUuid(),
		ShareLinkID:    shareLink.ID,
		CategoryID:     shareLink.CategoryID,
		CategoryItemID: shareLink.CategoryItemID,
		TargetType:     shareLink.TargetType,
		UserID:         shareLink.UserID,
		ViewerIP:       c.ClientIP(),
		UserAgent:      c.GetHeader("User-Agent"),
		Referer:        c.GetHeader("Referer"),
		CreatedAt:      util.GetTime(),
	}).Error
	shareLink.ViewCount++

	return model.ShareView{
		ShareLink:     shareLink,
		Category:      category,
		CategoryItem:  item,
		CategoryItems: itemList,
		User:          user,
		ResourceList:  resourceList,
		ShareURL:      buildShareURL(c, shareCode),
	}, true
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
		baseURL = strings.TrimRight(strings.TrimSpace(c.GetHeader("Origin")), "/")
	}
	if baseURL == "" {
		baseURL = strings.TrimRight(strings.TrimSpace(c.GetHeader("Referer")), "/")
		if index := strings.Index(baseURL, "/#/"); index >= 0 {
			baseURL = baseURL[:index]
		}
		if index := strings.Index(baseURL, "/share/"); index >= 0 {
			baseURL = baseURL[:index]
		}
	}
	if baseURL == "" {
		baseURL = requestBaseURL(c)
	}
	return baseURL + "/#/share/" + shareCode
}

func requestBaseURL(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s", scheme, c.Request.Host)
}

func shareTitle(req model.CreateShareLinkRequest, category model.Category, item *model.CategoryItem) string {
	if strings.TrimSpace(req.Title) != "" {
		return strings.TrimSpace(req.Title)
	}
	if item != nil {
		return item.Name
	}
	return category.Name
}

func shareDescription(req model.CreateShareLinkRequest, category model.Category, item *model.CategoryItem) string {
	if strings.TrimSpace(req.Description) != "" {
		return strings.TrimSpace(req.Description)
	}
	if item != nil && strings.TrimSpace(item.Description) != "" {
		return strings.TrimSpace(item.Description)
	}
	return strings.TrimSpace(category.Description)
}

func fallbackString(value string, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

func hasCategoryName(userID string, name string, excludeID string) bool {
	query := db.DB.Model(&model.Category{}).Where("user_id = ? AND name = ?", userID, strings.TrimSpace(name))
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	_ = query.Count(&count).Error
	return count > 0
}

func hasCategoryItemName(userID string, categoryID string, name string, excludeID string) bool {
	query := db.DB.Model(&model.CategoryItem{}).Where(
		"user_id = ? AND category_id = ? AND name = ?",
		userID,
		categoryID,
		strings.TrimSpace(name),
	)
	if excludeID != "" {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	_ = query.Count(&count).Error
	return count > 0
}

func collectCategoryResourcePaths(categoryID string, categoryItemID string) ([]string, error) {
	query := db.DB.Model(&model.Resource{})
	if categoryItemID != "" {
		query = query.Joins(
			"JOIN tbl_category_resource_relation rel ON rel.resource_id = tbl_resource.id",
		).Where("rel.category_item_id = ?", categoryItemID)
	} else {
		query = query.Joins(
			"JOIN tbl_category_resource_relation rel ON rel.resource_id = tbl_resource.id",
		).Where("rel.category_id = ?", categoryID)
	}
	var paths []string
	return paths, query.Pluck("tbl_resource.storage_path", &paths).Error
}

func removeLocalFiles(paths []string) {
	for _, path := range paths {
		cleanPath := strings.TrimSpace(path)
		if cleanPath == "" {
			continue
		}
		fullPath := filepath.Join(getUploadLocalDir(), filepath.FromSlash(cleanPath))
		if _, err := os.Stat(fullPath); err == nil {
			_ = os.Remove(fullPath)
		}
	}
}

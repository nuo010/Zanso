package service

import (
	"crypto/rand"
	"fmt"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
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

	now := util.GetTime()
	user := model.User{
		ID:           util.GetUuid(),
		Name:         strings.TrimSpace(req.Name),
		LoginName:    strings.TrimSpace(req.LoginName),
		PasswordHash: passwordHash,
		ContactName:  strings.TrimSpace(req.ContactName),
		ContactPhone: strings.TrimSpace(req.ContactPhone),
		Status:       model.UserStatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = db.DB.Transaction(func(tx *gorm.DB) error {
		if errTx := tx.Create(&user).Error; errTx != nil {
			return errTx
		}
		role, errTx := ensureDefaultRole(tx, model.RoleCodeUser, "普通用户")
		if errTx != nil {
			return errTx
		}
		return tx.Create(&model.UserRole{
			ID:        util.GetUuid(),
			UserID:    user.ID,
			RoleID:    role.ID,
			CreatedAt: now,
		}).Error
	}); err != nil {
		result.ErrSetMsg(c, "创建用户失败")
		return
	}
	userWithRoles, err := buildUserWithRoles(user)
	if err != nil {
		result.OkSetData(c, user)
		return
	}
	result.OkSetData(c, userWithRoles)
}

func GetUserList(c *gin.Context) {
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	if !isAdminUser(currentUser.ID) {
		result.ErrSetMsg(c, "只有管理员可以管理用户")
		return
	}

	var userList []model.User
	if err := db.DB.Order("created_at desc").Find(&userList).Error; err != nil {
		result.ErrSetMsg(c, "查询用户失败")
		return
	}
	userRoleList, err := buildUsersWithRoles(userList)
	if err != nil {
		result.ErrSetMsg(c, "查询用户角色失败")
		return
	}
	result.OkSetData(c, userRoleList)
}

func UpdateUserRole(c *gin.Context) {
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if !isAdminUser(currentUser.ID) {
		result.ErrSetMsg(c, "只有管理员可以修改用户角色")
		return
	}

	userID := strings.TrimSpace(c.Param("id"))
	if userID == "" {
		result.ErrSetMsg(c, "用户 ID 不能为空")
		return
	}
	if userID == currentUser.ID {
		result.ErrSetMsg(c, "不能修改自己的角色，别把自己锁门外了")
		return
	}

	var req model.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "角色参数错误")
		return
	}
	roleCode := strings.TrimSpace(req.RoleCode)
	if roleCode != model.RoleCodeAdmin && roleCode != model.RoleCodeUser {
		result.ErrSetMsg(c, "角色只能是 admin 或 user")
		return
	}

	var user model.User
	if err := db.DB.Where("id = ?", userID).Take(&user).Error; err != nil {
		result.ErrSetMsg(c, "用户不存在")
		return
	}

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		roleName := "普通用户"
		if roleCode == model.RoleCodeAdmin {
			roleName = "管理员"
		}
		role, errTx := ensureDefaultRole(tx, roleCode, roleName)
		if errTx != nil {
			return errTx
		}
		if errTx := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; errTx != nil {
			return errTx
		}
		return tx.Create(&model.UserRole{
			ID:        util.GetUuid(),
			UserID:    userID,
			RoleID:    role.ID,
			CreatedAt: util.GetTime(),
		}).Error
	}); err != nil {
		result.ErrSetMsg(c, "修改用户角色失败")
		return
	}

	userWithRoles, err := buildUserWithRoles(user)
	if err != nil {
		result.ErrSetMsg(c, "查询用户角色失败")
		return
	}
	result.OkSetData(c, userWithRoles)
}

func ensureDefaultRole(tx *gorm.DB, code string, name string) (model.Role, error) {
	var role model.Role
	if err := tx.Where("code = ?", code).Take(&role).Error; err == nil {
		return role, nil
	}
	now := util.GetTime()
	role = model.Role{
		ID:          util.GetUuid(),
		Name:        name,
		Code:        code,
		Description: name,
		Status:      model.UserStatusActive,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	return role, tx.Create(&role).Error
}

func isAdminUser(userID string) bool {
	var count int64
	db.DB.Table("tbl_user_role AS ur").
		Joins("JOIN tbl_role AS r ON r.id = ur.role_id").
		Where("ur.user_id = ? AND r.code = ?", userID, model.RoleCodeAdmin).
		Count(&count)
	return count > 0
}

func buildUserWithRoles(user model.User) (model.UserWithRoles, error) {
	users, err := buildUsersWithRoles([]model.User{user})
	if err != nil {
		return model.UserWithRoles{}, err
	}
	if len(users) == 0 {
		return model.UserWithRoles{}, nil
	}
	return users[0], nil
}

func buildUsersWithRoles(users []model.User) ([]model.UserWithRoles, error) {
	userIDs := make([]string, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	type roleRow struct {
		UserID string
		Code   string
		Name   string
	}
	var roleRows []roleRow
	if len(userIDs) > 0 {
		if err := db.DB.Table("tbl_user_role AS ur").
			Select("ur.user_id, r.code, r.name").
			Joins("JOIN tbl_role AS r ON r.id = ur.role_id").
			Where("ur.user_id IN ?", userIDs).
			Scan(&roleRows).Error; err != nil {
			return nil, err
		}
	}

	roleCodeMap := make(map[string][]string, len(users))
	roleNameMap := make(map[string][]string, len(users))
	for _, row := range roleRows {
		roleCodeMap[row.UserID] = append(roleCodeMap[row.UserID], row.Code)
		roleNameMap[row.UserID] = append(roleNameMap[row.UserID], row.Name)
	}

	resultList := make([]model.UserWithRoles, 0, len(users))
	for _, user := range users {
		resultList = append(resultList, model.UserWithRoles{
			ID:           user.ID,
			Name:         user.Name,
			LoginName:    user.LoginName,
			ContactName:  user.ContactName,
			ContactPhone: user.ContactPhone,
			Status:       user.Status,
			RoleCodes:    roleCodeMap[user.ID],
			RoleNames:    roleNameMap[user.ID],
			CreatedAt:    user.CreatedAt,
			UpdatedAt:    user.UpdatedAt,
		})
	}
	return resultList, nil
}

func CreateCategory(c *gin.Context) {
	var req model.CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "展册参数错误")
		return
	}
	currentUser, ok := util.GetCurrentUser(c)
	if !ok {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "展册名称不能为空")
		return
	}
	if hasCategoryName(currentUser.ID, req.Name, "") {
		result.ErrSetMsg(c, "展册名称不能重复")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.CategoryStatusDraft
	}
	visible := true
	if req.Visible != nil {
		visible = *req.Visible
	}

	category := model.Category{
		ID:          util.GetUuid(),
		UserID:      currentUser.ID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		CoverURL:    strings.TrimSpace(req.CoverURL),
		Visible:     visible,
		Status:      status,
		CreatedAt:   util.GetTime(),
		UpdatedAt:   util.GetTime(),
	}
	if err := db.DB.Create(&category).Error; err != nil {
		result.ErrSetMsg(c, "创建展册失败")
		return
	}
	result.OkSetData(c, category)
}

func UpdateCategory(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Param("id"))
	if categoryID == "" {
		result.ErrSetMsg(c, "展册 ID 不能为空")
		return
	}
	var req model.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "展册参数错误")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "展册名称不能为空")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在")
		return
	}
	if hasCategoryName(currentUserID, req.Name, categoryID) {
		result.ErrSetMsg(c, "展册名称不能重复")
		return
	}

	visible := category.Visible
	if req.Visible != nil {
		visible = *req.Visible
	}
	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = category.Status
	}
	now := util.GetTime()
	if err := db.DB.Model(&model.Category{}).Where("id = ?", categoryID).Updates(map[string]interface{}{
		"name":        strings.TrimSpace(req.Name),
		"description": strings.TrimSpace(req.Description),
		"visible":     visible,
		"status":      status,
		"updated_at":  now,
	}).Error; err != nil {
		result.ErrSetMsg(c, "更新展册失败")
		return
	}

	category.Name = strings.TrimSpace(req.Name)
	category.Description = strings.TrimSpace(req.Description)
	category.Visible = visible
	category.Status = status
	category.UpdatedAt = now
	result.OkSetData(c, category)
}

func CreateCategoryItem(c *gin.Context) {
	var req model.CreateCategoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分类参数错误")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	collectionID := strings.TrimSpace(req.CategoryID)
	if collectionID == "" || strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "展册 ID 和分类名称不能为空")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", collectionID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在")
		return
	}
	if hasCategoryItemName(currentUserID, collectionID, req.Name, "") {
		result.ErrSetMsg(c, "分类名称不能重复")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = model.InnerCategoryStatusDraft
	}
	visible := true
	if req.Visible != nil {
		visible = *req.Visible
	}

	item := model.CategoryItem{
		ID:          util.GetUuid(),
		UserID:      currentUserID,
		CategoryID:  collectionID,
		Name:        strings.TrimSpace(req.Name),
		Description: strings.TrimSpace(req.Description),
		CoverURL:    strings.TrimSpace(req.CoverURL),
		Visible:     visible,
		Status:      status,
		CreatedAt:   util.GetTime(),
		UpdatedAt:   util.GetTime(),
	}
	if err := db.DB.Create(&item).Error; err != nil {
		result.ErrSetMsg(c, "创建分类失败")
		return
	}
	result.OkSetData(c, item)
}

func UpdateCategoryItem(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}

	var req model.UpdateCategoryItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "分类参数错误")
		return
	}

	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		result.ErrSetMsg(c, "分类名称不能为空")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}
	if hasCategoryItemName(currentUserID, item.CategoryID, req.Name, itemID) {
		result.ErrSetMsg(c, "分类名称不能重复")
		return
	}

	status := strings.TrimSpace(req.Status)
	if status == "" {
		status = item.Status
	}
	visible := item.Visible
	if req.Visible != nil {
		visible = *req.Visible
	}
	now := util.GetTime()
	if err := db.DB.Model(&model.CategoryItem{}).Where("id = ?", itemID).Updates(map[string]interface{}{
		"name":        strings.TrimSpace(req.Name),
		"description": strings.TrimSpace(req.Description),
		"visible":     visible,
		"status":      status,
		"updated_at":  now,
	}).Error; err != nil {
		result.ErrSetMsg(c, "更新分类失败")
		return
	}

	item.Name = strings.TrimSpace(req.Name)
	item.Description = strings.TrimSpace(req.Description)
	item.Visible = visible
	item.Status = status
	item.UpdatedAt = now
	result.OkSetData(c, item)
}

func DeleteCategory(c *gin.Context) {
	categoryID := strings.TrimSpace(c.Param("id"))
	if categoryID == "" {
		result.ErrSetMsg(c, "展册 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在")
		return
	}

	resourcePaths, _ := collectCategoryResourcePaths(categoryID, "")
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		var itemIDs []string
		if err := tx.Model(&model.CategoryItem{}).Where("collection_id = ? AND user_id = ?", categoryID, currentUserID).Pluck("id", &itemIDs).Error; err != nil {
			return err
		}
		if len(itemIDs) > 0 {
			if err := tx.Where("category_id IN ?", itemIDs).Delete(&model.ShareLink{}).Error; err != nil {
				return err
			}
			if err := tx.Where("category_id IN ?", itemIDs).Delete(&model.ShareViewLog{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("collection_id = ?", categoryID).Delete(&model.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("collection_id = ?", categoryID).Delete(&model.ShareViewLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("collection_id = ?", categoryID).Delete(&model.CategoryResourceRelation{}).Error; err != nil {
			return err
		}
		if len(itemIDs) > 0 {
			if err := tx.Where("id IN ?", itemIDs).Delete(&model.CategoryItem{}).Error; err != nil {
				return err
			}
		}
		return tx.Where("id = ?", categoryID).Delete(&model.Category{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除展册失败")
		return
	}

	removeResourceFiles(c, resourcePaths)
	result.OkSetData(c, gin.H{"id": categoryID})
}

func DeleteCategoryItem(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}

	resourcePaths, _ := collectCategoryResourcePaths("", itemID)
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("category_id = ?", itemID).Delete(&model.ShareLink{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ?", itemID).Delete(&model.ShareViewLog{}).Error; err != nil {
			return err
		}
		if err := tx.Where("category_id = ?", itemID).Delete(&model.CategoryResourceRelation{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", itemID).Delete(&model.CategoryItem{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除分类失败")
		return
	}

	removeResourceFiles(c, resourcePaths)
	result.OkSetData(c, gin.H{"id": itemID})
}

func GetCurrentUserCategoryList(c *gin.Context) {
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	page, pageSize := parsePageParams(c)
	var categoryList []model.Category
	query := db.DB.Model(&model.Category{}).Where("user_id = ?", currentUserID)
	var total int64
	if err := query.Count(&total).Error; err != nil {
		result.ErrSetMsg(c, "查询展册失败")
		return
	}
	if err := query.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&categoryList).Error; err != nil {
		result.ErrSetMsg(c, "查询展册失败")
		return
	}

	result.OkSetData(c, model.PageResult{
		List:     categoryList,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

func GetDashboardStats(c *gin.Context) {
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	stats := model.DashboardStats{}
	if err := db.DB.Model(&model.Category{}).Where("user_id = ?", currentUserID).Count(&stats.CollectionCount).Error; err != nil {
		result.ErrSetMsg(c, "查询首页统计失败")
		return
	}
	if err := db.DB.Model(&model.Resource{}).Where("user_id = ?", currentUserID).Count(&stats.ResourceCount).Error; err != nil {
		result.ErrSetMsg(c, "查询首页统计失败")
		return
	}
	if err := db.DB.Model(&model.Resource{}).
		Where("user_id = ?", currentUserID).
		Select("COALESCE(SUM(file_size), 0)").
		Scan(&stats.FileSizeTotal).Error; err != nil {
		result.ErrSetMsg(c, "查询首页统计失败")
		return
	}

	result.OkSetData(c, stats)
}

func GetCategoryDetail(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		result.ErrSetMsg(c, "展册 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在")
		return
	}

	page, pageSize := parsePageParams(c)
	var user model.User
	var itemList []model.CategoryItem
	var resourceList []model.CategoryResourceRelation
	var shareLinks []model.ShareLink
	db.DB.Where("id = ?", category.UserID).Take(&user)
	itemQuery := db.DB.Model(&model.CategoryItem{}).Where("collection_id = ?", categoryID)
	var itemTotal int64
	if err := itemQuery.Count(&itemTotal).Error; err != nil {
		result.ErrSetMsg(c, "查询展册详情失败")
		return
	}
	itemQuery.Order("created_at desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&itemList)
	db.DB.Where("collection_id = ? AND (category_id = '' OR category_id IS NULL)", categoryID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Where("collection_id = ?", categoryID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, model.CategoryDetail{
		Collection:   category,
		User:         user,
		Categories:   itemList,
		Total:        itemTotal,
		Page:         page,
		PageSize:     pageSize,
		ResourceList: resourceList,
		ShareLinks:   shareLinks,
	})
}

func GetCategoryItemDetail(c *gin.Context) {
	itemID := c.Param("id")
	if itemID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}

	var resourceList []model.CategoryResourceRelation
	var shareLinks []model.ShareLink
	db.DB.Where("category_id = ?", itemID).Order("sort asc, created_at asc").Find(&resourceList)
	db.DB.Where("category_id = ?", itemID).Order("created_at desc").Find(&shareLinks)

	result.OkSetData(c, gin.H{
		"category":     item,
		"resourceList": resourceList,
		"shareLinks":   shareLinks,
	})
}

func UpdateCategoryResourceSort(c *gin.Context) {
	itemID := strings.TrimSpace(c.Param("id"))
	if itemID == "" {
		result.ErrSetMsg(c, "分类 ID 不能为空")
		return
	}

	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var req model.UpdateCategoryResourceSortRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		result.ErrSetMsg(c, "排序参数错误")
		return
	}
	if len(req.ResourceRelationIDs) == 0 {
		result.ErrSetMsg(c, "资源顺序不能为空")
		return
	}

	var item model.CategoryItem
	if err := db.DB.Where("id = ? AND user_id = ?", itemID, currentUserID).Take(&item).Error; err != nil {
		result.ErrSetMsg(c, "分类不存在")
		return
	}

	var relations []model.CategoryResourceRelation
	if err := db.DB.Where("category_id = ? AND user_id = ?", itemID, currentUserID).Find(&relations).Error; err != nil {
		result.ErrSetMsg(c, "查询资源失败")
		return
	}
	if len(relations) != len(req.ResourceRelationIDs) {
		result.ErrSetMsg(c, "资源顺序数据不完整")
		return
	}

	relationMap := make(map[string]model.CategoryResourceRelation, len(relations))
	for _, relation := range relations {
		relationMap[relation.ID] = relation
	}
	for _, relationID := range req.ResourceRelationIDs {
		if _, ok := relationMap[relationID]; !ok {
			result.ErrSetMsg(c, "资源顺序数据无效")
			return
		}
	}

	now := util.GetTime()
	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		for index, relationID := range req.ResourceRelationIDs {
			if err := tx.Model(&model.CategoryResourceRelation{}).
				Where("id = ? AND category_id = ? AND user_id = ?", relationID, itemID, currentUserID).
				Updates(map[string]interface{}{
					"sort":       index + 1,
					"updated_at": now,
				}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		result.ErrSetMsg(c, "保存资源顺序失败")
		return
	}

	result.OkSetData(c, gin.H{
		"categoryId": itemID,
	})
}

func UploadCategoryResource(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		result.ErrSetMsg(c, "展册 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", categoryID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在")
		return
	}

	categoryItemID := strings.TrimSpace(c.PostForm("categoryId"))
	if categoryItemID != "" {
		var item model.CategoryItem
		if err := db.DB.Where("id = ? AND collection_id = ? AND user_id = ?", categoryItemID, categoryID, currentUserID).Take(&item).Error; err != nil {
			result.ErrSetMsg(c, "分类不存在")
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
	storedFileName := randomResourceFileName(10) + ext
	relativePath := filepath.Join(
		"categories",
		now.Format("2006"),
		now.Format("01"),
		storedFileName,
	)
	storagePath := strings.ReplaceAll(relativePath, "\\", "/")
	publicURL, err := saveUploadedResource(c, file, storagePath)
	if err != nil {
		result.ErrSetMsg(c, "保存资源文件失败")
		return
	}

	resource := model.Resource{
		ID:           util.GetUuid(),
		UserID:       category.UserID,
		ResourceType: resourceType,
		FileName:     storedFileName,
		FileExt:      ext,
		FileSize:     file.Size,
		MimeType:     file.Header.Get("Content-Type"),
		StoragePath:  storagePath,
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
		FileName:       storedFileName,
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
		result.ErrSetMsg(c, "保存资源失败")
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

	removeResourceFiles(c, []string{resource.StoragePath})

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
	collectionID := strings.TrimSpace(req.CollectionID)
	categoryID := strings.TrimSpace(req.CategoryID)
	if collectionID == "" {
		result.ErrSetMsg(c, "展册 ID 不能为空")
		return
	}

	targetType := strings.TrimSpace(req.TargetType)
	if targetType == "" {
		targetType = model.ShareTargetCollection
	}
	if targetType != model.ShareTargetCollection && targetType != model.ShareTargetCategory {
		result.ErrSetMsg(c, "分享目标类型不正确")
		return
	}

	var category model.Category
	if err := db.DB.Where("id = ? AND user_id = ?", collectionID, currentUserID).Take(&category).Error; err != nil {
		result.ErrSetMsg(c, "展册不存在或不属于当前用户")
		return
	}
	if !category.Visible {
		result.ErrSetMsg(c, "当前展册不可分享")
		return
	}

	var item *model.CategoryItem
	if targetType == model.ShareTargetCategory {
		if categoryID == "" {
			result.ErrSetMsg(c, "分享分类时必须传分类 ID")
			return
		}
		var itemRecord model.CategoryItem
		if err := db.DB.Where("id = ? AND collection_id = ? AND user_id = ?", categoryID, collectionID, currentUserID).Take(&itemRecord).Error; err != nil {
			result.ErrSetMsg(c, "分类不存在或不属于当前展册")
			return
		}
		if !itemRecord.Visible {
			result.ErrSetMsg(c, "当前分类不可分享")
			return
		}
		item = &itemRecord
	}

	now := util.GetTime()
	shareLink := model.ShareLink{
		ID:             util.GetUuid(),
		UserID:         currentUserID,
		CategoryID:     collectionID,
		CategoryItemID: categoryID,
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

func GetShareLinkList(c *gin.Context) {
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	categoryID := strings.TrimSpace(c.Query("collectionId"))
	categoryItemID := strings.TrimSpace(c.Query("categoryId"))

	query := db.DB.Table("tbl_share_link AS sl").
		Select(`
			sl.id, sl.share_code, sl.title, sl.description, sl.target_type, sl.collection_id,
			c.name AS category_name, sl.category_id, ci.name AS category_item_name,
			sl.view_count, sl.status, sl.expires_at, sl.created_at, sl.updated_at
		`).
		Joins("LEFT JOIN tbl_collection c ON c.id = sl.collection_id").
		Joins("LEFT JOIN tbl_category ci ON ci.id = sl.category_id").
		Where("sl.user_id = ?", currentUserID)

	if categoryID != "" {
		query = query.Where("sl.collection_id = ?", categoryID)
	}
	if categoryItemID != "" {
		query = query.Where("sl.category_id = ?", categoryItemID)
	}

	type shareLinkRow struct {
		ID               string
		ShareCode        string
		Title            string
		Description      string
		TargetType       string
		CategoryID       string
		CategoryName     string
		CategoryItemID   string
		CategoryItemName string
		ViewCount        int64
		Status           string
		ExpiresAt        *time.Time
		CreatedAt        time.Time
		UpdatedAt        time.Time
	}
	var rows []shareLinkRow
	if err := query.Order("sl.created_at desc").Scan(&rows).Error; err != nil {
		result.ErrSetMsg(c, "查询分享链接失败")
		return
	}

	list := make([]model.ShareLinkListItem, 0, len(rows))
	for _, row := range rows {
		list = append(list, model.ShareLinkListItem{
			ID:             row.ID,
			ShareCode:      row.ShareCode,
			Title:          row.Title,
			Description:    row.Description,
			TargetType:     row.TargetType,
			CollectionName: row.CategoryName,
			CollectionID:   row.CategoryID,
			CategoryID:     row.CategoryItemID,
			CategoryName:   row.CategoryItemName,
			ViewCount:      row.ViewCount,
			Status:         row.Status,
			ExpiresAt:      row.ExpiresAt,
			CreatedAt:      row.CreatedAt,
			UpdatedAt:      row.UpdatedAt,
			ShareURL:       buildShareURL(c, row.ShareCode),
		})
	}

	result.OkSetData(c, list)
}

func DeleteShareLink(c *gin.Context) {
	shareLinkID := strings.TrimSpace(c.Param("id"))
	if shareLinkID == "" {
		result.ErrSetMsg(c, "分享链接 ID 不能为空")
		return
	}
	currentUserID := util.CurrentUserID(c)
	if currentUserID == "" {
		result.ErrSetMsg(c, "登录状态无效")
		return
	}

	var shareLink model.ShareLink
	if err := db.DB.Where("id = ? AND user_id = ?", shareLinkID, currentUserID).Take(&shareLink).Error; err != nil {
		result.ErrSetMsg(c, "分享链接不存在")
		return
	}

	if err := db.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("share_link_id = ?", shareLinkID).Delete(&model.ShareViewLog{}).Error; err != nil {
			return err
		}
		return tx.Where("id = ?", shareLinkID).Delete(&model.ShareLink{}).Error
	}); err != nil {
		result.ErrSetMsg(c, "删除分享链接失败")
		return
	}

	result.OkSetData(c, gin.H{"id": shareLinkID})
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
	if !category.Visible {
		result.ErrSetMsg(c, "分享内容当前不可查看")
		return model.ShareView{}, false
	}
	if shareLink.TargetType == model.ShareTargetCategory && shareLink.CategoryItemID != "" {
		var itemRecord model.CategoryItem
		if err := db.DB.Where("id = ? AND collection_id = ?", shareLink.CategoryItemID, shareLink.CategoryID).Take(&itemRecord).Error; err == nil {
			item = &itemRecord
			if !itemRecord.Visible {
				result.ErrSetMsg(c, "分享内容当前不可查看")
				return model.ShareView{}, false
			}
		}
		db.DB.Where("category_id = ?", shareLink.CategoryItemID).Order("sort asc, created_at asc").Find(&resourceList)
	} else {
		db.DB.Where("collection_id = ? AND visible = ?", shareLink.CategoryID, true).Order("created_at asc").Find(&itemList)
		db.DB.Where("collection_id = ?", shareLink.CategoryID).Order("sort asc, created_at asc").Find(&resourceList)
	}

	shareLink.Title = currentShareTitle(shareLink, category, item)
	shareLink.Description = currentShareDescription(shareLink, category, item)
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
		ShareLink:    shareLink,
		Collection:   category,
		Category:     item,
		Categories:   itemList,
		User:         user,
		ResourceList: resourceList,
		ShareURL:     buildShareURL(c, shareCode),
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

func saveUploadedResource(c *gin.Context, file *multipart.FileHeader, storagePath string) (string, error) {
	if db.UseMinioStorage() {
		src, err := file.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()

		objectName, err := db.UploadMinioObject(c.Request.Context(), storagePath, src, file.Size, file.Header.Get("Content-Type"))
		if err != nil {
			return "", err
		}
		return db.BuildMinioStoragePath(objectName), nil
	}

	localDir := getUploadLocalDir()
	fullPath := filepath.Join(localDir, filepath.FromSlash(storagePath))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}
	if err := c.SaveUploadedFile(file, fullPath); err != nil {
		return "", err
	}
	return buildPublicAssetURL(c, storagePath), nil
}

func randomResourceFileName(length int) string {
	const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if length <= 0 {
		length = 10
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return util.GetUuid()[:length]
	}
	for index, value := range bytes {
		bytes[index] = alphabet[int(value)%len(alphabet)]
	}
	return string(bytes)
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

func currentShareTitle(shareLink model.ShareLink, category model.Category, item *model.CategoryItem) string {
	if shareLink.TargetType == model.ShareTargetCategory && item != nil && strings.TrimSpace(item.Name) != "" {
		return strings.TrimSpace(item.Name)
	}
	if strings.TrimSpace(category.Name) != "" {
		return strings.TrimSpace(category.Name)
	}
	return strings.TrimSpace(shareLink.Title)
}

func currentShareDescription(shareLink model.ShareLink, category model.Category, item *model.CategoryItem) string {
	if shareLink.TargetType == model.ShareTargetCategory && item != nil && strings.TrimSpace(item.Description) != "" {
		return strings.TrimSpace(item.Description)
	}
	if strings.TrimSpace(category.Description) != "" {
		return strings.TrimSpace(category.Description)
	}
	return strings.TrimSpace(shareLink.Description)
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
		"user_id = ? AND collection_id = ? AND name = ?",
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
			"JOIN tbl_collection_resource_relation rel ON rel.resource_id = tbl_resource.id",
		).Where("rel.category_id = ?", categoryItemID)
	} else {
		query = query.Joins(
			"JOIN tbl_collection_resource_relation rel ON rel.resource_id = tbl_resource.id",
		).Where("rel.collection_id = ?", categoryID)
	}
	var paths []string
	return paths, query.Pluck("tbl_resource.storage_path", &paths).Error
}

func removeResourceFiles(c *gin.Context, paths []string) {
	if db.UseMinioStorage() {
		for _, path := range paths {
			cleanPath := strings.TrimSpace(path)
			if cleanPath == "" {
				continue
			}
			if err := db.DeleteMinioObject(c.Request.Context(), cleanPath); err != nil {
				util.Log().Error("删除 MinIO 对象失败: %v", err)
			}
		}
		return
	}
	removeLocalFiles(paths)
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

func parsePageParams(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("page", "1")))
	pageSize, _ := strconv.Atoi(strings.TrimSpace(c.DefaultQuery("pageSize", "20")))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

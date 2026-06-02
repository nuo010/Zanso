package model

import "time"

const (
	UserStatusActive          = "active"
	RoleCodeAdmin             = "admin"
	RoleCodeUser              = "user"
	CategoryStatusDraft       = "draft"
	CategoryStatusActive      = "active"
	InnerCategoryStatusDraft  = "draft"
	InnerCategoryStatusActive = "active"
	ShareStatusActive         = "active"
	ResourceStatusActive      = "active"
	ShareTargetCollection     = "collection"
	ShareTargetCategory       = "category"
)

type User struct {
	ID           string    `json:"id" gorm:"primarykey;size:32"`
	Name         string    `json:"name" gorm:"size:120;not null"`
	LoginName    string    `json:"loginName" gorm:"size:64;uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	ContactName  string    `json:"contactName" gorm:"size:64"`
	ContactPhone string    `json:"contactPhone" gorm:"size:32"`
	Status       string    `json:"status" gorm:"size:32;not null;default:active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (User) TableName() string {
	return "tbl_user"
}

type Role struct {
	ID          string    `json:"id" gorm:"primarykey;size:32"`
	Name        string    `json:"name" gorm:"size:64;not null"`
	Code        string    `json:"code" gorm:"size:64;uniqueIndex;not null"`
	Description string    `json:"description" gorm:"size:255"`
	Status      string    `json:"status" gorm:"size:32;not null;default:active"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Role) TableName() string {
	return "tbl_role"
}

type UserRole struct {
	ID        string    `json:"id" gorm:"primarykey;size:32"`
	UserID    string    `json:"userId" gorm:"size:32;index;not null"`
	RoleID    string    `json:"roleId" gorm:"size:32;index;not null"`
	CreatedAt time.Time `json:"createdAt"`
}

func (UserRole) TableName() string {
	return "tbl_user_role"
}

type Category struct {
	ID          string    `json:"id" gorm:"primarykey;size:32"`
	UserID      string    `json:"userId" gorm:"size:32;index;not null"`
	Name        string    `json:"name" gorm:"size:160;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CoverURL    string    `json:"coverUrl" gorm:"size:512"`
	Visible     bool      `json:"visible" gorm:"not null;default:true"`
	Status      string    `json:"status" gorm:"size:32;not null;default:draft"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Category) TableName() string {
	return "tbl_collection"
}

type CategoryItem struct {
	ID          string    `json:"id" gorm:"primarykey;size:32"`
	UserID      string    `json:"userId" gorm:"size:32;index;not null"`
	CategoryID  string    `json:"categoryId" gorm:"column:collection_id;size:32;index;not null"`
	Name        string    `json:"name" gorm:"size:160;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CoverURL    string    `json:"coverUrl" gorm:"size:512"`
	Visible     bool      `json:"visible" gorm:"not null;default:true"`
	Status      string    `json:"status" gorm:"size:32;not null;default:draft"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (CategoryItem) TableName() string {
	return "tbl_category"
}

type Resource struct {
	ID           string    `json:"id" gorm:"primarykey;size:32"`
	UserID       string    `json:"userId" gorm:"size:32;index;not null"`
	ResourceType string    `json:"resourceType" gorm:"size:32;not null"`
	FileName     string    `json:"fileName" gorm:"size:255"`
	FileExt      string    `json:"fileExt" gorm:"size:32"`
	FileSize     int64     `json:"fileSize"`
	MimeType     string    `json:"mimeType" gorm:"size:128"`
	StoragePath  string    `json:"storagePath" gorm:"size:512;not null"`
	URL          string    `json:"url" gorm:"size:512;not null"`
	Status       string    `json:"status" gorm:"size:32;not null;default:active"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (Resource) TableName() string {
	return "tbl_resource"
}

type CategoryResourceRelation struct {
	ID             string    `json:"id" gorm:"primarykey;size:32"`
	UserID         string    `json:"userId" gorm:"size:32;index;not null"`
	CategoryID     string    `json:"collectionId" gorm:"column:collection_id;size:32;index;not null"`
	CategoryItemID string    `json:"categoryId" gorm:"column:category_id;size:32;index"`
	ResourceID     string    `json:"resourceId" gorm:"size:32;index;not null"`
	ResourceType   string    `json:"resourceType" gorm:"size:32;not null"`
	FileName       string    `json:"fileName" gorm:"size:255"`
	FileSize       int64     `json:"fileSize"`
	MimeType       string    `json:"mimeType" gorm:"size:128"`
	URL            string    `json:"url" gorm:"size:512;not null"`
	Sort           int       `json:"sort"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (CategoryResourceRelation) TableName() string {
	return "tbl_collection_resource_relation"
}

type ShareLink struct {
	ID             string     `json:"id" gorm:"primarykey;size:32"`
	UserID         string     `json:"userId" gorm:"size:32;index;not null"`
	CategoryID     string     `json:"collectionId" gorm:"column:collection_id;size:32;index;not null"`
	CategoryItemID string     `json:"categoryId" gorm:"column:category_id;size:32;index"`
	TargetType     string     `json:"targetType" gorm:"size:32;index;not null;default:collection"`
	ShareCode      string     `json:"shareCode" gorm:"size:32;uniqueIndex;not null"`
	Title          string     `json:"title" gorm:"size:160"`
	Description    string     `json:"description" gorm:"type:text"`
	Status         string     `json:"status" gorm:"size:32;not null;default:active"`
	ViewCount      int64      `json:"viewCount"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
}

func (ShareLink) TableName() string {
	return "tbl_share_link"
}

type ShareViewLog struct {
	ID             string    `json:"id" gorm:"primarykey;size:32"`
	ShareLinkID    string    `json:"shareLinkId" gorm:"size:32;index;not null"`
	CategoryID     string    `json:"collectionId" gorm:"column:collection_id;size:32;index;not null"`
	CategoryItemID string    `json:"categoryId" gorm:"column:category_id;size:32;index"`
	TargetType     string    `json:"targetType" gorm:"size:32;index;not null"`
	UserID         string    `json:"userId" gorm:"size:32;index;not null"`
	ViewerIP       string    `json:"viewerIp" gorm:"size:64"`
	UserAgent      string    `json:"userAgent" gorm:"size:512"`
	Referer        string    `json:"referer" gorm:"size:512"`
	CreatedAt      time.Time `json:"createdAt"`
}

func (ShareViewLog) TableName() string {
	return "tbl_share_view_log"
}

type CreateUserRequest struct {
	Name         string `json:"name"`
	LoginName    string `json:"loginName"`
	Password     string `json:"password"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
}

type UserLoginRequest struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

type UpdateUserRoleRequest struct {
	RoleCode string `json:"roleCode"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	Visible     *bool  `json:"visible"`
	Status      string `json:"status"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visible     *bool  `json:"visible"`
	Status      string `json:"status"`
}

type CreateCategoryItemRequest struct {
	CategoryID  string `json:"collectionId"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	Visible     *bool  `json:"visible"`
	Status      string `json:"status"`
}

type UpdateCategoryItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Visible     *bool  `json:"visible"`
	Status      string `json:"status"`
}

type CreateShareLinkRequest struct {
	CollectionID string     `json:"collectionId"`
	CategoryID   string     `json:"categoryId"`
	TargetType   string     `json:"targetType"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	ExpiresAt    *time.Time `json:"expiresAt"`
}

type CategoryDetail struct {
	Collection   Category                   `json:"collection"`
	User         User                       `json:"user"`
	Categories   []CategoryItem             `json:"categories"`
	Total        int64                      `json:"total"`
	Page         int                        `json:"page"`
	PageSize     int                        `json:"pageSize"`
	ResourceList []CategoryResourceRelation `json:"resourceList"`
	ShareLinks   []ShareLink                `json:"shareLinks,omitempty"`
}

type PageResult struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

type UserWithRoles struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	LoginName    string    `json:"loginName"`
	ContactName  string    `json:"contactName"`
	ContactPhone string    `json:"contactPhone"`
	Status       string    `json:"status"`
	RoleCodes    []string  `json:"roleCodes"`
	RoleNames    []string  `json:"roleNames"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ShareLinkListItem struct {
	ID             string     `json:"id"`
	ShareCode      string     `json:"shareCode"`
	Title          string     `json:"title"`
	Description    string     `json:"description"`
	TargetType     string     `json:"targetType"`
	CollectionID   string     `json:"collectionId"`
	CollectionName string     `json:"collectionName"`
	CategoryID     string     `json:"categoryId"`
	CategoryName   string     `json:"categoryName"`
	ViewCount      int64      `json:"viewCount"`
	Status         string     `json:"status"`
	ExpiresAt      *time.Time `json:"expiresAt"`
	CreatedAt      time.Time  `json:"createdAt"`
	UpdatedAt      time.Time  `json:"updatedAt"`
	ShareURL       string     `json:"shareUrl"`
}

type ShareView struct {
	ShareLink    ShareLink                  `json:"shareLink"`
	Collection   Category                   `json:"collection"`
	Category     *CategoryItem              `json:"category,omitempty"`
	Categories   []CategoryItem             `json:"categories,omitempty"`
	User         User                       `json:"user"`
	ResourceList []CategoryResourceRelation `json:"resourceList"`
	ShareURL     string                     `json:"shareUrl"`
}

type UserAuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

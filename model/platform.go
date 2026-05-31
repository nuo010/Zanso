package model

import "time"

const (
	UserStatusActive     = "active"
	RoleCodeUser         = "user"
	CategoryStatusDraft  = "draft"
	CategoryStatusActive = "active"
	ShareStatusActive    = "active"
	ResourceStatusActive = "active"
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
	Status      string    `json:"status" gorm:"size:32;not null;default:draft"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Category) TableName() string {
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
	ID           string    `json:"id" gorm:"primarykey;size:32"`
	UserID       string    `json:"userId" gorm:"size:32;index;not null"`
	CategoryID   string    `json:"categoryId" gorm:"size:32;index;not null"`
	ResourceID   string    `json:"resourceId" gorm:"size:32;index;not null"`
	ResourceType string    `json:"resourceType" gorm:"size:32;not null"`
	FileName     string    `json:"fileName" gorm:"size:255"`
	FileSize     int64     `json:"fileSize"`
	MimeType     string    `json:"mimeType" gorm:"size:128"`
	URL          string    `json:"url" gorm:"size:512;not null"`
	Sort         int       `json:"sort"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (CategoryResourceRelation) TableName() string {
	return "tbl_category_resource_relation"
}

type ShareLink struct {
	ID          string     `json:"id" gorm:"primarykey;size:32"`
	UserID      string     `json:"userId" gorm:"size:32;index;not null"`
	CategoryID  string     `json:"categoryId" gorm:"size:32;index;not null"`
	ShareCode   string     `json:"shareCode" gorm:"size:32;uniqueIndex;not null"`
	Title       string     `json:"title" gorm:"size:160"`
	Description string     `json:"description" gorm:"type:text"`
	Status      string     `json:"status" gorm:"size:32;not null;default:active"`
	ViewCount   int64      `json:"viewCount"`
	ExpiresAt   *time.Time `json:"expiresAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (ShareLink) TableName() string {
	return "tbl_share_link"
}

type ShareViewLog struct {
	ID          string    `json:"id" gorm:"primarykey;size:32"`
	ShareLinkID string    `json:"shareLinkId" gorm:"size:32;index;not null"`
	CategoryID  string    `json:"categoryId" gorm:"size:32;index;not null"`
	UserID      string    `json:"userId" gorm:"size:32;index;not null"`
	ViewerIP    string    `json:"viewerIp" gorm:"size:64"`
	UserAgent   string    `json:"userAgent" gorm:"size:512"`
	Referer     string    `json:"referer" gorm:"size:512"`
	CreatedAt   time.Time `json:"createdAt"`
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

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	Status      string `json:"status"`
}

type CreateShareLinkRequest struct {
	CategoryID  string     `json:"categoryId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ExpiresAt   *time.Time `json:"expiresAt"`
}

type CategoryDetail struct {
	Category     Category                   `json:"category"`
	User         User                       `json:"user"`
	ResourceList []CategoryResourceRelation `json:"resourceList"`
	ShareLinks   []ShareLink                `json:"shareLinks,omitempty"`
}

type ShareView struct {
	ShareLink    ShareLink                  `json:"shareLink"`
	Category     Category                   `json:"category"`
	User         User                       `json:"user"`
	ResourceList []CategoryResourceRelation `json:"resourceList"`
	ShareURL     string                     `json:"shareUrl"`
}

type UserAuthResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

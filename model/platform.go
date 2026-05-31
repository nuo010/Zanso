package model

import "time"

const (
	MerchantStatusActive = "active"
	ProductStatusDraft   = "draft"
	ProductStatusActive  = "active"
	ShareStatusActive    = "active"
)

type Merchant struct {
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

func (Merchant) TableName() string {
	return "merchant"
}

type Product struct {
	ID          string    `json:"id" gorm:"primarykey;size:32"`
	MerchantID  string    `json:"merchantId" gorm:"size:32;index;not null"`
	Title       string    `json:"title" gorm:"size:160;not null"`
	Description string    `json:"description" gorm:"type:text"`
	CoverURL    string    `json:"coverUrl" gorm:"size:512"`
	Status      string    `json:"status" gorm:"size:32;not null;default:draft"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func (Product) TableName() string {
	return "product"
}

type MediaAsset struct {
	ID         string    `json:"id" gorm:"primarykey;size:32"`
	MerchantID string    `json:"merchantId" gorm:"size:32;index;not null"`
	ProductID  string    `json:"productId" gorm:"size:32;index;not null"`
	MediaType  string    `json:"mediaType" gorm:"size:32;not null"`
	FileName   string    `json:"fileName" gorm:"size:255"`
	FileSize   int64     `json:"fileSize"`
	MimeType   string    `json:"mimeType" gorm:"size:128"`
	URL        string    `json:"url" gorm:"size:512;not null"`
	Sort       int       `json:"sort"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

func (MediaAsset) TableName() string {
	return "media_asset"
}

type ShareLink struct {
	ID          string     `json:"id" gorm:"primarykey;size:32"`
	MerchantID  string     `json:"merchantId" gorm:"size:32;index;not null"`
	ProductID   string     `json:"productId" gorm:"size:32;index;not null"`
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
	return "share_link"
}

type CreateMerchantRequest struct {
	Name         string `json:"name"`
	LoginName    string `json:"loginName"`
	Password     string `json:"password"`
	ContactName  string `json:"contactName"`
	ContactPhone string `json:"contactPhone"`
}

type MerchantLoginRequest struct {
	LoginName string `json:"loginName"`
	Password  string `json:"password"`
}

type CreateProductRequest struct {
	MerchantID  string `json:"merchantId"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	Status      string `json:"status"`
}

type CreateShareLinkRequest struct {
	MerchantID  string     `json:"merchantId"`
	ProductID   string     `json:"productId"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	ExpiresAt   *time.Time `json:"expiresAt"`
}

type ProductDetail struct {
	Product    Product      `json:"product"`
	Merchant   Merchant     `json:"merchant"`
	MediaList  []MediaAsset `json:"mediaList"`
	ShareLinks []ShareLink  `json:"shareLinks,omitempty"`
}

type ShareView struct {
	ShareLink ShareLink    `json:"shareLink"`
	Product   Product      `json:"product"`
	Merchant  Merchant     `json:"merchant"`
	MediaList []MediaAsset `json:"mediaList"`
	ShareURL  string       `json:"shareUrl"`
}

type MerchantAuthResponse struct {
	Merchant Merchant `json:"merchant"`
	Token    string   `json:"token"`
}

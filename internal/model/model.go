package model

import (
	"time"
	"gorm.io/gorm"
)

type Model struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"size:100;not null" json:"name"`
	Version       string         `gorm:"size:50;not null" json:"version"`
	FileURL       string         `gorm:"size:255;not null" json:"file_url"`
	FileSize      uint64         `gorm:"not null" json:"file_size"`
	FileType      string         `gorm:"size:20;not null" json:"file_type"`
	Description   string         `gorm:"type:text" json:"description"`
	Category      string         `gorm:"size:50;not null" json:"category"`
	Tags          string         `gorm:"size:255" json:"tags"`
	Status        int8           `gorm:"not null;default:1" json:"status"`
	DownloadCount uint32         `gorm:"not null;default:0" json:"download_count"`
	CreatedBy     uint64         `gorm:"not null" json:"created_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type ModelCreateRequest struct {
	Name        string `json:"name" binding:"required,max=100"`
	Version     string `json:"version" binding:"required,max=50"`
	FileURL     string `json:"file_url" binding:"required,url"`
	FileSize    uint64 `json:"file_size" binding:"required"`
	FileType    string `json:"file_type" binding:"required,max=20"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"required,max=50"`
	Tags        string `json:"tags" binding:"max=255"`
}

type ModelUpdateRequest struct {
	Name        string `json:"name" binding:"omitempty,max=100"`
	Version     string `json:"version" binding:"omitempty,max=50"`
	Description string `json:"description"`
	Category    string `json:"category" binding:"omitempty,max=50"`
	Tags        string `json:"tags" binding:"omitempty,max=255"`
	Status      *int8  `json:"status" binding:"omitempty,oneof=0 1"`
}

// ModelFileResponse 模型文件响应
type ModelFileResponse struct {
	FileURL    string `json:"file_url"`    // 文件访问URL
	FileName   string `json:"file_name"`   // 文件名
	FileSize   uint64 `json:"file_size"`   // 文件大小
	FileType   string `json:"file_type"`   // 文件类型
	ExpireTime string `json:"expire_time"` // URL过期时间
}

// UploadURLRequest 获取上传URL请求
type UploadURLRequest struct {
	FileName string `json:"file_name" binding:"required"` // 文件名
	FileType string `json:"file_type" binding:"required"` // 文件类型
}

// UploadURLResponse 上传URL响应
type UploadURLResponse struct {
	UploadURL  string `json:"upload_url"`  // 上传URL
	FileURL    string `json:"file_url"`    // 文件访问URL
	ExpireTime string `json:"expire_time"` // URL过期时间
} 
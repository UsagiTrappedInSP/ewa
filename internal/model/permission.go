package model

import (
	"time"
	"gorm.io/gorm"
)

type UserModelPermission struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	UserID          uint64         `gorm:"not null" json:"user_id"`
	ModelID         uint64         `gorm:"not null" json:"model_id"`
	PermissionLevel int8           `gorm:"not null;default:1" json:"permission_level"`
	IsOwner         int8           `gorm:"not null;default:0" json:"is_owner"`
	ExpireTime      *time.Time     `json:"expire_time"`
	GrantedBy       uint64         `gorm:"not null" json:"granted_by"`
	GrantReason     string         `gorm:"size:255" json:"grant_reason"`
	Status          int8           `gorm:"not null;default:1" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type PermissionGrantRequest struct {
	UserID          uint64     `json:"user_id" binding:"required"`
	ModelID         uint64     `json:"model_id" binding:"required"`
	PermissionLevel int8       `json:"permission_level" binding:"required,min=1,max=4"`
	IsOwner         *int8      `json:"is_owner" binding:"omitempty,oneof=0 1"`
	ExpireTime      *time.Time `json:"expire_time"`
	GrantReason     string     `json:"grant_reason" binding:"max=255"`
}

type PermissionUpdateRequest struct {
	PermissionLevel int8       `json:"permission_level" binding:"omitempty,min=1,max=4"`
	IsOwner         *int8      `json:"is_owner" binding:"omitempty,oneof=0 1"`
	ExpireTime      *time.Time `json:"expire_time"`
	Status          *int8      `json:"status" binding:"omitempty,oneof=0 1"`
} 
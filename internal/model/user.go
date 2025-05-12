package model

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	ID            uint64         `gorm:"primaryKey" json:"id"`
	Username      string         `gorm:"size:50;not null;unique" json:"username"`
	Password      string         `gorm:"size:100;not null" json:"-"`
	Email         string         `gorm:"size:100;not null;unique" json:"email"`
	Phone         string         `gorm:"size:20;not null;unique" json:"phone"`
	Nickname      string         `gorm:"size:50" json:"nickname"`
	Avatar        string         `gorm:"size:255" json:"avatar"`
	Status        int8           `gorm:"not null;default:1" json:"status"`
	LastLoginTime *time.Time     `json:"last_login_time"`
	LastLoginIP   string         `gorm:"size:50" json:"last_login_ip"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Password string `json:"password" binding:"required,min=6,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required,len=11"`
	Nickname string `json:"nickname"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Phone    string `json:"phone" binding:"omitempty,len=11"`
	Email    string `json:"email" binding:"omitempty,email"`
} 
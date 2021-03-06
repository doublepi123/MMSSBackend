package entity

import (
	"gorm.io/gorm"
)

// Level 0 Administrator
// Level 1 Common User
// 用户Model
type UserEntity struct {
	gorm.Model
	Username      string
	Password      string
	Name          string
	RoleID        int
	Position      string
	WorkID        string `gorm:"not null"`
	Tittle        string
	Sex           string
	BirthDay      string
	Phone         string
	Address       string
	StartWorkTime string
}

// 返回给前端的用户列表
type SimpleUser struct {
	Username string
	Name     string
	RoleID   string
	Position string
	WorkID   string
}

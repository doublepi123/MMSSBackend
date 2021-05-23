package entity

import (
	"gorm.io/gorm"
	"time"
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
	WorkID        string
	Tittle        string
	Sex           string
	BirthDay      time.Time
	Phone         string
	Address       string
	StartWorkTime time.Time
}

// 返回给前端的用户列表
type SimpleUser struct {
	Username  string
	Name      string
	UserLevel int
	Position  string
	WorkID    string
}

const LevelAdministrator = 0
const LevelCommonUser = 1

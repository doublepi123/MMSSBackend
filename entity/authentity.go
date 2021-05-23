package entity

import "gorm.io/gorm"

type AuthEntity struct {
	gorm.Model
	RoleID int
	Auth   string
}

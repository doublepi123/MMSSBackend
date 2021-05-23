package entity

import "gorm.io/gorm"

type RoleEntity struct {
	gorm.Model
	Name string
}

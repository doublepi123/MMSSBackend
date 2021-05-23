package dao

import (
	"MMSSBackend/entity"
	"fmt"
)

type RoleDao struct {
	*Dao
}

func (roledao RoleDao) Exsit(id int) bool {
	roledao.db.DB.AutoMigrate(&entity.RoleEntity{})
	var count int64
	err := roledao.db.DB.Model(&entity.RoleEntity{}).Find(id).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (roledao RoleDao) GetRoleID(username string) int {
	var user entity.UserEntity
	roledao.db.DB.Model(&entity.UserEntity{}).Where("username = ? ", username).Find(&user)
	return user.RoleID
}

func (roledao RoleDao) CheckAuth(id int, Auth string) bool {
	if !roledao.Exsit(id) {
		return false
	}
	var count int64
	err := roledao.db.DB.Model(&entity.RoleEntity{}).Where("roleid = ? AND auth = ?", id, Auth).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count != 0
}

func (roledao RoleDao) AddAuth(id int, Auth string) bool {
	if !roledao.Exsit(id) {
		return false
	}

	if !roledao.CheckAuth(id, Auth) {
		return false
	}
	roledao.db.DB.AutoMigrate(&entity.AuthEntity{})
	err := roledao.db.DB.Create(&entity.AuthEntity{
		RoleID: id,
		Auth:   Auth,
	}).Error
	if err != nil {
		return false
	}
	return true
}

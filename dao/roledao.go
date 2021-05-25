package dao

import (
	"MMSSBackend/entity"
	"fmt"
)

type RoleDao struct {
	*Dao
}

func (roledao RoleDao) DelAuth(id int, Auth string) bool {
	if !roledao.ExistRoleID(id) {

		return false
	}
	if !roledao.CheckAuth(id, Auth) {
		return false
	}
	err := roledao.db.DB.Model(&entity.AuthEntity{}).Where("role_id = ? AND auth = ?", id, Auth).Delete(&entity.AuthEntity{}).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (roledao RoleDao) GetRoleName(id int) string {
	var role entity.RoleEntity
	roledao.db.DB.Where("id = ?", id).Find(&role)
	return role.Name
}

func (roledao RoleDao) ExistRoleID(id int) bool {
	roledao.db.DB.AutoMigrate(&entity.RoleEntity{})
	roledao.db.DB.AutoMigrate(&entity.AuthEntity{})
	var count int64
	roledao.db.DB.Model(&entity.RoleEntity{}).Find(id).Count(&count)
	return count != 0
}

func (RoleDao RoleDao) ExistRoleName(name string) bool {
	var count int64
	RoleDao.db.DB.Model(&entity.RoleEntity{}).Where("name = ?", name).Count(&count)
	return count != 0
}
func (roledao RoleDao) AddRole(name string) bool {
	roledao.db.DB.AutoMigrate(&entity.RoleEntity{})
	roledao.db.DB.AutoMigrate(&entity.AuthEntity{})
	roledao.db.DB.AutoMigrate(&entity.UserEntity{})
	if roledao.ExistRoleName(name) {
		return false
	}
	roledao.db.DB.Create(&entity.RoleEntity{
		Name: name,
	})
	return true
}

func (roledao RoleDao) GetRoleID(username string) int {
	var user entity.UserEntity
	roledao.db.DB.Model(&entity.UserEntity{}).Where("username = ? ", username).Find(&user)
	return user.RoleID
}

func (roledao RoleDao) CheckAuth(id int, Auth string) bool {
	if !roledao.ExistRoleID(id) {
		return false
	}
	var count int64
	err := roledao.db.DB.Model(&entity.AuthEntity{}).Where("role_id = ? AND auth = ?", id, Auth).Count(&count).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return count != 0
}

func (roledao RoleDao) AddAuth(id int, Auth string) bool {
	if !roledao.ExistRoleID(id) {
		return false
	}

	if roledao.CheckAuth(id, Auth) {
		return false
	}
	roledao.db.DB.AutoMigrate(&entity.AuthEntity{})
	err := roledao.db.DB.Create(&entity.AuthEntity{
		RoleID: id,
		Auth:   Auth,
	}).Error
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func (roledao RoleDao) Del(name string) bool {
	if !roledao.ExistRoleName(name) {
		return false
	}
	err := roledao.db.DB.Model(&entity.RoleEntity{}).Where("name = ?", name).Delete(&entity.RoleEntity{}).Error
	if err != nil {
		fmt.Println(err)
	}
	return err == nil
}

func (roledao RoleDao) RoleList() []entity.RoleEntity {
	var list []entity.RoleEntity
	roledao.db.DB.Model(&entity.RoleEntity{}).Find(&list)
	return list
}

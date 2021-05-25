package dao

import (
	"MMSSBackend/entity"
	"MMSSBackend/util"
	"errors"
)

type UserDao struct {
	*Dao
}

func (userdao UserDao) ChangePassword(username string, password string) error {
	return userdao.db.DB.Model(&entity.UserEntity{}).Where("username = ?", username).Update("password", util.GetPWD(password)).Error
}

func (userdao UserDao) Check(username string, password string) bool {
	var user entity.UserEntity
	userdao.db.DB.Where("Username = ?", username).First(&user)
	return util.CmpPWD(user.Password, password)
}
func (userdao UserDao) Exsit(username string) bool {
	var count int64
	userdao.db.DB.Model(&entity.UserEntity{}).Where("username = ?", username).Count(&count)
	return count != 0
}

func (userdao UserDao) Add(user *entity.UserEntity) error {
	userdao.db.DB.AutoMigrate(&entity.UserEntity{})
	userdao.db.DB.AutoMigrate(&entity.RoleEntity{})
	if userdao.Exsit(user.Username) {
		return errors.New("user exist")
	}
	return userdao.db.DB.Create(user).Error
}

func (userdao UserDao) Del(username string) error {
	userdao.db.DB.Model(&entity.AuthEntity{}).Delete("username = ?", username)
	return userdao.db.DB.Model(&entity.UserEntity{}).Delete("username = ?", username).Error
}

func (userdao UserDao) Find(username string) entity.UserEntity {
	var user entity.UserEntity
	userdao.db.DB.Where("username = ?", username).First(&user)
	return user
}

func (userdao UserDao) UserList() []entity.SimpleUser {
	var users []entity.SimpleUser
	userdao.db.DB.Model(&entity.UserEntity{}).Find(&users)
	return users
}

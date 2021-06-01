package dao

import (
	"MMSSBackend/entity"
	"errors"
	"fmt"
)

type PaperDao struct {
	*Dao
}

func (paperdao PaperDao) Add(paper *entity.PaperEntity) error {
	var count int64
	paperdao.db.DB.AutoMigrate(&entity.PaperEntity{})
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?",
		paper.UserName, paper.Tittle).Count(&count).Error
	if err != nil {
		return err
	}
	fmt.Println(count)
	if count >= 1 {
		return errors.New("paper exist")
	}
	err = paperdao.db.DB.Model(&entity.UserEntity{}).Where("username = ?", paper.UserName).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("user not exist")
	}
	return paperdao.db.DB.Create(&paper).Error

}

func (paperdao PaperDao) Find(username string, tittle string) (entity.PaperEntity, error) {
	var paper entity.PaperEntity
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?", username, tittle).Find(&paper).Error
	return paper, err
}

func (paperdao PaperDao) ADFind(q entity.PaperEntity) ([]entity.PaperEntity, error) {
	tx := paperdao.db.DB.Model(&entity.PaperEntity{})
	var paper []entity.PaperEntity
	err := tx.Where("user_name LIKE ? AND tittle LIKE ?", "%"+q.UserName+"%", "%"+q.Tittle+"%", "%").Find(&paper).Error
	return paper, err
}

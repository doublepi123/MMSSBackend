package dao

import (
	"MMSSBackend/entity"
	"errors"
	"fmt"
)

type PaperDao struct {
	*Dao
}

//添加paper
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

//普通查询
func (paperdao PaperDao) Find(username string, tittle string) (entity.PaperEntity, error) {
	var paper entity.PaperEntity
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?", username, tittle).Find(&paper).Error
	return paper, err
}

//模糊查询
func (paperdao PaperDao) ADFind(q entity.PaperEntity) ([]entity.PaperEntity, error) {
	tx := paperdao.db.DB.Model(&entity.PaperEntity{})
	var paper []entity.PaperEntity
	err := tx.Where("user_name LIKE ? AND tittle LIKE ?", "%"+q.UserName+"%", "%"+q.Tittle+"%", "%").Find(&paper).Error
	return paper, err
}

//为其他作者授权
func (paperdao PaperDao) Auth(WorkID string, PaperID uint, username string) error {
	var count int64
	err := paperdao.db.DB.Model(&entity.UserEntity{}).Where("work_id = ?", WorkID).Find(&entity.UserEntity{}).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("user not exsit")
	}
	err = paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ? AND user_name = ?", PaperID, username).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("paper not exsit")
	}
	paperdao.db.DB.AutoMigrate(&entity.PaperAuth{})
	err = paperdao.db.DB.Model(&entity.PaperAuth{}).Where("paper_id = ? AND work_id = ?", PaperID, WorkID).Count(&count).Error
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("paper auth exist")
	}
	err = paperdao.db.DB.Create(&entity.PaperAuth{
		PaperID: PaperID,
		WorkID:  WorkID,
	}).Error
	return err
}

func (paperdao PaperDao) AuthSelect(username string) ([]entity.PaperEntity, error) {
	var ans []entity.PaperEntity
	var user entity.UserEntity
	var auth []entity.PaperAuth
	err := paperdao.db.DB.Model(&entity.UserEntity{}).Where("username = ?", username).Find(&user).Error
	if err != nil {
		return ans, err
	}
	err = paperdao.db.DB.Model(&entity.PaperAuth{}).Where("work_id = ?", user.WorkID).Find(&auth).Error
	if err != nil {
		return ans, err
	}
	for i := range auth {
		var paper entity.PaperEntity
		err = paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ? ", auth[i].PaperID).Find(&paper).Error
		ans = append(ans, paper)
		if err != nil {
			return ans, err
		}
	}
	return ans, nil
}

func (paperdao PaperDao) AddFile(file *entity.PaperFile) error {
	tx := paperdao.db.DB.Begin()
	var count int64
	err := tx.Model(&entity.PaperEntity{}).Where("id = ?", file.PaperID).Count(&count).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	if count == 0 {
		tx.Rollback()
		return errors.New("paper not exist")
	}
	paperdao.db.DB.AutoMigrate(&entity.PaperFile{})
	tx.Model(&entity.PaperFile{}).Where("paper_id = ? ", file.PaperID).Delete(&entity.PaperFile{})
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Create(file).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Model(&entity.PaperEntity{}).Where("id = ?", file.PaperID).Update("filename", file.FileName).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit().Error

}

func (paperdao PaperDao) GetUncheckFile() ([]entity.PaperList, error) {
	var papers []entity.PaperList
	err := paperdao.db.DB.Debug().Model(&entity.PaperEntity{}).Where("hascheck = ?", false).Find(&papers).Error
	return papers, err
}

func (paperdao PaperDao) Check(paperid uint) error {
	var count int64
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ?", paperid).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("paper not exist")
	}
	return paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ?", paperid).Update("hascheck", true).Error
}

func (paperdao PaperDao) UnCheck(paperid uint) error {
	var count int64
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ?", paperid).Count(&count).Error
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("paper not exist")
	}
	return paperdao.db.DB.Model(&entity.PaperEntity{}).Where("id = ?", paperid).Update("hascheck", false).Error
}

func (paperdao PaperDao) Paperlist(username string) ([]entity.PaperList, error) {
	var paper []entity.PaperList
	err := paperdao.db.DB.Model(entity.PaperEntity{}).Where("user_name = ?", username).Find(&paper).Error
	return paper, err
}

func (paperdao PaperDao) Update(paper entity.PaperEntity) error {
	var old entity.PaperEntity
	err := paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?",
		paper.UserName, paper.Tittle).First(&old).Error
	if err != nil {
		return err
	}
	if old.Hascheck == true {
		return errors.New("canot be updated")
	}
	return paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?",
		paper.UserName, paper.Tittle).Updates(&paper).Error
}

func (paperdao PaperDao) GetPaper(paperid uint) (entity.PaperEntity, error) {
	var paper entity.PaperEntity
	err := paperdao.db.DB.Model(entity.PaperEntity{}).Find(&paper, paperid).Error
	return paper, err
}

func (paperdao PaperDao) GetFile(paperid uint) (entity.PaperFile, error) {
	var file entity.PaperFile
	err := paperdao.db.DB.Model(&entity.PaperFile{}).Where("paper_id = ?", paperid).Find(&file).Error
	return file, err
}

func (paperdao PaperDao) GetSomeoneAllPaper(username string) []entity.PaperList {
	var paperlist []entity.PaperList
	var ans []entity.PaperList
	var user entity.UserEntity
	var paperauth []entity.PaperAuth
	paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ?", username).Find(&paperlist)
	ans = append(ans, paperlist...)
	paperdao.db.DB.Model(&entity.UserEntity{}).Where("username = ?", username).Find(&user)
	paperdao.db.DB.Model(&entity.PaperAuth{}).Where("work_id = ?", user.WorkID).Find(&paperauth)
	for i := range paperauth {
		paperdao.db.DB.Model(&entity.PaperEntity{}).Find(&paperlist, paperauth[i].PaperID)
		ans = append(ans, paperlist...)
	}
	return ans
}

func (paperdao PaperDao) GetJournalList() []entity.Journal {
	var journallist []entity.Journal
	paperdao.db.DB.Model(&entity.Journal{}).Find(&journallist)
	return journallist
}

func (paperdao PaperDao) GetUncheckJournalList() []entity.Journal {
	var journallist []entity.Journal
	paperdao.db.DB.Model(&entity.Journal{}).Where("hascheck = ?", false).Find(&journallist)
	return journallist
}

func (paperdao PaperDao) CheckJournal(id uint) error {
	err := paperdao.db.DB.Model(&entity.Journal{}).Where("id = ?", id).Update("hascheck", true).Error
	if err != nil {
		return errors.New("id not find")
	}
	return err
}

func (paperdao PaperDao) DelJournal(id uint) error {
	return paperdao.db.DB.Model(&entity.Journal{}).Delete(&entity.Journal{}, id).Error
}

func (paperdao PaperDao) AddJournal(journal entity.Journal) error {
	journal.Hascheck = false
	paperdao.db.DB.AutoMigrate(&entity.Journal{})
	return paperdao.db.DB.Model(&entity.Journal{}).Create(&journal).Error
}

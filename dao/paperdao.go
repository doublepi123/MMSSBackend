package dao

import "MMSSBackend/entity"

type Paperdao struct {
	*Dao
}

func (paperdao Paperdao) AddPaper(paper *entity.PaperEntity) bool {
	paperdao.db.DB.AutoMigrate(&entity.PaperEntity{})
	paperdao.db.DB.AutoMigrate(&entity.PaperFile{})
	paperdao.db.DB.AutoMigrate(&entity.PaperAuth{})
	thispaper := entity.PaperEntity{}
	if paperdao.ExistPaper(paper) {
		return false
	}
	paperdao.db.DB.Create(paper)
	paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ？ ", paper.UserName, paper.Tittle).Find(&thispaper)
	paperdao.db.DB.Create(&entity.PaperAuth{
		PaperID:  thispaper.ID,
		Username: thispaper.UserName,
	})
	return true
}

func (paperdao Paperdao) ExistPaper(paper *entity.PaperEntity) bool {
	var count int64
	paperdao.db.DB.Model(&entity.PaperEntity{}).Where("user_name = ? AND tittle = ?", paper.UserName, paper.Tittle).Count(&count)

	return count == 1
}

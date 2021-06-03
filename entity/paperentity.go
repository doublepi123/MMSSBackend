package entity

import (
	"gorm.io/gorm"
	"time"
)

//论文Model
type PaperEntity struct {
	gorm.Model
	//用户名
	UserName string `gorm:"index"`
	//标题
	Tittle string `gorm:"index"`
	//发表日期
	Date string
	//期刊
	Journals string `gorm:"index"`
	//ISSN
	ISSN string
	//其他
	Other string
	//附件文件名
	Filename string
	//审核状态
	Hascheck bool
}

//论文列表
type PaperList struct {
	//PaperID
	ID uint
	//标题
	Tittle string
	//日期
	Date string
	//审核状态
	Hascheck bool
	//附件文件名
	Filename string
}

//论文附件
type PaperFile struct {
	CreatedAt time.Time
	//PaperID
	PaperID uint `gorm:"primarykey"`
	//二进制文件内容
	File []byte `gorm:"index"`
	//文件名
	FileName string `gorm:"index"`
}

//其他作者权限
type PaperAuth struct {
	gorm.Model
	PaperID uint
	WorkID  string
}

package entity

import (
	"gorm.io/gorm"
)

//论文Model
type PaperEntity struct {
	gorm.Model
	//用户名
	UserName string
	//标题
	Tittle string
	//发表日期
	Date string
	//期刊
	Journals string
	//ISSN
	ISSN string
	//学院
	College string
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
	//作者
	Author string
	//审核状态
	Hascheck bool
	//附件文件名
	Filename string
}

//论文附件
type PaperFile struct {
	gorm.Model
	//PaperID
	PaperID uint
	//二进制文件内容
	File []byte
	//文件名
	FileName string
}

//其他作者权限
type PaperAuth struct {
	gorm.Model
	PaperID uint
	WorkID  string
}

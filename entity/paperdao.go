package entity

import "gorm.io/gorm"

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
	//作者
	Author string
	//学院
	College string
	ISBN    string `gorm:"unique"`
	//其他
	Other string
}

//论文列表
type PaperList struct {
	//标题
	Tittle string
	//日期
	Date string
	//作者
	Author string
	//ISBN
	ISBN string
}

//论文附件
type PaperFile struct {
	gorm.Model
	//ISBN 不可重复 作为唯一标识
	ISBN string `gorm:"unique"`
	//二进制文件内容
	File []byte
	//文件名
	FileName string
}

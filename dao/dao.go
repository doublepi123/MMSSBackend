package dao

import "MMSSBackend/util"

type Dao struct {
	db util.Database
}

func (dao *Dao) Init() {
	dao.db.Init()
}

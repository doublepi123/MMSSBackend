package main

import (
	"MMSSBackend/dao"
	"MMSSBackend/service"
	"MMSSBackend/util"
	"time"
)

func main() {
	time.Sleep(time.Second * 5)
	basedao := &dao.Dao{}
	basedao.Init()
	userdao := &dao.UserDao{basedao}
	cookiedao := &dao.CookieDao{basedao}
	roledao := &dao.RoleDao{basedao}
	paperdao := &dao.PaperDao{basedao}
	server := service.Server{userdao, cookiedao, roledao, paperdao}
	go server.Run()
	util.PauseForRun()
}

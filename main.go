package main

import (
	"MMSSBackend/dao"
	"MMSSBackend/service"
	"MMSSBackend/util"
)

func main() {
	basedao := &dao.Dao{}
	basedao.Init()
	userdao := &dao.UserDao{basedao}
	cookiedao := &dao.CookieDao{basedao}
	roledao := &dao.RoleDao{basedao}
	server := service.Server{userdao, cookiedao, roledao}
	go server.Run()
	util.PauseForRun()
}

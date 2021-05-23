package service

import (
	"MMSSBackend/dao"
	"MMSSBackend/entity"
	"MMSSBackend/message"
	"MMSSBackend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Server struct {
	Userdao   *dao.UserDao
	Cookiedao *dao.CookieDao
	RoleDao   *dao.RoleDao
}

func (server Server) login(c *gin.Context) {
	var user entity.UserEntity
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, message.Fail())
		return
	}
	if server.Userdao.Check(user.Username, user.Password) {
		server.Cookiedao.SetCookie(user.Username, server.Cookiedao.GenerateUserID(user.Username), c)
		c.JSON(http.StatusOK, message.Success())
		return
	}
	c.JSON(http.StatusForbidden, message.Fail())
}

func (server Server) CheckLogin(c *gin.Context) {
	username, _ := c.Cookie("username")
	userid, _ := c.Cookie("userid")
	if !server.Cookiedao.CheckCookie(username, userid) {
		c.JSON(http.StatusForbidden, gin.H{"msg": "not login"})
		c.Abort()
		return
	}
	c.Next()
}

func (server Server) userAdmin(c *gin.Context) {
	username, _ := c.Cookie("username")
	roleid := server.RoleDao.GetRoleID(username)
	if !server.RoleDao.CheckAuth(roleid, "userManage") {
		fmt.Println(roleid)
		c.JSON(http.StatusForbidden, message.Fail())
		c.Abort()
		return
	}
	c.Next()
}

func (server Server) Run() {
	server.RoleDao.AddRole("Administrator")
	server.RoleDao.AddAuth(1, "userManage")
	server.Userdao.Add(&entity.UserEntity{
		Username: "root",
		Password: util.GetPWD("toor"),
		RoleID:   1,
	})
	r := gin.Default()
	r.POST("/api/login", server.login)
	api := r.Group("/api", server.CheckLogin)
	{
		api.GET("/username", func(c *gin.Context) {
			username, _ := c.Cookie("username")
			c.JSON(http.StatusOK, gin.H{"username": username})
		})
		user := api.Group("/user", server.userAdmin)
		{
			user.POST("/add", func(c *gin.Context) {
				var user entity.UserEntity
				err := c.ShouldBind(&user)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, message.Fail())
					return
				}
				err = server.Userdao.Add(&user)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusForbidden, gin.H{
						"err": err,
					})
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			user.GET("/userlist", func(context *gin.Context) {
				context.JSON(http.StatusOK, server.Userdao.UserList())
			})
			role := user.Group("/role")
			{
				role.GET("/rolelist", func(c *gin.Context) {
					c.JSON(http.StatusOK, server.RoleDao.RoleList())
				})
				role.POST("/add", func(c *gin.Context) {
					var input entity.RoleEntity
					err := c.ShouldBind(&input)
					if err != nil {
						c.JSON(http.StatusBadRequest, message.Fail())
						return
					}
					if server.RoleDao.AddRole(input.Name) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, message.Fail())
				})
				role.POST("/del", func(c *gin.Context) {
					var input entity.RoleEntity
					err := c.ShouldBind(&input)
					if err != nil {
						c.JSON(http.StatusBadRequest, message.Fail())
						return
					}
					if server.RoleDao.Del(input.Name) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, message.Fail())
				})
				role.POST("/permit", func(c *gin.Context) {
					m := struct {
						ID   int
						Auth string
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						c.JSON(http.StatusBadRequest, message.Fail())
						return
					}
					if server.RoleDao.AddAuth(m.ID, m.Auth) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, message.Fail())
				})
			}
		}
	}
	r.Run(":58888")
	util.PauseForRun()
}

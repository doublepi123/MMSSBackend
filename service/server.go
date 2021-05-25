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
	//新建角色Administrator
	server.RoleDao.AddRole("Administrator")
	//为Administrator角色赋予userManager权限
	server.RoleDao.AddAuth(1, "userManage")
	//新建用户root
	server.Userdao.Add(&entity.UserEntity{
		Username: "root",
		Password: util.GetPWD("toor"),
		RoleID:   1,
	})
	r := gin.Default()
	//登录接口
	r.POST("/api/login", server.login)
	//api的URL组
	api := r.Group("/api", server.CheckLogin)
	{
		//查询当前用户的用户名 /api/username
		api.GET("/username", func(c *gin.Context) {
			username, _ := c.Cookie("username")
			c.JSON(http.StatusOK, gin.H{"username": username})
		})
		self := api.Group("/self")
		{
			//修改自己的密码 /api/self/changepwd 	Oldpassword password
			self.POST("/changepwd", func(c *gin.Context) {
				m := struct {
					Oldpassword string
					Password    string
				}{}
				err := c.ShouldBind(&m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, message.Fail())
					return
				}
				username, err := c.Cookie("username")
				if !server.Userdao.Check(username, m.Oldpassword) {
					c.JSON(http.StatusForbidden, gin.H{
						"msg": "Oldpassword is wrong",
					})
					return
				}

				err = server.Userdao.ChangePassword(username, m.Password)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusForbidden, message.Fail())
				}
				c.JSON(http.StatusOK, message.Success())
			})
		}
		user := api.Group("/user", server.userAdmin)
		{
			//PATH /api/user
			//根据username查询某个用户 /api/user/find Username
			user.POST("/find", func(c *gin.Context) {
				m := struct {
					Username string
				}{}
				err := c.ShouldBind(&m)
				if err != nil {
					c.JSON(http.StatusBadRequest, message.Fail())
					return
				}
				ans := server.Userdao.Find(m.Username)
				ans.Password = "******"
				c.JSON(http.StatusOK, ans)
			})
			//查询权限列表 /api/user/authlist
			user.GET("/authlist", func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"auths": "userManager",
				})
			})
			//添加用户 /api/user/add
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
			//获取用户列表 /api/user/userlist
			user.GET("/userlist", func(context *gin.Context) {
				context.JSON(http.StatusOK, server.Userdao.UserList())
			})
			role := user.Group("/role")
			{
				//PATH /api/user/role
				//根据角色id查询角色名 /api/user/role/findname
				role.POST("/findname", func(c *gin.Context) {
					m := struct {
						ID int
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, message.Fail())
						return
					}
					c.JSON(http.StatusOK, gin.H{
						"name": server.RoleDao.GetRoleName(m.ID),
					})
				})
				//获取角色列表 /api/user/role/rolelist
				role.GET("/rolelist", func(c *gin.Context) {
					c.JSON(http.StatusOK, server.RoleDao.RoleList())
				})
				//添加角色	/api/user/role/add
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
				//删除角色	/api/user/role/del
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
				//为角色赋予权限	/api/user/role/permit
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
				//为角色移除权限	/api/user/role/ban
				role.POST("/ban", func(c *gin.Context) {
					m := struct {
						ID   int
						Auth string
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						c.JSON(http.StatusBadRequest, message.Fail())
						return
					}
					if server.RoleDao.DelAuth(m.ID, m.Auth) {
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

package service

import (
	"MMSSBackend/dao"
	"MMSSBackend/entity"
	"MMSSBackend/message"
	"MMSSBackend/util"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type Server struct {
	Userdao   *dao.UserDao
	Cookiedao *dao.CookieDao
	RoleDao   *dao.RoleDao
	PaperDao  *dao.PaperDao
}

func (server Server) login(c *gin.Context) {
	var user entity.UserEntity
	err := c.ShouldBind(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprint(err))
		return
	}
	if server.Userdao.Check(user.Username, user.Password) {
		server.Cookiedao.SetCookie(user.Username, server.Cookiedao.GenerateUserID(user.Username), c)
		c.JSON(http.StatusOK, message.Success())
		return
	}
	c.JSON(http.StatusForbidden, gin.H{"msg": "password wrong"})
}

func (server Server) CheckLogin(c *gin.Context) {
	username, _ := c.Cookie("username")
	userid, _ := c.Cookie("userid")
	if !server.Cookiedao.CheckCookie(username, userid) {
		c.JSON(http.StatusGone, gin.H{"msg": "not login"})
		c.Abort()
		return
	}
	server.Cookiedao.SetCookie(username, userid, c)
	c.Next()
}

func (server Server) userAdmin(c *gin.Context) {
	username, _ := c.Cookie("username")
	roleid := server.RoleDao.GetRoleID(username)
	if !server.RoleDao.CheckAuth(roleid, "userManage") {
		fmt.Println(roleid)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "not auth",
		})
		c.Abort()
		return
	}
	c.Next()
}

func (server Server) paperAdmin(c *gin.Context) {
	username, _ := c.Cookie("username")
	roleid := server.RoleDao.GetRoleID(username)
	if !server.RoleDao.CheckAuth(roleid, "paperManage") {
		fmt.Println(roleid)
		c.JSON(http.StatusForbidden, gin.H{
			"message": "not auth",
		})
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
	server.RoleDao.AddAuth(1, "paperManage")
	//新建用户root
	server.Userdao.Add(&entity.UserEntity{
		Username: "root",
		Password: util.GetPWD("toor"),
		RoleID:   1,
	})
	r := gin.Default()
	//登录接口
	r.POST("/api/login", server.login)
	//进入api路径前检查登入状态
	api := r.Group("/api", server.CheckLogin)
	{
		api.GET("/logout", func(c *gin.Context) {
			username, _ := c.Cookie("username")
			err := server.Userdao.Logout(username)
			if err != nil {
				util.MeetError(c, err)
				return
			}
			c.JSON(http.StatusOK, message.Success())
		})
		//查询当前用户的用户名 /api/username
		api.GET("/username", func(c *gin.Context) {
			username, _ := c.Cookie("username")
			c.JSON(http.StatusOK, gin.H{"username": username})
		})
		self := api.Group("/self")
		{
			//修改自己的密码 /api/self/changepwd 	仅两个字段：Oldpassword password
			self.POST("/changepwd", func(c *gin.Context) {
				m := struct {
					Oldpassword string
					Password    string
				}{}
				err := c.ShouldBind(&m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
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
					c.JSON(http.StatusForbidden, fmt.Sprint(err))
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//更新自己的的信息 /api/self/update 不需要传Username
			self.POST("/update", func(c *gin.Context) {
				username, _ := c.Cookie("username")
				m := entity.UserEntity{}
				err := c.ShouldBind(&m)
				if err != nil {
					fmt.Println(http.StatusBadRequest, fmt.Sprint(err))
					return
				}
				m.Username = username
				err = server.Userdao.Update(&m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
					return
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
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
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
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
					return
				}
				user.Password = util.GetPWD(user.Password)
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
			//更新用户信息/api/user/update
			user.POST("/update", func(c *gin.Context) {
				var user entity.UserEntity
				err := c.ShouldBind(&user)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
					return
				}
				err = server.Userdao.Update(&user)
				user.Password = ""
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//删除用户 /api/user/del	Username
			user.POST("/del", func(c *gin.Context) {
				var user entity.UserEntity
				err := c.ShouldBind(&user)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
					return
				}
				err = server.Userdao.Del(user.Username)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
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
						c.JSON(http.StatusBadRequest, fmt.Sprint(err))
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
						c.JSON(http.StatusBadRequest, fmt.Sprint(err))
						return
					}
					if server.RoleDao.AddRole(input.Name) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, fmt.Sprint(err))
				})
				//删除角色	/api/user/role/del
				role.POST("/del", func(c *gin.Context) {
					var input entity.RoleEntity
					err := c.ShouldBind(&input)
					if err != nil {
						c.JSON(http.StatusBadRequest, fmt.Sprint(err))
						return
					}
					if server.RoleDao.Del(input.Name) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, fmt.Sprint(err))
				})
				//为角色赋予权限	/api/user/role/permit
				role.POST("/permit", func(c *gin.Context) {
					m := struct {
						ID   int
						Auth string
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						c.JSON(http.StatusBadRequest, fmt.Sprint(err))
						return
					}
					if server.RoleDao.AddAuth(m.ID, m.Auth) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, fmt.Sprint(err))
				})
				//为角色移除权限	/api/user/role/ban
				role.POST("/ban", func(c *gin.Context) {
					m := struct {
						ID   int
						Auth string
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						c.JSON(http.StatusBadRequest, fmt.Sprint(err))
						return
					}
					if server.RoleDao.DelAuth(m.ID, m.Auth) {
						c.JSON(http.StatusOK, message.Success())
						return
					}
					c.JSON(http.StatusForbidden, fmt.Sprint(err))
				})
			}
		}
		paper := api.Group("/paper")
		{
			//该path下的接口仅可对自己的论文进行操作
			//添加paper /api/paper/add
			paper.POST("/add", func(c *gin.Context) {
				m := &entity.PaperEntity{}
				m.Hascheck = false
				username, _ := c.Cookie("username")
				err := c.ShouldBind(m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": err,
					})
					return
				}
				m.UserName = username
				m.Filename = ""
				err = server.PaperDao.Add(m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": fmt.Sprint(err),
					})
					return
				}

				c.JSON(http.StatusOK, message.Success())
			})
			//删除自己的论文 /api/paper/del
			paper.POST("/del", func(c *gin.Context) {
				username, _ := c.Cookie("username")
				m := struct {
					PaperID uint
				}{}
				paper, err := server.PaperDao.FindByID(m.PaperID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				if paper.UserName != username {
					c.JSON(http.StatusForbidden, gin.H{
						"msg": "not auth",
					})
				}
				err = server.PaperDao.Del(m.PaperID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//查看自己的paperList	/api/paper/list
			paper.GET("/list", func(c *gin.Context) {
				username, _ := c.Cookie("username")
				paper, err := server.PaperDao.Paperlist(username)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				c.JSON(http.StatusOK, paper)
			})
			//查找paper /api/paper/find
			paper.POST("/find", func(c *gin.Context) {
				m := &entity.PaperEntity{}
				username, _ := c.Cookie("username")
				err := c.ShouldBind(m)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": err,
					})
					return
				}
				m.UserName = username
				paper, err := server.PaperDao.Find(m.UserName, m.Tittle)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprint(err)})
					return
				}
				c.JSON(http.StatusOK, paper)
			})
			//模糊查找第一作者paper /api/paper/adfind
			paper.POST("/adfind", func(c *gin.Context) {
				m := &entity.PaperEntity{}
				username, _ := c.Cookie("username")
				err := c.ShouldBind(m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": fmt.Sprint(err),
					})
					return
				}
				m.UserName = username
				paper, err := server.PaperDao.ADFind(*m)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"msg": err})
					return
				}
				c.JSON(http.StatusOK, paper)
			})
			//添加其他作者 /api/paper/auth
			paper.POST("/auth", func(c *gin.Context) {
				m := &struct {
					PaperID uint
					WorkID  string
				}{}
				err := c.ShouldBind(m)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": fmt.Sprint(err),
					})
					return
				}
				username, _ := c.Cookie("username")
				err = server.PaperDao.Auth(m.WorkID, m.PaperID, username)
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{
						"msg": fmt.Sprint(err),
					})
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//更新文章信息 /api/paper/update
			paper.POST("/update", func(c *gin.Context) {
				var paper entity.PaperEntity
				err := c.ShouldBind(&paper)
				paper.UserName, _ = c.Cookie("username")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprint(err)})
					return
				}
				err = server.PaperDao.Update(paper)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprint(err)})
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//查询非第一作者paper /api/paper/findother
			paper.POST("/findother", func(c *gin.Context) {
				username, _ := c.Cookie("username")
				paper, err := server.PaperDao.AuthSelect(username)
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprint(err)})
					return
				}
				c.JSON(http.StatusOK, paper)
			})
			//添加附件 /api/paper/addfile
			paper.POST("/addfile", func(c *gin.Context) {
				paper := struct {
					PaperID uint
				}{}
				file, header, err := c.Request.FormFile("File")
				if err != nil {
					util.MeetError(c, err)
					return
				}
				content, err := ioutil.ReadAll(file)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				err = c.ShouldBind(&paper)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				username, _ := c.Cookie("username")
				u, err := server.PaperDao.GetPaper(paper.PaperID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				fmt.Println(u.UserName, username)
				if u.UserName != username {
					c.JSON(http.StatusForbidden, gin.H{
						"msg": "not auth",
					})
					return
				}
				err = server.PaperDao.AddFile(&entity.PaperFile{
					PaperID:  paper.PaperID,
					File:     content,
					FileName: header.Filename,
				})
				if err != nil {
					fmt.Println(err)
					c.JSON(http.StatusBadRequest, fmt.Sprint(err))
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//删除某人的文章 /api/paper/delete
			paper.DELETE("/delete", func(c *gin.Context) {
				m := struct {
					ID uint
				}{}
				err := c.ShouldBind(&m)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				err = server.PaperDao.DelJournal(m.ID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				c.JSON(http.StatusOK, message.Success())
			})
			//下载附件 /api/paper/getfile
			paper.POST("/getfile", func(c *gin.Context) {
				paper := struct {
					PaperID uint
				}{}
				c.ShouldBind(&paper)
				username, _ := c.Cookie("username")
				u, err := server.PaperDao.GetPaper(paper.PaperID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				if username != u.UserName {
					util.MeetError(c, errors.New("not auth"))
					return
				}
				file, err := server.PaperDao.GetFile(paper.PaperID)
				if err != nil {
					util.MeetError(c, err)
					return
				}
				c.Writer.WriteHeader(http.StatusOK)
				content := file.File
				c.Header("Content-Disposition", "attachment; filename="+file.FileName)
				c.Header("Content-Type", "application/text/plain")
				c.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
				c.Writer.Write([]byte(content))
			})
			//获取某人所有的paperlist（包括不是本人为第一作者的论文）/api/paper/alllist
			paper.GET("/alllist", func(c *gin.Context) {
				username, _ := c.Cookie("username")
				c.JSON(http.StatusOK, server.PaperDao.GetSomeoneAllPaper(username))
			})
			//期刊
			journal := paper.Group("/journal")
			{
				//获取期刊列表 /api/paper/journal/list
				journal.GET("/list", func(c *gin.Context) {
					c.JSON(http.StatusOK, server.PaperDao.GetJournalList())
				})
				//添加期刊	/api/paper/journal/add
				journal.POST("/add", func(c *gin.Context) {
					var journal entity.Journal
					err := c.ShouldBind(&journal)
					if err != nil {
						util.MeetError(c, err)
						return
					}
					err = server.PaperDao.AddJournal(journal)
					if err != nil {
						util.MeetError(c, err)
						return
					}
					c.JSON(http.StatusOK, message.Success())
				})
				ja := journal.Group("/admin", server.paperAdmin)
				{
					//获取期刊列表 /api/paper/journal/admin/list
					ja.GET("/list", func(c *gin.Context) {
						c.JSON(http.StatusOK, server.PaperDao.GetJournalList())
					})
					//获取未审核期刊列表	/api/paper/journal/admin/uncheck
					ja.GET("/uncheck", func(c *gin.Context) {
						c.JSON(http.StatusOK, server.PaperDao.GetUncheckJournalList())
					})
					//审核某个期刊信息	/api/paper/journal/admin/check 发送字段仅为 ID
					ja.POST("/check", func(c *gin.Context) {
						m := struct {
							ID uint
						}{}
						err := c.ShouldBind(&m)
						if err != nil {
							util.MeetError(c, err)
							return
						}
						err = server.PaperDao.CheckJournal(m.ID)
						if err != nil {
							util.MeetError(c, err)
							return
						}
						c.JSON(http.StatusOK, message.Success())
					})
					//删除某个期刊信息	/api/paper/journal/admin/delete 发送字段仅为 ID
					ja.DELETE("/delete", func(c *gin.Context) {
						m := struct {
							ID uint
						}{}
						err := c.ShouldBind(&m)
						if err != nil {
							util.MeetError(c, err)
							return
						}
						err = server.PaperDao.DelJournal(m.ID)
						if err != nil {
							util.MeetError(c, err)
							return
						}
						c.JSON(http.StatusOK, message.Success())
					})
				}
			}
			//管理员的操作
			papera := paper.Group("/admin", server.paperAdmin)
			{
				//模糊查询论文 /api/paper/admin/adfind
				papera.POST("/adfind", func(c *gin.Context) {
					m := &entity.PaperEntity{}
					err := c.ShouldBind(m)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"msg": fmt.Sprint(err),
						})
						return
					}
					paper, err := server.PaperDao.ADFind(*m)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{"msg": err})
						return
					}
					c.JSON(http.StatusOK, paper)
				})
				//查询某篇论文 /api/paper/admin/find
				papera.POST("/find", func(c *gin.Context) {
					m := &entity.PaperEntity{}
					err := c.ShouldBind(m)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{
							"msg": err,
						})
						return
					}
					paper, err := server.PaperDao.Find(m.UserName, m.Tittle)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{"msg": fmt.Sprint(err)})
						return
					}
					c.JSON(http.StatusOK, paper)
				})
				//获取未审核的论文列表	/api/paper/admin/uncheck
				papera.GET("/uncheck", func(c *gin.Context) {
					papers, err := server.PaperDao.GetUncheckFile()
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
						return
					}
					c.JSON(http.StatusOK, papers)
				})
				//审核某篇论文 /api/paper/admin/check	参数发送PaperID int
				papera.POST("/check", func(c *gin.Context) {
					m := struct {
						PaperID uint
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"msg": fmt.Sprint(err),
						})
					}
					err = server.PaperDao.Check(m.PaperID)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusForbidden, gin.H{"msg": fmt.Sprint(err)})
						return
					}
					c.JSON(http.StatusOK, message.Success())
				})
				//取消某篇论文的审核状态 /api/paper/admin/uncheck	参数同上
				papera.POST("/uncheck", func(c *gin.Context) {
					m := struct {
						PaperID uint
					}{}
					err := c.ShouldBind(&m)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusBadRequest, gin.H{
							"msg": fmt.Sprint(err),
						})
					}
					err = server.PaperDao.UnCheck(m.PaperID)
					if err != nil {
						fmt.Println(err)
						c.JSON(http.StatusForbidden, gin.H{"msg": fmt.Sprint(err)})
						return
					}
					c.JSON(http.StatusOK, message.Success())
				})
			}
		}
	}
	r.Run(":58888")
	util.PauseForRun()
}

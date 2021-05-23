package dao

import (
	"MMSSBackend/util"
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

const validTime = time.Minute * 15

type CookieDao struct {
	*Dao
}

func (cookie CookieDao) GenerateUserID(username string) string {
	return fmt.Sprint(time.Time.Unix) + util.GetPWD(username)
}

func (cookie CookieDao) SetCookie(username string, userid string, c *gin.Context) error {
	c.SetCookie("username", username, int(validTime), "/", c.Request.Host, false, true)
	c.SetCookie("userid", userid, int(validTime), "/", c.Request.Host, false, true)
	return cookie.db.Redis.Set(username, userid, validTime).Err()
}

func (cookie CookieDao) UpdateCookie(c *gin.Context) error {
	username, err := c.Cookie("username")
	if err != nil {
		return err
	}

	Userid, err := c.Cookie("userid")
	if err != nil {
		return err
	}
	return cookie.SetCookie(username, Userid, c)
}

func (cookie CookieDao) CheckCookie(username string, userid string) bool {
	if username == "" || userid == "" {
		return false
	}
	ud, err := cookie.db.Redis.Get(username).Result()
	if err != nil {
		return false
	}
	return ud == userid
}

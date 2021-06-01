package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"time"
)

func PauseForRun() {
	for {
		time.Sleep(time.Second)
	}
}

func MeetError(c *gin.Context, err error) {
	fmt.Println(err)
	c.JSON(http.StatusBadRequest, gin.H{
		"msg": fmt.Sprint(err),
	})
}

func ShowBody(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
}

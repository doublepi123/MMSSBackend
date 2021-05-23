package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"time"
)

func PauseForRun() {
	for {
		time.Sleep(time.Second)
	}
}

func ShowBody(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println(string(data))
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(data))
}

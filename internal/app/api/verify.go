package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gowechat/example/pkg/core"
	Log "github.com/gowechat/example/pkg/core"
	"github.com/gowechat/example/pkg/util/app"
)

func Verify(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	var pwd string
	err := core.Db.QueryRow("select MM from NEWJW.SYS_YHMMB where ZJH = :1", username).Scan(&pwd)
	if pwd != "" && err != nil {
		Log.Logrus().Println(err)
		return
	}
	if password == pwd {
		app.OK(c, map[string]interface{}{"info": "密码正确"}, "ok")
		return
	} else {
		app.PwdError(c, map[string]interface{}{"info": "密码错误"}, "false")
	}
}
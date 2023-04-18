package api

import (
	"github.com/gin-gonic/gin"
)

func PersonalCenter(c *gin.Context) {
	c.HTML(200, "index.html", gin.H{
		"title": "Gin 测试模板666",
	})
}
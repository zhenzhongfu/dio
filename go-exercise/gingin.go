package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 模拟私有数据
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	r := gin.Default()

	// 使用 gin.BasicAuth 中间件，设置授权用户
	authorized := r.Group("/admin", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))

	// 定义路由
	authorized.GET("/secrets", func(c *gin.Context) {
		fmt.Println("----------------- ", c.GetHeader("Authorization"))
		// 获取提交的用户名（AuthUserKey）
		user := c.MustGet(gin.AuthUserKey).(string)
		if secret, ok := secrets[user]; ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": secret})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "secret": "NO SECRET :("})
		}
	})

	// 使用浏览器打开。。。
	// Listen and serve on 0.0.0.0:8080
	r.Run(":28080")
}

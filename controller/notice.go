package controller

import (
	"github.com/gin-gonic/gin"
	"goflylivechat/common"
	"goflylivechat/models"
	"strings"
)

func GetNotice(c *gin.Context) {
	kefuId := c.Query("kefu_id")
	user := models.FindUser(kefuId)
	if user.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "user not found",
		})
		return
	}
	welcomeMessage := models.FindConfigByUserId(user.Name, "WelcomeMessage")
	offlineMessage := models.FindConfigByUserId(user.Name, "OfflineMessage")
	allNotice := models.FindConfigByUserId(user.Name, "AllNotice")

	// 动态处理头像路径
	basePath := common.GetDynamicBasePath(c)
	avatar := user.Avator
	if avatar != "" && !strings.HasPrefix(avatar, basePath) {
		// 如果头像路径没有前缀，添加动态前缀
		avatar = basePath + avatar
	}

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"welcome":   welcomeMessage.ConfValue,
			"offline":   offlineMessage.ConfValue,
			"avatar":    avatar,
			"nickname":  user.Nickname,
			"allNotice": allNotice.ConfValue,
		},
	})
}

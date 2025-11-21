package router

import (
	"goflylivechat/common"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitViewRouter(engine *gin.Engine) {
	// 注册带前缀的路由（代理访问）
	if common.IsPrefixEnabled() {
		prefix := common.GetPrefix()

		engine.GET(prefix+"/login", PageLogin)
		engine.GET(prefix+"/pannel", PagePannel)
		engine.GET(prefix+"/livechat", PageChat)
		engine.GET(prefix+"/main", PageMain)
		engine.GET(prefix+"/chat_main", PageChatMain)
		engine.GET(prefix+"/setting", PageSetting)
	}

	// 注册无前缀的路由（直接访问）
	engine.GET("/login", PageLogin)
	engine.GET("/pannel", PagePannel)
	engine.GET("/livechat", PageChat)
	engine.GET("/main", PageMain)
	engine.GET("/chat_main", PageChatMain)
	engine.GET("/setting", PageSetting)
}

// PageLogin Login page
func PageLogin(c *gin.Context) {
	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "login.html", gin.H{
		"BasePath": basePath,
	})
}

// PagePannel Dashboard
func PagePannel(c *gin.Context) {
	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "pannel.html", gin.H{
		"BasePath": basePath,
	})
}

// PageMain Admin console
func PageMain(c *gin.Context) {
	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "main.html", gin.H{
		"BasePath": basePath,
	})
}

// PageChat Customer chat interface
func PageChat(c *gin.Context) {
	referralSource := c.Query("refer") // More clear variable name

	if referralSource == "" {
		referralSource = c.Request.Referer()
	}
	if referralSource == "" {
		referralSource = "Direct access" // More natural English
	}

	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"Refer":    referralSource, // Keeping original template variable name
		"BasePath": basePath,
	})
}

// PageChatMain Support agent console
func PageChatMain(c *gin.Context) {
	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "chat_main.html", gin.H{
		"BasePath": basePath,
	})
}

// PageSetting Settings
func PageSetting(c *gin.Context) {
	basePath := common.GetDynamicBasePath(c)

	c.HTML(http.StatusOK, "setting.html", gin.H{
		"BasePath": basePath,
	})
}

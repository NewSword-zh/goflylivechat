package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitViewRouter(engine *gin.Engine) {
	//engine.GET("/", tmpl.PageIndex)
	engine.GET("/login", PageLogin)
	engine.GET("/pannel", PagePannel)
	engine.GET("/livechat", PageChat)
	engine.GET("/main", PageMain)
	engine.GET("/chat_main", PageChatMain)
	engine.GET("/setting", PageSetting)
}

// PageLogin Login page
func PageLogin(c *gin.Context) {
	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"BasePath": basePath,
	})
}

// PagePannel Dashboard
func PagePannel(c *gin.Context) {
	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

	c.HTML(http.StatusOK, "pannel.html", gin.H{
		"BasePath": basePath,
	})
}

// PageMain Admin console
func PageMain(c *gin.Context) {
	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

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

	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"Refer":     referralSource, // Keeping original template variable name
		"BasePath":  basePath,
	})
}

// PageChatMain Support agent console
func PageChatMain(c *gin.Context) {
	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

	c.HTML(http.StatusOK, "chat_main.html", gin.H{
		"BasePath": basePath,
	})
}

// PageSetting Settings
func PageSetting(c *gin.Context) {
	// 通过 HTTP Header 检测是否为代理模式
	basePath := ""
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		basePath = "/goflychat"
	}

	c.HTML(http.StatusOK, "setting.html", gin.H{
		"BasePath": basePath,
	})
}

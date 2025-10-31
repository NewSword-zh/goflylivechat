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
	c.HTML(http.StatusOK, "login.html", nil)
}

// PagePannel Dashboard
func PagePannel(c *gin.Context) {
	c.HTML(http.StatusOK, "pannel.html", nil)
}

// PageMain Admin console
func PageMain(c *gin.Context) {
	c.HTML(http.StatusOK, "main.html", nil)
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

	c.HTML(http.StatusOK, "chat_page.html", gin.H{
		"Refer": referralSource, // Keeping original template variable name
	})
}

// PageChatMain Support agent console
func PageChatMain(c *gin.Context) {
	c.HTML(http.StatusOK, "chat_main.html", nil)
}

// PageSetting Settings
func PageSetting(c *gin.Context) {
	c.HTML(http.StatusOK, "setting.html", nil)
}

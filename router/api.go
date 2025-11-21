package router

import (
	"goflylivechat/common"
	"goflylivechat/controller"
	"goflylivechat/middleware"
	"goflylivechat/ws"

	"github.com/gin-gonic/gin"
)

func InitApiRouter(engine *gin.Engine) {
	// 注册带前缀的API路由（代理访问）
	if common.IsPrefixEnabled() {
		prefix := common.GetPrefix()

		//路由分组
		v2WithPrefix := engine.Group(prefix + "/2")
		{
			//获取消息
			v2WithPrefix.GET("/messages", controller.GetMessagesV2)
			//发送单条信息
			v2WithPrefix.POST("/message", middleware.Ipblack, controller.SendMessageV2)
			//关闭连接
			v2WithPrefix.GET("/message_close", controller.SendCloseMessageV2)
			//分页查询消息
			v2WithPrefix.GET("/messagesPages", controller.GetMessagespages)
		}

		engine.GET(prefix+"/captcha", controller.GetCaptcha)
		engine.POST(prefix+"/check", controller.LoginCheckPass)

		engine.GET(prefix+"/userinfo", middleware.JwtApiMiddleware, controller.GetKefuInfoAll)
		engine.POST(prefix+"/register", middleware.Ipblack, controller.PostKefuRegister)
		engine.POST(prefix+"/install", controller.PostInstall)
		//前后聊天
		engine.GET(prefix+"/ws_kefu", middleware.JwtApiMiddleware, ws.NewKefuServer)
		engine.GET(prefix+"/ws_visitor", middleware.Ipblack, ws.NewVisitorServer)

		engine.GET(prefix+"/messages", controller.GetVisitorMessage)
		engine.GET(prefix+"/message_notice", controller.SendVisitorNotice)
		//上传文件
		engine.POST(prefix+"/uploadimg", middleware.Ipblack, controller.UploadImg)
		//上传文件
		engine.POST(prefix+"/uploadfile", middleware.Ipblack, controller.UploadFile)
		//获取未读消息数
		engine.GET(prefix+"/message_status", controller.GetVisitorMessage)
		//设置消息已读
		engine.POST(prefix+"/message_status", controller.GetVisitorMessage)

		//获取客服信息
		engine.POST(prefix+"/kefuinfo_client", middleware.JwtApiMiddleware, controller.PostKefuClient)
		engine.GET(prefix+"/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfo)
		engine.GET(prefix+"/kefuinfo_setting", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfoSetting)
		engine.POST(prefix+"/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuInfo)
		engine.DELETE(prefix+"/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DeleteKefuInfo)
		engine.GET(prefix+"/kefulist", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuList)
		engine.GET(prefix+"/other_kefulist", middleware.JwtApiMiddleware, controller.GetOtherKefuList)
		engine.GET(prefix+"/trans_kefu", middleware.JwtApiMiddleware, controller.PostTransKefu)
		engine.POST(prefix+"/modifypass", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuPass)
		engine.POST(prefix+"/modifyavator", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuAvator)
		//角色列表
		engine.GET(prefix+"/roles", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetRoleList)
		engine.POST(prefix+"/role", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostRole)

		engine.GET(prefix+"/visitors_online", controller.GetVisitorOnlines)
		engine.GET(prefix+"/visitors_kefu_online", middleware.JwtApiMiddleware, controller.GetKefusVisitorOnlines)
		engine.GET(prefix+"/clear_online_tcp", controller.DeleteOnlineTcp)
		engine.POST(prefix+"/visitor_login", middleware.Ipblack, controller.PostVisitorLogin)
		//engine.POST("/visitor", controller.PostVisitor)
		engine.GET(prefix+"/visitor", middleware.JwtApiMiddleware, controller.GetVisitor)
		engine.GET(prefix+"/visitors", middleware.JwtApiMiddleware, controller.GetVisitors)
		engine.GET(prefix+"/statistics", middleware.JwtApiMiddleware, controller.GetStatistics)
		//前台接口
		engine.GET(prefix+"/about", controller.GetAbout)
		engine.POST(prefix+"/about", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostAbout)
		engine.GET(prefix+"/aboutpages", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetAbouts)
		engine.GET(prefix+"/notice", controller.GetNotice)
		engine.POST(prefix+"/ipblack", middleware.JwtApiMiddleware, middleware.Ipblack, controller.PostIpblack)
		engine.DELETE(prefix+"/ipblack", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelIpblack)
		engine.GET(prefix+"/ipblacks_all", middleware.JwtApiMiddleware, controller.GetIpblacks)
		engine.GET(prefix+"/ipblacks", middleware.JwtApiMiddleware, controller.GetIpblacksByKefuId)
		engine.GET(prefix+"/configs", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetConfigs)
		engine.POST(prefix+"/config", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostConfig)
		engine.GET(prefix+"/config", controller.GetConfig)
		engine.GET(prefix+"/autoreply", controller.GetAutoReplys)
		engine.GET(prefix+"/replys", middleware.JwtApiMiddleware, controller.GetReplys)
		engine.POST(prefix+"/reply", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReply)
		engine.POST(prefix+"/reply_content", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReplyContent)
		engine.POST(prefix+"/reply_content_save", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReplyContentSave)
		engine.DELETE(prefix+"/reply_content", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelReplyContent)
		engine.DELETE(prefix+"/reply", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelReplyGroup)
		engine.POST(prefix+"/reply_search", middleware.JwtApiMiddleware, controller.PostReplySearch)
		//客服路由分组
		kefuGroup := engine.Group(prefix + "/kefu")
		kefuGroup.Use(middleware.JwtApiMiddleware)
		{
			kefuGroup.GET("/chartStatistics", controller.GetChartStatistic)
			kefuGroup.POST("/message", controller.SendKefuMessage)
		}
		//微信接口
		engine.GET(prefix+"/micro_program", middleware.JwtApiMiddleware, controller.GetCheckWeixinSign)
	}

	// 注册无前缀的API路由（直接访问）
	//路由分组
	v2 := engine.Group("/2")
	{
		//获取消息
		v2.GET("/messages", controller.GetMessagesV2)
		//发送单条信息
		v2.POST("/message", middleware.Ipblack, controller.SendMessageV2)
		//关闭连接
		v2.GET("/message_close", controller.SendCloseMessageV2)
		//分页查询消息
		v2.GET("/messagesPages", controller.GetMessagespages)
	}

	engine.GET("/captcha", controller.GetCaptcha)
	engine.POST("/check", controller.LoginCheckPass)

	engine.GET("/userinfo", middleware.JwtApiMiddleware, controller.GetKefuInfoAll)
	engine.POST("/register", middleware.Ipblack, controller.PostKefuRegister)
	engine.POST("/install", controller.PostInstall)
	//前后聊天
	engine.GET("/ws_kefu", middleware.JwtApiMiddleware, ws.NewKefuServer)
	engine.GET("/ws_visitor", middleware.Ipblack, ws.NewVisitorServer)

	engine.GET("/messages", controller.GetVisitorMessage)
	engine.GET("/message_notice", controller.SendVisitorNotice)
	//上传文件
	engine.POST("/uploadimg", middleware.Ipblack, controller.UploadImg)
	//上传文件
	engine.POST("/uploadfile", middleware.Ipblack, controller.UploadFile)
	//获取未读消息数
	engine.GET("/message_status", controller.GetVisitorMessage)
	//设置消息已读
	engine.POST("/message_status", controller.GetVisitorMessage)

	//获取客服信息
	engine.POST("/kefuinfo_client", middleware.JwtApiMiddleware, controller.PostKefuClient)
	engine.GET("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfo)
	engine.GET("/kefuinfo_setting", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfoSetting)
	engine.POST("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuInfo)
	engine.DELETE("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DeleteKefuInfo)
	engine.GET("/kefulist", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuList)
	engine.GET("/other_kefulist", middleware.JwtApiMiddleware, controller.GetOtherKefuList)
	engine.GET("/trans_kefu", middleware.JwtApiMiddleware, controller.PostTransKefu)
	engine.POST("/modifypass", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuPass)
	engine.POST("/modifyavator", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostKefuAvator)
	//角色列表
	engine.GET("/roles", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetRoleList)
	engine.POST("/role", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostRole)

	engine.GET("/visitors_online", controller.GetVisitorOnlines)
	engine.GET("/visitors_kefu_online", middleware.JwtApiMiddleware, controller.GetKefusVisitorOnlines)
	engine.GET("/clear_online_tcp", controller.DeleteOnlineTcp)
	engine.POST("/visitor_login", middleware.Ipblack, controller.PostVisitorLogin)
	//engine.POST("/visitor", controller.PostVisitor)
	engine.GET("/visitor", middleware.JwtApiMiddleware, controller.GetVisitor)
	engine.GET("/visitors", middleware.JwtApiMiddleware, controller.GetVisitors)
	engine.GET("/statistics", middleware.JwtApiMiddleware, controller.GetStatistics)
	//前台接口
	engine.GET("/about", controller.GetAbout)
	engine.POST("/about", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostAbout)
	engine.GET("/aboutpages", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetAbouts)
	engine.GET("/notice", controller.GetNotice)
	engine.POST("/ipblack", middleware.JwtApiMiddleware, middleware.Ipblack, controller.PostIpblack)
	engine.DELETE("/ipblack", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelIpblack)
	engine.GET("/ipblacks_all", middleware.JwtApiMiddleware, controller.GetIpblacks)
	engine.GET("/ipblacks", middleware.JwtApiMiddleware, controller.GetIpblacksByKefuId)
	engine.GET("/configs", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetConfigs)
	engine.POST("/config", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostConfig)
	engine.GET("/config", controller.GetConfig)
	engine.GET("/autoreply", controller.GetAutoReplys)
	engine.GET("/replys", middleware.JwtApiMiddleware, controller.GetReplys)
	engine.POST("/reply", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReply)
	engine.POST("/reply_content", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReplyContent)
	engine.POST("/reply_content_save", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.PostReplyContentSave)
	engine.DELETE("/reply_content", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelReplyContent)
	engine.DELETE("/reply", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.DelReplyGroup)
	engine.POST("/reply_search", middleware.JwtApiMiddleware, controller.PostReplySearch)
	//客服路由分组
	kefuGroup := engine.Group("/kefu")
	kefuGroup.Use(middleware.JwtApiMiddleware)
	{
		kefuGroup.GET("/chartStatistics", controller.GetChartStatistic)
		kefuGroup.POST("/message", controller.SendKefuMessage)
	}
	//微信接口
	engine.GET("/micro_program", middleware.JwtApiMiddleware, controller.GetCheckWeixinSign)
}

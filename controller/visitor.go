package controller

import (
	"encoding/json"
	"fmt"
	"goflylivechat/common"
	"goflylivechat/models"
	"goflylivechat/tools"
	"goflylivechat/ws"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

//	func PostVisitor(c *gin.Context) {
//		name := c.PostForm("name")
//		avator := c.PostForm("avator")
//		toId := c.PostForm("to_id")
//		id := c.PostForm("id")
//		refer := c.PostForm("refer")
//		city := c.PostForm("city")
//		client_ip := c.PostForm("client_ip")
//		if name == "" || avator == "" || toId == "" || id == "" || refer == "" || city == "" || client_ip == "" {
//			c.JSON(200, gin.H{
//				"code": 400,
//				"msg":  "error",
//			})
//			return
//		}
//		kefuInfo := models.FindUser(toId)
//		if kefuInfo.ID == 0 {
//			c.JSON(200, gin.H{
//				"code": 400,
//				"msg":  "用户不存在",
//			})
//			return
//		}
//		models.CreateVisitor(name, avator, c.ClientIP(), toId, id, refer, city, client_ip)
//
//		userInfo := make(map[string]string)
//		userInfo["uid"] = id
//		userInfo["username"] = name
//		userInfo["avator"] = avator
//		msg := TypeMessage{
//			Type: "userOnline",
//			Data: userInfo,
//		}
//		str, _ := json.Marshal(msg)
//		kefuConns := kefuList[toId]
//		if kefuConns != nil {
//			for k, kefuConn := range kefuConns {
//				log.Println(k, "xxxxxxxx")
//				kefuConn.WriteMessage(websocket.TextMessage, str)
//			}
//		}
//		c.JSON(200, gin.H{
//			"code": 200,
//			"msg":  "ok",
//		})
//	}
func PostVisitorLogin(c *gin.Context) {
	// 调试：检查所有相关的 Header
	fmt.Printf("DEBUG: Headers - X-Proxy-Mode: '%s', X-Forwarded-For: '%s', Host: '%s', Referer: '%s'\n",
		c.GetHeader("X-Proxy-Mode"),
		c.GetHeader("X-Forwarded-For"),
		c.GetHeader("Host"),
		c.GetHeader("Referer"))

	// 动态设置头像路径
	basePath := common.GetDynamicBasePath(c)
	fmt.Printf("DEBUG: 动态获取的基础路径: '%s'\n", basePath)

	avator := ""
	userAgent := c.GetHeader("User-Agent")
	if tools.IsMobile(userAgent) {
		avator = basePath + "/static/images/1.png"
	} else {
		avator = basePath + "/static/images/2.png"
	}

	toId := c.PostForm("to_id")
	id := c.PostForm("visitor_id")

	if id == "" {
		id = tools.Uuid()
	}
	refer := c.PostForm("refer")
	name := "Guest"
	city := ""
	countryname, cityname := tools.GetCity("./config/GeoLite2-City.mmdb", c.ClientIP())
	if countryname != "" || cityname != "" {
		city = fmt.Sprintf("%s %s", countryname, cityname)
		name = fmt.Sprintf("%s Guest", city)
	}

	client_ip := c.ClientIP()
	extra := c.PostForm("extra")
	extraJson := tools.Base64Decode(extra)
	if extraJson != "" {
		var extraObj VisitorExtra
		err := json.Unmarshal([]byte(extraJson), &extraObj)
		if err == nil {
			if extraObj.VisitorName != "" {
				name = extraObj.VisitorName
			}
			if extraObj.VisitorAvatar != "" {
				avator = extraObj.VisitorAvatar
			}
		}
	}
	//log.Println(name,avator,c.ClientIP(),toId,id,refer,city,client_ip)
	if name == "" || avator == "" || toId == "" || id == "" || refer == "" || client_ip == "" {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "error",
		})
		return
	}
	kefuInfo := models.FindUser(toId)
	if kefuInfo.ID == 0 {
		c.JSON(200, gin.H{
			"code": 400,
			"msg":  "The customer service account does not exist",
		})
		return
	}
	visitor := models.FindVisitorByVistorId(id)
	if visitor.Name != "" {
		// 检查数据库中的路径是否已经有前缀
		if !strings.HasPrefix(visitor.Avator, basePath) {
			// 如果没有前缀，使用动态生成的新路径
			avator = basePath + "/static/images/2.png"
		} else {
			// 如果已经有前缀，保持数据库中的路径
			avator = visitor.Avator
		}
		//更新状态上线，使用修正后的头像路径
		models.UpdateVisitor(name, avator, id, 1, c.ClientIP(), c.ClientIP(), refer, extra)
	} else {
		// 新访客，直接使用动态生成的路径
		models.CreateVisitor(name, avator, c.ClientIP(), toId, id, refer, city, client_ip, extra)
	}
	visitor.Name = name
	visitor.Avator = avator
	visitor.ToId = toId
	visitor.ClientIp = c.ClientIP()
	visitor.VisitorId = id

	//各种通知
	go SendNoticeEmail(visitor.Name, " incoming!")
	//go SendAppGetuiPush(kefuInfo.Name, visitor.Name, visitor.Name+" incoming!")
	go SendVisitorLoginNotice(kefuInfo.Name, visitor.Name, visitor.Avator, visitor.Name+" incoming!", visitor.VisitorId)
	go ws.VisitorOnline(kefuInfo.Name, visitor)
	//go SendServerJiang(visitor.Name, "来了", c.Request.Host)

	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": visitor,
	})
}
func GetVisitor(c *gin.Context) {
	visitorId := c.Query("visitorId")
	vistor := models.FindVisitorByVistorId(visitorId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": vistor,
	})
}

// @Summary 获取访客列表接口
// @Produce  json
// @Accept multipart/form-data
// @Param page query   string true "分页"
// @Param token header string true "认证token"
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors [get]
func GetVisitors(c *gin.Context) {
	page, _ := strconv.Atoi(c.Query("page"))
	pagesize, _ := strconv.Atoi(c.Query("pagesize"))
	if pagesize == 0 {
		pagesize = 10
	}
	kefuId, _ := c.Get("kefu_name")
	vistors := models.FindVisitorsByKefuId(uint(page), uint(pagesize), kefuId.(string))
	count := models.CountVisitorsByKefuId(kefuId.(string))
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"list":     vistors,
			"count":    count,
			"pagesize": common.PageSize,
		},
	})
}

// @Summary 获取访客聊天信息接口
// @Produce  json
// @Accept multipart/form-data
// @Param visitorId query   string true "访客ID"
// @Param token header string true "认证token"
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /messages [get]
func GetVisitorMessage(c *gin.Context) {
	visitorId := c.Query("visitorId")

	query := "message.visitor_id= ?"
	messages := models.FindMessageByWhere(query, visitorId)
	result := make([]map[string]interface{}, 0)
	for _, message := range messages {
		item := make(map[string]interface{})

		item["time"] = message.CreatedAt.Format("2006-01-02 15:04:05")
		item["content"] = message.Content
		item["mes_type"] = message.MesType
		item["visitor_name"] = message.VisitorName
		item["visitor_avator"] = message.VisitorAvator
		item["kefu_name"] = message.KefuName
		item["kefu_avator"] = message.KefuAvator
		result = append(result, item)

	}
	go models.ReadMessageByVisitorId(visitorId)
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": result,
	})
}

// @Summary 获取在线访客列表接口
// @Produce  json
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors_online [get]
func GetVisitorOnlines(c *gin.Context) {
	users := make([]map[string]string, 0)
	visitorIds := make([]string, 0)
	for uid, visitor := range ws.ClientList {
		userInfo := make(map[string]string)
		userInfo["uid"] = uid
		userInfo["name"] = visitor.Name
		userInfo["avator"] = visitor.Avator
		users = append(users, userInfo)
		visitorIds = append(visitorIds, visitor.Id)
	}

	//查询最新消息
	messages := models.FindLastMessage(visitorIds)
	temp := make(map[string]string, 0)
	for _, mes := range messages {
		temp[mes.VisitorId] = mes.Content
	}
	for _, user := range users {
		user["last_message"] = temp[user["uid"]]
	}

	tcps := make([]string, 0)
	for ip, _ := range clientTcpList {
		tcps = append(tcps, ip)
	}
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
		"result": gin.H{
			"ws":  users,
			"tcp": tcps,
		},
	})
}

// @Summary 获取客服的在线访客列表接口
// @Produce  json
// @Success 200 {object} controller.Response
// @Failure 200 {object} controller.Response
// @Router /visitors_kefu_online [get]
func GetKefusVisitorOnlines(c *gin.Context) {
	kefuName, _ := c.Get("kefu_name")
	users := make([]*VisitorOnline, 0)
	visitorIds := make([]string, 0)
	for uid, visitor := range ws.ClientList {
		if visitor.To_id != kefuName {
			continue
		}
		userInfo := new(VisitorOnline)
		userInfo.Uid = uid
		userInfo.Username = visitor.Name
		userInfo.Avator = visitor.Avator
		users = append(users, userInfo)
		visitorIds = append(visitorIds, visitor.Id)
	}

	//查询最新消息
	messages := models.FindLastMessage(visitorIds)
	temp := make(map[string]string, 0)
	for _, mes := range messages {
		temp[mes.VisitorId] = mes.Content
	}
	for _, user := range users {
		user.LastMessage = temp[user.Uid]
		if user.LastMessage == "" {
			user.LastMessage = "new visitor"
		}
	}

	tcps := make([]string, 0)
	for ip, _ := range clientTcpList {
		tcps = append(tcps, ip)
	}
	c.JSON(200, gin.H{
		"code":   200,
		"msg":    "ok",
		"result": users,
	})
}

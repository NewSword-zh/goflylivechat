package common

import (
	"encoding/json"
	"fmt"
	"goflylivechat/tools"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

type Mysql struct {
	Server   string
	Port     string
	Database string
	Username string
	Password string
}

type AppConfig struct {
	App App `json:"app"`
}

type App struct {
	Prefix       string `json:"prefix"`
	EnablePrefix bool   `json:"enable_prefix"`
}

func GetMysqlConf() *Mysql {
	var mysql = &Mysql{}
	isExist, _ := tools.IsFileExist(MysqlConf)
	if !isExist {
		return mysql
	}
	info, err := ioutil.ReadFile(MysqlConf)
	if err != nil {
		return mysql
	}
	err = json.Unmarshal(info, mysql)
	return mysql
}

func GetAppConf() *AppConfig {
	var appConfig = &AppConfig{
		App: App{
			Prefix:       "",
			EnablePrefix: false,
		},
	}
	isExist, _ := tools.IsFileExist(AppConf)
	if !isExist {
		return appConfig
	}
	info, err := ioutil.ReadFile(AppConf)
	if err != nil {
		return appConfig
	}
	err = json.Unmarshal(info, appConfig)
	return appConfig
}

// GetPrefix 获取配置的前缀
func GetPrefix() string {
	appConfig := GetAppConf()
	if appConfig.App.EnablePrefix {
		return appConfig.App.Prefix
	}
	return ""
}

// GetBasePath 获取基础路径，用于模板渲染
func GetBasePath() string {
	return GetPrefix()
}

// IsPrefixEnabled 检查是否启用前缀
func IsPrefixEnabled() bool {
	appConfig := GetAppConf()
	return appConfig.App.EnablePrefix
}

// GetDynamicBasePath 根据请求动态获取基础路径
func GetDynamicBasePath(c *gin.Context) string {
	// 优先检测代理模式Header
	if c.GetHeader("X-Proxy-Mode") == "goflychat" {
		fmt.Printf("DEBUG: 检测到代理模式 Header, 返回前缀: %s\n", GetPrefix())
		return GetPrefix()
	}

	// 检测其他代理特征
	if c.GetHeader("X-Forwarded-For") != "" || c.GetHeader("X-Real-IP") != "" {
		fmt.Printf("DEBUG: 检测到代理特征 Header, 返回前缀: %s\n", GetPrefix())
		return GetPrefix()
	}

	// 检测Referer中是否包含特定域名
	referer := c.GetHeader("Referer")
	if referer != "" && (strings.Contains(referer, "czliehuo.com") ||
		strings.Contains(referer, "localhost") && c.Request.URL.Path != "/") {
		fmt.Printf("DEBUG: 检测到代理 Referer: %s, 返回前缀: %s\n", referer, GetPrefix())
		return GetPrefix()
	}

	// 检测Host是否包含特定域名
	host := c.GetHeader("Host")
	if host != "" && strings.Contains(host, "czliehuo.com") {
		fmt.Printf("DEBUG: 检测到代理 Host: %s, 返回前缀: %s\n", host, GetPrefix())
		return GetPrefix()
	}

	// 直接访问，返回空路径
	fmt.Printf("DEBUG: 直接访问，返回空前缀, Host: %s, Path: %s\n", host, c.Request.URL.Path)
	return ""
}

// GetDynamicStaticPath 获取静态资源路径
func GetDynamicStaticPath(c *gin.Context) string {
	basePath := GetDynamicBasePath(c)
	if basePath != "" {
		return basePath + "/static"
	}
	return "/static"
}

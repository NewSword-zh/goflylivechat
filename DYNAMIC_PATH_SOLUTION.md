# GOFLY 聊天系统动态路径配置解决方案

## 📋 问题背景

GOFLY 聊天系统在部署时遇到了一个经典的嵌套代理架构问题：

### 部署架构
- **A服务器**: Go应用运行在8081端口，没有域名和SSL证书
- **B服务器**: 有域名和SSL证书，配置Nginx代理
- **访问需求**:
  - 嵌入式聊天组件通过B服务器访问
  - 客服后台需要支持直接访问和代理访问两种模式

### 核心问题
1. **嵌套代理路径不匹配**: 代理访问和直接访问需要不同的资源路径
2. **硬编码路径问题**: 静态资源、API接口、头像路径都缺少动态前缀支持
3. **后端配置固定**: 原有系统无法根据访问方式动态调整路径前缀

## 🔧 解决方案概览

### 核心思路
实现**智能路径检测机制**，让Go应用能够：
- **自动识别访问方式**（直接访问 vs 代理访问）
- **动态生成资源路径**（带/不带前缀）
- **支持双重路由配置**（同时支持两种访问模式）

## 🛠️ 具体实施步骤

### 第1步：配置文件系统

#### 1.1 创建应用配置文件
**文件**: `config/app.json`
```json
{
    "app": {
        "prefix": "/goflychat",
        "enable_prefix": true
    }
}
```

#### 1.2 实现配置管理模块
**文件**: `common/config.go`
```go
type AppConfig struct {
    App App `json:"app"`
}

type App struct {
    Prefix       string `json:"prefix"`
    EnablePrefix bool   `json:"enable_prefix"`
}

// 动态路径检测核心函数
func GetDynamicBasePath(c *gin.Context) string {
    // 优先检测代理模式Header
    proxyMode := c.GetHeader("X-Proxy-Proxy-Mode")
    fmt.Printf("DEBUG: X-Proxy-Mode Header: '%s'\n", proxyMode)

    if proxyMode == "goflychat" {
        prefix := GetPrefix()
        fmt.Printf("DEBUG: 检测到代理模式 Header, 返回前缀: %s\n", prefix)
        return prefix
    }

    // 检测其他代理特征
    if c.GetHeader("X-Forwarded-For") != "" || c.GetHeader("X-Real-IP") != "" {
        prefix := GetPrefix()
        fmt.Printf("DEBUG: 检测到代理特征 Header, 返回前缀: %s\n", prefix)
        return prefix
    }

    // 直接访问，返回空路径
    fmt.Printf("DEBUG: 直接访问，返回空前缀, Host: %s, Path: %s\n",
        c.GetHeader("Host"), c.Request.URL.Path)
    return ""
}
```

### 第2步：路由系统重构

#### 2.1 API路由双重支持
**文件**: `router/api.go`
```go
func InitApiRouter(engine *gin.Engine) {
    if common.IsPrefixEnabled() {
        prefix := common.GetPrefix()

        // 带前缀的API路由（代理访问）
        engine.GET(prefix+"/captcha", controller.GetCaptcha)
        engine.POST(prefix+"/check", controller.LoginCheckPass)
        engine.GET(prefix+"/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfo)
        // ... 其他API路由

        // WebSocket路由
        engine.GET(prefix+"/ws_kefu", middleware.JwtApiMiddleware, ws.NewKefuServer)
        engine.GET(prefix+"/ws_visitor", middleware.Ipblack, ws.NewVisitorServer)
    }

    // 无前缀的API路由（直接访问）
    engine.GET("/captcha", controller.GetCaptcha)
    engine.POST("/check", controller.LoginCheckPass)
    engine.GET("/kefuinfo", middleware.JwtApiMiddleware, middleware.RbacAuth, controller.GetKefuInfo)
    // ... 其他API路由
}
```

#### 2.2 页面路由双重支持
**文件**: `router/view.go`
```go
func InitViewRouter(engine *gin.Engine) {
    // 注册带前缀的路由（代理访问）
    if common.IsPrefixEnabled() {
        prefix := common.GetPrefix()
        engine.GET(prefix+"/login", PageLogin)
        engine.GET(prefix+"/main", PageMain)
        engine.GET(prefix+"/chat_main", PageChatMain)
        engine.GET(prefix+"/setting", PageSetting)
    }

    // 注册无前缀的路由（直接访问）
    engine.GET("/login", PageLogin)
    engine.GET("/main", PageMain)
    engine.GET("/chat_main", PageChatMain)
    engine.GET("/setting", PageSetting)
}

// 页面处理函数使用动态路径检测
func PageMain(c *gin.Context) {
    basePath := common.GetDynamicBasePath(c)
    c.HTML(http.StatusOK, "main.html", gin.H{
        "BasePath": basePath,
    })
}
```

#### 2.3 静态资源路由双重支持
**文件**: `cmd/server.go`
```go
if common.IsPrefixEnabled() {
    // 带前缀的静态资源（代理访问）
    staticPrefix := common.GetPrefix() + "/static"
    engine.Static(staticPrefix, "./static")
}

// 无前缀的静态资源（直接访问）
engine.Static("/static", "./static")
```

### 第3步：前端路径处理优化

#### 3.1 简化functions.js路径检测
```javascript
function getBaseUrl() {
    var ishttps = 'https:' == document.location.protocol ? true : false;
    var url = window.location.host;
    if (ishttps) {
        url = 'https://' + url;
    } else {
        url = 'http://' + url;
    }

    // 使用后端传递的基础路径变量
    var basePath = '';
    if (typeof window.APP_BASE_PATH !== 'undefined') {
        basePath = window.APP_BASE_PATH;
    }

    return url + basePath;
}
```

#### 3.2 所有模板初始化路径变量
在每个HTML模板中添加：
```html
<script>
    // 初始化应用基础路径变量
    window.APP_BASE_PATH = "{{.BasePath}}";
</script>
```

### 第4步：动态头像路径处理

#### 4.1 访客头像修复
**文件**: `controller/visitor.go`
```go
func PostVisitorLogin(c *gin.Context) {
    // 动态设置头像路径
    basePath := common.GetDynamicBasePath(c)

    avator := ""
    userAgent := c.GetHeader("User-Agent")
    if tools.IsMobile(userAgent) {
        avator = basePath + "/static/images/1.png"
    } else {
        avator = basePath + "/static/images/2.png"
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
        models.UpdateVisitor(name, avator, id, 1, c.ClientIP(), c.ClientIP(), refer, extra)
    }
    // ... 处理逻辑
}
```

#### 4.2 /notice接口头像修复
**文件**: `controller/notice.go`
```go
func GetNotice(c *gin.Context) {
    kefuId := c.Query("kefu_id")
    user := models.FindUser(kefuId)

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
            "avatar":    avatar,  // 使用动态处理的头像路径
            "nickname":  user.Nickname,
            "allNotice": allNotice.ConfValue,
        },
    })
}
```

#### 4.3 WebSocket消息头像修复

**WebSocket函数签名更新**:
```go
// ws/user.go
func KefuMessage(visitorId, content string, kefuInfo models.User, basePath ...string)

// ws/visitor.go
func VisitorMessage(visitorId, content string, kefuInfo models.User, basePath ...string)
```

**头像路径处理逻辑**:
```go
// 动态处理头像路径
avatar := kefuInfo.Avator
if len(basePath) > 0 && avatar != "" && !strings.HasPrefix(avatar, basePath[0]) {
    avatar = basePath[0] + avatar
}
```

**Controller调用更新**:
```go
// controller/message.go
func SendMessageV2(c *gin.Context) {
    // 动态获取基础路径
    basePath := common.GetDynamicBasePath(c)

    if cType == "kefu" {
        ws.VisitorMessage(vistorInfo.VisitorId, content, kefuInfo, basePath)
    }
    ws.KefuMessage(vistorInfo.VisitorId, content, kefuInfo, basePath)
}
```

## 🔍 B服务器Nginx配置

### 标准代理配置
```nginx
location /goflychat/ {
    proxy_pass http://47.122.20.1:8081;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # 添加代理标识（重要！）
    proxy_set_header X-Proxy-Mode "goflychat";

    # WebSocket支持
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";

    # 其他设置
    proxy_buffering off;
    proxy_max_temp_file_size 0;
    client_max_body_size 100M;
}
```

## 🧪 测试验证

### 1. 直接访问测试
```
http://47.122.20.1:8081/main
```
**预期结果**:
- 控制台日志：`DEBUG: 直接访问，返回空前缀`
- 静态资源：`/static/css/style.css`
- API接口：`/kefuinfo`, `/configs`
- 头像路径：`/static/images/4.jpg`

### 2. 代理访问测试
```
https://www.czliehuo.com/goflylychat/main
```
**预期结果**:
- 控制台日志：`DEBUG: 检测到代理特征 Header, 返回前缀: /goflychat`
- 静态资源：`/goflyflychat/static/css/style.css`
- API接口：`/goflychat/kefuinfo`, `/goflychat/configs`
- 头像路径：`/goflychat/static/images/4.jpg`

### 3. 嵌入式组件测试
```javascript
(function(window, document, scriptUrl, callback) {
    const head = document.getElementsByTagName('head')[0];
    const script = document.createElement('script');
    script.type = 'text/javascript';
    script.src = scriptUrl + "/static/js/chat-widget.js";
    script.onload = script.onreadystatechange = function () {
        if (!this.readyState || this.readyState === "loaded" || this.readyState === "complete") {
            callback(scriptUrl);
        }
    };
    head.appendChild(script);
})(window, document, "/goflychat", function(baseUrl) {
    CHAT_WIDGET.initialize({
        API_URL: baseUrl,
        AGENT_ID: "liujian",
    });
});
```
**预期结果**:
- 聊天组件正常加载
- 所有资源通过代理正常访问
- WebSocket连接：`wss://www.czliehuo.com/goflychat/ws_visitor`

## 📊 技术架构优势

### 1. 统一性
- 所有资源使用相同的路径模式
- 统一的配置管理
- 统一的路径检测逻辑

### 2. 灵活性
- 支持动态切换部署模式
- 配置文件可调整前缀
- 向后兼容旧版本访问

### 3. 可维护性
- 避免硬编码路径
- 集中的错误处理
- 详细的调试日志

### 4. 扩展性
- 易于添加新的路由类型
- 支持多种代理配置
- 支持多层代理架构

## 🎯 解决的核心问题

### ✅ 已解决的问题
1. **静态资源404错误** - 双重路由确保资源可访问
2. **API接口404错误** - 动态路径检测确保API调用正常
3. **头像显示问题** - 多个接口头像路径动态添加前缀
4. **WebSocket消息路径错误** - 实时消息头像正常显示
5. **嵌套代理复杂性** - 自动识别访问方式，简化部署

### ✅ 支持的访问模式
1. **直接访问**: `http://localhost:8081/main`
2. **代理访问**: `https://yourdomain.com/goflychat/main`
3. **嵌入式访问**: 通过JavaScript组件嵌入
4. **混合访问**: 同时支持多种方式

## 🚀 部署建议

### 开发环境
```bash
# 直接启动开发服务器
go run main.go server

# 或使用固定端口
go run main.go server -p 8081
```

### 生产环境
```bash
# 编译生产版本
go build -o gochat

# 启动生产服务
./gochat server

# 后台运行
./gochat server -d
```

### Docker部署
```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go build -o gochat
EXPOSE 8081
CMD ["./gochat", "server"]
```

## 📝 注意事项

### 配置文件管理
- 确保 `config/app.json` 文件存在
- 生产环境注意文件权限
- 配置修改需要重启应用

### Nginx配置要点
- 必须添加 `proxy_set_header X-Proxy-Mode "goflychat"`
- 确保 WebSocket代理配置正确
- 检查路径重写规则

### 数据库迁移
- 现有用户头像可能需要更新
- 可编写脚本批量更新数据库路径
- 新注册用户自动使用动态路径

### 调试技巧
- 使用 `fmt.Printf` 输出调试日志
- 检查浏览器开发者工具的请求头
- 观察A服务器控制台输出

## 🔮 故障排查

### 常见问题及解决方案

#### 1. 静态资源404
**症状**: 浏览器控制台显示静态资源404错误
**原因**: 静态资源路由配置问题
**解决**: 检查 `cmd/server.go` 中的静态路由配置

#### 2. API接口404
**症状**: AJAX请求返回404错误
**原因**: API路由没有正确注册
**解决**: 检查 `router/api.go` 中的路由定义

#### 3. 头像路径错误
**症状**: 头像显示破损或404
**原因**: 头像路径没有动态添加前缀
**解决**: 检查对应Controller中的路径处理逻辑

#### 4. WebSocket连接失败
**症状**: 无法建立WebSocket连接
**原因**: WebSocket路由或代理配置问题
**解决**: 检查WebSocket路由和Nginx代理配置

#### 5. 调试日志不显示
**症状**: 没有看到DEBUG日志输出
**原因**: 可能没有检测到预期的Header
**解决**: 检查Nginx代理配置和请求头传递

## 📈 性能优化建议

### 1. 缓存策略
- 静态资源设置适当的缓存头
- API响应根据类型设置缓存策略
- 启用浏览器缓存

### 2. 数据库优化
- 定期清理过期访客数据
- 为常用查询添加索引
- 优化头像路径存储

### 3. 网络优化
- 启用Gzip压缩
- 配置连接池
- 设置合适的超时时间

这个解决方案提供了一个完整的动态路径配置架构，能够很好地解决嵌套代理环境下的资源访问问题，同时保持了系统的灵活性和可维护性。
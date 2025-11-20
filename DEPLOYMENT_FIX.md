# 🚨 紧急修复部署指南

## 📋 问题总结

经过分析，发现了以下需要修复的问题：

### 问题 1: main.html 静态资源路径错误
- `GET http://47.122.20.1:8081/assets/js/functions.js net::ERR_ABORTED 404`
- **原因**: 路径错误，应该是 `/static/js/functions.js`
- **修复**: 已修复为 `{{.BasePath}}/static/js/functions.js`

### 问题 2: 客服后台 WebSocket 连接错误
- `ws://47.122.20.1:8081/goflychat/ws_kefu`
- **原因**: 直接访问时不应有 `/goflychat` 前缀
- **修复**: 通过 Header 检测，正确返回空 BasePath

### 问题 3: 嵌入聊天组件静态资源 404
- `GET https://czliehuo.com/static/images/2.png 404`
- **原因**: JavaScript 路径检测逻辑未生效
- **修复**: 优化路径检测逻辑，优先使用后端传递的变量

## 🔧 已完成的修复

### ✅ 后端修复
- **router/view.go**: 所有路由函数都支持通过 `X-Proxy-Mode` header 检测
- **模板变量**: 所有主要模板使用 `{{.BasePath}}` 动态路径

### ✅ 前端修复
- **functions.js**: `getBaseUrl()`, `getWsBaseUrl()`, `placeFace()`, `replaceAttachment()` 函数
- **优先级检测**: 后端变量 > 域名检测 > 路径检测

## 🚀 立即部署步骤

### 1. 更新 Nginx 配置
在你的 nginx 配置中添加这一行（如果还没有的话）：
```nginx
location ^~ /goflychat/ {
    # 添加这个 header
    proxy_set_header X-Proxy-Mode "goflychat";

    rewrite ^/goflychat/(.*)$ /$1 break;
    proxy_pass http://47.122.20.1:8081;

    # 其他配置保持不变...
}
```

### 2. 重新编译并重启服务
```bash
# 重新编译
go build -o gochat

# 重启服务
./gochat server
```

## 🧪 测试验证

### 测试场景 1: 代理访问
1. 访问：`https://czliehuo.com/goflychat/livechat?user_id=liujian`
2. 预期结果：
   - ✅ 静态资源路径：`/goflychat/static/images/2.png`
   - ✅ WebSocket 连接：`wss://czliehuo.com/goflychat/ws_visitor`
   - ✅ 所有功能正常

### 测试场景 2: 直接访问客服后台
1. 访问：`http://47.122.20.1:8081/main`
2. 预期结果：
   - ✅ 静态资源：`/static/js/functions.js`
   - ✅ WebSocket 连接：`ws://47.122.20.1:8081/ws_kefu`
   - ✅ 所有功能正常

### 测试场景 3: 直接访问聊天页面
1. 访问：`http://47.122.20.1:8081/livechat?user_id=liujian`
2. 预期结果：
   - ✅ 静态资源：`/static/images/2.png`
   - ✅ WebSocket 连接：`ws://47.122.20.1:8081/ws_visitor`
   - ✅ 所有功能正常

## 🔍 调试方法

### 1. 检查 Nginx Header
访问代理页面时，在浏览器开发者工具的 Network 标签中检查请求头是否包含 `X-Proxy-Mode: goflychat`

### 2. 检查页面源码
查看页面 HTML，确认 `{{.BasePath}}` 被正确替换：
- **代理模式**: `{{.BasePath}}` → `/goflychat`
- **直接模式**: `{{.BasePath}}` → 空字符串

### 3. 控制台调试
在浏览器控制台运行：
```javascript
console.log('Base Path:', window.APP_BASE_PATH);
console.log('WS Base URL:', getWsBaseUrl());
console.log('Face Path:', placeFace()[0]);
console.log('Ext Path:', extBasePath);
```

## 📊 预期效果

修复后，系统将完全支持：

- ✅ **代理访问**: 所有资源正确使用 `/goflychat` 前缀
- ✅ **直接访问**: 所有资源使用相对路径
- ✅ **智能检测**: 自动适应访问环境
- ✅ **无缝切换**: 用户无感知

## 🚨 如果仍有问题

如果测试后仍有问题，请检查：

1. **Nginx 配置**：确认 `proxy_set_header X-Proxy-Mode "goflychat";` 已正确添加
2. **服务重启**：确认 GOFLY 服务已重启并加载新代码
3. **浏览器缓存**：强制刷新浏览器缓存 (Ctrl+F5)

## 📞 紧急联系

如果在部署过程中遇到紧急问题，请提供：
- 错误信息截图
- 访问的完整 URL
- Nginx 配置相关内容
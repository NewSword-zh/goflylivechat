# GOFLY 智能路径检测 - 测试指南

## 🎯 功能说明
现已实施智能路径检测功能，让系统能够自动适应不同的访问方式：
- **代理访问**：通过 nginx `/goflychat/` 代理到 `http://47.122.20.1:8081/`
- **直接访问**：直接访问 `http://47.122.20.1:8081/`

## 🧪 测试场景

### 场景 1：通过 nginx 代理访问
**URL**: `https://www.yoursite.com/goflychat/livechat?user_id=liujian`

**预期行为**:
- ✅ 页面正常加载
- ✅ WebSocket 连接: `wss://www.yoursite.com/goflychat/ws_visitor`
- ✅ API 调用: `/goflychat/visitor_login`, `/goflychat/2/message`, `/goflychat/notice`
- ✅ 静态资源: `/goflychat/static/images/face/`, `/goflychat/static/images/ext/`
- ✅ 文件上传: `/goflychat/uploadimg`, `/goflychat/uploadfile`
- ✅ 音频文件: `/goflychat/static/images/alert2.ogg`

### 场景 2：直接访问服务器
**URL**: `http://47.122.20.1:8081/livechat?user_id=liujian`

**预期行为**:
- ✅ 页面正常加载
- ✅ WebSocket 连接: `ws://47.122.20.1:8081/ws_visitor`
- ✅ API 调用: `/visitor_login`, `/2/message`, `/notice`
- ✅ 静态资源: `/static/images/face/`, `/static/images/ext/`
- ✅ 文件上传: `/uploadimg`, `/uploadfile`
- ✅ 音频文件: `/static/images/alert2.ogg`

## 🔍 测试检查点

### 1. 页面加载测试
- [ ] 页面 HTML 正常渲染
- [ ] CSS 样式正常加载
- [ ] JavaScript 没有错误

### 2. WebSocket 连接测试
- [ ] 浏览器开发者工具 Network 标签显示 WebSocket 连接成功
- [ ] `ws:onopen` 事件正常触发
- [ ] `ws:onmessage` 事件正常接收消息

### 3. API 调用测试
- [ ] 用户登录 API (`/visitor_login`) 调用成功
- [ ] 获取通知 API (`/notice`) 调用成功
- [ ] 消息发送 API (`/2/message`) 调用成功

### 4. 静态资源测试
- [ ] 表情图片正常显示
- [ ] 文件类型图标正常显示
- [ ] 页面图标 (favicon) 正常显示

### 5. 功能测试
- [ ] 表情选择功能正常
- [ ] 图片上传功能正常
- [ ] 文件上传功能正常
- [ ] 消息发送接收正常
- [ ] 音频提示音正常

## 🐛 常见问题排查

### 问题 1: WebSocket 连接失败
**症状**: `WebSocket connection to 'ws://xxx/goflychat/ws_visitor' failed`
**原因**: 直接访问时仍然使用了代理路径
**解决方案**: 检查 `functions.js` 中的 `getWsBaseUrl()` 函数是否正确更新

### 问题 2: 静态资源 404
**症状**: 表情图片或图标显示为 404 错误
**原因**: 路径检测逻辑错误
**解决方案**: 检查 `functions.js` 中的 `placeFace()` 和 `replaceAttachment()` 函数

### 问题 3: API 调用 404
**症状**: AJAX 请求返回 404 错误
**原因**: `window.APP_BASE_PATH` 未正确设置
**解决方案**: 检查 `chat_page.html` 中的路径初始化代码

## 🚀 部署步骤

### 1. 重新编译并部署
```bash
# 重新编译
go build -o gochat

# 重启服务
./gochat server
```

### 2. 验证修改
修改完成后，系统会自动检测访问模式：
- **代理访问**: 如果 URL 包含 `/goflychat/`，所有资源使用 `/goflychat` 前缀
- **直接访问**: 如果 URL 不包含 `/goflychat/`，所有资源使用相对路径

## 🔧 调试方法

### 1. 控制台检查
在浏览器控制台中运行：
```javascript
console.log('Base Path:', window.APP_BASE_PATH);
console.log('WS Base URL:', getWsBaseUrl());
console.log('HTTP Base URL:', getBaseUrl());
```

### 2. 检查页面源码
查看页面HTML，确认静态资源URL：
- **代理模式**: `{{.BasePath}}` 被替换为 `/goflychat`
- **直接模式**: `{{.BasePath}}` 被替换为空字符串

### 3. 网络检查
- 打开浏览器开发者工具 → Network 标签
- 查看请求的 URL 是否正确包含或不包含 `/goflychat` 前缀

### 4. WebSocket 检查
- Network 标签 → 筛选 WebSocket
- 检查连接 URL 和连接状态

## 🎯 核心改进

### 后端改进
- **智能路径检测**: Go后端自动检测请求路径，传递正确的BasePath给模板
- **模板变量**: 使用 `{{.BasePath}}` 动态生成正确的资源路径
- **统一处理**: 所有主要页面都支持两种访问模式

### 前端改进
- **保持兼容**: JavaScript函数保持智能检测功能
- **全局变量**: `window.APP_BASE_PATH` 确保所有脚本使用正确的基础路径
- **无缝切换**: 用户无感知的环境自适应

## 📊 测试结果记录

请测试完成后记录结果：

| 测试场景 | 页面加载 | WebSocket | API 调用 | 静态资源 | 功能测试 | 整体状态 |
|---------|---------|----------|---------|---------|---------|---------|
| 代理访问 |         |          |         |         |         |         |
| 直接访问 |         |          |         |         |         |         |

## 🎉 成功标准
- 两种访问方式下所有功能都正常工作
- 用户无法感知到差异
- 无控制台错误信息
- 性能不受影响
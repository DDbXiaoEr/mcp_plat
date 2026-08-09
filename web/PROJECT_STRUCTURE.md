# 前端项目结构

```
web/
├── AGENTS.md              # 前端开发规范（每次会话开始时优先读取）
├── PROJECT_STRUCTURE.md   # 本文件 —— 前端项目结构速查
├── Makefile               # 前端构建脚本（install / build / dev / preview / clean）
├── index.html             # HTML 入口
├── package.json           # 依赖与脚本
├── package-lock.json
├── vite.config.js         # Vite 配置（dev proxy /api → localhost:8080）
│
├── public/
│   └── favicon.svg
│
├── dist/                  # 构建产物（embed 打包用，不入库）
│   ├── index.html
│   ├── favicon.svg
│   └── assets/
│
└── src/
    ├── main.js            # Vue 入口
    ├── App.vue            # 根组件（按登录状态条件渲染）
    ├── api.js             # HTTP API 封装
    │
    ├── styles/
    │   └── global.css     # 全局 CSS 变量（设计令牌）
    │
    ├── stores/
    │   ├── auth.js        # 鉴权状态（登录/登出/角色）
    │   ├── nav.js         # 侧边栏导航状态（无 Vue Router）
    │   └── settings.js    # 平台设置共享状态（名称/logo/跳转链接）
    │
    └── components/
        ├── LoginView.vue      # 登录页
        ├── TheHeader.vue      # 顶部栏
        ├── TheSidebar.vue     # 侧边栏
        ├── WelcomeView.vue    # 欢迎首页
        ├── ProfileView.vue    # 个人信息（user）
        ├── AccessKeyView.vue  # AccessKey 列表页（user）
        ├── AccessKeyDrawer.vue# AccessKey 编辑抽屉（user）
        ├── HistoryView.vue    # 使用历史页（user）
        ├── OverviewView.vue   # 平台概况（admin，四项统计卡片 + AI 调用趋势折线图）
        ├── ServersView.vue    # MCP 服务器管理（admin，列表+详情+新增对话框）
        └── SettingsView.vue   # 系统设置（admin，占位）
```

## 导航说明

前端不使用 Vue Router，通过 `stores/nav.js` 的响应式 `active` 状态切换组件：
- `App.vue` 使用 `<component :is="...">` 动态渲染
- 导航项按角色定义（admin / user）
- user 菜单：个人信息（profile）、AccessKey 管理（accesskey）、使用历史（history）
- admin 菜单：平台概况（overview）、MCP 服务器管理（servers）、系统设置（settings）

## 常用命令

```bash
npm install      # 安装依赖
npm run dev      # 本地开发，端口 5174
npm run build    # 生产构建（改动后请运行以验证）
npm run preview  # 预览构建产物
```

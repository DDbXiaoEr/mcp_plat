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
    ├── main.js            # Vue 入口（挂载 vue-i18n）
    ├── App.vue            # 根组件（按登录状态条件渲染）
    ├── api.js             # HTTP API 封装
    ├── i18n.js            # vue-i18n 实例（zh-CN/en-US）+ roleLabel 助手
    │
    ├── locales/           # 国际化词条（每视图一个领域文件，由 index.js 合并）
    │   ├── zh-CN/         #   中文（默认，原文为准）
    │   │   ├── index.js       #   合并入口
    │   │   ├── base.js        #   app/common/nav/role/auth 通用词条
    │   │   ├── login.js / welcome.js / profile.js
    │   │   └── overview / history / servers / accesskey / rbac / settings / quickaccess.js
    │   └── en-US/         # 英文（键结构与 zh-CN 完全一致）
    │
    ├── styles/
    │   └── global.css     # 全局 CSS 变量（设计令牌）
    │
    ├── stores/
    │   ├── auth.js        # 鉴权状态（登录/登出/角色 code，不再存中文角色名）
    │   ├── locale.js      # 语言状态（localStorage mcp-console-locale 持久化）
    │   ├── nav.js         # 侧边栏导航状态（无 Vue Router，菜单存 labelKey）
    │   ├── settings.js    # 平台设置共享状态（名称/logo/跳转链接）
    │   └── theme.js       # 亮暗主题状态（localStorage 持久化）
    │
    └── components/
        ├── LoginView.vue      # 登录页（含语言切换）
        ├── TheHeader.vue      # 顶部栏（右侧含语言切换）
        ├── TheSidebar.vue     # 侧边栏
        ├── WelcomeView.vue    # 欢迎首页
        ├── ProfileView.vue    # 个人信息（user）
        ├── AccessKeyView.vue  # AccessKey 列表页（user）
        ├── AccessKeyDrawer.vue# AccessKey 编辑抽屉（user）
        ├── HistoryView.vue    # 使用历史页（user）
        ├── QuickAccessView.vue# 快速接入页（user）
        ├── OverviewView.vue   # 平台概况（admin，四项统计卡片 + AI 调用趋势折线图）
        ├── ServersView.vue    # MCP 服务器管理（admin，列表+详情+新增对话框）
        ├── RbacView.vue       # RBAC 设置（admin，角色/用户管理）
        └── SettingsView.vue   # 系统设置（admin）
```

## 导航说明

前端不使用 Vue Router，通过 `stores/nav.js` 的响应式 `active` 状态切换组件：
- `App.vue` 使用 `<component :is="...">` 动态渲染
- 导航项按角色定义（admin / user），菜单文案使用 `nav.*` 词条键（`labelKey`），渲染时按当前语言取词
- user 菜单：个人信息（profile）、AccessKey 管理（accesskey）、快速接入（quickaccess）、使用历史（history）
- admin 菜单：平台概况（overview）、MCP 服务器管理（servers）、RBAC 设置（rbac）、系统设置（settings）

## 国际化约定

- 新增/修改界面文案时，必须在 `src/locales/zh-CN/<domain>.js` 与 `src/locales/en-US/<domain>.js` **同步维护**同名键（zh 为原文）。
- 通用词（确认/取消/保存/编辑/删除/复制/已复制/加载中等）优先复用 `common.*`；组件文案用该视图领域键（`servers.*`/`settings.*` 等）。
- 模板用 `{{ t('...') }}`，属性绑定 `:placeholder="t('...')"`，`<script setup>` 内 `const { t } = useI18n()`；渲染在列表中的文案优先存键、在渲染点翻译以保证切换语言即时生效。
- 语言默认跟随浏览器（中文优先），用户切换后以 `localStorage`（`mcp-console-locale`）记忆；切换入口在登录页与顶栏。

## 常用命令

```bash
npm install      # 安装依赖
npm run dev      # 本地开发，端口 5174
npm run build    # 生产构建（改动后请运行以验证）
npm run preview  # 预览构建产物
```

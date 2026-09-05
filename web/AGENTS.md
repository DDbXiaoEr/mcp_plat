# AGENTS.md

某某大学 MCP 服务平台 · 管理控制台（mcp-plat-console）。

> **项目结构速查见 `PROJECT_STRUCTURE.md`**（英文版 `PROJECT_STRUCTURE_EN.md`），包含完整目录树、导航说明、常用命令。每次会话开始时优先读取该文件了解项目全貌，无需重新探索整个代码库。

## 技术栈

- Vue 3（`<script setup>` 组合式 API）
- Vite 6
- vue-i18n（zh-CN/en-US 双语，`src/i18n.js` + `src/locales/`）
- 纯 CSS，无 UI 框架；样式写在组件 `<style scoped>` 内，全局变量在 `src/styles/global.css`

## 常用命令

```bash
npm install      # 安装依赖
npm run dev      # 本地开发，端口 5174
npm run build    # 生产构建（改动后请运行以验证）
npm run preview  # 预览构建产物
```

> 项目暂无 lint / 测试脚本；改动后以 `npm run build` 通过为验证标准。

## 目录结构

```
src/
  main.js                 # 入口
  App.vue                 # 根组件，按登录状态条件渲染
  styles/global.css       # 全局样式与 CSS 变量（设计令牌）
  stores/auth.js          # 鉴权状态（临时账号、登录/登出、角色）
  components/
    LoginView.vue         # 登录页
    TheHeader.vue         # 顶部栏
    TheSidebar.vue        # 侧边栏
    WelcomeView.vue       # 登录后主内容
```

## 约定

- 有文件新增或删除时，必须同步更新仓库根目录的 `PROJECT_STRUCTURE.md`（含英文版 `PROJECT_STRUCTURE_EN.md`）及本目录的 `PROJECT_STRUCTURE.md`（含英文版 `PROJECT_STRUCTURE_EN.md`）。
- 组件用 `<script setup>` + 组合式 API；单文件组件命名 PascalCase，布局类组件以 `The` 前缀（如 `TheHeader`）。
- 颜色、圆角、间距等统一用 `global.css` 里的 CSS 变量（`--xauat-blue`、`--radius`、`--header-height` 等），不要硬编码。
- 未完成或待接后端的地方用 `// TODO:` 标注。
- 除非明确要求，不要新增注释。
- 界面文案支持中英双语（vue-i18n）：新增/修改文案必须同步维护 `src/locales/zh-CN/<domain>.js` 与 `en-US/<domain>.js` 的同名键（zh 为原文）；通用词复用 `common.*`，页面文案用各自视图域键。禁止在模板/脚本中硬编码中文 UI 文案。

## 当前状态

- 尚无后端与路由。鉴权为临时实现，见 `src/stores/auth.js` 的 `TEMP_ACCOUNTS`。
- 临时账号：`admin / admin123`（管理员）、`user / user123`（用户）。
- 登录状态持久化在 `localStorage`（键 `mcp-console-auth`）。
- 登录后按角色显示内容，目前仅显示「欢迎，{角色名}登录」。
- 同级项目 `../mcp_plat_portal` 为对外门户站，本项目沿用其代码风格与设计令牌。

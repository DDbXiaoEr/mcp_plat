# AGENTS_EN.md

The MCP Service Platform Admin Console of a certain university (mcp-plat-console).

> **See `PROJECT_STRUCTURE.md` (Chinese) / `PROJECT_STRUCTURE_EN.md` (English) for the frontend structure quick reference** — full directory tree, navigation notes and common commands. Read them first at the start of each session to understand the whole project without re-exploring the codebase.

## Tech Stack

- Vue 3 (`<script setup>` Composition API)
- Vite 6
- vue-i18n (zh-CN/en-US bilingual, `src/i18n.js` + `src/locales/`)
- Plain CSS, no UI framework; styles live in each component's `<style scoped>`, global variables in `src/styles/global.css`

## Common Commands

```bash
npm install      # install dependencies
npm run dev      # local dev, port 5174
npm run build    # production build (run to verify after changes)
npm run preview  # preview the build output
```

> There is no lint/test script yet; passing `npm run build` is the verification standard after changes.

## Directory Structure

```
src/
  main.js                 # entry
  App.vue                 # root component, conditionally rendered by login state
  styles/global.css       # global styles and CSS variables (design tokens)
  stores/auth.js          # auth state (temp accounts, login/logout, role)
  components/
    LoginView.vue         # login page
    TheHeader.vue         # top bar
    TheSidebar.vue        # sidebar
    WelcomeView.vue       # main content after login
```

## Conventions

- When files are added or removed, you must sync `PROJECT_STRUCTURE.md` / `PROJECT_STRUCTURE_EN.md` at the repository root (and the ones in this directory).
- Components use `<script setup>` + Composition API; single-file components are named in PascalCase, layout components prefixed with `The` (e.g. `TheHeader`).
- Colors, radii, spacing, etc. must use the CSS variables from `global.css` (`--xauat-blue`, `--radius`, `--header-height`, …); do not hard-code values.
- Mark unfinished or not-yet-connected backend spots with `// TODO:`.
- Do not add comments unless explicitly requested.
- UI copy supports Chinese and English (vue-i18n): new/changed copy must maintain matching keys in `src/locales/zh-CN/<domain>.js` and `en-US/<domain>.js` (zh is the source). Reuse `common.*` for shared terms and per-view keys for page copy. Never hard-code Chinese UI text in templates or scripts.

## Current Status

- No backend or router yet. Auth is a temporary implementation, see `TEMP_ACCOUNTS` in `src/stores/auth.js`.
- Temporary accounts: `admin / admin123` (admin), `user / user123` (user).
- Login state is persisted in `localStorage` (key `mcp-console-auth`).
- After login, content is shown by role; currently only shows "Welcome, logged in as {role}".
- The sibling project `../mcp_plat_portal` is the public portal site; this project follows its code style and design tokens.

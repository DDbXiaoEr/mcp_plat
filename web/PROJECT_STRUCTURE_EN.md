# Frontend Project Structure (English)

```
web/
├── AGENTS.md                 # frontend development conventions (Chinese) — read first at the start of every session
├── AGENTS_EN.md              # frontend development conventions (English) — read first at the start of every session
├── PROJECT_STRUCTURE.md      # this file (Chinese) — frontend project structure quick reference
├── PROJECT_STRUCTURE_EN.md   # this file (English) — frontend project structure quick reference
├── Makefile                  # frontend build scripts (install / build / dev / preview / clean)
├── index.html                # HTML entry
├── package.json              # dependencies and scripts
├── package-lock.json
├── vite.config.js            # Vite config (dev proxy /api → localhost:8080)
│
├── public/
│   └── favicon.svg
│
├── dist/                     # build output (for embed packaging, not committed)
│   ├── index.html
│   ├── favicon.svg
│   └── assets/
│
└── src/
    ├── main.js               # Vue entry (mounts vue-i18n)
    ├── App.vue               # root component (conditionally rendered by login state)
    ├── api.js                # HTTP API wrapper
    ├── i18n.js               # vue-i18n instance (zh-CN/en-US) + roleLabel helper
    │
    ├── locales/              # i18n entries (one domain file per view, merged by index.js)
    │   ├── zh-CN/            #   Chinese (default, source of truth)
    │   │   ├── index.js      #   merge entry
    │   │   ├── base.js       #   shared entries for app/common/nav/role/auth
    │   │   ├── login.js / welcome.js / profile.js
    │   │   └── overview / history / servers / accesskey / rbac / settings / quickaccess.js
    │   └── en-US/            # English (key structure identical to zh-CN)
    │
    ├── styles/
    │   └── global.css        # global CSS variables (design tokens)
    │
    ├── stores/
    │   ├── auth.js           # auth state (login/logout/role code, no longer stores Chinese role names)
    │   ├── locale.js         # language state (persisted in localStorage mcp-console-locale)
    │   ├── nav.js            # sidebar nav state (no Vue Router; menu stores labelKey)
    │   ├── settings.js       # platform settings shared state (name/logo/jump links)
    │   └── theme.js          # light/dark theme state (persisted in localStorage)
    │
    └── components/
        ├── LoginView.vue         # login page (with language switcher)
        ├── TheHeader.vue         # top bar (language switcher on the right)
        ├── TheSidebar.vue        # sidebar
        ├── WelcomeView.vue       # welcome home page
        ├── ProfileView.vue       # personal info (user)
        ├── AccessKeyView.vue     # AccessKey list page (user)
        ├── AccessKeyDrawer.vue   # AccessKey edit drawer (user)
        ├── HistoryView.vue       # usage history page (user)
        ├── QuickAccessView.vue   # quick access page (user)
        ├── OverviewView.vue      # platform overview (admin; 4 stat cards + AI call-trend chart)
        ├── ServersView.vue       # MCP server management (admin; list + detail + create dialog)
        ├── RbacView.vue          # RBAC settings (admin; role/user management)
        └── SettingsView.vue      # system settings (admin)
```

## Navigation Notes

The frontend does not use Vue Router; it switches components through the reactive `active` state in `stores/nav.js`:
- `App.vue` renders dynamically with `<component :is="...">`
- Nav items are defined per role (admin / user); menu copy uses `nav.*` entry keys (`labelKey`) and is looked up in the current language when rendered
- user menu: Personal Info (profile), AccessKey (accesskey), Quick Access (quickaccess), Usage History (history)
- admin menu: Platform Overview (overview), MCP Servers (servers), RBAC (rbac), System Settings (settings)

## i18n Conventions

- When adding/changing UI copy you must **keep the same keys in sync** between `src/locales/zh-CN/<domain>.js` and `src/locales/en-US/<domain>.js` (zh is the source).
- Reuse `common.*` for shared terms (confirm/cancel/save/edit/delete/copy/copied/loading…); component copy uses the per-view domain keys (`servers.*`/`settings.*` etc.).
- In templates use `{{ t('...') }}`, for attribute binding `:placeholder="t('...')"`, and inside `<script setup>` `const { t } = useI18n()`. Store the key for copy rendered in lists and translate at the render point so the language switches take effect immediately.
- Language follows the browser by default (Chinese preferred); after the user switches, it is remembered via `localStorage` (`mcp-console-locale`). The switcher is on the login page and in the top bar.

## Common Commands

```bash
npm install      # install dependencies
npm run dev      # local dev, port 5174
npm run build    # production build (run to verify after changes)
npm run preview  # preview the build output
```

import { reactive, readonly } from 'vue'

// TODO: 后端/路由接入后，菜单可改为配置化或路由驱动
const MENUS = {
  user: [
    { key: 'profile', label: '个人信息' },
    { key: 'accesskey', label: 'AccessKey 管理' },
    { key: 'history', label: '使用历史' }
  ],
  admin: [
    { key: 'overview', label: '平台概况' },
    { key: 'servers', label: 'MCP 服务器管理' },
    { key: 'rbac', label: 'RBAC设置' },
    { key: 'settings', label: '系统设置' }
  ]
}

const state = reactive({
  active: ''
})

function menusFor(role) {
  return MENUS[role] || []
}

function setActive(key) {
  state.active = key
}

function resetNav() {
  state.active = ''
}

export const nav = readonly(state)
export { menusFor, setActive, resetNav }

/*
 * Copyright (C) 2026 Zhaoquan Wang
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Author: deepseek-v4-pro / opencode
import { reactive, readonly } from 'vue'

// TODO: 后端/路由接入后，菜单可改为配置化或路由驱动
const MENUS = {
  user: [
    { key: 'profile', label: '个人信息' },
    { key: 'accesskey', label: 'AccessKey 管理' },
    { key: 'quickaccess', label: '快速接入' },
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

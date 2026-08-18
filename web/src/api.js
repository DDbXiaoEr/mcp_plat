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
import { getToken } from './stores/auth.js'

const BASE = '/api'

async function request(path, options = {}) {
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  }
  if (getToken()) {
    headers.Authorization = `Bearer ${getToken()}`
  }
  const res = await fetch(`${BASE}${path}`, { ...options, headers })
  let json
  try {
    json = await res.json()
  } catch {
    const text = await res.text()
    throw new Error(`请求失败 (${res.status}): ${text || res.statusText}`)
  }
  if (json.code !== 200) {
    throw new Error(json.message || '请求失败')
  }
  return json.data
}

export function fetchAccessKeys() {
  return request('/access-keys')
}

export function createAccessKey(body) {
  return request('/access-keys', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function updateAccessKey(id, body) {
  return request(`/access-keys/${id}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function deleteAccessKey(id) {
  return request(`/access-keys/${id}`, { method: 'DELETE' })
}

export function fetchServers() {
  return request('/servers')
}

export function createServer(body) {
  return request('/servers', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function updateServer(id, body) {
  return request(`/servers/${id}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function deleteServer(id) {
  return request(`/servers/${id}`, { method: 'DELETE' })
}

export function fetchServerTools(body) {
  return request('/servers/fetch-tools', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function fetchRoles() {
  return request('/roles')
}

export function createRole(body) {
  return request('/roles', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function updateRole(id, body) {
  return request(`/roles/${id}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function deleteRole(id) {
  return request(`/roles/${id}`, { method: 'DELETE' })
}

export function fetchRoleUsers(id) {
  return request(`/roles/${id}/users`)
}

export function assignRoleUsers(id, userIds) {
  return request(`/roles/${id}/users`, {
    method: 'PUT',
    body: JSON.stringify({ user_ids: userIds })
  })
}

export function fetchUsers(roleId) {
  const params = roleId ? `?role_id=${roleId}` : ''
  return request(`/users${params}`)
}

export function createUser(body) {
  return request('/users', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

export function updateUser(id, body) {
  return request(`/users/${id}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function deleteUser(id) {
  return request(`/users/${id}`, { method: 'DELETE' })
}

export function fetchSettings() {
  return request('/settings')
}

export function fetchSetting(key) {
  return request(`/settings/${key}`)
}

export function fetchGatewayStatus() {
  return request('/settings/gateway-status')
}

export function fetchOverviewStats() {
  return request('/overview/stats')
}

export function fetchOverviewCallTrend(days = 30) {
  return request(`/overview/call-trend?days=${days}`)
}

export function publishServers(serverIds, enableAuth, accesskeyHeader, enableAuditLog) {
  return request('/servers/publish', {
    method: 'POST',
    body: JSON.stringify({ server_ids: serverIds, enable_auth: enableAuth, accesskey_header: accesskeyHeader, enable_audit_log: enableAuditLog })
  })
}

export function setServersMaintenance(serverIds, restoreIds) {
  return request('/servers/maintenance', {
    method: 'POST',
    body: JSON.stringify({ server_ids: serverIds, restore_ids: restoreIds })
  })
}

export function saveSetting(key, body) {
  return request(`/settings/${key}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

export function fetchHistory(params = {}) {
  const query = new URLSearchParams()
  Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== null && v !== '') query.set(k, v)
  })
  return request(`/history?${query.toString()}`)
}

export function testLdap(body) {
  return request('/settings/test-ldap', {
    method: 'POST',
    body: JSON.stringify(body)
  })
}

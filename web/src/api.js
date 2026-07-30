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
  const json = await res.json()
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

export function fetchGatewayStatus() {
  return request('/settings/gateway-status')
}

export function publishServers(serverIds, enableAuth, accesskeyHeader) {
  return request('/servers/publish', {
    method: 'POST',
    body: JSON.stringify({ server_ids: serverIds, enable_auth: enableAuth, accesskey_header: accesskeyHeader })
  })
}

export function saveSetting(key, body) {
  return request(`/settings/${key}`, {
    method: 'PUT',
    body: JSON.stringify(body)
  })
}

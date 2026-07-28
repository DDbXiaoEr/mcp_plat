// Author: deepseek-v4-pro / opencode
import { reactive, readonly } from 'vue'

const ROLE_LABELS = {
  admin: '管理员',
  user: '用户'
}

const STORAGE_KEY = 'mcp-console-auth'

function loadUser() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

const state = reactive({
  user: loadUser()
})

function roleLabel(role) {
  return ROLE_LABELS[role] || '未知角色'
}

export function getToken() {
  return state.user?.token || ''
}

export async function login(username, password) {
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ username, password })
    })
    const json = await res.json()
    if (json.code !== 200) {
      return { ok: false, message: json.message || '登录失败' }
    }
    saveLogin(json.data)
    return { ok: true }
  } catch {
    return { ok: false, message: '网络错误，请稍后重试' }
  }
}

export async function casLogin(ticket, serviceUrl) {
  try {
    const res = await fetch('/api/auth/cas/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ticket, serviceUrl })
    })
    const json = await res.json()
    if (json.code !== 200) {
      return { ok: false, message: json.message || 'CAS 登录失败' }
    }
    saveLogin(json.data)
    return { ok: true }
  } catch {
    return { ok: false, message: '网络错误，请稍后重试' }
  }
}

export async function fetchAuthMethod() {
  try {
    const res = await fetch('/api/auth/method')
    const json = await res.json()
    if (json.code !== 200) {
      return { method: 'local' }
    }
    return json.data
  } catch {
    return { method: 'local' }
  }
}

function saveLogin(data) {
  const user = {
    username: data.username,
    role: data.role,
    roleLabel: roleLabel(data.role),
    token: data.token
  }
  state.user = user
  localStorage.setItem(STORAGE_KEY, JSON.stringify(user))
}

export function logout() {
  state.user = null
  localStorage.removeItem(STORAGE_KEY)
}

export async function fetchProfile() {
  try {
    const res = await fetch('/api/auth/profile', {
      headers: { Authorization: `Bearer ${getToken()}` }
    })
    const json = await res.json()
    if (json.code !== 200) {
      return null
    }
    return json.data
  } catch {
    return null
  }
}

export const auth = readonly(state)

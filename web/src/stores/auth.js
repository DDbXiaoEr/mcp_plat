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
    const user = {
      username: json.data.username,
      role: json.data.role,
      roleLabel: roleLabel(json.data.role),
      token: json.data.token
    }
    state.user = user
    localStorage.setItem(STORAGE_KEY, JSON.stringify(user))
    return { ok: true }
  } catch {
    return { ok: false, message: '网络错误，请稍后重试' }
  }
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

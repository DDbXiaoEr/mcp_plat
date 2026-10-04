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
import { i18n } from '../i18n.js'

const state = reactive({
  user: null,
  ready: false
})

function t(key) {
  return i18n.global.t(key)
}

export async function login(username, password) {
  try {
    const res = await fetch('/api/auth/login', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ username, password })
    })
    const json = await res.json()
    if (json.code !== 200) {
      return { ok: false, message: json.message || t('auth.loginFailed') }
    }
    saveLogin(json.data)
    return { ok: true }
  } catch {
    return { ok: false, message: t('auth.networkError') }
  }
}

export async function casLogin(ticket, serviceUrl) {
  try {
    const res = await fetch('/api/auth/cas/validate', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials: 'include',
      body: JSON.stringify({ ticket, serviceUrl })
    })
    const json = await res.json()
    if (json.code !== 200) {
      return { ok: false, message: json.message || t('auth.casFailed') }
    }
    saveLogin(json.data)
    return { ok: true }
  } catch {
    return { ok: false, message: t('auth.networkError') }
  }
}

export async function fetchAuthMethod() {
  try {
    const res = await fetch('/api/auth/method', { credentials: 'include' })
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
  state.user = {
    username: data.username,
    name: data.name || data.username,
    role: data.role
  }
  state.ready = true
}

export async function logout() {
  state.user = null
  try {
    await fetch('/api/auth/logout', { method: 'POST', credentials: 'include' })
  } catch {
    // ignore
  }
}

export async function fetchProfile() {
  try {
    const res = await fetch('/api/auth/profile', { credentials: 'include' })
    const json = await res.json()
    if (json.code !== 200) {
      return null
    }
    return json.data
  } catch {
    return null
  }
}

export async function restoreSession() {
  const data = await fetchProfile()
  if (data && data.username) {
    saveLogin(data)
  } else {
    state.user = null
  }
  state.ready = true
}

export const auth = readonly(state)

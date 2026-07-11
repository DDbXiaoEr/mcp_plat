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

<!--
  Copyright (C) 2026 Zhaoquan Wang

  This program is free software: you can redistribute it and/or modify
  it under the terms of the GNU Affero General Public License as published
  by the Free Software Foundation, either version 3 of the License, or
  (at your option) any later version.

  This program is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU Affero General Public License for more details.

  You should have received a copy of the GNU Affero General Public License
  along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<script setup>

// Author: deepseek-v4-pro / opencode
import { ref, onMounted } from 'vue'
import { login, fetchAuthMethod } from '../stores/auth.js'

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

onMounted(async () => {
  const authMethod = await fetchAuthMethod()
  if (authMethod.method === 'cas' && authMethod.cas?.serverUrl) {
    const backUrl = window.location.origin + window.location.pathname
    const loginUrl = authMethod.cas.serverUrl.replace(/\/$/, '') +
      '/login?service=' + encodeURIComponent(backUrl)
    window.location.href = loginUrl
    return
  }
})

async function onSubmit() {
  error.value = ''
  if (!username.value || !password.value) {
    error.value = '请输入账号和密码'
    return
  }
  loading.value = true
  const result = await login(username.value.trim(), password.value)
  loading.value = false
  if (!result.ok) {
    error.value = result.message
  }
}
</script>

<template>
  <div class="login">
    <form class="login__card" @submit.prevent="onSubmit">
      <div class="login__brand">
        <span class="brand__mark" aria-hidden="true">MCP</span>
        <div class="brand__text">
          <strong>某某大学 MCP 服务平台</strong>
          <em>管理控制台</em>
        </div>
      </div>

      <label class="field">
        <span>账号</span>
        <input
          v-model="username"
          type="text"
          autocomplete="username"
          placeholder="请输入账号"
        />
      </label>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="请输入密码"
        />
      </label>

      <p v-if="error" class="login__error">{{ error }}</p>

      <button type="submit" class="login__submit" :disabled="loading">
        {{ loading ? '登录中…' : '登录' }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login {
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background: linear-gradient(135deg, #eef3fb, #f6f8fb);
}

.login__card {
  width: 100%;
  max-width: 380px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 36px 32px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.login__brand {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 6px;
}

.brand__mark {
  display: grid;
  place-items: center;
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--xauat-blue), var(--xauat-cyan));
  color: #fff;
  font-weight: 800;
  font-size: 14px;
}

.brand__text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}

.brand__text strong {
  font-size: 16px;
  color: var(--xauat-blue);
}

.brand__text em {
  font-style: normal;
  font-size: 13px;
  color: var(--text-muted);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 14px;
  color: var(--text-muted);
}

.field input {
  padding: 11px 14px;
  font-size: 15px;
  color: var(--text);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.field input:focus {
  border-color: var(--xauat-blue-light);
  box-shadow: 0 0 0 3px rgba(30, 95, 176, 0.12);
}

.login__error {
  margin: -4px 0 0;
  font-size: 13px;
  color: #d64545;
}

.login__submit {
  padding: 12px;
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  background: var(--xauat-blue);
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s, transform 0.15s;
}

.login__submit:hover:not(:disabled) {
  background: var(--xauat-blue-light);
  transform: translateY(-1px);
}

.login__submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>

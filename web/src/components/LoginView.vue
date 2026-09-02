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
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { login, fetchAuthMethod } from '../stores/auth.js'
import { platformSettings, loadPublicPlatform } from '../stores/settings.js'
import { locale, setLocale } from '../stores/locale.js'

const { t } = useI18n()

const currentLocale = computed(() => locale.locale)

function onLocaleChange(event) {
  setLocale(event.target.value)
}

const username = ref('')
const password = ref('')
const error = ref('')
const loading = ref(false)

const bgStyle = computed(() => {
  const url = platformSettings.platform.loginBackground
  if (!url) return {}
  return {
    backgroundImage: `url(${url})`,
    backgroundSize: 'cover',
    backgroundPosition: 'center',
    backgroundRepeat: 'no-repeat'
  }
})

onMounted(async () => {
  loadPublicPlatform()
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
    error.value = t('login.emptyHint')
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
  <div class="login" :style="bgStyle">
    <select class="login__lang" name="lang" :value="currentLocale" :aria-label="t('app.langLabel')" @change="onLocaleChange">
      <option value="zh-CN">中文</option>
      <option value="en-US">English</option>
    </select>
    <form class="login__card" @submit.prevent="onSubmit">
      <div class="login__brand">
        <span class="brand__mark" aria-hidden="true">MCP</span>
        <div class="brand__text">
          <strong>{{ t('login.brand') }}</strong>
          <em>{{ t('login.tagline') }}</em>
        </div>
      </div>

      <label class="field">
        <span>{{ t('login.username') }}</span>
        <input
          v-model="username"
          name="username"
          type="text"
          autocomplete="username"
          :placeholder="t('login.usernamePlaceholder')"
        />
      </label>

      <label class="field">
        <span>{{ t('login.password') }}</span>
        <input
          v-model="password"
          name="password"
          type="password"
          autocomplete="current-password"
          :placeholder="t('login.passwordPlaceholder')"
        />
      </label>

      <p v-if="error" class="login__error">{{ error }}</p>

      <button type="submit" class="login__submit" :disabled="loading">
        {{ loading ? t('login.submitting') : t('login.submit') }}
      </button>
    </form>
  </div>
</template>

<style scoped>
.login {
  position: relative;
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background: linear-gradient(135deg, #eef3fb, #f6f8fb);
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
}

.login__lang {
  position: absolute;
  top: 20px;
  right: 24px;
  padding: 5px 8px;
  font-size: 13px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  outline: none;
  cursor: pointer;
}

.login__card {
  width: 100%;
  max-width: 380px;
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 36px 32px;
  background: var(--glass-bg);
  border: 1px solid var(--glass-border);
  border-radius: var(--radius);
  box-shadow: 0 8px 30px rgba(10, 61, 122, 0.15);
  backdrop-filter: blur(18px) saturate(160%);
  -webkit-backdrop-filter: blur(18px) saturate(160%);
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

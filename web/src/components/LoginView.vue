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
    <div class="login__stage">
      <aside class="login__intro">
        <span class="login__mark" aria-hidden="true">MCP</span>
        <h1>{{ t('login.brand') }}</h1>
        <p>{{ t('login.stageLead') }}</p>
      </aside>
      <form class="login__card" @submit.prevent="onSubmit">
        <div class="login__brand">
          <strong>{{ t('login.tagline') }}</strong>
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
  </div>
</template>

<style scoped>
.login {
  position: relative;
  min-height: 100vh;
  display: grid;
  place-items: center;
  padding: 24px;
  background-color: transparent;
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
  border-radius: 6px;
  outline: none;
  cursor: pointer;
}

.login__stage {
  width: 100%;
  max-width: 760px;
  display: grid;
  grid-template-columns: 1.05fr 0.95fr;
  background: var(--surface);
  border: 1px solid var(--border);
}

.login__intro {
  padding: 40px 36px;
  background: var(--xauat-navy);
  color: #fff;
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  gap: 12px;
  min-height: 360px;
}

.login__mark {
  width: 44px;
  height: 44px;
  display: grid;
  place-items: center;
  margin-bottom: auto;
  border: 1px solid rgba(255, 255, 255, 0.28);
  color: var(--xauat-accent);
  font-weight: 800;
  font-size: 13px;
  letter-spacing: 0.04em;
}

.login__intro h1 {
  font-size: 26px;
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.3;
}

.login__intro p {
  font-size: 14px;
  line-height: 1.55;
  color: rgba(255, 255, 255, 0.78);
}

.login__card {
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 40px 32px;
  border-left: 3px solid var(--xauat-accent);
}

.login__brand strong {
  font-size: 16px;
  color: var(--xauat-blue);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 14px;
  color: var(--text-muted);
}

.field input {
  padding: 11px 12px;
  font-size: 15px;
  color: var(--text);
  background: var(--surface);
  border: none;
  border-bottom: 1px solid var(--border);
  border-radius: 0;
  outline: none;
  transition: border-color 0.15s;
}

.field input:focus {
  border-bottom-color: var(--xauat-blue);
}

.login__error {
  margin: -4px 0 0;
  font-size: 13px;
  color: var(--danger);
}

.login__submit {
  margin-top: 4px;
  padding: 12px;
  font-size: 15px;
  font-weight: 600;
  color: #fff;
  background: var(--xauat-blue);
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s;
}

.login__submit:hover:not(:disabled) {
  background: var(--xauat-blue-light);
}

.login__submit:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

@media (max-width: 720px) {
  .login__stage {
    grid-template-columns: 1fr;
    max-width: 400px;
  }

  .login__intro {
    min-height: 0;
    padding: 28px 28px 24px;
  }

  .login__mark {
    margin-bottom: 12px;
  }

  .login__card {
    border-left: none;
    border-top: 3px solid var(--xauat-accent);
    padding: 28px 24px 32px;
  }
}
</style>

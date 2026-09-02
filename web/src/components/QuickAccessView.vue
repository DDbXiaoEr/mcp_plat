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
import { ref, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchServers, fetchAccessKeys, fetchGatewayStatus, fetchSetting } from '../api.js'
import { setActive } from '../stores/nav.js'

const { t } = useI18n()

const CLIENTS = {
  cherrystudio: {
    key: 'cherrystudio',
    name: 'CherryStudio',
    descriptionKey: 'quickaccess.clients.cherrystudio.description',
    stepKeys: [
      'quickaccess.clients.cherrystudio.step1',
      'quickaccess.clients.cherrystudio.step2',
      'quickaccess.clients.cherrystudio.step3',
      'quickaccess.clients.cherrystudio.step4'
    ]
  }
}

const servers = ref([])
const keys = ref([])
const gateway = ref({})
const quickAccess = ref({ enabledClients: ['cherrystudio'], scheme: 'https' })
const loading = ref(true)
const selectedKeyId = ref('')
const selectedServerId = ref('')
const copied = ref(false)

onMounted(async () => {
  try {
    const [srv, ak, gw, qa] = await Promise.all([
      fetchServers(),
      fetchAccessKeys(),
      fetchGatewayStatus(),
      fetchSetting('quick_access').catch(() => null)
    ])
    servers.value = (srv || []).filter((s) => s.status === 'published')
    keys.value = (ak || []).filter((k) => k.enabled && !k.is_expired)
    gateway.value = gw || {}
    if (qa) quickAccess.value = qa
    if (keys.value.length) selectedKeyId.value = String(keys.value[0].id)
    const allowed = serverOptions.value
    if (allowed.length) selectedServerId.value = String(allowed[0].id)
  } catch {
    // ignore load error
  }
  loading.value = false
})

const enabledClients = computed(() =>
  Object.values(CLIENTS).filter((c) =>
    (quickAccess.value.enabledClients || []).includes(c.key)
  )
)

const client = computed(() => enabledClients.value[0] || null)

function parseKeyServers(raw) {
  if (!raw) return {}
  try {
    const obj = JSON.parse(raw)
    if (obj && typeof obj === 'object' && !Array.isArray(obj)) return obj
  } catch {
    // ignore
  }
  return {}
}

function authRequired(server) {
  return server.auth_enabled !== false
}

const openServers = computed(() => servers.value.filter((s) => !authRequired(s)))

const serverOptions = computed(() => {
  if (!keys.value.length) return openServers.value
  const key = keys.value.find((k) => String(k.id) === String(selectedKeyId.value))
  const perm = key ? parseKeyServers(key.servers) : {}
  const ids = Object.keys(perm)
  return servers.value.filter((s) => {
    if (!authRequired(s)) return true
    if (!key) return false
    if (ids.length) return ids.includes(String(s.id))
    return true
  })
})

watch(selectedKeyId, () => {
  const allowed = serverOptions.value
  if (!allowed.some((s) => String(s.id) === String(selectedServerId.value))) {
    selectedServerId.value = allowed.length ? String(allowed[0].id) : ''
  }
})

const selectedKey = computed(() =>
  keys.value.find((k) => String(k.id) === String(selectedKeyId.value))
)

const selectedServer = computed(() =>
  servers.value.find((s) => String(s.id) === String(selectedServerId.value))
)

const selectedServerIsOpen = computed(() =>
  selectedServer.value ? !authRequired(selectedServer.value) : false
)

const accesskeyHeader = computed(() =>
  (gateway.value.accesskeyHeader || '').trim() || 'X-Access-Key'
)

const endpointURL = computed(() => {
  const server = selectedServer.value
  if (!server) return ''
  const addr = (server.address || '').trim()
  if (/^https?:\/\//i.test(addr)) return addr
  let base = (gateway.value.defaultPublishDomain || '').trim()
  if (!base) return ''
  const scheme = quickAccess.value.scheme === 'http' ? 'http' : 'https'
  if (!/^https?:\/\//i.test(base)) base = scheme + '://' + base
  base = base.replace(/\/+$/, '')
  let path
  if (gateway.value.provider === 'kong') {
    path = addr ? (addr.startsWith('/') ? addr : '/' + addr) : '/' + server.id
  } else {
    path = '/' + server.id
  }
  return base + path
})

const mcpInstallType = computed(() =>
  (selectedServer.value?.protocol || '').toLowerCase().includes('sse') ? 'sse' : 'streamableHttp'
)

const configJSON = computed(() => {
  const server = selectedServer.value
  if (!server || !endpointURL.value) return ''
  const type = (server.protocol || '').toLowerCase().includes('sse') ? 'sse' : 'http'
  const entry = {
    type,
    url: endpointURL.value
  }
  if (!selectedServerIsOpen.value) {
    const key = selectedKey.value
    if (!key) return ''
    entry.headers = {
      [accesskeyHeader.value]: key.key
    }
  }
  const config = {
    mcpServers: {
      [server.name]: entry
    }
  }
  return JSON.stringify(config, null, 2)
})

function base64EncodeUtf8(text) {
  const bytes = new TextEncoder().encode(text)
  let bin = ''
  bytes.forEach((b) => { bin += String.fromCharCode(b) })
  return btoa(bin)
}

const cherryInstallLink = computed(() => {
  const server = selectedServer.value
  if (!client.value || client.value.key !== 'cherrystudio') return ''
  if (!server || !endpointURL.value) return ''
  const payload = {
    name: server.name,
    type: mcpInstallType.value,
    baseUrl: endpointURL.value
  }
  if (!selectedServerIsOpen.value) {
    const key = selectedKey.value
    if (!key) return ''
    payload.headers = {
      [accesskeyHeader.value]: key.key
    }
  }
  return 'cherrystudio://mcp/install?servers=' + base64EncodeUtf8(JSON.stringify(payload))
})

async function copyConfig() {
  if (!configJSON.value) return
  try {
    await navigator.clipboard.writeText(configJSON.value)
    copied.value = true
    setTimeout(() => { copied.value = false }, 1500)
  } catch {
    // ignore
  }
}

function maskKey(value) {
  if (!value) return ''
  if (value.length <= 8) return value
  return value.slice(0, 3) + '****' + value.slice(-4)
}

function goAccessKey() {
  setActive('accesskey')
}
</script>

<template>
  <section class="page">
    <h1 class="page__title">{{ t('quickaccess.title') }}</h1>
    <p class="page__desc">
      {{ t('quickaccess.description') }}
    </p>

    <div v-if="loading" class="qa__loading">{{ t('common.loading') }}</div>

    <template v-else>
      <div v-if="!gateway.configured" class="qa__notice qa__notice--warn">
        {{ t('quickaccess.gatewayNotConfigured') }}
      </div>

      <div v-else-if="enabledClients.length === 0" class="qa__notice qa__notice--warn">
        {{ t('quickaccess.noClients') }}
      </div>

      <div v-else-if="servers.length === 0" class="qa__notice">
        {{ t('quickaccess.noServers') }}
      </div>

      <div v-else-if="keys.length === 0 && openServers.length === 0" class="qa__notice">
        {{ t('quickaccess.noAccessKey') }}
        <button class="qa__link" type="button" @click="goAccessKey">{{ t('quickaccess.createKeyAction') }}</button>
      </div>

      <template v-else>
        <div class="qa__card">
          <h2 class="qa__card-title">{{ t('quickaccess.selectClient') }}</h2>
          <div class="qa__clients">
            <div
              v-for="c in enabledClients"
              :key="c.key"
              class="qa__client"
              :class="{ 'qa__client--active': client && c.key === client.key }"
            >
              <span class="qa__client-name">{{ c.name }}</span>
              <span class="qa__client-desc">{{ t(c.descriptionKey) }}</span>
            </div>
          </div>
        </div>

        <div class="qa__card">
          <h2 class="qa__card-title">{{ t('quickaccess.selectServerAndKey') }}</h2>
          <div class="qa__form">
            <label class="qa__field">
              <span class="qa__label">{{ t('quickaccess.serverLabel') }}</span>
              <select v-model="selectedServerId" class="qa__select">
                <option
                  v-for="s in serverOptions"
                  :key="s.id"
                  :value="s.id"
                >
                  {{ s.name }}（{{ s.protocol }}）
                </option>
              </select>
            </label>
            <label class="qa__field">
              <span class="qa__label">{{ t('quickaccess.accessKeyLabel') }}</span>
              <template v-if="selectedServerIsOpen">
                <span class="qa__note">{{ t('quickaccess.openServerNote') }}</span>
              </template>
              <select v-else v-model="selectedKeyId" class="qa__select">
                <option
                  v-for="k in keys"
                  :key="k.id"
                  :value="k.id"
                >
                  {{ k.name }}（{{ maskKey(k.key) }}）
                </option>
              </select>
            </label>
          </div>
          <dl class="qa__meta">
            <div class="qa__meta-row">
              <dt>{{ t('quickaccess.endpointLabel') }}</dt>
              <dd>{{ endpointURL || t('common.emptyDash') }}</dd>
            </div>
            <div class="qa__meta-row">
              <dt>{{ t('quickaccess.authLabel') }}</dt>
              <template v-if="selectedServerIsOpen">
                <dd>{{ t('quickaccess.authNone') }}</dd>
              </template>
              <template v-else>
                <dd>{{ accesskeyHeader }}: {{ selectedKey ? maskKey(selectedKey.key) : t('common.emptyDash') }}</dd>
              </template>
            </div>
          </dl>
        </div>

        <div class="qa__card">
          <h2 class="qa__card-title">{{ t('quickaccess.generateConfig') }}</h2>
          <div v-if="!endpointURL" class="qa__notice qa__notice--warn qa__notice--block">
            {{ t('quickaccess.noConfigWarn') }}
          </div>
          <template v-else>
            <div class="qa__code-head">
              <span class="qa__code-name">{{ t('quickaccess.importConfig', { client: client.name }) }}</span>
              <div class="qa__code-actions">
                <a
                  v-if="cherryInstallLink"
                  class="btn btn--primary"
                  :href="cherryInstallLink"
                  :title="t('quickaccess.openInClientTitle', { client: client.name })"
                >
                  {{ t('quickaccess.openInClient', { client: client.name }) }}
                </a>
                <button class="btn btn--ghost" type="button" :disabled="!configJSON" @click="copyConfig">
                  {{ copied ? t('common.copied') : t('common.copy') }}
                </button>
              </div>
            </div>
            <pre class="qa__code"><code>{{ configJSON }}</code></pre>
          </template>
        </div>

        <div class="qa__card">
          <h2 class="qa__card-title">{{ t('quickaccess.setupSteps', { client: client.name }) }}</h2>
          <ol class="qa__steps">
            <li v-for="(stepKey, i) in client.stepKeys" :key="i">{{ t(stepKey, { client: client.name }) }}</li>
          </ol>
        </div>
      </template>
    </template>
  </section>
</template>

<style scoped>
.page__desc {
  margin-top: 8px;
  font-size: 14px;
  color: var(--text-muted);
}

.qa__loading {
  margin-top: 24px;
  padding: 24px;
  font-size: 14px;
  color: var(--text-muted);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.qa__notice {
  margin-top: 20px;
  padding: 16px 20px;
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.qa__notice--warn {
  color: #856404;
  background: #fff3cd;
  border-color: #ffc107;
}

.qa__notice--block {
  margin-top: 0;
}

.qa__link {
  padding: 0;
  font-size: 14px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: transparent;
  border: none;
  cursor: pointer;
}

.qa__card {
  margin-top: 20px;
  padding: 24px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.qa__card-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--xauat-blue);
  margin-bottom: 16px;
}

.qa__clients {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.qa__client {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 220px;
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--bg);
}

.qa__client--active {
  border-color: var(--xauat-blue-light);
  background: var(--active-bg);
}

.qa__client-name {
  font-size: 15px;
  font-weight: 700;
  color: var(--text);
}

.qa__client-desc {
  font-size: 12px;
  color: var(--text-muted);
}

.qa__form {
  display: flex;
  gap: 20px;
  flex-wrap: wrap;
}

.qa__field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-width: 220px;
}

.qa__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.qa__note {
  padding: 9px 14px;
  font-size: 13px;
  line-height: 1.5;
  color: var(--text-muted);
  background: var(--bg);
  border: 1px dashed var(--border);
  border-radius: 10px;
}

.qa__select {
  width: 100%;
  padding: 9px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  cursor: pointer;
}

.qa__select:focus {
  border-color: var(--xauat-blue);
}

.qa__meta {
  margin-top: 18px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.qa__meta-row {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  padding: 12px 16px;
}

.qa__meta-row + .qa__meta-row {
  border-top: 1px solid var(--border);
}

.qa__meta-row dt {
  flex: 0 0 120px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.qa__meta-row dd {
  flex: 1;
  font-size: 13px;
  color: var(--text);
  word-break: break-all;
  font-family: "SF Mono", "Menlo", "Consolas", monospace;
}

.qa__code-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 12px;
}

.qa__code-actions {
  display: flex;
  align-items: center;
  gap: 10px;
}

.qa__code-name {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.qa__code {
  padding: 16px;
  overflow-x: auto;
  font-family: "SF Mono", "Menlo", "Consolas", monospace;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
}

.qa__steps {
  margin: 0;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.qa__steps li {
  font-size: 14px;
  color: var(--text);
}

.btn {
  display: inline-flex;
  align-items: center;
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 999px;
  border: 1px solid transparent;
  cursor: pointer;
  text-decoration: none;
  transition: background 0.2s, border-color 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn--primary {
  color: #fff;
  background: var(--xauat-blue);
}

.btn--primary:hover:not(:disabled) {
  background: var(--xauat-blue-light);
}

.btn--ghost {
  color: var(--xauat-blue);
  background: transparent;
  border-color: var(--border);
}

.btn--ghost:hover:not(:disabled) {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}
</style>

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
import { ref, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchRoles, fetchSettings, saveSetting, testLdap, testSmtp } from '../api.js'
import { applyPlatform } from '../stores/settings.js'

const { t } = useI18n()

const tab = ref('ops')

const logOpen = ref(true)
const smtpOpen = ref(true)
const apiGwOpen = ref(true)
const networkSecurityOpen = ref(true)
const auditLogOpen = ref(true)
const authOpen = ref(true)
const userOpsOpen = ref(true)
const platformOpen = ref(true)

const logForm = ref({
  syslogEnabled: false,
  logPath: './logs',
  logLevel: 'info',
  logPrefix: 'mcp_plat',
  maxSize: 100,
  maxBackups: 10,
  maxAge: 30,
  compress: false,
  syslogHost: '',
  syslogPort: 514,
  syslogProtocol: 'tcp'
})

const LOG_LEVELS = ['debug', 'info', 'warn', 'error']

const smtpForm = ref({
  enabled: false,
  host: '',
  port: 465,
  encryption: 'ssl',
  username: '',
  password: '',
  fromAddress: '',
  fromName: ''
})

const smtpTestTo = ref('')
const smtpTesting = ref(false)
const smtpTestResult = ref(null)

const apiGwForm = ref({
  provider: 'apisix',
  adminUrl: '',
  defaultPublishDomain: '',
  adminKey: '',
  accesskeyHeader: '',
  authGrpcAddrs: []
})

const API_GW_PROVIDERS = [
  { key: 'apisix', label: 'Apache APISIX' },
  { key: 'kong', label: 'Kong' },
  { key: 'tyk', label: 'Tyk' }
]

const showKongNotice = ref(false)

watch(
  () => apiGwForm.value.provider,
  (val, oldVal) => {
    if (val === 'kong' && oldVal !== 'kong') {
      showKongNotice.value = true
    }
  }
)

const networkSecurityForm = ref({
  allowlist: ''
})

const auditLogForm = ref({
  enabled: false,
  grpcAddrs: ['127.0.0.1:9091']
})

const authMethod = ref('cas')

const AUTH_METHODS = [
  { key: 'cas', label: 'CAS' },
  { key: 'ldap', label: 'LDAP' },
  { key: 'oauth', label: 'OAuth' }
]

const casForm = ref({
  serverUrl: '',
  serviceUrl: '',
  version: '3.0'
})

const savedCas = ref({})
const casAttrMapping = ref([])

const ldapForm = ref({
  host: '',
  port: 389,
  baseDn: '',
  bindDn: '',
  bindPassword: '',
  userFilter: ''
})

const savedLdap = ref({})

const PLATFORM_FIELDS = [
  { key: 'uid', labelKey: 'settings.fields.uid' },
  { key: 'name', labelKey: 'settings.fields.name' },
  { key: 'email', labelKey: 'settings.fields.email' },
  { key: 'phone', labelKey: 'settings.fields.phone' },
  { key: 'organization', labelKey: 'settings.fields.organization' }
]

const ldapAttrMapping = ref([])

function addAttrMappingRow(list) {
  list.push({ field: '', attr: '' })
}

function removeAttrMappingRow(list, index) {
  list.splice(index, 1)
}

function mappingToObject(list) {
  const obj = {}
  for (const item of list) {
    if (item.field && item.attr) {
      obj[item.field] = item.attr
    }
  }
  return obj
}

function objectToMapping(obj) {
  if (!obj) return []
  return Object.entries(obj).map(([field, attr]) => ({ field, attr }))
}

function unusedFields(list, rowIndex) {
  const used = new Set(
    list
      .filter((_, i) => i !== rowIndex)
      .map(item => item.field)
      .filter(Boolean)
  )
  return PLATFORM_FIELDS.filter(f => !used.has(f.key))
}

function platformFieldLabel(key) {
  const f = PLATFORM_FIELDS.find(item => item.key === key)
  return f ? t(f.labelKey) : key
}

const ldapTestUser = ref('')
const ldapTestResult = ref(null)
const ldapTesting = ref(false)

async function runLdapTest() {
  if (!ldapTestUser.value.trim()) return
  ldapTesting.value = true
  ldapTestResult.value = null
  try {
    const saved = savedLdap.value || {}
    const ldap = {}
    for (const key of ['host', 'baseDn', 'bindDn', 'bindPassword', 'userFilter']) {
      ldap[key] = ldapForm.value[key] !== '' ? ldapForm.value[key] : (saved[key] || '')
    }
    ldap.port = ldapForm.value.port !== '' && ldapForm.value.port != null ? ldapForm.value.port : (saved.port ?? 389)
    ldap.attrMapping = mappingToObject(ldapAttrMapping.value)

    const data = await testLdap({
      ldap,
      username: ldapTestUser.value.trim()
    })
    ldapTestResult.value = data
  } catch (e) {
    ldapTestResult.value = { success: false, message: e.message || t('settings.auth.ldap.testFailed') }
  }
  ldapTesting.value = false
}

const oauthForm = ref({
  authorizeUrl: '',
  tokenUrl: '',
  userinfoUrl: '',
  clientId: '',
  clientSecret: '',
  redirectUrl: '',
  scope: ''
})

const platformForm = ref({
  name: '某某大学',
  logoUrl: '',
  siteUrl: '',
  loginBackground: ''
})

const quickAccessOpen = ref(true)

const QUICK_ACCESS_CLIENTS = [
  { key: 'cherrystudio', label: 'CherryStudio' }
]

const quickAccessForm = ref({
  enabledClients: ['cherrystudio'],
  scheme: 'https'
})

const roles = ref([])
const filterAttribute = ref('username')
const roleRules = ref([])
const maxAccessKeys = ref(5)
const accessKeyCron = ref('')

const ROLE_ATTRS = [
  { key: 'username', labelKey: 'settings.fields.username' },
  { key: 'uid', labelKey: 'settings.fields.uid' },
  { key: 'name', labelKey: 'settings.fields.name' },
  { key: 'email', labelKey: 'settings.fields.email' },
  { key: 'phone', labelKey: 'settings.fields.phone' },
  { key: 'organization', labelKey: 'settings.fields.organization' }
]

const apiGwConfigured = ref(false)

const saving = ref('')
const tips = ref({})

async function saveSection(key, payload) {
  saving.value = key
  try {
    await saveSetting(key, payload)
    tips.value = { ...tips.value, [key]: t('settings.general.saved') }
  } catch (e) {
    tips.value = { ...tips.value, [key]: e.message || t('settings.general.saveFailed') }
  }
  saving.value = ''
  setTimeout(() => {
    tips.value = { ...tips.value, [key]: '' }
  }, 2000)
}

function saveLogSettings() {
  saveSection('log', logForm.value)
}

function saveSmtpSettings() {
  saveSection('smtp', smtpForm.value)
}

async function runSmtpTest() {
  if (!smtpTestTo.value.trim()) return
  smtpTesting.value = true
  smtpTestResult.value = null
  try {
    await testSmtp({
      smtp: {
        host: smtpForm.value.host,
        port: smtpForm.value.port,
        encryption: smtpForm.value.encryption,
        username: smtpForm.value.username,
        password: smtpForm.value.password,
        fromAddress: smtpForm.value.fromAddress,
        fromName: smtpForm.value.fromName
      },
      to: smtpTestTo.value.trim()
    })
    smtpTestResult.value = { ok: true, message: t('settings.mail.testOk') }
  } catch (e) {
    smtpTestResult.value = { ok: false, message: e.message || t('settings.mail.sendFailed') }
  }
  smtpTesting.value = false
}

function saveApiGwSettings() {
  saveSection('api_gateway', apiGwForm.value)
}

function saveNetworkSecuritySettings() {
  saveSection('network_security', {
    allowlist: networkSecurityForm.value.allowlist
      .split('\n')
      .map(line => line.trim())
      .filter(line => line !== '')
  })
}

function saveAuditLogSettings() {
  saveSection('audit_log', auditLogForm.value)
}

function saveAuthSettings() {
  const saved = savedLdap.value || {}
  const ldap = {}
  for (const key of ['host', 'baseDn', 'bindDn', 'bindPassword', 'userFilter']) {
    ldap[key] = ldapForm.value[key] !== '' ? ldapForm.value[key] : (saved[key] || '')
  }
  ldap.port = ldapForm.value.port !== '' && ldapForm.value.port != null ? ldapForm.value.port : (saved.port ?? 389)
  ldap.attrMapping = mappingToObject(ldapAttrMapping.value)
  saveSection('auth', {
    method: authMethod.value,
    cas: { ...casForm.value, attrMapping: mappingToObject(casAttrMapping.value) },
    ldap,
    oauth: oauthForm.value
  })
}

function saveUserOpsSettings() {
  const rules = roleRules.value
    .filter(item => item.pattern.trim() !== '' && item.roleId !== '')
    .map(item => ({
      roleId: item.roleId,
      pattern: item.pattern.trim()
    }))
  saveSection('user_ops', {
    filterAttribute: filterAttribute.value,
    roleRules: rules,
    maxAccessKeys: maxAccessKeys.value,
    accessKeyCron: accessKeyCron.value
  })
}

function addRoleRule() {
  roleRules.value.push({ pattern: '', roleId: '' })
}

function removeRoleRule(index) {
  roleRules.value.splice(index, 1)
}

function savePlatformSettings() {
  saveSection('platform', platformForm.value)
  applyPlatform(platformForm.value)
}

function saveQuickAccessSettings() {
  saveSection('quick_access', quickAccessForm.value)
}

function toAddrList(value) {
  if (Array.isArray(value)) return value.slice()
  if (typeof value === 'string' && value.trim()) {
    return value
      .split(/[,\n]/)
      .map(s => s.trim())
      .filter(Boolean)
  }
  return []
}

function addAddr(list) {
  list.push('')
}

function removeAddr(list, index) {
  list.splice(index, 1)
}

function toggleQuickAccessClient(key) {
  const list = quickAccessForm.value.enabledClients
  const idx = list.indexOf(key)
  if (idx === -1) list.push(key)
  else list.splice(idx, 1)
}

function applySettings(data) {
  if (data.log) Object.assign(logForm.value, data.log)
  if (data.smtp) Object.assign(smtpForm.value, data.smtp)
  if (data.api_gateway) {
    apiGwConfigured.value = data.api_gateway.configured || false
    if (data.api_gateway.provider) apiGwForm.value.provider = data.api_gateway.provider
    if (data.api_gateway.adminUrl != null) apiGwForm.value.adminUrl = data.api_gateway.adminUrl
    if (data.api_gateway.defaultPublishDomain != null) apiGwForm.value.defaultPublishDomain = data.api_gateway.defaultPublishDomain
    if (data.api_gateway.accesskeyHeader != null) apiGwForm.value.accesskeyHeader = data.api_gateway.accesskeyHeader
    if (data.api_gateway.authGrpcAddrs != null || data.api_gateway.authGrpcAddr != null) {
      apiGwForm.value.authGrpcAddrs = toAddrList(data.api_gateway.authGrpcAddrs ?? data.api_gateway.authGrpcAddr)
    }
  }
  if (data.network_security && data.network_security.allowlist) {
    networkSecurityForm.value.allowlist = data.network_security.allowlist.join('\n')
  }
  if (data.audit_log) {
    if (data.audit_log.enabled != null) auditLogForm.value.enabled = data.audit_log.enabled
    if (data.audit_log.grpcAddrs != null || data.audit_log.grpcAddr != null) {
      auditLogForm.value.grpcAddrs = toAddrList(data.audit_log.grpcAddrs ?? data.audit_log.grpcAddr)
    }
  }
  if (data.auth) {
    if (data.auth.method) authMethod.value = data.auth.method
    if (data.auth.cas) {
      savedCas.value = { ...data.auth.cas }
      casForm.value.serverUrl = data.auth.cas.serverUrl || ''
      casForm.value.serviceUrl = data.auth.cas.serviceUrl || ''
      casForm.value.version = data.auth.cas.version || '3.0'
      casAttrMapping.value = objectToMapping(data.auth.cas.attrMapping)
    }
    if (data.auth.ldap) {
      savedLdap.value = { ...data.auth.ldap }
      if (authMethod.value === 'ldap') {
        ldapForm.value.host = ''
        ldapForm.value.port = ''
        ldapForm.value.baseDn = ''
        ldapForm.value.bindDn = ''
        ldapForm.value.bindPassword = ''
        ldapForm.value.userFilter = ''
      } else {
        ldapForm.value.host = data.auth.ldap.host || ''
        ldapForm.value.port = data.auth.ldap.port ?? 389
        ldapForm.value.baseDn = data.auth.ldap.baseDn || ''
        ldapForm.value.bindDn = data.auth.ldap.bindDn || ''
        ldapForm.value.bindPassword = data.auth.ldap.bindPassword || ''
        ldapForm.value.userFilter = data.auth.ldap.userFilter || ''
      }
      ldapAttrMapping.value = objectToMapping(data.auth.ldap.attrMapping)
    }
    if (data.auth.oauth) Object.assign(oauthForm.value, data.auth.oauth)
  }
  if (data.user_ops) {
    if (data.user_ops.filterAttribute) filterAttribute.value = data.user_ops.filterAttribute
    if (data.user_ops.maxAccessKeys != null) maxAccessKeys.value = data.user_ops.maxAccessKeys
    accessKeyCron.value = data.user_ops.accessKeyCron || ''
    roleRules.value = Array.isArray(data.user_ops.roleRules)
      ? data.user_ops.roleRules.map(item => ({
          pattern: item.pattern || '',
          roleId: item.roleId != null ? item.roleId : ''
        }))
      : []
  }
  if (data.platform) {
    Object.assign(platformForm.value, data.platform)
    applyPlatform(data.platform)
  }
  if (data.quick_access) {
    if (Array.isArray(data.quick_access.enabledClients)) {
      quickAccessForm.value.enabledClients = data.quick_access.enabledClients
    }
    if (data.quick_access.scheme) {
      quickAccessForm.value.scheme = data.quick_access.scheme
    }
  }
}

onMounted(async () => {
  await Promise.all([
    fetchRoles()
      .then((d) => { roles.value = d })
      .catch(() => {}),
    fetchSettings()
      .then((d) => { applySettings(d || {}) })
      .catch(() => {})
  ])
})
</script>

<template>
  <section class="settings">
    <div class="settings__head">
      <h1 class="page__title">{{ t('nav.settings') }}</h1>
    </div>

    <div class="settings__tabs">
      <button
        class="settings__tab"
        :class="{ 'settings__tab--active': tab === 'ops' }"
        type="button"
        @click="tab = 'ops'"
      >
        {{ t('settings.tabs.ops') }}
      </button>
      <button
        class="settings__tab"
        :class="{ 'settings__tab--active': tab === 'operation' }"
        type="button"
        @click="tab = 'operation'"
      >
        {{ t('settings.tabs.operation') }}
      </button>
    </div>

    <div v-if="tab === 'ops'" class="settings__panel">
      <div class="collapse">
        <button class="collapse__head" type="button" @click="logOpen = !logOpen">
          <span class="collapse__title">{{ t('settings.log.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': logOpen }">▾</span>
        </button>
        <div v-show="logOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">{{ t('settings.log.enableSyslog') }}</span>
            <label class="switch">
              <input type="checkbox" v-model="logForm.syslogEnabled" />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ logForm.syslogEnabled ? t('common.enabled') : t('common.disabled') }}</span>
            </label>
          </div>

          <template v-if="!logForm.syslogEnabled">
            <div class="ldap-notice">
              {{ t('settings.log.offNotice') }}
            </div>
            <label class="field">
              <span class="field__label">{{ t('settings.log.pathLabel') }}</span>
              <input
                v-model="logForm.logPath"
                class="field__input"
                type="text"
                placeholder="./logs"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.log.levelLabel') }}</span>
              <select v-model="logForm.logLevel" class="field__input">
                <option v-for="lv in LOG_LEVELS" :key="lv" :value="lv">{{ lv }}</option>
              </select>
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.log.fileNameLabel') }}</span>
              <input
                v-model="logForm.logPrefix"
                class="field__input"
                type="text"
                placeholder="mcp_plat"
              />
              <span class="field__help-text">{{ t('settings.log.fileNameHelp') }}</span>
            </label>

            <p class="field__section-title">{{ t('settings.log.rotation') }}</p>
            <div class="field-row">
              <label class="field">
                <span class="field__label">{{ t('settings.log.maxSizeLabel') }}</span>
                <input
                  v-model.number="logForm.maxSize"
                  class="field__input"
                  type="number"
                  min="1"
                  placeholder="100"
                />
              </label>
              <label class="field">
                <span class="field__label">{{ t('settings.log.maxBackupsLabel') }}</span>
                <input
                  v-model.number="logForm.maxBackups"
                  class="field__input"
                  type="number"
                  min="1"
                  placeholder="10"
                />
              </label>
            </div>
            <div class="field-row">
              <label class="field">
                <span class="field__label">{{ t('settings.log.maxAgeLabel') }}</span>
                <input
                  v-model.number="logForm.maxAge"
                  class="field__input"
                  type="number"
                  min="1"
                  placeholder="30"
                />
              </label>
              <label class="field">
                <span class="field__label">{{ t('settings.log.compressLabel') }}</span>
                <select v-model="logForm.compress" class="field__input">
                  <option :value="false">{{ t('settings.log.compressNone') }}</option>
                  <option :value="true">{{ t('settings.log.compressGzip') }}</option>
                </select>
              </label>
            </div>
          </template>

          <template v-else>
            <div class="ldap-notice">
              {{ t('settings.log.onNotice') }}
            </div>
            <label class="field">
              <span class="field__label">{{ t('settings.log.syslogHostLabel') }}</span>
              <input
                v-model="logForm.syslogHost"
                class="field__input"
                type="text"
                :placeholder="t('settings.log.syslogHostPlaceholder')"
              />
            </label>
            <div class="field-row">
              <label class="field">
                <span class="field__label">{{ t('settings.log.portLabel') }}</span>
                <input
                  v-model.number="logForm.syslogPort"
                  class="field__input"
                  type="number"
                  placeholder="514"
                />
              </label>
              <label class="field">
                <span class="field__label">{{ t('settings.log.protocolLabel') }}</span>
                <select v-model="logForm.syslogProtocol" class="field__input">
                  <option value="tcp">TCP</option>
                  <option value="udp">UDP</option>
                </select>
              </label>
            </div>
          </template>

          <div class="collapse__actions">
            <span v-if="tips.log" class="save-tip">{{ tips.log }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'log'" @click="saveLogSettings">
              {{ saving === 'log' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="smtpOpen = !smtpOpen">
          <span class="collapse__title">{{ t('settings.mail.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': smtpOpen }">▾</span>
        </button>
        <div v-show="smtpOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">{{ t('settings.mail.enable') }}</span>
            <label class="switch">
              <input type="checkbox" v-model="smtpForm.enabled" />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ smtpForm.enabled ? t('common.enabled') : t('common.disabled') }}</span>
            </label>
            <span class="field__help-text">{{ t('settings.mail.toggleHint') }}</span>
          </div>

          <div class="field-row">
            <label class="field">
              <span class="field__label">{{ t('settings.mail.hostLabel') }}</span>
              <input
                v-model="smtpForm.host"
                class="field__input"
                type="text"
                placeholder="smtp.example.edu.cn"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.mail.portLabel') }}</span>
              <input
                v-model.number="smtpForm.port"
                class="field__input"
                type="number"
                placeholder="465"
              />
            </label>
          </div>
          <label class="field">
            <span class="field__label">{{ t('settings.mail.encryptionLabel') }}</span>
            <select v-model="smtpForm.encryption" class="field__input">
              <option value="none">{{ t('settings.mail.encryptionNone') }}</option>
              <option value="ssl">SSL/TLS</option>
              <option value="starttls">STARTTLS</option>
            </select>
          </label>
          <label class="field">
            <span class="field__label">{{ t('settings.mail.usernameLabel') }}</span>
            <input
              v-model="smtpForm.username"
              class="field__input"
              type="text"
              :placeholder="t('settings.mail.usernamePlaceholder')"
            />
          </label>
          <label class="field">
            <span class="field__label">{{ t('settings.mail.passwordLabel') }}</span>
            <input
              v-model="smtpForm.password"
              class="field__input"
              type="password"
              :placeholder="t('settings.mail.passwordPlaceholder')"
            />
          </label>
          <div class="field-row">
            <label class="field">
              <span class="field__label">{{ t('settings.mail.fromAddressLabel') }}</span>
              <input
                v-model="smtpForm.fromAddress"
                class="field__input"
                type="text"
                placeholder="noreply@example.edu.cn"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.mail.fromNameLabel') }}</span>
              <input
                v-model="smtpForm.fromName"
                class="field__input"
                type="text"
                :placeholder="t('settings.mail.fromNamePlaceholder')"
              />
            </label>
          </div>

          <div class="smtp-test">
            <p class="field__section-title">{{ t('settings.mail.sendTest') }}</p>
            <div class="smtp-test__row">
              <input
                v-model="smtpTestTo"
                class="field__input smtp-test__input"
                type="text"
                :placeholder="t('settings.mail.testToPlaceholder')"
                @keyup.enter="runSmtpTest"
              />
              <button
                class="btn btn--primary smtp-test__btn"
                type="button"
                :disabled="!smtpTestTo.trim() || smtpTesting"
                @click="runSmtpTest"
              >
                {{ smtpTesting ? t('settings.mail.sending') : t('settings.mail.sendTest') }}
              </button>
            </div>
            <div
              v-if="smtpTestResult"
              class="smtp-test__result"
              :class="smtpTestResult.ok ? 'smtp-test__result--ok' : 'smtp-test__result--err'"
            >
              {{ smtpTestResult.message }}
            </div>
          </div>

          <div class="collapse__actions">
            <span v-if="tips.smtp" class="save-tip">{{ tips.smtp }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'smtp'" @click="saveSmtpSettings">
              {{ saving === 'smtp' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="apiGwOpen = !apiGwOpen">
          <span class="collapse__title">{{ t('settings.apiGw.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': apiGwOpen }">▾</span>
        </button>
        <div v-show="apiGwOpen" class="collapse__body">
          <div v-if="apiGwConfigured" class="api-gw-warning">
            <span class="api-gw-warning__icon">!</span>
            {{ t('settings.apiGw.configuredWarning') }}
          </div>
          <label class="field">
            <span class="field__label">{{ t('settings.apiGw.providerLabel') }}</span>
            <select v-model="apiGwForm.provider" class="field__input">
              <option
                v-for="gw in API_GW_PROVIDERS"
                :key="gw.key"
                :value="gw.key"
              >
                {{ gw.label }}
              </option>
            </select>
            <span v-if="apiGwForm.provider === 'kong'" class="field__help-text">
              {{ t('settings.apiGw.kongHelpPre') }}
              <code>http://127.0.0.1:8001</code>{{ t('settings.apiGw.kongHelpPost') }}
            </span>
          </label>

          <label class="field">
            <span class="field__label">
              {{ t('settings.apiGw.adminUrlLabel') }}
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  <template v-if="apiGwForm.provider === 'kong'">
                    {{ t('settings.apiGw.adminUrlTooltipKongPre') }}
                    <code>http://127.0.0.1:8001</code>{{ t('settings.apiGw.adminUrlTooltipKongPost') }}
                  </template>
                  <template v-else>
                    {{ t('settings.apiGw.adminUrlTooltipDefaultPre') }}
                    <code>http://127.0.0.1:9180</code>{{ t('settings.apiGw.adminUrlTooltipDefaultPost') }}
                  </template>
                </span>
              </span>
            </span>
            <input
              v-model="apiGwForm.adminUrl"
              class="field__input"
              type="text"
              :placeholder="apiGwForm.provider === 'kong' ? 'http://127.0.0.1:8001' : 'http://127.0.0.1:9180'"
            />
          </label>

          <label class="field">
            <span class="field__label">{{ t('settings.apiGw.defaultPublishLabel') }}</span>
            <input
              v-model="apiGwForm.defaultPublishDomain"
              class="field__input"
              type="text"
              :placeholder="t('settings.apiGw.defaultPublishPlaceholder')"
            />
          </label>

          <label class="field">
            <span class="field__label">{{ t('settings.apiGw.adminKeyLabel') }}</span>
            <input
              v-model="apiGwForm.adminKey"
              class="field__input"
              type="password"
              :placeholder="t('settings.apiGw.adminKeyPlaceholder')"
            />
            <span v-if="apiGwForm.provider === 'kong'" class="field__help-text">
              {{ t('settings.apiGw.adminKeyKongHint') }}
            </span>
          </label>

          <div class="field">
            <span class="field__label">
              {{ t('settings.apiGw.authGrpcLabel') }}
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  {{ t('settings.apiGw.authGrpcTip1') }}
                  <code>:9090</code>{{ t('settings.apiGw.authGrpcTip2') }}<code>127.0.0.1:9090</code>{{ t('settings.apiGw.authGrpcTip3') }}
                </span>
              </span>
            </span>
            <div
              v-for="(addr, i) in apiGwForm.authGrpcAddrs"
              :key="i"
              class="addr-row"
            >
              <input
                v-model="apiGwForm.authGrpcAddrs[i]"
                class="field__input addr-row__input"
                type="text"
                placeholder="127.0.0.1:9090"
              />
              <button
                class="addr-row__remove"
                type="button"
                :title="t('settings.general.remove')"
                @click="removeAddr(apiGwForm.authGrpcAddrs, i)"
              >
                ×
              </button>
            </div>
            <button
              class="addr-row__add"
              type="button"
              @click="addAddr(apiGwForm.authGrpcAddrs)"
            >
              {{ t('settings.general.addAddress') }}
            </button>
          </div>

          <label class="field">
            <span class="field__label">
              {{ t('settings.apiGw.headerLabel') }}
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  {{ t('settings.apiGw.headerTip1') }}<br />
                  {{ t('settings.apiGw.headerTip2') }}
                  <code>X-Access-Key</code>{{ t('settings.apiGw.headerTip3') }}
                </span>
              </span>
            </span>
            <input
              v-model="apiGwForm.accesskeyHeader"
              class="field__input"
              type="text"
              placeholder="X-Access-Key"
            />
          </label>

          <div class="collapse__actions">
            <span v-if="tips.api_gateway" class="save-tip">{{ tips.api_gateway }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'api_gateway'" @click="saveApiGwSettings">
              {{ saving === 'api_gateway' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showKongNotice" class="dialog-overlay" @click.self="showKongNotice = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('settings.kongDialog.title') }}</h2>
          <p class="dialog__desc">
            {{ t('settings.kongDialog.intro1') }}<strong>Kong</strong>{{ t('settings.kongDialog.intro2') }}
          </p>
          <ul class="kong-notice__list">
            <li>{{ t('settings.kongDialog.item1a') }}<strong>{{ t('settings.kongDialog.item1b') }}</strong>{{ t('settings.kongDialog.item1c') }}</li>
            <li>{{ t('settings.kongDialog.item2') }}</li>
            <li>{{ t('settings.kongDialog.item3') }}</li>
          </ul>
          <p class="dialog__desc">
            {{ t('settings.kongDialog.dockerPre') }}
            <code>http://127.0.0.1:8001</code>{{ t('settings.kongDialog.dockerPost') }}
          </p>
          <div class="dialog__actions">
            <button class="btn btn--primary" type="button" @click="showKongNotice = false">
              {{ t('settings.general.gotIt') }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="networkSecurityOpen = !networkSecurityOpen">
          <span class="collapse__title">{{ t('settings.network.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': networkSecurityOpen }">▾</span>
        </button>
        <div v-show="networkSecurityOpen" class="collapse__body">
          <label class="field">
            <span class="field__label">
              {{ t('settings.network.allowLabel') }}
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  {{ t('settings.network.tip1') }}
                  <code>10.0.0.0/8</code>{{ t('settings.network.tip2') }}<br />
                  {{ t('settings.network.tip3') }}<br />
                  {{ t('settings.network.tip4') }}
                </span>
              </span>
            </span>
            <textarea
              v-model="networkSecurityForm.allowlist"
              class="field__input field__textarea"
              rows="5"
              placeholder="10.0.0.0/8&#10;172.16.0.0/12&#10;192.168.1.0/24"
            ></textarea>
          </label>

          <div class="collapse__actions">
            <span v-if="tips.network_security" class="save-tip">{{ tips.network_security }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'network_security'" @click="saveNetworkSecuritySettings">
              {{ saving === 'network_security' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>
      <div class="collapse">
        <button class="collapse__head" type="button" @click="auditLogOpen = !auditLogOpen">
          <span class="collapse__title">{{ t('settings.audit.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': auditLogOpen }">▾</span>
        </button>
        <div v-show="auditLogOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">{{ t('settings.audit.enable') }}</span>
            <label class="switch">
              <input type="checkbox" v-model="auditLogForm.enabled" />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ auditLogForm.enabled ? t('common.enabled') : t('common.disabled') }}</span>
            </label>
          </div>

          <div class="field">
            <span class="field__label">
              {{ t('settings.audit.grpcAddrLabel') }}
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  {{ t('settings.audit.tip1') }}
                  <code>:9091</code>{{ t('settings.audit.tip2') }}<code>127.0.0.1:9091</code>{{ t('settings.audit.tip3') }}
                </span>
              </span>
            </span>
            <div
              v-for="(addr, i) in auditLogForm.grpcAddrs"
              :key="i"
              class="addr-row"
            >
              <input
                v-model="auditLogForm.grpcAddrs[i]"
                class="field__input addr-row__input"
                type="text"
                placeholder="127.0.0.1:9091"
              />
              <button
                class="addr-row__remove"
                type="button"
                :title="t('settings.general.remove')"
                @click="removeAddr(auditLogForm.grpcAddrs, i)"
              >
                ×
              </button>
            </div>
            <button
              class="addr-row__add"
              type="button"
              @click="addAddr(auditLogForm.grpcAddrs)"
            >
              {{ t('settings.general.addAddress') }}
            </button>
          </div>

          <div class="collapse__actions">
            <span v-if="tips.audit_log" class="save-tip">{{ tips.audit_log }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'audit_log'" @click="saveAuditLogSettings">
              {{ saving === 'audit_log' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'operation'" class="settings__panel">
      <div class="collapse">
        <button class="collapse__head" type="button" @click="platformOpen = !platformOpen">
          <span class="collapse__title">{{ t('settings.platform.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': platformOpen }">▾</span>
        </button>
        <div v-show="platformOpen" class="collapse__body">
          <label class="field">
            <span class="field__label">{{ t('settings.platform.nameLabel') }}</span>
            <input
              v-model="platformForm.name"
              class="field__input"
              type="text"
              :placeholder="t('settings.platform.namePlaceholder')"
            />
          </label>

          <label class="field">
            <span class="field__label">{{ t('settings.platform.logoLabel') }}</span>
            <input
              v-model="platformForm.logoUrl"
              class="field__input"
              type="text"
              :placeholder="t('settings.platform.logoPlaceholder')"
            />
          </label>

          <label class="field">
            <span class="field__label">{{ t('settings.platform.loginBgLabel') }}</span>
            <input
              v-model="platformForm.loginBackground"
              class="field__input"
              type="text"
              :placeholder="t('settings.platform.loginBgPlaceholder')"
            />
            <span class="field__help-text">{{ t('settings.platform.loginBgHint') }}</span>
          </label>

          <label class="field">
            <span class="field__label">{{ t('settings.platform.siteUrlLabel') }}</span>
            <input
              v-model="platformForm.siteUrl"
              class="field__input"
              type="text"
              :placeholder="t('settings.platform.siteUrlPlaceholder')"
            />
          </label>

          <div class="collapse__actions">
            <span v-if="tips.platform" class="save-tip">{{ tips.platform }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'platform'" @click="savePlatformSettings">
              {{ saving === 'platform' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="authOpen = !authOpen">
          <span class="collapse__title">{{ t('settings.auth.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': authOpen }">▾</span>
        </button>
        <div v-show="authOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">{{ t('settings.auth.methodLabel') }}</span>
            <div class="seg">
              <button
                v-for="m in AUTH_METHODS"
                :key="m.key"
                class="seg__item"
                :class="{ 'seg__item--active': authMethod === m.key }"
                type="button"
                @click="authMethod = m.key"
              >
                {{ m.label }}
              </button>
            </div>
          </div>

          <template v-if="authMethod === 'cas'">
            <label class="field">
              <span class="field__label">{{ t('settings.auth.cas.serverUrlLabel') }}</span>
              <input
                v-model="casForm.serverUrl"
                class="field__input"
                type="text"
                placeholder="https://cas.example.edu.cn/cas"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.cas.serviceUrlLabel') }}</span>
              <input
                v-model="casForm.serviceUrl"
                class="field__input"
                type="text"
                placeholder="https://mcp.example.edu.cn/login/cas"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.cas.versionLabel') }}</span>
              <select v-model="casForm.version" class="field__input">
                <option value="2.0">2.0</option>
                <option value="3.0">3.0</option>
              </select>
            </label>

            <p class="field__section-title">{{ t('settings.auth.cas.mappingTitle') }}</p>
            <div
              v-for="(item, i) in casAttrMapping"
              :key="i"
              class="mapping-row"
            >
              <select
                v-model="item.field"
                class="field__input mapping-row__select"
              >
                <option value="" disabled>{{ t('settings.auth.cas.selectPlatformField') }}</option>
                <option
                  v-for="f in unusedFields(casAttrMapping, i)"
                  :key="f.key"
                  :value="f.key"
                >
                  {{ t(f.labelKey) }}
                </option>
              </select>
              <span class="mapping-row__arrow">→</span>
              <input
                v-model="item.attr"
                class="field__input mapping-row__input"
                type="text"
                :placeholder="t('settings.auth.cas.attrPlaceholder')"
              />
              <button
                class="mapping-row__remove"
                type="button"
                @click="removeAttrMappingRow(casAttrMapping, i)"
                :title="t('settings.general.remove')"
              >
                ×
              </button>
            </div>
            <button
              v-if="casAttrMapping.length < PLATFORM_FIELDS.length"
              class="mapping-row__add"
              type="button"
              @click="addAttrMappingRow(casAttrMapping)"
            >
              {{ t('settings.auth.cas.addMapping') }}
            </button>
          </template>

          <template v-if="authMethod === 'ldap'">
            <div class="ldap-notice">
              {{ t('settings.auth.ldap.enabledNotice') }}
            </div>
            <div class="field-row">
              <label class="field">
                <span class="field__label">{{ t('settings.auth.ldap.hostLabel') }}</span>
                <input
                  v-model="ldapForm.host"
                  class="field__input"
                  type="text"
                  placeholder="ldap.example.edu.cn"
                />
              </label>
              <label class="field">
                <span class="field__label">{{ t('settings.auth.ldap.portLabel') }}</span>
                <input
                  v-model.number="ldapForm.port"
                  class="field__input"
                  type="number"
                  placeholder="389"
                />
              </label>
            </div>
            <label class="field">
              <span class="field__label">Base DN</span>
              <input
                v-model="ldapForm.baseDn"
                class="field__input"
                type="text"
                placeholder="dc=example,dc=edu,dc=cn"
              />
            </label>
            <label class="field">
              <span class="field__label">Bind DN</span>
              <input
                v-model="ldapForm.bindDn"
                class="field__input"
                type="text"
                placeholder="cn=admin,dc=example,dc=edu,dc=cn"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.ldap.bindPasswordLabel') }}</span>
              <input
                v-model="ldapForm.bindPassword"
                class="field__input"
                type="password"
                :placeholder="t('settings.auth.ldap.bindPasswordPlaceholder')"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.ldap.userFilterLabel') }}</span>
              <input
                v-model="ldapForm.userFilter"
                class="field__input"
                type="text"
                placeholder="(uid=%s)"
              />
            </label>

            <p class="field__section-title">{{ t('settings.auth.ldap.mappingTitle') }}</p>
            <div
              v-for="(item, i) in ldapAttrMapping"
              :key="i"
              class="mapping-row"
            >
              <select
                v-model="item.field"
                class="field__input mapping-row__select"
              >
                <option value="" disabled>{{ t('settings.auth.ldap.selectPlatformField') }}</option>
                <option
                  v-for="f in unusedFields(ldapAttrMapping, i)"
                  :key="f.key"
                  :value="f.key"
                >
                  {{ t(f.labelKey) }}
                </option>
              </select>
              <span class="mapping-row__arrow">→</span>
              <input
                v-model="item.attr"
                class="field__input mapping-row__input"
                type="text"
                :placeholder="t('settings.auth.ldap.ldapAttrPlaceholder')"
              />
              <button
                class="mapping-row__remove"
                type="button"
                @click="removeAttrMappingRow(ldapAttrMapping, i)"
                :title="t('settings.general.remove')"
              >
                ×
              </button>
            </div>
            <button
              v-if="ldapAttrMapping.length < PLATFORM_FIELDS.length"
              class="mapping-row__add"
              type="button"
              @click="addAttrMappingRow(ldapAttrMapping)"
            >
              {{ t('settings.auth.ldap.addMapping') }}
            </button>

            <div class="ldap-test">
              <p class="field__section-title">{{ t('settings.auth.ldap.testTitle') }}</p>
              <div class="ldap-test__row">
                <input
                  v-model="ldapTestUser"
                  class="field__input ldap-test__input"
                  type="text"
                  :placeholder="t('settings.auth.ldap.testPlaceholder')"
                  @keyup.enter="runLdapTest"
                />
                <button
                  class="btn btn--primary ldap-test__btn"
                  type="button"
                  :disabled="!ldapTestUser.trim() || ldapTesting"
                  @click="runLdapTest"
                >
                  {{ ldapTesting ? t('settings.auth.ldap.testing') : t('settings.auth.ldap.testTitle') }}
                </button>
              </div>

              <div v-if="ldapTestResult" class="ldap-result">
                <div
                  class="ldap-result__banner"
                  :class="ldapTestResult.success ? 'ldap-result__banner--ok' : 'ldap-result__banner--err'"
                >
                  {{ ldapTestResult.message }}
                </div>

                <template v-if="ldapTestResult.success">
                  <div class="ldap-result__section">
                    <p class="ldap-result__subtitle">{{ t('settings.auth.ldap.allAttributes', { count: ldapTestResult.attributes.length }) }}</p>
                    <div class="ldap-result__tags">
                      <code
                        v-for="attr in ldapTestResult.attributes"
                        :key="attr.name"
                        class="ldap-result__tag"
                      >
                        <strong>{{ attr.name }}</strong>
                        <span>{{ attr.values.join(', ') || t('common.emptyDash') }}</span>
                      </code>
                    </div>
                  </div>

                  <div class="ldap-result__section">
                    <p class="ldap-result__subtitle">{{ t('settings.auth.ldap.currentMapping') }}</p>
                    <div class="ldap-result__table">
                      <div
                        v-for="(_v, field) in ldapTestResult.mapped"
                        :key="field"
                        class="ldap-result__row"
                      >
                        <span class="ldap-result__field">{{ platformFieldLabel(field) }}</span>
                        <span class="ldap-result__arrow">=</span>
                        <code class="ldap-result__value">{{ _v || t('common.emptyDash') }}</code>
                      </div>
                      <div
                        v-if="Object.keys(ldapTestResult.mapped).length === 0"
                        class="ldap-result__row ldap-result__row--empty"
                      >
                        {{ t('settings.auth.ldap.noMapping') }}
                    </div>
                  </div>
                  </div>
                </template>
              </div>
            </div>
          </template>

          <template v-if="authMethod === 'oauth'">
            <label class="field">
              <span class="field__label">{{ t('settings.auth.oauth.authorizeUrlLabel') }}</span>
              <input
                v-model="oauthForm.authorizeUrl"
                class="field__input"
                type="text"
                placeholder="https://oauth.example.edu.cn/authorize"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.oauth.tokenUrlLabel') }}</span>
              <input
                v-model="oauthForm.tokenUrl"
                class="field__input"
                type="text"
                placeholder="https://oauth.example.edu.cn/token"
              />
            </label>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.oauth.userinfoUrlLabel') }}</span>
              <input
                v-model="oauthForm.userinfoUrl"
                class="field__input"
                type="text"
                placeholder="https://oauth.example.edu.cn/userinfo"
              />
            </label>
            <div class="field-row">
              <label class="field">
                <span class="field__label">Client ID</span>
                <input
                  v-model="oauthForm.clientId"
                  class="field__input"
                  type="text"
                  :placeholder="t('settings.auth.oauth.clientIdPlaceholder')"
                />
              </label>
              <label class="field">
                <span class="field__label">Client Secret</span>
                <input
                  v-model="oauthForm.clientSecret"
                  class="field__input"
                  type="password"
                  :placeholder="t('settings.auth.oauth.clientSecretPlaceholder')"
                />
              </label>
            </div>
            <label class="field">
              <span class="field__label">{{ t('settings.auth.oauth.redirectUrlLabel') }}</span>
              <input
                v-model="oauthForm.redirectUrl"
                class="field__input"
                type="text"
                placeholder="https://mcp.example.edu.cn/login/oauth/callback"
              />
            </label>
            <label class="field">
              <span class="field__label">Scope</span>
              <input
                v-model="oauthForm.scope"
                class="field__input"
                type="text"
                placeholder="openid profile"
              />
            </label>
          </template>

          <div class="collapse__actions">
            <span v-if="tips.auth" class="save-tip">{{ tips.auth }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'auth'" @click="saveAuthSettings">
              {{ saving === 'auth' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
          <div
            v-if="ldapTestResult && ldapTestResult.success"
            class="ldap-result__hint ldap-result__hint--under"
          >
            <p>{{ t('settings.auth.ldap.saveHint1') }}<span class="ldap-result__em">{{ t('settings.auth.ldap.saveHint2') }}</span>{{ t('settings.auth.ldap.saveHint3') }}</p>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="userOpsOpen = !userOpsOpen">
          <span class="collapse__title">{{ t('settings.userOps.title') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': userOpsOpen }">▾</span>
        </button>
        <div v-show="userOpsOpen" class="collapse__body">
          <div class="user-ops-grid">
            <div class="user-ops-col">
              <label class="field">
                <span class="field__label">{{ t('settings.userOps.maxAccessKeysLabel') }}</span>
                <input
                  v-model.number="maxAccessKeys"
                  class="field__input"
                  type="number"
                  min="1"
                  placeholder="5"
                />
              </label>

              <label class="field">
                <span class="field__label">{{ t('settings.userOps.accessKeyCronLabel') }}</span>
                <input
                  v-model="accessKeyCron"
                  class="field__input"
                  type="text"
                  placeholder="0 */6 * * *"
                />
                <span class="field__help-text">{{ t('settings.userOps.accessKeyCronHint') }}</span>
              </label>
            </div>

            <div class="user-ops-col">
              <p class="field__section-title">{{ t('settings.userOps.roleRuleTitle') }}</p>
              <label class="field">
                <span class="field__label">{{ t('settings.userOps.filterAttrLabel') }}</span>
                <select v-model="filterAttribute" class="field__input">
                  <option v-for="a in ROLE_ATTRS" :key="a.key" :value="a.key">
                    {{ t(a.labelKey) }}
                  </option>
                </select>
                <span class="field__help-text">{{ t('settings.userOps.filterAttrHint') }}</span>
              </label>

              <div
                v-for="(rule, i) in roleRules"
                :key="i"
                class="role-rule-row"
              >
                <span class="role-rule-row__index">{{ t('settings.userOps.conditionIndex', { n: i + 1 }) }}</span>
                <span class="role-rule-row__label">{{ t('settings.userOps.satisfies') }}</span>
                <input
                  v-model="rule.pattern"
                  class="field__input role-rule-row__input"
                  type="text"
                  :placeholder="t('settings.userOps.patternPlaceholder')"
                />
                <span class="role-rule-row__label">{{ t('settings.userOps.assignRole') }}</span>
                <select v-model="rule.roleId" class="field__input role-rule-row__select">
                  <option value="" disabled>{{ t('settings.userOps.selectRole') }}</option>
                  <option v-for="r in roles" :key="r.id" :value="r.id">
                    {{ r.name }}
                  </option>
                </select>
                <button
                  class="role-rule-row__remove"
                  type="button"
                  :title="t('settings.userOps.removeCondition')"
                  @click="removeRoleRule(i)"
                >
                  ×
                </button>
              </div>

              <button class="role-rule-row__add" type="button" @click="addRoleRule">
                {{ t('settings.userOps.addCondition') }}
              </button>
            </div>
          </div>

          <div class="collapse__actions">
            <span v-if="tips.user_ops" class="save-tip">{{ tips.user_ops }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'user_ops'" @click="saveUserOpsSettings">
              {{ saving === 'user_ops' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="quickAccessOpen = !quickAccessOpen">
          <span class="collapse__title">{{ t('nav.quickaccess') }}</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': quickAccessOpen }">▾</span>
        </button>
        <div v-show="quickAccessOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">{{ t('settings.quickAccess.supportedClients') }}</span>
            <span class="field__help-text">{{ t('settings.quickAccess.supportedClientsHint') }}</span>
            <div class="qa-clients">
              <label
                v-for="c in QUICK_ACCESS_CLIENTS"
                :key="c.key"
                class="qa-client"
              >
                <input
                  type="checkbox"
                  :checked="quickAccessForm.enabledClients.includes(c.key)"
                  @change="toggleQuickAccessClient(c.key)"
                />
                <span>{{ c.label }}</span>
              </label>
            </div>
          </div>

          <div class="field">
            <span class="field__label">{{ t('settings.quickAccess.schemeLabel') }}</span>
            <span class="field__help-text">{{ t('settings.quickAccess.schemeHint') }}</span>
            <div class="seg">
              <button
                type="button"
                :class="['seg__item', { 'seg__item--active': quickAccessForm.scheme === 'http' }]"
                @click="quickAccessForm.scheme = 'http'"
              >HTTP</button>
              <button
                type="button"
                :class="['seg__item', { 'seg__item--active': quickAccessForm.scheme === 'https' }]"
                @click="quickAccessForm.scheme = 'https'"
              >HTTPS</button>
            </div>
          </div>

          <div class="collapse__actions">
            <span v-if="tips.quick_access" class="save-tip">{{ tips.quick_access }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'quick_access'" @click="saveQuickAccessSettings">
              {{ saving === 'quick_access' ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </section>
</template>

<style scoped>
.settings {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.settings__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.settings__tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid var(--border);
}

.settings__tab {
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  cursor: pointer;
  transition: color 0.2s, border-color 0.2s;
}

.settings__tab:hover {
  color: var(--xauat-blue);
}

.settings__tab--active {
  color: var(--xauat-blue);
  border-bottom-color: var(--xauat-blue);
}

.settings__panel {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.collapse {
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--surface);
  overflow: hidden;
}

.collapse__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: 16px 20px;
  background: transparent;
  border: none;
  cursor: pointer;
}

.collapse__title {
  font-size: 15px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.collapse__arrow {
  font-size: 14px;
  color: var(--text-muted);
  transition: transform 0.2s;
}

.collapse__arrow--open {
  transform: rotate(180deg);
}

.collapse__body {
  display: flex;
  flex-direction: column;
  gap: 20px;
  padding: 4px 20px 20px;
  border-top: 1px solid var(--border);
  padding-top: 20px;
}

.collapse__actions {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 12px;
}

.save-tip {
  font-size: 13px;
  color: var(--text-muted);
}

.field {
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 480px;
}

.field-row {
  display: flex;
  gap: 16px;
  max-width: 480px;
}

.field-row .field {
  flex: 1;
}

.field__label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.field__input {
  padding: 10px 14px;
  font-size: 15px;
  color: var(--text);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.field__input:focus {
  border-color: var(--xauat-blue-light);
  box-shadow: 0 0 0 3px rgba(30, 95, 176, 0.12);
}

.field__help-text {
  display: block;
  margin-top: 4px;
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.5;
}

.field__textarea {
  resize: vertical;
  min-height: 100px;
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
  line-height: 1.6;
}

.field__muted {
  font-size: 13px;
  color: var(--text-muted);
}

.ldap-notice {
  padding: 10px 14px;
  font-size: 13px;
  line-height: 1.6;
  color: #92400e;
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 10px;
}

.field__help {
  position: relative;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 15px;
  height: 15px;
  margin-left: 4px;
  font-size: 11px;
  font-weight: 700;
  color: var(--text-muted);
  background: rgba(10, 61, 122, 0.08);
  border: 1px solid var(--border);
  border-radius: 50%;
  cursor: help;
  vertical-align: middle;
}

.field__tooltip {
  position: absolute;
  bottom: calc(100% + 8px);
  left: 50%;
  transform: translateX(-50%);
  width: 280px;
  padding: 10px 12px;
  font-size: 12px;
  font-weight: 400;
  line-height: 1.6;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  box-shadow: var(--shadow);
  opacity: 0;
  visibility: hidden;
  transition: opacity 0.2s, visibility 0.2s;
  z-index: 10;
}

.field__tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-top-color: var(--border);
}

.field__tooltip code {
  padding: 1px 5px;
  font-size: 11px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.1);
  border-radius: 4px;
}

.field__help:hover .field__tooltip {
  opacity: 1;
  visibility: visible;
}

.field__section-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  margin: 0 0 12px;
  padding-top: 8px;
  border-top: 1px solid var(--border);
}

.mapping-row {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 10px;
}

.mapping-row__select {
  width: 160px;
  flex-shrink: 0;
}

.mapping-row__arrow {
  color: var(--text-muted);
  font-size: 14px;
  flex-shrink: 0;
}

.mapping-row__input {
  flex: 1;
}

.mapping-row__remove {
  flex-shrink: 0;
  width: 28px;
  height: 28px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  color: var(--text-muted);
  background: none;
  border: 1px solid var(--border);
  border-radius: 6px;
  cursor: pointer;
  transition: color 0.2s, border-color 0.2s;
  line-height: 1;
  padding: 0;
}

.mapping-row__remove:hover {
  color: #d93025;
  border-color: #d93025;
}

.mapping-row__add {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: var(--xauat-blue);
  background: none;
  border: 1px dashed var(--border);
  border-radius: 8px;
  padding: 6px 14px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.mapping-row__add:hover {
  border-color: var(--xauat-blue);
  background: rgba(10, 61, 122, 0.04);
}

.addr-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.addr-row__input {
  flex: 1;
}

.addr-row__remove {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  line-height: 1;
  color: var(--text-muted);
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s, border-color 0.2s;
}

.addr-row__remove:hover {
  color: #d93025;
  border-color: #d93025;
}

.addr-row__add {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  align-self: flex-start;
  width: fit-content;
  font-size: 13px;
  color: var(--xauat-blue);
  background: none;
  border: 1px dashed var(--border);
  border-radius: 8px;
  padding: 6px 14px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.addr-row__add:hover {
  border-color: var(--xauat-blue);
  background: rgba(10, 61, 122, 0.04);
}

.api-gw-warning {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 12px 16px;
  margin-bottom: 16px;
  font-size: 13px;
  line-height: 1.6;
  color: #92400e;
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 10px;
}

.qa-clients {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.qa-client {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.qa-client:hover {
  border-color: var(--xauat-blue-light);
}

.qa-client input {
  width: 16px;
  height: 16px;
  accent-color: var(--xauat-blue);
  cursor: pointer;
}

.dialog-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.dialog {
  background: var(--surface);
  border-radius: var(--radius);
  padding: 32px;
  max-width: 520px;
  width: 90%;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: var(--shadow);
}

.dialog__title {
  font-size: 20px;
  font-weight: 700;
  color: var(--xauat-blue);
  margin-bottom: 20px;
  flex-shrink: 0;
}

.dialog__desc {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 16px;
  line-height: 1.7;
}

.dialog__desc code {
  padding: 1px 5px;
  font-size: 12px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.1);
  border-radius: 4px;
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
  flex-shrink: 0;
}

.kong-notice__list {
  margin: 0 0 16px;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  font-size: 14px;
  line-height: 1.7;
  color: var(--text);
}

.kong-notice__list strong {
  color: #d93025;
}

.api-gw-warning__icon {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  margin-top: 1px;
  font-size: 12px;
  font-weight: 700;
  color: #fff;
  background: #f59e0b;
  border-radius: 50%;
}

.switch {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.switch input {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
}

.switch__track {
  position: relative;
  width: 40px;
  height: 22px;
  border-radius: 999px;
  background: var(--border);
  transition: background 0.2s;
}

.switch__thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s;
}

.switch input:checked + .switch__track {
  background: var(--xauat-blue);
}

.switch input:checked + .switch__track .switch__thumb {
  transform: translateX(18px);
}

.switch__label {
  font-size: 13px;
  color: var(--text-muted);
}

.seg {
  display: inline-flex;
  border: 1px solid var(--border);
  border-radius: 999px;
  overflow: hidden;
  width: fit-content;
}

.seg__item {
  padding: 8px 22px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.seg__item + .seg__item {
  border-left: 1px solid var(--border);
}

.seg__item:hover {
  color: var(--xauat-blue);
}

.seg__item--active {
  color: #fff;
  background: var(--xauat-blue);
}

.seg__item--active:hover {
  color: #fff;
}

.btn {
  padding: 8px 20px;
  font-size: 14px;
  font-weight: 600;
  border-radius: 999px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.btn--primary {
  color: #fff;
  background: var(--xauat-blue);
}

.btn--primary:hover {
  background: var(--xauat-blue-light);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.ldap-test {
  margin-top: 16px;
}

.smtp-test {
  margin-top: 16px;
}

.smtp-test__row {
  display: flex;
  gap: 8px;
}

.smtp-test__input {
  flex: 1;
}

.smtp-test__btn {
  flex-shrink: 0;
  padding: 8px 16px;
  font-size: 13px;
}

.smtp-test__result {
  margin-top: 12px;
  padding: 8px 14px;
  font-size: 13px;
  border-radius: 8px;
}

.smtp-test__result--ok {
  color: #1a7d3a;
  background: #e8f5e9;
}

.smtp-test__result--err {
  color: #b71c1c;
  background: #fbe9e7;
}

.ldap-test__row {
  display: flex;
  gap: 8px;
}

.ldap-test__input {
  flex: 1;
}

.ldap-test__btn {
  flex-shrink: 0;
  padding: 8px 16px;
  font-size: 13px;
}

.ldap-result {
  margin-top: 12px;
}

.ldap-result__banner {
  padding: 8px 14px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 8px;
}

.ldap-result__banner--ok {
  color: #1a7d3a;
  background: #e8f5e9;
}

.ldap-result__banner--err {
  color: #b71c1c;
  background: #fbe9e7;
}

.ldap-result__section {
  margin-top: 14px;
}

.ldap-result__subtitle {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.ldap-result__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.ldap-result__tag {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  font-size: 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 6px;
}

.ldap-result__tag strong {
  color: var(--xauat-blue);
}

.ldap-result__tag span {
  color: var(--text-muted);
  max-width: 200px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.ldap-result__table {
  border: 1px solid var(--border);
  border-radius: 8px;
  overflow: hidden;
}

.ldap-result__row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 14px;
  font-size: 13px;
  border-bottom: 1px solid var(--border);
}

.ldap-result__row:last-child {
  border-bottom: none;
}

.ldap-result__row--empty {
  color: var(--text-muted);
  justify-content: center;
}

.ldap-result__field {
  width: 80px;
  flex-shrink: 0;
  color: var(--text);
  font-weight: 500;
}

.ldap-result__arrow {
  color: var(--text-muted);
  flex-shrink: 0;
}

.ldap-result__value {
  color: var(--xauat-blue);
  font-weight: 600;
}

.ldap-result__hint {
  margin-top: 14px;
  padding: 10px 14px;
  font-size: 12px;
  line-height: 1.6;
  color: #92400e;
  background: #fef3c7;
  border: 1px solid #f59e0b;
  border-radius: 8px;
}

.ldap-result__hint p {
  margin: 0;
}

.ldap-result__hint--under {
  margin-top: 0;
}

.ldap-result__em {
  font-weight: 700;
  color: #b45309;
}

.user-ops-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  align-items: start;
}

.user-ops-col {
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-width: 0;
}

.role-rule-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  padding: 12px 14px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
}

.role-rule-row__index {
  flex-shrink: 0;
  min-width: 52px;
  padding: 4px 10px;
  font-size: 12px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.08);
  border-radius: 999px;
  text-align: center;
}

.role-rule-row__label {
  flex-shrink: 0;
  font-size: 13px;
  color: var(--text-muted);
}

.role-rule-row__select {
  width: 140px;
  flex-shrink: 0;
}

.role-rule-row__input {
  flex: 1;
  min-width: 140px;
  font-family: 'SF Mono', 'Fira Code', 'Fira Mono', Menlo, Consolas, monospace;
}

.role-rule-row__remove {
  flex-shrink: 0;
  width: 32px;
  height: 32px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
  line-height: 1;
  color: var(--text-muted);
  background: none;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s, border-color 0.2s;
}

.role-rule-row__remove:hover {
  color: #d93025;
  border-color: #d93025;
}

.role-rule-row__add {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  align-self: flex-start;
  width: fit-content;
  font-size: 13px;
  color: var(--xauat-blue);
  background: none;
  border: 1px dashed var(--border);
  border-radius: 8px;
  padding: 6px 14px;
  cursor: pointer;
  transition: border-color 0.2s, background 0.2s;
}

.role-rule-row__add:hover {
  border-color: var(--xauat-blue);
  background: rgba(10, 61, 122, 0.04);
}

@media (max-width: 900px) {
  .user-ops-grid {
    grid-template-columns: 1fr;
  }
}
</style>

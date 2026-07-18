<script setup>
import { ref, onMounted } from 'vue'
import { fetchRoles, fetchSettings, saveSetting } from '../api.js'

const tab = ref('ops')

const logOpen = ref(true)
const smtpOpen = ref(true)
const authOpen = ref(true)
const userOpsOpen = ref(true)

const logForm = ref({
  syslogEnabled: false,
  logPath: '/var/log/mcp-plat',
  logLevel: 'info',
  logPrefix: 'mcp-plat',
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

const ldapForm = ref({
  host: '',
  port: 389,
  baseDn: '',
  bindDn: '',
  bindPassword: '',
  userFilter: ''
})

const oauthForm = ref({
  authorizeUrl: '',
  tokenUrl: '',
  userinfoUrl: '',
  clientId: '',
  clientSecret: '',
  redirectUrl: '',
  scope: ''
})

const roles = ref([])
const defaultRoleId = ref('')
const maxAccessKeys = ref(5)

const saving = ref('')
const tips = ref({})

async function saveSection(key, payload) {
  saving.value = key
  try {
    await saveSetting(key, payload)
    tips.value = { ...tips.value, [key]: '已保存' }
  } catch (e) {
    tips.value = { ...tips.value, [key]: e.message || '保存失败' }
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

function saveAuthSettings() {
  saveSection('auth', {
    method: authMethod.value,
    cas: casForm.value,
    ldap: ldapForm.value,
    oauth: oauthForm.value
  })
}

function saveUserOpsSettings() {
  saveSection('user_ops', {
    defaultRoleId: defaultRoleId.value === '' ? null : defaultRoleId.value,
    maxAccessKeys: maxAccessKeys.value
  })
}

function applySettings(data) {
  if (data.log) Object.assign(logForm.value, data.log)
  if (data.smtp) Object.assign(smtpForm.value, data.smtp)
  if (data.auth) {
    if (data.auth.method) authMethod.value = data.auth.method
    if (data.auth.cas) Object.assign(casForm.value, data.auth.cas)
    if (data.auth.ldap) Object.assign(ldapForm.value, data.auth.ldap)
    if (data.auth.oauth) Object.assign(oauthForm.value, data.auth.oauth)
  }
  if (data.user_ops) {
    if (data.user_ops.defaultRoleId != null) defaultRoleId.value = data.user_ops.defaultRoleId
    if (data.user_ops.maxAccessKeys != null) maxAccessKeys.value = data.user_ops.maxAccessKeys
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
      <h1 class="page__title">系统设置</h1>
    </div>

    <div class="settings__tabs">
      <button
        class="settings__tab"
        :class="{ 'settings__tab--active': tab === 'ops' }"
        type="button"
        @click="tab = 'ops'"
      >
        运维设置
      </button>
      <button
        class="settings__tab"
        :class="{ 'settings__tab--active': tab === 'operation' }"
        type="button"
        @click="tab = 'operation'"
      >
        运营设置
      </button>
    </div>

    <div v-if="tab === 'ops'" class="settings__panel">
      <div class="collapse">
        <button class="collapse__head" type="button" @click="logOpen = !logOpen">
          <span class="collapse__title">日志设置</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': logOpen }">▾</span>
        </button>
        <div v-show="logOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">启用 Syslog</span>
            <label class="switch">
              <input type="checkbox" v-model="logForm.syslogEnabled" />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ logForm.syslogEnabled ? '启用' : '禁用' }}</span>
            </label>
          </div>

          <template v-if="!logForm.syslogEnabled">
            <label class="field">
              <span class="field__label">日志路径</span>
              <input
                v-model="logForm.logPath"
                class="field__input"
                type="text"
                placeholder="请输入日志路径"
              />
            </label>
            <label class="field">
              <span class="field__label">日志级别</span>
              <select v-model="logForm.logLevel" class="field__input">
                <option v-for="lv in LOG_LEVELS" :key="lv" :value="lv">{{ lv }}</option>
              </select>
            </label>
            <label class="field">
              <span class="field__label">日志文件前缀</span>
              <input
                v-model="logForm.logPrefix"
                class="field__input"
                type="text"
                placeholder="请输入日志文件前缀"
              />
            </label>
          </template>

          <template v-else>
            <label class="field">
              <span class="field__label">Syslog 服务器地址</span>
              <input
                v-model="logForm.syslogHost"
                class="field__input"
                type="text"
                placeholder="请输入 Syslog 服务器地址"
              />
            </label>
            <div class="field-row">
              <label class="field">
                <span class="field__label">端口</span>
                <input
                  v-model.number="logForm.syslogPort"
                  class="field__input"
                  type="number"
                  placeholder="514"
                />
              </label>
              <label class="field">
                <span class="field__label">协议</span>
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
              {{ saving === 'log' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="smtpOpen = !smtpOpen">
          <span class="collapse__title">SMTP 设置</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': smtpOpen }">▾</span>
        </button>
        <div v-show="smtpOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">启用邮件通知</span>
            <label class="switch">
              <input type="checkbox" v-model="smtpForm.enabled" />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ smtpForm.enabled ? '启用' : '禁用' }}</span>
            </label>
          </div>

          <template v-if="smtpForm.enabled">
            <div class="field-row">
              <label class="field">
                <span class="field__label">SMTP 服务器地址</span>
                <input
                  v-model="smtpForm.host"
                  class="field__input"
                  type="text"
                  placeholder="smtp.example.edu.cn"
                />
              </label>
              <label class="field">
                <span class="field__label">端口</span>
                <input
                  v-model.number="smtpForm.port"
                  class="field__input"
                  type="number"
                  placeholder="465"
                />
              </label>
            </div>
            <label class="field">
              <span class="field__label">加密方式</span>
              <select v-model="smtpForm.encryption" class="field__input">
                <option value="none">无</option>
                <option value="ssl">SSL/TLS</option>
                <option value="starttls">STARTTLS</option>
              </select>
            </label>
            <label class="field">
              <span class="field__label">用户名</span>
              <input
                v-model="smtpForm.username"
                class="field__input"
                type="text"
                placeholder="请输入 SMTP 用户名"
              />
            </label>
            <label class="field">
              <span class="field__label">密码</span>
              <input
                v-model="smtpForm.password"
                class="field__input"
                type="password"
                placeholder="请输入 SMTP 密码"
              />
            </label>
            <div class="field-row">
              <label class="field">
                <span class="field__label">发件人地址</span>
                <input
                  v-model="smtpForm.fromAddress"
                  class="field__input"
                  type="text"
                  placeholder="noreply@example.edu.cn"
                />
              </label>
              <label class="field">
                <span class="field__label">发件人名称</span>
                <input
                  v-model="smtpForm.fromName"
                  class="field__input"
                  type="text"
                  placeholder="MCP 服务平台"
                />
              </label>
            </div>
          </template>

          <div class="collapse__actions">
            <span v-if="tips.smtp" class="save-tip">{{ tips.smtp }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'smtp'" @click="saveSmtpSettings">
              {{ saving === 'smtp' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'operation'" class="settings__panel">
      <div class="collapse">
        <button class="collapse__head" type="button" @click="authOpen = !authOpen">
          <span class="collapse__title">用户认证</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': authOpen }">▾</span>
        </button>
        <div v-show="authOpen" class="collapse__body">
          <div class="field">
            <span class="field__label">认证方式</span>
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
              <span class="field__label">CAS 服务器地址</span>
              <input
                v-model="casForm.serverUrl"
                class="field__input"
                type="text"
                placeholder="https://cas.example.edu.cn/cas"
              />
            </label>
            <label class="field">
              <span class="field__label">服务回调地址</span>
              <input
                v-model="casForm.serviceUrl"
                class="field__input"
                type="text"
                placeholder="https://mcp.example.edu.cn/login/cas"
              />
            </label>
            <label class="field">
              <span class="field__label">协议版本</span>
              <select v-model="casForm.version" class="field__input">
                <option value="2.0">2.0</option>
                <option value="3.0">3.0</option>
              </select>
            </label>
          </template>

          <template v-if="authMethod === 'ldap'">
            <div class="field-row">
              <label class="field">
                <span class="field__label">服务器地址</span>
                <input
                  v-model="ldapForm.host"
                  class="field__input"
                  type="text"
                  placeholder="ldap.example.edu.cn"
                />
              </label>
              <label class="field">
                <span class="field__label">端口</span>
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
              <span class="field__label">Bind 密码</span>
              <input
                v-model="ldapForm.bindPassword"
                class="field__input"
                type="password"
                placeholder="请输入 Bind 密码"
              />
            </label>
            <label class="field">
              <span class="field__label">用户过滤器</span>
              <input
                v-model="ldapForm.userFilter"
                class="field__input"
                type="text"
                placeholder="(uid=%s)"
              />
            </label>
          </template>

          <template v-if="authMethod === 'oauth'">
            <label class="field">
              <span class="field__label">授权地址</span>
              <input
                v-model="oauthForm.authorizeUrl"
                class="field__input"
                type="text"
                placeholder="https://oauth.example.edu.cn/authorize"
              />
            </label>
            <label class="field">
              <span class="field__label">Token 地址</span>
              <input
                v-model="oauthForm.tokenUrl"
                class="field__input"
                type="text"
                placeholder="https://oauth.example.edu.cn/token"
              />
            </label>
            <label class="field">
              <span class="field__label">用户信息地址</span>
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
                  placeholder="请输入 Client ID"
                />
              </label>
              <label class="field">
                <span class="field__label">Client Secret</span>
                <input
                  v-model="oauthForm.clientSecret"
                  class="field__input"
                  type="password"
                  placeholder="请输入 Client Secret"
                />
              </label>
            </div>
            <label class="field">
              <span class="field__label">回调地址</span>
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
              {{ saving === 'auth' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="userOpsOpen = !userOpsOpen">
          <span class="collapse__title">用户运营</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': userOpsOpen }">▾</span>
        </button>
        <div v-show="userOpsOpen" class="collapse__body">
          <label class="field">
            <span class="field__label">默认新用户角色</span>
            <select v-model="defaultRoleId" class="field__input">
              <option value="" disabled>请选择角色</option>
              <option v-for="r in roles" :key="r.id" :value="r.id">
                {{ r.name }}
              </option>
            </select>
            <span v-if="roles.length === 0" class="field__muted">暂无可选角色</span>
          </label>

          <label class="field">
            <span class="field__label">每用户最大 AccessKey 数量</span>
            <input
              v-model.number="maxAccessKeys"
              class="field__input"
              type="number"
              min="1"
              placeholder="5"
            />
          </label>

          <div class="collapse__actions">
            <span v-if="tips.user_ops" class="save-tip">{{ tips.user_ops }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'user_ops'" @click="saveUserOpsSettings">
              {{ saving === 'user_ops' ? '保存中…' : '保存' }}
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

.field__muted {
  font-size: 13px;
  color: var(--text-muted);
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
</style>

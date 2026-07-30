<script setup>

// Author: deepseek-v4-pro / opencode
import { ref, onMounted } from 'vue'
import { fetchRoles, fetchSettings, saveSetting } from '../api.js'
import { applyPlatform } from '../stores/settings.js'

const tab = ref('ops')

const logOpen = ref(true)
const smtpOpen = ref(true)
const apiGwOpen = ref(true)
const networkSecurityOpen = ref(true)
const authOpen = ref(true)
const userOpsOpen = ref(true)
const platformOpen = ref(true)

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

const apiGwForm = ref({
  provider: 'apisix',
  adminUrl: '',
  defaultPublishDomain: '',
  adminKey: '',
  accesskeyHeader: ''
})

const API_GW_PROVIDERS = [
  { key: 'apisix', label: 'Apache APISIX' },
  { key: 'kong', label: 'Kong' },
  { key: 'tyk', label: 'Tyk' }
]

const networkSecurityForm = ref({
  allowlist: ''
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

const platformForm = ref({
  name: '某某大学',
  logoUrl: '',
  siteUrl: ''
})

const roles = ref([])
const defaultRoleId = ref('')
const maxAccessKeys = ref(5)

const apiGwConfigured = ref(false)

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

function savePlatformSettings() {
  saveSection('platform', platformForm.value)
  applyPlatform(platformForm.value)
}

function applySettings(data) {
  if (data.log) Object.assign(logForm.value, data.log)
  if (data.smtp) Object.assign(smtpForm.value, data.smtp)
  if (data.api_gateway) {
    apiGwConfigured.value = data.api_gateway.configured || false
  }
  if (data.network_security && data.network_security.allowlist) {
    networkSecurityForm.value.allowlist = data.network_security.allowlist.join('\n')
  }
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
  if (data.platform) {
    Object.assign(platformForm.value, data.platform)
    applyPlatform(data.platform)
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

      <div class="collapse">
        <button class="collapse__head" type="button" @click="apiGwOpen = !apiGwOpen">
          <span class="collapse__title">API 网关设置</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': apiGwOpen }">▾</span>
        </button>
        <div v-show="apiGwOpen" class="collapse__body">
          <div v-if="apiGwConfigured" class="api-gw-warning">
            <span class="api-gw-warning__icon">!</span>
            已检测到当前系统配置了 API 网关，修改保存后会覆盖当前配置。
          </div>
          <label class="field">
            <span class="field__label">API 网关</span>
            <select v-model="apiGwForm.provider" class="field__input">
              <option
                v-for="gw in API_GW_PROVIDERS"
                :key="gw.key"
                :value="gw.key"
              >
                {{ gw.label }}
              </option>
            </select>
          </label>

          <label class="field">
            <span class="field__label">
              网关 Admin API 地址
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  填写网关 Admin API 的基础地址，不要附带路径。例如填写
                  <code>http://127.0.0.1:9180</code> 而非 <code>http://127.0.0.1:9180/apisix/admin</code>。
                </span>
              </span>
            </span>
            <input
              v-model="apiGwForm.adminUrl"
              class="field__input"
              type="text"
              placeholder="http://127.0.0.1:9180"
            />
          </label>

          <label class="field">
            <span class="field__label">默认发布域名</span>
            <input
              v-model="apiGwForm.defaultPublishDomain"
              class="field__input"
              type="text"
              placeholder="例如 mcp.xauat.edu.cn"
            />
          </label>

          <label class="field">
            <span class="field__label">Admin API Key</span>
            <input
              v-model="apiGwForm.adminKey"
              class="field__input"
              type="password"
              placeholder="请输入管理员 API Key"
            />
          </label>

          <label class="field">
            <span class="field__label">
              Access Key 认证 gRPC 地址
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  对应 accesskey-auth-server 的 gRPC 监听地址，例如
                  <code>:9090</code> 或 <code>127.0.0.1:9090</code>。
                  发布时若启用认证，将下发 accesskey_verify 插件并指向该地址。
                </span>
              </span>
            </span>
            <input
              v-model="apiGwForm.authGrpcAddr"
              class="field__input"
              type="text"
              placeholder=":9090"
            />
          </label>

          <label class="field">
            <span class="field__label">
              Access Key Header 名称
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  指定请求中携带 Access Key 的 HTTP Header 名称。<br />
                  发布路由时若已配置则使用自定义 Header，未配置则默认使用
                  <code>X-Access-Key</code>。
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
              {{ saving === 'api_gateway' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>

      <div class="collapse">
        <button class="collapse__head" type="button" @click="networkSecurityOpen = !networkSecurityOpen">
          <span class="collapse__title">网络安全</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': networkSecurityOpen }">▾</span>
        </button>
        <div v-show="networkSecurityOpen" class="collapse__body">
          <label class="field">
            <span class="field__label">
              允许连接的内网 CIDR
              <span class="field__help">
                ?
                <span class="field__tooltip">
                  每行一个 CIDR 地址段，例如 <code>10.0.0.0/8</code>。<br />
                  配置后，获取工具列表时将允许连接这些内网地址。<br />
                  如需连接 MCP 服务器内网 IP，在此添加对应的地址段即可。
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
              {{ saving === 'network_security' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="tab === 'operation'" class="settings__panel">
      <div class="collapse">
        <button class="collapse__head" type="button" @click="platformOpen = !platformOpen">
          <span class="collapse__title">平台设置</span>
          <span class="collapse__arrow" :class="{ 'collapse__arrow--open': platformOpen }">▾</span>
        </button>
        <div v-show="platformOpen" class="collapse__body">
          <label class="field">
            <span class="field__label">平台名称</span>
            <input
              v-model="platformForm.name"
              class="field__input"
              type="text"
              placeholder="请输入平台名称，如：某某大学"
            />
          </label>

          <label class="field">
            <span class="field__label">Logo 图片地址</span>
            <input
              v-model="platformForm.logoUrl"
              class="field__input"
              type="text"
              placeholder="请输入 Logo 图片的 URL 地址"
            />
          </label>

          <label class="field">
            <span class="field__label">跳转链接</span>
            <input
              v-model="platformForm.siteUrl"
              class="field__input"
              type="text"
              placeholder="点击 Logo 时跳转到的网址，如：https://www.example.edu.cn"
            />
          </label>

          <div class="collapse__actions">
            <span v-if="tips.platform" class="save-tip">{{ tips.platform }}</span>
            <button class="btn btn--primary" type="button" :disabled="saving === 'platform'" @click="savePlatformSettings">
              {{ saving === 'platform' ? '保存中…' : '保存' }}
            </button>
          </div>
        </div>
      </div>

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
</style>

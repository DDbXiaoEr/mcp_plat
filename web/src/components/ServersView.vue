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
import { fetchServers, createServer, fetchServerTools, publishServers, setServersMaintenance, deleteServer, updateServer, fetchGatewayStatus } from '../api.js'

const { t } = useI18n()

const servers = ref([])
const loading = ref(true)
const gatewayConfigured = ref(false)
const gatewayProvider = ref('apisix')

const isKongGateway = computed(() => gatewayProvider.value === 'kong')

const editingDesc = ref(false)
const editDescValue = ref('')
const savingDesc = ref(false)

onMounted(async () => {
  await Promise.all([
    loadServers(),
    loadGatewaySettings()
  ])
  loading.value = false
})

async function loadGatewaySettings() {
  try {
    const data = await fetchGatewayStatus()
    gatewayConfigured.value = data.configured
    gatewayProvider.value = data.provider || 'apisix'
  } catch (e) {
    console.error('加载 API 网关状态失败:', e)
  }
}

async function loadServers() {
  try {
    const data = await fetchServers()
    servers.value = data.map((s) => ({
      ...s,
      tools: parseTools(s.tools)
    }))
  } catch {
    // ignore load error
  }
}

function parseTools(raw) {
  if (!raw) return []
  if (Array.isArray(raw)) return raw.map(normalizeTool)
  try {
    const parsed = JSON.parse(raw)
    if (!parsed) return []
    return (Array.isArray(parsed) ? parsed : []).map(normalizeTool)
  } catch { return [] }
}

function normalizeTool(t) {
  return typeof t === 'string' ? { name: t, description: '' } : t
}

function splitAddresses(text) {
  return (text || '')
    .split(/[,，\n]/)
    .map((s) => s.trim())
    .filter(Boolean)
}

function formatServiceAddresses(raw) {
  if (!raw) return ''
  if (raw.startsWith('[')) {
    try {
      const arr = JSON.parse(raw)
      if (Array.isArray(arr)) return arr.join('\n')
    } catch { /**/ }
  }
  return splitAddresses(raw).join('\n')
}

function packServiceAddresses(text) {
  return JSON.stringify(splitAddresses(text))
}

const PROTOCOLS = ['SSE', 'Streamable HTTP']
const PROTOCOL_VERSIONS = ['2026-07-28', '2025-06-18', '2025-03-26']

const selectedId = ref('')
const selected = ref(null)

function select(server) {
  selectedId.value = server.id
  selected.value = server
}

function openDeleteConfirm(server) {
  deleteTarget.value = server
  showingDeleteConfirm.value = true
}

async function removeServer() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await deleteServer(deleteTarget.value.id)
    servers.value = servers.value.filter((s) => s.id !== deleteTarget.value.id)
    if (selectedId.value === deleteTarget.value.id) {
      selectedId.value = ''
      selected.value = null
    }
    showingDeleteConfirm.value = false
    deleteTarget.value = null
    alert(t('servers.alert.deleteSuccess'))
  } catch {
    // ignore
  }
  deleting.value = false
}

const selectedTools = computed(() => parseTools(selected.value?.tools))

function startEditDesc() {
  editDescValue.value = selected.value?.description || ''
  editingDesc.value = true
}

async function saveDescription() {
  if (!selected.value) return
  savingDesc.value = true
  try {
    await updateServer(selected.value.id, { description: editDescValue.value })
    selected.value.description = editDescValue.value
    editingDesc.value = false
  } catch {
    // ignore
  }
  savingDesc.value = false
}

function cancelEditDesc() {
  editingDesc.value = false
}

const showingCreate = ref(false)
const createForm = ref({ name: '', address: '', service_address: '', department: '', protocol: 'SSE', protocol_version: '2026-07-28', tools: '', description: '' })
const fetchingTools = ref(false)
const fetchToolsError = ref('')
const useHttps = ref(false)

const serviceAddressList = computed(() => splitAddresses(createForm.value.service_address))

const fetchAddressValue = ref('')
const customAddressInput = ref('')

const isCustomFetchAddress = computed(() => fetchAddressValue.value === '__custom__')

function randomizeFetchAddress() {
  const list = serviceAddressList.value
  fetchAddressValue.value = list.length > 0 ? list[Math.floor(Math.random() * list.length)] : ''
  customAddressInput.value = ''
}

watch(() => createForm.value.service_address, () => {
  if (!serviceAddressList.value.includes(fetchAddressValue.value) && fetchAddressValue.value !== '__custom__') {
    randomizeFetchAddress()
  }
})

function openCreate() {
  createForm.value = { name: '', address: '', service_address: '', department: '', protocol: 'SSE', protocol_version: '2026-07-28', tools: '', description: '' }
  fetchToolsError.value = ''
  useHttps.value = false
  fetchAddressValue.value = ''
  customAddressInput.value = ''
  showingCreate.value = true
}

async function fetchTools() {
  const host = (isCustomFetchAddress.value ? customAddressInput.value : fetchAddressValue.value).trim()
  const path = createForm.value.address.trim()
  fetchToolsError.value = ''
  if (!host) {
    fetchToolsError.value = t('servers.fetch.requireAddress')
    return
  }
  fetchingTools.value = true
  try {
    const scheme = useHttps.value ? 'https://' : 'http://'
    const address = scheme + host.replace(/\/+$/, '') + (path.startsWith('/') ? path : '/' + path)
    const data = await fetchServerTools({
      address,
      protocol: createForm.value.protocol,
      protocol_version: createForm.value.protocol_version
    })
    createForm.value.tools = JSON.stringify(data.tools || [])
  } catch (err) {
    fetchToolsError.value = err.message || t('servers.fetch.failed')
  }
  fetchingTools.value = false
}

async function confirmCreate() {
  const name = createForm.value.name.trim()
  if (!name) return
  const raw = createForm.value.tools.trim()
  let tools
  try {
    const parsed = JSON.parse(raw)
    tools = Array.isArray(parsed) ? parsed.map(normalizeTool) : []
  } catch {
    tools = raw
      .split(/[,，\n]/)
      .map((t) => t.trim())
      .filter(Boolean)
      .map((name) => ({ name, description: '' }))
  }
  try {
    const server = await createServer({
      name,
      address: createForm.value.address.trim(),
      service_address: packServiceAddresses(createForm.value.service_address),
      department: createForm.value.department.trim(),
      protocol: createForm.value.protocol,
      protocol_version: createForm.value.protocol_version,
      tools: JSON.stringify(tools),
      description: createForm.value.description.trim()
    })
    server.tools = tools
    servers.value.push(server)
    showingCreate.value = false
    select(server)
  } catch {
    // ignore
  }
}

const showingDeleteConfirm = ref(false)
const deleting = ref(false)
const deleteTarget = ref(null)

const showingPublish = ref(false)
const publishSelected = ref([])
const showingPublishConfirm = ref(false)
const publishing = ref(false)
const enableAuth = ref(false)
const enableAuditLog = ref(false)
const publishSearchLeft = ref('')
const publishSearchRight = ref('')
const publishSelectedIds = ref([])

const publishTargets = computed(() =>
  servers.value.filter((s) => publishSelected.value.includes(s.id))
)

const leftPublishServers = computed(() => {
  const list = servers.value.filter((s) => !publishSelected.value.includes(s.id))
  const kw = publishSearchLeft.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((s) =>
    (s.name || '').toLowerCase().includes(kw) || (s.id || '').toLowerCase().includes(kw)
  )
})

const rightPublishServers = computed(() => {
  const list = servers.value.filter((s) => publishSelected.value.includes(s.id))
  const kw = publishSearchRight.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((s) =>
    (s.name || '').toLowerCase().includes(kw) || (s.id || '').toLowerCase().includes(kw)
  )
})

const isAllLeftPublishChecked = computed(() => {
  if (leftPublishServers.value.length === 0) return false
  return leftPublishServers.value.every((s) => publishSelectedIds.value.includes(s.id))
})

const isAllRightPublishChecked = computed(() => {
  if (rightPublishServers.value.length === 0) return false
  return rightPublishServers.value.every((s) => publishSelectedIds.value.includes(s.id))
})

function openPublish() {
  publishSelected.value = []
  publishSelectedIds.value = []
  publishSearchLeft.value = ''
  publishSearchRight.value = ''
  enableAuth.value = false
  enableAuditLog.value = false
  showingPublish.value = true
}

function togglePublishSelect(id) {
  const idx = publishSelectedIds.value.indexOf(id)
  if (idx >= 0) publishSelectedIds.value.splice(idx, 1)
  else publishSelectedIds.value.push(id)
}

function togglePublishSelectAll(list) {
  const ids = list.map((s) => s.id)
  if (ids.every((id) => publishSelectedIds.value.includes(id))) {
    publishSelectedIds.value = publishSelectedIds.value.filter((id) => !ids.includes(id))
  } else {
    ids.forEach((id) => {
      if (!publishSelectedIds.value.includes(id)) publishSelectedIds.value.push(id)
    })
  }
}

function movePublishRight() {
  const selected = publishSelectedIds.value.filter((id) =>
    leftPublishServers.value.some((s) => s.id === id)
  )
  selected.forEach((id) => {
    if (!publishSelected.value.includes(id)) publishSelected.value.push(id)
  })
  publishSelectedIds.value = publishSelectedIds.value.filter((id) => !selected.includes(id))
}

function movePublishLeft() {
  const selected = publishSelectedIds.value.filter((id) =>
    rightPublishServers.value.some((s) => s.id === id)
  )
  publishSelected.value = publishSelected.value.filter((id) => !selected.includes(id))
  publishSelectedIds.value = publishSelectedIds.value.filter((id) => !selected.includes(id))
}

function confirmPublish() {
  showingPublish.value = false
  showingPublishConfirm.value = true
}

function backToPublish() {
  showingPublishConfirm.value = false
  showingPublish.value = true
}

async function executePublish() {
  publishing.value = true
  try {
    await publishServers(publishSelected.value, enableAuth.value, '', enableAuditLog.value)
    alert(t('servers.alert.publishSuccess'))
    showingPublishConfirm.value = false
    await loadServers()
  } catch (e) {
    alert(e.message || t('servers.alert.publishFailed'))
  }
  publishing.value = false
}

const showingMaintenance = ref(false)
const showingMaintenanceConfirm = ref(false)
const maintaining = ref(false)
const initialMaintenanceIds = ref([])
const maintenanceSelected = ref([])
const maintenanceSelectedIds = ref([])
const maintenanceSearchLeft = ref('')
const maintenanceSearchRight = ref('')

const leftMaintenanceServers = computed(() => {
  const list = servers.value.filter((s) =>
    (s.status === 'published' || s.status === 'maintenance') && !maintenanceSelected.value.includes(s.id)
  )
  const kw = maintenanceSearchLeft.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((s) =>
    (s.name || '').toLowerCase().includes(kw) || (s.id || '').toLowerCase().includes(kw)
  )
})

const rightMaintenanceServers = computed(() => {
  const list = servers.value.filter((s) => maintenanceSelected.value.includes(s.id))
  const kw = maintenanceSearchRight.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((s) =>
    (s.name || '').toLowerCase().includes(kw) || (s.id || '').toLowerCase().includes(kw)
  )
})

const isAllLeftMaintenanceChecked = computed(() => {
  if (leftMaintenanceServers.value.length === 0) return false
  return leftMaintenanceServers.value.every((s) => maintenanceSelectedIds.value.includes(s.id))
})

const isAllRightMaintenanceChecked = computed(() => {
  if (rightMaintenanceServers.value.length === 0) return false
  return rightMaintenanceServers.value.every((s) => maintenanceSelectedIds.value.includes(s.id))
})

const maintenanceTargets = computed(() => {
  const current = maintenanceSelected.value
  const initial = initialMaintenanceIds.value
  const enter = current.filter((id) => !initial.includes(id))
  const restore = initial.filter((id) => !current.includes(id))
  return {
    enter: servers.value.filter((s) => enter.includes(s.id)),
    restore: servers.value.filter((s) => restore.includes(s.id))
  }
})

const maintenanceConfirmLabel = computed(() => {
  const { enter, restore } = maintenanceTargets.value
  if (enter.length > 0 && restore.length > 0) return t('servers.maintenance.confirmBoth')
  if (restore.length > 0) return t('servers.maintenance.confirmCancel')
  return t('servers.maintenance.confirm')
})

function openMaintenance() {
  initialMaintenanceIds.value = servers.value.filter((s) => s.status === 'maintenance').map((s) => s.id)
  maintenanceSelected.value = [...initialMaintenanceIds.value]
  maintenanceSelectedIds.value = []
  maintenanceSearchLeft.value = ''
  maintenanceSearchRight.value = ''
  showingMaintenance.value = true
}

function toggleMaintenanceSelect(id) {
  const idx = maintenanceSelectedIds.value.indexOf(id)
  if (idx >= 0) maintenanceSelectedIds.value.splice(idx, 1)
  else maintenanceSelectedIds.value.push(id)
}

function toggleMaintenanceSelectAll(list) {
  const ids = list.map((s) => s.id)
  if (ids.every((id) => maintenanceSelectedIds.value.includes(id))) {
    maintenanceSelectedIds.value = maintenanceSelectedIds.value.filter((id) => !ids.includes(id))
  } else {
    ids.forEach((id) => {
      if (!maintenanceSelectedIds.value.includes(id)) maintenanceSelectedIds.value.push(id)
    })
  }
}

function moveMaintenanceRight() {
  const selected = maintenanceSelectedIds.value.filter((id) =>
    leftMaintenanceServers.value.some((s) => s.id === id)
  )
  selected.forEach((id) => {
    if (!maintenanceSelected.value.includes(id)) maintenanceSelected.value.push(id)
  })
  maintenanceSelectedIds.value = maintenanceSelectedIds.value.filter((id) => !selected.includes(id))
}

function moveMaintenanceLeft() {
  const selected = maintenanceSelectedIds.value.filter((id) =>
    rightMaintenanceServers.value.some((s) => s.id === id)
  )
  maintenanceSelected.value = maintenanceSelected.value.filter((id) => !selected.includes(id))
  maintenanceSelectedIds.value = maintenanceSelectedIds.value.filter((id) => !selected.includes(id))
}

function confirmMaintenance() {
  showingMaintenance.value = false
  showingMaintenanceConfirm.value = true
}

function backToMaintenance() {
  showingMaintenanceConfirm.value = false
  showingMaintenance.value = true
}

async function executeMaintenance() {
  const { enter, restore } = maintenanceTargets.value
  if (enter.length === 0 && restore.length === 0) return
  maintaining.value = true
  try {
    await setServersMaintenance(enter.map((s) => s.id), restore.map((s) => s.id))
    alert(t('servers.alert.maintenanceSuccess'))
    showingMaintenanceConfirm.value = false
    showingMaintenance.value = false
    await loadServers()
  } catch (e) {
    alert(e.message || t('servers.alert.operationFailed'))
  }
  maintaining.value = false
}
</script>

<template>
  <section class="servers">
    <div class="servers__head">
      <h1 class="page__title">{{ t('servers.title') }}</h1>
      <div class="servers__actions">
        <button class="btn btn--primary" type="button" @click="openCreate">
          {{ t('servers.add') }}
        </button>
        <button class="btn btn--secondary" type="button" @click="openPublish">
          {{ t('servers.publish.button') }}
        </button>
        <button class="btn btn--secondary" type="button" @click="openMaintenance">
          {{ t('servers.maintenance.button') }}
        </button>
      </div>
    </div>

    <div class="server-list">
       <div class="server-list__header">
        <span class="server-list__header-cell server-list__name-col">{{ t('common.name') }}</span>
        <span class="server-list__header-cell server-list__uuid-col">UUID</span>
        <span class="server-list__header-cell server-list__dept-col">{{ t('servers.department') }}</span>
        <span class="server-list__header-cell server-list__status-col">{{ t('common.status') }}</span>
        <span class="server-list__header-cell server-list__action-col">{{ t('common.actions') }}</span>
      </div>
      <button
        v-for="server in servers"
        :key="server.id"
        class="server-list__item"
        :class="{ 'server-list__item--active': selectedId === server.id }"
        type="button"
        @click="select(server)"
      >
        <div class="server-list__info server-list__name-col">
          <span class="server-list__name">{{ server.name }}</span>
        </div>
        <span class="server-list__uuid server-list__uuid-col">{{ server.id }}</span>
        <span class="server-list__dept server-list__dept-col">{{ server.department }}</span>
        <span class="server-list__status-col">
          <span class="server-list__status" :class="server.status === 'published' ? 'server-list__status--on' : server.status === 'maintenance' ? 'server-list__status--maint' : 'server-list__status--off'">
            {{ server.status === 'published' ? t('servers.status.published') : server.status === 'maintenance' ? t('servers.status.maintenance') : t('servers.status.unpublished') }}
          </span>
        </span>
        <span class="server-list__action-col">
          <button
            class="server-list__del"
            type="button"
            @click.stop="openDeleteConfirm(server)"
          >
            {{ t('common.delete') }}
          </button>
        </span>
      </button>
    </div>

    <div v-if="selected" class="drawer">
      <div class="drawer__overlay" @click="selected = null"></div>
      <aside class="drawer__panel">
        <header class="drawer__head">
          <h2 class="drawer__title">{{ selected.name }}</h2>
          <button class="drawer__close" type="button" :aria-label="t('common.close')" @click="selected = null">
            ×
          </button>
        </header>
        <div class="drawer__body">
          <dl class="server-detail__grid">
            <div class="server-detail__row">
              <dt>UUID</dt>
              <dd class="server-detail__uuid">{{ selected.id }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.uriPath') }}</dt>
              <dd>{{ selected.address }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.serviceAddress') }}</dt>
              <dd>
                <template v-if="formatServiceAddresses(selected.service_address)">
                  <div
                    v-for="(addr, i) in formatServiceAddresses(selected.service_address).split('\n')"
                    :key="i"
                  >{{ addr }}</div>
                </template>
                <span v-else>{{ t('servers.notFilled') }}</span>
              </dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.responsibleDepartment') }}</dt>
              <dd>{{ selected.department }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.protocolType') }}</dt>
              <dd>{{ selected.protocol }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.protocolVersion') }}</dt>
              <dd>{{ selected.protocol_version || '2026-07-28' }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.tools') }}</dt>
              <dd>
                <ul class="server-detail__tools">
                  <li
                    v-for="tool in selectedTools"
                    :key="tool.name"
                  >
                    <span class="server-detail__tool-name">{{ tool.name }}</span>
                    <span v-if="tool.description" class="server-detail__tool-desc">{{ tool.description }}</span>
                  </li>
                </ul>
              </dd>
            </div>
            <div class="server-detail__row">
              <dt>{{ t('servers.detail.publishStatus') }}</dt>
              <dd>
                <span class="server-detail__status" :class="selected.status === 'published' ? 'server-detail__status--on' : selected.status === 'maintenance' ? 'server-detail__status--maint' : 'server-detail__status--off'">
                  {{ selected.status === 'published' ? t('servers.status.published') : selected.status === 'maintenance' ? t('servers.status.maintenance') : t('servers.status.unpublished') }}
                </span>
              </dd>
            </div>
            <div class="server-detail__row">
              <div class="server-detail__label-row">
                <dt>{{ t('common.description') }}</dt>
                <button
                  v-if="!editingDesc"
                  class="server-detail__edit-btn"
                  type="button"
                  @click="startEditDesc"
                >
                  {{ t('common.edit') }}
                </button>
              </div>
              <dd v-if="!editingDesc">{{ selected.description || t('servers.detail.noDescription') }}</dd>
              <dd v-else class="server-detail__edit">
                <textarea
                  v-model="editDescValue"
                  class="dialog__input dialog__textarea"
                  rows="3"
                  :placeholder="t('servers.detail.descriptionPlaceholder')"
                ></textarea>
                <div class="server-detail__edit-actions">
                  <button class="btn btn--ghost" type="button" @click="cancelEditDesc">{{ t('common.cancel') }}</button>
                  <button class="btn btn--primary" type="button" :disabled="savingDesc" @click="saveDescription">
                    {{ savingDesc ? t('common.saving') : t('common.save') }}
                  </button>
                </div>
              </dd>
            </div>
          </dl>
        </div>
      </aside>
    </div>

    <Teleport to="body">
      <div v-if="showingCreate" class="dialog-overlay" @click.self="showingCreate = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('servers.form.title') }}</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">{{ t('common.name') }}</label>
              <input
                v-model="createForm.name"
                class="dialog__input"
                type="text"
                :placeholder="t('servers.form.namePlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">
                {{ t('servers.form.uriPathLabel') }}
                <span class="dialog__help">
                  ?
                  <span class="dialog__tooltip">
                    <i18n-t keypath="servers.form.uriPathHelp">
                      <template #addr><code>192.168.1.100:8080/mcp</code></template>
                      <template #path><code>/mcp</code></template>
                    </i18n-t>
                  </span>
                </span>
              </label>
              <input
                v-model="createForm.address"
                class="dialog__input"
                type="text"
                :placeholder="t('servers.form.uriPathPlaceholder', { path: '/mcp/server' })"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">
                {{ t('servers.detail.serviceAddress') }}
                <span class="dialog__help">
                  ?
                  <span class="dialog__tooltip">
                    <i18n-t keypath="servers.form.serviceAddressHelp">
                      <template #addr><code>192.168.1.100:8081</code></template>
                    </i18n-t>
                  </span>
                </span>
              </label>
              <textarea
                v-model="createForm.service_address"
                class="dialog__input dialog__textarea"
                rows="3"
                placeholder="192.168.1.100:8081"
              ></textarea>
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('servers.detail.responsibleDepartment') }}</label>
              <input
                v-model="createForm.department"
                class="dialog__input"
                type="text"
                :placeholder="t('servers.form.departmentPlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('common.description') }}</label>
              <textarea
                v-model="createForm.description"
                class="dialog__input dialog__textarea"
                rows="2"
                :placeholder="t('servers.form.descriptionPlaceholder')"
              ></textarea>
            </div>
            <hr class="dialog__divider" />
            <p class="dialog__section-label">{{ t('servers.fetch.sectionTitle') }}</p>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('servers.form.selectAddress') }}</label>
              <select v-model="fetchAddressValue" class="dialog__input">
                <option value="" disabled>{{ t('servers.form.selectAddressPlaceholder') }}</option>
                <option
                  v-for="addr in serviceAddressList"
                  :key="addr"
                  :value="addr"
                >{{ addr }}</option>
                <option value="__custom__">{{ t('servers.form.customAddressOption') }}</option>
              </select>
            </div>
            <div v-if="isCustomFetchAddress" class="dialog__group">
              <label class="dialog__label">{{ t('servers.form.customAddressLabel') }}</label>
              <input
                v-model="customAddressInput"
                class="dialog__input"
                type="text"
                :placeholder="t('servers.form.customAddressPlaceholder', { example: '192.168.1.100:8081' })"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('servers.form.connectionMethod') }}</label>
              <div class="dialog__toggle">
                <button
                  type="button"
                  :class="['dialog__toggle-btn', { 'dialog__toggle-btn--active': !useHttps }]"
                  @click="useHttps = false"
                >HTTP</button>
                <button
                  type="button"
                  :class="['dialog__toggle-btn', { 'dialog__toggle-btn--active': useHttps }]"
                  @click="useHttps = true"
                >HTTPS</button>
              </div>
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('servers.detail.protocolType') }}</label>
              <select v-model="createForm.protocol" class="dialog__input">
                <option v-for="p in PROTOCOLS" :key="p" :value="p">{{ p }}</option>
              </select>
            </div>
            <div class="dialog__group">
              <label class="dialog__label">
                {{ t('servers.detail.protocolVersion') }}
                <span class="dialog__help">
                  ?
                  <span class="dialog__tooltip">
                    {{ t('servers.form.protocolVersionHelp', { version: '2026-07-28' }) }}
                  </span>
                </span>
              </label>
              <select v-model="createForm.protocol_version" class="dialog__input">
                <option v-for="v in PROTOCOL_VERSIONS" :key="v" :value="v">{{ v }}</option>
              </select>
            </div>
            <div class="dialog__group">
              <div class="dialog__label-row">
                <label class="dialog__label">{{ t('servers.detail.tools') }}</label>
                <button
                  class="dialog__fetch"
                  type="button"
                  :disabled="fetchingTools"
                  @click="fetchTools"
                >
                  {{ fetchingTools ? t('servers.fetch.fetching') : t('servers.fetch.button') }}
                </button>
              </div>
              <textarea
                v-model="createForm.tools"
                class="dialog__input dialog__textarea"
                :placeholder="t('servers.form.toolsPlaceholder')"
              ></textarea>
              <p v-if="fetchToolsError" class="dialog__error">{{ fetchToolsError }}</p>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingCreate = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!createForm.name.trim()"
              @click="confirmCreate"
            >
              {{ t('servers.form.confirm') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showingPublish" class="dialog-overlay" @click.self="showingPublish = false">
        <div class="dialog dialog--wide dialog--publish">
          <h2 class="dialog__title">{{ t('servers.publish.title') }}</h2>
          <p class="dialog__desc">{{ t('servers.publish.desc') }}</p>
          <p v-if="gatewayConfigured" class="publish-gateway__hint">
            {{ t('servers.publish.gatewayConfigured', { kong: gatewayProvider === 'kong' ? t('servers.publish.kongSuffix') : '' }) }}
          </p>
          <p v-if="isKongGateway" class="publish-gateway__hint publish-gateway__hint--warn">
            {{ t('servers.publish.kongHint') }}
          </p>
          <div class="shuttle">
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">{{ t('servers.shuttle.available') }}</span>
                <span class="shuttle__count">{{ leftPublishServers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="publishSearchLeft"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('servers.shuttle.filterPlaceholder')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllLeftPublishChecked"
                    @change="togglePublishSelectAll(leftPublishServers)"
                  />
                  <span>{{ t('servers.shuttle.selectAll') }}</span>
                </label>
                <label
                  v-for="server in leftPublishServers"
                  :key="server.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="publishSelectedIds.includes(server.id)"
                    @change="togglePublishSelect(server.id)"
                  />
                  <span class="shuttle__name">{{ server.name }}</span>
                  <span class="shuttle__meta">{{ server.address }}</span>
                </label>
                <span v-if="leftPublishServers.length === 0" class="table__muted shuttle__empty">
                  {{ t('servers.shuttle.emptyAvailable') }}
                </span>
              </div>
            </div>
            <div class="shuttle__actions">
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!publishSelectedIds.some((id) => leftPublishServers.some((s) => s.id === id))"
                @click="movePublishRight"
              >
                &gt;
              </button>
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!publishSelectedIds.some((id) => rightPublishServers.some((s) => s.id === id))"
                @click="movePublishLeft"
              >
                &lt;
              </button>
            </div>
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">{{ t('servers.shuttle.selected') }}</span>
                <span class="shuttle__count">{{ rightPublishServers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="publishSearchRight"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('servers.shuttle.filterPlaceholder')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllRightPublishChecked"
                    @change="togglePublishSelectAll(rightPublishServers)"
                  />
                  <span>{{ t('servers.shuttle.selectAll') }}</span>
                </label>
                <label
                  v-for="server in rightPublishServers"
                  :key="server.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="publishSelectedIds.includes(server.id)"
                    @change="togglePublishSelect(server.id)"
                  />
                  <span class="shuttle__name">{{ server.name }}</span>
                  <span class="shuttle__meta">{{ server.address }}</span>
                </label>
                <span v-if="rightPublishServers.length === 0" class="table__muted shuttle__empty">
                  {{ t('servers.shuttle.emptySelected') }}
                </span>
              </div>
            </div>
          </div>
          <div class="publish-auth">
            <label class="publish-auth__label">
              <input type="checkbox" v-model="enableAuth" :disabled="isKongGateway" />
              <span class="publish-auth__text">{{ t('servers.publish.optionAuth') }}</span>
            </label>
            <p class="publish-auth__hint">
              {{ t('servers.publish.authHint') }}
            </p>
          </div>
          <div class="publish-auth">
            <label class="publish-auth__label">
              <input type="checkbox" v-model="enableAuditLog" :disabled="isKongGateway" />
              <span class="publish-auth__text">{{ t('servers.publish.optionAudit') }}</span>
            </label>
            <p class="publish-auth__hint">
              {{ t('servers.publish.auditHint') }}
            </p>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingPublish = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="publishSelected.length === 0"
              @click="confirmPublish"
            >
              {{ t('servers.publish.confirm') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showingPublishConfirm" class="dialog-overlay" @click.self="showingPublishConfirm = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('servers.publish.confirmTitle') }}</h2>
          <p class="dialog__desc">
            {{ t('servers.publish.confirmDesc') }}
          </p>
          <div class="publish-preview">
            <div
              v-for="server in publishTargets"
              :key="server.id"
              class="publish-preview__card"
            >
              <div class="publish-preview__row">
                <span class="publish-preview__label">{{ t('servers.preview.serviceName') }}</span>
                <span class="publish-preview__value">{{ server.name }}</span>
              </div>
              <div class="publish-preview__row">
                <span class="publish-preview__label">{{ t('servers.preview.gatewayPath') }}</span>
                <span class="publish-preview__value">{{ server.address }}</span>
              </div>
              <div class="publish-preview__row">
                <span class="publish-preview__label">{{ t('servers.preview.backendAddress') }}</span>
                <span class="publish-preview__value">
                  <template v-if="formatServiceAddresses(server.service_address)">{{ formatServiceAddresses(server.service_address).split('\n').join(', ') }}</template>
                  <span v-else>{{ t('servers.notFilled') }}</span>
                </span>
              </div>
              <div class="publish-preview__row">
                <span class="publish-preview__label">{{ t('servers.preview.authStatus') }}</span>
                <span class="publish-preview__value" :class="{ 'publish-preview__auth-on': enableAuth }">
                  {{ enableAuth ? t('servers.preview.authEnabled') : t('servers.preview.disabled') }}
                </span>
              </div>
              <div class="publish-preview__row">
                <span class="publish-preview__label">{{ t('servers.preview.auditLog') }}</span>
                <span class="publish-preview__value" :class="{ 'publish-preview__auth-on': enableAuditLog }">
                  {{ enableAuditLog ? t('servers.preview.enabled') : t('servers.preview.disabled') }}
                </span>
              </div>
          </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="backToPublish">
              {{ t('servers.backToEdit') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="publishing"
              @click="executePublish"
            >
              {{ publishing ? t('servers.publish.progress') : t('servers.confirmPush') }}
            </button>
          </div>
        </div>
      </div>
      <div v-if="showingMaintenance" class="dialog-overlay" @click.self="showingMaintenance = false">
        <div class="dialog dialog--wide dialog--publish">
          <h2 class="dialog__title">{{ t('servers.maintenance.title') }}</h2>
          <p class="dialog__desc">{{ t('servers.maintenance.desc') }}</p>
          <div class="shuttle">
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">{{ t('servers.shuttle.available') }}</span>
                <span class="shuttle__count">{{ leftMaintenanceServers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="maintenanceSearchLeft"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('servers.shuttle.filterPlaceholder')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllLeftMaintenanceChecked"
                    @change="toggleMaintenanceSelectAll(leftMaintenanceServers)"
                  />
                  <span>{{ t('servers.shuttle.selectAll') }}</span>
                </label>
                <label
                  v-for="server in leftMaintenanceServers"
                  :key="server.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="maintenanceSelectedIds.includes(server.id)"
                    @change="toggleMaintenanceSelect(server.id)"
                  />
                  <span class="shuttle__name">{{ server.name }}</span>
                  <span class="shuttle__meta">{{ server.address }}</span>
                </label>
                <span v-if="leftMaintenanceServers.length === 0" class="table__muted shuttle__empty">
                  {{ t('servers.shuttle.emptyAvailable') }}
                </span>
              </div>
            </div>
            <div class="shuttle__actions">
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!maintenanceSelectedIds.some((id) => leftMaintenanceServers.some((s) => s.id === id))"
                @click="moveMaintenanceRight"
              >
                &gt;
              </button>
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!maintenanceSelectedIds.some((id) => rightMaintenanceServers.some((s) => s.id === id))"
                @click="moveMaintenanceLeft"
              >
                &lt;
              </button>
            </div>
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">{{ t('servers.maintenance.underMaintenance') }}</span>
                <span class="shuttle__count">{{ rightMaintenanceServers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="maintenanceSearchRight"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('servers.shuttle.filterPlaceholder')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllRightMaintenanceChecked"
                    @change="toggleMaintenanceSelectAll(rightMaintenanceServers)"
                  />
                  <span>{{ t('servers.shuttle.selectAll') }}</span>
                </label>
                <label
                  v-for="server in rightMaintenanceServers"
                  :key="server.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="maintenanceSelectedIds.includes(server.id)"
                    @change="toggleMaintenanceSelect(server.id)"
                  />
                  <span class="shuttle__name">{{ server.name }}</span>
                  <span class="shuttle__meta">{{ server.address }}</span>
                </label>
                <span v-if="rightMaintenanceServers.length === 0" class="table__muted shuttle__empty">
                  {{ t('servers.maintenance.emptyUnderMaintenance') }}
                </span>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingMaintenance = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="maintenanceTargets.enter.length === 0 && maintenanceTargets.restore.length === 0"
              @click="confirmMaintenance"
            >
              {{ maintenanceConfirmLabel }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showingMaintenanceConfirm" class="dialog-overlay" @click.self="showingMaintenanceConfirm = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('servers.maintenance.confirmBoth') }}</h2>
          <p class="dialog__desc">
            {{ t('servers.maintenance.confirmDesc') }}
          </p>
          <div class="publish-preview">
            <template v-if="maintenanceTargets.enter.length > 0">
              <p class="maintenance-preview__section">{{ t('servers.maintenance.sectionEnter') }}</p>
              <div
                v-for="server in maintenanceTargets.enter"
                :key="server.id"
                class="publish-preview__card"
              >
                <div class="publish-preview__row">
                  <span class="publish-preview__label">{{ t('servers.preview.serviceName') }}</span>
                  <span class="publish-preview__value">{{ server.name }}</span>
                </div>
                <div class="publish-preview__row">
                  <span class="publish-preview__label">{{ t('servers.preview.gatewayPath') }}</span>
                  <span class="publish-preview__value">{{ server.address }}</span>
                </div>
              </div>
            </template>
            <template v-if="maintenanceTargets.restore.length > 0">
              <p class="maintenance-preview__section">{{ t('servers.maintenance.sectionRestore') }}</p>
              <div
                v-for="server in maintenanceTargets.restore"
                :key="server.id"
                class="publish-preview__card"
              >
                <div class="publish-preview__row">
                  <span class="publish-preview__label">{{ t('servers.preview.serviceName') }}</span>
                  <span class="publish-preview__value">{{ server.name }}</span>
                </div>
                <div class="publish-preview__row">
                  <span class="publish-preview__label">{{ t('servers.preview.gatewayPath') }}</span>
                  <span class="publish-preview__value">{{ server.address }}</span>
                </div>
              </div>
            </template>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="backToMaintenance">
              {{ t('servers.backToEdit') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="maintaining"
              @click="executeMaintenance"
            >
              {{ maintaining ? t('servers.maintenance.submitting') : t('servers.confirmPush') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="showingDeleteConfirm" class="dialog-overlay" @click.self="showingDeleteConfirm = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('servers.alert.deleteConfirm') }}</h2>
          <p class="dialog__desc">
            {{ t('servers.alert.deleteBody', { name: deleteTarget?.name }) }}
          </p>
          <p class="dialog__warn">
            {{ t('servers.alert.deleteWarn') }}
          </p>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingDeleteConfirm = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--danger"
              type="button"
              :disabled="deleting"
              @click="removeServer"
            >
              {{ deleting ? t('servers.alert.deleting') : t('servers.alert.deleteConfirm') }}
            </button>
          </div>
        </div>
      </div>

    </Teleport>
  </section>
</template>

<style scoped>
.servers {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.servers__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.servers__actions {
  display: flex;
  gap: 10px;
}

.btn {
  padding: 8px 18px;
  font-size: 14px;
  font-weight: 600;
  border-radius: 999px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 0.2s, transform 0.15s, border-color 0.2s;
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
  transform: translateY(-1px);
}

.btn--secondary {
  color: var(--xauat-blue);
  background: transparent;
  border-color: var(--xauat-blue);
}

.btn--secondary:hover:not(:disabled) {
  background: rgba(10, 61, 122, 0.08);
}

.btn--ghost {
  color: var(--xauat-blue);
  background: transparent;
  border-color: var(--border);
}

.btn--ghost:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.btn--danger {
  color: #fff;
  background: #c0392b;
}

.btn--danger:hover:not(:disabled) {
  background: #e74c3c;
  transform: translateY(-1px);
}

.server-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 12px;
}

.server-list__header {
  display: flex;
  align-items: center;
  padding: 8px 16px;
  border-bottom: 2px solid var(--border);
  margin-bottom: 4px;
}

.server-list__header-cell {
  font-size: 13px;
  font-weight: 700;
  color: var(--text-muted);
}

.server-list__name-col {
  flex: 0 0 160px;
  min-width: 0;
}

.server-list__uuid-col {
  flex: 0 0 320px;
  margin-left: 24px;
}

.server-list__dept-col {
  flex: 0 0 120px;
  margin-left: 24px;
}

.server-list__status-col {
  flex: 0 0 80px;
  margin-left: 24px;
}

.server-list__status {
  display: inline-block;
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  border-radius: 999px;
}

.server-list__status--on {
  color: #1a7a1a;
  background: #d4edda;
}

.server-list__status--off {
  color: #6c757d;
  background: #e9ecef;
}

.server-list__status--maint {
  color: #92400e;
  background: #fef3c7;
}

.server-list__action-col {
  flex: 0 0 64px;
  margin-left: 12px;
}

.server-list__item {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  text-align: left;
  transition: background 0.2s, border-color 0.2s;
}

.server-list__info {
  display: flex;
  align-items: center;
  min-width: 0;
}

.server-list__item:hover {
  background: rgba(10, 61, 122, 0.06);
}

.server-list__item--active {
  border-color: var(--xauat-blue-light);
  background: rgba(30, 95, 176, 0.1);
}

.server-list__name {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
}

.server-list__uuid {
  font-size: 12px;
  font-family: monospace;
  color: var(--text-muted);
  opacity: 0.75;
}

.server-list__dept {
  font-size: 13px;
  color: var(--text-muted);
}

.server-list__del {
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 600;
  color: #fff;
  background: #c0392b;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.server-list__del:hover {
  background: #e74c3c;
}

.drawer {
  position: fixed;
  inset: 0;
  z-index: 200;
}

.drawer__overlay {
  position: absolute;
  inset: 0;
  background: rgba(10, 61, 122, 0.28);
}

.drawer__panel {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 420px;
  max-width: 90vw;
  display: flex;
  flex-direction: column;
  background: var(--surface);
  box-shadow: -8px 0 30px rgba(10, 61, 122, 0.12);
}

.drawer__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid var(--border);
}

.drawer__title {
  font-size: 17px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.drawer__close {
  width: 32px;
  height: 32px;
  font-size: 22px;
  line-height: 1;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.drawer__close:hover {
  background: rgba(10, 61, 122, 0.06);
}

.drawer__body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
}

.server-detail__grid {
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.server-detail__row {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.server-detail__row dt {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.server-detail__row dd {
  font-size: 14px;
  color: var(--text);
  word-break: break-all;
}

.server-detail__uuid {
  font-family: monospace;
  font-size: 13px;
  color: var(--text-muted);
  opacity: 0.85;
}

.server-detail__tools {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  list-style: none;
}

.server-detail__tools li {
  display: flex;
  flex-direction: column;
  padding: 8px 14px;
  background: rgba(30, 95, 176, 0.1);
  border-radius: 10px;
  font-size: 13px;
  color: var(--xauat-blue);
  cursor: default;
  gap: 4px;
}

.server-detail__tool-name {
  font-weight: 600;
}

.server-detail__tool-desc {
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.4;
}

.server-detail__label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.server-detail__edit-btn {
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.server-detail__edit-btn:hover {
  background: rgba(10, 61, 122, 0.06);
}

.server-detail__status {
  display: inline-block;
  padding: 2px 12px;
  font-size: 13px;
  font-weight: 600;
  border-radius: 999px;
}

.server-detail__status--on {
  color: #1a7a1a;
  background: #d4edda;
}

.server-detail__status--off {
  color: #6c757d;
  background: #e9ecef;
}

.server-detail__status--maint {
  color: #92400e;
  background: #fef3c7;
}

.maintenance-preview__section {
  font-size: 13px;
  font-weight: 700;
  color: var(--xauat-blue);
  margin: 4px 0;
}

.server-detail__edit {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.server-detail__edit-actions {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
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

.dialog__form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

.dialog__group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.dialog__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.dialog__help {
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

.dialog__tooltip {
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

.dialog__tooltip::after {
  content: '';
  position: absolute;
  top: 100%;
  left: 50%;
  transform: translateX(-50%);
  border: 6px solid transparent;
  border-top-color: var(--border);
}

.dialog__tooltip code {
  padding: 1px 5px;
  font-size: 11px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.1);
  border-radius: 4px;
}

.dialog__help:hover .dialog__tooltip {
  opacity: 1;
  visibility: visible;
}

.dialog__label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.dialog__fetch {
  padding: 4px 12px;
  font-size: 12px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.dialog__fetch:hover:not(:disabled) {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.dialog__fetch:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dialog__toggle {
  display: flex;
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid var(--border);
}

.dialog__toggle-btn {
  flex: 1;
  padding: 8px 0;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: color 0.2s, background 0.2s;
}

.dialog__toggle-btn--active {
  color: #fff;
  background: var(--xauat-blue);
}

.dialog__input {
  width: 100%;
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s;
}

.dialog__input:focus {
  border-color: var(--xauat-blue);
}

.dialog__textarea {
  min-height: 72px;
  resize: vertical;
  font-family: inherit;
}

.dialog__error {
  font-size: 12px;
  color: #c0392b;
}

.dialog__divider {
  border: none;
  border-top: 1px solid var(--border);
  margin: 0;
}

.dialog__section-label {
  font-size: 13px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
  flex-shrink: 0;
}

.dialog__desc {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.dialog__warn {
  padding: 12px 16px;
  font-size: 13px;
  line-height: 1.6;
  color: #856404;
  background: #fff3cd;
  border: 1px solid #ffc107;
  border-radius: 10px;
  margin-bottom: 8px;
}

.dialog--publish,
.dialog--wide {
  max-width: 640px;
}

.dialog--publish {
  min-height: 620px;
  height: auto;
}

.dialog--publish .dialog__title,
.dialog--publish .dialog__desc,
.dialog--publish .publish-gateway__hint,
.dialog--publish .publish-auth {
  flex-shrink: 0;
}

.shuttle {
  display: flex;
  gap: 12px;
  flex: 1 1 auto;
  min-height: 0;
}

.shuttle__panel {
  flex: 1;
  min-width: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.shuttle__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg);
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.shuttle__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.shuttle__search {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.shuttle__search-input {
  width: 100%;
  padding: 6px 10px;
  font-size: 13px;
  border: 1px solid var(--border);
  border-radius: 6px;
  outline: none;
  background: var(--surface);
  color: var(--text);
  box-sizing: border-box;
}

.shuttle__search-input::placeholder {
  color: var(--text-muted);
}

.shuttle__search-input:focus {
  border-color: var(--xauat-blue);
}

.shuttle__count {
  font-size: 12px;
  color: var(--text-muted);
  background: var(--surface);
  padding: 2px 8px;
  border-radius: 999px;
}

.shuttle__list {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.shuttle__item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  font-size: 14px;
  color: var(--text);
  cursor: pointer;
  border-radius: 6px;
  transition: background 0.15s;
}

.shuttle__item:hover {
  background: var(--bg);
}

.shuttle__item input {
  accent-color: var(--xauat-blue);
  flex-shrink: 0;
}

.shuttle__item--all {
  padding-bottom: 8px;
  margin-bottom: 4px;
  border-bottom: 1px solid var(--border);
  font-weight: 600;
}

.shuttle__name {
  font-size: 14px;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.shuttle__meta {
  font-size: 12px;
  color: var(--text-muted);
  margin-left: auto;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 120px;
  flex-shrink: 0;
}

.shuttle__empty {
  padding: 16px 0;
  text-align: center;
}

.shuttle__actions {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 10px;
  flex-shrink: 0;
}

.shuttle__btn {
  width: 32px;
  height: 32px;
  font-size: 16px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background 0.2s, border-color 0.2s;
}

.shuttle__btn:hover:not(:disabled) {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.shuttle__btn:disabled {
  opacity: 0.35;
  cursor: not-allowed;
}

.publish-preview {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 320px;
  overflow-y: auto;
}

.publish-preview__card {
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--bg);
}

.publish-preview__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 6px 0;
}

.publish-preview__row + .publish-preview__row {
  border-top: 1px solid rgba(10, 61, 122, 0.06);
}

.publish-preview__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  flex-shrink: 0;
}

.publish-preview__value {
  font-size: 13px;
  color: var(--text);
  word-break: break-all;
  text-align: right;
}

.publish-preview__auth-on {
  color: var(--xauat-blue);
  font-weight: 600;
}

.publish-auth {
  margin-top: 16px;
  padding: 12px 14px;
  border: 1px solid var(--border);
  border-radius: 10px;
  background: var(--bg);
}

.publish-auth__label {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.publish-auth__label input {
  width: 16px;
  height: 16px;
  accent-color: var(--xauat-blue);
  cursor: pointer;
  flex-shrink: 0;
}

.publish-auth__text {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.publish-auth__hint {
  margin-top: 6px;
  margin-left: 24px;
  font-size: 12px;
  color: var(--text-muted);
  line-height: 1.5;
}

.publish-gateway__hint {
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 600;
  color: #1a7a1a;
  background: #d4edda;
  border-radius: 8px;
  display: inline-block;
}

.publish-gateway__hint--warn {
  color: #92400e;
  background: #fef3c7;
}
</style>

<script setup>
import { ref, onMounted } from 'vue'
import { setActive } from '../stores/nav.js'
import { fetchAccessKeys, createAccessKey, updateAccessKey, deleteAccessKey, fetchServers } from '../api.js'
import AccessKeyDrawer from './AccessKeyDrawer.vue'

const servers = ref([])

function parseTools(raw) {
  try { return JSON.parse(raw) || [] } catch { return [] }
}

function defaultServersJSON() {
  const obj = {}
  servers.value.forEach((s) => {
    const tools = parseTools(s.tools)
    if (tools.length) obj[s.id] = tools
  })
  return JSON.stringify(obj)
}

onMounted(async () => {
  try {
    keys.value = await fetchAccessKeys()
    const data = await fetchServers()
    servers.value = data
  } catch {
    // ignore load error
  }
  loading.value = false
})

const EXPIRATION_OPTIONS = [
  { label: '7 天', value: 7 },
  { label: '30 天', value: 30 },
  { label: '90 天', value: 90 },
  { label: '180 天', value: 180 },
  { label: '365 天', value: 365 },
  { label: '永不过期', value: -1 }
]

const keys = ref([])
const loading = ref(true)
const copiedId = ref(null)
const editingKey = ref(null)

const showingCreate = ref(false)
const createForm = ref({ name: '', expireDays: 30 })
const createdKey = ref(null)
const dialogCopyId = ref(null)

function openCreate() {
  createForm.value = { name: '', expireDays: 30 }
  showingCreate.value = true
}

async function confirmCreate() {
  if (!createForm.value.name.trim()) return
  const body = {
    name: createForm.value.name.trim(),
    servers: defaultServersJSON()
  }
  if (createForm.value.expireDays !== -1) {
    const d = new Date()
    d.setDate(d.getDate() + createForm.value.expireDays)
    body.expired_at = d.toISOString()
  }
  try {
    const key = await createAccessKey(body)
    createdKey.value = key
    showingCreate.value = false
  } catch {
    // ignore
  }
}

function closeDialog() {
  if (createdKey.value) {
    keys.value.unshift(createdKey.value)
  }
  createdKey.value = null
  dialogCopyId.value = null
}

async function copyDialogKey() {
  if (!createdKey.value) return
  try {
    await navigator.clipboard.writeText(createdKey.value.key)
    dialogCopyId.value = createdKey.value.id
    setTimeout(() => {
      if (dialogCopyId.value === createdKey.value.id) dialogCopyId.value = null
    }, 1500)
  } catch {
    // ignore
  }
}

function formatExpire(key) {
  if (!key.expired_at) return '永不过期'
  const d = new Date(key.expired_at)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: '2-digit', day: '2-digit' })
}

function maskKey(value) {
  if (!value) return ''
  const prefix = value.slice(0, 3)
  const body = value.slice(3)
  if (body.length <= 8) return value
  return `${prefix}${body.slice(0, 4)}****${body.slice(-4)}`
}

async function toggle(key) {
  try {
    await updateAccessKey(key.id, { enabled: !key.enabled })
    key.enabled = !key.enabled
  } catch {
    // ignore
  }
}

async function copyKey(key) {
  try {
    await navigator.clipboard.writeText(key.key)
    copiedId.value = key.id
    setTimeout(() => {
      if (copiedId.value === key.id) copiedId.value = null
    }, 1500)
  } catch {
    // ignore
  }
}

function openEdit(key) {
  editingKey.value = { ...key }
}

async function removeKey(key) {
  try {
    await deleteAccessKey(key.id)
    keys.value = keys.value.filter((k) => k.id !== key.id)
  } catch {
    // ignore
  }
}

async function onSave(payload) {
  try {
    await updateAccessKey(editingKey.value.id, payload)
    Object.assign(editingKey.value, payload)
    const idx = keys.value.findIndex((k) => k.id === editingKey.value.id)
    if (idx !== -1) {
      editingKey.value.updated_at = new Date().toISOString()
      keys.value[idx] = { ...editingKey.value }
    }
  } catch {
    // ignore
  }
  editingKey.value = null
}

function goHistory() {
  setActive('history')
}
</script>

<template>
  <section class="page">
    <div class="page__head">
      <h1 class="page__title">AccessKey 管理</h1>
      <div class="page__actions">
        <button class="btn btn--ghost" type="button" @click="goHistory">
          使用历史
        </button>
        <button class="btn btn--primary" type="button" @click="openCreate">
          创建 AccessKey
        </button>
      </div>
    </div>

    <div v-if="showingCreate" class="create-form">
      <div class="create-form__group">
        <label class="create-form__label">名称</label>
        <input
          v-model="createForm.name"
          class="create-form__name"
          type="text"
          placeholder="请输入 AccessKey 名称"
          @keyup.enter="confirmCreate"
        />
      </div>
      <div class="create-form__group">
        <label class="create-form__label">过期时间</label>
        <select v-model="createForm.expireDays" class="create-form__select">
          <option
            v-for="opt in EXPIRATION_OPTIONS"
            :key="opt.value"
            :value="opt.value"
          >
            {{ opt.label }}
          </option>
        </select>
      </div>
      <div class="create-form__actions">
        <button class="btn btn--ghost" type="button" @click="showingCreate = false">
          取消
        </button>
        <button
          class="btn btn--primary"
          type="button"
          :disabled="!createForm.name.trim()"
          @click="confirmCreate"
        >
          确认创建
        </button>
      </div>
    </div>

    <table v-if="!loading" class="keys">
      <thead>
        <tr>
          <th>名称</th>
          <th>AccessKey</th>
          <th class="keys__col-expire">过期时间</th>
          <th class="keys__col-status">状态</th>
          <th class="keys__col-action">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="keys.length === 0">
          <td colspan="5" class="keys__empty">暂无 AccessKey</td>
        </tr>
        <tr v-for="key in keys" :key="key.id">
          <td>{{ key.name }}</td>
          <td>
            <div class="keys__key">
              <code
                class="keys__value"
                :class="{ 'keys__value--off': !key.enabled }"
              >{{ maskKey(key.key) }}</code>
              <button class="keys__copy" type="button" @click="copyKey(key)">
                {{ copiedId === key.id ? '已复制' : '复制' }}
              </button>
            </div>
          </td>
          <td class="keys__col-expire">{{ formatExpire(key) }}</td>
          <td class="keys__col-status">
            <label class="switch">
              <input
                type="checkbox"
                :checked="key.enabled"
                @change="toggle(key)"
              />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ key.enabled ? '启用' : '禁用' }}</span>
            </label>
          </td>
          <td class="keys__col-action">
            <button class="keys__edit" type="button" @click="openEdit(key)">
              编辑
            </button>
            <button class="keys__del" type="button" @click="removeKey(key)">
              删除
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <AccessKeyDrawer
      :item="editingKey"
      :servers="servers"
      @save="onSave"
      @close="editingKey = null"
    />

    <Teleport to="body">
      <div v-if="createdKey" class="dialog-overlay" @click.self="closeDialog">
        <div class="dialog">
          <p class="dialog__warn">请妥善保存以下 AccessKey，关闭后将无法再次查看完整密钥。</p>
          <div class="dialog__key">
            <code class="dialog__value">{{ createdKey.key }}</code>
            <button class="keys__copy" type="button" @click="copyDialogKey">
              {{ dialogCopyId === createdKey.id ? '已复制' : '复制' }}
            </button>
          </div>
          <button class="btn btn--primary" type="button" @click="closeDialog">
            我知道了
          </button>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<style scoped>
.page__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.page__actions {
  display: flex;
  gap: 12px;
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

.btn--ghost {
  color: var(--xauat-blue);
  background: transparent;
  border-color: var(--border);
}

.btn--ghost:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.create-form {
  margin-top: 20px;
  padding: 24px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  display: flex;
  align-items: flex-end;
  gap: 16px;
  flex-wrap: wrap;
}

.create-form__group {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.create-form__label {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.create-form__name {
  width: 240px;
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s;
}

.create-form__name:focus {
  border-color: var(--xauat-blue);
}

.create-form__select {
  width: 160px;
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  cursor: pointer;
  appearance: auto;
}

.create-form__select:focus {
  border-color: var(--xauat-blue);
}

.create-form__actions {
  display: flex;
  gap: 10px;
  margin-left: auto;
}

.keys {
  width: 100%;
  margin-top: 20px;
  border-collapse: collapse;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.keys th,
.keys td {
  padding: 14px 20px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.keys thead th {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg);
}

.keys tbody tr:last-child td {
  border-bottom: none;
}

.keys__empty {
  text-align: center;
  color: var(--text-muted);
  padding: 32px 20px !important;
}

.keys__col-status {
  width: 140px;
}

.keys__col-expire {
  width: 130px;
}

.keys__col-action {
  width: 150px;
}

.keys__key {
  display: flex;
  align-items: center;
  gap: 12px;
}

.keys__value {
  font-family: "SF Mono", "Menlo", "Consolas", monospace;
  font-size: 14px;
  color: var(--text);
  word-break: break-all;
}

.keys__value--off {
  color: var(--text-muted);
  text-decoration: line-through;
}

.keys__copy,
.keys__edit,
.keys__del {
  flex: none;
  padding: 4px 12px;
  font-size: 13px;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.keys__copy:hover,
.keys__edit:hover,
.keys__del:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
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
  text-align: center;
  box-shadow: var(--shadow);
}

.dialog__warn {
  font-size: 15px;
  color: var(--xauat-accent);
  font-weight: 600;
  margin-bottom: 20px;
  line-height: 1.6;
}

.dialog__key {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px;
  margin-bottom: 24px;
}

.dialog__value {
  flex: 1;
  font-family: "SF Mono", "Menlo", "Consolas", monospace;
  font-size: 14px;
  color: var(--text);
  word-break: break-all;
  text-align: left;
}
</style>

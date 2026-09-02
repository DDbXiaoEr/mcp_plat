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
import { ref, onMounted, watch, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import { setActive } from '../stores/nav.js'
import { fetchAccessKeys, createAccessKey, updateAccessKey, deleteAccessKey, fetchServers, fetchSetting } from '../api.js'
import AccessKeyDrawer from './AccessKeyDrawer.vue'

const { t, locale } = useI18n()

const servers = ref([])

function parseTools(raw) {
  try {
    const parsed = JSON.parse(raw)
    if (Array.isArray(parsed)) {
      return parsed.map((t) => {
        if (typeof t === 'object' && t.name) return t.name
        return t
      })
    }
    return []
  } catch {
    return []
  }
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
    const ops = await fetchSetting('user_ops')
    if (ops && ops.maxAccessKeys > 0) maxKeys.value = ops.maxAccessKeys
  } catch {
    // ignore load error
  }
  loading.value = false
})

const EXPIRATION_OPTIONS = [
  { value: 7 },
  { value: 30 },
  { value: 90 },
  { value: 180 },
  { value: 365 },
  { value: -1 }
]

function expiryLabel(value) {
  if (value === -1) return t('accesskey.neverExpires')
  return t('accesskey.expiryDays', { days: value })
}

const keys = ref([])
const maxKeys = ref(0)
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
  } catch (e) {
    alert(e.message || t('accesskey.createFailed'))
  }
}

function closeDialog() {
  if (createdKey.value) {
    keys.value.unshift(createdKey.value)
  }
  createdKey.value = null
  dialogCopyId.value = null
}

function onDialogKeydown(e) {
  if (e.key === 'Enter') closeDialog()
}

watch(createdKey, (val) => {
  if (val) {
    window.addEventListener('keydown', onDialogKeydown)
  } else {
    window.removeEventListener('keydown', onDialogKeydown)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onDialogKeydown)
})

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
  if (!key.expired_at) return t('accesskey.neverExpires')
  const d = new Date(key.expired_at)
  return d.toLocaleDateString(locale.value, { year: 'numeric', month: '2-digit', day: '2-digit' })
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
  } catch (e) {
    alert(e.message || t('accesskey.operateFailed'))
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
  } catch (e) {
    alert(e.message || t('accesskey.saveFailed'))
    return
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
      <h1 class="page__title">{{ t('nav.accesskey') }}</h1>
      <div class="page__actions">
        <span v-if="maxKeys > 0" class="keys__limit">
          {{ keys.length }} / {{ maxKeys }}
        </span>
        <button class="btn btn--ghost" type="button" @click="goHistory">
          {{ t('nav.history') }}
        </button>
        <button
          class="btn btn--primary"
          type="button"
          :disabled="maxKeys > 0 && keys.length >= maxKeys"
          :title="maxKeys > 0 && keys.length >= maxKeys ? t('accesskey.limitReached') : ''"
          @click="openCreate"
        >
          {{ t('accesskey.create') }}
        </button>
      </div>
    </div>

    <div v-if="showingCreate" class="create-form">
      <div class="create-form__group">
        <label class="create-form__label">{{ t('common.name') }}</label>
        <input
          v-model="createForm.name"
          class="create-form__name"
          type="text"
          :placeholder="t('accesskey.namePlaceholder')"
          @keyup.enter="confirmCreate"
        />
      </div>
      <div class="create-form__group">
        <label class="create-form__label">{{ t('accesskey.expireTime') }}</label>
        <select v-model="createForm.expireDays" class="create-form__select">
          <option
            v-for="opt in EXPIRATION_OPTIONS"
            :key="opt.value"
            :value="opt.value"
          >
            {{ expiryLabel(opt.value) }}
          </option>
        </select>
      </div>
      <div class="create-form__actions">
        <button class="btn btn--ghost" type="button" @click="showingCreate = false">
          {{ t('common.cancel') }}
        </button>
        <button
          class="btn btn--primary"
          type="button"
          :disabled="!createForm.name.trim()"
          @click="confirmCreate"
        >
          {{ t('accesskey.confirmCreate') }}
        </button>
      </div>
    </div>

    <table v-if="!loading" class="keys">
      <thead>
        <tr>
          <th>{{ t('common.name') }}</th>
          <th>AccessKey</th>
          <th class="keys__col-expire">{{ t('accesskey.expireTime') }}</th>
          <th class="keys__col-status">{{ t('common.status') }}</th>
          <th class="keys__col-action">{{ t('common.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="keys.length === 0">
          <td colspan="5" class="keys__empty">{{ t('accesskey.empty') }}</td>
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
                {{ copiedId === key.id ? t('common.copied') : t('common.copy') }}
              </button>
            </div>
          </td>
          <td class="keys__col-expire">{{ formatExpire(key) }}</td>
          <td class="keys__col-status">
            <span v-if="key.is_expired" class="keys__expired-badge">{{ t('accesskey.expired') }}</span>
            <label v-else class="switch">
              <input
                type="checkbox"
                :checked="key.enabled"
                @change="toggle(key)"
              />
              <span class="switch__track"><span class="switch__thumb"></span></span>
              <span class="switch__label">{{ key.enabled ? t('common.enabled') : t('common.disabled') }}</span>
            </label>
          </td>
          <td class="keys__col-action">
            <button class="keys__edit" type="button" @click="openEdit(key)">
              {{ t('common.edit') }}
            </button>
            <button class="keys__del" type="button" @click="removeKey(key)">
              {{ t('common.delete') }}
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
          <p class="dialog__warn">{{ t('accesskey.createdWarn') }}</p>
          <div class="dialog__key">
            <code class="dialog__value">{{ createdKey.key }}</code>
            <button class="keys__copy" type="button" @click="copyDialogKey">
              {{ dialogCopyId === createdKey.id ? t('common.copied') : t('common.copy') }}
            </button>
          </div>
          <button class="btn btn--primary" type="button" @click="closeDialog">
            {{ t('accesskey.gotIt') }}
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
  align-items: center;
  gap: 12px;
}

.keys__limit {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  padding: 4px 10px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 999px;
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

.keys__expired-badge {
  display: inline-block;
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  color: #b71c1c;
  background: #fbe9e7;
  border: 1px solid #ef9a9a;
  border-radius: 999px;
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

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

// Author: deepseek-v4-flash / opencode
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  checkGatewayServer,
  createGatewayServer,
  deleteGatewayServer,
  fetchGatewayServers,
  setDefaultGatewayServer,
  testGatewayServer,
  updateGatewayServer
} from '../api.js'

const { t } = useI18n()

const PROVIDERS = ['apisix', 'kong']

const gateways = ref([])
const loading = ref(true)
const testingId = ref('')
const defaultingId = ref('')

const notice = ref(null)
let noticeTimer = null

function showNotice(type, text) {
  notice.value = { type, text }
  if (noticeTimer) clearTimeout(noticeTimer)
  noticeTimer = setTimeout(() => {
    notice.value = null
  }, 4000)
}

function providerLabel(provider) {
  return t(`gateway.provider.${provider === 'kong' ? 'kong' : 'apisix'}`)
}

function statusLabel(status) {
  if (status === 'online') return t('gateway.status.online')
  if (status === 'offline') return t('gateway.status.offline')
  return t('gateway.status.unknown')
}

async function load() {
  loading.value = true
  try {
    gateways.value = await fetchGatewayServers()
  } catch (e) {
    showNotice('error', e.message || t('gateway.alert.loadFailed'))
  }
  loading.value = false
}

async function runTest(gw) {
  testingId.value = gw.id
  try {
    const result = await testGatewayServer(gw.id)
    gw.status = result.online ? 'online' : 'offline'
    gw.status_message = result.message
    gw.latency_ms = result.latency_ms
    if (result.online) {
      showNotice('success', `${gw.name}: ${result.message}`)
    } else {
      showNotice('error', `${gw.name}: ${result.message}`)
    }
  } catch (e) {
    showNotice('error', e.message || t('gateway.alert.loadFailed'))
  }
  testingId.value = ''
}

async function makeDefault(gw) {
  defaultingId.value = gw.id
  try {
    await setDefaultGatewayServer(gw.id)
    gateways.value.forEach((item) => {
      item.is_default = item.id === gw.id
    })
    showNotice('success', t('gateway.alert.setDefaultSuccess'))
  } catch (e) {
    showNotice('error', e.message || t('gateway.alert.saveFailed'))
  }
  defaultingId.value = ''
}

const showingForm = ref(false)
const editingId = ref('')
const editingHasKey = ref(false)
const saving = ref(false)
const formError = ref('')
const formTest = ref(null)
const formTesting = ref(false)

const emptyForm = () => ({
  name: '',
  provider: 'apisix',
  admin_url: '',
  admin_key: '',
  accesskey_header: '',
  auth_grpc_addrs: [''],
  audit_grpc_addrs: [''],
  description: ''
})

const form = ref(emptyForm())

const adminUrlPlaceholder = computed(() =>
  form.value.provider === 'kong' ? 'http://127.0.0.1:8001' : 'http://127.0.0.1:9180'
)

const providerHint = computed(() =>
  form.value.provider === 'kong' ? t('gateway.form.providerHintKong') : t('gateway.form.providerHintApisix')
)

const canSave = computed(() => form.value.name.trim() !== '' && form.value.admin_url.trim() !== '')

function resetTestResult() {
  formTest.value = null
}

function openCreate() {
  editingId.value = ''
  editingHasKey.value = false
  form.value = emptyForm()
  formError.value = ''
  resetTestResult()
  showingForm.value = true
}

function openEdit(gw) {
  editingId.value = gw.id
  editingHasKey.value = !!gw.has_admin_key
  form.value = {
    name: gw.name,
    provider: gw.provider === 'kong' ? 'kong' : 'apisix',
    admin_url: gw.admin_url,
    admin_key: '',
    accesskey_header: gw.accesskey_header || '',
    auth_grpc_addrs: gw.auth_grpc_addrs.length ? gw.auth_grpc_addrs.slice() : [''],
    audit_grpc_addrs: gw.audit_grpc_addrs.length ? gw.audit_grpc_addrs.slice() : [''],
    description: gw.description || ''
  }
  formError.value = ''
  resetTestResult()
  showingForm.value = true
}

function closeForm() {
  showingForm.value = false
}

function cleanAddrs(list) {
  return list.map((addr) => addr.trim()).filter((addr) => addr !== '')
}

function payload() {
  return {
    name: form.value.name.trim(),
    provider: form.value.provider,
    admin_url: form.value.admin_url.trim(),
    admin_key: form.value.admin_key,
    accesskey_header: form.value.accesskey_header.trim(),
    auth_grpc_addrs: cleanAddrs(form.value.auth_grpc_addrs),
    audit_grpc_addrs: cleanAddrs(form.value.audit_grpc_addrs),
    description: form.value.description.trim()
  }
}

function addAddr(field) {
  form.value[field].push('')
}

function removeAddr(field, index) {
  form.value[field].splice(index, 1)
  if (form.value[field].length === 0) form.value[field].push('')
}

async function runFormTest() {
  if (!canSave.value) return
  formTesting.value = true
  formTest.value = null
  try {
    const result = await checkGatewayServer({ id: editingId.value, ...payload() })
    const suffix = result.latency_ms ? ` (${result.latency_ms} ms)` : ''
    formTest.value = { ok: result.online, message: `${result.message}${suffix}` }
  } catch (e) {
    formTest.value = { ok: false, message: e.message || t('gateway.form.testFail') }
  }
  formTesting.value = false
}

async function saveForm() {
  if (!canSave.value || saving.value) return
  const isEdit = !!editingId.value
  saving.value = true
  formError.value = ''
  try {
    if (isEdit) {
      await updateGatewayServer(editingId.value, payload())
    } else {
      await createGatewayServer(payload())
    }
    closeForm()
    showNotice('success', isEdit ? t('gateway.alert.updateSuccess') : t('gateway.alert.createSuccess'))
    await load()
  } catch (e) {
    formError.value = e.message || t('gateway.alert.saveFailed')
  }
  saving.value = false
}

const deleteTarget = ref(null)
const deleting = ref(false)

function askDelete(gw) {
  deleteTarget.value = gw
}

async function confirmDelete() {
  if (!deleteTarget.value || deleting.value) return
  deleting.value = true
  try {
    await deleteGatewayServer(deleteTarget.value.id)
    deleteTarget.value = null
    showNotice('success', t('gateway.alert.deleteSuccess'))
    await load()
  } catch (e) {
    showNotice('error', e.message || t('gateway.alert.deleteFailed'))
  }
  deleting.value = false
}

onMounted(load)

onBeforeUnmount(() => {
  if (noticeTimer) clearTimeout(noticeTimer)
})
</script>

<template>
  <section class="gateways">
    <header class="gateways__head">
      <div class="gateways__heading">
        <h1 class="page__title">{{ t('gateway.title') }}</h1>
        <p class="gateways__desc">{{ t('gateway.desc') }}</p>
      </div>
      <div class="gateways__actions">
        <button class="btn btn--ghost" type="button" :disabled="loading" @click="load">
          {{ loading ? t('gateway.refreshing') : t('gateway.refresh') }}
        </button>
        <button class="btn btn--primary" type="button" @click="openCreate">
          {{ t('gateway.add') }}
        </button>
      </div>
    </header>

    <p
      v-if="notice"
      class="gateways__notice"
      :class="`gateways__notice--${notice.type}`"
    >
      {{ notice.text }}
    </p>

    <div v-if="loading && gateways.length === 0" class="gateways__empty">
      {{ t('gateway.loading') }}
    </div>
    <div v-else-if="gateways.length === 0" class="gateways__empty">
      <p class="gateways__empty-title">{{ t('gateway.empty') }}</p>
      <p class="gateways__empty-hint">{{ t('gateway.emptyHint') }}</p>
      <button class="btn btn--primary" type="button" @click="openCreate">
        {{ t('gateway.add') }}
      </button>
    </div>

    <div v-else class="gateway-grid">
      <article
        v-for="gw in gateways"
        :key="gw.id"
        class="gateway-card"
        :class="{ 'gateway-card--default': gw.is_default }"
      >
        <div class="gateway-card__head">
          <span class="gateway-card__provider">{{ providerLabel(gw.provider) }}</span>
          <span class="gateway-card__status" :class="`gateway-card__status--${gw.status}`">
            <span class="gateway-card__dot"></span>
            {{ statusLabel(gw.status) }}
            <span v-if="gw.status === 'online' && gw.latency_ms" class="gateway-card__latency">
              {{ t('gateway.card.latency', { ms: gw.latency_ms }) }}
            </span>
          </span>
        </div>

        <div class="gateway-card__title-row">
          <h2 class="gateway-card__name">{{ gw.name }}</h2>
          <span v-if="gw.is_default" class="gateway-card__tag">{{ t('gateway.defaultTag') }}</span>
        </div>

        <p class="gateway-card__desc">
          {{ gw.description || t('gateway.card.noDescription') }}
        </p>

        <dl class="gateway-card__meta">
          <div class="gateway-card__row">
            <dt>{{ t('gateway.card.adminUrl') }}</dt>
            <dd class="gateway-card__mono">{{ gw.admin_url }}</dd>
          </div>
          <div class="gateway-card__row">
            <dt>{{ t('gateway.card.accessKeyHeader') }}</dt>
            <dd>{{ gw.accesskey_header || 'X-Access-Key' }}</dd>
          </div>
          <div class="gateway-card__row">
            <dt>{{ t('gateway.card.authGrpcAddrs') }}</dt>
            <dd class="gateway-card__addrs">
              <template v-if="gw.auth_grpc_addrs.length">
                <span v-for="addr in gw.auth_grpc_addrs" :key="addr" class="gateway-card__addr">{{ addr }}</span>
              </template>
              <span v-else class="gateway-card__muted">{{ t('gateway.card.noGrpc') }}</span>
            </dd>
          </div>
          <div class="gateway-card__row">
            <dt>{{ t('gateway.card.auditGrpcAddrs') }}</dt>
            <dd class="gateway-card__addrs">
              <template v-if="gw.audit_grpc_addrs.length">
                <span v-for="addr in gw.audit_grpc_addrs" :key="addr" class="gateway-card__addr">{{ addr }}</span>
              </template>
              <span v-else class="gateway-card__muted">{{ t('gateway.card.noGrpc') }}</span>
            </dd>
          </div>
        </dl>

        <p v-if="gw.status === 'offline' && gw.status_message" class="gateway-card__error">
          {{ gw.status_message }}
        </p>

        <div class="gateway-card__actions">
          <button
            class="btn btn--ghost btn--sm"
            type="button"
            :disabled="testingId === gw.id"
            @click="runTest(gw)"
          >
            {{ testingId === gw.id ? t('gateway.actions.testing') : t('gateway.actions.test') }}
          </button>
          <button
            v-if="!gw.is_default"
            class="btn btn--ghost btn--sm"
            type="button"
            :disabled="defaultingId === gw.id"
            @click="makeDefault(gw)"
          >
            {{ t('gateway.actions.setDefault') }}
          </button>
          <span class="gateway-card__spacer"></span>
          <button class="btn btn--ghost btn--sm" type="button" @click="openEdit(gw)">
            {{ t('gateway.actions.edit') }}
          </button>
          <button class="btn btn--danger-ghost btn--sm" type="button" @click="askDelete(gw)">
            {{ t('common.delete') }}
          </button>
        </div>
      </article>
    </div>

    <Teleport to="body">
      <div v-if="showingForm" class="dialog-overlay" @click.self="closeForm">
        <div class="dialog">
          <h2 class="dialog__title">
            {{ editingId ? t('gateway.form.editTitle') : t('gateway.form.createTitle') }}
          </h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.name') }}</label>
              <input
                v-model="form.name"
                class="dialog__input"
                type="text"
                :placeholder="t('gateway.form.namePlaceholder')"
              />
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.provider') }}</label>
              <select v-model="form.provider" class="dialog__input">
                <option v-for="p in PROVIDERS" :key="p" :value="p">
                  {{ t(`gateway.provider.${p}`) }}
                </option>
              </select>
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.adminUrl') }}</label>
              <input
                v-model="form.admin_url"
                class="dialog__input"
                type="text"
                :placeholder="adminUrlPlaceholder"
              />
              <span class="dialog__hint">{{ providerHint }}</span>
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.adminKey') }}</label>
              <input
                v-model="form.admin_key"
                class="dialog__input"
                type="password"
                :placeholder="editingHasKey ? t('gateway.form.adminKeyKeep') : t('gateway.form.adminKeyPlaceholder')"
              />
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.accessKeyHeader') }}</label>
              <input
                v-model="form.accesskey_header"
                class="dialog__input"
                type="text"
                :placeholder="t('gateway.form.accessKeyHeaderPlaceholder')"
              />
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.authGrpcAddrs') }}</label>
              <div v-for="(_, i) in form.auth_grpc_addrs" :key="i" class="addr-row">
                <input
                  v-model="form.auth_grpc_addrs[i]"
                  class="dialog__input addr-row__input"
                  type="text"
                  :placeholder="t('gateway.form.grpcPlaceholder')"
                  @input="resetTestResult"
                />
                <button
                  class="addr-row__remove"
                  type="button"
                  :title="t('common.delete')"
                  @click="removeAddr('auth_grpc_addrs', i)"
                >
                  ×
                </button>
              </div>
              <button class="addr-row__add" type="button" @click="addAddr('auth_grpc_addrs')">
                {{ t('gateway.form.addAddress') }}
              </button>
              <span class="dialog__hint">{{ t('gateway.form.grpcHint') }}</span>
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.auditGrpcAddrs') }}</label>
              <div v-for="(_, i) in form.audit_grpc_addrs" :key="i" class="addr-row">
                <input
                  v-model="form.audit_grpc_addrs[i]"
                  class="dialog__input addr-row__input"
                  type="text"
                  :placeholder="t('gateway.form.auditGrpcPlaceholder')"
                  @input="resetTestResult"
                />
                <button
                  class="addr-row__remove"
                  type="button"
                  :title="t('common.delete')"
                  @click="removeAddr('audit_grpc_addrs', i)"
                >
                  ×
                </button>
              </div>
              <button class="addr-row__add" type="button" @click="addAddr('audit_grpc_addrs')">
                {{ t('gateway.form.addAddress') }}
              </button>
              <span class="dialog__hint">{{ t('gateway.form.auditGrpcHint') }}</span>
            </div>

            <div class="dialog__group">
              <label class="dialog__label">{{ t('gateway.form.description') }}</label>
              <textarea
                v-model="form.description"
                class="dialog__input dialog__textarea"
                rows="2"
                :placeholder="t('gateway.form.descriptionPlaceholder')"
              ></textarea>
            </div>

            <p
              v-if="formTest"
              class="dialog__test-result"
              :class="formTest.ok ? 'dialog__test-result--ok' : 'dialog__test-result--err'"
            >
              {{ formTest.message }}
            </p>
            <p v-if="formError" class="dialog__error">{{ formError }}</p>
          </div>

          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="closeForm">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--secondary"
              type="button"
              :disabled="!canSave || formTesting"
              @click="runFormTest"
            >
              {{ formTesting ? t('gateway.form.testing') : t('gateway.form.test') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!canSave || saving"
              @click="saveForm"
            >
              {{ saving ? t('common.saving') : t('common.save') }}
            </button>
          </div>
        </div>
      </div>

      <div v-if="deleteTarget" class="dialog-overlay" @click.self="deleteTarget = null">
        <div class="dialog dialog--narrow">
          <h2 class="dialog__title">{{ t('gateway.delete.title') }}</h2>
          <p class="dialog__desc">{{ t('gateway.delete.desc', { name: deleteTarget.name }) }}</p>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="deleteTarget = null">
              {{ t('common.cancel') }}
            </button>
            <button class="btn btn--danger" type="button" :disabled="deleting" @click="confirmDelete">
              {{ deleting ? t('gateway.delete.deleting') : t('gateway.delete.confirm') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<style scoped>
.gateways {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.gateways__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.gateways__heading {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.gateways__desc {
  max-width: 640px;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-muted);
}

.gateways__actions {
  display: flex;
  gap: 10px;
}

.gateways__notice {
  padding: 10px 14px;
  font-size: 13px;
  border-radius: var(--radius);
  border: 1px solid var(--border);
}

.gateways__notice--success {
  color: #1c6b4a;
  background: rgba(28, 107, 74, 0.08);
  border-color: rgba(28, 107, 74, 0.35);
}

.gateways__notice--error {
  color: var(--danger);
  background: rgba(180, 35, 24, 0.06);
  border-color: rgba(180, 35, 24, 0.35);
}

.gateways__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  padding: 56px 24px;
  text-align: center;
  background: var(--surface);
  border: 1px dashed var(--border);
  border-radius: var(--radius);
  font-size: 14px;
  color: var(--text-muted);
}

.gateways__empty-title {
  font-size: 16px;
  font-weight: 600;
  color: var(--text);
}

.gateways__empty-hint {
  font-size: 13px;
  color: var(--text-muted);
  margin-bottom: 6px;
}

.gateway-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
  gap: 16px;
}

.gateway-card {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-top: 3px solid var(--border);
  border-radius: var(--radius);
  transition: border-color 0.15s, transform 0.15s;
}

.gateway-card:hover {
  border-color: rgba(10, 61, 122, 0.25);
}

.gateway-card--default {
  border-top-color: var(--xauat-accent);
}

.gateway-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
}

.gateway-card__provider {
  padding: 3px 10px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.03em;
  color: var(--xauat-blue);
  background: var(--active-bg);
  border-radius: 999px;
}

.gateway-card__status {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  font-weight: 600;
  color: var(--text-muted);
}

.gateway-card__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  background: var(--text-muted);
}

.gateway-card__status--online {
  color: #1c8f5a;
}

.gateway-card__status--online .gateway-card__dot {
  background: #1c8f5a;
}

.gateway-card__status--offline {
  color: var(--danger);
}

.gateway-card__status--offline .gateway-card__dot {
  background: var(--danger);
}

.gateway-card__latency {
  font-weight: 400;
  color: var(--text-muted);
}

.gateway-card__title-row {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.gateway-card__name {
  font-size: 17px;
  font-weight: 700;
  color: var(--text);
  word-break: break-word;
}

.gateway-card__tag {
  padding: 2px 8px;
  font-size: 11px;
  font-weight: 700;
  color: #7a5300;
  background: rgba(242, 169, 0, 0.16);
  border-radius: 999px;
}

html[data-theme="dark"] .gateway-card__tag {
  color: #f6c76a;
}

.gateway-card__desc {
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-muted);
  min-height: 21px;
}

.gateway-card__meta {
  display: flex;
  flex-direction: column;
  gap: 8px;
  padding: 12px 0;
  border-top: 1px solid var(--border);
  border-bottom: 1px solid var(--border);
}

.gateway-card__row {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 13px;
}

.gateway-card__row dt {
  flex: 0 0 112px;
  color: var(--text-muted);
}

.gateway-card__row dd {
  flex: 1;
  min-width: 0;
  color: var(--text);
  word-break: break-all;
}

.gateway-card__mono {
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 12px;
}

.gateway-card__addrs {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
}

.gateway-card__addr {
  padding: 1px 6px;
  font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
  font-size: 11px;
  color: var(--xauat-blue);
  background: var(--panel-bg);
  border: 1px solid var(--border);
  border-radius: 4px;
}

.gateway-card__muted {
  color: var(--text-muted);
}

.gateway-card__error {
  font-size: 12px;
  line-height: 1.5;
  color: var(--danger);
  word-break: break-word;
}

.gateway-card__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.gateway-card__spacer {
  flex: 1;
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

.btn--sm {
  padding: 5px 12px;
  font-size: 12px;
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

.btn--ghost:hover:not(:disabled) {
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

.btn--danger-ghost {
  color: var(--danger);
  background: transparent;
  border-color: var(--border);
}

.btn--danger-ghost:hover:not(:disabled) {
  border-color: var(--danger);
  background: rgba(180, 35, 24, 0.06);
}

.dialog-overlay {
  position: fixed;
  inset: 0;
  z-index: 1000;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.35);
}

.dialog {
  display: flex;
  flex-direction: column;
  width: 90%;
  max-width: 560px;
  max-height: 90vh;
  padding: 32px;
  background: var(--surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.dialog--narrow {
  max-width: 460px;
}

.dialog__title {
  flex-shrink: 0;
  margin-bottom: 20px;
  font-size: 20px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.dialog__form {
  display: flex;
  flex-direction: column;
  gap: 16px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
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
  min-height: 64px;
  resize: vertical;
  font-family: inherit;
}

.dialog__hint {
  font-size: 12px;
  line-height: 1.5;
  color: var(--text-muted);
}

.dialog__test-result {
  font-size: 13px;
  line-height: 1.5;
  word-break: break-word;
}

.dialog__test-result--ok {
  color: #1c8f5a;
}

.dialog__test-result--err {
  color: var(--danger);
}

.dialog__error {
  font-size: 12px;
  color: var(--danger);
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  flex-shrink: 0;
  margin-top: 24px;
}

.dialog__desc {
  margin-bottom: 16px;
  font-size: 14px;
  line-height: 1.6;
  color: var(--text-muted);
}

.addr-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.addr-row__input {
  flex: 1;
  min-width: 0;
}

.addr-row__remove {
  flex-shrink: 0;
  width: 30px;
  height: 30px;
  font-size: 18px;
  line-height: 1;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 50%;
  cursor: pointer;
  transition: color 0.15s, border-color 0.15s;
}

.addr-row__remove:hover {
  color: var(--danger);
  border-color: var(--danger);
}

.addr-row__add {
  align-self: flex-start;
  padding: 3px 10px;
  font-size: 12px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px dashed var(--border);
  border-radius: 999px;
  cursor: pointer;
  transition: border-color 0.15s, background 0.15s;
}

.addr-row__add:hover {
  background: var(--hover-bg);
  border-color: rgba(10, 61, 122, 0.35);
}

@media (max-width: 720px) {
  .gateway-grid {
    grid-template-columns: 1fr;
  }
}
</style>

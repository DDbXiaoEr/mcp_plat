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
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { fetchHistory, fetchAccessKeys, fetchServers } from '../api.js'

const { t, n } = useI18n()

const accessKeys = ref([])
const servers = ref([])
const history = ref({ list: [], total: 0, page: 1, page_size: 20 })
const loading = ref(true)
const page = ref(1)

const filters = ref({
  start: '',
  end: '',
  access_key: '',
  server_id: ''
})

function formatTime(dateStr) {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  const pad = (n) => String(n).padStart(2, '0')
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(d.getMinutes())}`
}

function getKeyLabel(r) {
  return r.access_key_name || '—'
}

function getServerName(id) {
  const s = servers.value.find((s) => s.id === id)
  return s ? s.name : id
}

async function loadHistory() {
  loading.value = true
  try {
    const params = { page: page.value, page_size: 20 }
    if (filters.value.access_key) params.access_key = filters.value.access_key
    if (filters.value.server_id) params.server_id = filters.value.server_id
    if (filters.value.start) params.start = filters.value.start
    if (filters.value.end) params.end = filters.value.end
    history.value = await fetchHistory(params)
  } catch {
    history.value = { list: [], total: 0, page: 1, page_size: 20 }
  }
  loading.value = false
}

function search() {
  page.value = 1
  loadHistory()
}

function reset() {
  filters.value = { start: '', end: '', access_key: '', server_id: '' }
  search()
}

function goPage(p) {
  page.value = p
  loadHistory()
}

const totalPages = computed(() => Math.ceil(history.value.total / history.value.page_size) || 1)

const pageNumbers = computed(() => {
  const total = totalPages.value
  const current = page.value
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  const set = new Set([1, total, current - 1, current, current + 1])
  const sorted = [...set].filter((p) => p >= 1 && p <= total).sort((a, b) => a - b)
  const pages = []
  let prev = 0
  for (const p of sorted) {
    if (p - prev > 1) pages.push('...')
    pages.push(p)
    prev = p
  }
  return pages
})

const jumpPage = ref('')

function goJump() {
  const p = Math.min(totalPages.value, Math.max(1, Number(jumpPage.value) || 1))
  goPage(p)
  jumpPage.value = ''
}

onMounted(async () => {
  try {
    const [keys, svrs] = await Promise.all([fetchAccessKeys(), fetchServers()])
    accessKeys.value = keys
    servers.value = svrs
  } catch {
    // ignore
  }
  await loadHistory()
})
</script>

<template>
  <section class="page">
    <h1 class="page__title">{{ t('history.title') }}</h1>
    <p class="page__hint">{{ t('history.hint') }}</p>

    <div class="filters">
      <label class="filters__field">
        <span>{{ t('history.startTime') }}</span>
        <input v-model="filters.start" type="date" />
      </label>
      <label class="filters__field">
        <span>{{ t('history.endTime') }}</span>
        <input v-model="filters.end" type="date" />
      </label>
      <label class="filters__field">
        <span>AccessKey</span>
        <select v-model="filters.access_key">
          <option value="">{{ t('common.all') }}</option>
          <option v-for="key in accessKeys" :key="key.id" :value="key.key">
            {{ key.name }}
          </option>
        </select>
      </label>
      <label class="filters__field">
        <span>{{ t('history.server') }}</span>
        <select v-model="filters.server_id">
          <option value="">{{ t('common.all') }}</option>
          <option v-for="s in servers" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
      </label>
      <button class="filters__reset" type="button" @click="reset">{{ t('history.reset') }}</button>
      <button class="filters__search" type="button" @click="search">{{ t('history.query') }}</button>
    </div>

    <div v-if="loading" class="history__status">{{ t('history.loading') }}</div>

    <template v-else>
      <table class="history">
        <thead>
          <tr>
            <th>{{ t('history.time') }}</th>
            <th>AccessKey</th>
            <th>{{ t('history.server') }}</th>
            <th>{{ t('history.tool') }}</th>
            <th>{{ t('history.sourceIp') }}</th>
            <th>{{ t('common.status') }}</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in history.list" :key="r.id">
            <td>{{ formatTime(r.created_at) }}</td>
            <td>{{ getKeyLabel(r) }}</td>
            <td>{{ getServerName(r.server_id) }}</td>
            <td>{{ r.tool_name }}</td>
            <td>{{ r.client_ip || '—' }}</td>
            <td>{{ r.success ? t('history.success') : t('history.failed') }}</td>
          </tr>
          <tr v-if="!history.list.length">
            <td class="history__empty" colspan="6">{{ t('history.empty') }}</td>
          </tr>
        </tbody>
      </table>

      <div v-if="totalPages > 1" class="history__pager">
        <button :disabled="page <= 1" @click="goPage(page - 1)">{{ t('history.prevPage') }}</button>
        <template v-for="(p, i) in pageNumbers" :key="i">
          <span v-if="p === '...'" class="history__ellipsis">…</span>
          <button
            v-else
            class="history__page"
            :class="{ 'is-active': p === page }"
            @click="goPage(p)"
          >
            {{ p }}
          </button>
        </template>
        <button :disabled="page >= totalPages" @click="goPage(page + 1)">{{ t('history.nextPage') }}</button>
        <span class="history__jump">
          {{ t('history.jumpPrefix') }}
          <input v-model="jumpPage" type="number" min="1" :max="totalPages" @keyup.enter="goJump" />
          {{ t('history.jumpUnit') }}
          <button type="button" @click="goJump">{{ t('history.jumpBtn') }}</button>
        </span>
        <span>{{ t('history.total', { count: n(history.total) }) }}</span>
      </div>
    </template>
  </section>
</template>

<style scoped>
.page__hint {
  margin-top: 6px;
  font-size: 14px;
  color: var(--text-muted);
}

.history__status {
  margin-top: 40px;
  text-align: center;
  color: var(--text-muted);
  font-size: 14px;
}

.filters {
  display: flex;
  align-items: flex-end;
  flex-wrap: wrap;
  gap: 16px;
  margin-top: 20px;
  padding: 18px 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.filters__field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  font-size: 13px;
  color: var(--text-muted);
}

.filters__field input,
.filters__field select {
  padding: 8px 12px;
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.filters__field input:focus,
.filters__field select:focus {
  border-color: var(--xauat-blue-light);
  box-shadow: 0 0 0 3px rgba(30, 95, 176, 0.12);
}

.filters__reset {
  padding: 9px 20px;
  font-size: 14px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.filters__reset:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.filters__search {
  padding: 9px 20px;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  background: var(--xauat-blue);
  border: 1px solid var(--xauat-blue);
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s;
}

.filters__search:hover {
  background: var(--xauat-blue-dark);
}

.history {
  width: 100%;
  margin-top: 20px;
  border-collapse: collapse;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.history th,
.history td {
  padding: 14px 20px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.history thead th {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg);
}

.history tbody tr:last-child td {
  border-bottom: none;
}

.history__empty {
  text-align: center;
  color: var(--text-muted);
}

.history__pager {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 20px;
  font-size: 14px;
  color: var(--text-muted);
}

.history__pager button {
  padding: 8px 14px;
  font-size: 13px;
  font-weight: 600;
  color: var(--xauat-blue);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s;
}

.history__pager button:hover:not(:disabled) {
  background: rgba(10, 61, 122, 0.06);
}

.history__pager button:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.history__ellipsis {
  padding: 0 2px;
  color: var(--text-muted);
}

.history__page.is-active {
  color: #fff;
  background: var(--xauat-blue);
  border-color: var(--xauat-blue);
}

.history__jump {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-left: 8px;
}

.history__jump input {
  width: 56px;
  padding: 7px 8px;
  font-size: 13px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  outline: none;
  text-align: center;
}

.history__jump input:focus {
  border-color: var(--xauat-blue-light);
  box-shadow: 0 0 0 3px rgba(30, 95, 176, 0.12);
}
</style>

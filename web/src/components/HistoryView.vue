<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchHistory, fetchAccessKeys, fetchServers } from '../api.js'

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

function getKeyNameByKey(keyStr) {
  const k = accessKeys.value.find((k) => k.key === keyStr)
  return k ? k.name : (keyStr ? keyStr.substring(0, 16) + '...' : '')
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
    <h1 class="page__title">使用历史</h1>
    <p class="page__hint">近一周的调用记录</p>

    <div class="filters">
      <label class="filters__field">
        <span>开始时间</span>
        <input v-model="filters.start" type="date" />
      </label>
      <label class="filters__field">
        <span>结束时间</span>
        <input v-model="filters.end" type="date" />
      </label>
      <label class="filters__field">
        <span>AccessKey</span>
        <select v-model="filters.access_key">
          <option value="">全部</option>
          <option v-for="key in accessKeys" :key="key.id" :value="key.key">
            {{ key.name }}
          </option>
        </select>
      </label>
      <label class="filters__field">
        <span>MCP 服务器</span>
        <select v-model="filters.server_id">
          <option value="">全部</option>
          <option v-for="s in servers" :key="s.id" :value="s.id">{{ s.name }}</option>
        </select>
      </label>
      <button class="filters__reset" type="button" @click="reset">重置</button>
      <button class="filters__search" type="button" @click="search">查询</button>
    </div>

    <div v-if="loading" class="history__status">加载中...</div>

    <template v-else>
      <table class="history">
        <thead>
          <tr>
            <th>时间</th>
            <th>AccessKey</th>
            <th>MCP 服务器</th>
            <th>工具</th>
            <th>状态</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="r in history.list" :key="r.id">
            <td>{{ formatTime(r.created_at) }}</td>
            <td>{{ getKeyNameByKey(r.access_key) }}</td>
            <td>{{ getServerName(r.server_id) }}</td>
            <td>{{ r.tool_name }}</td>
            <td>{{ r.success ? '成功' : '失败' }}</td>
          </tr>
          <tr v-if="!history.list.length">
            <td class="history__empty" colspan="5">暂无匹配的记录</td>
          </tr>
        </tbody>
      </table>

      <div v-if="totalPages > 1" class="history__pager">
        <button :disabled="page <= 1" @click="goPage(page - 1)">上一页</button>
        <span>第 {{ page }} / {{ totalPages }} 页（共 {{ history.total }} 条）</span>
        <button :disabled="page >= totalPages" @click="goPage(page + 1)">下一页</button>
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
  gap: 16px;
  margin-top: 20px;
  font-size: 14px;
  color: var(--text-muted);
}

.history__pager button {
  padding: 8px 16px;
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
</style>

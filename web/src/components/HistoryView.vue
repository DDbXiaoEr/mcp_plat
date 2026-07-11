<script setup>
import { ref, computed } from 'vue'

// TODO: 后端接入后，替换为真实的使用历史接口数据
const KEY_NAMES = ['默认密钥 1', '默认密钥 2', '默认密钥 3']
const USAGES = [
  { server: '图书馆 MCP', tool: '文献检索' },
  { server: '图书馆 MCP', tool: '借阅查询' },
  { server: '教务 MCP', tool: '课表查询' },
  { server: '教务 MCP', tool: '成绩查询' },
  { server: '校园服务 MCP', tool: '通知公告' },
  { server: '一卡通 MCP', tool: '余额查询' },
  { server: '一卡通 MCP', tool: '消费记录查询' }
]
const SERVERS = [...new Set(USAGES.map((u) => u.server))]

function pick(list) {
  return list[Math.floor(Math.random() * list.length)]
}

function formatTime(date) {
  const pad = (n) => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

const WEEK_MS = 7 * 24 * 60 * 60 * 1000

function createRecords(n) {
  const now = Date.now()
  return Array.from({ length: n }, (_, i) => {
    const usage = pick(USAGES)
    return {
      id: i + 1,
      keyName: pick(KEY_NAMES),
      server: usage.server,
      tool: usage.tool,
      time: new Date(now - Math.floor(Math.random() * WEEK_MS))
    }
  }).sort((a, b) => b.time - a.time)
}

const records = ref(createRecords(12))

const filters = ref({
  start: '',
  end: '',
  keyName: '',
  server: ''
})

const filtered = computed(() =>
  records.value.filter((r) => {
    const { start, end, keyName, server } = filters.value
    if (keyName && r.keyName !== keyName) return false
    if (server && r.server !== server) return false
    if (start && r.time < new Date(`${start}T00:00:00`)) return false
    if (end && r.time > new Date(`${end}T23:59:59`)) return false
    return true
  })
)

function reset() {
  filters.value = { start: '', end: '', keyName: '', server: '' }
}
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
        <select v-model="filters.keyName">
          <option value="">全部</option>
          <option v-for="name in KEY_NAMES" :key="name" :value="name">
            {{ name }}
          </option>
        </select>
      </label>
      <label class="filters__field">
        <span>MCP 服务器</span>
        <select v-model="filters.server">
          <option value="">全部</option>
          <option v-for="s in SERVERS" :key="s" :value="s">{{ s }}</option>
        </select>
      </label>
      <button class="filters__reset" type="button" @click="reset">重置</button>
    </div>

    <table class="history">
      <thead>
        <tr>
          <th>时间</th>
          <th>AccessKey</th>
          <th>MCP 服务器</th>
          <th>工具</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="r in filtered" :key="r.id">
          <td>{{ formatTime(r.time) }}</td>
          <td>{{ r.keyName }}</td>
          <td>{{ r.server }}</td>
          <td>{{ r.tool }}</td>
        </tr>
        <tr v-if="!filtered.length">
          <td class="history__empty" colspan="4">暂无匹配的记录</td>
        </tr>
      </tbody>
    </table>
  </section>
</template>

<style scoped>
.page__hint {
  margin-top: 6px;
  font-size: 14px;
  color: var(--text-muted);
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
</style>

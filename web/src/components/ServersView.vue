<script setup>
import { ref, computed, onMounted } from 'vue'
import { fetchServers, createServer, fetchServerTools } from '../api.js'

const servers = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    const data = await fetchServers()
    servers.value = data.map((s) => ({
      ...s,
      tools: parseTools(s.tools)
    }))
  } catch {
    // ignore load error
  }
  loading.value = false
})

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

const PROTOCOLS = ['SSE', 'Streamable HTTP']

const selectedId = ref('')
const selected = ref(null)

function select(server) {
  selectedId.value = server.id
  selected.value = server
}

const selectedTools = computed(() => parseTools(selected.value?.tools))

const showingCreate = ref(false)
const createForm = ref({ name: '', address: '', department: '', protocol: 'SSE', tools: '' })
const fetchingTools = ref(false)
const fetchToolsError = ref('')

function openCreate() {
  createForm.value = { name: '', address: '', department: '', protocol: 'SSE', tools: '' }
  fetchToolsError.value = ''
  showingCreate.value = true
}

async function fetchTools() {
  const address = createForm.value.address.trim()
  fetchToolsError.value = ''
  if (!address) {
    fetchToolsError.value = '请先填写 MCP 服务器路径'
    return
  }
  fetchingTools.value = true
  try {
    const data = await fetchServerTools({
      address,
      protocol: createForm.value.protocol
    })
    createForm.value.tools = JSON.stringify(data.tools || [])
  } catch (err) {
    fetchToolsError.value = err.message || '获取工具列表失败'
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
      department: createForm.value.department.trim(),
      protocol: createForm.value.protocol,
      tools: JSON.stringify(tools)
    })
    server.tools = tools
    servers.value.push(server)
    showingCreate.value = false
    select(server)
  } catch {
    // ignore
  }
}

const showingPublish = ref(false)
const publishSelected = ref([])

function openPublish() {
  publishSelected.value = []
  showingPublish.value = true
}

function togglePublishServer(id) {
  const i = publishSelected.value.indexOf(id)
  if (i === -1) publishSelected.value.push(id)
  else publishSelected.value.splice(i, 1)
}

function confirmPublish() {
  alert('发布成功')
  showingPublish.value = false
}
</script>

<template>
  <section class="servers">
    <div class="servers__head">
      <h1 class="page__title">MCP 服务器管理</h1>
      <div class="servers__actions">
        <button class="btn btn--primary" type="button" @click="openCreate">
          增加
        </button>
        <button class="btn btn--secondary" type="button" @click="openPublish">
          发布
        </button>
      </div>
    </div>

    <div class="server-list">
      <button
        v-for="server in servers"
        :key="server.id"
        class="server-list__item"
        :class="{ 'server-list__item--active': selectedId === server.id }"
        type="button"
        @click="select(server)"
      >
        <span class="server-list__name">{{ server.name }}</span>
        <span class="server-list__dept">{{ server.department }}</span>
      </button>
    </div>

    <div v-if="selected" class="drawer">
      <div class="drawer__overlay" @click="selected = null"></div>
      <aside class="drawer__panel">
        <header class="drawer__head">
          <h2 class="drawer__title">{{ selected.name }}</h2>
          <button class="drawer__close" type="button" aria-label="关闭" @click="selected = null">
            ×
          </button>
        </header>
        <div class="drawer__body">
          <dl class="server-detail__grid">
            <div class="server-detail__row">
              <dt>MCP 服务器地址</dt>
              <dd>{{ selected.address }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>负责部门</dt>
              <dd>{{ selected.department }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>协议类型</dt>
              <dd>{{ selected.protocol }}</dd>
            </div>
            <div class="server-detail__row">
              <dt>工具列表</dt>
              <dd>
                <ul class="server-detail__tools">
                  <li
                    v-for="tool in selectedTools"
                    :key="tool.name"
                    :title="tool.description || undefined"
                  >
                    {{ tool.name }}
                  </li>
                </ul>
              </dd>
            </div>
          </dl>
        </div>
      </aside>
    </div>

    <Teleport to="body">
      <div v-if="showingCreate" class="dialog-overlay" @click.self="showingCreate = false">
        <div class="dialog">
          <h2 class="dialog__title">新增 MCP 服务器</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">名称</label>
              <input
                v-model="createForm.name"
                class="dialog__input"
                type="text"
                placeholder="请输入服务器名称"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">
                MCP 服务器地址
                <span class="dialog__help">
                  ?
                  <span class="dialog__tooltip">
                    在 API 网关统一管理的场景下，地址应填写对应的 URL 路径，例如
                    <code>/jwc/mcp</code>（由网关转发到实际服务）；独立部署时填写完整地址，例如
                    <code>https://mcp.xauat.edu.cn/jwc</code>。
                  </span>
                </span>
              </label>
              <input
                v-model="createForm.address"
                class="dialog__input"
                type="text"
                placeholder="https://mcp.xauat.edu.cn/..."
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">负责部门</label>
              <input
                v-model="createForm.department"
                class="dialog__input"
                type="text"
                placeholder="请输入负责部门"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">协议类型</label>
              <select v-model="createForm.protocol" class="dialog__input">
                <option v-for="p in PROTOCOLS" :key="p" :value="p">{{ p }}</option>
              </select>
            </div>
            <div class="dialog__group">
              <div class="dialog__label-row">
                <label class="dialog__label">工具列表</label>
                <button
                  class="dialog__fetch"
                  type="button"
                  :disabled="fetchingTools"
                  @click="fetchTools"
                >
                  {{ fetchingTools ? '获取中...' : '自动获取工具列表' }}
                </button>
              </div>
              <textarea
                v-model="createForm.tools"
                class="dialog__input dialog__textarea"
                placeholder="多个工具以逗号或换行分隔"
              ></textarea>
              <p v-if="fetchToolsError" class="dialog__error">{{ fetchToolsError }}</p>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingCreate = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!createForm.name.trim()"
              @click="confirmCreate"
            >
              确认新增
            </button>
          </div>
        </div>
      </div>

      <div v-if="showingPublish" class="dialog-overlay" @click.self="showingPublish = false">
        <div class="dialog">
          <h2 class="dialog__title">发布到 API 网关</h2>
          <p class="dialog__desc">选择要发布的 MCP 服务器：</p>
          <div class="publish-list">
            <label
              v-for="server in servers"
              :key="server.id"
              class="publish-list__item"
            >
              <input
                type="checkbox"
                :checked="publishSelected.includes(server.id)"
                @change="togglePublishServer(server.id)"
              />
              <span class="publish-list__name">{{ server.name }}</span>
              <span class="publish-list__addr">{{ server.address }}</span>
            </label>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingPublish = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="publishSelected.length === 0"
              @click="confirmPublish"
            >
              确认发布
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

.server-list__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  text-align: left;
  transition: background 0.2s, border-color 0.2s;
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

.server-list__dept {
  font-size: 13px;
  color: var(--text-muted);
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

.server-detail__tools {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  list-style: none;
}

.server-detail__tools li {
  padding: 4px 12px;
  background: rgba(30, 95, 176, 0.1);
  border-radius: 999px;
  font-size: 13px;
  color: var(--xauat-blue);
  cursor: default;
}

.server-detail__tools li[title] {
  cursor: help;
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
  box-shadow: var(--shadow);
}

.dialog__title {
  font-size: 20px;
  font-weight: 700;
  color: var(--xauat-blue);
  margin-bottom: 20px;
}

.dialog__form {
  display: flex;
  flex-direction: column;
  gap: 16px;
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

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}

.dialog__desc {
  font-size: 14px;
  color: var(--text-muted);
  margin-bottom: 16px;
}

.publish-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 300px;
  overflow-y: auto;
}

.publish-list__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.publish-list__item:hover {
  background: rgba(10, 61, 122, 0.04);
  border-color: rgba(10, 61, 122, 0.25);
}

.publish-list__item input {
  width: 16px;
  height: 16px;
  accent-color: var(--xauat-blue);
  cursor: pointer;
  flex-shrink: 0;
}

.publish-list__name {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.publish-list__addr {
  font-size: 12px;
  color: var(--text-muted);
  margin-left: auto;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 200px;
}
</style>

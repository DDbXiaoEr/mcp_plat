<script setup>
import { ref } from 'vue'

// TODO: 接入后端后替换为真实数据
const servers = ref([
  {
    id: 'srv-jwc',
    name: '教务系统 MCP',
    address: 'https://mcp.xauat.edu.cn/jwc',
    department: '教务处',
    protocol: 'SSE',
    tools: ['查询课表', '成绩查询', '选课信息', '考试安排']
  },
  {
    id: 'srv-lib',
    name: '图书馆 MCP',
    address: 'https://mcp.xauat.edu.cn/library',
    department: '图书馆',
    protocol: 'Streamable HTTP',
    tools: ['馆藏检索', '借阅记录', '座位预约']
  },
  {
    id: 'srv-hr',
    name: '人事系统 MCP',
    address: 'https://mcp.xauat.edu.cn/hr',
    department: '人事处',
    protocol: 'stdio',
    tools: ['职工信息查询', '考勤统计']
  },
  {
    id: 'srv-fin',
    name: '财务系统 MCP',
    address: 'https://mcp.xauat.edu.cn/finance',
    department: '财务处',
    protocol: 'SSE',
    tools: ['报销进度', '工资查询', '经费余额']
  }
])

const PROTOCOLS = ['SSE', 'Streamable HTTP', 'stdio']

const selectedId = ref('')
const selected = ref(null)

function select(server) {
  selectedId.value = server.id
  selected.value = server
}

const showingCreate = ref(false)
const createForm = ref({ name: '', address: '', department: '', protocol: 'SSE', tools: '' })
const fetchingTools = ref(false)

function openCreate() {
  createForm.value = { name: '', address: '', department: '', protocol: 'SSE', tools: '' }
  showingCreate.value = true
}

function fetchTools() {
  fetchingTools.value = true
  // TODO: 接入后端后改为从 MCP 服务器地址实际拉取工具列表
  setTimeout(() => {
    createForm.value.tools = ['校园网账号查询', '网络套餐办理', '在线故障报修', '流量使用统计', '宽带缴费'].join('\n')
    fetchingTools.value = false
  }, 600)
}

function confirmCreate() {
  const name = createForm.value.name.trim()
  if (!name) return
  // TODO: 接入后端后改为调用创建接口
  const server = {
    id: `srv-${Date.now()}`,
    name,
    address: createForm.value.address.trim(),
    department: createForm.value.department.trim(),
    protocol: createForm.value.protocol,
    tools: createForm.value.tools
      .split(/[,，\n]/)
      .map((t) => t.trim())
      .filter(Boolean)
  }
  servers.value.push(server)
  showingCreate.value = false
  select(server)
}
</script>

<template>
  <section class="servers">
    <div class="servers__head">
      <h1 class="page__title">MCP 服务器管理</h1>
      <button class="btn btn--primary" type="button" @click="openCreate">
        增加
      </button>
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

    <div v-if="selected" class="server-detail">
      <h2 class="server-detail__title">{{ selected.name }}</h2>
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
              <li v-for="tool in selected.tools" :key="tool">{{ tool }}</li>
            </ul>
          </dd>
        </div>
      </dl>
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
              <label class="dialog__label">MCP 服务器地址</label>
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

.server-detail {
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  padding: 24px;
}

.server-detail__title {
  font-size: 20px;
  font-weight: 700;
  color: var(--xauat-blue);
  margin-bottom: 16px;
}

.server-detail__grid {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.server-detail__row {
  display: grid;
  grid-template-columns: 140px 1fr;
  gap: 16px;
  align-items: start;
}

.server-detail__row dt {
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
}

.server-detail__row dd {
  font-size: 14px;
  color: var(--text);
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

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}
</style>

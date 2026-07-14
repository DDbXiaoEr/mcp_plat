<script setup>
import { ref, computed } from 'vue'

const roles = ref([
  { id: 1, name: '超级管理员', description: '拥有全部权限', servers: [1, 2, 3, 4, 5, 6, 7, 8, 9] },
  { id: 2, name: 'MCP管理员', description: '管理MCP服务器与工具', servers: [1, 2, 7] },
  { id: 3, name: '普通用户', description: '基础访问权限', servers: [3] }
])

const allServers = ref([
  { id: 1, name: '校园网络管理' },
  { id: 2, name: '教学资源服务' },
  { id: 3, name: '学生事务服务' },
  { id: 4, name: '财务管理服务' },
  { id: 5, name: '人事管理服务' },
  { id: 6, name: '科研管理系统' },
  { id: 7, name: '图书馆服务' },
  { id: 8, name: '一卡通服务' },
  { id: 9, name: '邮件与消息服务' }
])

const users = ref([
  { id: 1, username: 'admin', role: '超级管理员' },
  { id: 2, username: 'zhangsan', role: 'MCP管理员' },
  { id: 3, username: 'lisi', role: '普通用户' },
  { id: 4, username: 'wangwu', role: '普通用户' },
  { id: 5, username: 'zhaoliu', role: '' },
  { id: 6, username: 'sunqi', role: '' },
  { id: 7, username: 'zhouba', role: '' }
])

const showingRoleCreate = ref(false)
const showingUserCreate = ref(false)
const rbacTab = ref('roles')

const roleForm = ref({ name: '', description: '' })
const userForm = ref({ username: '', role: '' })

const editingRole = ref(null)
const editingUser = ref(null)

const shuttleLeftSelected = ref([])
const shuttleRightSelected = ref([])

const shuttleLeft = computed(() => {
  if (!editingRole.value) return []
  const selectedIds = editingRole.value.servers
  return allServers.value.filter((s) => !selectedIds.includes(s.id))
})

const shuttleRight = computed(() => {
  if (!editingRole.value) return []
  const selectedIds = editingRole.value.servers
  return allServers.value.filter((s) => selectedIds.includes(s.id))
})

function toggleShuttleLeft(id) {
  const idx = shuttleLeftSelected.value.indexOf(id)
  if (idx === -1) {
    shuttleLeftSelected.value.push(id)
  } else {
    shuttleLeftSelected.value.splice(idx, 1)
  }
}

function toggleShuttleRight(id) {
  const idx = shuttleRightSelected.value.indexOf(id)
  if (idx === -1) {
    shuttleRightSelected.value.push(id)
  } else {
    shuttleRightSelected.value.splice(idx, 1)
  }
}

function moveToRight() {
  if (shuttleLeftSelected.value.length === 0) return
  editingRole.value.servers.push(...shuttleLeftSelected.value)
  shuttleLeftSelected.value = []
}

function moveToLeft() {
  if (shuttleRightSelected.value.length === 0) return
  editingRole.value.servers = editingRole.value.servers.filter(
    (id) => !shuttleRightSelected.value.includes(id)
  )
  shuttleRightSelected.value = []
}

function moveAllToRight() {
  shuttleLeft.value.forEach((s) => {
    if (!editingRole.value.servers.includes(s.id)) {
      editingRole.value.servers.push(s.id)
    }
  })
  shuttleLeftSelected.value = []
}

function moveAllToLeft() {
  editingRole.value.servers = []
  shuttleRightSelected.value = []
}

const assigningRole = ref(null)
const userAssignLeftSel = ref([])
const userAssignRightSel = ref([])

const userAssignLeft = computed(() => {
  if (!assigningRole.value) return []
  return users.value.filter((u) => u.role !== assigningRole.value.name)
})

const userAssignRight = computed(() => {
  if (!assigningRole.value) return []
  return users.value.filter((u) => u.role === assigningRole.value.name)
})

function toggleUserAssignLeft(id) {
  const idx = userAssignLeftSel.value.indexOf(id)
  if (idx === -1) {
    userAssignLeftSel.value.push(id)
  } else {
    userAssignLeftSel.value.splice(idx, 1)
  }
}

function toggleUserAssignRight(id) {
  const idx = userAssignRightSel.value.indexOf(id)
  if (idx === -1) {
    userAssignRightSel.value.push(id)
  } else {
    userAssignRightSel.value.splice(idx, 1)
  }
}

function userMoveToRight() {
  if (userAssignLeftSel.value.length === 0) return
  userAssignLeftSel.value.forEach((id) => {
    const u = users.value.find((u) => u.id === id)
    if (u) u.role = assigningRole.value.name
  })
  userAssignLeftSel.value = []
}

function userMoveToLeft() {
  if (userAssignRightSel.value.length === 0) return
  userAssignRightSel.value.forEach((id) => {
    const u = users.value.find((u) => u.id === id)
    if (u) u.role = ''
  })
  userAssignRightSel.value = []
}

function userMoveAllToRight() {
  userAssignLeft.value.forEach((u) => {
    u.role = assigningRole.value.name
  })
  userAssignLeftSel.value = []
}

function userMoveAllToLeft() {
  userAssignRight.value.forEach((u) => {
    u.role = ''
  })
  userAssignRightSel.value = []
}

function openAssignUsers(role) {
  assigningRole.value = role
  userAssignLeftSel.value = []
  userAssignRightSel.value = []
}

function openRoleCreate() {
  roleForm.value = { name: '', description: '' }
  showingRoleCreate.value = true
}

function confirmRoleCreate() {
  const name = roleForm.value.name.trim()
  if (!name) return
  roles.value.push({
    id: Date.now(),
    name,
    description: roleForm.value.description.trim()
  })
  showingRoleCreate.value = false
}

function startEditRole(role) {
  editingRole.value = { ...role, servers: [...(role.servers || [])] }
  shuttleLeftSelected.value = []
  shuttleRightSelected.value = []
}

function saveRole() {
  if (!editingRole.value.name.trim()) return
  const idx = roles.value.findIndex((r) => r.id === editingRole.value.id)
  if (idx !== -1) {
    roles.value[idx] = { ...editingRole.value }
  }
  editingRole.value = null
}

function removeRole(role) {
  const idx = roles.value.findIndex((r) => r.id === role.id)
  if (idx !== -1) roles.value.splice(idx, 1)
}

function openUserCreate() {
  userForm.value = { username: '', role: '' }
  showingUserCreate.value = true
}

function confirmUserCreate() {
  const username = userForm.value.username.trim()
  if (!username) return
  users.value.push({
    id: Date.now(),
    username,
    role: userForm.value.role
  })
  showingUserCreate.value = false
}

function startEditUser(user) {
  editingUser.value = { ...user }
}

function saveUser() {
  if (!editingUser.value.username.trim()) return
  const idx = users.value.findIndex((u) => u.id === editingUser.value.id)
  if (idx !== -1) {
    users.value[idx] = { ...editingUser.value }
  }
  editingUser.value = null
}

function removeUser(user) {
  const idx = users.value.findIndex((u) => u.id === user.id)
  if (idx !== -1) users.value.splice(idx, 1)
}
</script>

<template>
  <section class="rbac">
    <h1 class="page__title">RBAC设置</h1>

    <div class="rbac__tabs">
      <button
        class="rbac__tab"
        :class="{ 'rbac__tab--active': rbacTab === 'roles' }"
        type="button"
        @click="rbacTab = 'roles'"
      >角色管理</button>
      <button
        class="rbac__tab"
        :class="{ 'rbac__tab--active': rbacTab === 'users' }"
        type="button"
        @click="rbacTab = 'users'"
      >用户管理</button>
    </div>

    <!-- 角色管理 -->
    <div v-if="rbacTab === 'roles'" class="rbac__section">
      <div class="rbac__head">
        <h2 class="rbac__subtitle">角色管理</h2>
        <button class="btn btn--primary" type="button" @click="openRoleCreate">
          新增角色
        </button>
      </div>

      <table class="rbac__table">
        <thead>
          <tr>
            <th>角色名称</th>
            <th>描述</th>
            <th>可使用的MCP服务器</th>
            <th class="rbac__col-action">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="roles.length === 0">
            <td colspan="4" class="rbac__empty">暂无角色</td>
          </tr>
          <tr v-for="role in roles" :key="role.id">
            <td>{{ role.name }}</td>
            <td>{{ role.description }}</td>
            <td>
              <div v-if="role.servers && role.servers.length" class="rbac__tags">
                <span
                  v-for="sId in role.servers"
                  :key="sId"
                  class="rbac__tag"
                >{{ allServers.find((s) => s.id === sId)?.name || '' }}</span>
              </div>
              <span v-else class="rbac__tag-empty">无</span>
            </td>
            <td class="rbac__col-action">
              <button class="rbac__edit" type="button" @click="startEditRole(role)">
                编辑
              </button>
              <button class="rbac__edit" type="button" @click="openAssignUsers(role)">
                分配用户
              </button>
              <button class="rbac__del" type="button" @click="removeRole(role)">
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 用户管理 -->
    <div v-if="rbacTab === 'users'" class="rbac__section">
      <div class="rbac__head">
        <h2 class="rbac__subtitle">用户管理</h2>
        <button class="btn btn--primary" type="button" @click="openUserCreate">
          新增用户
        </button>
      </div>

      <table class="rbac__table">
        <thead>
          <tr>
            <th>用户名</th>
            <th>角色</th>
            <th class="rbac__col-action">操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-if="users.length === 0">
            <td colspan="3" class="rbac__empty">暂无用户</td>
          </tr>
          <tr v-for="user in users" :key="user.id">
            <td>{{ user.username }}</td>
            <td>{{ user.role || '未分配' }}</td>
            <td class="rbac__col-action">
              <button class="rbac__edit" type="button" @click="startEditUser(user)">
                编辑
              </button>
              <button class="rbac__del" type="button" @click="removeUser(user)">
                删除
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- 新增/编辑角色对话框 -->
    <Teleport to="body">
      <div v-if="showingRoleCreate" class="dialog-overlay" @click.self="showingRoleCreate = false">
        <div class="dialog">
          <h2 class="dialog__title">新增角色</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">角色名称</label>
              <input
                v-model="roleForm.name"
                class="dialog__input"
                type="text"
                placeholder="请输入角色名称"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">描述</label>
              <input
                v-model="roleForm.description"
                class="dialog__input"
                type="text"
                placeholder="请输入角色描述"
              />
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingRoleCreate = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!roleForm.name.trim()"
              @click="confirmRoleCreate"
            >
              确认新增
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="editingRole" class="dialog-overlay" @click.self="editingRole = null">
        <div class="dialog dialog--shuttle">
          <h2 class="dialog__title">编辑角色</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">角色名称</label>
              <input
                v-model="editingRole.name"
                class="dialog__input"
                type="text"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">描述</label>
              <input
                v-model="editingRole.description"
                class="dialog__input"
                type="text"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">可使用的MCP服务器</label>
              <div class="shuttle">
                <div class="shuttle__panel">
                  <div class="shuttle__head">可选服务器</div>
                  <div class="shuttle__list">
                    <div
                      v-for="srv in shuttleLeft"
                      :key="srv.id"
                      class="shuttle__item"
                      :class="{ 'shuttle__item--sel': shuttleLeftSelected.includes(srv.id) }"
                      @click="toggleShuttleLeft(srv.id)"
                    >
                      {{ srv.name }}
                    </div>
                  </div>
                </div>
                <div class="shuttle__controls">
                  <button
                    class="shuttle__btn"
                    type="button"
                    title="添加选中"
                    @click="moveToRight"
                  >&gt;</button>
                  <button
                    class="shuttle__btn"
                    type="button"
                    title="全部添加"
                    @click="moveAllToRight"
                  >&gt;&gt;</button>
                  <button
                    class="shuttle__btn"
                    type="button"
                    title="移除选中"
                    @click="moveToLeft"
                  >&lt;</button>
                  <button
                    class="shuttle__btn"
                    type="button"
                    title="全部移除"
                    @click="moveAllToLeft"
                  >&lt;&lt;</button>
                </div>
                <div class="shuttle__panel">
                  <div class="shuttle__head">已选服务器</div>
                  <div class="shuttle__list">
                    <div
                      v-for="srv in shuttleRight"
                      :key="srv.id"
                      class="shuttle__item"
                      :class="{ 'shuttle__item--sel': shuttleRightSelected.includes(srv.id) }"
                      @click="toggleShuttleRight(srv.id)"
                    >
                      {{ srv.name }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="editingRole = null">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!editingRole.name.trim()"
              @click="saveRole"
            >
              保存
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 新增/编辑用户对话框 -->
    <Teleport to="body">
      <div v-if="showingUserCreate" class="dialog-overlay" @click.self="showingUserCreate = false">
        <div class="dialog">
          <h2 class="dialog__title">新增用户</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">用户名</label>
              <input
                v-model="userForm.username"
                class="dialog__input"
                type="text"
                placeholder="请输入用户名"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">角色</label>
              <select v-model="userForm.role" class="dialog__input">
                <option value="">请选择角色</option>
                <option v-for="role in roles" :key="role.id" :value="role.name">
                  {{ role.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingUserCreate = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!userForm.username.trim()"
              @click="confirmUserCreate"
            >
              确认新增
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="editingUser" class="dialog-overlay" @click.self="editingUser = null">
        <div class="dialog">
          <h2 class="dialog__title">编辑用户</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">用户名</label>
              <input
                v-model="editingUser.username"
                class="dialog__input"
                type="text"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">角色</label>
              <select v-model="editingUser.role" class="dialog__input">
                <option value="">请选择角色</option>
                <option v-for="role in roles" :key="role.id" :value="role.name">
                  {{ role.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="editingUser = null">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!editingUser.username.trim()"
              @click="saveUser"
            >
              保存
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 分配用户对话框 -->
    <Teleport to="body">
      <div v-if="assigningRole" class="dialog-overlay" @click.self="assigningRole = null">
        <div class="dialog dialog--shuttle">
          <h2 class="dialog__title">分配用户 - {{ assigningRole.name }}</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">选择用户</label>
              <div class="shuttle">
                <div class="shuttle__panel">
                  <div class="shuttle__head">未分配用户</div>
                  <div class="shuttle__list">
                    <div
                      v-for="u in userAssignLeft"
                      :key="u.id"
                      class="shuttle__item"
                      :class="{ 'shuttle__item--sel': userAssignLeftSel.includes(u.id) }"
                      @click="toggleUserAssignLeft(u.id)"
                    >
                      {{ u.username }}
                    </div>
                  </div>
                </div>
                <div class="shuttle__controls">
                  <button class="shuttle__btn" type="button" title="添加选中" @click="userMoveToRight">&gt;</button>
                  <button class="shuttle__btn" type="button" title="全部添加" @click="userMoveAllToRight">&gt;&gt;</button>
                  <button class="shuttle__btn" type="button" title="移除选中" @click="userMoveToLeft">&lt;</button>
                  <button class="shuttle__btn" type="button" title="全部移除" @click="userMoveAllToLeft">&lt;&lt;</button>
                </div>
                <div class="shuttle__panel">
                  <div class="shuttle__head">已分配用户</div>
                  <div class="shuttle__list">
                    <div
                      v-for="u in userAssignRight"
                      :key="u.id"
                      class="shuttle__item"
                      :class="{ 'shuttle__item--sel': userAssignRightSel.includes(u.id) }"
                      @click="toggleUserAssignRight(u.id)"
                    >
                      {{ u.username }}
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="assigningRole = null">
              取消
            </button>
            <button class="btn btn--primary" type="button" @click="assigningRole = null">
              完成
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </section>
</template>

<style scoped>
.rbac {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.rbac__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 24px;
  align-items: start;
}

.rbac__tabs {
  display: flex;
  gap: 4px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  padding: 4px;
  width: fit-content;
}

.rbac__tab {
  padding: 8px 20px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.rbac__tab:hover {
  color: var(--text);
  background: rgba(10, 61, 122, 0.04);
}

.rbac__tab--active {
  color: #fff;
  background: var(--xauat-blue);
}

.rbac__tab--active:hover {
  color: #fff;
  background: var(--xauat-blue-light);
}

.rbac__section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.rbac__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.rbac__subtitle {
  font-size: 18px;
  font-weight: 700;
  color: var(--text);
}

.rbac__table {
  width: 100%;
  border-collapse: collapse;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.rbac__table th,
.rbac__table td {
  padding: 14px 20px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.rbac__table thead th {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg);
}

.rbac__table tbody tr:last-child td {
  border-bottom: none;
}

.rbac__empty {
  text-align: center;
  color: var(--text-muted);
  padding: 32px 20px !important;
}

.rbac__col-action {
  width: 220px;
  white-space: nowrap;
}

.rbac__edit,
.rbac__del {
  flex: none;
  padding: 4px 12px;
  font-size: 13px;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
  margin-right: 8px;
}

.rbac__edit:hover,
.rbac__del:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.rbac__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.rbac__tag {
  padding: 2px 10px;
  font-size: 12px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.08);
  border-radius: 999px;
  white-space: nowrap;
}

.rbac__tag-empty {
  font-size: 12px;
  color: var(--text-muted);
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

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
}

.dialog--shuttle {
  max-width: 680px;
}

.shuttle {
  display: flex;
  align-items: stretch;
  gap: 12px;
}

.shuttle__panel {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.shuttle__head {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  padding: 8px 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: 10px 10px 0 0;
  border-bottom: none;
}

.shuttle__list {
  flex: 1;
  min-height: 180px;
  max-height: 220px;
  overflow-y: auto;
  border: 1px solid var(--border);
  border-radius: 0 0 10px 10px;
  background: var(--surface);
}

.shuttle__item {
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.shuttle__item:hover {
  background: rgba(10, 61, 122, 0.04);
}

.shuttle__item--sel {
  color: #fff;
  background: var(--xauat-blue);
}

.shuttle__item--sel:hover {
  background: var(--xauat-blue-light);
}

.shuttle__controls {
  display: flex;
  flex-direction: column;
  justify-content: center;
  gap: 8px;
  padding: 0 4px;
}

.shuttle__btn {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  color: var(--xauat-blue);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.15s, border-color 0.15s;
}

.shuttle__btn:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}
</style>

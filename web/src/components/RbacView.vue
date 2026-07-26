<script setup>
import { ref, onMounted, computed } from 'vue'
import {
  fetchRoles, createRole, updateRole, deleteRole,
  fetchUsers, createUser, updateUser, deleteUser,
  fetchRoleUsers, assignRoleUsers, fetchServers
} from '../api.js'

const tab = ref('roles')

const servers = ref([])
const serverMap = computed(() => {
  const m = {}
  servers.value.forEach((s) => { m[s.id] = s.name })
  return m
})

const roles = ref([])
const loadingRoles = ref(true)

const users = ref([])
const loadingUsers = ref(true)
const userFilter = ref('all')

const filteredUsers = computed(() => {
  if (userFilter.value === 'assigned') return users.value.filter((u) => u.role_name)
  if (userFilter.value === 'unassigned') return users.value.filter((u) => !u.role_name)
  return users.value
})

async function loadRoles() {
  try {
    roles.value = await fetchRoles()
  } catch {
    // ignore
  }
  loadingRoles.value = false
}

async function loadUsers() {
  try {
    users.value = await fetchUsers()
  } catch {
    // ignore
  }
  loadingUsers.value = false
}

onMounted(async () => {
  await Promise.all([
    loadRoles(),
    loadUsers(),
    fetchServers().then((d) => { servers.value = d }).catch(() => {})
  ])
})

function serverNames(role) {
  return (role.server_ids || []).map((id) => serverMap.value[id] || id.toString())
}

const showingRoleForm = ref(false)
const roleFormMode = ref('create')
const roleForm = ref({ name: '', description: '', servers: {} })

function openRoleCreate() {
  roleFormMode.value = 'create'
  roleForm.value = { name: '', description: '', servers: {} }
  showingRoleForm.value = true
}

function openRoleEdit(role) {
  roleFormMode.value = 'edit'
  roleForm.value = {
    id: role.id,
    name: role.name,
    description: role.description,
    servers: {}
  }
  ;(role.server_ids || []).forEach((id) => { roleForm.value.servers[id] = true })
  showingRoleForm.value = true
}

async function confirmRole() {
  const body = {
    name: roleForm.value.name.trim(),
    description: roleForm.value.description.trim(),
    server_ids: Object.keys(roleForm.value.servers)
      .filter((k) => roleForm.value.servers[k])
  }
  if (!body.name) return
  try {
    if (roleFormMode.value === 'create') {
      const role = await createRole(body)
      roles.value.unshift(role)
    } else {
      await updateRole(roleForm.value.id, body)
      await loadRoles()
    }
  } catch {
    // ignore
  }
  showingRoleForm.value = false
}

async function removeRole(role) {
  try {
    await deleteRole(role.id)
    roles.value = roles.value.filter((r) => r.id !== role.id)
  } catch {
    // ignore
  }
}

const showingAssign = ref(false)
const assigningRole = ref(null)
const assignUserIds = ref([])

const assignableUsers = computed(() => users.value.filter((u) => !u.role_name))

const leftUsers = computed(() =>
  assignableUsers.value.filter((u) => !assignUserIds.value.includes(u.id))
)

const rightUsers = computed(() =>
  users.value.filter((u) => assignUserIds.value.includes(u.id))
)

const isAllLeftChecked = computed(() => {
  if (leftUsers.value.length === 0) return false
  return leftUsers.value.every((u) => assignSelectedIds.value.includes(u.id))
})

const isAllRightChecked = computed(() => {
  if (rightUsers.value.length === 0) return false
  return rightUsers.value.every((u) => assignSelectedIds.value.includes(u.id))
})

const assignSelectedIds = ref([])

function toggleSelect(id) {
  const idx = assignSelectedIds.value.indexOf(id)
  if (idx >= 0) {
    assignSelectedIds.value.splice(idx, 1)
  } else {
    assignSelectedIds.value.push(id)
  }
}

function toggleSelectAll(list) {
  const ids = list.map((u) => u.id)
  if (ids.every((id) => assignSelectedIds.value.includes(id))) {
    assignSelectedIds.value = assignSelectedIds.value.filter((id) => !ids.includes(id))
  } else {
    ids.forEach((id) => {
      if (!assignSelectedIds.value.includes(id)) assignSelectedIds.value.push(id)
    })
  }
}

function moveRight() {
  const selected = assignSelectedIds.value.filter((id) =>
    leftUsers.value.some((u) => u.id === id)
  )
  selected.forEach((id) => {
    if (!assignUserIds.value.includes(id)) assignUserIds.value.push(id)
  })
  assignSelectedIds.value = assignSelectedIds.value.filter((id) => !selected.includes(id))
}

function moveLeft() {
  const selected = assignSelectedIds.value.filter((id) =>
    rightUsers.value.some((u) => u.id === id)
  )
  assignUserIds.value = assignUserIds.value.filter((id) => !selected.includes(id))
  assignSelectedIds.value = assignSelectedIds.value.filter((id) => !selected.includes(id))
}

async function openAssign(role) {
  assigningRole.value = role
  try {
    const assigned = await fetchRoleUsers(role.id)
    assignUserIds.value = (assigned || []).map((u) => u.id)
  } catch {
    assignUserIds.value = []
  }
  assignSelectedIds.value = []
  await loadUsers()
  showingAssign.value = true
}

function toggleAssignUser(id) {
  const idx = assignUserIds.value.indexOf(id)
  if (idx >= 0) {
    assignUserIds.value.splice(idx, 1)
  } else {
    assignUserIds.value.push(id)
  }
}

async function confirmAssign() {
  try {
    await assignRoleUsers(assigningRole.value.id, assignUserIds.value)
    await loadRoles()
    await loadUsers()
  } catch {
    // ignore
  }
  showingAssign.value = false
  assigningRole.value = null
}

const showingUserForm = ref(false)
const userFormMode = ref('create')
const userForm = ref({ username: '', password: '', role_id: 0 })

function openUserCreate() {
  userFormMode.value = 'create'
  userForm.value = { username: '', password: '', role_id: 0 }
  showingUserForm.value = true
}

function openUserEdit(user) {
  userFormMode.value = 'edit'
  userForm.value = {
    id: user.id,
    username: user.username,
    password: '',
    role_id: user.role_id || 0
  }
  showingUserForm.value = true
}

async function confirmUser() {
  const body = {
    username: userForm.value.username.trim()
  }
  if (!body.username) return

  if (userForm.value.password) {
    body.password = userForm.value.password.trim()
  }

  if (userFormMode.value === 'create') {
    body.password = body.password || userForm.value.password.trim()
    if (!body.password) return
    body.role_id = userForm.value.role_id > 0 ? userForm.value.role_id : null
    try {
      const user = await createUser(body)
      users.value.unshift(user)
    } catch {
      // ignore
    }
  } else {
    body.role_id = userForm.value.role_id > 0 ? userForm.value.role_id : 0
    try {
      await updateUser(userForm.value.id, body)
      await loadUsers()
      await loadRoles()
    } catch {
      // ignore
    }
  }
  showingUserForm.value = false
}

async function removeUser(user) {
  try {
    await deleteUser(user.id)
    users.value = users.value.filter((u) => u.id !== user.id)
    await loadRoles()
  } catch {
    // ignore
  }
}
</script>

<template>
  <section class="rbac">
    <div class="rbac__head">
      <h1 class="page__title">RBAC 设置</h1>
      <div class="rbac__actions">
        <button
          v-if="tab === 'roles'"
          class="btn btn--primary"
          type="button"
          @click="openRoleCreate"
        >
          增加角色
        </button>
        <button
          v-if="tab === 'users'"
          class="btn btn--primary"
          type="button"
          @click="openUserCreate"
        >
          增加用户
        </button>
      </div>
    </div>

    <div class="rbac__tabs">
      <button
        class="rbac__tab"
        :class="{ 'rbac__tab--active': tab === 'roles' }"
        type="button"
        @click="tab = 'roles'"
      >
        角色管理
      </button>
      <button
        class="rbac__tab"
        :class="{ 'rbac__tab--active': tab === 'users' }"
        type="button"
        @click="tab = 'users'"
      >
        用户管理
      </button>
    </div>

    <table v-if="tab === 'roles' && !loadingRoles" class="table">
      <thead>
        <tr>
          <th>名称</th>
          <th>描述</th>
          <th>MCP 服务器</th>
          <th class="table__col-num">用户数</th>
          <th class="table__col-action">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="roles.length === 0">
          <td colspan="5" class="table__empty">暂无角色</td>
        </tr>
        <tr v-for="role in roles" :key="role.id">
          <td class="table__name">{{ role.name }}</td>
          <td class="table__desc">{{ role.description || '—' }}</td>
          <td>
            <div class="table__tags">
              <template v-if="serverNames(role).length">
                <span
                  v-for="name in serverNames(role)"
                  :key="name"
                  class="table__tag"
                >{{ name }}</span>
              </template>
              <span v-else class="table__muted">—</span>
            </div>
          </td>
          <td class="table__col-num">{{ role.user_count }}</td>
          <td class="table__col-action">
            <button class="table__btn" type="button" @click="openAssign(role)">
              分配用户
            </button>
            <button class="table__btn" type="button" @click="openRoleEdit(role)">
              编辑
            </button>
            <button class="table__btn" type="button" @click="removeRole(role)">
              删除
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <table v-if="tab === 'users' && !loadingUsers" class="table">
      <thead>
        <tr>
          <th>
            学号/工号
            <select v-model="userFilter" class="table__filter">
              <option value="all">全部</option>
              <option value="assigned">已分配角色</option>
              <option value="unassigned">未分配角色</option>
            </select>
          </th>
          <th>角色</th>
          <th class="table__col-action">操作</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="filteredUsers.length === 0">
          <td colspan="3" class="table__empty">暂无用户</td>
        </tr>
        <tr v-for="user in filteredUsers" :key="user.id">
          <td class="table__name">{{ user.uid }}</td>
          <td>
            <span v-if="user.role_name" class="table__tag">{{ user.role_name }}</span>
            <span v-else class="table__muted">未分配</span>
          </td>
          <td class="table__col-action">
            <button class="table__btn" type="button" @click="openUserEdit(user)">
              编辑
            </button>
            <button class="table__btn" type="button" @click="removeUser(user)">
              删除
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <Teleport to="body">
      <div v-if="showingRoleForm" class="dialog-overlay" @click.self="showingRoleForm = false">
        <div class="dialog">
          <h2 class="dialog__title">
            {{ roleFormMode === 'create' ? '新增角色' : '编辑角色' }}
          </h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">名称</label>
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
            <div class="dialog__group">
              <label class="dialog__label">MCP 服务器权限</label>
              <div class="dialog__checklist">
                <label
                  v-for="s in servers"
                  :key="s.id"
                  class="dialog__check"
                >
                  <input
                    type="checkbox"
                    :checked="!!roleForm.servers[s.id]"
                    @change="roleForm.servers[s.id] = !roleForm.servers[s.id]"
                  />
                  <span>{{ s.name }}</span>
                </label>
                <span v-if="servers.length === 0" class="table__muted">暂无 MCP 服务器</span>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingRoleForm = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!roleForm.name.trim()"
              @click="confirmRole"
            >
              保存
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showingAssign" class="dialog-overlay" @click.self="showingAssign = false">
        <div class="dialog dialog--wide">
          <h2 class="dialog__title">分配用户 - {{ assigningRole?.name }}</h2>
          <div class="shuttle">
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">可选用户</span>
                <span class="shuttle__count">{{ leftUsers.length }}</span>
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllLeftChecked"
                    @change="toggleSelectAll(leftUsers)"
                  />
                  <span>全选</span>
                </label>
                <label
                  v-for="user in leftUsers"
                  :key="user.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="assignSelectedIds.includes(user.id)"
                    @change="toggleSelect(user.id)"
                  />
                  <span>{{ user.uid }}</span>
                </label>
                <span v-if="leftUsers.length === 0" class="table__muted shuttle__empty">
                  暂无可选用户
                </span>
              </div>
            </div>
            <div class="shuttle__actions">
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!assignSelectedIds.some((id) => leftUsers.some((u) => u.id === id))"
                @click="moveRight"
              >
                &gt;
              </button>
              <button
                class="shuttle__btn"
                type="button"
                :disabled="!assignSelectedIds.some((id) => rightUsers.some((u) => u.id === id))"
                @click="moveLeft"
              >
                &lt;
              </button>
            </div>
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">已选用户</span>
                <span class="shuttle__count">{{ rightUsers.length }}</span>
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllRightChecked"
                    @change="toggleSelectAll(rightUsers)"
                  />
                  <span>全选</span>
                </label>
                <label
                  v-for="user in rightUsers"
                  :key="user.id"
                  class="shuttle__item"
                >
                  <input
                    type="checkbox"
                    :checked="assignSelectedIds.includes(user.id)"
                    @change="toggleSelect(user.id)"
                  />
                  <span>{{ user.uid }}</span>
                </label>
                <span v-if="rightUsers.length === 0" class="table__muted shuttle__empty">
                  暂未选择
                </span>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingAssign = false">
              取消
            </button>
            <button class="btn btn--primary" type="button" @click="confirmAssign">
              保存
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showingUserForm" class="dialog-overlay" @click.self="showingUserForm = false">
        <div class="dialog">
          <h2 class="dialog__title">
            {{ userFormMode === 'create' ? '新增用户' : '编辑用户' }}
          </h2>
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
              <label class="dialog__label">{{ userFormMode === 'create' ? '密码' : '新密码（留空不修改）' }}</label>
              <input
                v-model="userForm.password"
                class="dialog__input"
                type="password"
                placeholder="请输入密码"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">角色</label>
              <select v-model="userForm.role_id" class="dialog__input">
                <option :value="0">无角色</option>
                <option v-for="role in roles" :key="role.id" :value="role.id">
                  {{ role.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingUserForm = false">
              取消
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!userForm.username.trim() || (userFormMode === 'create' && !userForm.password.trim())"
              @click="confirmUser"
            >
              保存
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

.rbac__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.rbac__tabs {
  display: flex;
  gap: 0;
  border-bottom: 2px solid var(--border);
}

.rbac__tab {
  padding: 10px 24px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-bottom: 2px solid transparent;
  margin-bottom: -2px;
  cursor: pointer;
  transition: color 0.2s, border-color 0.2s;
}

.rbac__tab:hover {
  color: var(--xauat-blue);
}

.rbac__tab--active {
  color: var(--xauat-blue);
  border-bottom-color: var(--xauat-blue);
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

.table {
  width: 100%;
  border-collapse: collapse;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.table th,
.table td {
  padding: 14px 20px;
  text-align: left;
  border-bottom: 1px solid var(--border);
}

.table thead th {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: var(--bg);
}

.table tbody tr:last-child td {
  border-bottom: none;
}

.table__empty {
  text-align: center;
  color: var(--text-muted);
  padding: 32px 20px !important;
}

.table__col-num {
  width: 80px;
}

.table__col-action {
  width: 240px;
}

.table__filter {
  margin-left: 10px;
  padding: 4px 8px;
  font-size: 12px;
  font-weight: 400;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  outline: none;
  cursor: pointer;
}

.table__name {
  font-weight: 600;
}

.table__desc {
  color: var(--text-muted);
  max-width: 260px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.table__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.table__tag {
  padding: 3px 10px;
  background: rgba(30, 95, 176, 0.1);
  border-radius: 999px;
  font-size: 12px;
  color: var(--xauat-blue);
}

.table__muted {
  font-size: 13px;
  color: var(--text-muted);
}

.table__btn {
  padding: 4px 12px;
  font-size: 13px;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.table__btn:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}

.table__btn + .table__btn {
  margin-left: 6px;
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
  max-height: 85vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
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

.dialog__checklist {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-height: 160px;
  overflow-y: auto;
}

.dialog__checklist--scroll {
  flex: 1;
  min-height: 0;
  max-height: none;
}

.dialog__check {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: var(--text);
  cursor: pointer;
}

.dialog__check input {
  accent-color: var(--xauat-blue);
}

.dialog__check--all {
  padding-bottom: 8px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 4px;
  font-weight: 600;
}

.dialog__actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 24px;
  flex-shrink: 0;
}

.dialog--wide {
  max-width: 640px;
}

.shuttle {
  display: flex;
  gap: 12px;
  flex: 1;
  min-height: 0;
}

.shuttle__panel {
  flex: 1;
  min-width: 0;
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
</style>

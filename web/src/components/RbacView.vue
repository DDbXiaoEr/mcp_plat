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
import { ref, onMounted, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import {
  fetchRoles, createRole, updateRole, deleteRole,
  fetchUsers, updateUser, deleteUser, batchCreateUsers,
  fetchRoleUsers, assignRoleUsers, fetchServers
} from '../api.js'

const { t } = useI18n()

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
const userKeyword = ref('')

const filteredUsers = computed(() => {
  let list = users.value
  if (userFilter.value === 'assigned') list = list.filter((u) => u.role_name)
  if (userFilter.value === 'unassigned') list = list.filter((u) => !u.role_name)
  return list
})

async function loadUsers() {
  try {
    users.value = await fetchUsers(undefined, userKeyword.value.trim())
  } catch {
    // ignore
  }
  loadingUsers.value = false
}

async function loadRoles() {
  try {
    roles.value = await fetchRoles()
  } catch {
    // ignore
  }
  loadingRoles.value = false
}

let searchTimer = null
watch(userKeyword, () => {
  clearTimeout(searchTimer)
  searchTimer = setTimeout(loadUsers, 300)
})

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
const assignSearchLeft = ref('')
const assignSearchRight = ref('')

const assignableUsers = computed(() => users.value.filter((u) => !u.role_name))

const leftUsers = computed(() => {
  const list = assignableUsers.value.filter((u) => !assignUserIds.value.includes(u.id))
  const kw = assignSearchLeft.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((u) => (u.uid || '').toLowerCase().startsWith(kw))
})

const rightUsers = computed(() => {
  const list = users.value.filter((u) => assignUserIds.value.includes(u.id))
  const kw = assignSearchRight.value.trim().toLowerCase()
  if (!kw) return list
  return list.filter((u) => (u.uid || '').toLowerCase().startsWith(kw))
})

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
  assignSearchLeft.value = ''
  assignSearchRight.value = ''
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
const userForm = ref({ id: null, username: '', password: '', role_id: 0 })

function openUserEdit(user) {
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
  body.role_id = userForm.value.role_id > 0 ? userForm.value.role_id : 0
  try {
    await updateUser(userForm.value.id, body)
    await loadUsers()
    await loadRoles()
  } catch {
    // ignore
  }
  showingUserForm.value = false
}

function emptyUserEntry() {
  return { username: '', password: '', name: '', email: '', uid: '', phone: '', organization: '' }
}

const showingBatchCreate = ref(false)
const batchSubmitting = ref(false)
const batchRoleId = ref(0)
const batchDraft = ref([])
const entryForm = ref(emptyUserEntry())

function openBatchCreate() {
  batchDraft.value = []
  batchRoleId.value = 0
  entryForm.value = emptyUserEntry()
  showingBatchCreate.value = true
}

function addToDraft() {
  const e = entryForm.value
  if (!e.username.trim() || !e.password.trim()) {
    alert(t('rbac.missingRequired'))
    return
  }
  if (batchDraft.value.some((u) => u.username === e.username.trim())) {
    alert(t('rbac.duplicateInList'))
    return
  }
  batchDraft.value.push({
    username: e.username.trim(),
    password: e.password.trim(),
    name: e.name.trim(),
    email: e.email.trim(),
    uid: e.uid.trim(),
    phone: e.phone.trim(),
    organization: e.organization.trim()
  })
  entryForm.value = emptyUserEntry()
}

function removeDraftItem(index) {
  batchDraft.value.splice(index, 1)
}

function clearDraft() {
  batchDraft.value = []
}

async function submitBatchCreate() {
  if (batchDraft.value.length === 0 || batchSubmitting.value) return
  batchSubmitting.value = true
  try {
    await batchCreateUsers(batchDraft.value, batchRoleId.value)
    await loadUsers()
    await loadRoles()
    showingBatchCreate.value = false
  } catch (e) {
    alert(e.message || t('rbac.createFailed'))
  } finally {
    batchSubmitting.value = false
  }
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
      <h1 class="page__title">{{ t('nav.rbac') }}</h1>
      <div class="rbac__actions">
        <button
          v-if="tab === 'roles'"
          class="btn btn--primary"
          type="button"
          @click="openRoleCreate"
        >
          {{ t('rbac.addRole') }}
        </button>
        <button
          v-if="tab === 'users'"
          class="btn btn--primary"
          type="button"
          @click="openBatchCreate"
        >
          {{ t('rbac.addUser') }}
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
        {{ t('rbac.rolesTab') }}
      </button>
      <button
        class="rbac__tab"
        :class="{ 'rbac__tab--active': tab === 'users' }"
        type="button"
        @click="tab = 'users'"
      >
        {{ t('rbac.usersTab') }}
      </button>
    </div>

    <table v-if="tab === 'roles' && !loadingRoles" class="table">
      <thead>
        <tr>
          <th>{{ t('common.name') }}</th>
          <th>{{ t('common.description') }}</th>
          <th>{{ t('rbac.mcpServers') }}</th>
          <th class="table__col-num">{{ t('rbac.userCount') }}</th>
          <th class="table__col-action">{{ t('common.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="roles.length === 0">
          <td colspan="5" class="table__empty">{{ t('rbac.noRoles') }}</td>
        </tr>
        <tr v-for="role in roles" :key="role.id">
          <td class="table__name">{{ role.name }}</td>
          <td class="table__desc">{{ role.description || t('common.emptyDash') }}</td>
          <td>
            <div class="table__tags">
              <template v-if="serverNames(role).length">
                <span
                  v-for="name in serverNames(role)"
                  :key="name"
                  class="table__tag"
                >{{ name }}</span>
              </template>
              <span v-else class="table__muted">{{ t('common.emptyDash') }}</span>
            </div>
          </td>
          <td class="table__col-num">{{ role.user_count }}</td>
          <td class="table__col-action">
            <button class="table__btn" type="button" @click="openAssign(role)">
              {{ t('rbac.assign') }}
            </button>
            <button class="table__btn" type="button" @click="openRoleEdit(role)">
              {{ t('common.edit') }}
            </button>
            <button class="table__btn" type="button" @click="removeRole(role)">
              {{ t('common.delete') }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <div v-if="tab === 'users'" class="rbac__search">
      <input
        v-model="userKeyword"
        class="rbac__search-input"
        type="text"
        :placeholder="t('rbac.uidSearchPlaceholder')"
      />
    </div>

    <table v-if="tab === 'users' && !loadingUsers" class="table">
      <thead>
        <tr>
          <th>
            {{ t('rbac.uidHeader') }}
            <select v-model="userFilter" class="table__filter">
              <option value="all">{{ t('common.all') }}</option>
              <option value="assigned">{{ t('rbac.filterAssigned') }}</option>
              <option value="unassigned">{{ t('rbac.filterUnassigned') }}</option>
            </select>
          </th>
          <th>{{ t('rbac.role') }}</th>
          <th class="table__col-action">{{ t('common.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="filteredUsers.length === 0">
          <td colspan="3" class="table__empty">{{ t('rbac.noUsers') }}</td>
        </tr>
        <tr v-for="user in filteredUsers" :key="user.id">
          <td class="table__name">{{ user.uid }}</td>
          <td>
            <span v-if="user.role_name" class="table__tag">{{ user.role_name }}</span>
            <span v-else class="table__muted">{{ t('rbac.unassigned') }}</span>
          </td>
          <td class="table__col-action">
            <button class="table__btn" type="button" @click="openUserEdit(user)">
              {{ t('common.edit') }}
            </button>
            <button class="table__btn" type="button" @click="removeUser(user)">
              {{ t('common.delete') }}
            </button>
          </td>
        </tr>
      </tbody>
    </table>

    <Teleport to="body">
      <div v-if="showingRoleForm" class="dialog-overlay" @click.self="showingRoleForm = false">
        <div class="dialog">
          <h2 class="dialog__title">
            {{ roleFormMode === 'create' ? t('rbac.roleDialogCreate') : t('rbac.roleDialogEdit') }}
          </h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">{{ t('common.name') }}</label>
              <input
                v-model="roleForm.name"
                class="dialog__input"
                type="text"
                :placeholder="t('rbac.namePlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('common.description') }}</label>
              <input
                v-model="roleForm.description"
                class="dialog__input"
                type="text"
                :placeholder="t('rbac.descriptionPlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('rbac.serverPermLabel') }}</label>
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
                <span v-if="servers.length === 0" class="table__muted">{{ t('rbac.noServers') }}</span>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingRoleForm = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!roleForm.name.trim()"
              @click="confirmRole"
            >
              {{ t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showingAssign" class="dialog-overlay" @click.self="showingAssign = false">
        <div class="dialog dialog--wide">
          <h2 class="dialog__title">{{ t('rbac.assignDialogTitle', { name: assigningRole?.name }) }}</h2>
          <div class="shuttle">
            <div class="shuttle__panel">
              <div class="shuttle__head">
                <span class="shuttle__label">{{ t('rbac.availableUsers') }}</span>
                <span class="shuttle__count">{{ leftUsers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="assignSearchLeft"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('rbac.filterAccount')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllLeftChecked"
                    @change="toggleSelectAll(leftUsers)"
                  />
                  <span>{{ t('rbac.selectAll') }}</span>
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
                  {{ t('rbac.noAvailableUsers') }}
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
                <span class="shuttle__label">{{ t('rbac.selectedUsers') }}</span>
                <span class="shuttle__count">{{ rightUsers.length }}</span>
              </div>
              <div class="shuttle__search">
                <input
                  v-model="assignSearchRight"
                  type="text"
                  class="shuttle__search-input"
                  :placeholder="t('rbac.filterAccount')"
                />
              </div>
              <div class="shuttle__list">
                <label class="shuttle__item shuttle__item--all">
                  <input
                    type="checkbox"
                    :checked="isAllRightChecked"
                    @change="toggleSelectAll(rightUsers)"
                  />
                  <span>{{ t('rbac.selectAll') }}</span>
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
                  {{ t('rbac.noneSelected') }}
                </span>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingAssign = false">
              {{ t('common.cancel') }}
            </button>
            <button class="btn btn--primary" type="button" @click="confirmAssign">
              {{ t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showingUserForm" class="dialog-overlay" @click.self="showingUserForm = false">
        <div class="dialog">
          <h2 class="dialog__title">{{ t('rbac.userDialogEdit') }}</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">{{ t('rbac.username') }}</label>
              <input
                v-model="userForm.username"
                class="dialog__input"
                type="text"
                :placeholder="t('rbac.usernamePlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('rbac.newPasswordHint') }}</label>
              <input
                v-model="userForm.password"
                class="dialog__input"
                type="password"
                :placeholder="t('rbac.passwordPlaceholder')"
              />
            </div>
            <div class="dialog__group">
              <label class="dialog__label">{{ t('rbac.role') }}</label>
              <select v-model="userForm.role_id" class="dialog__input">
                <option :value="0">{{ t('rbac.noRole') }}</option>
                <option v-for="role in roles" :key="role.id" :value="role.id">
                  {{ role.name }}
                </option>
              </select>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingUserForm = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="!userForm.username.trim()"
              @click="confirmUser"
            >
              {{ t('common.save') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <Teleport to="body">
      <div v-if="showingBatchCreate" class="dialog-overlay" @click.self="showingBatchCreate = false">
        <div class="dialog dialog--wide">
          <h2 class="dialog__title">{{ t('rbac.batchCreateTitle') }}</h2>
          <div class="dialog__form">
            <div class="dialog__group">
              <label class="dialog__label">
                {{ t('rbac.batchRole') }}
                <span class="dialog__label-hint">{{ t('rbac.batchRoleHint') }}</span>
              </label>
              <select v-model="batchRoleId" class="dialog__input">
                <option :value="0">{{ t('rbac.noRole') }}</option>
                <option v-for="role in roles" :key="role.id" :value="role.id">
                  {{ role.name }}
                </option>
              </select>
            </div>
            <div class="batch__grid">
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.username') }}</label>
                <input
                  v-model="entryForm.username"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.usernamePlaceholder')"
                />
              </div>
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.password') }}</label>
                <input
                  v-model="entryForm.password"
                  class="dialog__input"
                  type="password"
                  :placeholder="t('rbac.passwordPlaceholder')"
                />
              </div>
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.fullName') }}</label>
                <input
                  v-model="entryForm.name"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.namePlaceholder')"
                />
              </div>
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.email') }}</label>
                <input
                  v-model="entryForm.email"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.emailPlaceholder')"
                />
              </div>
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.uidLabel') }}</label>
                <input
                  v-model="entryForm.uid"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.uidPlaceholder')"
                />
              </div>
              <div class="dialog__group">
                <label class="dialog__label">{{ t('rbac.phone') }}</label>
                <input
                  v-model="entryForm.phone"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.phonePlaceholder')"
                />
              </div>
              <div class="dialog__group batch__grid-full">
                <label class="dialog__label">{{ t('rbac.organization') }}</label>
                <input
                  v-model="entryForm.organization"
                  class="dialog__input"
                  type="text"
                  :placeholder="t('rbac.organizationPlaceholder')"
                />
              </div>
            </div>
            <div class="batch__add-row">
              <button
                class="btn btn--primary"
                type="button"
                @click="addToDraft"
              >
                {{ t('rbac.addToList') }}
              </button>
            </div>
            <div class="batch__pending">
              <div class="batch__pending-head">
                <span class="batch__pending-title">{{ t('rbac.pendingList', { count: batchDraft.length }) }}</span>
                <button
                  v-if="batchDraft.length"
                  class="batch__clear"
                  type="button"
                  @click="clearDraft"
                >
                  {{ t('rbac.clearList') }}
                </button>
              </div>
              <div class="batch__list">
                <div v-if="batchDraft.length === 0" class="batch__empty">{{ t('rbac.pendingEmpty') }}</div>
                <div v-for="(item, index) in batchDraft" :key="item.username" class="batch__item">
                  <span class="batch__item-idx">{{ index + 1 }}</span>
                  <span class="batch__item-main">
                    <span class="batch__item-username">{{ item.username }}</span>
                    <span class="batch__item-sub">{{ item.name || item.email || item.uid || t('common.emptyDash') }}</span>
                  </span>
                  <button class="table__btn" type="button" @click="removeDraftItem(index)">
                    {{ t('common.delete') }}
                  </button>
                </div>
              </div>
            </div>
          </div>
          <div class="dialog__actions">
            <button class="btn btn--ghost" type="button" @click="showingBatchCreate = false">
              {{ t('common.cancel') }}
            </button>
            <button
              class="btn btn--primary"
              type="button"
              :disabled="batchDraft.length === 0 || batchSubmitting"
              @click="submitBatchCreate"
            >
              {{ batchSubmitting ? t('common.saving') : t('rbac.submitCreate') }}
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

.rbac__search {
  display: flex;
  align-items: center;
  gap: 8px;
}

.rbac__search-input {
  max-width: 280px;
  padding: 8px 14px;
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s;
}

.rbac__search-input::placeholder {
  color: var(--text-muted);
}

.rbac__search-input:focus {
  border-color: var(--xauat-blue);
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

.dialog__label-hint {
  margin-left: 6px;
  font-weight: 400;
  color: var(--text-muted);
}

.batch__grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 14px 16px;
}

.batch__grid-full {
  grid-column: 1 / -1;
}

.batch__add-row {
  display: flex;
  justify-content: flex-end;
}

.batch__pending {
  border: 1px solid var(--border);
  border-radius: 10px;
  overflow: hidden;
}

.batch__pending-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background: var(--bg);
  border-bottom: 1px solid var(--border);
}

.batch__pending-title {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}

.batch__clear {
  padding: 2px 10px;
  font-size: 12px;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 6px;
  cursor: pointer;
}

.batch__clear:hover {
  color: var(--xauat-blue);
  border-color: rgba(10, 61, 122, 0.25);
}

.batch__list {
  max-height: 220px;
  overflow-y: auto;
}

.batch__empty {
  padding: 20px;
  text-align: center;
  font-size: 13px;
  color: var(--text-muted);
}

.batch__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
}

.batch__item:last-child {
  border-bottom: none;
}

.batch__item-idx {
  width: 22px;
  height: 22px;
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.1);
  border-radius: 999px;
}

.batch__item-main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

.batch__item-username {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.batch__item-sub {
  font-size: 12px;
  color: var(--text-muted);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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

.shuttle__search {
  padding: 8px 12px;
  border-bottom: 1px solid var(--border);
  flex-shrink: 0;
}

.shuttle__search-input {
  width: 100%;
  padding: 6px 10px;
  font-size: 13px;
  border: 1px solid var(--border);
  border-radius: 6px;
  outline: none;
  background: var(--surface);
  color: var(--text);
  box-sizing: border-box;
}

.shuttle__search-input::placeholder {
  color: var(--text-muted);
}

.shuttle__search-input:focus {
  border-color: var(--xauat-blue);
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

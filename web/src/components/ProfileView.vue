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
import { ref, onMounted, computed } from 'vue'
import { auth, fetchProfile } from '../stores/auth.js'
import { updateProfile } from '../api.js'

const loading = ref(true)
const profile = ref(null)

const EMAIL_RE = /^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$/

onMounted(async () => {
  const data = await fetchProfile()
  profile.value = data
  loading.value = false
})

const isAdmin = computed(() => auth.user?.role === 'admin')

const adminFields = computed(() => {
  if (!profile.value) return []
  return [
    { label: '账号', value: profile.value.username },
    { label: '角色', value: '管理员' }
  ]
})

const userFields = computed(() => {
  if (!profile.value) return []
  return [
    { label: '学号/工号', value: profile.value.uid },
    { label: '姓名', value: profile.value.name },
    { label: '手机', value: profile.value.phone },
    { label: '所属部门/学院', value: profile.value.organization }
  ]
})

const email = computed(() => profile.value?.email || '')
const editingEmail = ref(false)
const emailDraft = ref('')
const savingEmail = ref(false)
const emailError = ref('')

function startEditEmail() {
  emailDraft.value = email.value
  emailError.value = ''
  editingEmail.value = true
}

function cancelEditEmail() {
  editingEmail.value = false
  emailError.value = ''
}

async function saveEmail() {
  const value = emailDraft.value.trim()
  if (!value) {
    emailError.value = '邮箱不能为空'
    return
  }
  if (!EMAIL_RE.test(value)) {
    emailError.value = '邮箱格式不正确'
    return
  }
  if (value === email.value) {
    editingEmail.value = false
    return
  }
  savingEmail.value = true
  emailError.value = ''
  try {
    profile.value = await updateProfile({ email: value })
    editingEmail.value = false
  } catch (e) {
    emailError.value = e.message
  } finally {
    savingEmail.value = false
  }
}
</script>

<template>
  <section class="page">
    <h1 class="page__title">个人信息</h1>
    <dl v-if="loading" class="profile profile--loading">
      <div class="profile__row">加载中…</div>
    </dl>
    <dl v-else-if="!profile" class="profile profile--error">
      <div class="profile__row">获取信息失败</div>
    </dl>
    <dl v-else class="profile">
      <template v-if="isAdmin">
        <div v-for="field in adminFields" :key="field.label" class="profile__row">
          <dt class="profile__label">{{ field.label }}</dt>
          <dd class="profile__value">{{ field.value || '—' }}</dd>
        </div>
      </template>
      <template v-else>
        <div v-for="field in userFields" :key="field.label" class="profile__row">
          <dt class="profile__label">{{ field.label }}</dt>
          <dd class="profile__value">{{ field.value || '—' }}</dd>
        </div>

        <div class="profile__row">
          <dt class="profile__label">邮箱</dt>
          <dd v-if="!editingEmail" class="profile__value profile__email">
            <span>{{ profile.email || '—' }}</span>
            <button class="profile__btn" type="button" @click="startEditEmail">修改</button>
          </dd>
          <dd v-else class="profile__value profile__email">
            <input
              v-model="emailDraft"
              class="profile__input"
              type="text"
              placeholder="请输入邮箱"
              @keyup.enter="saveEmail"
              @keyup.esc="cancelEditEmail"
            />
            <button
              class="profile__btn profile__btn--primary"
              type="button"
              :disabled="savingEmail"
              @click="saveEmail"
            >
              {{ savingEmail ? '保存中…' : '保存' }}
            </button>
            <button class="profile__btn" type="button" :disabled="savingEmail" @click="cancelEditEmail">取消</button>
            <p v-if="emailError" class="profile__error">{{ emailError }}</p>
          </dd>
        </div>
      </template>
    </dl>
  </section>
</template>

<style scoped>
.profile {
  margin-top: 20px;
  max-width: 560px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  overflow: hidden;
}

.profile--loading .profile__row,
.profile--error .profile__row {
  padding: 24px 20px;
  color: var(--text-muted);
  font-size: 14px;
}

.profile__row {
  display: flex;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid var(--border);
}

.profile__row:last-child {
  border-bottom: none;
}

.profile__label {
  flex: 0 0 140px;
  font-size: 14px;
  color: var(--text-muted);
}

.profile__value {
  margin: 0;
  font-size: 15px;
  color: var(--text);
}

.profile__email {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.profile__btn {
  padding: 4px 12px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 6px;
  cursor: pointer;
}

.profile__btn:hover:not(:disabled) {
  border-color: var(--xauat-blue);
}

.profile__btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.profile__btn--primary {
  color: #fff;
  background: var(--xauat-blue);
  border-color: var(--xauat-blue);
}

.profile__btn--primary:hover:not(:disabled) {
  background: var(--xauat-blue-light);
  border-color: var(--xauat-blue-light);
}

.profile__input {
  width: 220px;
  padding: 6px 10px;
  font-size: 14px;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 6px;
  outline: none;
}

.profile__input:focus {
  border-color: var(--xauat-blue);
}

.profile__error {
  flex-basis: 100%;
  margin: 0;
  font-size: 12px;
  color: #dc2626;
}
</style>

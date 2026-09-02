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
import { useI18n } from 'vue-i18n'
import { auth, fetchProfile } from '../stores/auth.js'
import { roleLabel } from '../i18n.js'
import { updateProfile } from '../api.js'

const { t } = useI18n()

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
    { labelKey: 'profile.account', value: profile.value.username },
    { labelKey: 'profile.roleField', value: roleLabel(auth.user?.role) }
  ]
})

const userFields = computed(() => {
  if (!profile.value) return []
  return [
    { labelKey: 'profile.uid', value: profile.value.uid },
    { labelKey: 'profile.name', value: profile.value.name },
    { labelKey: 'profile.phone', value: profile.value.phone },
    { labelKey: 'profile.organization', value: profile.value.organization }
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
    emailError.value = t('profile.emailEmpty')
    return
  }
  if (!EMAIL_RE.test(value)) {
    emailError.value = t('profile.emailInvalid')
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
    <h1 class="page__title">{{ t('profile.title') }}</h1>
    <dl v-if="loading" class="profile profile--loading">
      <div class="profile__row">{{ t('common.loading') }}</div>
    </dl>
    <dl v-else-if="!profile" class="profile profile--error">
      <div class="profile__row">{{ t('profile.loadError') }}</div>
    </dl>
    <dl v-else class="profile">
      <template v-if="isAdmin">
        <div v-for="field in adminFields" :key="field.labelKey" class="profile__row">
          <dt class="profile__label">{{ t(field.labelKey) }}</dt>
          <dd class="profile__value">{{ field.value || t('common.emptyDash') }}</dd>
        </div>
      </template>
      <template v-else>
        <div v-for="field in userFields" :key="field.labelKey" class="profile__row">
          <dt class="profile__label">{{ t(field.labelKey) }}</dt>
          <dd class="profile__value">{{ field.value || t('common.emptyDash') }}</dd>
        </div>

        <div class="profile__row">
          <dt class="profile__label">{{ t('profile.email') }}</dt>
          <dd v-if="!editingEmail" class="profile__value profile__email">
            <span>{{ profile.email || t('common.emptyDash') }}</span>
            <button class="profile__btn" type="button" @click="startEditEmail">{{ t('profile.change') }}</button>
          </dd>
          <dd v-else class="profile__value profile__email">
            <input
              v-model="emailDraft"
              class="profile__input"
              type="text"
              :placeholder="t('profile.emailPlaceholder')"
              @keyup.enter="saveEmail"
              @keyup.esc="cancelEditEmail"
            />
            <button
              class="profile__btn profile__btn--primary"
              type="button"
              :disabled="savingEmail"
              @click="saveEmail"
            >
              {{ savingEmail ? t('common.saving') : t('common.save') }}
            </button>
            <button class="profile__btn" type="button" :disabled="savingEmail" @click="cancelEditEmail">{{ t('common.cancel') }}</button>
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

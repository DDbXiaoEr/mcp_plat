<script setup>

// Author: deepseek-v4-pro / opencode
import { ref, onMounted, computed } from 'vue'
import { auth, fetchProfile } from '../stores/auth.js'

const loading = ref(true)
const profile = ref(null)

onMounted(async () => {
  const data = await fetchProfile()
  profile.value = data
  loading.value = false
})

const isAdmin = computed(() => auth.user?.role === 'admin')

const fields = computed(() => {
  if (!profile.value) return []
  if (isAdmin.value) {
    return [
      { label: '账号', value: profile.value.username },
      { label: '角色', value: '管理员' }
    ]
  }
  return [
    { label: '学号/工号', value: profile.value.uid },
    { label: '用户名', value: profile.value.username },
    { label: '邮箱', value: profile.value.email },
    { label: '手机', value: profile.value.phone },
    { label: '所属部门/学院', value: profile.value.organization }
  ]
})
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
      <div v-for="field in fields" :key="field.label" class="profile__row">
        <dt class="profile__label">{{ field.label }}</dt>
        <dd class="profile__value">{{ field.value || '—' }}</dd>
      </div>
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
</style>

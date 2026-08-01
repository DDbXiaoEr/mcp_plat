<script setup>

// Author: deepseek-v4-pro / opencode
import { onMounted } from 'vue'
import { auth, logout } from '../stores/auth.js'
import { platformSettings, loadSettings } from '../stores/settings.js'

onMounted(() => loadSettings())
</script>

<template>
  <header class="header">
    <div class="header__brand">
      <a
        class="brand__link"
        :href="platformSettings.platform.siteUrl || undefined"
        target="_blank"
        rel="noopener noreferrer"
      >
        <span class="brand__mark" aria-hidden="true">
          <img v-if="platformSettings.platform.logoUrl" :src="platformSettings.platform.logoUrl" alt="" class="brand__logo" />
        </span>
        <span class="brand__text">
          <strong>{{ platformSettings.platform.name }}</strong>
          <em>MCP 服务平台 · 控制台</em>
        </span>
      </a>
    </div>

    <div class="header__actions">
      <span v-if="auth.user" class="header__user">
        <span class="header__role">{{ auth.user.roleLabel }}</span>
        <span class="header__name">{{ auth.user.name || auth.user.username }}</span>
      </span>
      <button class="header__logout" type="button" @click="logout">
        退出登录
      </button>
    </div>
  </header>
</template>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: var(--header-height);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
  background: var(--surface);
  border-bottom: 1px solid var(--border);
}

.header__brand {
  display: flex;
  align-items: center;
}

.brand__link {
  display: flex;
  align-items: center;
  gap: 12px;
}

.brand__mark {
  display: grid;
  place-items: center;
  flex-shrink: 0;
}

.brand__logo {
  height: 36px;
  width: auto;
  object-fit: contain;
}

.brand__text {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.brand__text strong {
  font-size: 15px;
  color: var(--xauat-blue);
}

.brand__text em {
  font-style: normal;
  font-size: 12px;
  color: var(--text-muted);
}

.header__actions {
  display: flex;
  align-items: center;
  gap: 16px;
}

.header__user {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
}

.header__role {
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  color: var(--xauat-blue);
  background: rgba(30, 95, 176, 0.1);
}

.header__name {
  color: var(--text);
  font-weight: 500;
}

.header__logout {
  padding: 7px 16px;
  font-size: 14px;
  color: var(--xauat-blue);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.header__logout:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}
</style>

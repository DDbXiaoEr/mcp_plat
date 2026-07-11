<script setup>
import { auth, logout } from '../stores/auth.js'

// TODO: 后期可配置的 logo 地址，留空则显示占位块
const logoUrl = ''
</script>

<template>
  <header class="header">
    <div class="header__brand">
      <span class="brand__mark" aria-hidden="true">
        <img v-if="logoUrl" :src="logoUrl" alt="" class="brand__logo" />
      </span>
      <span class="brand__text">
        <strong>某某大学</strong>
        <em>MCP 服务平台 · 控制台</em>
      </span>
    </div>

    <div class="header__actions">
      <span v-if="auth.user" class="header__user">
        <span class="header__role">{{ auth.user.roleLabel }}</span>
        <span class="header__name">{{ auth.user.username }}</span>
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
  gap: 12px;
}

.brand__mark {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border-radius: 12px;
  background: linear-gradient(135deg, var(--xauat-blue), var(--xauat-cyan));
  color: #fff;
  font-weight: 800;
  overflow: hidden;
}

.brand__logo {
  width: 100%;
  height: 100%;
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

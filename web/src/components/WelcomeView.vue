<script setup>
import { auth } from '../stores/auth.js'
</script>

<template>
  <!-- TODO: 后期按角色拆分为不同的功能页面 -->
  <section v-if="auth.user" class="welcome">
    <h1 class="welcome__title">欢迎，{{ auth.user.roleLabel }}登录</h1>
  </section>
</template>

<style scoped>
.welcome {
  display: grid;
  place-items: center;
  min-height: calc(100vh - var(--header-height) - 48px);
}

.welcome__title {
  font-size: 28px;
  font-weight: 700;
  color: var(--xauat-blue);
}
</style>

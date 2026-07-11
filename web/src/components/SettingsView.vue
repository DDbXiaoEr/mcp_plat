<script setup></script>

<template>
  <section class="settings">
    <h1 class="settings__title">系统设置</h1>
    <p class="settings__hint">暂无内容</p>
  </section>
</template>

<style scoped>
.settings {
  display: grid;
  place-items: center;
  gap: 8px;
  min-height: calc(100vh - var(--header-height) - 48px);
}

.settings__title {
  font-size: 24px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.settings__hint {
  font-size: 14px;
  color: var(--text-muted);
}
</style>

<script setup>
// TODO: 接入后端后替换为真实统计数据
const stats = [
  { key: 'users', label: '平台用户数', value: '1,286' },
  { key: 'systems', label: '接入业务系统数（MCP 服务器数）', value: '24' },
  { key: 'tools', label: 'MCP 工具数', value: '68' },
  { key: 'calls', label: '当天 AI 调用次数', value: '3,472' }
]
</script>

<template>
  <section class="overview">
    <h1 class="page__title">平台概况</h1>
    <div class="overview__stats">
      <div v-for="item in stats" :key="item.key" class="stat-card">
        <span class="stat-card__label">{{ item.label }}</span>
        <span class="stat-card__value">{{ item.value }}</span>
      </div>
    </div>
  </section>
</template>

<style scoped>
.overview {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.overview__stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.stat-card {
  display: flex;
  flex-direction: column;
  gap: 6px;
  aspect-ratio: 2 / 1;
  align-items: center;
  justify-content: center;
  text-align: center;
  padding: 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.stat-card__value {
  font-size: 44px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.stat-card__label {
  font-size: 17px;
  font-weight: 600;
  color: var(--text);
}

@media (max-width: 900px) {
  .overview__stats {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>

<script setup>

// Author: deepseek-v4-pro / opencode
// TODO: 接入后端后替换为真实统计数据
const stats = [
  { key: 'users', label: '平台用户数', value: '1,286' },
  { key: 'systems', label: '接入业务系统数（MCP 服务器数）', value: '24' },
  { key: 'tools', label: 'MCP 工具数', value: '68' },
  { key: 'calls', label: '当天 AI 调用次数', value: '3,472' }
]

const version = __APP_VERSION__
const gitCommit = __GIT_COMMIT__
const buildTime = __BUILD_TIME__
</script>

<template>
  <section class="overview">
    <div class="overview__head">
      <h1 class="page__title">平台概况</h1>
      <span class="overview__version" :title="`commit ${gitCommit}, built at ${buildTime}`">{{ version }}</span>
    </div>
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

.overview__head {
  display: flex;
  align-items: baseline;
  gap: 12px;
}

.overview__version {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
}
</style>

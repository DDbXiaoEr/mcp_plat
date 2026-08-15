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
import { ref, watch, onMounted, onBeforeUnmount, nextTick } from 'vue'
import * as echarts from 'echarts/core'
import { LineChart, BarChart, PieChart } from 'echarts/charts'
import { GridComponent, TooltipComponent, LegendComponent } from 'echarts/components'
import { CanvasRenderer } from 'echarts/renderers'
import { fetchOverviewStats, fetchOverviewCallTrend } from '../api.js'

echarts.use([LineChart, BarChart, PieChart, GridComponent, TooltipComponent, LegendComponent, CanvasRenderer])

// Author: deepseek-v4-pro / opencode
const stats = ref([])
const statsLoading = ref(true)

const STATS_DEF = [
  { key: 'users', label: '平台用户数' },
  { key: 'servers', label: '接入业务系统数（MCP 服务器数）' },
  { key: 'tools', label: 'MCP 工具数' },
  { key: 'today_calls', label: '当天 AI 调用次数' }
]

function formatCount(n) {
  return (n ?? 0).toLocaleString()
}

async function loadStats(silent = false) {
  if (!silent) statsLoading.value = true
  try {
    stats.value = await fetchOverviewStats()
  } catch {
    if (!silent) stats.value = {}
  }
  if (!silent) statsLoading.value = false
}

const AUTO_REFRESH_OPTIONS = [
  { value: 0, label: '自动刷新：关闭' },
  { value: 10000, label: '自动刷新：10 秒' },
  { value: 30000, label: '自动刷新：30 秒' },
  { value: 60000, label: '自动刷新：1 分钟' },
  { value: 300000, label: '自动刷新：5 分钟' }
]
const autoRefresh = ref(0)
const refreshing = ref(false)
const lastUpdated = ref(null)
let refreshTimer = null

function formatTime(d) {
  return d.toLocaleTimeString('zh-CN', { hour12: false })
}

async function refreshAll() {
  if (refreshing.value) return
  refreshing.value = true
  try {
    await Promise.all([loadStats(true), loadTrend()])
    lastUpdated.value = new Date()
  } catch {
    // 刷新失败时保留已有数据
  }
  refreshing.value = false
}

const version = __APP_VERSION__
const gitCommit = __GIT_COMMIT__
const buildTime = __BUILD_TIME__

const periods = [
  { key: '7', label: '近7天' },
  { key: '30', label: '近30天' }
]
const period = ref('30')
const trend = ref({ dates: [], counts: [], server_calls: [], user_groups: [], server_names: [], server_counts_by_day: [] })
const chartEl = ref(null)
let chart = null

const serverPieEl = ref(null)
const groupPieEl = ref(null)
const serverPieChart = ref(null)
const groupPieChart = ref(null)

function cssVar(name, fallback) {
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return value || fallback
}

const SERVER_BAR_PALETTE = ['#0a3d7a', '#4a90d9', '#7fb4ea', '#35a08f', '#e09b3d', '#c25e5e', '#9b6bb8', '#5b6b80', '#8a97a6']

function renderChart() {
  if (!chartEl.value) return
  if (!chart) {
    chart = echarts.init(chartEl.value)
  }
  if (!trend.value.dates.length) {
    chart.clear()
    return
  }
  const serverNames = trend.value.server_names || []
  const byDay = trend.value.server_counts_by_day || []
  const series = []
  serverNames.forEach((name, idx) => {
    const isLast = idx === serverNames.length - 1
    series.push({
      name,
      type: 'bar',
      stack: 'calls',
      barMaxWidth: 28,
      itemStyle: {
        color: SERVER_BAR_PALETTE[idx % SERVER_BAR_PALETTE.length],
        borderRadius: isLast ? [3, 3, 0, 0] : 0,
        borderColor: '#fff',
        borderWidth: 1
      },
      data: byDay.map((row) => row[idx] || 0)
    })
  })
  series.push({
    name: 'AI 调用次数',
    type: 'line',
    symbol: 'circle',
    symbolSize: 6,
    data: trend.value.counts,
    lineStyle: { width: 3, color: cssVar('--xauat-blue', '#0a3d7a') },
    itemStyle: { color: cssVar('--xauat-blue', '#0a3d7a') },
    areaStyle: {
      color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
        { offset: 0, color: 'rgba(10, 61, 122, 0.28)' },
        { offset: 1, color: 'rgba(10, 61, 122, 0.02)' }
      ])
    }
  })
  chart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    legend: {
      top: 0,
      icon: 'circle',
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 12,
      textStyle: { color: cssVar('--text-muted', '#5b6b80'), fontSize: 12 }
    },
    grid: { left: 12, right: 16, top: 36, bottom: 8, containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: true,
      data: trend.value.dates,
      axisLine: { lineStyle: { color: cssVar('--border', '#e3e9f2') } },
      axisLabel: { color: cssVar('--text-muted', '#5b6b80'), fontSize: 12 }
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: '#eef2f8' } },
      axisLabel: { color: cssVar('--text-muted', '#5b6b80'), fontSize: 12 }
    },
    series
  })
}

function renderPie(el, holder, data) {
  if (!el) return
  if (!holder.value) holder.value = echarts.init(el)
  if (!data.length) {
    holder.value.clear()
    return
  }
  holder.value.setOption({
    tooltip: { trigger: 'item', formatter: '{b}: {c}（{d}%）' },
    legend: {
      bottom: 0,
      icon: 'circle',
      itemWidth: 10,
      itemHeight: 10,
      itemGap: 12,
      textStyle: { color: cssVar('--text-muted', '#5b6b80'), fontSize: 12 }
    },
    series: [
      {
        type: 'pie',
        radius: ['42%', '68%'],
        center: ['50%', '44%'],
        avoidLabelOverlap: true,
        itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
        label: { color: cssVar('--text', '#1f2d3d'), fontSize: 12, formatter: '{b} {d}%' },
        labelLine: { length: 12, length2: 8 },
        data: data.map((d) => ({ name: d.name, value: d.count }))
      }
    ]
  })
}

function renderServerPie() {
  renderPie(serverPieEl.value, serverPieChart, trend.value.server_calls || [])
}

function renderGroupPie() {
  renderPie(groupPieEl.value, groupPieChart, trend.value.user_groups || [])
}

async function loadTrend() {
  const data = await fetchOverviewCallTrend(Number(period.value))
  trend.value = data
  await nextTick()
  renderChart()
  renderServerPie()
  renderGroupPie()
}

async function refreshTrend() {
  try {
    await loadTrend()
  } catch {
    trend.value = { dates: [], counts: [], server_calls: [], user_groups: [], server_names: [], server_counts_by_day: [] }
    await nextTick()
    renderChart()
    renderServerPie()
    renderGroupPie()
  }
}

function handleResize() {
  chart && chart.resize()
  serverPieChart.value && serverPieChart.value.resize()
  groupPieChart.value && groupPieChart.value.resize()
}

watch(period, refreshTrend)

watch(autoRefresh, (ms) => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
  if (ms > 0) {
    refreshTimer = setInterval(refreshAll, ms)
  }
})

onMounted(() => {
  loadStats()
  refreshTrend()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  if (refreshTimer) {
    clearInterval(refreshTimer)
    refreshTimer = null
  }
  window.removeEventListener('resize', handleResize)
  chart && chart.dispose()
  chart = null
  serverPieChart.value && serverPieChart.value.dispose()
  serverPieChart.value = null
  groupPieChart.value && groupPieChart.value.dispose()
  groupPieChart.value = null
})
</script>

<template>
  <section class="overview">
    <div class="overview__head">
      <h1 class="page__title">平台概况</h1>
      <span class="overview__version" :title="`commit ${gitCommit}, built at ${buildTime}`">{{ version }}</span>
      <div class="overview__controls">
        <select v-model="autoRefresh" class="overview__select" title="自动刷新间隔">
          <option v-for="opt in AUTO_REFRESH_OPTIONS" :key="opt.value" :value="opt.value">{{ opt.label }}</option>
        </select>
        <button type="button" class="overview__refresh" :disabled="refreshing" @click="refreshAll">
          <span class="overview__refresh-icon" :class="{ 'is-spinning': refreshing }">↻</span>
          <span>{{ refreshing ? '刷新中' : '刷新' }}</span>
        </button>
        <span v-if="lastUpdated" class="overview__updated">更新于 {{ formatTime(lastUpdated) }}</span>
      </div>
    </div>
    <div class="overview__stats">
      <div v-for="item in STATS_DEF" :key="item.key" class="stat-card">
        <span class="stat-card__label">{{ item.label }}</span>
        <span class="stat-card__value">{{ statsLoading ? '—' : formatCount(stats[item.key]) }}</span>
      </div>
    </div>
    <div class="trend">
      <div class="trend__head">
        <h2 class="trend__title">AI 调用历史趋势</h2>
        <div class="trend__tabs">
          <button
            v-for="p in periods"
            :key="p.key"
            type="button"
            class="trend__tab"
            :class="{ 'is-active': period === p.key }"
            @click="period = p.key"
          >
            {{ p.label }}
          </button>
        </div>
      </div>
      <div class="trend__body">
        <div ref="chartEl" class="trend__canvas"></div>
        <div v-show="!trend.dates.length" class="trend__empty">暂无数据</div>
      </div>
      <div class="trend__pies">
        <div class="pie-card">
          <h3 class="pie-card__title">服务器调用量占比</h3>
          <div class="pie-card__body">
            <div ref="serverPieEl" class="pie-card__canvas"></div>
            <div v-show="!trend.server_calls?.length" class="pie-card__empty">暂无数据</div>
          </div>
        </div>
        <div class="pie-card">
          <h3 class="pie-card__title">用户组用户数占比</h3>
          <div class="pie-card__body">
            <div ref="groupPieEl" class="pie-card__canvas"></div>
            <div v-show="!trend.user_groups?.length" class="pie-card__empty">暂无数据</div>
          </div>
        </div>
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

.overview__controls {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 10px;
}

.overview__select {
  padding: 5px 10px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  outline: none;
}

.overview__refresh {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text);
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.overview__refresh:hover:not(:disabled) {
  border-color: rgba(10, 61, 122, 0.3);
}

.overview__refresh:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.overview__refresh-icon {
  display: inline-block;
}

.overview__refresh-icon.is-spinning {
  animation: overview-spin 0.8s linear infinite;
}

@keyframes overview-spin {
  to {
    transform: rotate(360deg);
  }
}

.overview__updated {
  font-size: 12px;
  color: var(--text-muted);
}

.trend {
  padding: 20px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.trend__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.trend__title {
  font-size: 17px;
  font-weight: 600;
  color: var(--text);
}

.trend__tabs {
  display: flex;
  gap: 8px;
}

.trend__tab {
  padding: 6px 16px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-muted);
  background: transparent;
  border: 1px solid var(--border);
  border-radius: 999px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s, border-color 0.2s;
}

.trend__tab:hover {
  border-color: rgba(10, 61, 122, 0.3);
}

.trend__tab.is-active {
  color: #fff;
  background: var(--xauat-blue);
  border-color: var(--xauat-blue);
}

.trend__body {
  position: relative;
  margin-top: 16px;
}

.trend__canvas {
  height: 300px;
}

.trend__empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--text-muted);
  background: var(--surface);
}

.trend__pies {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  margin-top: 16px;
}

@media (max-width: 900px) {
  .trend__pies {
    grid-template-columns: 1fr;
  }
}

.pie-card {
  padding: 16px;
  background: #f7fafd;
  border: 1px solid var(--border);
  border-radius: var(--radius);
}

.pie-card__title {
  font-size: 15px;
  font-weight: 600;
  color: var(--text);
}

.pie-card__body {
  position: relative;
  margin-top: 8px;
}

.pie-card__canvas {
  height: 280px;
}

.pie-card__empty {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  color: var(--text-muted);
  background: #f7fafd;
}
</style>

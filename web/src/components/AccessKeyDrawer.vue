<script setup>
import { reactive, ref, computed, watch } from 'vue'

const EXPIRATION_OPTIONS = [
  { label: '7 天', value: 7 },
  { label: '30 天', value: 30 },
  { label: '90 天', value: 90 },
  { label: '180 天', value: 180 },
  { label: '365 天', value: 365 },
  { label: '永不过期', value: -1 }
]

const props = defineProps({
  item: { type: Object, default: null },
  servers: { type: Array, default: () => [] }
})
const emit = defineEmits(['save', 'close'])

const form = reactive({
  name: '',
  enabled: true,
  expireDays: 30,
  tools: {}
})

const selectedServer = ref('')

const currentTools = computed(() => {
  const server = props.servers.find((s) => s.key === selectedServer.value)
  return server ? server.tools : []
})

function parseServers(raw) {
  try {
    return JSON.parse(raw) || {}
  } catch {
    return {}
  }
}

function cloneTools(tools) {
  return Object.fromEntries(
    Object.entries(tools || {}).map(([k, v]) => [k, [...v]])
  )
}

watch(
  () => props.item,
  (item) => {
    if (item) {
      form.name = item.name
      form.enabled = item.enabled
      form.tools = parseServers(item.servers)
      if (item.expired_at) {
        const now = new Date()
        const exp = new Date(item.expired_at)
        const diff = Math.ceil((exp.getTime() - now.getTime()) / 86400000)
        form.expireDays = diff > 0 ? diff : 30
      } else {
        form.expireDays = -1
      }
      selectedServer.value = props.servers[0]?.key || ''
    }
  },
  { immediate: true }
)

function isChecked(toolKey) {
  return form.tools[selectedServer.value]?.includes(toolKey)
}

function toggleTool(toolKey) {
  const list = form.tools[selectedServer.value] || (form.tools[selectedServer.value] = [])
  const i = list.indexOf(toolKey)
  if (i === -1) list.push(toolKey)
  else list.splice(i, 1)
}

function onSave() {
  const payload = {
    name: form.name.trim(),
    enabled: form.enabled,
    servers: JSON.stringify(cloneTools(form.tools))
  }
  if (form.expireDays !== -1) {
    const d = new Date()
    d.setDate(d.getDate() + form.expireDays)
    payload.expired_at = d.toISOString()
  } else {
    payload.expired_at = null
  }
  emit('save', payload)
}
</script>

<template>
  <div class="drawer" :class="{ 'drawer--open': item }">
    <div class="drawer__overlay" @click="emit('close')"></div>
    <aside class="drawer__panel">
      <header class="drawer__head">
        <h2 class="drawer__title">编辑 AccessKey</h2>
        <button class="drawer__close" type="button" aria-label="关闭" @click="emit('close')">
          ×
        </button>
      </header>

      <div v-if="item" class="drawer__body">
        <label class="field">
          <span class="field__label">名称</span>
          <input
            v-model="form.name"
            type="text"
            class="field__input"
            placeholder="请输入名称"
          />
        </label>

        <div class="field">
          <span class="field__label">可用状态</span>
          <label class="switch">
            <input type="checkbox" v-model="form.enabled" />
            <span class="switch__track"><span class="switch__thumb"></span></span>
            <span class="switch__label">{{ form.enabled ? '启用' : '禁用' }}</span>
          </label>
        </div>

        <div class="field">
          <span class="field__label">过期时间</span>
          <select v-model="form.expireDays" class="field__input">
            <option
              v-for="opt in EXPIRATION_OPTIONS"
              :key="opt.value"
              :value="opt.value"
            >
              {{ opt.label }}
            </option>
          </select>
        </div>

        <div class="field">
          <span class="field__label">工具使用权限</span>
          <select v-model="selectedServer" class="field__input">
            <option v-for="s in servers" :key="s.key" :value="s.key">
              {{ s.label }}
            </option>
          </select>
          <div class="tools">
            <label v-for="tool in currentTools" :key="tool.key" class="tool">
              <input
                type="checkbox"
                :checked="isChecked(tool.key)"
                @change="toggleTool(tool.key)"
              />
              <span>{{ tool.label }}</span>
            </label>
          </div>
        </div>
      </div>

      <footer class="drawer__foot">
        <button class="btn btn--ghost" type="button" @click="emit('close')">
          取消
        </button>
        <button class="btn btn--primary" type="button" @click="onSave">
          保存
        </button>
      </footer>
    </aside>
  </div>
</template>

<style scoped>
.drawer {
  position: fixed;
  inset: 0;
  z-index: 200;
  visibility: hidden;
  pointer-events: none;
}

.drawer--open {
  visibility: visible;
  pointer-events: auto;
}

.drawer__overlay {
  position: absolute;
  inset: 0;
  background: rgba(10, 61, 122, 0.28);
  opacity: 0;
  transition: opacity 0.25s ease;
}

.drawer--open .drawer__overlay {
  opacity: 1;
}

.drawer__panel {
  position: absolute;
  top: 0;
  right: 0;
  bottom: 0;
  width: 380px;
  max-width: 90vw;
  display: flex;
  flex-direction: column;
  background: var(--surface);
  box-shadow: -8px 0 30px rgba(10, 61, 122, 0.12);
  transform: translateX(100%);
  transition: transform 0.25s ease;
}

.drawer--open .drawer__panel {
  transform: translateX(0);
}

.drawer__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 18px 24px;
  border-bottom: 1px solid var(--border);
}

.drawer__title {
  font-size: 17px;
  font-weight: 700;
  color: var(--xauat-blue);
}

.drawer__close {
  width: 32px;
  height: 32px;
  font-size: 22px;
  line-height: 1;
  color: var(--text-muted);
  background: transparent;
  border: none;
  border-radius: 8px;
  cursor: pointer;
}

.drawer__close:hover {
  background: rgba(10, 61, 122, 0.06);
}

.drawer__body {
  flex: 1;
  overflow-y: auto;
  padding: 24px;
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.field__label {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
}

.field__input {
  padding: 10px 14px;
  font-size: 15px;
  color: var(--text);
  border: 1px solid var(--border);
  border-radius: 10px;
  outline: none;
  transition: border-color 0.2s, box-shadow 0.2s;
}

.field__input:focus {
  border-color: var(--xauat-blue-light);
  box-shadow: 0 0 0 3px rgba(30, 95, 176, 0.12);
}

.tools {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.tool {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 14px;
  color: var(--text);
  cursor: pointer;
}

.tool input {
  width: 16px;
  height: 16px;
  accent-color: var(--xauat-blue);
  cursor: pointer;
}

.switch {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.switch input {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
}

.switch__track {
  position: relative;
  width: 40px;
  height: 22px;
  border-radius: 999px;
  background: var(--border);
  transition: background 0.2s;
}

.switch__thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  width: 18px;
  height: 18px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
  transition: transform 0.2s;
}

.switch input:checked + .switch__track {
  background: var(--xauat-blue);
}

.switch input:checked + .switch__track .switch__thumb {
  transform: translateX(18px);
}

.switch__label {
  font-size: 13px;
  color: var(--text-muted);
}

.drawer__foot {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  padding: 16px 24px;
  border-top: 1px solid var(--border);
}

.btn {
  padding: 8px 20px;
  font-size: 14px;
  font-weight: 600;
  border-radius: 999px;
  border: 1px solid transparent;
  cursor: pointer;
  transition: background 0.2s, border-color 0.2s;
}

.btn--primary {
  color: #fff;
  background: var(--xauat-blue);
}

.btn--primary:hover {
  background: var(--xauat-blue-light);
}

.btn--ghost {
  color: var(--xauat-blue);
  background: transparent;
  border-color: var(--border);
}

.btn--ghost:hover {
  background: rgba(10, 61, 122, 0.06);
  border-color: rgba(10, 61, 122, 0.25);
}
</style>

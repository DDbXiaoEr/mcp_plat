// Author: deepseek-v4-pro / opencode
import { reactive, readonly } from 'vue'
import { fetchSettings } from '../api.js'

const state = reactive({
  platform: {
    name: '某某大学',
    logoUrl: '',
    siteUrl: ''
  },
  loaded: false
})

async function loadSettings() {
  if (state.loaded) return
  try {
    const data = await fetchSettings()
    if (data && data.platform) {
      if (data.platform.name) state.platform.name = data.platform.name
      if (data.platform.logoUrl != null) state.platform.logoUrl = data.platform.logoUrl
      if (data.platform.siteUrl != null) state.platform.siteUrl = data.platform.siteUrl
    }
  } catch (_) { /* ignore */ }
  state.loaded = true
}

function applyPlatform(data) {
  if (!data) return
  state.platform.name = data.name || '某某大学'
  state.platform.logoUrl = data.logoUrl || ''
  state.platform.siteUrl = data.siteUrl || ''
}

export const platformSettings = readonly(state)
export { loadSettings, applyPlatform }

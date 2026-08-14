/*
 * Copyright (C) 2026 Zhaoquan Wang
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

// Author: deepseek-v4-pro / opencode
import { reactive, readonly } from 'vue'
import { fetchSetting } from '../api.js'

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
    const data = await fetchSetting('platform')
    applyPlatform(data)
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

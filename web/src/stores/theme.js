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

const STORAGE_KEY = 'mcp-console-theme'

function systemTheme() {
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
}

function loadTheme() {
  const saved = localStorage.getItem(STORAGE_KEY)
  return saved === 'light' || saved === 'dark' ? saved : systemTheme()
}

const state = reactive({
  theme: loadTheme()
})

function apply(theme) {
  document.documentElement.dataset.theme = theme
}

apply(state.theme)

function toggleTheme() {
  state.theme = state.theme === 'dark' ? 'light' : 'dark'
  localStorage.setItem(STORAGE_KEY, state.theme)
  apply(state.theme)
}

export const theme = readonly(state)
export { toggleTheme }

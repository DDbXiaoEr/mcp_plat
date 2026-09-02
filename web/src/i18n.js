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
import { createI18n } from 'vue-i18n'
import zhCN from './locales/zh-CN/index.js'
import enUS from './locales/en-US/index.js'

export const SUPPORTED_LOCALES = ['zh-CN', 'en-US']

const STORAGE_KEY = 'mcp-console-locale'

function detectLocale() {
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (SUPPORTED_LOCALES.includes(saved)) return saved
  } catch (_) { /* ignore */ }
  const lang = (navigator.language || navigator.languages?.[0] || '').toLowerCase()
  return lang.startsWith('zh') ? 'zh-CN' : 'en-US'
}

export const initialLocale = detectLocale()

export const i18n = createI18n({
  legacy: false,
  globalInjection: true,
  locale: initialLocale,
  fallbackLocale: 'zh-CN',
  messages: {
    'zh-CN': zhCN,
    'en-US': enUS
  },
  missingWarn: false,
  fallbackWarn: false
})

// roleLabel resolves a role code to its localized label at render time.
export function roleLabel(role) {
  if (role === 'admin') return i18n.global.t('role.admin')
  if (role === 'user') return i18n.global.t('role.user')
  return i18n.global.t('role.unknown')
}

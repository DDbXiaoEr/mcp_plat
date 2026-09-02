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
import { i18n, initialLocale, SUPPORTED_LOCALES } from '../i18n.js'

const STORAGE_KEY = 'mcp-console-locale'

const state = reactive({
  locale: initialLocale
})

function applyLocale(locale) {
  i18n.global.locale.value = locale
  document.documentElement.lang = locale
  document.title = i18n.global.t('app.title')
  const meta = document.querySelector('meta[name="description"]')
  if (meta) meta.setAttribute('content', i18n.global.t('app.metaDescription'))
}

applyLocale(state.locale)

function setLocale(locale) {
  if (!SUPPORTED_LOCALES.includes(locale)) locale = 'zh-CN'
  state.locale = locale
  try {
    localStorage.setItem(STORAGE_KEY, locale)
  } catch (_) { /* ignore */ }
  applyLocale(locale)
}

export const locale = readonly(state)
export { setLocale }

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

// Author: deepseek-v4-pro / opencode
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { auth } from '../stores/auth.js'
import { nav, menusFor, setActive } from '../stores/nav.js'
import { theme, toggleTheme } from '../stores/theme.js'

const { t } = useI18n()

const menus = computed(() => menusFor(auth.user?.role))
</script>

<template>
  <aside class="sidebar">
    <nav class="sidebar__nav">
      <button
        v-for="item in menus"
        :key="item.key"
        class="sidebar__item"
        :class="{ 'sidebar__item--active': nav.active === item.key }"
        type="button"
        @click="setActive(item.key)"
      >
        {{ t(item.labelKey) }}
      </button>
    </nav>
    <div class="sidebar__footer">
      <button
        class="theme-toggle"
        type="button"
        role="switch"
        :aria-checked="theme.theme === 'dark'"
        @click="toggleTheme()"
      >
        <span class="theme-toggle__track">
          <span class="theme-toggle__thumb"></span>
        </span>
        <span class="theme-toggle__label">
          {{ theme.theme === 'dark' ? t('app.themeDark') : t('app.themeLight') }}
        </span>
      </button>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: var(--header-height);
  left: 0;
  bottom: 0;
  z-index: 90;
  display: flex;
  flex-direction: column;
  width: var(--sidebar-width);
  background: var(--surface);
  border-right: 1px solid var(--border);
}

.sidebar__nav {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px 12px;
  overflow-y: auto;
}

.sidebar__item {
  display: flex;
  align-items: center;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--text);
  text-align: left;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s, color 0.2s;
}

.sidebar__item:hover {
  background: var(--hover-bg);
}

.sidebar__item--active {
  color: var(--xauat-blue);
  font-weight: 600;
  background: var(--active-bg);
}

.sidebar__footer {
  padding: 12px;
  border-top: 1px solid var(--border);
}

.theme-toggle {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  padding: 10px 14px;
  font-size: 14px;
  color: var(--text);
  text-align: left;
  background: transparent;
  border: none;
  border-radius: 10px;
  cursor: pointer;
  transition: background 0.2s;
}

.theme-toggle:hover {
  background: var(--hover-bg);
}

.theme-toggle__track {
  position: relative;
  flex-shrink: 0;
  width: 40px;
  height: 22px;
  border-radius: 999px;
  background: var(--border);
  transition: background 0.2s;
}

.theme-toggle__thumb {
  position: absolute;
  top: 3px;
  left: 3px;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: var(--surface);
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.25);
  transition: transform 0.2s;
}

.theme-toggle[aria-checked="true"] .theme-toggle__track {
  background: var(--xauat-blue-light);
}

.theme-toggle[aria-checked="true"] .theme-toggle__thumb {
  transform: translateX(18px);
}

.theme-toggle__label {
  color: var(--text-muted);
}
</style>

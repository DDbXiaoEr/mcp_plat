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
import { computed, h } from 'vue'
import { useI18n } from 'vue-i18n'
import { auth } from '../stores/auth.js'
import { nav, menusFor, setActive } from '../stores/nav.js'
import { theme, toggleTheme } from '../stores/theme.js'

const { t } = useI18n()

const menus = computed(() => menusFor(auth.user?.role))

const icons = {
  user: [
    ['circle', { cx: '12', cy: '8', r: '3.25' }],
    ['path', { d: 'M5.5 19.2c.9-3.1 3.4-5 6.5-5s5.6 1.9 6.5 5' }]
  ],
  key: [
    ['circle', { cx: '8', cy: '12', r: '3.25' }],
    ['path', { d: 'M11 12h8.5M16.5 12v3M19.5 12v2' }]
  ],
  spark: [
    ['path', { d: 'M12 3.5v4M12 16.5v4M3.5 12h4M16.5 12h4M7.1 7.1l2.8 2.8M14.1 14.1l2.8 2.8M16.9 7.1l-2.8 2.8M9.9 14.1l-2.8 2.8' }]
  ],
  clock: [
    ['circle', { cx: '12', cy: '12', r: '7.25' }],
    ['path', { d: 'M12 8.2V12l2.6 2.2' }]
  ],
  layout: [
    ['rect', { x: '4.5', y: '5', width: '15', height: '14', rx: '1.5' }],
    ['path', { d: 'M4.5 10h15M10 10v9' }]
  ],
  server: [
    ['rect', { x: '5', y: '4.5', width: '14', height: '5.5', rx: '1.2' }],
    ['rect', { x: '5', y: '14', width: '14', height: '5.5', rx: '1.2' }],
    ['path', { d: 'M8 7.2h.01M8 16.7h.01' }]
  ],
  users: [
    ['circle', { cx: '9', cy: '8.5', r: '2.6' }],
    ['path', { d: 'M4.5 18c.7-2.6 2.6-4 4.5-4s3.8 1.4 4.5 4' }],
    ['circle', { cx: '16', cy: '9', r: '2.1' }],
    ['path', { d: 'M14.6 18c.4-1.7 1.6-2.8 3-2.8 1.5 0 2.7 1.2 3.1 3' }]
  ],
  sliders: [
    ['path', { d: 'M5 7h14M5 12h14M5 17h14' }],
    ['circle', { cx: '9', cy: '7', r: '1.6' }],
    ['circle', { cx: '15', cy: '12', r: '1.6' }],
    ['circle', { cx: '11', cy: '17', r: '1.6' }]
  ]
}

const NavIcon = {
  props: { name: { type: String, required: true } },
  setup(props) {
    return () => {
      const parts = icons[props.name] || icons.layout
      return h(
        'svg',
        {
          class: 'sidebar__icon',
          viewBox: '0 0 24 24',
          fill: 'none',
          stroke: 'currentColor',
          strokeWidth: '1.7',
          strokeLinecap: 'round',
          strokeLinejoin: 'round',
          'aria-hidden': 'true'
        },
        parts.map(([tag, attrs]) => h(tag, attrs))
      )
    }
  }
}
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
        <NavIcon :name="item.icon" />
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
  gap: 2px;
  padding: 12px 10px;
  overflow-y: auto;
}

.sidebar__item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 9px 12px;
  font-size: 14px;
  color: var(--text);
  text-align: left;
  background: transparent;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  transition: background 0.15s, color 0.15s;
}

.sidebar__item:hover {
  background: var(--hover-bg);
}

.sidebar__item--active {
  color: var(--xauat-blue);
  font-weight: 600;
  background: var(--active-bg);
  box-shadow: inset 3px 0 0 var(--xauat-accent);
}

.sidebar__icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
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
  border-radius: 6px;
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
  transition: background 0.15s;
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

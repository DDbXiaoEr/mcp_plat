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
import { auth } from '../stores/auth.js'
import { nav, menusFor, setActive } from '../stores/nav.js'

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
        {{ item.label }}
      </button>
    </nav>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: var(--header-height);
  left: 0;
  bottom: 0;
  z-index: 90;
  width: var(--sidebar-width);
  background: var(--surface);
  border-right: 1px solid var(--border);
  overflow-y: auto;
}

.sidebar__nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 16px 12px;
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
  background: rgba(10, 61, 122, 0.06);
}

.sidebar__item--active {
  color: var(--xauat-blue);
  font-weight: 600;
  background: rgba(30, 95, 176, 0.1);
}
</style>

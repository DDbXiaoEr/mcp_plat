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
import { menusFor, setActive } from '../stores/nav.js'
import { roleLabel } from '../i18n.js'

const { t } = useI18n()

const menus = computed(() => menusFor(auth.user?.role))
</script>

<template>
  <section v-if="auth.user" class="welcome">
    <header class="welcome__head">
      <h1 class="welcome__title">{{ t('welcome.title', { role: roleLabel(auth.user.role) }) }}</h1>
      <p class="welcome__lead">{{ t('welcome.lead') }}</p>
    </header>
    <nav class="welcome__index" :aria-label="t('welcome.indexLabel')">
      <button
        v-for="item in menus"
        :key="item.key"
        class="welcome__item"
        type="button"
        @click="setActive(item.key)"
      >
        <span class="welcome__name">{{ t(item.labelKey) }}</span>
        <span class="welcome__hint">{{ t('welcome.hint.' + item.key) }}</span>
      </button>
    </nav>
  </section>
</template>

<style scoped>
.welcome {
  max-width: 720px;
}

.welcome__head {
  margin-bottom: 28px;
}

.welcome__title {
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.03em;
  line-height: 1.25;
  color: var(--xauat-blue);
}

.welcome__lead {
  margin-top: 8px;
  font-size: 15px;
  color: var(--text-muted);
}

.welcome__index {
  display: flex;
  flex-direction: column;
  border: 1px solid var(--border);
  border-left: 3px solid var(--xauat-accent);
  background: var(--surface);
}

.welcome__item {
  display: grid;
  grid-template-columns: minmax(140px, 0.42fr) 1fr;
  gap: 16px;
  align-items: baseline;
  width: 100%;
  padding: 16px 20px;
  text-align: left;
  color: var(--text);
  background: transparent;
  border: none;
  border-bottom: 1px solid var(--border);
  cursor: pointer;
  transition: background 0.15s;
}

.welcome__item:last-child {
  border-bottom: none;
}

.welcome__item:hover {
  background: var(--hover-bg);
}

.welcome__item:hover .welcome__name {
  color: var(--xauat-blue);
}

.welcome__name {
  font-size: 15px;
  font-weight: 600;
}

.welcome__hint {
  font-size: 14px;
  color: var(--text-muted);
}

@media (max-width: 560px) {
  .welcome__item {
    grid-template-columns: 1fr;
    gap: 4px;
  }
}
</style>

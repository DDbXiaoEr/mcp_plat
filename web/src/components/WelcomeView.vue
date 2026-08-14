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
import { auth } from '../stores/auth.js'
</script>

<template>
  <!-- TODO: 后期按角色拆分为不同的功能页面 -->
  <section v-if="auth.user" class="welcome">
    <h1 class="welcome__title">欢迎，{{ auth.user.roleLabel }}登录</h1>
  </section>
</template>

<style scoped>
.welcome {
  display: grid;
  place-items: center;
  min-height: calc(100vh - var(--header-height) - 48px);
}

.welcome__title {
  font-size: 28px;
  font-weight: 700;
  color: var(--xauat-blue);
}
</style>

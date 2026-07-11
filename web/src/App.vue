<script setup>
import { computed, watch } from 'vue'
import { auth } from './stores/auth.js'
import { nav, resetNav } from './stores/nav.js'
import LoginView from './components/LoginView.vue'
import TheHeader from './components/TheHeader.vue'
import TheSidebar from './components/TheSidebar.vue'
import WelcomeView from './components/WelcomeView.vue'
import ProfileView from './components/ProfileView.vue'
import AccessKeyView from './components/AccessKeyView.vue'
import HistoryView from './components/HistoryView.vue'

const views = {
  profile: ProfileView,
  accesskey: AccessKeyView,
  history: HistoryView
}

const currentView = computed(() => views[nav.active] || WelcomeView)

watch(() => auth.user, () => resetNav())
</script>

<template>
  <LoginView v-if="!auth.user" />
  <template v-else>
    <TheHeader />
    <TheSidebar />
    <main class="content">
      <component :is="currentView" />
    </main>
  </template>
</template>

<style scoped>
.content {
  margin-top: var(--header-height);
  margin-left: var(--sidebar-width);
  min-height: calc(100vh - var(--header-height));
  padding: 24px;
}
</style>

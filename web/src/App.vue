<script setup>
import { computed, watch, onMounted, ref } from 'vue'
import { auth, casLogin } from './stores/auth.js'
import { nav, resetNav } from './stores/nav.js'
import LoginView from './components/LoginView.vue'
import TheHeader from './components/TheHeader.vue'
import TheSidebar from './components/TheSidebar.vue'
import WelcomeView from './components/WelcomeView.vue'
import ProfileView from './components/ProfileView.vue'
import AccessKeyView from './components/AccessKeyView.vue'
import HistoryView from './components/HistoryView.vue'
import SettingsView from './components/SettingsView.vue'
import RbacView from './components/RbacView.vue'
import OverviewView from './components/OverviewView.vue'
import ServersView from './components/ServersView.vue'

const views = {
  profile: ProfileView,
  accesskey: AccessKeyView,
  history: HistoryView,
  settings: SettingsView,
  rbac: RbacView,
  overview: OverviewView,
  servers: ServersView
}

const currentView = computed(() => views[nav.active] || WelcomeView)

watch(() => auth.user, () => resetNav())

const casError = ref('')

onMounted(async () => {
  const params = new URLSearchParams(window.location.search)
  const ticket = params.get('ticket')
  if (ticket && !auth.user) {
    const serviceUrl = window.location.origin + window.location.pathname
    const result = await casLogin(ticket, serviceUrl)
    if (!result.ok) {
      casError.value = result.message
    } else {
      const url = new URL(window.location.href)
      url.searchParams.delete('ticket')
      window.history.replaceState({}, '', url.toString())
    }
  }
})
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

<template>
  <div class="app-shell">
    <NavBar @open-search="searchOpen = true" />
    <main class="page-content">
      <RouterView v-slot="{ Component, route }">
        <Transition name="page" mode="out-in">
          <component :is="Component" :key="route.path" />
        </Transition>
      </RouterView>
    </main>
    <FooterSection />
    <SearchModal :open="searchOpen" @close="searchOpen = false" />
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import FooterSection from './components/FooterSection.vue'
import NavBar from './components/NavBar.vue'
import SearchModal from './components/SearchModal.vue'

const searchOpen = ref(false)

function onKeydown(e) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    searchOpen.value = !searchOpen.value
  }
}

onMounted(() => window.addEventListener('keydown', onKeydown))
onUnmounted(() => window.removeEventListener('keydown', onKeydown))
</script>

<style>
.app-shell {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

/* ── Page transition ────────────────────────────────────────────────────────── */
.page-enter-active,
.page-leave-active {
  transition:
    opacity 180ms var(--ease),
    transform 180ms var(--ease);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(8px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>

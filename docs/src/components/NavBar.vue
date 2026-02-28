<template>
  <header
    class="navbar"
    :class="{ 'navbar--scrolled': scrolled, 'navbar--open': menuOpen }"
    @keydown.escape="closeMenu"
  >
    <div class="navbar__inner">
      <!-- Logo -->
      <RouterLink to="/" class="navbar__logo" @click="closeMenu">
        <span class="navbar__logo-icon">🧪</span>
        <span class="navbar__logo-text">gest</span>
        <span class="navbar__version">v1.1.0</span>
      </RouterLink>

      <!-- Desktop Nav -->
      <nav class="navbar__nav" aria-label="Main navigation">
        <RouterLink
          v-for="link in links"
          :key="link.to"
          :to="link.to"
          class="navbar__link"
          :class="{ 'navbar__link--active': isActive(link.to) }"
        >
          <span class="navbar__link-icon">{{ link.icon }}</span>
          {{ link.label }}
        </RouterLink>
      </nav>

      <!-- Right side -->
      <div class="navbar__actions">
        <!-- Search trigger -->
        <button
          class="navbar__search"
          @click="$emit('openSearch')"
          aria-label="Search documentation (Ctrl+K)"
        >
          <svg
            width="14"
            height="14"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            aria-hidden="true"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
          <span class="navbar__search-placeholder">Search…</span>
          <span class="navbar__search-kbd" aria-hidden="true">
            <kbd>{{ metaKey }}</kbd
            ><kbd>K</kbd>
          </span>
        </button>

        <a
          href="https://github.com/caiolandgraf/gest"
          target="_blank"
          rel="noopener noreferrer"
          class="navbar__github btn btn--ghost btn--sm"
          aria-label="View gest on GitHub"
        >
          <svg
            width="16"
            height="16"
            viewBox="0 0 24 24"
            fill="currentColor"
            aria-hidden="true"
          >
            <path
              d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
            />
          </svg>
          GitHub
        </a>

        <!-- Mobile search icon -->
        <button
          class="navbar__search-icon"
          @click="$emit('openSearch')"
          aria-label="Search documentation"
        >
          <svg
            width="18"
            height="18"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            aria-hidden="true"
          >
            <circle cx="11" cy="11" r="8" />
            <path d="m21 21-4.35-4.35" />
          </svg>
        </button>

        <!-- Mobile toggle -->
        <button
          class="navbar__burger"
          :class="{ 'navbar__burger--open': menuOpen }"
          @click="toggleMenu"
          :aria-expanded="menuOpen"
          aria-controls="mobile-menu"
          aria-label="Toggle menu"
        >
          <span></span>
          <span></span>
          <span></span>
        </button>
      </div>
    </div>

    <!-- Mobile Menu -->
    <Transition name="mobile-menu">
      <div v-if="menuOpen" id="mobile-menu" class="navbar__mobile">
        <nav class="navbar__mobile-nav" aria-label="Mobile navigation">
          <RouterLink
            v-for="link in links"
            :key="link.to"
            :to="link.to"
            class="navbar__mobile-link"
            :class="{ 'navbar__mobile-link--active': isActive(link.to) }"
            @click="closeMenu"
          >
            <span class="navbar__mobile-link-icon">{{ link.icon }}</span>
            <span>{{ link.label }}</span>
            <svg
              class="navbar__mobile-arrow"
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
              aria-hidden="true"
            >
              <path d="M9 18l6-6-6-6" />
            </svg>
          </RouterLink>
        </nav>

        <div class="navbar__mobile-footer">
          <a
            href="https://github.com/caiolandgraf/gest"
            target="_blank"
            rel="noopener noreferrer"
            class="btn btn--ghost"
            style="width: 100%; justify-content: center"
            @click="closeMenu"
          >
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="currentColor"
              aria-hidden="true"
            >
              <path
                d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
              />
            </svg>
            View on GitHub
          </a>
        </div>
      </div>
    </Transition>
  </header>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'

defineEmits(['openSearch'])

const route = useRoute()
const scrolled = ref(false)
const menuOpen = ref(false)

const metaKey = computed(() =>
  typeof navigator !== 'undefined' &&
  /Mac|iPhone|iPad|iPod/.test(navigator.platform)
    ? '⌘'
    : 'Ctrl'
)

const links = [
  { to: '/', label: 'Home', icon: '⌂' },
  { to: '/docs', label: 'Docs', icon: '📖' },
  { to: '/contributors', label: 'Contributors', icon: '👥' }
]

function isActive(path) {
  if (path === '/') return route.path === '/'
  return route.path.startsWith(path)
}

function toggleMenu() {
  menuOpen.value = !menuOpen.value
  document.body.style.overflow = menuOpen.value ? 'hidden' : ''
}

function closeMenu() {
  menuOpen.value = false
  document.body.style.overflow = ''
}

function onScroll() {
  scrolled.value = window.scrollY > 12
}

onMounted(() => {
  window.addEventListener('scroll', onScroll, { passive: true })
  onScroll()
})

onUnmounted(() => {
  window.removeEventListener('scroll', onScroll)
  document.body.style.overflow = ''
})
</script>

<style scoped>
/* ── Base ───────────────────────────────────────────────────────────────────── */
.navbar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: var(--nav-height);
  background: rgba(10, 12, 16, 0.75);
  backdrop-filter: blur(16px) saturate(180%);
  -webkit-backdrop-filter: blur(16px) saturate(180%);
  border-bottom: 1px solid transparent;
  transition:
    border-color var(--duration-normal) var(--ease),
    background var(--duration-normal) var(--ease),
    box-shadow var(--duration-normal) var(--ease);
}

.navbar--scrolled {
  background: rgba(10, 12, 16, 0.92);
  border-bottom-color: var(--border-subtle);
  box-shadow: 0 1px 24px rgba(0, 0, 0, 0.4);
}

.navbar--open {
  background: var(--bg-1);
  border-bottom-color: var(--border-subtle);
}

/* ── Inner layout ───────────────────────────────────────────────────────────── */
.navbar__inner {
  max-width: var(--content-max);
  margin: 0 auto;
  height: var(--nav-height);
  padding: 0 1.5rem;
  display: flex;
  align-items: center;
  gap: 2rem;
}

/* ── Logo ───────────────────────────────────────────────────────────────────── */
.navbar__logo {
  display: flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  flex-shrink: 0;
}

.navbar__logo-icon {
  font-size: 1.35rem;
  line-height: 1;
  filter: drop-shadow(0 0 8px rgba(63, 185, 80, 0.4));
}

.navbar__logo-text {
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text-1);
  letter-spacing: -0.03em;
  font-family: var(--font-mono);
  transition: color var(--duration-fast) var(--ease);
}

.navbar__logo:hover .navbar__logo-text {
  color: var(--green);
}

.navbar__version {
  font-size: 0.65rem;
  font-weight: 600;
  font-family: var(--font-mono);
  color: var(--green);
  background: var(--green-glow);
  border: 1px solid rgba(63, 185, 80, 0.2);
  border-radius: var(--radius-full);
  padding: 0.1rem 0.45rem;
  margin-left: 2px;
  letter-spacing: 0.02em;
}

/* ── Desktop nav ────────────────────────────────────────────────────────────── */
.navbar__nav {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  flex: 1;
}

.navbar__link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 0.4rem 0.85rem;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-3);
  text-decoration: none;
  transition:
    color var(--duration-fast) var(--ease),
    background var(--duration-fast) var(--ease);
  position: relative;
}

.navbar__link:hover {
  color: var(--text-1);
  background: var(--bg-3);
}

.navbar__link--active {
  color: var(--text-1) !important;
  background: var(--bg-3);
}

.navbar__link--active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 50%;
  transform: translateX(-50%);
  width: 20px;
  height: 2px;
  background: var(--green);
  border-radius: var(--radius-full);
  box-shadow: 0 0 8px var(--green);
}

.navbar__link-icon {
  font-size: 0.85rem;
  opacity: 0.7;
}

/* ── Search trigger ─────────────────────────────────────────────────────────── */
.navbar__search {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 0.35rem 0.65rem 0.35rem 0.6rem;
  background: var(--bg-3);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-sm);
  color: var(--text-4);
  font-family: var(--font-sans);
  font-size: 0.82rem;
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease);
  white-space: nowrap;
  min-width: 180px;
}

.navbar__search:hover {
  background: var(--bg-4);
  border-color: var(--border);
  color: var(--text-3);
}

.navbar__search-placeholder {
  flex: 1;
  text-align: left;
}

.navbar__search-kbd {
  display: flex;
  align-items: center;
  gap: 2px;
  margin-left: auto;
}

.navbar__search-kbd kbd {
  font-family: var(--font-mono);
  font-size: 0.62rem;
  background: var(--bg-5);
  border: 1px solid var(--border);
  border-radius: 3px;
  padding: 0.1rem 0.3rem;
  color: var(--text-4);
  line-height: 1.6;
}

/* Mobile search icon (shown on small screens) */
.navbar__search-icon {
  display: none;
  align-items: center;
  justify-content: center;
  width: 36px;
  height: 36px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid var(--border);
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--duration-fast) var(--ease);
}

.navbar__search-icon:hover {
  background: var(--bg-3);
  color: var(--text-1);
}

/* ── Actions ────────────────────────────────────────────────────────────────── */
.navbar__actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  margin-left: auto;
}

.navbar__github {
  display: flex;
  align-items: center;
  gap: 7px;
}

/* ── Burger ─────────────────────────────────────────────────────────────────── */
.navbar__burger {
  display: none;
  flex-direction: column;
  justify-content: center;
  gap: 5px;
  width: 36px;
  height: 36px;
  padding: 6px;
  border-radius: var(--radius-sm);
  background: transparent;
  border: 1px solid var(--border);
  cursor: pointer;
  transition:
    background var(--duration-fast) var(--ease),
    border-color var(--duration-fast) var(--ease);
}

.navbar__burger:hover {
  background: var(--bg-3);
  border-color: var(--text-4);
}

.navbar__burger span {
  display: block;
  width: 100%;
  height: 1.5px;
  background: var(--text-2);
  border-radius: 2px;
  transform-origin: center;
  transition:
    transform var(--duration-normal) var(--ease),
    opacity var(--duration-normal) var(--ease);
}

.navbar__burger--open span:nth-child(1) {
  transform: translateY(6.5px) rotate(45deg);
}

.navbar__burger--open span:nth-child(2) {
  opacity: 0;
  transform: scaleX(0);
}

.navbar__burger--open span:nth-child(3) {
  transform: translateY(-6.5px) rotate(-45deg);
}

/* ── Mobile menu ────────────────────────────────────────────────────────────── */
.navbar__mobile {
  position: fixed;
  top: var(--nav-height);
  left: 0;
  right: 0;
  bottom: 0;
  background: var(--bg-1);
  border-top: 1px solid var(--border-subtle);
  display: flex;
  flex-direction: column;
  padding: 1.5rem;
  overflow-y: auto;
  z-index: 99;
}

.navbar__mobile-nav {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  flex: 1;
}

.navbar__mobile-link {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 1rem 1.25rem;
  border-radius: var(--radius-md);
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-2);
  text-decoration: none;
  border: 1px solid transparent;
  transition: all var(--duration-fast) var(--ease);
}

.navbar__mobile-link:hover,
.navbar__mobile-link--active {
  color: var(--text-1);
  background: var(--bg-3);
  border-color: var(--border-subtle);
}

.navbar__mobile-link--active {
  color: var(--green) !important;
  border-color: rgba(63, 185, 80, 0.2) !important;
  background: var(--green-glow) !important;
}

.navbar__mobile-link-icon {
  font-size: 1.1rem;
  width: 24px;
  text-align: center;
  flex-shrink: 0;
}

.navbar__mobile-arrow {
  margin-left: auto;
  color: var(--text-4);
  flex-shrink: 0;
}

.navbar__mobile-footer {
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-subtle);
  margin-top: 1.5rem;
}

/* ── Transition ─────────────────────────────────────────────────────────────── */
.mobile-menu-enter-active,
.mobile-menu-leave-active {
  transition:
    opacity var(--duration-normal) var(--ease),
    transform var(--duration-normal) var(--ease);
}

.mobile-menu-enter-from,
.mobile-menu-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 768px) {
  .navbar__nav,
  .navbar__github,
  .navbar__search {
    display: none;
  }

  .navbar__search-icon {
    display: flex;
  }

  .navbar__burger {
    display: flex;
  }
}
</style>

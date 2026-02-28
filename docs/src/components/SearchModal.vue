<template>
  <Teleport to="body">
    <Transition name="s-overlay">
      <div
        v-if="open"
        class="s-overlay"
        @mousedown.self="$emit('close')"
        aria-hidden="true"
      ></div>
    </Transition>

    <Transition name="s-modal">
      <div
        v-if="open"
        class="s-wrap"
        role="dialog"
        aria-modal="true"
        aria-label="Search documentation"
        @keydown.escape.prevent="$emit('close')"
      >
        <div class="s-modal">
          <!-- ── Input ──────────────────────────────────────────────────── -->
          <div class="s-header">
            <svg
              class="s-header__icon"
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

            <input
              ref="inputRef"
              v-model="query"
              type="text"
              class="s-input"
              placeholder="Search docs, matchers, API…"
              autocomplete="off"
              spellcheck="false"
              @keydown.up.prevent="moveUp"
              @keydown.down.prevent="moveDown"
              @keydown.enter.prevent="selectActive"
            />

            <button class="s-esc" @click="$emit('close')" tabindex="-1">
              Esc
            </button>
          </div>

          <!-- ── Body ──────────────────────────────────────────────────── -->
          <div class="s-body" ref="bodyRef">
            <!-- Results -->
            <template v-if="groups.length">
              <div v-for="group in groups" :key="group.label" class="s-group">
                <p class="s-group__label">{{ group.label }}</p>
                <ul class="s-group__list" role="listbox">
                  <li
                    v-for="item in group.items"
                    :key="item._idx"
                    role="option"
                    :aria-selected="activeIndex === item._idx"
                    class="s-result"
                    :class="{ 's-result--active': activeIndex === item._idx }"
                    @mouseenter="activeIndex = item._idx"
                    @mousedown.prevent="navigate(item)"
                  >
                    <span class="s-result__icon" aria-hidden="true">{{
                      item.icon
                    }}</span>
                    <span class="s-result__body">
                      <span
                        class="s-result__title"
                        v-html="highlight(item.title)"
                      ></span>
                      <span v-if="item.desc" class="s-result__desc">{{
                        item.desc
                      }}</span>
                    </span>
                    <span class="s-result__badge" :data-type="item.type">{{
                      item.type
                    }}</span>
                    <svg
                      class="s-result__arrow"
                      width="14"
                      height="14"
                      viewBox="0 0 24 24"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="2"
                      aria-hidden="true"
                    >
                      <path d="M9 18l6-6-6-6" />
                    </svg>
                  </li>
                </ul>
              </div>
            </template>

            <!-- No results -->
            <div v-else-if="query.trim()" class="s-empty">
              <svg
                width="32"
                height="32"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="1.5"
                aria-hidden="true"
              >
                <circle cx="11" cy="11" r="8" />
                <path d="m21 21-4.35-4.35" />
              </svg>
              <p>
                No results for <strong>"{{ query }}"</strong>
              </p>
              <span>Try searching for a matcher name or doc section.</span>
            </div>

            <!-- Initial hint -->
            <div v-else class="s-hint">
              <div class="s-hint__grid">
                <button
                  v-for="quick in quickLinks"
                  :key="quick.title"
                  class="s-hint__item"
                  @mousedown.prevent="navigate(quick)"
                >
                  <span class="s-hint__icon">{{ quick.icon }}</span>
                  <span>{{ quick.title }}</span>
                </button>
              </div>
            </div>
          </div>

          <!-- ── Footer ─────────────────────────────────────────────────── -->
          <div class="s-footer">
            <span class="s-footer__shortcuts">
              <span class="s-footer__key-group">
                <kbd>↑</kbd><kbd>↓</kbd> navigate
              </span>
              <span class="s-footer__key-group"> <kbd>↵</kbd> open </span>
              <span class="s-footer__key-group"> <kbd>Esc</kbd> close </span>
            </span>
            <span class="s-footer__brand">
              <span>gest</span> <span class="s-footer__emoji">🧪</span> docs
            </span>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'

/* ── Props / emits ─────────────────────────────────────────────────────────── */
const props = defineProps({ open: Boolean })
const emit = defineEmits(['close'])

const router = useRouter()
const inputRef = ref(null)
const bodyRef = ref(null)
const query = ref('')
const activeIndex = ref(0)

/* ── Reset when modal opens / closes ──────────────────────────────────────── */
watch(
  () => props.open,
  val => {
    if (val) {
      query.value = ''
      activeIndex.value = 0
      nextTick(() => inputRef.value?.focus())
    }
  }
)

/* ── Search index ──────────────────────────────────────────────────────────── */
const allItems = [
  // ── Pages ──────────────────────────────────────────────────────────────
  {
    type: 'page',
    icon: '⌂',
    title: 'Home',
    desc: 'Landing page — features and quick start',
    route: '/'
  },
  {
    type: 'page',
    icon: '📖',
    title: 'Documentation',
    desc: 'Full reference — installation, matchers, API',
    route: '/docs'
  },
  {
    type: 'page',
    icon: '👥',
    title: 'Contributors',
    desc: 'Meet the people behind gest',
    route: '/contributors'
  },

  // ── Doc sections ───────────────────────────────────────────────────────
  {
    type: 'section',
    icon: '📦',
    title: 'Installation',
    desc: 'go get github.com/caiolandgraf/gest',
    route: '/docs',
    hash: 'installation'
  },
  {
    type: 'section',
    icon: '📁',
    title: 'Project Structure',
    desc: '_spec.go files alongside your code',
    route: '/docs',
    hash: 'structure'
  },
  {
    type: 'section',
    icon: '🚀',
    title: 'Basic Usage',
    desc: 'Describe, It, Register, RunRegistered',
    route: '/docs',
    hash: 'usage'
  },
  {
    type: 'section',
    icon: '▶️',
    title: 'Running Tests',
    desc: 'go run . and go run . -c',
    route: '/docs',
    hash: 'running'
  },
  {
    type: 'section',
    icon: '❌',
    title: 'Failure Messages',
    desc: 'Expected vs Received diff with code snippet',
    route: '/docs',
    hash: 'failures'
  },
  {
    type: 'section',
    icon: '📊',
    title: 'Coverage',
    desc: 'Per-suite pass-rate table with progress bars',
    route: '/docs',
    hash: 'coverage'
  },
  {
    type: 'section',
    icon: '🔗',
    title: 'Matchers',
    desc: 'All 10 built-in assertion matchers',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'section',
    icon: '🔄',
    title: 'Negation',
    desc: '.Not() negates any matcher',
    route: '/docs',
    hash: 'negation'
  },
  {
    type: 'section',
    icon: '📚',
    title: 'Full API',
    desc: 'Describe, It, Register, RunAll, Expect, Not',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'section',
    icon: '👁',
    title: 'Watch Mode',
    desc: '--watch re-runs tests on every .go file change',
    route: '/docs',
    hash: 'watch'
  },
  {
    type: 'section',
    icon: '💡',
    title: 'Philosophy',
    desc: 'Minimal deps, zero config, auto-discovery',
    route: '/docs',
    hash: 'philosophy'
  },
  {
    type: 'section',
    icon: '🧪',
    title: 'Example Project',
    desc: 'Runnable example in example/ directory',
    route: '/docs',
    hash: 'example'
  },

  // ── Matchers ───────────────────────────────────────────────────────────
  {
    type: 'matcher',
    icon: '=',
    title: '.ToBe(v)',
    desc: 'Strict equality (==) — best for primitives',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '≡',
    title: '.ToEqual(v)',
    desc: 'Deep equality via reflect.DeepEqual',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '∅',
    title: '.ToBeNil()',
    desc: 'Value is nil — works with pointers and interfaces',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '✓',
    title: '.ToBeTrue()',
    desc: 'Value is exactly true',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '✗',
    title: '.ToBeFalse()',
    desc: 'Value is exactly false',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '⊃',
    title: '.ToContain(s)',
    desc: 'String contains the given substring',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '#',
    title: '.ToHaveLength(n)',
    desc: 'Length of string, slice, array or map',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '>',
    title: '.ToBeGreaterThan(n)',
    desc: 'Number strictly greater than n',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '<',
    title: '.ToBeLessThan(n)',
    desc: 'Number strictly less than n',
    route: '/docs',
    hash: 'matchers'
  },
  {
    type: 'matcher',
    icon: '≈',
    title: '.ToBeCloseTo(n, d?)',
    desc: 'Float approximately equal (default delta ±0.001)',
    route: '/docs',
    hash: 'matchers'
  },

  // ── API symbols ────────────────────────────────────────────────────────
  {
    type: 'api',
    icon: '◎',
    title: 'gest.Describe(name)',
    desc: 'Create a new test suite',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: '(*Suite).It(name, fn)',
    desc: 'Add a test case to a suite',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: 'gest.Register(s)',
    desc: 'Add suite to the global registry',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: 'gest.RunRegistered()',
    desc: 'Run all registered suites from main()',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: 'gest.RunAll(...)',
    desc: 'Run specific suites manually',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: '(*T).Expect(v)',
    desc: 'Start an assertion chain',
    route: '/docs',
    hash: 'api'
  },
  {
    type: 'api',
    icon: '◎',
    title: '(*Expectation).Not()',
    desc: 'Negate the next matcher',
    route: '/docs',
    hash: 'api'
  }
].map((item, i) => ({ ...item, _idx: i }))

const typeOrder = { page: 0, section: 1, matcher: 2, api: 3 }

/* ── Filtered results ──────────────────────────────────────────────────────── */
const results = computed(() => {
  const q = query.value.trim().toLowerCase()
  if (!q) return []

  return allItems
    .filter(
      item =>
        item.title.toLowerCase().includes(q) ||
        item.desc?.toLowerCase().includes(q) ||
        item.type.toLowerCase().includes(q)
    )
    .sort((a, b) => {
      // Exact title match first
      const aExact = a.title.toLowerCase().startsWith(q)
      const bExact = b.title.toLowerCase().startsWith(q)
      if (aExact !== bExact) return aExact ? -1 : 1
      return (typeOrder[a.type] ?? 9) - (typeOrder[b.type] ?? 9)
    })
    .slice(0, 10)
})

/* ── Grouped results ───────────────────────────────────────────────────────── */
const GROUP_LABELS = {
  page: 'Pages',
  section: 'Sections',
  matcher: 'Matchers',
  api: 'API'
}

const groups = computed(() => {
  if (!results.value.length) return []
  const map = new Map()
  for (const item of results.value) {
    if (!map.has(item.type)) map.set(item.type, [])
    map.get(item.type).push(item)
  }
  return Array.from(map.entries()).map(([type, items]) => ({
    label: GROUP_LABELS[type] ?? type,
    items
  }))
})

/* ── Quick links shown when query is empty ────────────────────────────────── */
const quickLinks = allItems.filter(
  i =>
    i.type === 'page' ||
    ['installation', 'usage', 'matchers', 'api'].includes(i.hash)
)

/* ── Keyboard navigation ───────────────────────────────────────────────────── */
const flatResults = computed(() => (results.value.length ? results.value : []))

function moveDown() {
  const max = flatResults.value.length - 1
  if (max < 0) return
  activeIndex.value = activeIndex.value < max ? activeIndex.value + 1 : 0
  scrollActiveIntoView()
}

function moveUp() {
  const max = flatResults.value.length - 1
  if (max < 0) return
  activeIndex.value = activeIndex.value > 0 ? activeIndex.value - 1 : max
  scrollActiveIntoView()
}

function selectActive() {
  const item =
    flatResults.value.find(i => i._idx === activeIndex.value) ??
    flatResults.value[0]
  if (item) navigate(item)
}

function scrollActiveIntoView() {
  nextTick(() => {
    bodyRef.value
      ?.querySelector('.s-result--active')
      ?.scrollIntoView({ block: 'nearest' })
  })
}

/* ── Navigation ────────────────────────────────────────────────────────────── */
async function navigate(item) {
  emit('close')
  await router.push(item.route ?? '/')
  if (item.hash) {
    await nextTick()
    setTimeout(() => {
      const el = document.getElementById(item.hash)
      if (!el) return
      window.scrollTo({
        top: el.getBoundingClientRect().top + window.scrollY - 88,
        behavior: 'smooth'
      })
    }, 120)
  }
}

/* ── Highlight matching text ───────────────────────────────────────────────── */
function escapeRegex(s) {
  return s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function highlight(text) {
  const q = query.value.trim()
  if (!q) return text
  return text.replace(
    new RegExp(`(${escapeRegex(q)})`, 'gi'),
    '<mark>$1</mark>'
  )
}

/* ── Reset active when results change ──────────────────────────────────────── */
watch(results, () => {
  activeIndex.value = flatResults.value[0]?._idx ?? 0
})
</script>

<style scoped>
/* ── Overlay ────────────────────────────────────────────────────────────────── */
.s-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
  z-index: 200;
}

/* ── Wrapper (centers the modal) ────────────────────────────────────────────── */
.s-wrap {
  position: fixed;
  inset: 0;
  z-index: 201;
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: clamp(60px, 12vh, 120px);
  padding-left: 1rem;
  padding-right: 1rem;
}

/* ── Modal ──────────────────────────────────────────────────────────────────── */
.s-modal {
  width: 100%;
  max-width: 600px;
  background: var(--bg-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow:
    0 24px 64px rgba(0, 0, 0, 0.7),
    0 0 0 1px rgba(255, 255, 255, 0.04);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  max-height: calc(100vh - 160px);
}

/* ── Header ─────────────────────────────────────────────────────────────────── */
.s-header {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 1rem;
  border-bottom: 1px solid var(--border-subtle);
  flex-shrink: 0;
}

.s-header__icon {
  color: var(--text-4);
  flex-shrink: 0;
}

.s-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  font-family: var(--font-sans);
  font-size: 1rem;
  color: var(--text-1);
  padding: 1rem 0;
  caret-color: var(--green);
}

.s-input::placeholder {
  color: var(--text-4);
}

.s-esc {
  font-family: var(--font-mono);
  font-size: 0.7rem;
  color: var(--text-4);
  background: var(--bg-4);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-xs);
  padding: 0.2rem 0.45rem;
  cursor: pointer;
  flex-shrink: 0;
  transition: all var(--duration-fast) var(--ease);
  white-space: nowrap;
}

.s-esc:hover {
  color: var(--text-2);
  border-color: var(--border);
}

/* ── Body ───────────────────────────────────────────────────────────────────── */
.s-body {
  overflow-y: auto;
  overscroll-behavior: contain;
  flex: 1;
  min-height: 0;
}

/* ── Group ──────────────────────────────────────────────────────────────────── */
.s-group {
  padding: 0.5rem 0;
}

.s-group + .s-group {
  border-top: 1px solid var(--border-subtle);
}

.s-group__label {
  font-size: 0.68rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-4);
  padding: 0.4rem 1rem 0.2rem;
}

.s-group__list {
  list-style: none;
}

/* ── Result item ────────────────────────────────────────────────────────────── */
.s-result {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0.6rem 1rem;
  cursor: pointer;
  border-radius: 0;
  transition: background var(--duration-fast) var(--ease);
}

.s-result--active {
  background: var(--bg-3);
}

.s-result__icon {
  font-size: 0.95rem;
  width: 22px;
  text-align: center;
  flex-shrink: 0;
  color: var(--text-3);
  font-family: var(--font-mono);
  font-style: normal;
}

.s-result__body {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.s-result__title {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.s-result__title :deep(mark) {
  background: transparent;
  color: var(--green-bright);
  font-weight: 700;
}

.s-result__desc {
  font-size: 0.78rem;
  color: var(--text-4);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.s-result__badge {
  font-size: 0.65rem;
  font-weight: 700;
  font-family: var(--font-mono);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 0.15rem 0.45rem;
  border-radius: var(--radius-xs);
  flex-shrink: 0;
  opacity: 0.85;
}

.s-result__badge[data-type='page'] {
  background: rgba(88, 166, 255, 0.12);
  color: var(--blue);
}
.s-result__badge[data-type='section'] {
  background: rgba(63, 185, 80, 0.1);
  color: var(--green-bright);
}
.s-result__badge[data-type='matcher'] {
  background: rgba(188, 140, 255, 0.1);
  color: var(--purple);
}
.s-result__badge[data-type='api'] {
  background: rgba(227, 179, 65, 0.1);
  color: var(--yellow);
}

.s-result__arrow {
  color: var(--text-4);
  flex-shrink: 0;
  opacity: 0;
  transform: translateX(-4px);
  transition:
    opacity var(--duration-fast) var(--ease),
    transform var(--duration-fast) var(--ease);
}

.s-result--active .s-result__arrow {
  opacity: 1;
  transform: translateX(0);
}

/* ── Empty state ────────────────────────────────────────────────────────────── */
.s-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
  padding: 2.5rem 1.5rem;
  text-align: center;
  color: var(--text-4);
}

.s-empty svg {
  opacity: 0.35;
  margin-bottom: 0.25rem;
}

.s-empty p {
  font-size: 0.9rem;
  color: var(--text-3);
  line-height: 1.5;
}

.s-empty p strong {
  color: var(--text-1);
}

.s-empty span {
  font-size: 0.78rem;
  color: var(--text-4);
}

/* ── Hint / quick links ─────────────────────────────────────────────────────── */
.s-hint {
  padding: 0.75rem;
}

.s-hint__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 4px;
}

.s-hint__item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 0.55rem 0.75rem;
  background: transparent;
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  font-size: 0.82rem;
  font-family: var(--font-sans);
  color: var(--text-3);
  cursor: pointer;
  text-align: left;
  transition: all var(--duration-fast) var(--ease);
}

.s-hint__item:hover {
  background: var(--bg-3);
  border-color: var(--border-subtle);
  color: var(--text-1);
}

.s-hint__icon {
  font-size: 0.9rem;
  flex-shrink: 0;
}

/* ── Footer ─────────────────────────────────────────────────────────────────── */
.s-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.55rem 1rem;
  border-top: 1px solid var(--border-subtle);
  background: var(--bg-1);
  flex-shrink: 0;
}

.s-footer__shortcuts {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.s-footer__key-group {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 0.72rem;
  color: var(--text-4);
  white-space: nowrap;
}

.s-footer kbd {
  font-family: var(--font-mono);
  font-size: 0.65rem;
  background: var(--bg-4);
  border: 1px solid var(--border-subtle);
  border-radius: 3px;
  padding: 0.1rem 0.35rem;
  color: var(--text-3);
  line-height: 1.6;
}

.s-footer__brand {
  font-size: 0.72rem;
  color: var(--text-4);
  display: flex;
  align-items: center;
  gap: 4px;
  font-family: var(--font-mono);
}

.s-footer__emoji {
  font-style: normal;
}

/* ── Transitions ────────────────────────────────────────────────────────────── */
.s-overlay-enter-active,
.s-overlay-leave-active {
  transition: opacity var(--duration-normal) var(--ease);
}
.s-overlay-enter-from,
.s-overlay-leave-to {
  opacity: 0;
}

.s-modal-enter-active {
  transition:
    opacity var(--duration-normal) var(--ease),
    transform var(--duration-normal) var(--ease-spring);
}
.s-modal-leave-active {
  transition:
    opacity 150ms var(--ease),
    transform 150ms var(--ease);
}
.s-modal-enter-from {
  opacity: 0;
  transform: translateY(-12px) scale(0.97);
}
.s-modal-leave-to {
  opacity: 0;
  transform: translateY(-6px) scale(0.98);
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 640px) {
  .s-wrap {
    align-items: flex-end;
    padding: 0;
  }

  .s-modal {
    max-width: 100%;
    max-height: 80vh;
    border-bottom-left-radius: 0;
    border-bottom-right-radius: 0;
    border-bottom: none;
  }

  .s-modal-enter-from {
    opacity: 0;
    transform: translateY(20px);
  }

  .s-modal-leave-to {
    opacity: 0;
    transform: translateY(20px);
  }

  .s-footer__shortcuts {
    gap: 0.6rem;
  }
}
</style>

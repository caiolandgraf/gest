<template>
  <div class="docs-page">
    <!-- Mobile sidebar toggle -->
    <button
      class="docs-sidebar-toggle"
      :class="{ open: sidebarOpen }"
      @click="sidebarOpen = !sidebarOpen"
      :aria-expanded="sidebarOpen"
      aria-controls="docs-sidebar"
    >
      <svg
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        aria-hidden="true"
      >
        <line x1="3" y1="6" x2="21" y2="6" />
        <line x1="3" y1="12" x2="21" y2="12" />
        <line x1="3" y1="18" x2="21" y2="18" />
      </svg>
      On this page
    </button>

    <div class="docs-layout container">
      <!-- ── Sidebar ──────────────────────────────────────────────────────── -->
      <aside
        id="docs-sidebar"
        class="docs-sidebar"
        :class="{ 'docs-sidebar--open': sidebarOpen }"
        aria-label="Documentation navigation"
      >
        <div class="docs-sidebar__inner">
          <div
            class="docs-sidebar__section"
            v-for="group in navGroups"
            :key="group.label"
          >
            <p class="docs-sidebar__group-label">{{ group.label }}</p>
            <ul class="docs-sidebar__list">
              <li v-for="item in group.items" :key="item.id">
                <a
                  :href="`#${item.id}`"
                  class="docs-sidebar__link"
                  :class="{
                    'docs-sidebar__link--active': activeSection === item.id
                  }"
                  @click.prevent="scrollTo(item.id)"
                >
                  {{ item.label }}
                </a>
              </li>
            </ul>
          </div>
        </div>
      </aside>

      <!-- Overlay for mobile -->
      <Transition name="overlay">
        <div
          v-if="sidebarOpen"
          class="docs-overlay"
          @click="sidebarOpen = false"
          aria-hidden="true"
        ></div>
      </Transition>

      <!-- ── Main content ──────────────────────────────────────────────────── -->
      <main class="docs-main" id="docs-content">
        <!-- Page header -->
        <div class="docs-header">
          <div class="docs-header__meta">
            <span class="badge badge--green">Docs</span>
            <span class="docs-header__sep" aria-hidden="true">/</span>
            <span class="docs-header__page">Getting Started</span>
          </div>
          <h1 class="docs-header__title">Documentation</h1>
          <p class="docs-header__desc">
            Everything you need to know to use <strong>gest</strong> — from
            installation to the full API reference.
          </p>
        </div>

        <!-- ── Installation ────────────────────────────────────────────────── -->
        <section class="docs-section" id="installation">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#installation"
                class="docs-anchor"
                @click.prevent="scrollTo('installation')"
                aria-label="Link to Installation"
                >#</a
              >
              Installation
            </h2>
          </div>
          <p>
            gest has <strong>zero external dependencies</strong>. Installing it
            is a single command and it will never pollute your
            <code>go.sum</code>.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">bash</span>
              <button
                class="code-block__copy"
                @click="copy('go get github.com/caiolandgraf/gest', 'install')"
              >
                {{ copied === 'install' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-bash">go get github.com/caiolandgraf/gest</code></pre>
          </div>
          <div class="docs-callout docs-callout--info">
            <span class="docs-callout__icon">💡</span>
            <div>
              <strong>Requires Go 1.18+.</strong> The core runner is stdlib
              only. <code>fsnotify</code> is added as the single external
              dependency to power watch mode.
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Project structure ───────────────────────────────────────────── -->
        <section class="docs-section" id="structure">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#structure"
                class="docs-anchor"
                @click.prevent="scrollTo('structure')"
                aria-label="Link to Project Structure"
                >#</a
              >
              Project Structure
            </h2>
          </div>
          <p>
            gest is designed to live alongside your production code, not in a
            separate test binary. Use <code>_spec.go</code> files instead of
            <code>_test.go</code> — they're picked up by <code>go run</code>.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">text</span>
            </div>
            <pre><code>my-project/
├── go.mod
├── main.go              ← just calls gest.RunRegistered()
├── calculator.go
├── calculator_spec.go   ← self-registers via init()
├── user.go
└── user_spec.go         ← self-registers via init()</code></pre>
          </div>
          <div class="docs-callout docs-callout--warning">
            <span class="docs-callout__icon">⚠️</span>
            <div>
              <strong
                >Why <code>_spec.go</code> instead of
                <code>_test.go</code>?</strong
              ><br />
              The <code>_test.go</code> suffix is reserved by the Go toolchain
              for <code>go test</code>. gest uses <code>go run</code>, so any
              other suffix works — <code>_spec.go</code> is the recommended
              convention.
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Basic Usage ─────────────────────────────────────────────────── -->
        <section class="docs-section" id="usage">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#usage"
                class="docs-anchor"
                @click.prevent="scrollTo('usage')"
                aria-label="Link to Basic Usage"
                >#</a
              >
              Basic Usage
            </h2>
          </div>
          <p>
            Your <code>main.go</code> stays simple forever — it never needs to
            import individual spec files. Each <code>_spec.go</code> registers
            itself via <code>init()</code>.
          </p>

          <h3 class="docs-subsection-title">main.go</h3>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">go</span>
              <button
                class="code-block__copy"
                @click="copy(mainGoCode, 'main')"
              >
                {{ copied === 'main' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-go">{{ mainGoCode }}</code></pre>
          </div>

          <h3 class="docs-subsection-title">calculator_spec.go</h3>
          <p>
            Each spec file self-registers via <code>init()</code>. The suite is
            created with <code>Describe</code>, test cases are added with
            <code>It</code>, and the whole suite is handed to the global runner
            with <code>Register</code>.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">go</span>
              <button class="code-block__copy" @click="copy(specCode, 'spec')">
                {{ copied === 'spec' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-go">{{ specCode }}</code></pre>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Running tests ───────────────────────────────────────────────── -->
        <section class="docs-section" id="running">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#running"
                class="docs-anchor"
                @click.prevent="scrollTo('running')"
                aria-label="Link to Running Tests"
                >#</a
              >
              Running Tests
            </h2>
          </div>
          <p>
            Use <code>go run .</code> from your project root. Pass
            <code>-c</code> or <code>--coverage</code> to also display the
            coverage table.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">bash</span>
              <button class="code-block__copy" @click="copy(runCode, 'run')">
                {{ copied === 'run' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-bash">{{ runCode }}</code></pre>
          </div>

          <h3 class="docs-subsection-title">Passing output</h3>
          <div class="screenshot-wrapper">
            <img
              :src="`${baseUrl}images/passing.png`"
              alt="gest passing output — two test suites, all green checkmarks, timing and summary"
              class="screenshot"
              width="1280"
              height="892"
              loading="lazy"
            />
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Watch Mode ─────────────────────────────────────────────────── -->
        <section class="docs-section" id="watch">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#watch"
                class="docs-anchor"
                @click.prevent="scrollTo('watch')"
                aria-label="Link to Watch Mode"
                >#</a
              >
              Watch Mode
            </h2>
          </div>
          <p>
            Pass <code>--watch</code> (or <code>-w</code>) to enter watch mode.
            gest compiles your project into a temporary binary in your OS temp
            directory, runs it immediately, then re-runs it automatically
            whenever any <code>.go</code> file changes.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">bash</span>
              <button
                class="code-block__copy"
                @click="copy(watchCode, 'watch')"
              >
                {{ copied === 'watch' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-bash">{{ watchCode }}</code></pre>
          </div>
          <div class="docs-callout docs-callout--info">
            <span class="docs-callout__icon">⚡</span>
            <div>
              <strong>Debounced re-runs.</strong> Rapid saves (e.g. auto-format
              on save) are collapsed into a single re-run via a 200 ms debounce
              — so you always see one clean result, never a burst.
            </div>
          </div>
          <div class="docs-callout docs-callout--tip">
            <span class="docs-callout__icon">🧹</span>
            <div>
              The terminal is <strong>fully cleared</strong> (including
              scrollback) before each re-run. The compiled binary is written to
              your OS temp directory and
              <strong>automatically removed</strong> when you press
              <kbd>Ctrl+C</kbd> — no artifacts are left in your project.
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Failure messages ────────────────────────────────────────────── -->
        <section class="docs-section" id="failures">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#failures"
                class="docs-anchor"
                @click.prevent="scrollTo('failures')"
                aria-label="Link to Failure Messages"
                >#</a
              >
              Failure Messages
            </h2>
          </div>
          <p>
            When a test fails, gest shows you exactly what went wrong — the
            matcher that failed, the <strong>expected</strong> value, the
            <strong>received</strong>
            value, and a code snippet from the exact failing assertion.
          </p>

          <div class="screenshot-wrapper">
            <img
              :src="`${baseUrl}images/failure.png`"
              alt="gest failure output — FAIL badge, Expected vs Received diff and inline code snippet"
              class="screenshot"
              width="1181"
              height="945"
              loading="lazy"
            />
          </div>

          <h3 class="docs-subsection-title">Panic recovery</h3>
          <p>
            If a test panics, gest catches it gracefully and reports the panic
            message rather than crashing the entire run.
          </p>
          <div class="terminal">
            <div class="terminal__bar">
              <span class="terminal__dot terminal__dot--red"></span>
              <span class="terminal__dot terminal__dot--yellow"></span>
              <span class="terminal__dot terminal__dot--green"></span>
              <span class="terminal__title">panic recovery</span>
            </div>
            <div class="terminal__body">
              <div v-for="(line, i) in panicLines" :key="i">
                <span v-html="line"></span>
              </div>
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Coverage ────────────────────────────────────────────────────── -->
        <section class="docs-section" id="coverage">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#coverage"
                class="docs-anchor"
                @click.prevent="scrollTo('coverage')"
                aria-label="Link to Coverage"
                >#</a
              >
              Coverage
            </h2>
          </div>
          <p>
            Pass <code>-c</code> or <code>--coverage</code> to display a
            per-suite pass-rate table with pip-style progress bars after the
            test run.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">bash</span>
            </div>
            <pre><code class="language-bash">go run . -c
go run . --coverage</code></pre>
          </div>

          <div class="screenshot-wrapper">
            <img
              :src="`${baseUrl}images/coverage.png`"
              alt="gest coverage table — per-suite pass rate with pip-style progress bars"
              class="screenshot"
              width="1280"
              height="892"
              loading="lazy"
            />
          </div>

          <div class="docs-callout docs-callout--info">
            <span class="docs-callout__icon">🎨</span>
            <div>
              <strong>Bar colors:</strong>
              <span style="color: #3fb950"> ● green</span> ≥ 80% ·
              <span style="color: #e3b341"> ● yellow</span> ≥ 50% ·
              <span style="color: #f85149"> ● red</span> &lt; 50%
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Matchers ────────────────────────────────────────────────────── -->
        <section class="docs-section" id="matchers">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#matchers"
                class="docs-anchor"
                @click.prevent="scrollTo('matchers')"
                aria-label="Link to Matchers"
                >#</a
              >
              Matchers
            </h2>
          </div>
          <p>
            Every assertion starts with <code>t.Expect(value)</code> which
            returns an <code>*Expectation</code>. Chain one of the matchers
            below to complete the assertion.
          </p>

          <div class="table-wrapper">
            <table>
              <thead>
                <tr>
                  <th>Matcher</th>
                  <th>Description</th>
                  <th>Example</th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="matcher in matchers" :key="matcher.name">
                  <td>
                    <code>{{ matcher.name }}</code>
                  </td>
                  <td>{{ matcher.desc }}</td>
                  <td>
                    <code class="example-code">{{ matcher.example }}</code>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <h3 class="docs-subsection-title">Usage examples</h3>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">go</span>
              <button
                class="code-block__copy"
                @click="copy(matchersExampleCode, 'matchers-ex')"
              >
                {{ copied === 'matchers-ex' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-go">{{ matchersExampleCode }}</code></pre>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Negation ────────────────────────────────────────────────────── -->
        <section class="docs-section" id="negation">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#negation"
                class="docs-anchor"
                @click.prevent="scrollTo('negation')"
                aria-label="Link to Negation"
                >#</a
              >
              Negation
            </h2>
          </div>
          <p>
            Any matcher can be negated by calling <code>.Not()</code> before it.
            <code>.Not()</code> returns the same <code>*Expectation</code>, so
            the chain reads naturally.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">go</span>
              <button
                class="code-block__copy"
                @click="copy(negationCode, 'negation')"
              >
                {{ copied === 'negation' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-go">{{ negationCode }}</code></pre>
          </div>

          <div class="docs-callout docs-callout--tip">
            <span class="docs-callout__icon">✨</span>
            <div>
              When a negated assertion fails, the error message automatically
              prefixes the matcher name with <code>not.</code> — e.g.
              <code>expect(received).not.toBeNil(expected)</code> — so failures
              are still crystal clear.
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Full API ────────────────────────────────────────────────────── -->
        <section class="docs-section" id="api">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#api"
                class="docs-anchor"
                @click.prevent="scrollTo('api')"
                aria-label="Link to Full API"
                >#</a
              >
              Full API Reference
            </h2>
          </div>
          <p>
            All exported symbols in
            <code>github.com/caiolandgraf/gest/gest</code>:
          </p>

          <div class="api-grid">
            <div
              class="api-card"
              v-for="fn in apiFunctions"
              :key="fn.signature"
            >
              <div class="api-card__sig">
                <code>{{ fn.signature }}</code>
              </div>
              <p class="api-card__desc">{{ fn.desc }}</p>
              <div v-if="fn.example" class="code-block api-card__example">
                <div class="code-block__header">
                  <span class="code-block__lang">go</span>
                </div>
                <pre><code class="language-go">{{ fn.example }}</code></pre>
              </div>
            </div>
          </div>

          <h3 class="docs-subsection-title">Full example wiring</h3>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">go</span>
              <button
                class="code-block__copy"
                @click="copy(fullApiCode, 'full-api')"
              >
                {{ copied === 'full-api' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-go">{{ fullApiCode }}</code></pre>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Philosophy ──────────────────────────────────────────────────── -->
        <section class="docs-section" id="philosophy">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#philosophy"
                class="docs-anchor"
                @click.prevent="scrollTo('philosophy')"
                aria-label="Link to Philosophy"
                >#</a
              >
              Philosophy
            </h2>
          </div>
          <p>
            gest is opinionated by design. These five principles guide every
            decision in the codebase:
          </p>
          <div class="philosophy-list">
            <div class="philosophy-item" v-for="p in philosophy" :key="p.title">
              <span class="philosophy-item__icon">{{ p.icon }}</span>
              <div>
                <h4 class="philosophy-item__title">{{ p.title }}</h4>
                <p class="philosophy-item__desc">{{ p.desc }}</p>
              </div>
            </div>
          </div>
        </section>

        <hr class="docs-divider" />

        <!-- ── Example project ─────────────────────────────────────────────── -->
        <section class="docs-section" id="example">
          <div class="docs-section__anchor-header">
            <h2 class="docs-section__title">
              <a
                href="#example"
                class="docs-anchor"
                @click.prevent="scrollTo('example')"
                aria-label="Link to Example Project"
                >#</a
              >
              Example Project
            </h2>
          </div>
          <p>
            A complete, runnable example is included in the repository under the
            <code>example/</code> directory.
          </p>
          <div class="code-block">
            <div class="code-block__header">
              <span class="code-block__lang">bash</span>
              <button
                class="code-block__copy"
                @click="copy(exampleRunCode, 'example-run')"
              >
                {{ copied === 'example-run' ? '✓ Copied' : 'Copy' }}
              </button>
            </div>
            <pre><code class="language-bash">{{ exampleRunCode }}</code></pre>
          </div>

          <div class="docs-nav-links">
            <RouterLink
              to="/contributors"
              class="docs-nav-link docs-nav-link--next"
            >
              <div>
                <span class="docs-nav-link__label">Next</span>
                <span class="docs-nav-link__title">Contributors</span>
              </div>
              <svg
                width="20"
                height="20"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                aria-hidden="true"
              >
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </RouterLink>
          </div>
        </section>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, nextTick } from 'vue'
const baseUrl = import.meta.env.BASE_URL

import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import 'highlight.js/styles/github-dark.css'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)

onMounted(async () => {
  await nextTick()
  hljs.highlightAll()
  setupIntersection()
})

onUnmounted(() => {
  if (observer) observer.disconnect()
})

/* ── Sidebar state ─────────────────────────────────────────────────────────── */
const sidebarOpen = ref(false)
const activeSection = ref('installation')

/* ── Copy state ────────────────────────────────────────────────────────────── */
const copied = ref(null)

function copy(text, key) {
  navigator.clipboard.writeText(text).then(() => {
    copied.value = key
    setTimeout(() => {
      copied.value = null
    }, 2000)
  })
}

/* ── Scroll to section ─────────────────────────────────────────────────────── */
function scrollTo(id) {
  sidebarOpen.value = false
  const el = document.getElementById(id)
  if (!el) return
  const offset = 88
  const top = el.getBoundingClientRect().top + window.scrollY - offset
  window.scrollTo({ top, behavior: 'smooth' })
  activeSection.value = id
}

/* ── Intersection Observer for active section ──────────────────────────────── */
let observer = null

function setupIntersection() {
  const allIds = navGroups.flatMap(g => g.items.map(i => i.id))
  const sections = allIds.map(id => document.getElementById(id)).filter(Boolean)

  observer = new IntersectionObserver(
    entries => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          activeSection.value = entry.target.id
          break
        }
      }
    },
    { rootMargin: '-20% 0px -70% 0px', threshold: 0 }
  )

  sections.forEach(s => observer.observe(s))
}

/* ── Nav groups ────────────────────────────────────────────────────────────── */
const navGroups = [
  {
    label: 'Getting Started',
    items: [
      { id: 'installation', label: 'Installation' },
      { id: 'structure', label: 'Project Structure' },
      { id: 'usage', label: 'Basic Usage' },
      { id: 'running', label: 'Running Tests' }
    ]
  },
  {
    label: 'Features',
    items: [
      { id: 'watch', label: 'Watch Mode' },
      { id: 'failures', label: 'Failure Messages' },
      { id: 'coverage', label: 'Coverage' },
      { id: 'matchers', label: 'Matchers' },
      { id: 'negation', label: 'Negation' }
    ]
  },
  {
    label: 'Reference',
    items: [
      { id: 'api', label: 'Full API' },
      { id: 'philosophy', label: 'Philosophy' },
      { id: 'example', label: 'Example Project' }
    ]
  }
]

/* ── Terminal lines ────────────────────────────────────────────────────────── */

const panicLines = [
  ``,
  `<span class="t-badge-fail"> FAIL </span>  <span class="t-bold t-white">math</span>`,
  `  <span class="t-red">✕</span> <span class="t-red t-bold">should not panic (0ms)</span>`,
  ``,
  `  <span class="t-red t-bold">● should not panic</span>`,
  ``,
  `    <span class="t-red">panic: runtime error: index out of range [3] with length 2</span>`,
  ``
]

/* ── Matchers data ─────────────────────────────────────────────────────────── */
const matchers = [
  {
    name: '.ToBe(v)',
    desc: 'Strict equality (==). Best for primitives and identical references.',
    example: 't.Expect(2+2).ToBe(float64(4))'
  },
  {
    name: '.ToEqual(v)',
    desc: 'Deep equality via reflect.DeepEqual. Best for structs, slices, maps.',
    example: 't.Expect(got).ToEqual(expected)'
  },
  {
    name: '.ToBeNil()',
    desc: 'Asserts the value is nil. Works with pointers and interfaces.',
    example: 't.Expect(err).ToBeNil()'
  },
  {
    name: '.ToBeTrue()',
    desc: 'Asserts the value is exactly true.',
    example: 't.Expect(ok).ToBeTrue()'
  },
  {
    name: '.ToBeFalse()',
    desc: 'Asserts the value is exactly false.',
    example: 't.Expect(done).ToBeFalse()'
  },
  {
    name: '.ToContain(s)',
    desc: 'Asserts a string contains the given substring.',
    example: 't.Expect(msg).ToContain("error")'
  },
  {
    name: '.ToHaveLength(n)',
    desc: 'Asserts the length of a string, slice, array, or map.',
    example: 't.Expect(items).ToHaveLength(3)'
  },
  {
    name: '.ToBeGreaterThan(n)',
    desc: 'Asserts a number is strictly greater than n.',
    example: 't.Expect(score).ToBeGreaterThan(0)'
  },
  {
    name: '.ToBeLessThan(n)',
    desc: 'Asserts a number is strictly less than n.',
    example: 't.Expect(err).ToBeLessThan(0.5)'
  },
  {
    name: '.ToBeCloseTo(n, d?)',
    desc: 'Asserts a float is approximately equal to n (default delta ±0.001).',
    example: 't.Expect(pi).ToBeCloseTo(3.14, 0.01)'
  }
]

/* ── API functions ─────────────────────────────────────────────────────────── */
const apiFunctions = [
  {
    signature: 'gest.Describe(name string) *Suite',
    desc: 'Creates a new test suite with the given name. Returns a *Suite you can add test cases to.',
    example: `s := gest.Describe("calculator")`
  },
  {
    signature: '(*Suite).It(name string, fn func(*T))',
    desc: 'Adds a test case to the suite. The function receives a *T used to make assertions.',
    example: `s.It("adds correctly", func(t *gest.T) {\n    t.Expect(add(1, 2)).ToBe(float64(3))\n})`
  },
  {
    signature: 'gest.Register(s *Suite)',
    desc: 'Adds a suite to the global registry. Call this at the end of init() in each spec file.',
    example: `gest.Register(s)`
  },
  {
    signature: 'gest.RunRegistered() bool',
    desc: 'Runs all suites registered via Register(). Returns true if all tests passed. Call this in main().',
    example: `func main() { gest.RunRegistered() }`
  },
  {
    signature: 'gest.RunAll(suites ...*Suite) bool',
    desc: 'Runs the provided suites manually, without the global registry. Useful for fine-grained control.',
    example: `gest.RunAll(suite1, suite2)`
  },
  {
    signature: '(*T).Expect(v interface{}) *Expectation',
    desc: 'Starts an assertion chain on the given value. Returns an *Expectation to chain a matcher on.',
    example: `t.Expect(result).ToBe(float64(42))`
  },
  {
    signature: '(*Expectation).Not() *Expectation',
    desc: 'Negates the next matcher in the chain. Returns the same *Expectation for chaining.',
    example: `t.Expect(err).Not().ToBeNil()`
  }
]

/* ── Philosophy ────────────────────────────────────────────────────────────── */
const philosophy = [
  {
    icon: '📦',
    title: 'Minimal dependencies',
    desc: 'The core runner is stdlib only. The single external dependency (fsnotify) exists solely to power watch mode.'
  },
  {
    icon: '⚙️',
    title: 'Zero config',
    desc: 'No config files, no separate CLI, no setup rituals. Drop in the file and run.'
  },
  {
    icon: '🔮',
    title: 'Auto-discovery',
    desc: 'Spec files self-register via init(). main.go never needs to import individual spec files.'
  },
  {
    icon: '🎨',
    title: 'Beautiful output',
    desc: 'Colors, code snippets, timing, and progress bars. Your testing experience should be joyful.'
  },
  {
    icon: '🤝',
    title: 'Familiar API',
    desc: 'If you know Jest, you already know gest. The mental model is identical.'
  }
]

/* ── Code samples ──────────────────────────────────────────────────────────── */
const mainGoCode = `package main

import "github.com/caiolandgraf/gest/gest"

func main() {
    gest.RunRegistered()
}`

const specCode = `package main

import "github.com/caiolandgraf/gest/gest"

func init() {
    calc := Calculator{}
    s := gest.Describe("calculator")

    s.It("adding 2 + 2 should return 4", func(t *gest.T) {
        t.Expect(calc.Add(2, 2)).ToBe(float64(4))
    })

    s.It("dividing by zero should return an error", func(t *gest.T) {
        _, err := calc.Divide(10, 0)
        t.Expect(err).Not().ToBeNil()
    })

    gest.Register(s)
}`

const runCode = `go run .           # run all tests
go run . -c        # run with coverage table
go run . --coverage`

const watchCode = `go run . --watch         # watch mode — re-runs on every .go change
go run . -w              # shorthand
go run . --watch -c      # watch + coverage`

const matchersExampleCode = `s.It("demonstrates all matchers", func(t *gest.T) {
    // Strict equality
    t.Expect(2 + 2).ToBe(float64(4))

    // Deep equality (structs, slices, maps)
    t.Expect([]int{1, 2, 3}).ToEqual([]int{1, 2, 3})

    // Nil check
    var err error
    t.Expect(err).ToBeNil()

    // Boolean
    t.Expect(true).ToBeTrue()
    t.Expect(false).ToBeFalse()

    // String contains
    t.Expect("hello world").ToContain("world")

    // Length
    t.Expect([]string{"a", "b"}).ToHaveLength(2)

    // Numeric comparisons
    t.Expect(float64(10)).ToBeGreaterThan(5)
    t.Expect(float64(3)).ToBeLessThan(10)

    // Float approximation (default delta: ±0.001)
    t.Expect(math.Pi).ToBeCloseTo(3.14159, 0.001)
})`

const negationCode = `s.It("uses negation", func(t *gest.T) {
    // Any matcher can be negated with .Not()
    _, err := calc.Divide(10, 0)
    t.Expect(err).Not().ToBeNil()

    t.Expect("gest").Not().ToContain("jest")

    result, _ := calc.Divide(10, 2)
    t.Expect(result).Not().ToBe(float64(0))

    t.Expect([]int{1, 2}).Not().ToHaveLength(5)
})`

const fullApiCode = `// ── main.go ─────────────────────────────────────────────────────────
package main

import "github.com/caiolandgraf/gest/gest"

func main() {
    // Option A: run everything registered via init()
    gest.RunRegistered()

    // Option B: manually specify which suites to run
    // gest.RunAll(suite1, suite2, suite3)
}

// ── calculator_spec.go ───────────────────────────────────────────────
package main

import "github.com/caiolandgraf/gest/gest"

func init() {
    calc := Calculator{}
    s := gest.Describe("calculator")

    s.It("adds two numbers", func(t *gest.T) {
        t.Expect(calc.Add(1, 2)).ToBe(float64(3))
    })

    s.It("divides correctly", func(t *gest.T) {
        result, err := calc.Divide(10, 2)
        t.Expect(err).ToBeNil()
        t.Expect(result).ToBe(float64(5))
    })

    s.It("errors on divide by zero", func(t *gest.T) {
        _, err := calc.Divide(10, 0)
        t.Expect(err).Not().ToBeNil()
    })

    gest.Register(s)
}`

const exampleRunCode = `# Clone the repository
git clone https://github.com/caiolandgraf/gest.git
cd gest/example

# Run the example tests
go run .

# Run with coverage table
go run . -c`
</script>

<style scoped>
/* ── Layout ─────────────────────────────────────────────────────────────────── */
.docs-page {
  min-height: calc(100vh - var(--nav-height));
}

.docs-layout {
  display: grid;
  grid-template-columns: var(--docs-sidebar) 1fr;
  gap: 0;
  align-items: start;
  padding: 0 1.5rem;
}

/* ── Sidebar toggle (mobile) ────────────────────────────────────────────────── */
.docs-sidebar-toggle {
  display: none;
  align-items: center;
  gap: 8px;
  width: 100%;
  padding: 0.75rem 1.5rem;
  background: var(--bg-2);
  border: none;
  border-bottom: 1px solid var(--border-subtle);
  color: var(--text-2);
  font-size: 0.875rem;
  font-weight: 500;
  font-family: var(--font-sans);
  cursor: pointer;
  transition: background var(--duration-fast) var(--ease);
}

.docs-sidebar-toggle:hover {
  background: var(--bg-3);
}

.docs-sidebar-toggle.open {
  color: var(--green);
  background: var(--green-glow);
}

/* ── Sidebar ────────────────────────────────────────────────────────────────── */
.docs-sidebar {
  position: sticky;
  top: calc(var(--nav-height) + 1rem);
  max-height: calc(100vh - var(--nav-height) - 2rem);
  overflow-y: auto;
  padding: 2rem 0 2rem 0;
  border-right: 1px solid var(--border-subtle);
}

.docs-sidebar__inner {
  padding-right: 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

.docs-sidebar__group-label {
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-4);
  margin-bottom: 0.5rem;
}

.docs-sidebar__list {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.docs-sidebar__link {
  display: block;
  padding: 0.35rem 0.75rem;
  border-radius: var(--radius-sm);
  font-size: 0.875rem;
  color: var(--text-3);
  text-decoration: none;
  border-left: 2px solid transparent;
  transition: all var(--duration-fast) var(--ease);
}

.docs-sidebar__link:hover {
  color: var(--text-1);
  background: var(--bg-3);
}

.docs-sidebar__link--active {
  color: var(--green) !important;
  border-left-color: var(--green);
  background: var(--green-glow);
}

/* ── Overlay ────────────────────────────────────────────────────────────────── */
.docs-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  z-index: 49;
  backdrop-filter: blur(2px);
}

.overlay-enter-active,
.overlay-leave-active {
  transition: opacity var(--duration-normal) var(--ease);
}

.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}

/* ── Main content ───────────────────────────────────────────────────────────── */
.docs-main {
  padding: 2.5rem 0 4rem 3rem;
  min-width: 0;
  max-width: 820px;
}

/* ── Page header ────────────────────────────────────────────────────────────── */
.docs-header {
  padding-bottom: 2.5rem;
  margin-bottom: 2.5rem;
  border-bottom: 1px solid var(--border-subtle);
}

.docs-header__meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 1.25rem;
}

.docs-header__sep {
  color: var(--text-4);
  font-size: 0.85rem;
}

.docs-header__page {
  font-size: 0.85rem;
  color: var(--text-3);
}

.docs-header__title {
  font-size: clamp(2rem, 4vw, 2.75rem);
  margin-bottom: 0.75rem;
}

.docs-header__desc {
  font-size: 1.05rem;
  color: var(--text-2);
  line-height: 1.7;
  max-width: 600px;
}

/* ── Section ────────────────────────────────────────────────────────────────── */
.docs-section {
  padding: 2.5rem 0;
  scroll-margin-top: 88px;
}

.docs-section__anchor-header {
  margin-bottom: 1rem;
}

.docs-section__title {
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--text-1);
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.docs-anchor {
  color: var(--text-4);
  font-size: 1rem;
  text-decoration: none;
  opacity: 0;
  transition: opacity var(--duration-fast) var(--ease);
  flex-shrink: 0;
}

.docs-section__title:hover .docs-anchor {
  opacity: 1;
}

.docs-anchor:hover {
  color: var(--green);
}

.docs-section > p {
  margin-bottom: 1.25rem;
  font-size: 0.975rem;
  line-height: 1.8;
}

.docs-subsection-title {
  font-size: 1.05rem;
  font-weight: 600;
  color: var(--text-1);
  margin: 1.75rem 0 0.75rem;
}

.docs-divider {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, var(--border-subtle), transparent);
  margin: 0;
}

/* ── Callout ────────────────────────────────────────────────────────────────── */
.docs-callout {
  display: flex;
  gap: 12px;
  padding: 1rem 1.25rem;
  border-radius: var(--radius-md);
  border-left: 3px solid;
  margin-top: 1.25rem;
  font-size: 0.9rem;
  line-height: 1.7;
}

.docs-callout--info {
  background: rgba(88, 166, 255, 0.07);
  border-color: #58a6ff;
  color: var(--text-2);
}

.docs-callout--warning {
  background: rgba(227, 179, 65, 0.07);
  border-color: #e3b341;
  color: var(--text-2);
}

.docs-callout--tip {
  background: var(--green-glow);
  border-color: var(--green);
  color: var(--text-2);
}

.docs-callout__icon {
  font-size: 1.1rem;
  flex-shrink: 0;
  margin-top: 1px;
}

/* ── Table enhancements ─────────────────────────────────────────────────────── */
.example-code {
  font-size: 0.78rem;
  white-space: nowrap;
  color: var(--text-3);
}

/* ── API cards ──────────────────────────────────────────────────────────────── */
.api-grid {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-bottom: 1.75rem;
}

.api-card {
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  padding: 1.25rem 1.5rem;
  transition: border-color var(--duration-normal) var(--ease);
}

.api-card:hover {
  border-color: var(--border);
}

.api-card__sig {
  margin-bottom: 0.5rem;
}

.api-card__sig code {
  font-size: 0.875rem;
  color: var(--green-bright);
  background: transparent;
  border: none;
  padding: 0;
}

.api-card__desc {
  font-size: 0.875rem;
  color: var(--text-3);
  margin-bottom: 0.75rem;
  line-height: 1.65;
}

.api-card__example {
  margin-top: 0.75rem;
}

.api-card__example .code-block__header {
  padding: 0.4rem 0.85rem;
}

.api-card__example pre {
  padding: 0.85rem 1rem;
}

/* ── Philosophy list ────────────────────────────────────────────────────────── */
.philosophy-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  margin-top: 1.25rem;
}

.philosophy-item {
  display: flex;
  gap: 1rem;
  padding: 1.25rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition:
    border-color var(--duration-normal) var(--ease),
    background var(--duration-normal) var(--ease);
}

.philosophy-item:hover {
  border-color: var(--border);
  background: var(--bg-3);
}

.philosophy-item__icon {
  font-size: 1.5rem;
  line-height: 1;
  flex-shrink: 0;
  margin-top: 2px;
}

.philosophy-item__title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-1);
  margin-bottom: 0.3rem;
}

.philosophy-item__desc {
  font-size: 0.875rem;
  color: var(--text-3);
  line-height: 1.65;
}

/* ── Prev / Next navigation ─────────────────────────────────────────────────── */
.docs-nav-links {
  display: flex;
  gap: 1rem;
  margin-top: 3rem;
  justify-content: flex-end;
}

.docs-nav-link {
  flex: 1;
  max-width: 280px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem 1.25rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: all var(--duration-normal) var(--ease);
}

.docs-nav-link:hover {
  border-color: var(--green-dim);
  background: var(--green-glow);
}

.docs-nav-link--next {
  text-align: right;
}

.docs-nav-link__label {
  display: block;
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--green);
  margin-bottom: 2px;
}

.docs-nav-link__title {
  display: block;
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-1);
}

.docs-nav-link svg {
  color: var(--text-3);
  flex-shrink: 0;
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 900px) {
  .docs-layout {
    grid-template-columns: 1fr;
  }

  .docs-sidebar-toggle {
    display: flex;
  }

  .docs-sidebar {
    position: fixed;
    top: var(--nav-height);
    left: 0;
    bottom: 0;
    width: min(var(--docs-sidebar), 85vw);
    background: var(--bg-1);
    border-right: 1px solid var(--border);
    z-index: 50;
    transform: translateX(-100%);
    transition: transform var(--duration-normal) var(--ease);
    max-height: none;
    padding: 1.5rem 0;
  }

  .docs-sidebar--open {
    transform: translateX(0);
  }

  .docs-sidebar__inner {
    padding: 0 1.25rem;
  }

  .docs-main {
    padding: 2rem 0 3rem;
  }
}

@media (max-width: 600px) {
  .docs-header__title {
    font-size: 1.75rem;
  }

  .docs-section__title {
    font-size: 1.25rem;
  }

  .table-wrapper {
    font-size: 0.82rem;
  }

  .example-code {
    display: none;
  }
}
</style>

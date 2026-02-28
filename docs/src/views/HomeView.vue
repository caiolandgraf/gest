<template>
  <div class="home">
    <!-- ── Hero ──────────────────────────────────────────────────────────── -->
    <section class="hero">
      <div class="hero__bg" aria-hidden="true">
        <div class="hero__grid"></div>
        <div class="hero__radial"></div>
      </div>

      <div class="hero__inner container">
        <div class="hero__content">
          <div
            class="hero__badge-row animate-fade-up"
            style="animation-delay: 0ms"
          >
            <span class="badge badge--green badge--dot"
              >v1.1.0 — Now available</span
            >
            <a
              href="https://github.com/caiolandgraf/gest"
              target="_blank"
              rel="noopener noreferrer"
              class="hero__gh-stars"
            >
              <svg
                width="13"
                height="13"
                viewBox="0 0 24 24"
                fill="currentColor"
                aria-hidden="true"
              >
                <path
                  d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                />
              </svg>
              Star on GitHub
              <svg
                width="11"
                height="11"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
                aria-hidden="true"
              >
                <path d="M7 17L17 7M17 7H7M17 7v10" />
              </svg>
            </a>
          </div>

          <h1 class="hero__title animate-fade-up" style="animation-delay: 60ms">
            Testing in Go,<br />
            <span class="gradient-text">the Jest way.</span>
          </h1>

          <p
            class="hero__description animate-fade-up"
            style="animation-delay: 120ms"
          >
            <strong>gest</strong> 🧪 brings Jest's legendary developer
            experience to Go — colored output, descriptive failure messages with
            code snippets, a fluent assertion API, and a beautiful coverage
            table. All in a single file with
            <strong>minimal dependencies</strong>.
          </p>

          <div
            class="hero__ctas animate-fade-up"
            style="animation-delay: 180ms"
          >
            <RouterLink to="/docs" class="btn btn--primary hero__cta-main">
              <svg
                width="16"
                height="16"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                aria-hidden="true"
              >
                <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                <path
                  d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"
                />
              </svg>
              Get Started
            </RouterLink>
            <a
              href="https://github.com/caiolandgraf/gest"
              target="_blank"
              rel="noopener noreferrer"
              class="btn btn--ghost"
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

          <div
            class="hero__install animate-fade-up"
            style="animation-delay: 240ms"
          >
            <div class="hero__install-inner" @click="copyInstall">
              <span class="hero__install-prompt">$</span>
              <code class="hero__install-cmd"
                >go get github.com/caiolandgraf/gest</code
              >
              <span
                class="hero__install-copy"
                :class="{ copied: installCopied }"
              >
                <svg
                  v-if="!installCopied"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2"
                  aria-hidden="true"
                >
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2" />
                  <path
                    d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"
                  />
                </svg>
                <svg
                  v-else
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  stroke-width="2.5"
                  aria-hidden="true"
                >
                  <polyline points="20 6 9 17 4 12" />
                </svg>
                {{ installCopied ? 'Copied!' : 'Copy' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Terminal demo -->
        <div
          class="hero__terminal animate-fade-up"
          style="animation-delay: 300ms"
        >
          <div class="screenshot-wrapper screenshot-wrapper--float">
            <img
              :src="`${baseUrl}images/passing.png`"
              alt="gest output — two test suites passing with colored checkmarks and timing"
              class="screenshot"
              width="1280"
              height="892"
              loading="eager"
            />
          </div>
        </div>
      </div>

      <!-- Scroll indicator -->
      <div class="hero__scroll" aria-hidden="true">
        <div class="hero__scroll-dot"></div>
      </div>
    </section>

    <!-- ── Stats strip ────────────────────────────────────────────────────── -->
    <section class="stats" aria-label="Project statistics">
      <div class="stats__inner container">
        <div class="stats__item" v-for="stat in stats" :key="stat.label">
          <span class="stats__value">{{ stat.value }}</span>
          <span class="stats__label">{{ stat.label }}</span>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Features ──────────────────────────────────────────────────────── -->
    <section class="section features" aria-labelledby="features-heading">
      <div class="container">
        <div class="features__header">
          <p class="section-label">Why gest?</p>
          <h2 id="features-heading">
            Everything you love about Jest,<br />now in Go.
          </h2>
          <p class="features__subtitle">
            Built for developers who care about their testing experience. No
            setup, no config files, no compromises.
          </p>
        </div>

        <div class="features__grid">
          <div
            v-for="feature in features"
            :key="feature.title"
            class="feature-card"
          >
            <div
              class="feature-card__icon"
              :style="{ background: feature.iconBg }"
            >
              {{ feature.icon }}
            </div>
            <h3 class="feature-card__title">{{ feature.title }}</h3>
            <p class="feature-card__desc">{{ feature.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Matchers preview ───────────────────────────────────────────────── -->
    <section
      class="section matchers-section"
      aria-labelledby="matchers-heading"
    >
      <div class="container">
        <div class="matchers-section__layout">
          <div class="matchers-section__text">
            <p class="section-label">Assertions</p>
            <h2 id="matchers-heading">A fluent API you already know.</h2>
            <p>
              If you've used Jest before, gest's assertion API will feel
              instantly familiar. Chain matchers, negate with
              <code>.Not()</code>, and get clear failure messages with inline
              code snippets when things go wrong.
            </p>
            <RouterLink
              to="/docs#matchers"
              class="btn btn--ghost matchers-section__cta"
            >
              View all matchers
              <svg
                width="14"
                height="14"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2.5"
                aria-hidden="true"
              >
                <path d="M5 12h14M12 5l7 7-7 7" />
              </svg>
            </RouterLink>
          </div>

          <div class="matchers-section__code">
            <div class="code-block">
              <div class="code-block__header">
                <span class="code-block__lang">Go</span>
                <button
                  class="code-block__copy"
                  @click="copyCode(matcherCode, 'matchers')"
                  :class="{ copied: copiedKey === 'matchers' }"
                >
                  {{ copiedKey === 'matchers' ? '✓ Copied' : 'Copy' }}
                </button>
              </div>
              <pre><code class="language-go">{{ matcherCode }}</code></pre>
            </div>
          </div>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Failure messages preview ──────────────────────────────────────── -->
    <section class="section failure-section" aria-labelledby="failure-heading">
      <div class="container">
        <div class="failure-section__layout">
          <div class="failure-section__terminal">
            <div class="screenshot-wrapper screenshot-wrapper--float">
              <img
                :src="`${baseUrl}images/failure.png`"
                alt="gest failure output — showing Expected vs Received diff and inline code snippet"
                class="screenshot"
                width="1181"
                height="945"
                loading="lazy"
              />
            </div>
          </div>

          <div class="failure-section__text">
            <p class="section-label">Failure messages</p>
            <h2 id="failure-heading">Know exactly what broke, and why.</h2>
            <p>
              When a test fails, gest pinpoints the problem immediately —
              showing you what was <strong>expected</strong>, what was
              <strong>received</strong>, and even a code snippet from the exact
              line where the assertion failed.
            </p>
            <ul class="failure-section__list">
              <li v-for="item in failureFeatures" :key="item">
                <span class="failure-section__check">✓</span>
                {{ item }}
              </li>
            </ul>
          </div>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Quick start ────────────────────────────────────────────────────── -->
    <section class="section quickstart" aria-labelledby="quickstart-heading">
      <div class="container">
        <div class="quickstart__header">
          <p class="section-label">Quick start</p>
          <h2 id="quickstart-heading">Up and running in 60 seconds.</h2>
          <p class="quickstart__subtitle">
            No config files. No CLI to install. Just write tests.
          </p>
        </div>

        <div class="quickstart__steps">
          <div
            class="quickstart__step"
            v-for="(step, i) in quickstartSteps"
            :key="i"
          >
            <div class="quickstart__step-header">
              <span class="quickstart__step-num">{{ i + 1 }}</span>
              <h3 class="quickstart__step-title">{{ step.title }}</h3>
            </div>
            <div class="code-block">
              <div class="code-block__header">
                <span class="code-block__lang">{{ step.lang }}</span>
                <button
                  class="code-block__copy"
                  @click="copyCode(step.code, `step-${i}`)"
                  :class="{ copied: copiedKey === `step-${i}` }"
                >
                  {{ copiedKey === `step-${i}` ? '✓ Copied' : 'Copy' }}
                </button>
              </div>
              <pre><code :class="`language-${step.lang.toLowerCase()}`">{{ step.code }}</code></pre>
            </div>
          </div>
        </div>

        <div class="quickstart__cta">
          <RouterLink to="/docs" class="btn btn--primary">
            Read the full documentation
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2.5"
              aria-hidden="true"
            >
              <path d="M5 12h14M12 5l7 7-7 7" />
            </svg>
          </RouterLink>
          <a
            href="https://github.com/caiolandgraf/gest/tree/main/example"
            target="_blank"
            rel="noopener noreferrer"
            class="btn btn--ghost"
          >
            View examples on GitHub
          </a>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Philosophy strip ───────────────────────────────────────────────── -->
    <section class="section philosophy" aria-labelledby="philosophy-heading">
      <div class="container">
        <div class="philosophy__header">
          <p class="section-label">Philosophy</p>
          <h2 id="philosophy-heading">Designed with purpose.</h2>
        </div>

        <div class="philosophy__grid">
          <div
            v-for="pillar in pillars"
            :key="pillar.title"
            class="philosophy__item"
          >
            <span class="philosophy__item-icon">{{ pillar.icon }}</span>
            <div>
              <h4 class="philosophy__item-title">{{ pillar.title }}</h4>
              <p class="philosophy__item-desc">{{ pillar.desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <!-- ── CTA banner ─────────────────────────────────────────────────────── -->
    <section class="cta-banner" aria-labelledby="cta-heading">
      <div class="cta-banner__glow" aria-hidden="true"></div>
      <div class="cta-banner__inner container">
        <h2 id="cta-heading" class="cta-banner__title">
          Ready to write better Go tests?
        </h2>
        <p class="cta-banner__subtitle">
          Join the developers who've brought the Jest DX to Go.
        </p>
        <div class="cta-banner__actions">
          <RouterLink to="/docs" class="btn btn--primary">
            Get Started — it's free
          </RouterLink>
          <RouterLink to="/contributors" class="btn btn--ghost">
            Meet the contributors
          </RouterLink>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const baseUrl = import.meta.env.BASE_URL

import hljs from 'highlight.js/lib/core'
import go from 'highlight.js/lib/languages/go'
import bash from 'highlight.js/lib/languages/bash'
import 'highlight.js/styles/github-dark.css'

hljs.registerLanguage('go', go)
hljs.registerLanguage('bash', bash)

onMounted(() => {
  hljs.highlightAll()
})

/* ── Copy helpers ──────────────────────────────────────────────────────────── */
const installCopied = ref(false)
const copiedKey = ref(null)

function copyInstall() {
  navigator.clipboard
    .writeText('go get github.com/caiolandgraf/gest')
    .then(() => {
      installCopied.value = true
      setTimeout(() => {
        installCopied.value = false
      }, 2000)
    })
}

function copyCode(text, key) {
  navigator.clipboard.writeText(text).then(() => {
    copiedKey.value = key
    setTimeout(() => {
      copiedKey.value = null
    }, 2000)
  })
}

/* ── Stats ─────────────────────────────────────────────────────────────────── */
const stats = [
  { value: '10', label: 'Built-in matchers' },
  { value: '0', label: 'External dependencies' },
  { value: '1', label: 'File to add' },
  { value: '∞', label: 'Tests you can write' }
]

/* ── Features ──────────────────────────────────────────────────────────────── */
const features = [
  {
    icon: '🎨',
    iconBg: 'rgba(63, 185, 80, 0.1)',
    title: 'Beautiful output',
    desc: 'Colorized, structured terminal output inspired by Jest. Pass/fail badges, test names, timing — all perfectly formatted.'
  },
  {
    icon: '📦',
    iconBg: 'rgba(88, 166, 255, 0.1)',
    title: 'Minimal dependencies',
    desc: 'Core runner is stdlib only. fsnotify is the single external dep, used exclusively for watch mode.'
  },
  {
    icon: '🔗',
    iconBg: 'rgba(188, 140, 255, 0.1)',
    title: 'Fluent assertion API',
    desc: 'Chain matchers naturally: t.Expect(val).ToBe(...), .ToEqual(...), .Not().ToBeNil() — reads like plain English.'
  },
  {
    icon: '🔍',
    iconBg: 'rgba(227, 179, 65, 0.1)',
    title: 'Inline code snippets',
    desc: 'Failed assertions print the exact source line that triggered them, with context lines around it — no guessing.'
  },
  {
    icon: '📊',
    iconBg: 'rgba(63, 185, 80, 0.1)',
    title: 'Coverage table',
    desc: 'Run with -c to get a per-suite pass-rate table with pip-style progress bars and color-coded thresholds.'
  },
  {
    icon: '⚡',
    iconBg: 'rgba(248, 81, 73, 0.1)',
    title: 'Auto-discovery via init()',
    desc: 'Add a _spec.go file, call gest.Register(s) in init(), and it just runs. main.go never needs to know what changed.'
  }
]

/* ── Matcher code sample ───────────────────────────────────────────────────── */
const matcherCode = `s := gest.Describe("calculator")

s.It("adds correctly", func(t *gest.T) {
  t.Expect(calc.Add(2, 2)).ToBe(float64(4))
})

s.It("does not divide by zero", func(t *gest.T) {
  _, err := calc.Divide(10, 0)
  t.Expect(err).Not().ToBeNil()
})

s.It("result is within range", func(t *gest.T) {
  r := calc.Add(1.0, 1.0)
  t.Expect(r).ToBeGreaterThan(1.5)
  t.Expect(r).ToBeLessThan(2.5)
  t.Expect(r).ToBeCloseTo(2.0, 0.01)
})`

/* ── Failure features list ─────────────────────────────────────────────────── */
const failureFeatures = [
  'Shows expected vs received values side-by-side',
  'Highlights the exact failing assertion in your code',
  'Displays surrounding lines for context',
  'Reports the file path and line number',
  'Catches panics and reports them gracefully'
]

/* ── Quick start steps ─────────────────────────────────────────────────────── */
const quickstartSteps = [
  {
    title: 'Install gest',
    lang: 'bash',
    code: `go get github.com/caiolandgraf/gest`
  },
  {
    title: 'Create your main.go',
    lang: 'Go',
    code: `package main

import "github.com/caiolandgraf/gest/gest"

func main() {
    gest.RunRegistered()
}`
  },
  {
    title: 'Write a spec file',
    lang: 'Go',
    code: `package main

import "github.com/caiolandgraf/gest/gest"

func init() {
    s := gest.Describe("calculator")

    s.It("adds 2 + 2", func(t *gest.T) {
        t.Expect(add(2, 2)).ToBe(float64(4))
    })

    s.It("does not divide by zero", func(t *gest.T) {
        _, err := divide(10, 0)
        t.Expect(err).Not().ToBeNil()
    })

    gest.Register(s)
}`
  },
  {
    title: 'Run your tests',
    lang: 'bash',
    code: `go run .           # run all tests
go run . -c        # run with coverage table
go run . --watch   # watch mode — re-run on every .go change`
  }
]

/* ── Philosophy pillars ────────────────────────────────────────────────────── */
const pillars = [
  {
    icon: '📦',
    title: 'Minimal dependencies',
    desc: 'Core runner is stdlib only. The single external dep (fsnotify) powers watch mode and nothing else.'
  },
  {
    icon: '⚙️',
    title: 'Zero config',
    desc: 'No config files, no separate CLI, no flags to learn. It just works.'
  },
  {
    icon: '🔮',
    title: 'Auto-discovery',
    desc: 'Spec files self-register via init(). Adding a test is as simple as adding a file.'
  },
  {
    icon: '🎨',
    title: 'Beautiful output',
    desc: 'Colors, code snippets, and progress bars. Testing should be a pleasure.'
  },
  {
    icon: '🤝',
    title: 'Familiar API',
    desc: 'If you know Jest, you already know gest. The learning curve is a speed bump.'
  }
]
</script>

<style scoped>
/* ── Hero ───────────────────────────────────────────────────────────────────── */
.hero {
  position: relative;
  min-height: calc(100vh - var(--nav-height));
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.hero__bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  user-select: none;
}

.hero__grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(var(--border-muted) 1px, transparent 1px),
    linear-gradient(90deg, var(--border-muted) 1px, transparent 1px);
  background-size: 48px 48px;
  opacity: 0.35;
  mask-image: radial-gradient(
    ellipse 80% 60% at 50% 0%,
    black 40%,
    transparent 100%
  );
  -webkit-mask-image: radial-gradient(
    ellipse 80% 60% at 50% 0%,
    black 40%,
    transparent 100%
  );
}

.hero__radial {
  position: absolute;
  top: -20%;
  left: 50%;
  transform: translateX(-50%);
  width: 900px;
  height: 600px;
  background: radial-gradient(
    ellipse at top,
    rgba(63, 185, 80, 0.08) 0%,
    transparent 60%
  );
  pointer-events: none;
}

/* ── Hero inner ─────────────────────────────────────────────────────────────── */
.hero__inner {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: center;
  padding: 5rem 1.5rem 4rem;
  flex: 1;
}

/* ── Hero content ───────────────────────────────────────────────────────────── */
.hero__badge-row {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 10px;
  margin-bottom: 1.5rem;
}

.hero__gh-stars {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 0.78rem;
  font-weight: 500;
  color: var(--text-3);
  padding: 0.25rem 0.65rem;
  border-radius: var(--radius-full);
  border: 1px solid var(--border-subtle);
  background: var(--bg-3);
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease);
}

.hero__gh-stars:hover {
  color: var(--text-1);
  border-color: var(--border);
  background: var(--bg-4);
}

.hero__title {
  margin-bottom: 1.25rem;
  line-height: 1.1;
}

.hero__description {
  font-size: 1.05rem;
  color: var(--text-2);
  line-height: 1.75;
  margin-bottom: 2rem;
  max-width: 520px;
}

.hero__ctas {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-bottom: 2rem;
}

.hero__cta-main {
  font-size: 0.95rem;
  padding: 0.75rem 1.6rem;
}

/* ── Hero install bar ───────────────────────────────────────────────────────── */
.hero__install {
  display: inline-flex;
}

.hero__install-inner {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0.6rem 1rem;
  background: var(--bg-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition:
    border-color var(--duration-fast) var(--ease),
    background var(--duration-fast) var(--ease);
  max-width: 100%;
}

.hero__install-inner:hover {
  border-color: var(--text-4);
  background: var(--bg-3);
}

.hero__install-prompt {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  color: var(--green);
  flex-shrink: 0;
  user-select: none;
}

.hero__install-cmd {
  font-family: var(--font-mono);
  font-size: 0.85rem;
  color: var(--text-2);
  background: none;
  border: none;
  padding: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.hero__install-copy {
  display: flex;
  align-items: center;
  gap: 5px;
  font-size: 0.75rem;
  font-family: var(--font-sans);
  color: var(--text-3);
  flex-shrink: 0;
  margin-left: 4px;
  transition: color var(--duration-fast) var(--ease);
  user-select: none;
}

.hero__install-copy.copied {
  color: var(--green);
}

/* ── Hero terminal ──────────────────────────────────────────────────────────── */
.hero__terminal {
  width: 100%;
}

/* ── Terminal lines animation ───────────────────────────────────────────────── */
.terminal__line {
  animation: fadeInUp 0.35s var(--ease) both;
  min-height: 1.4em;
}

/* ── Scroll indicator ───────────────────────────────────────────────────────── */
.hero__scroll {
  display: flex;
  justify-content: center;
  padding-bottom: 2rem;
}

.hero__scroll-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--text-4);
  animation: pulse 2.5s ease-in-out infinite;
}

/* ── Stats ──────────────────────────────────────────────────────────────────── */
.stats {
  background: var(--bg-0);
  border-top: 1px solid var(--border-subtle);
  border-bottom: 1px solid var(--border-subtle);
}

.stats__inner {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  padding: 0 1.5rem;
}

.stats__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  padding: 2rem 1rem;
  border-right: 1px solid var(--border-subtle);
  transition: background var(--duration-fast) var(--ease);
}

.stats__item:last-child {
  border-right: none;
}

.stats__item:hover {
  background: var(--bg-2);
}

.stats__value {
  font-family: var(--font-mono);
  font-size: 2rem;
  font-weight: 700;
  color: var(--green);
  line-height: 1;
  letter-spacing: -0.03em;
}

.stats__label {
  font-size: 0.78rem;
  color: var(--text-3);
  text-align: center;
  letter-spacing: 0.01em;
}

/* ── Features ───────────────────────────────────────────────────────────────── */
.features__header {
  text-align: center;
  max-width: 640px;
  margin: 0 auto 3.5rem;
}

.features__subtitle {
  font-size: 1.05rem;
  margin-top: 1rem;
  color: var(--text-3);
}

.features__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.25rem;
}

.feature-card {
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  padding: 1.75rem;
  transition:
    border-color var(--duration-normal) var(--ease),
    box-shadow var(--duration-normal) var(--ease),
    transform var(--duration-normal) var(--ease);
}

.feature-card:hover {
  border-color: var(--border);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.feature-card__icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.4rem;
  margin-bottom: 1rem;
}

.feature-card__title {
  font-size: 1rem;
  font-weight: 700;
  color: var(--text-1);
  margin-bottom: 0.5rem;
}

.feature-card__desc {
  font-size: 0.875rem;
  color: var(--text-3);
  line-height: 1.65;
}

/* ── Matchers section ───────────────────────────────────────────────────────── */
.matchers-section__layout {
  display: grid;
  grid-template-columns: 1fr 1.4fr;
  gap: 4rem;
  align-items: center;
}

.matchers-section__text h2 {
  margin-bottom: 1rem;
}

.matchers-section__text p {
  margin-bottom: 1.5rem;
  font-size: 1rem;
}

.matchers-section__cta {
  display: inline-flex;
  margin-top: 0.5rem;
}

.matchers-section__code {
  width: 100%;
}

/* ── Failure section ────────────────────────────────────────────────────────── */
.failure-section__layout {
  display: grid;
  grid-template-columns: 1.3fr 1fr;
  gap: 4rem;
  align-items: center;
}

.failure-section__text h2 {
  margin-bottom: 1rem;
}

.failure-section__text p {
  margin-bottom: 1.5rem;
  font-size: 1rem;
}

.failure-section__list {
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
  list-style: none;
}

.failure-section__list li {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: 0.9rem;
  color: var(--text-2);
}

.failure-section__check {
  color: var(--green);
  font-weight: 700;
  flex-shrink: 0;
  margin-top: 1px;
}

/* ── Quick start ────────────────────────────────────────────────────────────── */
.quickstart__header {
  text-align: center;
  max-width: 560px;
  margin: 0 auto 3rem;
}

.quickstart__subtitle {
  font-size: 1.05rem;
  margin-top: 1rem;
  color: var(--text-3);
}

.quickstart__steps {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.quickstart__step {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.quickstart__step-header {
  display: flex;
  align-items: center;
  gap: 10px;
}

.quickstart__step-num {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: var(--green-glow);
  border: 1px solid rgba(63, 185, 80, 0.25);
  color: var(--green);
  font-size: 0.8rem;
  font-weight: 700;
  font-family: var(--font-mono);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.quickstart__step-title {
  font-size: 0.95rem;
  font-weight: 600;
  color: var(--text-1);
}

.quickstart__cta {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 1rem;
}

/* ── Philosophy ─────────────────────────────────────────────────────────────── */
.philosophy__header {
  text-align: center;
  margin-bottom: 3rem;
}

.philosophy__grid {
  display: grid;
  grid-template-columns: repeat(5, 1fr);
  gap: 1rem;
}

.philosophy__item {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 0.75rem;
  padding: 1.5rem 1rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition:
    border-color var(--duration-normal) var(--ease),
    background var(--duration-normal) var(--ease);
}

.philosophy__item:hover {
  border-color: var(--border);
  background: var(--bg-3);
}

.philosophy__item-icon {
  font-size: 1.75rem;
  line-height: 1;
}

.philosophy__item-title {
  font-size: 0.875rem;
  font-weight: 700;
  color: var(--text-1);
}

.philosophy__item-desc {
  font-size: 0.8rem;
  color: var(--text-3);
  line-height: 1.6;
}

/* ── CTA banner ─────────────────────────────────────────────────────────────── */
.cta-banner {
  position: relative;
  padding: 6rem 1.5rem;
  text-align: center;
  overflow: hidden;
  background: var(--bg-0);
  border-top: 1px solid var(--border-subtle);
}

.cta-banner__glow {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 700px;
  height: 400px;
  background: radial-gradient(
    ellipse,
    rgba(63, 185, 80, 0.1) 0%,
    transparent 65%
  );
  pointer-events: none;
}

.cta-banner__inner {
  position: relative;
  z-index: 1;
}

.cta-banner__title {
  font-size: clamp(1.75rem, 4vw, 2.5rem);
  margin-bottom: 1rem;
}

.cta-banner__subtitle {
  font-size: 1.05rem;
  color: var(--text-3);
  margin-bottom: 2.5rem;
}

.cta-banner__actions {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 1rem;
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 1100px) {
  .philosophy__grid {
    grid-template-columns: repeat(3, 1fr);
  }
}

@media (max-width: 960px) {
  .hero__inner {
    grid-template-columns: 1fr;
    gap: 3rem;
    padding-top: 3rem;
    text-align: center;
  }

  .hero__badge-row,
  .hero__ctas {
    justify-content: center;
  }

  .hero__description {
    max-width: 100%;
    margin-left: auto;
    margin-right: auto;
  }

  .hero__install {
    justify-content: center;
  }

  .hero__terminal {
    max-width: 580px;
    margin: 0 auto;
  }

  .features__grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .matchers-section__layout,
  .failure-section__layout {
    grid-template-columns: 1fr;
    gap: 2.5rem;
  }

  .failure-section__layout {
    direction: rtl;
  }

  .failure-section__layout > * {
    direction: ltr;
  }
}

@media (max-width: 768px) {
  .stats__inner {
    grid-template-columns: repeat(2, 1fr);
  }

  .stats__item:nth-child(2) {
    border-right: none;
  }

  .stats__item:nth-child(3) {
    border-bottom: none;
  }

  .quickstart__steps {
    grid-template-columns: 1fr;
  }

  .philosophy__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 560px) {
  .features__grid {
    grid-template-columns: 1fr;
  }

  .hero__install-cmd {
    font-size: 0.78rem;
  }

  .philosophy__grid {
    grid-template-columns: 1fr;
  }
}
</style>

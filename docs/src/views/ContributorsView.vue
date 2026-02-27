<template>
  <div class="contributors-page">
    <!-- ── Hero ──────────────────────────────────────────────────────────── -->
    <section class="contrib-hero">
      <div class="contrib-hero__glow" aria-hidden="true"></div>
      <div class="contrib-hero__inner container">
        <span class="badge badge--green badge--dot">Open Source</span>
        <h1 class="contrib-hero__title">
          The people behind<br />
          <span class="gradient-text">gest</span>
        </h1>
        <p class="contrib-hero__desc">
          gest is built and maintained by the community. Every contribution —
          from a typo fix to a new matcher — makes it better for everyone.
        </p>
        <div class="contrib-hero__stats">
          <div
            class="contrib-hero__stat"
            v-for="stat in heroStats"
            :key="stat.label"
          >
            <span class="contrib-hero__stat-value">{{ stat.value }}</span>
            <span class="contrib-hero__stat-label">{{ stat.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Contributors grid ─────────────────────────────────────────────── -->
    <section
      class="section contrib-section"
      aria-labelledby="contributors-heading"
    >
      <div class="container">
        <div class="contrib-section__header">
          <p class="section-label">Contributors</p>
          <h2 id="contributors-heading">Meet the team</h2>
          <p class="contrib-section__subtitle">
            Everyone who has contributed code, documentation, or ideas to gest.
          </p>
        </div>

        <!-- Loading state -->
        <div
          v-if="loading"
          class="contrib-loading"
          aria-live="polite"
          aria-label="Loading contributors"
        >
          <div class="contrib-skeleton" v-for="n in 4" :key="n"></div>
        </div>

        <!-- Error state -->
        <div v-else-if="error" class="contrib-error">
          <span class="contrib-error__icon">⚠️</span>
          <p>{{ error }}</p>
        </div>

        <!-- Contributors grid -->
        <div v-else class="contrib-grid">
          <a
            v-for="(contributor, i) in contributors"
            :key="contributor.login"
            :href="contributor.html_url"
            target="_blank"
            rel="noopener noreferrer"
            class="contrib-card"
            :style="{ animationDelay: `${i * 60}ms` }"
          >
            <div class="contrib-card__avatar-wrapper">
              <img
                :src="contributor.avatar_url"
                :alt="`${contributor.login}'s avatar`"
                class="contrib-card__avatar"
                loading="lazy"
                width="96"
                height="96"
              />
            </div>

            <div class="contrib-card__info">
              <h3 class="contrib-card__name">
                {{ contributor.name || contributor.login }}
              </h3>
              <span class="contrib-card__login">@{{ contributor.login }}</span>

              <div v-if="contributor.bio" class="contrib-card__bio">
                {{ truncate(contributor.bio, 80) }}
              </div>

              <div class="contrib-card__meta">
                <span class="contrib-card__commits">
                  <svg
                    width="13"
                    height="13"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    aria-hidden="true"
                  >
                    <circle cx="12" cy="12" r="3" />
                    <line x1="3" y1="12" x2="9" y2="12" />
                    <line x1="15" y1="12" x2="21" y2="12" />
                  </svg>
                  {{ contributor.contributions }}
                  {{ contributor.contributions === 1 ? 'commit' : 'commits' }}
                </span>

                <span
                  v-if="contributor.location"
                  class="contrib-card__location"
                >
                  <svg
                    width="11"
                    height="11"
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2"
                    aria-hidden="true"
                  >
                    <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
                    <circle cx="12" cy="10" r="3" />
                  </svg>
                  {{ contributor.location }}
                </span>
              </div>
            </div>

            <div class="contrib-card__links">
              <span
                v-if="contributor.blog"
                class="contrib-card__link-icon"
                title="Website"
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
                  <circle cx="12" cy="12" r="10" />
                  <line x1="2" y1="12" x2="22" y2="12" />
                  <path
                    d="M12 2a15.3 15.3 0 0 1 4 10 15.3 15.3 0 0 1-4 10 15.3 15.3 0 0 1-4-10 15.3 15.3 0 0 1 4-10z"
                  />
                </svg>
              </span>
              <span class="contrib-card__link-icon contrib-card__link-icon--gh">
                <svg
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                  aria-hidden="true"
                >
                  <path
                    d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"
                  />
                </svg>
              </span>
            </div>
          </a>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Contribute CTA ─────────────────────────────────────────────────── -->
    <section class="section contribute" aria-labelledby="contribute-heading">
      <div class="container">
        <div class="contribute__layout">
          <div class="contribute__text">
            <p class="section-label">Get involved</p>
            <h2 id="contribute-heading">Want to see your name here?</h2>
            <p>
              gest is open source and contributions of all kinds are warmly
              welcomed. Whether you want to fix a bug, add a matcher, improve
              the docs, or just share feedback — every bit helps.
            </p>
            <div class="contribute__ctas">
              <a
                href="https://github.com/caiolandgraf/gest/fork"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn--primary"
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
                  <circle cx="18" cy="18" r="3" />
                  <circle cx="6" cy="6" r="3" />
                  <path d="M6 21V9a9 9 0 0 0 9 9" />
                </svg>
                Fork on GitHub
              </a>
              <a
                href="https://github.com/caiolandgraf/gest/issues/new"
                target="_blank"
                rel="noopener noreferrer"
                class="btn btn--ghost"
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
                  <circle cx="12" cy="12" r="10" />
                  <line x1="12" y1="8" x2="12" y2="12" />
                  <line x1="12" y1="16" x2="12.01" y2="16" />
                </svg>
                Open an issue
              </a>
            </div>
          </div>

          <div class="contribute__steps">
            <div
              v-for="(step, i) in howToContribute"
              :key="step.title"
              class="contribute__step"
            >
              <div class="contribute__step-num">{{ i + 1 }}</div>
              <div>
                <h4 class="contribute__step-title">{{ step.title }}</h4>
                <p class="contribute__step-desc">{{ step.desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </section>

    <hr class="divider" />

    <!-- ── Good first issues ──────────────────────────────────────────────── -->
    <section class="section good-first" aria-labelledby="first-issues-heading">
      <div class="container">
        <div class="good-first__header">
          <p class="section-label">Good first issues</p>
          <h2 id="first-issues-heading">Not sure where to start?</h2>
          <p class="good-first__subtitle">
            Here are some ideas for newcomers — no contribution is too small.
          </p>
        </div>

        <div class="good-first__grid">
          <div
            v-for="idea in firstIssues"
            :key="idea.title"
            class="first-issue-card"
          >
            <div class="first-issue-card__icon">{{ idea.icon }}</div>
            <div class="first-issue-card__content">
              <div class="first-issue-card__header">
                <h4 class="first-issue-card__title">{{ idea.title }}</h4>
                <span
                  class="first-issue-card__difficulty"
                  :class="`first-issue-card__difficulty--${idea.level}`"
                >
                  {{ idea.levelLabel }}
                </span>
              </div>
              <p class="first-issue-card__desc">{{ idea.desc }}</p>
            </div>
          </div>
        </div>

        <div class="good-first__cta">
          <a
            href="https://github.com/caiolandgraf/gest/issues?q=is%3Aopen+is%3Aissue+label%3A%22good+first+issue%22"
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
            Browse good first issues on GitHub
          </a>
        </div>
      </div>
    </section>

    <!-- ── Recognition banner ─────────────────────────────────────────────── -->
    <section class="recognition-banner" aria-labelledby="recognition-heading">
      <div class="recognition-banner__glow" aria-hidden="true"></div>
      <div class="recognition-banner__inner container">
        <span class="recognition-banner__emoji" aria-hidden="true">🧪</span>
        <h2 id="recognition-heading" class="recognition-banner__title">
          Every commit tells a story.
        </h2>
        <p class="recognition-banner__subtitle">
          Your name, your code, your impact — permanently in the history of
          gest.
        </p>
        <a
          href="https://github.com/caiolandgraf/gest"
          target="_blank"
          rel="noopener noreferrer"
          class="btn btn--primary"
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
          Star the repo &amp; get started
        </a>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

/* ── State ─────────────────────────────────────────────────────────────────── */
const contributors = ref([])
const loading = ref(true)
const error = ref(null)

/* ── Fetch contributors from GitHub API ────────────────────────────────────── */
onMounted(async () => {
  try {
    // Fetch contributor list (sorted by contributions)
    const listRes = await fetch(
      'https://api.github.com/repos/caiolandgraf/gest/contributors?per_page=30',
      { headers: { Accept: 'application/vnd.github+json' } }
    )

    if (!listRes.ok) throw new Error(`GitHub API returned ${listRes.status}`)

    const list = await listRes.json()

    // Enrich each contributor with full profile data
    const enriched = await Promise.allSettled(
      list.map(c =>
        fetch(`https://api.github.com/users/${c.login}`, {
          headers: { Accept: 'application/vnd.github+json' }
        })
          .then(r => r.json())
          .then(profile => ({
            login: c.login,
            avatar_url: c.avatar_url,
            html_url: c.html_url,
            contributions: c.contributions,
            name: profile.name || null,
            bio: profile.bio || null,
            location: profile.location || null,
            blog: profile.blog || null
          }))
      )
    )

    contributors.value = enriched
      .filter(r => r.status === 'fulfilled')
      .map(r => r.value)
  } catch (e) {
    // Fallback to known static data so the page is never empty
    contributors.value = [
      {
        login: 'caiolandgraf',
        avatar_url: 'https://github.com/caiolandgraf.png',
        html_url: 'https://github.com/caiolandgraf',
        contributions: 1,
        name: 'Caio Landgraf',
        bio: 'Author of gest — bringing Jest DX to Go.',
        location: null,
        blog: null
      }
    ]

    if (e.message.includes('429')) {
      error.value = 'GitHub API rate limit reached. Showing cached data below.'
    }
  } finally {
    loading.value = false
  }
})

/* ── Helpers ───────────────────────────────────────────────────────────────── */
function truncate(str, max) {
  if (!str || str.length <= max) return str
  return str.slice(0, max).trimEnd() + '…'
}

/* ── Hero stats ────────────────────────────────────────────────────────────── */
const heroStats = [
  { value: '1', label: 'Packages' },
  { value: '10', label: 'Matchers' },
  { value: '0', label: 'Dependencies' },
  { value: 'MIT', label: 'License' }
]

/* ── How to contribute ─────────────────────────────────────────────────────── */
const howToContribute = [
  {
    title: 'Fork the repository',
    desc: 'Click "Fork" on GitHub to create your own copy of the gest repo.'
  },
  {
    title: 'Make your changes',
    desc: 'Fix a bug, add a matcher, improve documentation — all contributions are welcome.'
  },
  {
    title: 'Run the example tests',
    desc: 'Verify your changes work by running `go run .` inside the example/ directory.'
  },
  {
    title: 'Open a Pull Request',
    desc: 'Push your branch and open a PR describing what you changed and why.'
  }
]

/* ── Good first issues ─────────────────────────────────────────────────────── */
const firstIssues = [
  {
    icon: '🔧',
    title: 'Add a new matcher',
    desc: "Think of a useful assertion that's missing? Add a ToMatchRegex, ToBeEmpty, or ToBePositive matcher — the pattern is well-established in the codebase.",
    level: 'easy',
    levelLabel: 'Easy'
  },
  {
    icon: '📖',
    title: 'Improve the docs',
    desc: 'Found something confusing or missing in the README or this documentation site? Fix it! Docs PRs are some of the most impactful contributions.',
    level: 'easy',
    levelLabel: 'Easy'
  },
  {
    icon: '🧪',
    title: 'Add example spec files',
    desc: 'Extend the example/ directory with more real-world spec file patterns to help new users learn how to structure their tests.',
    level: 'easy',
    levelLabel: 'Easy'
  },
  {
    icon: '🌍',
    title: 'Improve error messages',
    desc: 'Make failure output even more helpful. Add context, improve formatting for specific types, or handle edge cases in existing matchers.',
    level: 'medium',
    levelLabel: 'Medium'
  },
  {
    icon: '⚡',
    title: 'Performance benchmarks',
    desc: 'Add benchmark tests to measure and track the performance of the test runner itself under large numbers of suites and test cases.',
    level: 'medium',
    levelLabel: 'Medium'
  },
  {
    icon: '🔌',
    title: 'go test compatibility',
    desc: 'Explore ways to make gest output compatible with go test -json so it can integrate with existing CI tooling and test result parsers.',
    level: 'hard',
    levelLabel: 'Advanced'
  }
]
</script>

<style scoped>
/* ── Hero ───────────────────────────────────────────────────────────────────── */
.contrib-hero {
  position: relative;
  padding: 6rem 1.5rem 4rem;
  text-align: center;
  overflow: hidden;
}

.contrib-hero__glow {
  position: absolute;
  top: -20%;
  left: 50%;
  transform: translateX(-50%);
  width: 800px;
  height: 500px;
  background: radial-gradient(
    ellipse at top,
    rgba(63, 185, 80, 0.09) 0%,
    transparent 60%
  );
  pointer-events: none;
  user-select: none;
}

.contrib-hero__inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.contrib-hero__title {
  font-size: clamp(2.25rem, 5vw, 3.25rem);
  margin-top: 1rem;
  line-height: 1.1;
}

.contrib-hero__desc {
  font-size: 1.05rem;
  color: var(--text-2);
  line-height: 1.75;
  max-width: 560px;
  margin-top: 0.25rem;
}

.contrib-hero__stats {
  display: flex;
  align-items: center;
  gap: 0;
  margin-top: 2rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.contrib-hero__stat {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 1.25rem 2.5rem;
  border-right: 1px solid var(--border-subtle);
  transition: background var(--duration-fast) var(--ease);
}

.contrib-hero__stat:last-child {
  border-right: none;
}

.contrib-hero__stat:hover {
  background: var(--bg-3);
}

.contrib-hero__stat-value {
  font-family: var(--font-mono);
  font-size: 1.5rem;
  font-weight: 700;
  color: var(--green);
  line-height: 1;
}

.contrib-hero__stat-label {
  font-size: 0.72rem;
  color: var(--text-4);
  letter-spacing: 0.04em;
  text-transform: uppercase;
  font-weight: 600;
}

/* ── Contributors section ───────────────────────────────────────────────────── */
.contrib-section__header {
  text-align: center;
  margin-bottom: 3rem;
}

.contrib-section__subtitle {
  font-size: 1.05rem;
  color: var(--text-3);
  margin-top: 0.75rem;
}

/* ── Loading skeleton ───────────────────────────────────────────────────────── */
.contrib-loading {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

@keyframes skeleton-pulse {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0.4;
  }
}

.contrib-skeleton {
  height: 220px;
  background: var(--bg-3);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  animation: skeleton-pulse 1.8s ease-in-out infinite;
}

/* ── Error state ────────────────────────────────────────────────────────────── */
.contrib-error {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 1rem 1.25rem;
  background: rgba(248, 81, 73, 0.07);
  border: 1px solid rgba(248, 81, 73, 0.2);
  border-radius: var(--radius-md);
  color: var(--red);
  font-size: 0.9rem;
  margin-bottom: 2rem;
}

.contrib-error__icon {
  font-size: 1.1rem;
  flex-shrink: 0;
}

/* ── Contributors grid ──────────────────────────────────────────────────────── */
.contrib-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 1.5rem;
}

/* ── Contributor card ───────────────────────────────────────────────────────── */
.contrib-card {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
  padding: 2rem 1.5rem 1.5rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-lg);
  text-decoration: none;
  text-align: center;
  overflow: hidden;
  animation: fadeInUp 0.45s var(--ease) both;
  transition:
    border-color var(--duration-normal) var(--ease),
    box-shadow var(--duration-normal) var(--ease),
    transform var(--duration-normal) var(--ease);
}

.contrib-card::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: linear-gradient(
    90deg,
    transparent,
    var(--green-dim),
    transparent
  );
  opacity: 0;
  transition: opacity var(--duration-normal) var(--ease);
}

.contrib-card:hover {
  border-color: rgba(63, 185, 80, 0.3);
  box-shadow: var(--shadow-md), var(--shadow-glow);
  transform: translateY(-3px);
}

.contrib-card:hover::before {
  opacity: 1;
}

/* ── Avatar ─────────────────────────────────────────────────────────────────── */
.contrib-card__avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.contrib-card__avatar {
  width: 96px;
  height: 96px;
  border-radius: 50%;
  border: 3px solid var(--bg-4);
  object-fit: cover;
  display: block;
  position: relative;
  z-index: 1;
  transition:
    border-color var(--duration-normal) var(--ease),
    box-shadow var(--duration-normal) var(--ease);
}

.contrib-card:hover .contrib-card__avatar {
  border-color: var(--green-dim);
  box-shadow:
    0 0 0 4px rgba(63, 185, 80, 0.12),
    0 0 20px rgba(63, 185, 80, 0.2);
}

/* ── Card info ──────────────────────────────────────────────────────────────── */
.contrib-card__info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.3rem;
  width: 100%;
}

.contrib-card__name {
  font-size: 1.05rem;
  font-weight: 700;
  color: var(--text-1);
  line-height: 1.25;
}

.contrib-card__login {
  font-family: var(--font-mono);
  font-size: 0.8rem;
  color: var(--green);
  font-weight: 500;
}

.contrib-card__bio {
  font-size: 0.825rem;
  color: var(--text-3);
  line-height: 1.55;
  margin-top: 0.25rem;
  max-width: 220px;
}

.contrib-card__meta {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 0.5rem;
}

.contrib-card__commits,
.contrib-card__location {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 0.78rem;
  color: var(--text-4);
  font-family: var(--font-mono);
}

.contrib-card__commits {
  color: var(--green-dim);
}

/* ── Card links ─────────────────────────────────────────────────────────────── */
.contrib-card__links {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: auto;
  padding-top: 0.75rem;
  border-top: 1px solid var(--border-subtle);
  width: 100%;
}

.contrib-card__link-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: var(--radius-sm);
  color: var(--text-4);
  background: var(--bg-4);
  border: 1px solid var(--border-subtle);
  transition: all var(--duration-fast) var(--ease);
}

.contrib-card:hover .contrib-card__link-icon {
  color: var(--text-2);
  border-color: var(--border);
}

.contrib-card__link-icon--gh:hover {
  color: var(--text-1) !important;
  background: var(--bg-5) !important;
}

/* ── Contribute section ─────────────────────────────────────────────────────── */
.contribute__layout {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 4rem;
  align-items: start;
}

.contribute__text h2 {
  margin-bottom: 1rem;
}

.contribute__text p {
  margin-bottom: 1.5rem;
  font-size: 1rem;
  line-height: 1.75;
}

.contribute__ctas {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.contribute__steps {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.contribute__step {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.25rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition:
    border-color var(--duration-normal) var(--ease),
    background var(--duration-normal) var(--ease);
}

.contribute__step:hover {
  border-color: var(--border);
  background: var(--bg-3);
}

.contribute__step-num {
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
  margin-top: 2px;
}

.contribute__step-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-1);
  margin-bottom: 0.3rem;
}

.contribute__step-desc {
  font-size: 0.875rem;
  color: var(--text-3);
  line-height: 1.6;
}

/* ── Good first issues ──────────────────────────────────────────────────────── */
.good-first__header {
  text-align: center;
  max-width: 560px;
  margin: 0 auto 2.5rem;
}

.good-first__subtitle {
  font-size: 1.05rem;
  color: var(--text-3);
  margin-top: 0.75rem;
}

.good-first__grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1.25rem;
  margin-bottom: 2.5rem;
}

.first-issue-card {
  display: flex;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-2);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-md);
  transition:
    border-color var(--duration-normal) var(--ease),
    box-shadow var(--duration-normal) var(--ease),
    transform var(--duration-normal) var(--ease);
}

.first-issue-card:hover {
  border-color: var(--border);
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

.first-issue-card__icon {
  font-size: 1.5rem;
  line-height: 1;
  flex-shrink: 0;
  margin-top: 2px;
}

.first-issue-card__content {
  min-width: 0;
}

.first-issue-card__header {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 0.5rem;
}

.first-issue-card__title {
  font-size: 0.95rem;
  font-weight: 700;
  color: var(--text-1);
}

.first-issue-card__difficulty {
  font-size: 0.68rem;
  font-weight: 700;
  font-family: var(--font-mono);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  padding: 0.15rem 0.5rem;
  border-radius: var(--radius-full);
}

.first-issue-card__difficulty--easy {
  background: rgba(63, 185, 80, 0.12);
  color: var(--green-bright);
  border: 1px solid rgba(63, 185, 80, 0.2);
}

.first-issue-card__difficulty--medium {
  background: rgba(227, 179, 65, 0.12);
  color: var(--yellow);
  border: 1px solid rgba(227, 179, 65, 0.2);
}

.first-issue-card__difficulty--hard {
  background: rgba(248, 81, 73, 0.1);
  color: var(--red);
  border: 1px solid rgba(248, 81, 73, 0.2);
}

.first-issue-card__desc {
  font-size: 0.85rem;
  color: var(--text-3);
  line-height: 1.65;
}

.good-first__cta {
  display: flex;
  justify-content: center;
}

/* ── Recognition banner ─────────────────────────────────────────────────────── */
.recognition-banner {
  position: relative;
  padding: 6rem 1.5rem;
  text-align: center;
  overflow: hidden;
  background: var(--bg-0);
  border-top: 1px solid var(--border-subtle);
}

.recognition-banner__glow {
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

.recognition-banner__inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.recognition-banner__emoji {
  font-size: 3rem;
  line-height: 1;
  filter: drop-shadow(0 0 16px rgba(63, 185, 80, 0.4));
  display: block;
  animation: pulse 3s ease-in-out infinite;
}

.recognition-banner__title {
  font-size: clamp(1.75rem, 4vw, 2.5rem);
  margin: 0.5rem 0 0;
}

.recognition-banner__subtitle {
  font-size: 1.05rem;
  color: var(--text-3);
  max-width: 480px;
  margin-bottom: 1rem;
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 1024px) {
  .good-first__grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 860px) {
  .contribute__layout {
    grid-template-columns: 1fr;
    gap: 2.5rem;
  }
}

@media (max-width: 768px) {
  .contrib-hero__stats {
    flex-wrap: wrap;
    border-radius: var(--radius-md);
  }

  .contrib-hero__stat {
    padding: 1rem 1.75rem;
    flex: 1 0 calc(50% - 1px);
  }

  .contrib-hero__stat:nth-child(2) {
    border-right: none;
  }

  .contrib-hero__stat:nth-child(3),
  .contrib-hero__stat:nth-child(4) {
    border-top: 1px solid var(--border-subtle);
  }

  .contrib-hero__stat:nth-child(4) {
    border-right: none;
  }

  .good-first__grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 480px) {
  .contrib-grid {
    grid-template-columns: 1fr;
  }

  .contrib-hero__stat {
    padding: 0.85rem 1.25rem;
  }
}
</style>

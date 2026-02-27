<template>
  <footer class="footer">
    <div class="footer__glow" aria-hidden="true"></div>

    <div class="footer__inner container">

      <!-- Top row -->
      <div class="footer__top">

        <!-- Brand -->
        <div class="footer__brand">
          <RouterLink to="/" class="footer__logo">
            <span class="footer__logo-icon">🧪</span>
            <span class="footer__logo-text">gest</span>
          </RouterLink>
          <p class="footer__tagline">
            A Jest-inspired testing framework for Go.<br />
            Beautiful output. Zero dependencies.
          </p>
          <div class="footer__badges">
            <a
              href="https://pkg.go.dev/github.com/caiolandgraf/gest"
              target="_blank"
              rel="noopener noreferrer"
              class="footer__badge"
              aria-label="Go reference"
            >
              <img
                src="https://pkg.go.dev/badge/github.com/caiolandgraf/gest.svg"
                alt="Go Reference"
                width="100"
                height="20"
              />
            </a>
            <a
              href="https://github.com/caiolandgraf/gest/blob/main/LICENSE"
              target="_blank"
              rel="noopener noreferrer"
              class="footer__tag"
            >
              MIT License
            </a>
          </div>
        </div>

        <!-- Links -->
        <nav class="footer__links" aria-label="Footer navigation">
          <div class="footer__col">
            <h3 class="footer__col-title">Navigation</h3>
            <ul class="footer__col-list">
              <li><RouterLink to="/" class="footer__link">Home</RouterLink></li>
              <li><RouterLink to="/docs" class="footer__link">Documentation</RouterLink></li>
              <li><RouterLink to="/contributors" class="footer__link">Contributors</RouterLink></li>
            </ul>
          </div>

          <div class="footer__col">
            <h3 class="footer__col-title">Docs</h3>
            <ul class="footer__col-list">
              <li><RouterLink to="/docs#installation" class="footer__link">Installation</RouterLink></li>
              <li><RouterLink to="/docs#usage" class="footer__link">Basic Usage</RouterLink></li>
              <li><RouterLink to="/docs#matchers" class="footer__link">Matchers</RouterLink></li>
              <li><RouterLink to="/docs#coverage" class="footer__link">Coverage</RouterLink></li>
              <li><RouterLink to="/docs#api" class="footer__link">Full API</RouterLink></li>
            </ul>
          </div>

          <div class="footer__col">
            <h3 class="footer__col-title">Community</h3>
            <ul class="footer__col-list">
              <li>
                <a
                  href="https://github.com/caiolandgraf/gest"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="footer__link footer__link--external"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                    <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z" />
                  </svg>
                  GitHub
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/caiolandgraf/gest/issues"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="footer__link footer__link--external"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <circle cx="12" cy="12" r="10" />
                    <line x1="12" y1="8" x2="12" y2="12" />
                    <line x1="12" y1="16" x2="12.01" y2="16" />
                  </svg>
                  Issues
                </a>
              </li>
              <li>
                <a
                  href="https://github.com/caiolandgraf/gest/pulls"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="footer__link footer__link--external"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <circle cx="18" cy="18" r="3" />
                    <circle cx="6" cy="6" r="3" />
                    <path d="M6 21V9a9 9 0 0 0 9 9" />
                  </svg>
                  Pull Requests
                </a>
              </li>
              <li>
                <a
                  href="https://pkg.go.dev/github.com/caiolandgraf/gest"
                  target="_blank"
                  rel="noopener noreferrer"
                  class="footer__link footer__link--external"
                >
                  <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true">
                    <path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20" />
                    <path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z" />
                  </svg>
                  pkg.go.dev
                </a>
              </li>
            </ul>
          </div>
        </nav>
      </div>

      <hr class="footer__divider" />

      <!-- Bottom row -->
      <div class="footer__bottom">
        <p class="footer__copy">
          &copy; {{ year }} gest — MIT License
        </p>
        <p class="footer__made">
          Made with 🧪 by
          <a
            href="https://github.com/caiolandgraf"
            target="_blank"
            rel="noopener noreferrer"
            class="footer__author"
          >@caiolandgraf</a>
        </p>
      </div>

    </div>
  </footer>
</template>

<script setup>
const year = new Date().getFullYear()
</script>

<style scoped>
/* ── Base ───────────────────────────────────────────────────────────────────── */
.footer {
  position: relative;
  background: var(--bg-0);
  border-top: 1px solid var(--border-subtle);
  overflow: hidden;
}

/* ── Glow ───────────────────────────────────────────────────────────────────── */
.footer__glow {
  position: absolute;
  bottom: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 600px;
  height: 200px;
  background: radial-gradient(ellipse at bottom, rgba(63, 185, 80, 0.06) 0%, transparent 70%);
  pointer-events: none;
  user-select: none;
}

/* ── Inner ──────────────────────────────────────────────────────────────────── */
.footer__inner {
  padding: 4rem 1.5rem 2.5rem;
}

/* ── Top ────────────────────────────────────────────────────────────────────── */
.footer__top {
  display: grid;
  grid-template-columns: 1fr 2fr;
  gap: 4rem;
  align-items: flex-start;
}

/* ── Brand ──────────────────────────────────────────────────────────────────── */
.footer__logo {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  text-decoration: none;
  margin-bottom: 1rem;
}

.footer__logo-icon {
  font-size: 1.3rem;
  filter: drop-shadow(0 0 6px rgba(63, 185, 80, 0.35));
}

.footer__logo-text {
  font-family: var(--font-mono);
  font-size: 1.15rem;
  font-weight: 800;
  color: var(--text-1);
  letter-spacing: -0.03em;
  transition: color var(--duration-fast) var(--ease);
}

.footer__logo:hover .footer__logo-text {
  color: var(--green);
}

.footer__tagline {
  font-size: 0.875rem;
  color: var(--text-3);
  line-height: 1.65;
  margin-bottom: 1.25rem;
  max-width: 260px;
}

.footer__badges {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
}

.footer__badge {
  display: inline-flex;
  opacity: 0.85;
  transition: opacity var(--duration-fast) var(--ease);
}

.footer__badge:hover {
  opacity: 1;
}

.footer__tag {
  display: inline-flex;
  align-items: center;
  font-size: 0.72rem;
  font-weight: 600;
  font-family: var(--font-mono);
  color: var(--text-3);
  background: var(--bg-4);
  border: 1px solid var(--border-subtle);
  border-radius: var(--radius-full);
  padding: 0.18rem 0.65rem;
  text-decoration: none;
  transition: all var(--duration-fast) var(--ease);
  white-space: nowrap;
}

.footer__tag:hover {
  color: var(--text-1);
  border-color: var(--border);
}

/* ── Links ──────────────────────────────────────────────────────────────────── */
.footer__links {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 2rem;
}

.footer__col-title {
  font-size: 0.72rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.1em;
  color: var(--text-3);
  margin-bottom: 1rem;
}

.footer__col-list {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.footer__link {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  font-size: 0.875rem;
  color: var(--text-3);
  text-decoration: none;
  transition: color var(--duration-fast) var(--ease);
  width: fit-content;
}

.footer__link:hover {
  color: var(--text-1);
}

.footer__link--external svg {
  opacity: 0.6;
  flex-shrink: 0;
}

/* ── Divider ────────────────────────────────────────────────────────────────── */
.footer__divider {
  border: none;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border-subtle), transparent);
  margin: 2.5rem 0 1.75rem;
}

/* ── Bottom ─────────────────────────────────────────────────────────────────── */
.footer__bottom {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.footer__copy,
.footer__made {
  font-size: 0.8rem;
  color: var(--text-4);
  line-height: 1;
}

.footer__author {
  color: var(--text-3);
  text-decoration: none;
  font-weight: 500;
  transition: color var(--duration-fast) var(--ease);
}

.footer__author:hover {
  color: var(--green);
}

/* ── Responsive ─────────────────────────────────────────────────────────────── */
@media (max-width: 900px) {
  .footer__top {
    grid-template-columns: 1fr;
    gap: 2.5rem;
  }

  .footer__tagline {
    max-width: 100%;
  }
}

@media (max-width: 600px) {
  .footer__links {
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }

  .footer__bottom {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.4rem;
  }
}

@media (max-width: 400px) {
  .footer__links {
    grid-template-columns: 1fr;
  }
}
</style>

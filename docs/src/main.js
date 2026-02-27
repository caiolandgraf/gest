import { createApp } from 'vue'
import { createRouter, createWebHashHistory } from 'vue-router'
import App from './App.vue'
import './assets/main.css'

import ContributorsView from './views/ContributorsView.vue'
import DocsView from './views/DocsView.vue'
import HomeView from './views/HomeView.vue'

const routes = [
  { path: '/',             component: HomeView,         meta: { title: 'gest — Jest-inspired testing for Go' } },
  { path: '/docs',         component: DocsView,         meta: { title: 'Docs — gest' } },
  { path: '/contributors', component: ContributorsView, meta: { title: 'Contributors — gest' } },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
  scrollBehavior(to, _from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) {
      return {
        el: to.hash,
        behavior: 'smooth',
        top: 80,
      }
    }
    return { top: 0, behavior: 'smooth' }
  },
})

router.afterEach((to) => {
  document.title = to.meta?.title ?? 'gest'
})

createApp(App).use(router).mount('#app')

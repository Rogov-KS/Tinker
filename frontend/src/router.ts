import { getState, subscribe, clearAuth } from './state'
import { renderLoginPage } from './pages/login'
import { renderSignupPage } from './pages/signup'
import { renderFeedPage } from './pages/feed'
import { renderProfilePage } from './pages/profile'
import { renderSearchPage } from './pages/search'
import { renderSettingsPage } from './pages/settings'
import { renderLandingPage } from './pages/landing'

type RouteHandler = (params: Record<string, string>) => Promise<void> | void

interface RouteConfig {
  pattern: RegExp
  auth?: boolean
  handler: RouteHandler
}

const routes: RouteConfig[] = [
  { pattern: /^#\/login$/, handler: () => renderLoginPage() },
  { pattern: /^#\/signup$/, handler: () => renderSignupPage() },
  { pattern: /^#\/feed$/, auth: true, handler: () => renderFeedPage() },
  {
    pattern: /^#\/profile\/(?<id>[\w-]+)$/,
    auth: false,
    handler: (params) => renderProfilePage(params.id),
  },
  { pattern: /^#\/profile$/, auth: true, handler: () => {
    const state = getState()
    if (state.user) {
      return renderProfilePage(state.user.id)
    }
  } },
  { pattern: /^#\/search$/, auth: true, handler: () => renderSearchPage() },
  { pattern: /^#\/settings$/, auth: true, handler: () => renderSettingsPage() },
]

export function initRouter() {
  window.addEventListener('hashchange', handleRoute)
  subscribe(handleRoute)
  if (!window.location.hash) {
    window.location.hash = getState().token ? '#/feed' : '#/'
  } else {
    handleRoute()
  }
}

export function navigateTo(hash: string) {
  if (window.location.hash === hash) {
    handleRoute()
  } else {
    window.location.hash = hash
  }
}

async function handleRoute() {
  const hash = window.location.hash || '#/login'
  const app = document.getElementById('app')

  if (!app) return

  const state = getState()

  if (!state.token && (hash === '#/' || hash === '#/landing')) {
    await renderLandingPage()
    return
  }

  const match = routes.find((route) => route.pattern.test(hash))

  if (!match) {
    if (state.token) {
      navigateTo('#/feed')
    } else {
      navigateTo('#/login')
    }
    return
  }

  if (match.auth && !state.token) {
    clearAuth()
    navigateTo('#/login')
    return
  }

  const result = match.pattern.exec(hash)
  const params = result?.groups ?? {}

  await match.handler(params)
}



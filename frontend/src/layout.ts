import { getState, clearAuth } from './state'
import { navigateTo } from './router'
import { createElement, clearChildren, createAvatar } from './utils/dom'

export function renderWithShell(
  content: HTMLElement,
  options: { showHeader?: boolean } = { showHeader: true }
) {
  const root = document.getElementById('app')
  if (!root) return

  clearChildren(root)

  const wrapper = createElement('div', { className: 'app-shell' })

  if (options.showHeader !== false) {
    wrapper.appendChild(createHeader())
  }

  const main = createElement('main', { className: 'app-main' })
  main.appendChild(content)
  wrapper.appendChild(main)

  root.appendChild(wrapper)
}

function createHeader() {
  const state = getState()
  const header = createElement('header', { className: 'app-header' })

  const left = createElement('div', { className: 'header-left' })
  const logo = createElement('button', {
    className: 'logo-button',
    text: 'T-News',
  })
  logo.addEventListener('click', () => navigateTo('#/feed'))
  left.appendChild(logo)

  const nav = createElement('nav', { className: 'header-nav' })
  const links = [
    { label: 'Лента', hash: '#/feed' },
    { label: 'Поиск', hash: '#/search' },
    { label: 'Профиль', hash: '#/profile' },
    { label: 'Настройки', hash: '#/settings' },
  ]

  links.forEach((link) => {
    const btn = createElement('button', {
      className: 'nav-link',
      text: link.label,
    })
    btn.addEventListener('click', () => navigateTo(link.hash))
    nav.appendChild(btn)
  })

  left.appendChild(nav)

  const right = createElement('div', { className: 'header-right' })
  if (state.user) {
    const avatar = createAvatar(state.user.username, state.user.avatar)
    avatar.addEventListener('click', () =>
      navigateTo(`#/profile/${state.user?.id}`)
    )
    right.appendChild(avatar)
  }

  const logoutBtn = createElement('button', {
    className: 'ghost-button',
    text: 'Выйти',
  })
  logoutBtn.addEventListener('click', () => {
    clearAuth()
    navigateTo('#/login')
  })
  right.appendChild(logoutBtn)

  header.appendChild(left)
  header.appendChild(right)

  return header
}



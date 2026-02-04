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

  // Форма поиска в хедере
  const searchForm = createElement('form', { className: 'header-search-form' })
  const searchInput = createElement('input', {
    attrs: {
      type: 'text',
      name: 'query',
      placeholder: 'Поиск...',
    },
  }) as HTMLInputElement
  
  const searchButton = createElement('button', {
    className: 'header-search-button',
    attrs: {
      type: 'submit',
    },
  })
  searchButton.innerHTML = '🔍'
  searchButton.setAttribute('aria-label', 'Поиск')
  
  const performSearch = () => {
    const query = searchInput.value.trim()
    if (!query) return
    
    // Переходим на страницу поиска
    navigateTo(`#/search?query=${encodeURIComponent(query)}&type=posts`)
    
    // После перехода находим форму поиска на странице и программно устанавливаем значения и вызываем submit
    setTimeout(() => {
      const searchPageForm = document.querySelector('.search-form') as HTMLFormElement
      if (searchPageForm) {
        // Устанавливаем значение поля query
        const queryInput = searchPageForm.querySelector('input[name="query"]') as HTMLInputElement
        if (queryInput) {
          queryInput.value = query
        }
        
        // Устанавливаем тип поиска на "posts"
        const postsRadio = searchPageForm.querySelector('input[type="radio"][name="type"][value="posts"]') as HTMLInputElement
        if (postsRadio) {
          postsRadio.checked = true
        }
        
        // Вызываем submit формы
        if (searchPageForm.requestSubmit) {
          searchPageForm.requestSubmit()
        } else {
          const submitEvent = new Event('submit', { bubbles: true, cancelable: true })
          searchPageForm.dispatchEvent(submitEvent)
        }
      }
    }, 100) // Небольшая задержка чтобы страница успела отрендериться
  }
  
  searchForm.addEventListener('submit', (event) => {
    event.preventDefault()
    performSearch()
  })
  
  searchButton.addEventListener('click', (event) => {
    event.preventDefault()
    performSearch()
  })
  
  searchForm.appendChild(searchInput)
  searchForm.appendChild(searchButton)

  const center = createElement('div', { className: 'header-center' })
  center.appendChild(searchForm)

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
  header.appendChild(center)
  header.appendChild(right)

  return header
}



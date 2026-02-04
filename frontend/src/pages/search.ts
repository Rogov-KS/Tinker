import { renderWithShell } from '../layout'
import { createElement } from '../utils/dom'
import { searchEntities, getAllUsers, getAllPosts } from '../api'
import { createPostCard } from '../components/post-card'
import { navigateTo } from '../router'
import type { Post, User } from '../types'

export function renderSearchPage(initialQuery: string = '', initialType: 'users' | 'posts' = 'users') {
  const container = createElement('div', { className: 'page-section' })
  const title = createElement('h1', { text: 'Поиск' })
  const form = document.createElement('form')
  form.className = 'search-form'
  form.innerHTML = `
    <input type="text" name="query" placeholder="Кого или что ищем?" value="${initialQuery}" />
    <div class="toggle-group">
      <label>
        <input type="radio" name="type" value="users" ${initialType === 'users' ? 'checked' : ''} />
        <span>Пользователи</span>
      </label>
      <label>
        <input type="radio" name="type" value="posts" ${initialType === 'posts' ? 'checked' : ''} />
        <span>Посты</span>
      </label>
    </div>
    <div class="search-buttons">
      <button type="submit" class="primary-button">Найти</button>
      <button type="button" class="primary-button" id="show-all-btn">Показать всех</button>
    </div>
  `

  const results = createElement('div', { className: 'search-results' })
  
  // Автоматический поиск при переключении радиобаттона
  const radioButtons = form.querySelectorAll('input[type="radio"][name="type"]')
  radioButtons.forEach((radio) => {
    radio.addEventListener('change', async () => {
      const formData = new FormData(form)
      const query = String(formData.get('query') || '').trim()
      const type = formData.get('type') === 'posts' ? 'posts' : 'users'
      if (query) {
        results.innerHTML = '<p class="muted">Поиск...</p>'
        try {
          const data = await searchEntities(query, type)
          renderResults(results, data, type)
        } catch (error) {
          results.innerHTML = ''
          results.appendChild(
            createElement('p', {
              className: 'error-text',
              text: (error as Error).message,
            })
          )
        }
      }
    })
  })

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    const query = String(formData.get('query') || '').trim()
    const type = formData.get('type') === 'posts' ? 'posts' : 'users'
    if (!query) return
    results.innerHTML = '<p class="muted">Поиск...</p>'
    try {
      const data = await searchEntities(query, type)
      renderResults(results, data, type)
    } catch (error) {
      results.innerHTML = ''
      results.appendChild(
        createElement('p', {
          className: 'error-text',
          text: (error as Error).message,
        })
      )
    }
  })

  const showAllBtn = form.querySelector('#show-all-btn') as HTMLButtonElement
  showAllBtn.addEventListener('click', async () => {
    const formData = new FormData(form)
    const type = formData.get('type') === 'posts' ? 'posts' : 'users'
    results.innerHTML = '<p class="muted">Загрузка...</p>'
    try {
      let data: Array<User | Post>
      if (type === 'users') {
        data = await getAllUsers()
      } else {
        data = await getAllPosts()
      }
      renderResults(results, data, type)
    } catch (error) {
      results.innerHTML = ''
      results.appendChild(
        createElement('p', {
          className: 'error-text',
          text: (error as Error).message,
        })
      )
    }
  })

  container.appendChild(title)
  container.appendChild(form)
  container.appendChild(results)

  renderWithShell(container)
  
  // Автоматический поиск при наличии начального запроса
  if (initialQuery) {
    results.innerHTML = '<p class="muted">Поиск...</p>'
    searchEntities(initialQuery, initialType)
      .then((data) => {
        renderResults(results, data, initialType)
      })
      .catch((error) => {
        results.innerHTML = ''
        results.appendChild(
          createElement('p', {
            className: 'error-text',
            text: (error as Error).message,
          })
        )
      })
  }
}

function renderResults(
  container: HTMLElement,
  data: Array<User | Post>,
  type: 'users' | 'posts'
) {
  container.innerHTML = ''

  if (!data.length) {
    container.appendChild(
      createElement('p', { className: 'muted', text: 'Ничего не найдено' })
    )
    return
  }

  if (type === 'users') {
    const list = createElement('div', { className: 'user-results' })
    ;(data as User[]).forEach((user) => {
      const item = createElement('button', {
        className: 'user-result',
        text: user.username,
      })
      item.addEventListener('click', () => navigateTo(`#/profile/${user.id}`))
      list.appendChild(item)
    })
    container.appendChild(list)
  } else {
    const postsList = createElement('div', { className: 'post-list' })
    ;(data as Post[]).forEach((post) => {
      postsList.appendChild(createPostCard(post))
    })
    container.appendChild(postsList)
  }
}



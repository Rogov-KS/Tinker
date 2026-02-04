import { updateProfile, fetchFollowing, unfollowUser } from '../api'
import { renderWithShell } from '../layout'
import { getState, updateCurrentUser } from '../state'
import { createElement, createAvatar } from '../utils/dom'
import { navigateTo } from '../router'

export async function renderSettingsPage() {
  const state = getState()
  if (!state.user) {
    renderWithShell(
      createElement('p', {
        className: 'error-text',
        text: 'Необходима авторизация',
      })
    )
    return
  }

  const currentUser = state.user // Store reference to avoid null check issues

  const container = createElement('div', { className: 'page-section narrow' })
  const title = createElement('h1', { text: 'Настройки профиля' })
  const status = createElement('p', { className: 'muted hidden' })

  const form = document.createElement('form')
  form.innerHTML = `
    <label class="field">
      <span>Имя пользователя</span>
      <input type="text" name="username" value="${currentUser.username}" required />
    </label>
    <label class="field">
      <span>Avatar URL</span>
      <input type="url" name="avatar" value="${currentUser.avatar ?? ''}" placeholder="https://..." />
    </label>
    <label class="field">
      <span>Био</span>
      <textarea name="bio" rows="4" placeholder="Расскажите о себе">${currentUser.bio ?? ''}</textarea>
    </label>
    <button type="submit" class="primary-button">Сохранить</button>
  `

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    status.className = 'muted'
    const payload = {
      username: String(formData.get('username') || '').trim(),
      avatar: String(formData.get('avatar') || '').trim() || undefined,
      bio: String(formData.get('bio') || '').trim() || undefined,
    }
    try {
      const updated = await updateProfile(currentUser.id, payload)
      updateCurrentUser(updated)
      status.textContent = 'Данные сохранены'
    } catch (error) {
      status.className = 'error-text'
      status.textContent = (error as Error).message
    }
    status.classList.remove('hidden')
  })

  container.appendChild(title)
  container.appendChild(status)
  container.appendChild(form)

  // Секция подписок
  const followingSection = createElement('div', { className: 'following-section' })
  const followingTitle = createElement('h2', { text: 'Мои подписки' })
  followingSection.appendChild(followingTitle)
  
  const followingList = createElement('div', { className: 'following-list' })
  followingList.appendChild(createElement('p', { className: 'muted', text: 'Загрузка...' }))
  followingSection.appendChild(followingList)
  
  container.appendChild(followingSection)

  renderWithShell(container)

  // Загружаем подписки
  try {
    const following = await fetchFollowing(currentUser.id)
    followingList.innerHTML = ''
    if (following.length === 0) {
      followingList.appendChild(
        createElement('p', { className: 'muted', text: 'Вы пока ни на кого не подписаны' })
      )
    } else {
      following.forEach((user) => {
        const item = createElement('div', { className: 'following-item' })
        const avatar = createAvatar(user.username, user.avatar, 'small')
        avatar.style.cursor = 'pointer'
        avatar.addEventListener('click', () => navigateTo(`#/profile/${user.id}`))
        
        const info = createElement('div', { className: 'following-info' })
        const username = createElement('strong', { text: user.username })
        username.style.cursor = 'pointer'
        username.addEventListener('click', () => navigateTo(`#/profile/${user.id}`))
        info.appendChild(username)
        
        const unfollowBtn = createElement('button', {
          className: 'ghost-button tiny',
          text: 'Отписаться',
        })
        unfollowBtn.addEventListener('click', async () => {
          if (confirm(`Отписаться от ${user.username}?`)) {
            try {
              await unfollowUser(user.id)
              renderSettingsPage() // Перезагружаем страницу
            } catch (error) {
              alert((error as Error).message)
            }
          }
        })
        
        item.appendChild(avatar)
        item.appendChild(info)
        item.appendChild(unfollowBtn)
        followingList.appendChild(item)
      })
    }
  } catch (error) {
    followingList.innerHTML = ''
    followingList.appendChild(
      createElement('p', {
        className: 'error-text',
        text: (error as Error).message,
      })
    )
  }
}



import { updateProfile } from '../api'
import { renderWithShell } from '../layout'
import { getState, updateCurrentUser } from '../state'
import { createElement } from '../utils/dom'

export function renderSettingsPage() {
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

  renderWithShell(container)
}



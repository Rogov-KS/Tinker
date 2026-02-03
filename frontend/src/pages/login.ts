import { login } from '../api'
import { renderWithShell } from '../layout'
import { navigateTo } from '../router'
import { setAuth } from '../state'
import { createElement } from '../utils/dom'

export function renderLoginPage() {
  const container = createElement('div', { className: 'auth-page' })
  const card = createElement('div', { className: 'auth-card' })
  
  const logoContainer = createElement('div', { className: 'auth-logo' })
  const logo = createElement('img', { 
    attrs: {
      src: '/logo.png',
      alt: 'T-News Logo'
    }
  })
  logoContainer.appendChild(logo)
  
  const title = createElement('h1', { text: 'Вход' })
  const errorBox = createElement('p', { className: 'error-text hidden' })

  const form = document.createElement('form')
  form.innerHTML = `
    <label class="field">
      <span>Логин</span>
      <input type="text" name="username" required />
    </label>
    <label class="field">
      <span>Пароль</span>
      <input type="password" name="password" required />
    </label>
    <div class="auth-actions">
      <button type="button" class="ghost-button" id="signupLink">Зарегистрироваться</button>
      <button type="submit" class="primary-button">Войти</button>
    </div>
  `

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    const username = String(formData.get('username') || '').trim()
    const password = String(formData.get('password') || '').trim()
    errorBox.classList.add('hidden')

    try {
      const { token, user } = await login(username, password)
      setAuth(token, user)
      navigateTo('#/feed')
    } catch (error) {
      errorBox.textContent = (error as Error).message
      errorBox.classList.remove('hidden')
    }
  })

  form
    .querySelector('#signupLink')
    ?.addEventListener('click', () => navigateTo('#/signup'))

  card.appendChild(logoContainer)
  card.appendChild(title)
  card.appendChild(errorBox)
  card.appendChild(form)
  container.appendChild(card)

  renderWithShell(container, { showHeader: false })
}



import { signup, login } from '../api'
import { renderWithShell } from '../layout'
import { navigateTo } from '../router'
import { setAuth } from '../state'
import { createElement } from '../utils/dom'

export function renderSignupPage() {
  const container = createElement('div', { className: 'auth-page' })
  const card = createElement('div', { className: 'auth-card' })
  const title = createElement('h1', { text: 'Регистрация' })
  const errorBox = createElement('p', { className: 'error-text hidden' })
  const successBox = createElement('p', { className: 'success-text hidden' })

  const form = document.createElement('form')
  form.innerHTML = `
    <label class="field">
      <span>Логин</span>
      <input type="text" name="username" required minlength="3" />
    </label>
    <label class="field">
      <span>Пароль</span>
      <input type="password" name="password" required minlength="6" />
    </label>
    <label class="field">
      <span>Повторите пароль</span>
      <input type="password" name="confirmPassword" required minlength="6" />
    </label>
    <div class="auth-actions">
      <button type="button" class="ghost-button" id="loginLink">Войти</button>
      <button type="submit" class="primary-button">Зарегистрироваться</button>
    </div>
  `

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const formData = new FormData(form)
    const username = String(formData.get('username') || '').trim()
    const password = String(formData.get('password') || '').trim()
    const confirmPassword = String(
      formData.get('confirmPassword') || ''
    ).trim()

    errorBox.classList.add('hidden')
    successBox.classList.add('hidden')

    if (password !== confirmPassword) {
      errorBox.textContent = 'Пароли не совпадают'
      errorBox.classList.remove('hidden')
      return
    }

    try {
      await signup(username, password)
      // Автоматически входим после успешной регистрации
      try {
        const { token, user } = await login(username, password)
        setAuth(token, user)
        navigateTo('#/feed')
      } catch (loginError) {
        // Если автоматический вход не удался, показываем сообщение об успехе
        successBox.textContent = 'Успешно! Теперь войдите в систему.'
        successBox.classList.remove('hidden')
        setTimeout(() => navigateTo('#/login'), 800)
      }
    } catch (error) {
      errorBox.textContent = (error as Error).message
      errorBox.classList.remove('hidden')
    }
  })

  form
    .querySelector('#loginLink')
    ?.addEventListener('click', () => navigateTo('#/login'))

  card.appendChild(title)
  card.appendChild(errorBox)
  card.appendChild(successBox)
  card.appendChild(form)
  container.appendChild(card)

  renderWithShell(container, { showHeader: false })
}



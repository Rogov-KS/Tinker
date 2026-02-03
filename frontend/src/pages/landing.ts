import { renderWithShell } from '../layout'
import { navigateTo } from '../router'
import { createElement } from '../utils/dom'

export function renderLandingPage() {
  const container = createElement('div', { className: 'landing-page' })
  const card = createElement('div', { className: 'landing-card' })
  const title = createElement('h1', { text: 'T-News' })
  const subtitle = createElement('p', {
    text: 'Читайте и делитесь короткими постами, комментируйте и подписывайтесь на друзей.',
  })

  const actions = createElement('div', { className: 'landing-actions' })
  const loginBtn = createElement('button', {
    className: 'primary-button',
    text: 'Войти',
  })
  loginBtn.addEventListener('click', () => navigateTo('#/login'))

  const signupBtn = createElement('button', {
    className: 'ghost-button',
    text: 'Зарегистрироваться',
  })
  signupBtn.addEventListener('click', () => navigateTo('#/signup'))

  actions.appendChild(loginBtn)
  actions.appendChild(signupBtn)

  card.appendChild(title)
  card.appendChild(subtitle)
  card.appendChild(actions)
  container.appendChild(card)

  renderWithShell(container, { showHeader: false })
}



import { createPost } from '../api'
import { renderWithShell } from '../layout'
import { getState } from '../state'
import { createPostCard } from '../components/post-card'
import { fetchFeed } from '../api'
import { createElement } from '../utils/dom'

export async function renderFeedPage() {
  const container = createElement('div', { className: 'page-section' })
  const title = createElement('h1', { text: 'Ваша лента' })
  container.appendChild(title)

  const state = getState()
  if (!state.user) {
    renderWithShell(createElement('p', { className: 'error-text', text: 'Необходима авторизация' }))
    return
  }

  const composer = createPostComposer(state.user.id, async () => {
    await renderFeedPage()
  })
  container.appendChild(composer)

  const list = createElement('div', { className: 'post-list' })
  container.appendChild(list)

  renderWithShell(container)

  try {
    const posts = await fetchFeed()
    if (!posts.length) {
      list.appendChild(
        createElement('p', {
          className: 'muted',
          text: 'Подпишитесь на кого-нибудь, чтобы увидеть посты.',
        })
      )
      return
    }

    posts.forEach((post) => {
      list.appendChild(
        createPostCard(post, {
          onRefresh: renderFeedPage,
        })
      )
    })
  } catch (error) {
    list.appendChild(
      createElement('p', {
        className: 'error-text',
        text: (error as Error).message,
      })
    )
  }
}

function createPostComposer(userId: string, onSuccess: () => void) {
  const state = getState()
  const card = createElement('div', { className: 'composer-card' })
  const form = document.createElement('form')
  form.innerHTML = `
    <label class="field">
      <span>Новый пост</span>
      <textarea name="content" rows="3" placeholder="Поделитесь мыслями, ${state.user?.username || ''}" required></textarea>
    </label>
    <div class="right">
      <button type="submit" class="primary-button">Опубликовать</button>
    </div>
  `

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const textarea = form.querySelector('textarea[name="content"]') as HTMLTextAreaElement
    const content = textarea.value.trim()
    if (!content) return
    await createPost(userId, content)
    textarea.value = ''
    onSuccess()
  })

  card.appendChild(form)
  return card
}



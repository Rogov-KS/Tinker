import {
  createPost,
  fetchProfile,
  fetchUserPosts,
  followUser,
  unfollowUser,
  fetchFollowing,
} from '../api'
import { renderWithShell } from '../layout'
import { getState } from '../state'
import { createPostCard } from '../components/post-card'
import { createElement, createAvatar } from '../utils/dom'

export async function renderProfilePage(userId?: string) {
  const state = getState()
  const targetId = userId ?? state.user?.id

  if (!targetId) {
    renderWithShell(
      createElement('p', {
        className: 'error-text',
        text: 'Пользователь не найден',
      })
    )
    return
  }

  const container = createElement('div', { className: 'page-section' })
  container.appendChild(createElement('p', { text: 'Загрузка профиля...' }))
  renderWithShell(container)

  try {
    const [profile, posts, following] = await Promise.all([
      fetchProfile(targetId),
      fetchUserPosts(targetId),
      state.token && state.user ? fetchFollowing(state.user.id) : Promise.resolve([]),
    ])

    const isOwnProfile = state.user?.id === profile.id
    const isFollowing = !!following.find((user) => user.id === profile.id)

    renderProfileContent(container, profile, posts, {
      isOwnProfile,
      isFollowing,
    })
  } catch (error) {
    container.innerHTML = ''
    container.appendChild(
      createElement('p', {
        className: 'error-text',
        text: (error as Error).message,
      })
    )
  }
}

function renderProfileContent(
  container: HTMLElement,
  profile: Awaited<ReturnType<typeof fetchProfile>>,
  posts: Awaited<ReturnType<typeof fetchUserPosts>>,
  options: { isOwnProfile: boolean; isFollowing: boolean }
) {
  container.innerHTML = ''

  const headerCard = createElement('div', { className: 'profile-card' })
  const avatar = createAvatar(profile.username, profile.avatar, 'large')

  const info = createElement('div', { className: 'profile-info' })
  info.innerHTML = `
    <h1>${profile.username}</h1>
    <p class="muted">${profile.bio || 'Без описания'}</p>
  `

  const actions = createElement('div', { className: 'profile-actions' })

  if (options.isOwnProfile) {
    actions.appendChild(createElement('span', { text: 'Это ваш профиль' }))
  } else if (getState().token) {
    const followBtn = createElement('button', {
      className: options.isFollowing ? 'ghost-button' : 'primary-button',
      text: options.isFollowing ? 'Отписаться' : 'Подписаться',
    })
    followBtn.addEventListener('click', async () => {
      if (options.isFollowing) {
        await unfollowUser(profile.id)
      } else {
        await followUser(profile.id)
      }
      renderProfilePage(profile.id)
    })
    actions.appendChild(followBtn)
  }

  headerCard.appendChild(avatar)
  headerCard.appendChild(info)
  headerCard.appendChild(actions)

  container.appendChild(headerCard)

  if (options.isOwnProfile) {
    container.appendChild(
      createComposer(profile.id, () => renderProfilePage(profile.id))
    )
  }

  const list = createElement('div', { className: 'post-list' })
  if (!posts.length) {
    list.appendChild(
      createElement('p', {
        className: 'muted',
        text: 'Пока нет постов.',
      })
    )
  } else {
    posts.forEach((post) => {
      list.appendChild(
        createPostCard(post, {
          onRefresh: () => renderProfilePage(profile.id),
        })
      )
    })
  }

  container.appendChild(list)
}

function createComposer(userId: string, onSubmit: () => void) {
  const card = createElement('div', { className: 'composer-card' })
  const form = document.createElement('form')
  form.innerHTML = `
    <label class="field">
      <span>Напишите пост</span>
      <textarea name="content" rows="3" placeholder="Поделитесь новостью" required></textarea>
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
    onSubmit()
  })

  card.appendChild(form)
  return card
}



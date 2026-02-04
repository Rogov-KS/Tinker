import {
  createComment,
  deleteComment,
  deletePost,
  fetchComments,
  likePost,
  unlikePost,
} from '../api'
import { getState } from '../state'
import type { Comment, Post } from '../types'
import { createElement, createAvatar } from '../utils/dom'

interface PostCardOptions {
  onRefresh?: () => void
  showActions?: boolean
}

export function createPostCard(
  post: Post,
  options: PostCardOptions = {}
): HTMLElement {
  const state = getState()
  const card = createElement('article', { className: 'post-card' })
  let liked = !!post.likedByCurrentUser
  let likesCount = post.likes ?? 0

  const header = createElement('div', { className: 'post-header' })
  const avatar = createAvatar(post.author.username, post.author.avatar)
  header.appendChild(avatar)

  const info = createElement('div', { className: 'post-header-info' })
  const name = createElement('strong', { text: post.author.username })
  const meta = createElement('span', {
    className: 'muted',
    text: post.createdAt
      ? new Date(post.createdAt).toLocaleString('ru-RU')
      : '',
  })
  info.appendChild(name)
  info.appendChild(meta)
  header.appendChild(info)

  if (state.user?.id === post.userId) {
    const deleteBtn = createElement('button', {
      className: 'ghost-button danger',
      text: 'Удалить',
    })
    deleteBtn.addEventListener('click', async () => {
      if (!confirm('Удалить пост?')) return
      await deletePost(post.id)
      if (options.onRefresh) {
        options.onRefresh()
      } else {
        card.remove()
      }
    })
    header.appendChild(deleteBtn)
  }

  const content = createElement('p', {
    className: 'post-content',
    text: post.content,
  })

  const actions = createElement('div', { className: 'post-actions' })
  const likeButton = createElement('button', {
    className: `icon-button ${liked ? 'active' : ''}`,
  })
  const updateLikeButton = () => {
    likeButton.classList.toggle('active', liked)
    likeButton.innerHTML = `${liked ? '❤️' : '🤍'} ${likesCount}`
  }
  updateLikeButton()

  likeButton.addEventListener('click', async () => {
    try {
      likeButton.disabled = true
      if (liked) {
        await unlikePost(post.id)
        likesCount = Math.max(0, likesCount - 1)
        liked = false
      } else {
        await likePost(post.id)
        likesCount += 1
        liked = true
      }
      updateLikeButton()
      options.onRefresh?.()
    } catch (error) {
      alert((error as Error).message)
    } finally {
      likeButton.disabled = false
    }
  })

  const commentsButton = createElement('button', {
    className: 'ghost-button',
  })
  
  const updateCommentsButton = (count?: number) => {
    const commentCount = count ?? post.comments?.length ?? 0
    commentsButton.textContent = `Комментарии (${commentCount})`
  }
  updateCommentsButton()

  const commentsSection = createElement('div', {
    className: 'comments-section hidden',
  })
  let commentsLoaded = false

  commentsButton.addEventListener('click', async () => {
    const wasHidden = commentsSection.classList.contains('hidden')
    commentsSection.classList.toggle('hidden')
    if (!commentsLoaded && !commentsSection.classList.contains('hidden')) {
      const comments = await fetchComments(post.id)
      updateCommentsButton(comments.length)
      renderComments(commentsSection, comments, post.id, () => {
        // Не сворачиваем блок комментариев при обновлении
        if (!wasHidden) {
          commentsSection.classList.remove('hidden')
        }
        options.onRefresh?.()
      })
      commentsLoaded = true
    }
  })

  actions.appendChild(likeButton)
  actions.appendChild(commentsButton)

  card.appendChild(header)
  card.appendChild(content)
  card.appendChild(actions)
  card.appendChild(commentsSection)

  return card
}

function renderComments(
  container: HTMLElement,
  comments: Comment[],
  postId: string,
  onRefresh?: () => void
) {
  container.innerHTML = ''
  // Убеждаемся, что блок комментариев открыт
  container.classList.remove('hidden')

  const list = createElement('div', { className: 'comments-list' })
  comments.forEach((comment) => {
    list.appendChild(createCommentItem(comment, () => {
      // Update count after deletion
      const count = list.querySelectorAll('.comment-row').length
      const button = container.closest('.post-card')?.querySelector('.ghost-button:not(.danger)') as HTMLElement
      if (button && button.textContent?.startsWith('Комментарии')) {
        button.textContent = `Комментарии (${count})`
      }
      // Не сворачиваем блок при удалении комментария
      container.classList.remove('hidden')
      onRefresh?.()
    }))
  })

  container.appendChild(list)
  container.appendChild(createCommentForm(postId, list, () => {
    // Update count after adding
    const count = list.querySelectorAll('.comment-row').length
    const button = container.closest('.post-card')?.querySelector('.ghost-button:not(.danger)') as HTMLElement
    if (button && button.textContent?.startsWith('Комментарии')) {
      button.textContent = `Комментарии (${count})`
    }
    // Не сворачиваем блок при добавлении комментария
    container.classList.remove('hidden')
    onRefresh?.()
  }))
}

function createCommentItem(comment: Comment, onRefresh?: () => void) {
  const state = getState()
  const row = createElement('div', { className: 'comment-row' })
  const avatar = createAvatar(comment.author.username, comment.author.avatar, 'small')
  const body = createElement('div', { className: 'comment-body' })
  const header = createElement('div', { className: 'comment-header' })
  header.innerHTML = `<strong>${comment.author.username}</strong>
    <span class="muted">${comment.createdAt ? new Date(comment.createdAt).toLocaleString('ru-RU') : ''}</span>`
  const text = createElement('p', { text: comment.content })
  body.appendChild(header)
  body.appendChild(text)

  if (state.user?.id === comment.userId) {
    const del = createElement('button', {
      className: 'ghost-button danger tiny',
      text: 'Удалить',
    })
    del.addEventListener('click', async () => {
      if (!confirm('Удалить комментарий?')) return
      await deleteComment(comment.id)
      if (onRefresh) {
        onRefresh()
      } else {
        row.remove()
      }
    })
    body.appendChild(del)
  }

  row.appendChild(avatar)
  row.appendChild(body)
  return row
}

function createCommentForm(
  postId: string,
  list: HTMLElement,
  onCommentAdded?: () => void
) {
  const form = document.createElement('form')
  form.className = 'comment-form'
  form.innerHTML = `
    <input type="text" name="content" placeholder="Введите свой комментарий" required />
    <button type="submit" class="primary-button">Отправить</button>
  `

  form.addEventListener('submit', async (event) => {
    event.preventDefault()
    const input = form.querySelector('input[name="content"]') as HTMLInputElement
    const content = input.value.trim()
    if (!content) return
    const comment = await createComment(postId, content)
    list.appendChild(createCommentItem(comment, () => {
      // Update count after deletion
      const count = list.querySelectorAll('.comment-row').length
      const button = form.closest('.post-card')?.querySelector('.ghost-button:not(.danger)') as HTMLElement
      if (button && button.textContent?.startsWith('Комментарии')) {
        button.textContent = `Комментарии (${count})`
      }
      onCommentAdded?.()
    }))
    input.value = ''
    // Update count after adding
    onCommentAdded?.()
  })

  return form
}



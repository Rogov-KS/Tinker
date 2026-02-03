import { getState, clearAuth } from './state'
import type { AuthPayload, Comment, Post, User } from './types'

const API_BASE = import.meta.env.VITE_API_URL ?? '/api'

async function apiFetch<T>(
  path: string,
  options: RequestInit = {},
  requiresAuth = true
): Promise<T> {
  const { token } = getState()
  const headers = new Headers(options.headers)

  if (!headers.has('Content-Type') && options.body) {
    headers.set('Content-Type', 'application/json')
  }

  if (token) {
    headers.set('Authorization', `Bearer ${token}`)
  } else if (requiresAuth) {
    throw new Error('Необходима авторизация. Пожалуйста, войдите.')
  }

  const response = await fetch(`${API_BASE}${path}`, {
    ...options,
    headers,
  })

  if (response.status === 401) {
    clearAuth()
    throw new Error('Необходима авторизация. Пожалуйста, войдите заново.')
  }

  if (!response.ok) {
    const errorBody = await response.json().catch(() => ({}))
    const message = errorBody.message
      ? Array.isArray(errorBody.message)
        ? errorBody.message.join(', ')
        : errorBody.message
      : errorBody.error
    throw new Error(message || 'Произошла ошибка')
  }

  if (response.status === 204) {
    return null as T
  }

  return response.json()
}

export function login(username: string, password: string) {
  return apiFetch<AuthPayload>(
    '/auth/login',
    {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    },
    false
  )
}

export function signup(username: string, password: string) {
  return apiFetch<User>(
    '/users',
    {
      method: 'POST',
      body: JSON.stringify({ username, password }),
    },
    false
  )
}

export function fetchFeed() {
  return apiFetch<Post[]>('/feed')
}

export function fetchProfile(userId: string) {
  return apiFetch<User>(
    `/users/${encodeURIComponent(userId)}`,
    {},
    false
  )
}

export function updateProfile(userId: string, data: Partial<User>) {
  return apiFetch<User>(`/users/${encodeURIComponent(userId)}`, {
    method: 'PATCH',
    body: JSON.stringify(data),
  })
}

export function fetchUserPosts(userId: string) {
  return apiFetch<Post[]>(
    `/users/${encodeURIComponent(userId)}/posts`,
    {},
    false
  )
}

export function createPost(userId: string, content: string) {
  return apiFetch<Post>(`/users/${encodeURIComponent(userId)}/posts`, {
    method: 'POST',
    body: JSON.stringify({ content }),
  })
}

export function deletePost(postId: string) {
  return apiFetch<void>(`/posts/${encodeURIComponent(postId)}`, {
    method: 'DELETE',
  })
}

export function likePost(postId: string) {
  return apiFetch<void>(`/posts/${encodeURIComponent(postId)}/likes`, {
    method: 'POST',
  })
}

export function unlikePost(postId: string) {
  return apiFetch<void>(`/posts/${encodeURIComponent(postId)}/likes`, {
    method: 'DELETE',
  })
}

export function fetchComments(postId: string) {
  return apiFetch<Comment[]>(`/posts/${encodeURIComponent(postId)}/comments`, {}, false)
}

export function createComment(postId: string, content: string) {
  return apiFetch<Comment>(`/posts/${encodeURIComponent(postId)}/comments`, {
    method: 'POST',
    body: JSON.stringify({ content }),
  })
}

export function deleteComment(commentId: string) {
  return apiFetch<void>(`/comments/${encodeURIComponent(commentId)}`, {
    method: 'DELETE',
  })
}

export function followUser(userId: string) {
  return apiFetch<void>(`/users/${encodeURIComponent(userId)}/follow`, {
    method: 'POST',
  })
}

export function unfollowUser(userId: string) {
  return apiFetch<void>(`/users/${encodeURIComponent(userId)}/follow`, {
    method: 'DELETE',
  })
}

export function fetchFollowing(userId: string) {
  return apiFetch<User[]>(`/users/${encodeURIComponent(userId)}/following`)
}

export function searchEntities(query: string, type: 'users' | 'posts') {
  const params = new URLSearchParams({ query, type })
  return apiFetch<Array<User | Post>>(
    `/search?${params.toString()}`,
    {},
    false
  )
}

export function getAllUsers() {
  return apiFetch<User[]>('/users', {}, false)
}

export function getAllPosts() {
  return apiFetch<Post[]>('/posts', {}, false)
}



import type { User } from './types'

export interface AppState {
  token: string | null
  user: User | null
}

type Listener = () => void

const STORAGE_KEY = 'tnews-auth'

const listeners: Listener[] = []

const initial = (() => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (!raw) {
      return { token: null, user: null }
    }
    return JSON.parse(raw) as AppState
  } catch {
    return { token: null, user: null }
  }
})()

const state: AppState = {
  token: initial.token,
  user: initial.user,
}

export function getState(): AppState {
  return { ...state }
}

export function setAuth(token: string, user: User) {
  state.token = token
  state.user = user
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  notify()
}

export function updateCurrentUser(user: User) {
  state.user = user
  localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
  notify()
}

export function clearAuth() {
  state.token = null
  state.user = null
  localStorage.removeItem(STORAGE_KEY)
  notify()
}

export function subscribe(listener: Listener) {
  listeners.push(listener)
  return () => {
    const index = listeners.indexOf(listener)
    if (index >= 0) {
      listeners.splice(index, 1)
    }
  }
}

function notify() {
  listeners.forEach((listener) => listener())
}



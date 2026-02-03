export interface User {
  id: string
  username: string
  avatar?: string
  bio?: string
}

export interface PostAuthor {
  id: string
  username: string
  avatar?: string
}

export interface CommentAuthor extends PostAuthor {}

export interface Comment {
  id: string
  userId: string
  content: string
  createdAt?: string
  author: CommentAuthor
}

export interface Post {
  id: string
  userId: string
  author: PostAuthor
  content: string
  likes?: number
  likedByCurrentUser?: boolean
  createdAt?: string
  comments?: Comment[]
}

export interface AuthPayload {
  token: string
  user: User
}

export interface SearchResult {
  id: string
  username?: string
  content?: string
}



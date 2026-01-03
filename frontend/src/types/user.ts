export type UserRole = 'superadmin' | 'user'

export interface User {
  id: string
  email: string
  name: string
  role: UserRole
  avatar_url?: string
  is_active: boolean
  created_at: string
  updated_at: string
}

export interface UpdateUserRequest {
  name?: string
  avatar_url?: string
}

export interface UpdatePasswordRequest {
  current_password: string
  new_password: string
}

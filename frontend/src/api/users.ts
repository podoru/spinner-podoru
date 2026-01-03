import { apiClient } from './client'
import type { ApiResponse, User, UpdateUserRequest, UpdatePasswordRequest } from '@/types'

export const usersApi = {
  getMe: () =>
    apiClient.get<ApiResponse<User>>('/users/me'),

  updateMe: (data: UpdateUserRequest) =>
    apiClient.put<ApiResponse<User>>('/users/me', data),

  updatePassword: (data: UpdatePasswordRequest) =>
    apiClient.put('/users/me/password', data),
}

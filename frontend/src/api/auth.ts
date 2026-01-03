import { apiClient } from './client'
import type { ApiResponse, AuthResponse, TokenResponse, LoginRequest, RegisterRequest, RefreshRequest } from '@/types'

export const authApi = {
  login: (data: LoginRequest) =>
    apiClient.post<ApiResponse<AuthResponse>>('/auth/login', data),

  register: (data: RegisterRequest) =>
    apiClient.post<ApiResponse<AuthResponse>>('/auth/register', data),

  refresh: (data: RefreshRequest) =>
    apiClient.post<ApiResponse<TokenResponse>>('/auth/refresh', data),

  logout: (data: RefreshRequest) =>
    apiClient.post('/auth/logout', data),
}

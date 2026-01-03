import { useMutation, useQueryClient } from '@tanstack/react-query'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/authStore'
import type { LoginRequest, RegisterRequest } from '@/types'

export function useLogin() {
  const queryClient = useQueryClient()
  const { setAuth } = useAuthStore()

  return useMutation({
    mutationFn: (data: LoginRequest) => authApi.login(data),
    onSuccess: (response) => {
      const { user, tokens } = response.data.data!
      setAuth(user, tokens.access_token, tokens.refresh_token)
      queryClient.invalidateQueries({ queryKey: ['user'] })
    },
  })
}

export function useRegister() {
  const { setAuth } = useAuthStore()

  return useMutation({
    mutationFn: (data: RegisterRequest) => authApi.register(data),
    onSuccess: (response) => {
      const { user, tokens } = response.data.data!
      setAuth(user, tokens.access_token, tokens.refresh_token)
    },
  })
}

export function useLogout() {
  const queryClient = useQueryClient()
  const { refreshToken, logout } = useAuthStore()

  return useMutation({
    mutationFn: () => authApi.logout({ refresh_token: refreshToken! }),
    onSettled: () => {
      logout()
      queryClient.clear()
    },
  })
}

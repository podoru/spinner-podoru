import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
})

// We'll set up the auth interceptors after importing the auth store
// to avoid circular dependencies
let getAccessToken: (() => string | null) | null = null
let getRefreshToken: (() => string | null) | null = null
let setTokens: ((access: string, refresh: string) => void) | null = null
let logout: (() => void) | null = null

export function setupAuthInterceptors(
  accessTokenGetter: () => string | null,
  refreshTokenGetter: () => string | null,
  tokensSetter: (access: string, refresh: string) => void,
  logoutFn: () => void
) {
  getAccessToken = accessTokenGetter
  getRefreshToken = refreshTokenGetter
  setTokens = tokensSetter
  logout = logoutFn
}

// Request interceptor - add auth token
apiClient.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = getAccessToken?.()
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => Promise.reject(error)
)

// Response interceptor - handle token refresh
apiClient.interceptors.response.use(
  (response) => response,
  async (error: AxiosError) => {
    const originalRequest = error.config as InternalAxiosRequestConfig & { _retry?: boolean }

    if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
      originalRequest._retry = true

      try {
        const refreshToken = getRefreshToken?.()
        if (refreshToken && setTokens) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken,
          })

          const { access_token, refresh_token } = response.data.data
          setTokens(access_token, refresh_token)

          originalRequest.headers.Authorization = `Bearer ${access_token}`
          return apiClient(originalRequest)
        }
      } catch {
        logout?.()
        window.location.href = '/login'
      }
    }

    return Promise.reject(error)
  }
)

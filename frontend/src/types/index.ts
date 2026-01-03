// API Response types
export interface ApiResponse<T> {
  success: boolean
  data?: T
  error?: ErrorInfo
  meta?: Meta
}

export interface ErrorInfo {
  code: string
  message: string
  details?: Record<string, string>
}

export interface Meta {
  page?: number
  per_page?: number
  total?: number
  total_pages?: number
}

export interface MessageResponse {
  message: string
}

// Re-export all types
export * from './auth'
export * from './user'
export * from './team'
export * from './project'
export * from './service'
export * from './domain'

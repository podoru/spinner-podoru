export interface Domain {
  id: string
  service_id: string
  domain: string
  ssl_enabled: boolean
  ssl_auto: boolean
  created_at: string
}

export interface CreateDomainRequest {
  domain: string
  ssl_enabled?: boolean
  ssl_auto?: boolean
}

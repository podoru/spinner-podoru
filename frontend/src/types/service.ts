export type DeployType = 'image' | 'dockerfile' | 'compose'
export type ServiceStatus = 'stopped' | 'running' | 'deploying' | 'failed'
export type RestartPolicy = 'no' | 'always' | 'on-failure' | 'unless-stopped'

export interface Service {
  id: string
  project_id: string
  name: string
  slug: string
  deploy_type: DeployType
  image?: string
  dockerfile_path: string
  build_context: string
  replicas: number
  cpu_limit?: number
  memory_limit?: number
  health_check_path?: string
  health_check_interval: number
  restart_policy: RestartPolicy
  status: ServiceStatus
  container_id?: string
  swarm_service_id?: string
  created_at: string
  updated_at: string
}

export interface EnvVar {
  key: string
  value: string
}

export interface CreateServiceRequest {
  name: string
  slug: string
  deploy_type: DeployType
  image?: string
  dockerfile_path?: string
  build_context?: string
  compose_file?: string
  env_vars?: EnvVar[]
  replicas?: number
  cpu_limit?: number
  memory_limit?: number
  health_check_path?: string
  health_check_interval?: number
  restart_policy?: RestartPolicy
}

export interface UpdateServiceRequest {
  name?: string
  image?: string
  dockerfile_path?: string
  build_context?: string
  env_vars?: EnvVar[]
  replicas?: number
  cpu_limit?: number
  memory_limit?: number
  health_check_path?: string
  health_check_interval?: number
  restart_policy?: RestartPolicy
}

export interface ScaleServiceRequest {
  replicas: number
}

export interface ServiceLogsResponse {
  service_id: string
  logs: string
  timestamp: string
}

export interface DeploymentResponse {
  id: string
  service_id: string
  triggered_by?: string
  commit_sha?: string
  commit_message?: string
  status: string
  logs?: string
  started_at: string
  finished_at?: string
}

import { apiClient } from './client'
import type {
  ApiResponse, Project, Service,
  UpdateProjectRequest, CreateServiceRequest
} from '@/types'

export const projectsApi = {
  get: (projectId: string) =>
    apiClient.get<ApiResponse<Project>>(`/projects/${projectId}`),

  update: (projectId: string, data: UpdateProjectRequest) =>
    apiClient.put<ApiResponse<Project>>(`/projects/${projectId}`, data),

  delete: (projectId: string) =>
    apiClient.delete(`/projects/${projectId}`),

  // Services in project
  listServices: (projectId: string) =>
    apiClient.get<ApiResponse<Service[]>>(`/projects/${projectId}/services`),

  createService: (projectId: string, data: CreateServiceRequest) =>
    apiClient.post<ApiResponse<Service>>(`/projects/${projectId}/services`, data),
}

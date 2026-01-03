import { apiClient } from './client'
import type {
  ApiResponse, Service, Domain, DeploymentResponse, ServiceLogsResponse,
  UpdateServiceRequest, ScaleServiceRequest, CreateDomainRequest, MessageResponse
} from '@/types'

export const servicesApi = {
  get: (serviceId: string) =>
    apiClient.get<ApiResponse<Service>>(`/services/${serviceId}`),

  update: (serviceId: string, data: UpdateServiceRequest) =>
    apiClient.put<ApiResponse<Service>>(`/services/${serviceId}`, data),

  delete: (serviceId: string) =>
    apiClient.delete(`/services/${serviceId}`),

  // Service actions
  deploy: (serviceId: string) =>
    apiClient.post<ApiResponse<DeploymentResponse>>(`/services/${serviceId}/deploy`),

  start: (serviceId: string) =>
    apiClient.post<ApiResponse<MessageResponse>>(`/services/${serviceId}/start`),

  stop: (serviceId: string) =>
    apiClient.post<ApiResponse<MessageResponse>>(`/services/${serviceId}/stop`),

  restart: (serviceId: string) =>
    apiClient.post<ApiResponse<MessageResponse>>(`/services/${serviceId}/restart`),

  scale: (serviceId: string, data: ScaleServiceRequest) =>
    apiClient.post<ApiResponse<Service>>(`/services/${serviceId}/scale`, data),

  getLogs: (serviceId: string, params?: { tail?: number; since?: string }) =>
    apiClient.get<ApiResponse<ServiceLogsResponse>>(`/services/${serviceId}/logs`, { params }),

  // Domains
  listDomains: (serviceId: string) =>
    apiClient.get<ApiResponse<Domain[]>>(`/services/${serviceId}/domains`),

  addDomain: (serviceId: string, data: CreateDomainRequest) =>
    apiClient.post<ApiResponse<Domain>>(`/services/${serviceId}/domains`, data),

  deleteDomain: (serviceId: string, domainId: string) =>
    apiClient.delete(`/services/${serviceId}/domains/${domainId}`),
}

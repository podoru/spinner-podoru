import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { projectsApi } from '@/api/projects'
import { servicesApi } from '@/api/services'
import type { CreateServiceRequest, UpdateServiceRequest, ScaleServiceRequest } from '@/types'

export function useProjectServices(projectId: string) {
  return useQuery({
    queryKey: ['projects', projectId, 'services'],
    queryFn: async () => {
      const response = await projectsApi.listServices(projectId)
      return response.data.data!
    },
    enabled: !!projectId,
  })
}

export function useService(serviceId: string) {
  return useQuery({
    queryKey: ['services', serviceId],
    queryFn: async () => {
      const response = await servicesApi.get(serviceId)
      return response.data.data!
    },
    enabled: !!serviceId,
    refetchInterval: 10000, // Poll every 10 seconds for status updates
  })
}

export function useCreateService(projectId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateServiceRequest) => projectsApi.createService(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', projectId, 'services'] })
    },
  })
}

export function useUpdateService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: UpdateServiceRequest) => servicesApi.update(serviceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

export function useDeleteService() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (serviceId: string) => servicesApi.delete(serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services'] })
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

// Service Actions
export function useDeployService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => servicesApi.deploy(serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

export function useStartService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => servicesApi.start(serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

export function useStopService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => servicesApi.stop(serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

export function useRestartService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: () => servicesApi.restart(serviceId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

export function useScaleService(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: ScaleServiceRequest) => servicesApi.scale(serviceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId] })
    },
  })
}

// Service Logs
export function useServiceLogs(
  serviceId: string,
  options?: { tail?: number; since?: string; enabled?: boolean }
) {
  return useQuery({
    queryKey: ['services', serviceId, 'logs', options?.tail, options?.since],
    queryFn: async () => {
      const response = await servicesApi.getLogs(serviceId, {
        tail: options?.tail,
        since: options?.since,
      })
      return response.data.data!
    },
    enabled: options?.enabled !== false && !!serviceId,
    refetchInterval: 5000, // Poll every 5 seconds for new logs
  })
}

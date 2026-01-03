import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { servicesApi } from '@/api/services'
import type { CreateDomainRequest } from '@/types'

export function useServiceDomains(serviceId: string) {
  return useQuery({
    queryKey: ['services', serviceId, 'domains'],
    queryFn: async () => {
      const response = await servicesApi.listDomains(serviceId)
      return response.data.data!
    },
    enabled: !!serviceId,
  })
}

export function useAddDomain(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateDomainRequest) => servicesApi.addDomain(serviceId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId, 'domains'] })
    },
  })
}

export function useDeleteDomain(serviceId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (domainId: string) => servicesApi.deleteDomain(serviceId, domainId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['services', serviceId, 'domains'] })
    },
  })
}

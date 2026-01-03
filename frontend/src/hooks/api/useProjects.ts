import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { teamsApi } from '@/api/teams'
import { projectsApi } from '@/api/projects'
import type { CreateProjectRequest, UpdateProjectRequest } from '@/types'

export function useTeamProjects(teamId: string) {
  return useQuery({
    queryKey: ['teams', teamId, 'projects'],
    queryFn: async () => {
      const response = await teamsApi.listProjects(teamId)
      return response.data.data!
    },
    enabled: !!teamId,
  })
}

export function useProject(projectId: string) {
  return useQuery({
    queryKey: ['projects', projectId],
    queryFn: async () => {
      const response = await projectsApi.get(projectId)
      return response.data.data!
    },
    enabled: !!projectId,
  })
}

export function useCreateProject(teamId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: CreateProjectRequest) => teamsApi.createProject(teamId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams', teamId, 'projects'] })
    },
  })
}

export function useUpdateProject(projectId: string) {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (data: UpdateProjectRequest) => projectsApi.update(projectId, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects', projectId] })
    },
  })
}

export function useDeleteProject() {
  const queryClient = useQueryClient()

  return useMutation({
    mutationFn: (projectId: string) => projectsApi.delete(projectId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['teams'] })
      queryClient.invalidateQueries({ queryKey: ['projects'] })
    },
  })
}

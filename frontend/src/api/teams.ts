import { apiClient } from './client'
import type {
  ApiResponse, Team, TeamWithRole, TeamMember, Project,
  CreateTeamRequest, UpdateTeamRequest,
  AddTeamMemberRequest, UpdateTeamMemberRequest,
  CreateProjectRequest
} from '@/types'

export const teamsApi = {
  list: () =>
    apiClient.get<ApiResponse<TeamWithRole[]>>('/teams'),

  get: (teamId: string) =>
    apiClient.get<ApiResponse<Team>>(`/teams/${teamId}`),

  create: (data: CreateTeamRequest) =>
    apiClient.post<ApiResponse<Team>>('/teams', data),

  update: (teamId: string, data: UpdateTeamRequest) =>
    apiClient.put<ApiResponse<Team>>(`/teams/${teamId}`, data),

  delete: (teamId: string) =>
    apiClient.delete(`/teams/${teamId}`),

  // Members
  listMembers: (teamId: string) =>
    apiClient.get<ApiResponse<TeamMember[]>>(`/teams/${teamId}/members`),

  addMember: (teamId: string, data: AddTeamMemberRequest) =>
    apiClient.post<ApiResponse<TeamMember>>(`/teams/${teamId}/members`, data),

  updateMember: (teamId: string, userId: string, data: UpdateTeamMemberRequest) =>
    apiClient.put<ApiResponse<TeamMember>>(`/teams/${teamId}/members/${userId}`, data),

  removeMember: (teamId: string, userId: string) =>
    apiClient.delete(`/teams/${teamId}/members/${userId}`),

  // Projects in team
  listProjects: (teamId: string) =>
    apiClient.get<ApiResponse<Project[]>>(`/teams/${teamId}/projects`),

  createProject: (teamId: string, data: CreateProjectRequest) =>
    apiClient.post<ApiResponse<Project>>(`/teams/${teamId}/projects`, data),
}

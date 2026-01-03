import type { User } from './user'

export type TeamRole = 'owner' | 'admin' | 'member'

export interface Team {
  id: string
  name: string
  slug: string
  description?: string
  owner_id: string
  created_at: string
  updated_at: string
}

export interface TeamWithRole extends Team {
  role: TeamRole
}

export interface TeamMember {
  id: string
  team_id: string
  user_id: string
  role: TeamRole
  user?: User
  created_at: string
}

export interface CreateTeamRequest {
  name: string
  slug: string
  description?: string
}

export interface UpdateTeamRequest {
  name?: string
  description?: string
}

export interface AddTeamMemberRequest {
  email: string
  role: 'admin' | 'member'
}

export interface UpdateTeamMemberRequest {
  role: 'admin' | 'member'
}

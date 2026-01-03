export interface Project {
  id: string
  team_id: string
  name: string
  slug: string
  description?: string
  github_repo?: string
  github_branch: string
  auto_deploy: boolean
  created_at: string
  updated_at: string
}

export interface CreateProjectRequest {
  name: string
  slug: string
  description?: string
  github_repo?: string
  github_branch?: string
  github_token?: string
  auto_deploy?: boolean
}

export interface UpdateProjectRequest {
  name?: string
  description?: string
  github_repo?: string
  github_branch?: string
  github_token?: string
  auto_deploy?: boolean
}

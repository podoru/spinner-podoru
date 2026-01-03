import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { ArrowLeft, Plus, Server, Settings, Trash2, Github, Pencil } from 'lucide-react'
import { Breadcrumb } from '@/components/common/Breadcrumb'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { EmptyState } from '@/components/common/EmptyState'
import { StatusBadge } from '@/components/common/StatusBadge'
import { ConfirmDialog } from '@/components/common/ConfirmDialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { useProject, useDeleteProject } from '@/hooks/api/useProjects'
import { useProjectServices } from '@/hooks/api/useServices'
import { useTeam } from '@/hooks/api/useTeams'
import { toast } from 'sonner'

export function ProjectDetailPage() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const { data: project, isLoading: projectLoading } = useProject(projectId!)
  const { data: team } = useTeam(project?.team_id || '')
  const { data: services, isLoading: servicesLoading } = useProjectServices(projectId!)
  const deleteProject = useDeleteProject()

  const handleDelete = async () => {
    try {
      await deleteProject.mutateAsync(projectId!)
      toast.success('Project deleted successfully')
      navigate('/teams')
    } catch {
      toast.error('Failed to delete project')
    }
  }

  if (projectLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!project) {
    return (
      <EmptyState
        title="Project not found"
        description="The project you're looking for doesn't exist."
      >
        <Button asChild>
          <Link to="/teams">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Teams
          </Link>
        </Button>
      </EmptyState>
    )
  }

  return (
    <div className="space-y-6">
      <Breadcrumb
        items={[
          { label: 'Teams', href: '/teams' },
          { label: team?.name || 'Team', href: `/teams/${project.team_id}` },
          { label: project.name },
        ]}
      />

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to={`/teams/${project.team_id}`}>
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">{project.name}</h1>
            <p className="text-muted-foreground">{project.description || project.slug}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link to={`/projects/${projectId}/edit`}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </Link>
          </Button>
          <Button variant="destructive" size="sm" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="mr-2 h-4 w-4" />
            Delete
          </Button>
        </div>
      </div>

      <Tabs defaultValue="services">
        <TabsList>
          <TabsTrigger value="services">
            <Server className="mr-2 h-4 w-4" />
            Services
          </TabsTrigger>
          <TabsTrigger value="settings">
            <Settings className="mr-2 h-4 w-4" />
            Settings
          </TabsTrigger>
        </TabsList>

        <TabsContent value="services" className="mt-6">
          <div className="flex justify-end mb-4">
            <Button asChild>
              <Link to={`/projects/${projectId}/services/new`}>
                <Plus className="mr-2 h-4 w-4" />
                Create Service
              </Link>
            </Button>
          </div>

          {servicesLoading ? (
            <LoadingSpinner className="h-32" />
          ) : services?.length === 0 ? (
            <EmptyState
              icon={Server}
              title="No services yet"
              description="Create your first service in this project."
            >
              <Button asChild>
                <Link to={`/projects/${projectId}/services/new`}>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Service
                </Link>
              </Button>
            </EmptyState>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {services?.map((service) => (
                <Link key={service.id} to={`/services/${service.id}`}>
                  <Card className="hover:bg-muted/50 transition-colors">
                    <CardHeader className="pb-3">
                      <div className="flex items-start justify-between">
                        <div>
                          <CardTitle className="text-lg">{service.name}</CardTitle>
                          <CardDescription>{service.slug}</CardDescription>
                        </div>
                        <StatusBadge status={service.status} />
                      </div>
                    </CardHeader>
                    <CardContent>
                      <div className="flex items-center gap-2 text-sm text-muted-foreground">
                        <Badge variant="outline">{service.deploy_type}</Badge>
                        {service.image && (
                          <span className="truncate">{service.image}</span>
                        )}
                      </div>
                    </CardContent>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="settings" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Project Settings</CardTitle>
              <CardDescription>Configure your project settings</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="font-medium">Name</h4>
                <p className="text-sm text-muted-foreground">{project.name}</p>
              </div>
              <div>
                <h4 className="font-medium">Slug</h4>
                <p className="text-sm text-muted-foreground">{project.slug}</p>
              </div>
              {project.github_repo && (
                <div>
                  <h4 className="font-medium flex items-center gap-2">
                    <Github className="h-4 w-4" />
                    GitHub Repository
                  </h4>
                  <p className="text-sm text-muted-foreground">{project.github_repo}</p>
                  <p className="text-sm text-muted-foreground">Branch: {project.github_branch}</p>
                </div>
              )}
              <div>
                <h4 className="font-medium">Auto Deploy</h4>
                <p className="text-sm text-muted-foreground">
                  {project.auto_deploy ? 'Enabled' : 'Disabled'}
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Delete Project"
        description="Are you sure you want to delete this project? This action cannot be undone and will delete all services."
        confirmText="Delete"
        variant="destructive"
        isLoading={deleteProject.isPending}
        onConfirm={handleDelete}
      />
    </div>
  )
}

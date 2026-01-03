import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import { ArrowLeft, Plus, Users, FolderKanban, Trash2, Pencil, UserPlus } from 'lucide-react'
import { Breadcrumb } from '@/components/common/Breadcrumb'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { EmptyState } from '@/components/common/EmptyState'
import { ConfirmDialog } from '@/components/common/ConfirmDialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { useTeam, useDeleteTeam } from '@/hooks/api/useTeams'
import { useTeamMembers } from '@/hooks/api/useTeamMembers'
import { useTeamProjects } from '@/hooks/api/useProjects'
import { toast } from 'sonner'

export function TeamDetailPage() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const { data: team, isLoading: teamLoading } = useTeam(teamId!)
  const { data: members, isLoading: membersLoading } = useTeamMembers(teamId!)
  const { data: projects, isLoading: projectsLoading } = useTeamProjects(teamId!)
  const deleteTeam = useDeleteTeam()

  const handleDelete = async () => {
    try {
      await deleteTeam.mutateAsync(teamId!)
      toast.success('Team deleted successfully')
      navigate('/teams')
    } catch {
      toast.error('Failed to delete team')
    }
  }

  if (teamLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!team) {
    return (
      <EmptyState
        title="Team not found"
        description="The team you're looking for doesn't exist."
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
          { label: team.name },
        ]}
      />

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to="/teams">
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <h1 className="text-2xl font-bold tracking-tight">{team.name}</h1>
            <p className="text-muted-foreground">{team.description || team.slug}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link to={`/teams/${teamId}/edit`}>
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

      <Tabs defaultValue="projects">
        <TabsList>
          <TabsTrigger value="projects">
            <FolderKanban className="mr-2 h-4 w-4" />
            Projects
          </TabsTrigger>
          <TabsTrigger value="members">
            <Users className="mr-2 h-4 w-4" />
            Members
          </TabsTrigger>
        </TabsList>

        <TabsContent value="projects" className="mt-6">
          <div className="flex justify-end mb-4">
            <Button asChild>
              <Link to={`/teams/${teamId}/projects/new`}>
                <Plus className="mr-2 h-4 w-4" />
                Create Project
              </Link>
            </Button>
          </div>

          {projectsLoading ? (
            <LoadingSpinner className="h-32" />
          ) : projects?.length === 0 ? (
            <EmptyState
              icon={FolderKanban}
              title="No projects yet"
              description="Create your first project in this team."
            >
              <Button asChild>
                <Link to={`/teams/${teamId}/projects/new`}>
                  <Plus className="mr-2 h-4 w-4" />
                  Create Project
                </Link>
              </Button>
            </EmptyState>
          ) : (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {projects?.map((project) => (
                <Link key={project.id} to={`/projects/${project.id}`}>
                  <Card className="hover:bg-muted/50 transition-colors">
                    <CardHeader>
                      <CardTitle className="text-lg">{project.name}</CardTitle>
                      <CardDescription>{project.slug}</CardDescription>
                    </CardHeader>
                    <CardContent>
                      {project.description && (
                        <p className="text-sm text-muted-foreground line-clamp-2">
                          {project.description}
                        </p>
                      )}
                    </CardContent>
                  </Card>
                </Link>
              ))}
            </div>
          )}
        </TabsContent>

        <TabsContent value="members" className="mt-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Team Members</CardTitle>
                <CardDescription>Manage who has access to this team</CardDescription>
              </div>
              <Button asChild size="sm">
                <Link to={`/teams/${teamId}/members/add`}>
                  <UserPlus className="mr-2 h-4 w-4" />
                  Add Member
                </Link>
              </Button>
            </CardHeader>
            <CardContent>
              {membersLoading ? (
                <LoadingSpinner className="h-32" />
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Member</TableHead>
                      <TableHead>Role</TableHead>
                      <TableHead>Joined</TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {members?.map((member) => (
                      <TableRow key={member.id}>
                        <TableCell>
                          <div className="flex items-center gap-3">
                            <Avatar className="h-8 w-8">
                              <AvatarImage src={member.user?.avatar_url} />
                              <AvatarFallback>
                                {member.user?.name?.charAt(0).toUpperCase() || 'U'}
                              </AvatarFallback>
                            </Avatar>
                            <div>
                              <p className="font-medium">{member.user?.name}</p>
                              <p className="text-sm text-muted-foreground">{member.user?.email}</p>
                            </div>
                          </div>
                        </TableCell>
                        <TableCell>
                          <Badge variant="outline" className="capitalize">
                            {member.role}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {new Date(member.created_at).toLocaleDateString()}
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Delete Team"
        description="Are you sure you want to delete this team? This action cannot be undone and will delete all projects and services."
        confirmText="Delete"
        variant="destructive"
        isLoading={deleteTeam.isPending}
        onConfirm={handleDelete}
      />
    </div>
  )
}

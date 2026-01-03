import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { FormPageLayout } from '@/components/common/FormPageLayout'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { Switch } from '@/components/ui/switch'
import { useTeam } from '@/hooks/api/useTeams'
import { useCreateProject } from '@/hooks/api/useProjects'
import { toast } from 'sonner'

const createProjectSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  slug: z.string().min(2, 'Slug must be at least 2 characters').regex(/^[a-z0-9-]+$/, 'Slug must be lowercase alphanumeric with hyphens'),
  description: z.string().optional(),
  github_repo: z.string().optional(),
  github_branch: z.string().optional(),
  auto_deploy: z.boolean().optional(),
})

type CreateProjectFormData = z.infer<typeof createProjectSchema>

export function CreateProjectPage() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const { data: team, isLoading } = useTeam(teamId!)
  const createProject = useCreateProject(teamId!)

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateProjectFormData>({
    resolver: zodResolver(createProjectSchema),
    defaultValues: {
      auto_deploy: false,
      github_branch: 'main',
    },
  })

  const autoDeploy = watch('auto_deploy')

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setValue('name', value)
    setValue('slug', value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, ''))
  }

  const onSubmit = async (data: CreateProjectFormData) => {
    try {
      await createProject.mutateAsync(data)
      toast.success('Project created successfully')
      navigate(`/teams/${teamId}`)
    } catch {
      toast.error('Failed to create project')
    }
  }

  if (isLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!team) {
    return null
  }

  return (
    <FormPageLayout
      title="Create Project"
      description={`Create a new project in ${team.name}`}
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: team.name, href: `/teams/${teamId}` },
        { label: 'New Project' },
      ]}
      backHref={`/teams/${teamId}`}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        {/* Basic Info */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Basic Information</h3>

          <div className="space-y-2">
            <Label htmlFor="name">Project Name</Label>
            <Input
              id="name"
              placeholder="My Web App"
              {...register('name')}
              onChange={handleNameChange}
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="slug">Slug</Label>
            <Input
              id="slug"
              placeholder="my-web-app"
              {...register('slug')}
            />
            <p className="text-xs text-muted-foreground">
              URL-friendly identifier. Auto-generated from name.
            </p>
            {errors.slug && (
              <p className="text-sm text-destructive">{errors.slug.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="description">Description (optional)</Label>
            <Textarea
              id="description"
              placeholder="Describe your project..."
              rows={3}
              {...register('description')}
            />
          </div>
        </div>

        {/* GitHub Integration */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">GitHub Integration (optional)</h3>
          <p className="text-sm text-muted-foreground">
            Connect a GitHub repository to enable automatic deployments.
          </p>

          <div className="space-y-2">
            <Label htmlFor="github_repo">Repository</Label>
            <Input
              id="github_repo"
              placeholder="owner/repo"
              {...register('github_repo')}
            />
            <p className="text-xs text-muted-foreground">e.g., github.com/owner/repo</p>
          </div>

          <div className="space-y-2">
            <Label htmlFor="github_branch">Branch</Label>
            <Input
              id="github_branch"
              placeholder="main"
              {...register('github_branch')}
            />
          </div>

          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="space-y-0.5">
              <Label htmlFor="auto_deploy">Auto Deploy</Label>
              <p className="text-xs text-muted-foreground">
                Automatically deploy when code is pushed to the branch
              </p>
            </div>
            <Switch
              id="auto_deploy"
              checked={autoDeploy}
              onCheckedChange={(checked) => setValue('auto_deploy', checked)}
            />
          </div>
        </div>

        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate(`/teams/${teamId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={createProject.isPending}>
            {createProject.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Create Project
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

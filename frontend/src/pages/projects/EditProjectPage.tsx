import { useEffect } from 'react'
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
import { useProject, useUpdateProject } from '@/hooks/api/useProjects'
import { toast } from 'sonner'

const editProjectSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().optional(),
  github_repo: z.string().optional(),
  github_branch: z.string().optional(),
  auto_deploy: z.boolean().optional(),
})

type EditProjectFormData = z.infer<typeof editProjectSchema>

export function EditProjectPage() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const { data: project, isLoading } = useProject(projectId!)
  const updateProject = useUpdateProject(projectId!)

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<EditProjectFormData>({
    resolver: zodResolver(editProjectSchema),
  })

  const autoDeploy = watch('auto_deploy')

  useEffect(() => {
    if (project) {
      reset({
        name: project.name,
        description: project.description || '',
        github_repo: project.github_repo || '',
        github_branch: project.github_branch || 'main',
        auto_deploy: project.auto_deploy,
      })
    }
  }, [project, reset])

  const onSubmit = async (data: EditProjectFormData) => {
    try {
      await updateProject.mutateAsync(data)
      toast.success('Project updated successfully')
      navigate(`/projects/${projectId}`)
    } catch {
      toast.error('Failed to update project')
    }
  }

  if (isLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!project) {
    return null
  }

  return (
    <FormPageLayout
      title="Edit Project"
      description={`Editing ${project.name}`}
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: project.name, href: `/projects/${projectId}` },
        { label: 'Edit' },
      ]}
      backHref={`/projects/${projectId}`}
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
            />
            {errors.name && (
              <p className="text-sm text-destructive">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label>Slug</Label>
            <Input value={project.slug} disabled className="bg-muted" />
            <p className="text-xs text-muted-foreground">
              Slug cannot be changed after creation.
            </p>
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
          <h3 className="text-lg font-medium">GitHub Integration</h3>

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
            onClick={() => navigate(`/projects/${projectId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={updateProject.isPending}>
            {updateProject.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Changes
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

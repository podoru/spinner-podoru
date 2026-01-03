import { useState } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2, Plus, Trash2 } from 'lucide-react'
import { FormPageLayout } from '@/components/common/FormPageLayout'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useProject } from '@/hooks/api/useProjects'
import { useCreateService } from '@/hooks/api/useServices'
import { toast } from 'sonner'

const createServiceSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  slug: z.string().min(2, 'Slug must be at least 2 characters').regex(/^[a-z0-9-]+$/, 'Slug must be lowercase alphanumeric with hyphens'),
  deploy_type: z.enum(['image', 'dockerfile', 'compose']),
  image: z.string().optional(),
  dockerfile_path: z.string().optional(),
  build_context: z.string().optional(),
  compose_file: z.string().optional(),
  replicas: z.number().min(1).optional(),
  cpu_limit: z.number().optional(),
  memory_limit: z.number().optional(),
  health_check_path: z.string().optional(),
  health_check_interval: z.number().optional(),
  restart_policy: z.enum(['no', 'always', 'on-failure', 'unless-stopped']).optional(),
})

type CreateServiceFormData = z.infer<typeof createServiceSchema>

interface EnvVar {
  key: string
  value: string
}

export function CreateServicePage() {
  const { projectId } = useParams<{ projectId: string }>()
  const navigate = useNavigate()
  const { data: project, isLoading } = useProject(projectId!)
  const createService = useCreateService(projectId!)
  const [envVars, setEnvVars] = useState<EnvVar[]>([])

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<CreateServiceFormData>({
    resolver: zodResolver(createServiceSchema),
    defaultValues: {
      deploy_type: 'image',
      replicas: 1,
      restart_policy: 'unless-stopped',
      build_context: '.',
      dockerfile_path: 'Dockerfile',
    },
  })

  const deployType = watch('deploy_type')

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setValue('name', value)
    setValue('slug', value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, ''))
  }

  const addEnvVar = () => {
    setEnvVars([...envVars, { key: '', value: '' }])
  }

  const removeEnvVar = (index: number) => {
    setEnvVars(envVars.filter((_, i) => i !== index))
  }

  const updateEnvVar = (index: number, field: 'key' | 'value', value: string) => {
    const updated = [...envVars]
    updated[index][field] = value
    setEnvVars(updated)
  }

  const onSubmit = async (data: CreateServiceFormData) => {
    try {
      const validEnvVars = envVars.filter(ev => ev.key.trim() !== '')
      await createService.mutateAsync({
        ...data,
        env_vars: validEnvVars.length > 0 ? validEnvVars : undefined,
      })
      toast.success('Service created successfully')
      navigate(`/projects/${projectId}`)
    } catch {
      toast.error('Failed to create service')
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
      title="Create Service"
      description={`Create a new service in ${project.name}`}
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: project.name, href: `/projects/${projectId}` },
        { label: 'New Service' },
      ]}
      backHref={`/projects/${projectId}`}
      maxWidth="2xl"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-8">
        {/* Basic Info */}
        <div className="space-y-4">
          <h3 className="text-lg font-medium">Basic Information</h3>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="name">Service Name</Label>
              <Input
                id="name"
                placeholder="api-server"
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
                placeholder="api-server"
                {...register('slug')}
              />
              {errors.slug && (
                <p className="text-sm text-destructive">{errors.slug.message}</p>
              )}
            </div>
          </div>
        </div>

        {/* Deploy Configuration */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">Deployment Configuration</h3>

          <div className="space-y-2">
            <Label>Deploy Type</Label>
            <Select
              value={deployType}
              onValueChange={(value: 'image' | 'dockerfile' | 'compose') => setValue('deploy_type', value)}
            >
              <SelectTrigger>
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="image">Docker Image</SelectItem>
                <SelectItem value="dockerfile">Dockerfile</SelectItem>
                <SelectItem value="compose">Docker Compose</SelectItem>
              </SelectContent>
            </Select>
          </div>

          {deployType === 'image' && (
            <div className="space-y-2">
              <Label htmlFor="image">Image</Label>
              <Input
                id="image"
                placeholder="nginx:latest"
                {...register('image')}
              />
              <p className="text-xs text-muted-foreground">Docker image to deploy (e.g., nginx:latest)</p>
            </div>
          )}

          {deployType === 'dockerfile' && (
            <div className="grid gap-4 sm:grid-cols-2">
              <div className="space-y-2">
                <Label htmlFor="dockerfile_path">Dockerfile Path</Label>
                <Input
                  id="dockerfile_path"
                  placeholder="Dockerfile"
                  {...register('dockerfile_path')}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="build_context">Build Context</Label>
                <Input
                  id="build_context"
                  placeholder="."
                  {...register('build_context')}
                />
              </div>
            </div>
          )}

          {deployType === 'compose' && (
            <div className="space-y-2">
              <Label htmlFor="compose_file">Compose File</Label>
              <Input
                id="compose_file"
                placeholder="docker-compose.yml"
                {...register('compose_file')}
              />
            </div>
          )}
        </div>

        {/* Resources */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">Resources</h3>

          <div className="grid gap-4 sm:grid-cols-3">
            <div className="space-y-2">
              <Label htmlFor="replicas">Replicas</Label>
              <Input
                id="replicas"
                type="number"
                min={1}
                {...register('replicas', { valueAsNumber: true })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="cpu_limit">CPU Limit (cores)</Label>
              <Input
                id="cpu_limit"
                type="number"
                step="0.1"
                placeholder="0.5"
                {...register('cpu_limit', { valueAsNumber: true })}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="memory_limit">Memory Limit (MB)</Label>
              <Input
                id="memory_limit"
                type="number"
                placeholder="512"
                {...register('memory_limit', { valueAsNumber: true })}
              />
            </div>
          </div>
        </div>

        {/* Health Check */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">Health Check (optional)</h3>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="health_check_path">Health Check Path</Label>
              <Input
                id="health_check_path"
                placeholder="/health"
                {...register('health_check_path')}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="health_check_interval">Interval (seconds)</Label>
              <Input
                id="health_check_interval"
                type="number"
                placeholder="30"
                {...register('health_check_interval', { valueAsNumber: true })}
              />
            </div>
          </div>
        </div>

        {/* Restart Policy */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">Restart Policy</h3>

          <div className="space-y-2">
            <Select
              value={watch('restart_policy')}
              onValueChange={(value: 'no' | 'always' | 'on-failure' | 'unless-stopped') => setValue('restart_policy', value)}
            >
              <SelectTrigger className="w-full sm:w-64">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="no">No</SelectItem>
                <SelectItem value="always">Always</SelectItem>
                <SelectItem value="on-failure">On Failure</SelectItem>
                <SelectItem value="unless-stopped">Unless Stopped</SelectItem>
              </SelectContent>
            </Select>
          </div>
        </div>

        {/* Environment Variables */}
        <div className="space-y-4 pt-4 border-t">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Environment Variables</h3>
            <Button type="button" variant="outline" size="sm" onClick={addEnvVar}>
              <Plus className="mr-2 h-4 w-4" />
              Add Variable
            </Button>
          </div>

          {envVars.length === 0 ? (
            <p className="text-sm text-muted-foreground">No environment variables defined.</p>
          ) : (
            <div className="space-y-2">
              {envVars.map((envVar, index) => (
                <div key={index} className="flex gap-2">
                  <Input
                    placeholder="KEY"
                    value={envVar.key}
                    onChange={(e) => updateEnvVar(index, 'key', e.target.value)}
                    className="font-mono"
                  />
                  <Input
                    placeholder="value"
                    value={envVar.value}
                    onChange={(e) => updateEnvVar(index, 'value', e.target.value)}
                    className="font-mono"
                  />
                  <Button
                    type="button"
                    variant="ghost"
                    size="icon"
                    onClick={() => removeEnvVar(index)}
                  >
                    <Trash2 className="h-4 w-4" />
                  </Button>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate(`/projects/${projectId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={createService.isPending}>
            {createService.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Create Service
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

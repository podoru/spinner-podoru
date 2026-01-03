import { useEffect, useState } from 'react'
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
import { useService, useUpdateService } from '@/hooks/api/useServices'
import { toast } from 'sonner'

const editServiceSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
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

type EditServiceFormData = z.infer<typeof editServiceSchema>

interface EnvVar {
  key: string
  value: string
}

export function EditServicePage() {
  const { serviceId } = useParams<{ serviceId: string }>()
  const navigate = useNavigate()
  const { data: service, isLoading } = useService(serviceId!)
  const updateService = useUpdateService(serviceId!)
  const [envVars, setEnvVars] = useState<EnvVar[]>([])

  const {
    register,
    handleSubmit,
    reset,
    setValue,
    watch,
    formState: { errors },
  } = useForm<EditServiceFormData>({
    resolver: zodResolver(editServiceSchema),
  })

  useEffect(() => {
    if (service) {
      reset({
        name: service.name,
        image: service.image || '',
        dockerfile_path: service.dockerfile_path || '',
        build_context: service.build_context || '.',
        compose_file: service.compose_file || '',
        replicas: service.replicas,
        cpu_limit: service.cpu_limit || undefined,
        memory_limit: service.memory_limit || undefined,
        health_check_path: service.health_check_path || '',
        health_check_interval: service.health_check_interval || undefined,
        restart_policy: service.restart_policy as 'no' | 'always' | 'on-failure' | 'unless-stopped',
      })
      setEnvVars(service.env_vars || [])
    }
  }, [service, reset])

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

  const onSubmit = async (data: EditServiceFormData) => {
    try {
      const validEnvVars = envVars.filter(ev => ev.key.trim() !== '')
      await updateService.mutateAsync({
        ...data,
        env_vars: validEnvVars,
      })
      toast.success('Service updated successfully')
      navigate(`/services/${serviceId}`)
    } catch {
      toast.error('Failed to update service')
    }
  }

  if (isLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!service) {
    return null
  }

  return (
    <FormPageLayout
      title="Edit Service"
      description={`Editing ${service.name}`}
      breadcrumbs={[
        { label: 'Projects', href: `/projects/${service.project_id}` },
        { label: service.name, href: `/services/${serviceId}` },
        { label: 'Edit' },
      ]}
      backHref={`/services/${serviceId}`}
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
              />
              {errors.name && (
                <p className="text-sm text-destructive">{errors.name.message}</p>
              )}
            </div>

            <div className="space-y-2">
              <Label>Slug</Label>
              <Input value={service.slug} disabled className="bg-muted" />
              <p className="text-xs text-muted-foreground">Slug cannot be changed.</p>
            </div>
          </div>

          <div className="space-y-2">
            <Label>Deploy Type</Label>
            <Input value={service.deploy_type} disabled className="bg-muted w-full sm:w-48" />
            <p className="text-xs text-muted-foreground">Deploy type cannot be changed after creation.</p>
          </div>
        </div>

        {/* Deploy Configuration */}
        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">Deployment Configuration</h3>

          {service.deploy_type === 'image' && (
            <div className="space-y-2">
              <Label htmlFor="image">Image</Label>
              <Input
                id="image"
                placeholder="nginx:latest"
                {...register('image')}
              />
            </div>
          )}

          {service.deploy_type === 'dockerfile' && (
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

          {service.deploy_type === 'compose' && (
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
          <h3 className="text-lg font-medium">Health Check</h3>

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
            onClick={() => navigate(`/services/${serviceId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={updateService.isPending}>
            {updateService.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Changes
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

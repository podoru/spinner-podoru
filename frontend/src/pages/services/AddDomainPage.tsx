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
import { Switch } from '@/components/ui/switch'
import { useService } from '@/hooks/api/useServices'
import { useAddDomain } from '@/hooks/api/useDomains'
import { toast } from 'sonner'

const addDomainSchema = z.object({
  domain: z.string().min(3, 'Domain is required'),
  ssl_enabled: z.boolean().optional(),
  ssl_auto: z.boolean().optional(),
})

type AddDomainFormData = z.infer<typeof addDomainSchema>

export function AddDomainPage() {
  const { serviceId } = useParams<{ serviceId: string }>()
  const navigate = useNavigate()
  const { data: service, isLoading } = useService(serviceId!)
  const addDomain = useAddDomain(serviceId!)

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<AddDomainFormData>({
    resolver: zodResolver(addDomainSchema),
    defaultValues: {
      ssl_enabled: true,
      ssl_auto: true,
    },
  })

  const sslEnabled = watch('ssl_enabled')
  const sslAuto = watch('ssl_auto')

  const onSubmit = async (data: AddDomainFormData) => {
    try {
      await addDomain.mutateAsync(data)
      toast.success('Domain added successfully')
      navigate(`/services/${serviceId}`)
    } catch {
      toast.error('Failed to add domain')
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
      title="Add Domain"
      description={`Add a custom domain to ${service.name}`}
      breadcrumbs={[
        { label: 'Services', href: `/services/${serviceId}` },
        { label: service.name, href: `/services/${serviceId}` },
        { label: 'Add Domain' },
      ]}
      backHref={`/services/${serviceId}`}
      maxWidth="lg"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="domain">Domain Name</Label>
          <Input
            id="domain"
            placeholder="api.example.com"
            {...register('domain')}
          />
          <p className="text-xs text-muted-foreground">
            Enter the domain name that should point to this service. Make sure to update your DNS records.
          </p>
          {errors.domain && (
            <p className="text-sm text-destructive">{errors.domain.message}</p>
          )}
        </div>

        <div className="space-y-4 pt-4 border-t">
          <h3 className="text-lg font-medium">SSL Configuration</h3>

          <div className="flex items-center justify-between rounded-lg border p-4">
            <div className="space-y-0.5">
              <Label htmlFor="ssl_enabled">Enable SSL</Label>
              <p className="text-xs text-muted-foreground">
                Secure the domain with HTTPS
              </p>
            </div>
            <Switch
              id="ssl_enabled"
              checked={sslEnabled}
              onCheckedChange={(checked) => setValue('ssl_enabled', checked)}
            />
          </div>

          {sslEnabled && (
            <div className="flex items-center justify-between rounded-lg border p-4">
              <div className="space-y-0.5">
                <Label htmlFor="ssl_auto">Auto SSL (Let's Encrypt)</Label>
                <p className="text-xs text-muted-foreground">
                  Automatically provision and renew SSL certificates
                </p>
              </div>
              <Switch
                id="ssl_auto"
                checked={sslAuto}
                onCheckedChange={(checked) => setValue('ssl_auto', checked)}
              />
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
          <Button type="submit" disabled={addDomain.isPending}>
            {addDomain.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Add Domain
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

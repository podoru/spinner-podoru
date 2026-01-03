import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Loader2 } from 'lucide-react'
import { FormPageLayout } from '@/components/common/FormPageLayout'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Textarea } from '@/components/ui/textarea'
import { useCreateTeam } from '@/hooks/api/useTeams'
import { toast } from 'sonner'

const createTeamSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  slug: z.string().min(2, 'Slug must be at least 2 characters').regex(/^[a-z0-9-]+$/, 'Slug must be lowercase alphanumeric with hyphens'),
  description: z.string().optional(),
})

type CreateTeamFormData = z.infer<typeof createTeamSchema>

export function CreateTeamPage() {
  const navigate = useNavigate()
  const createTeam = useCreateTeam()

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm<CreateTeamFormData>({
    resolver: zodResolver(createTeamSchema),
  })

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value
    setValue('name', value)
    setValue('slug', value.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, ''))
  }

  const onSubmit = async (data: CreateTeamFormData) => {
    try {
      await createTeam.mutateAsync(data)
      toast.success('Team created successfully')
      navigate('/teams')
    } catch {
      toast.error('Failed to create team')
    }
  }

  return (
    <FormPageLayout
      title="Create Team"
      description="Create a new team to organize your projects and services."
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: 'New Team' },
      ]}
      backHref="/teams"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="name">Team Name</Label>
          <Input
            id="name"
            placeholder="My Awesome Team"
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
            placeholder="my-awesome-team"
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
            placeholder="Describe your team's purpose..."
            rows={3}
            {...register('description')}
          />
        </div>

        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate('/teams')}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={createTeam.isPending}>
            {createTeam.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Create Team
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

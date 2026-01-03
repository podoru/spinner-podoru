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
import { useTeam, useUpdateTeam } from '@/hooks/api/useTeams'
import { toast } from 'sonner'

const editTeamSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().optional(),
})

type EditTeamFormData = z.infer<typeof editTeamSchema>

export function EditTeamPage() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const { data: team, isLoading } = useTeam(teamId!)
  const updateTeam = useUpdateTeam(teamId!)

  const {
    register,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm<EditTeamFormData>({
    resolver: zodResolver(editTeamSchema),
  })

  useEffect(() => {
    if (team) {
      reset({
        name: team.name,
        description: team.description || '',
      })
    }
  }, [team, reset])

  const onSubmit = async (data: EditTeamFormData) => {
    try {
      await updateTeam.mutateAsync(data)
      toast.success('Team updated successfully')
      navigate(`/teams/${teamId}`)
    } catch {
      toast.error('Failed to update team')
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
      title="Edit Team"
      description={`Editing ${team.name}`}
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: team.name, href: `/teams/${teamId}` },
        { label: 'Edit' },
      ]}
      backHref={`/teams/${teamId}`}
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="name">Team Name</Label>
          <Input
            id="name"
            placeholder="My Awesome Team"
            {...register('name')}
          />
          {errors.name && (
            <p className="text-sm text-destructive">{errors.name.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label>Slug</Label>
          <Input value={team.slug} disabled className="bg-muted" />
          <p className="text-xs text-muted-foreground">
            Slug cannot be changed after creation.
          </p>
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
            onClick={() => navigate(`/teams/${teamId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={updateTeam.isPending}>
            {updateTeam.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Save Changes
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

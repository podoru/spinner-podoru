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
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { useTeam } from '@/hooks/api/useTeams'
import { useAddTeamMember } from '@/hooks/api/useTeamMembers'
import { toast } from 'sonner'

const addMemberSchema = z.object({
  email: z.string().email('Valid email is required'),
  role: z.enum(['admin', 'member'], { required_error: 'Please select a role' }),
})

type AddMemberFormData = z.infer<typeof addMemberSchema>

export function AddMemberPage() {
  const { teamId } = useParams<{ teamId: string }>()
  const navigate = useNavigate()
  const { data: team, isLoading } = useTeam(teamId!)
  const addMember = useAddTeamMember(teamId!)

  const {
    register,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm<AddMemberFormData>({
    resolver: zodResolver(addMemberSchema),
    defaultValues: {
      role: 'member',
    },
  })

  const selectedRole = watch('role')

  const onSubmit = async (data: AddMemberFormData) => {
    try {
      await addMember.mutateAsync(data)
      toast.success('Member added successfully')
      navigate(`/teams/${teamId}`)
    } catch {
      toast.error('Failed to add member. User may not exist or is already a member.')
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
      title="Add Team Member"
      description={`Add a new member to ${team.name}`}
      breadcrumbs={[
        { label: 'Teams', href: '/teams' },
        { label: team.name, href: `/teams/${teamId}` },
        { label: 'Add Member' },
      ]}
      backHref={`/teams/${teamId}`}
      maxWidth="lg"
    >
      <form onSubmit={handleSubmit(onSubmit)} className="space-y-6">
        <div className="space-y-2">
          <Label htmlFor="email">Email Address</Label>
          <Input
            id="email"
            type="email"
            placeholder="member@example.com"
            {...register('email')}
          />
          <p className="text-xs text-muted-foreground">
            Enter the email address of the user you want to add. They must have an account.
          </p>
          {errors.email && (
            <p className="text-sm text-destructive">{errors.email.message}</p>
          )}
        </div>

        <div className="space-y-2">
          <Label>Role</Label>
          <Select
            value={selectedRole}
            onValueChange={(value: 'admin' | 'member') => setValue('role', value)}
          >
            <SelectTrigger>
              <SelectValue placeholder="Select a role" />
            </SelectTrigger>
            <SelectContent>
              <SelectItem value="member">Member</SelectItem>
              <SelectItem value="admin">Admin</SelectItem>
            </SelectContent>
          </Select>
          <p className="text-xs text-muted-foreground">
            Admins can manage team settings and members. Members can manage projects and services.
          </p>
          {errors.role && (
            <p className="text-sm text-destructive">{errors.role.message}</p>
          )}
        </div>

        <div className="flex gap-3 pt-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => navigate(`/teams/${teamId}`)}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={addMember.isPending}>
            {addMember.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
            Add Member
          </Button>
        </div>
      </form>
    </FormPageLayout>
  )
}

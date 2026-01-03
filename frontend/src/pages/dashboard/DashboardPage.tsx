import { Link } from 'react-router-dom'
import { Users, Plus } from 'lucide-react'
import { PageHeader } from '@/components/common/PageHeader'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Skeleton } from '@/components/ui/skeleton'
import { useTeams } from '@/hooks/api/useTeams'
import { useAuthStore } from '@/stores/authStore'

export function DashboardPage() {
  const { user } = useAuthStore()
  const { data: teams, isLoading } = useTeams()

  const stats = [
    {
      title: 'Teams',
      value: teams?.length ?? 0,
      icon: Users,
      href: '/teams',
    },
  ]

  return (
    <div className="space-y-6">
      <PageHeader
        title={`Welcome back, ${user?.name?.split(' ')[0] ?? 'User'}`}
        description="Here's an overview of your deployment platform"
      />

      {/* Stats */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        {isLoading ? (
          Array.from({ length: 4 }).map((_, i) => (
            <Card key={i}>
              <CardHeader className="flex flex-row items-center justify-between pb-2">
                <Skeleton className="h-4 w-20" />
                <Skeleton className="h-4 w-4" />
              </CardHeader>
              <CardContent>
                <Skeleton className="h-8 w-16" />
              </CardContent>
            </Card>
          ))
        ) : (
          stats.map((stat) => (
            <Link key={stat.title} to={stat.href}>
              <Card className="hover:bg-muted/50 transition-colors">
                <CardHeader className="flex flex-row items-center justify-between pb-2">
                  <CardTitle className="text-sm font-medium text-muted-foreground">
                    {stat.title}
                  </CardTitle>
                  <stat.icon className="h-4 w-4 text-muted-foreground" />
                </CardHeader>
                <CardContent>
                  <div className="text-2xl font-bold">{stat.value}</div>
                </CardContent>
              </Card>
            </Link>
          ))
        )}
      </div>

      {/* Quick Actions */}
      <Card>
        <CardHeader>
          <CardTitle>Quick Actions</CardTitle>
          <CardDescription>Get started with common tasks</CardDescription>
        </CardHeader>
        <CardContent className="flex flex-wrap gap-2">
          <Button asChild>
            <Link to="/teams">
              <Users className="mr-2 h-4 w-4" />
              View Teams
            </Link>
          </Button>
        </CardContent>
      </Card>

      {/* Recent Teams */}
      <Card>
        <CardHeader>
          <CardTitle>Your Teams</CardTitle>
          <CardDescription>Teams you're a member of</CardDescription>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="space-y-2">
              {Array.from({ length: 3 }).map((_, i) => (
                <Skeleton key={i} className="h-12 w-full" />
              ))}
            </div>
          ) : teams?.length === 0 ? (
            <div className="text-center py-6 text-muted-foreground">
              <p>You're not a member of any teams yet.</p>
              <Button asChild className="mt-4">
                <Link to="/teams">
                  <Plus className="mr-2 h-4 w-4" />
                  Create a Team
                </Link>
              </Button>
            </div>
          ) : (
            <div className="space-y-2">
              {teams?.slice(0, 5).map((team) => (
                <Link
                  key={team.id}
                  to={`/teams/${team.id}`}
                  className="flex items-center justify-between rounded-lg border p-3 hover:bg-muted/50 transition-colors"
                >
                  <div>
                    <p className="font-medium">{team.name}</p>
                    <p className="text-sm text-muted-foreground">{team.slug}</p>
                  </div>
                  <span className="text-xs text-muted-foreground capitalize">{team.role}</span>
                </Link>
              ))}
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  )
}

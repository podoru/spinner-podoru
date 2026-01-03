import { Link } from 'react-router-dom'
import { Users, ChevronRight } from 'lucide-react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import type { TeamWithRole } from '@/types'

interface TeamCardProps {
  team: TeamWithRole
}

export function TeamCard({ team }: TeamCardProps) {
  return (
    <Link to={`/teams/${team.id}`}>
      <Card className="hover:bg-muted/50 transition-colors">
        <CardHeader className="pb-3">
          <div className="flex items-start justify-between">
            <div>
              <CardTitle className="text-lg">{team.name}</CardTitle>
              <CardDescription className="text-sm">{team.slug}</CardDescription>
            </div>
            <Badge variant="outline" className="capitalize">
              {team.role}
            </Badge>
          </div>
        </CardHeader>
        <CardContent>
          {team.description && (
            <p className="text-sm text-muted-foreground line-clamp-2 mb-4">
              {team.description}
            </p>
          )}
          <div className="flex items-center justify-between text-sm text-muted-foreground">
            <div className="flex items-center gap-4">
              <span className="flex items-center gap-1">
                <Users className="h-4 w-4" />
                Members
              </span>
            </div>
            <ChevronRight className="h-4 w-4" />
          </div>
        </CardContent>
      </Card>
    </Link>
  )
}

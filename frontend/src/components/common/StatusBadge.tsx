import { cn } from '@/lib/utils'
import { Badge } from '@/components/ui/badge'
import type { ServiceStatus } from '@/types'

interface StatusBadgeProps {
  status: ServiceStatus
  className?: string
}

const statusConfig: Record<ServiceStatus, { label: string; className: string }> = {
  running: {
    label: 'Running',
    className: 'bg-green-500/15 text-green-700 border-green-500/20',
  },
  stopped: {
    label: 'Stopped',
    className: 'bg-gray-500/15 text-gray-700 border-gray-500/20',
  },
  deploying: {
    label: 'Deploying',
    className: 'bg-blue-500/15 text-blue-700 border-blue-500/20 animate-pulse',
  },
  failed: {
    label: 'Failed',
    className: 'bg-red-500/15 text-red-700 border-red-500/20',
  },
}

export function StatusBadge({ status, className }: StatusBadgeProps) {
  const config = statusConfig[status]

  return (
    <Badge variant="outline" className={cn(config.className, className)}>
      {status === 'running' && (
        <span className="mr-1.5 h-2 w-2 rounded-full bg-green-500 animate-pulse" />
      )}
      {config.label}
    </Badge>
  )
}

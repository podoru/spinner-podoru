import { useState } from 'react'
import { useParams, Link, useNavigate } from 'react-router-dom'
import {
  ArrowLeft, Play, Square, RotateCw, Rocket, Trash2,
  Globe, Terminal, Settings, Pencil, Plus
} from 'lucide-react'
import { Breadcrumb } from '@/components/common/Breadcrumb'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'
import { EmptyState } from '@/components/common/EmptyState'
import { StatusBadge } from '@/components/common/StatusBadge'
import { ConfirmDialog } from '@/components/common/ConfirmDialog'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  useService, useDeleteService,
  useDeployService, useStartService, useStopService, useRestartService,
  useServiceLogs
} from '@/hooks/api/useServices'
import { useProject } from '@/hooks/api/useProjects'
import { useServiceDomains, useDeleteDomain } from '@/hooks/api/useDomains'
import { toast } from 'sonner'

export function ServiceDetailPage() {
  const { serviceId } = useParams<{ serviceId: string }>()
  const navigate = useNavigate()
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
  const [logTail, setLogTail] = useState(100)

  const { data: service, isLoading: serviceLoading } = useService(serviceId!)
  const { data: project } = useProject(service?.project_id || '')
  const { data: domains, isLoading: domainsLoading } = useServiceDomains(serviceId!)
  const { data: logsData, isLoading: logsLoading } = useServiceLogs(serviceId!, { tail: logTail })

  const deleteService = useDeleteService()
  const deployService = useDeployService(serviceId!)
  const startService = useStartService(serviceId!)
  const stopService = useStopService(serviceId!)
  const restartService = useRestartService(serviceId!)
  const deleteDomain = useDeleteDomain(serviceId!)

  const handleDelete = async () => {
    try {
      await deleteService.mutateAsync(serviceId!)
      toast.success('Service deleted successfully')
      navigate(-1)
    } catch {
      toast.error('Failed to delete service')
    }
  }

  const handleDeploy = async () => {
    try {
      await deployService.mutateAsync()
      toast.success('Deployment started')
    } catch {
      toast.error('Failed to start deployment')
    }
  }

  const handleStart = async () => {
    try {
      await startService.mutateAsync()
      toast.success('Service started')
    } catch {
      toast.error('Failed to start service')
    }
  }

  const handleStop = async () => {
    try {
      await stopService.mutateAsync()
      toast.success('Service stopped')
    } catch {
      toast.error('Failed to stop service')
    }
  }

  const handleRestart = async () => {
    try {
      await restartService.mutateAsync()
      toast.success('Service restarted')
    } catch {
      toast.error('Failed to restart service')
    }
  }

  if (serviceLoading) {
    return <LoadingSpinner className="h-64" />
  }

  if (!service) {
    return (
      <EmptyState
        title="Service not found"
        description="The service you're looking for doesn't exist."
      >
        <Button asChild>
          <Link to="/teams">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Teams
          </Link>
        </Button>
      </EmptyState>
    )
  }

  const isRunning = service.status === 'running'
  const isDeploying = service.status === 'deploying'

  return (
    <div className="space-y-6">
      <Breadcrumb
        items={[
          { label: 'Projects', href: `/projects/${service.project_id}` },
          { label: project?.name || 'Project', href: `/projects/${service.project_id}` },
          { label: service.name },
        ]}
      />

      <div className="flex items-center justify-between">
        <div className="flex items-center gap-4">
          <Button variant="ghost" size="icon" asChild>
            <Link to={`/projects/${service.project_id}`}>
              <ArrowLeft className="h-4 w-4" />
            </Link>
          </Button>
          <div>
            <div className="flex items-center gap-3">
              <h1 className="text-2xl font-bold tracking-tight">{service.name}</h1>
              <StatusBadge status={service.status} />
            </div>
            <p className="text-muted-foreground">{service.slug}</p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm" asChild>
            <Link to={`/services/${serviceId}/edit`}>
              <Pencil className="mr-2 h-4 w-4" />
              Edit
            </Link>
          </Button>
          <Button variant="destructive" size="sm" onClick={() => setDeleteDialogOpen(true)}>
            <Trash2 className="mr-2 h-4 w-4" />
            Delete
          </Button>
        </div>
      </div>

      {/* Service Controls */}
      <Card>
        <CardHeader>
          <CardTitle>Controls</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="flex flex-wrap gap-2">
            <Button onClick={handleDeploy} disabled={isDeploying || deployService.isPending}>
              <Rocket className="mr-2 h-4 w-4" />
              Deploy
            </Button>
            <Button
              variant="outline"
              onClick={handleStart}
              disabled={isRunning || isDeploying || startService.isPending}
            >
              <Play className="mr-2 h-4 w-4" />
              Start
            </Button>
            <Button
              variant="outline"
              onClick={handleStop}
              disabled={!isRunning || stopService.isPending}
            >
              <Square className="mr-2 h-4 w-4" />
              Stop
            </Button>
            <Button
              variant="outline"
              onClick={handleRestart}
              disabled={!isRunning || restartService.isPending}
            >
              <RotateCw className="mr-2 h-4 w-4" />
              Restart
            </Button>
          </div>
        </CardContent>
      </Card>

      <Tabs defaultValue="overview">
        <TabsList>
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="logs">
            <Terminal className="mr-2 h-4 w-4" />
            Logs
          </TabsTrigger>
          <TabsTrigger value="domains">
            <Globe className="mr-2 h-4 w-4" />
            Domains
          </TabsTrigger>
          <TabsTrigger value="settings">
            <Settings className="mr-2 h-4 w-4" />
            Settings
          </TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="mt-6">
          <div className="grid gap-4 md:grid-cols-2">
            <Card>
              <CardHeader>
                <CardTitle>Deployment</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Type</span>
                  <Badge variant="outline">{service.deploy_type}</Badge>
                </div>
                {service.image && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Image</span>
                    <span className="font-mono text-sm">{service.image}</span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Replicas</span>
                  <span>{service.replicas}</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-muted-foreground">Restart Policy</span>
                  <span>{service.restart_policy}</span>
                </div>
              </CardContent>
            </Card>

            <Card>
              <CardHeader>
                <CardTitle>Resources</CardTitle>
              </CardHeader>
              <CardContent className="space-y-2">
                {service.cpu_limit && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">CPU Limit</span>
                    <span>{service.cpu_limit} cores</span>
                  </div>
                )}
                {service.memory_limit && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Memory Limit</span>
                    <span>{service.memory_limit} MB</span>
                  </div>
                )}
                {service.health_check_path && (
                  <div className="flex justify-between">
                    <span className="text-muted-foreground">Health Check</span>
                    <span className="font-mono text-sm">{service.health_check_path}</span>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="logs" className="mt-6">
          <Card>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle>Container Logs</CardTitle>
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">Tail:</span>
                  {[100, 500, 1000].map((n) => (
                    <Button
                      key={n}
                      variant={logTail === n ? 'default' : 'outline'}
                      size="sm"
                      onClick={() => setLogTail(n)}
                    >
                      {n}
                    </Button>
                  ))}
                </div>
              </div>
            </CardHeader>
            <CardContent>
              {logsLoading ? (
                <LoadingSpinner className="h-32" />
              ) : (
                <ScrollArea className="h-96 rounded-md border bg-zinc-950 p-4">
                  <pre className="text-sm text-zinc-100 font-mono whitespace-pre-wrap">
                    {logsData?.logs || 'No logs available'}
                  </pre>
                </ScrollArea>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="domains" className="mt-6">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <div>
                <CardTitle>Domains</CardTitle>
                <CardDescription>Manage custom domains for this service</CardDescription>
              </div>
              <Button asChild size="sm">
                <Link to={`/services/${serviceId}/domains/add`}>
                  <Plus className="mr-2 h-4 w-4" />
                  Add Domain
                </Link>
              </Button>
            </CardHeader>
            <CardContent>
              {domainsLoading ? (
                <LoadingSpinner className="h-32" />
              ) : domains?.length === 0 ? (
                <EmptyState
                  icon={Globe}
                  title="No domains configured"
                  description="Add a custom domain to access this service."
                >
                  <Button asChild>
                    <Link to={`/services/${serviceId}/domains/add`}>
                      <Plus className="mr-2 h-4 w-4" />
                      Add Domain
                    </Link>
                  </Button>
                </EmptyState>
              ) : (
                <Table>
                  <TableHeader>
                    <TableRow>
                      <TableHead>Domain</TableHead>
                      <TableHead>SSL</TableHead>
                      <TableHead>Created</TableHead>
                      <TableHead></TableHead>
                    </TableRow>
                  </TableHeader>
                  <TableBody>
                    {domains?.map((domain) => (
                      <TableRow key={domain.id}>
                        <TableCell className="font-mono">{domain.domain}</TableCell>
                        <TableCell>
                          {domain.ssl_enabled ? (
                            <Badge variant="outline" className="bg-green-500/15 text-green-700">
                              {domain.ssl_auto ? 'Auto SSL' : 'SSL'}
                            </Badge>
                          ) : (
                            <Badge variant="outline">No SSL</Badge>
                          )}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {new Date(domain.created_at).toLocaleDateString()}
                        </TableCell>
                        <TableCell>
                          <Button
                            variant="ghost"
                            size="sm"
                            onClick={() => deleteDomain.mutate(domain.id)}
                          >
                            <Trash2 className="h-4 w-4" />
                          </Button>
                        </TableCell>
                      </TableRow>
                    ))}
                  </TableBody>
                </Table>
              )}
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="settings" className="mt-6">
          <Card>
            <CardHeader>
              <CardTitle>Service Settings</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <h4 className="font-medium">Name</h4>
                <p className="text-sm text-muted-foreground">{service.name}</p>
              </div>
              <div>
                <h4 className="font-medium">Slug</h4>
                <p className="text-sm text-muted-foreground">{service.slug}</p>
              </div>
              <div>
                <h4 className="font-medium">Created</h4>
                <p className="text-sm text-muted-foreground">
                  {new Date(service.created_at).toLocaleString()}
                </p>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>

      <ConfirmDialog
        open={deleteDialogOpen}
        onOpenChange={setDeleteDialogOpen}
        title="Delete Service"
        description="Are you sure you want to delete this service? This will stop and remove the container."
        confirmText="Delete"
        variant="destructive"
        isLoading={deleteService.isPending}
        onConfirm={handleDelete}
      />
    </div>
  )
}

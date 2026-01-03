import { ArrowLeft } from 'lucide-react'
import { Link } from 'react-router-dom'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Breadcrumb } from './Breadcrumb'
import type { BreadcrumbItem } from './Breadcrumb'

interface FormPageLayoutProps {
  title: string
  description?: string
  breadcrumbs: BreadcrumbItem[]
  backHref: string
  children: React.ReactNode
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl' | '2xl'
}

export function FormPageLayout({
  title,
  description,
  breadcrumbs,
  backHref,
  children,
  maxWidth = '2xl'
}: FormPageLayoutProps) {
  const maxWidthClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl',
    '2xl': 'max-w-2xl',
  }

  return (
    <div className="space-y-6">
      <Breadcrumb items={breadcrumbs} />

      <div className="flex items-center gap-4">
        <Button variant="ghost" size="icon" asChild>
          <Link to={backHref}>
            <ArrowLeft className="h-4 w-4" />
          </Link>
        </Button>
        <div>
          <h1 className="text-2xl font-bold tracking-tight">{title}</h1>
          {description && (
            <p className="text-muted-foreground">{description}</p>
          )}
        </div>
      </div>

      <Card className={maxWidthClasses[maxWidth]}>
        <CardContent className="pt-6">
          {children}
        </CardContent>
      </Card>
    </div>
  )
}

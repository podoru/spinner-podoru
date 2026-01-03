import { createBrowserRouter, Navigate } from 'react-router-dom'
import { AppLayout } from '@/components/layout/AppLayout'
import { AuthLayout } from '@/components/layout/AuthLayout'
import { ProtectedRoute } from '@/components/auth/ProtectedRoute'

// Auth pages
import { LoginPage } from '@/pages/auth/LoginPage'
import { RegisterPage } from '@/pages/auth/RegisterPage'

// Protected pages
import { DashboardPage } from '@/pages/dashboard/DashboardPage'
import { NotFoundPage } from '@/pages/NotFoundPage'

// Lazy load heavier pages
import { lazy, Suspense } from 'react'
import { LoadingSpinner } from '@/components/common/LoadingSpinner'

// Teams
const TeamsPage = lazy(() => import('@/pages/teams/TeamsPage').then(m => ({ default: m.TeamsPage })))
const TeamDetailPage = lazy(() => import('@/pages/teams/TeamDetailPage').then(m => ({ default: m.TeamDetailPage })))
const CreateTeamPage = lazy(() => import('@/pages/teams/CreateTeamPage').then(m => ({ default: m.CreateTeamPage })))
const EditTeamPage = lazy(() => import('@/pages/teams/EditTeamPage').then(m => ({ default: m.EditTeamPage })))
const AddMemberPage = lazy(() => import('@/pages/teams/AddMemberPage').then(m => ({ default: m.AddMemberPage })))

// Projects
const ProjectDetailPage = lazy(() => import('@/pages/projects/ProjectDetailPage').then(m => ({ default: m.ProjectDetailPage })))
const CreateProjectPage = lazy(() => import('@/pages/projects/CreateProjectPage').then(m => ({ default: m.CreateProjectPage })))
const EditProjectPage = lazy(() => import('@/pages/projects/EditProjectPage').then(m => ({ default: m.EditProjectPage })))

// Services
const ServiceDetailPage = lazy(() => import('@/pages/services/ServiceDetailPage').then(m => ({ default: m.ServiceDetailPage })))
const CreateServicePage = lazy(() => import('@/pages/services/CreateServicePage').then(m => ({ default: m.CreateServicePage })))
const EditServicePage = lazy(() => import('@/pages/services/EditServicePage').then(m => ({ default: m.EditServicePage })))
const AddDomainPage = lazy(() => import('@/pages/services/AddDomainPage').then(m => ({ default: m.AddDomainPage })))

// Profile
const ProfilePage = lazy(() => import('@/pages/profile/ProfilePage').then(m => ({ default: m.ProfilePage })))

function LazyWrapper({ children }: { children: React.ReactNode }) {
  return (
    <Suspense fallback={<LoadingSpinner className="h-64" />}>
      {children}
    </Suspense>
  )
}

export const router = createBrowserRouter([
  {
    path: '/',
    element: <Navigate to="/dashboard" replace />,
  },
  {
    element: <AuthLayout />,
    children: [
      {
        path: '/login',
        element: <LoginPage />,
      },
      {
        path: '/register',
        element: <RegisterPage />,
      },
    ],
  },
  {
    element: (
      <ProtectedRoute>
        <AppLayout />
      </ProtectedRoute>
    ),
    children: [
      {
        path: '/dashboard',
        element: <DashboardPage />,
      },
      // Teams
      {
        path: '/teams',
        element: <LazyWrapper><TeamsPage /></LazyWrapper>,
      },
      {
        path: '/teams/new',
        element: <LazyWrapper><CreateTeamPage /></LazyWrapper>,
      },
      {
        path: '/teams/:teamId',
        element: <LazyWrapper><TeamDetailPage /></LazyWrapper>,
      },
      {
        path: '/teams/:teamId/edit',
        element: <LazyWrapper><EditTeamPage /></LazyWrapper>,
      },
      {
        path: '/teams/:teamId/members/add',
        element: <LazyWrapper><AddMemberPage /></LazyWrapper>,
      },
      {
        path: '/teams/:teamId/projects/new',
        element: <LazyWrapper><CreateProjectPage /></LazyWrapper>,
      },
      // Projects
      {
        path: '/projects/:projectId',
        element: <LazyWrapper><ProjectDetailPage /></LazyWrapper>,
      },
      {
        path: '/projects/:projectId/edit',
        element: <LazyWrapper><EditProjectPage /></LazyWrapper>,
      },
      {
        path: '/projects/:projectId/services/new',
        element: <LazyWrapper><CreateServicePage /></LazyWrapper>,
      },
      // Services
      {
        path: '/services/:serviceId',
        element: <LazyWrapper><ServiceDetailPage /></LazyWrapper>,
      },
      {
        path: '/services/:serviceId/edit',
        element: <LazyWrapper><EditServicePage /></LazyWrapper>,
      },
      {
        path: '/services/:serviceId/domains/add',
        element: <LazyWrapper><AddDomainPage /></LazyWrapper>,
      },
      // Profile
      {
        path: '/profile',
        element: <LazyWrapper><ProfilePage /></LazyWrapper>,
      },
    ],
  },
  {
    path: '*',
    element: <NotFoundPage />,
  },
])

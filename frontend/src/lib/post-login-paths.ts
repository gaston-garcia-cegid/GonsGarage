import { UserRole } from '@/types';

/**
 * Canonical first screen after login. All paths must exist under `src/app`
 * so post-login navigation never hits a Next.js 404.
 */
export function getPostLoginPath(role: string): string {
  switch (role) {
    case UserRole.ADMIN:
    case UserRole.MANAGER:
      return '/employees';
    case UserRole.EMPLOYEE:
      return '/employee/dashboard';
    case UserRole.CLIENT:
      return '/client';
    case 'technician':
      return '/employee/dashboard';
    default:
      return '/';
  }
}

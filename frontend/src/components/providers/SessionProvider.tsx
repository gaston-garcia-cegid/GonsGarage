'use client';

import { useEffect, type ReactNode } from 'react';
import { useAuthStore } from '@/stores/auth.store';

/**
 * Validates persisted session against the API on app load.
 * Auth state and API calls live in {@link useAuthStore} — do not use React Context for auth.
 */
export function SessionProvider({ children }: { children: ReactNode }) {
  useEffect(() => {
    void useAuthStore.getState().checkAuthStatus();
  }, []);

  return <>{children}</>;
}

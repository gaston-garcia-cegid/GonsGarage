'use client';

import { useEffect, useState } from 'react';
import { useAuthStore } from '@/stores/auth.store';

/**
 * Zustand `persist` rehydrates from localStorage after the first client render.
 * Until then, `user` / `token` may be null even when the session exists — avoid
 * treating that as "logged out" (e.g. redirecting to /auth/login).
 */
export function useAuthHydrationReady(): boolean {
  const [ready, setReady] = useState(() =>
    typeof window === 'undefined' ? false : useAuthStore.persist.hasHydrated(),
  );

  useEffect(() => {
    if (useAuthStore.persist.hasHydrated()) {
      let cancelled = false;
      queueMicrotask(() => {
        if (!cancelled) setReady(true);
      });
      return () => {
        cancelled = true;
      };
    }
    return useAuthStore.persist.onFinishHydration(() => {
      setReady(true);
    });
  }, []);

  return ready;
}

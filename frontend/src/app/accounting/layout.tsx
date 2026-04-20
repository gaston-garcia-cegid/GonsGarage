'use client';

import React, { useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import { isClient } from '@/types/user';

export default function AccountingLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  const { user } = useAuth();
  const router = useRouter();
  const authHydrated = useAuthHydrationReady();

  useEffect(() => {
    if (!authHydrated) return;
    if (!user) {
      router.replace('/auth/login');
      return;
    }
    if (isClient(user)) {
      router.replace('/dashboard');
    }
  }, [authHydrated, user, router]);

  if (!authHydrated || !user || isClient(user)) {
    return (
      <div className="loadingScreen">
        <div className="spinnerLg" aria-hidden />
      </div>
    );
  }

  return children;
}

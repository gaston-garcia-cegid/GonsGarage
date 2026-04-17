'use client';

import React, { Suspense, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import CarsContainer from '@/app/cars/components/CarsContainer';
import AppShell from '@/components/layout/AppShell';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import styles from './cars.module.css';

export default function CarsPage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const authHydrated = useAuthHydrationReady();

  useEffect(() => {
    if (!authHydrated) return;
    if (!user) {
      router.replace('/auth/login');
    }
  }, [authHydrated, user, router]);

  if (!authHydrated || !user) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    );
  }

  return (
    <AppShell
      user={user}
      subtitle="My Cars"
      activeNav="cars"
      onLogout={logout}
      logoVariant="branded"
    >
      <Suspense
        fallback={
          <div className="flex flex-col items-center justify-center gap-3 py-12 text-gray-600">
            <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600" />
            <span>Loading…</span>
          </div>
        }
      >
        <CarsContainer
          headerTitle="My Cars"
          headerSubtitle="Manage your registered cars"
          addButtonText="Add Car"
          className={styles.carsSection}
        />
      </Suspense>
    </AppShell>
  );
}

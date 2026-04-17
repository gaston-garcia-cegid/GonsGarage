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
      <div className="loadingScreen">
        <div className="spinnerLg" aria-hidden />
      </div>
    );
  }

  return (
    <AppShell
      user={user}
      subtitle="Os meus automóveis"
      activeNav="cars"
      onLogout={logout}
      logoVariant="branded"
    >
      <Suspense
        fallback={
          <div className="loadingStack">
            <div className="spinnerMd" aria-hidden />
            <span>A carregar…</span>
          </div>
        }
      >
        <CarsContainer
          headerTitle="Os meus automóveis"
          headerSubtitle="Gerir os seus automóveis registados"
          addButtonText="Adicionar automóvel"
          className={styles.carsSection}
        />
      </Suspense>
    </AppShell>
  );
}

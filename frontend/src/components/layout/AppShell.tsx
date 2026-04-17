'use client';

import React from 'react';
import Image from 'next/image';
import { useRouter } from 'next/navigation';
import type { User } from '@/types';
import styles from './AppShell.module.css';

export type AppShellNavId = 'dashboard' | 'cars' | 'appointments';

export interface AppShellProps {
  user: User;
  subtitle: string;
  activeNav: AppShellNavId;
  carsNavLabel?: string;
  onLogout: () => void;
  children: React.ReactNode;
  /** Use workshop logo on home link (same as cars page). */
  logoVariant?: 'svg' | 'branded';
  /** Optional row below nav (e.g. page title + primary action). */
  toolbar?: React.ReactNode;
}

export default function AppShell({
  user,
  subtitle,
  activeNav,
  carsNavLabel = 'My Cars',
  onLogout,
  children,
  logoVariant = 'svg',
  toolbar,
}: Readonly<AppShellProps>) {
  const router = useRouter();

  return (
    <div className={styles.container}>
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <button
            type="button"
            className={styles.logoSection}
            onClick={() => router.push('/')}
            aria-label="Go to home"
          >
            <div className={styles.logoIcon}>
              {logoVariant === 'branded' ? (
                <Image
                  src="/images/LogoGonsGarage.jpg"
                  alt=""
                  width={32}
                  height={32}
                  style={{ objectFit: 'cover', width: '100%', height: '100%' }}
                />
              ) : (
                <svg fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-2m-2 0H7m5 0v-5a2 2 0 012-2h2a2 2 0 012 2v5"
                  />
                </svg>
              )}
            </div>
            <div>
              <h1>GonsGarage</h1>
              <p>{subtitle}</p>
            </div>
          </button>
          <div className={styles.userSection}>
            <span>
              Welcome, {user.firstName} {user.lastName}
            </span>
            <button type="button" onClick={onLogout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      <nav className={styles.navigation} aria-label="Main">
        <button
          type="button"
          onClick={() => router.push('/dashboard')}
          className={`${styles.navButton} ${activeNav === 'dashboard' ? styles.active : ''}`}
        >
          Dashboard
        </button>
        <button
          type="button"
          onClick={() => router.push('/cars')}
          className={`${styles.navButton} ${activeNav === 'cars' ? styles.active : ''}`}
        >
          {carsNavLabel}
        </button>
        <button
          type="button"
          onClick={() => router.push('/appointments')}
          className={`${styles.navButton} ${activeNav === 'appointments' ? styles.active : ''}`}
        >
          Appointments
        </button>
      </nav>

      <main className={styles.main}>
        {toolbar ? <div className={styles.pageToolbar}>{toolbar}</div> : null}
        {children}
      </main>
    </div>
  );
}

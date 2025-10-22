// src/components/layouts/DashboardLayout.tsx
'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import Image from 'next/image';
import styles from './DashboardLayout.module.css';

interface DashboardLayoutProps {
  children: React.ReactNode;
  title: string;
  subtitle: string;
  activeTab: string;
  navigationItems: Array<{
    key: string;
    label: string;
    href: string;
  }>;
}

export default function DashboardLayout({
  children,
  title,
  subtitle,
  activeTab,
  navigationItems
}: DashboardLayoutProps) {
  const { user, logout } = useAuth();
  const router = useRouter();

  return (
    <div className={styles.container}>
      {/* Header Reutilizável */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div 
            className={styles.logoSection} 
            onClick={() => router.push('/')} 
            style={{ cursor: 'pointer' }}
          >
            <div className={styles.logoIcon}>
              <Image
                src="/images/LogoGonsGarage.jpg"
                alt="GonsGarage Logo"
                width={24}
                height={24}
                style={{ objectFit: 'contain' }}
              />
            </div>
            <div>
              <h1>GonsGarage</h1>
              <p>{subtitle}</p>
            </div>
          </div>
          
          <div className={styles.userSection}>
            <span>Welcome, {user?.first_name} {user?.last_name}</span>
            <button onClick={logout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation Dinâmica */}
      <nav className={styles.navigation}>
        {navigationItems.map((item) => (
          <button
            key={item.key}
            onClick={() => router.push(item.href)}
            className={`${styles.navButton} ${
              activeTab === item.key ? styles.active : ''
            }`}
          >
            {item.label}
          </button>
        ))}
      </nav>

      {/* Conteúdo */}
      <main className={styles.main}>
        <div className={styles.pageHeader}>
          <h2>{title}</h2>
        </div>
        {children}
      </main>
    </div>
  );
}
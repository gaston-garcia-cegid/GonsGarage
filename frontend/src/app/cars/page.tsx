'use client';

import React from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import CarsContainer from '@/app/cars/components/CarsContainer';
import Image from 'next/image';
import styles from './cars.module.css';

// Main Cars page component following Agent.md component conventions
export default function CarsPage() {
  const { user, logout } = useAuth();
  const router = useRouter();

  // Redirect if not authenticated
  React.useEffect(() => {
    if (!user) {
      router.push('/auth/login');
    }
  }, [user, router]);

  if (!user) {
    return null;
  }

  return (
    <div className={styles.container}>
      {/* Header - following Agent.md UI conventions */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logoSection} onClick={() => router.push('/')}>
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
              <p>My Cars</p>
            </div>
          </div>
          <div className={styles.userSection}>
            <span>Welcome, {user?.firstName} {user?.lastName}</span>
            <button onClick={logout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation - following Agent.md navigation patterns */}
      <nav className={styles.navigation}>
        <button 
          onClick={() => router.push('/dashboard')}
          className={styles.navButton}
        >
          Dashboard
        </button>
        <button 
          onClick={() => router.push('/cars')}
          className={`${styles.navButton} ${styles.active}`}
        >
          My Cars
        </button>
        <button 
          onClick={() => router.push('/appointments')}
          className={styles.navButton}
        >
          Appointments
        </button>
      </nav>

      {/* Main Content */}
      <main className={styles.main}>
        {/* ✅ Use shared container */}
        <CarsContainer
          headerTitle="My Cars"
          headerSubtitle="Manage your registered cars"
          addButtonText="Add Car"
          className={styles.carsSection}
        />
      </main>
    </div>
  );
}
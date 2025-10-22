// src/app/client/page.tsx (dashboard principal)
'use client';

import React, { useState, useEffect } from 'react';
import DashboardLayout from '@/components/layouts/DashboardLayout';
import { useAuthGuard } from '@/hooks/useAuthGuard';
import { getNavigationConfig } from '@/lib/navigation';
import styles from './client.module.css';

export default function ClientPage() {
  const { user, isLoading, isAuthorized } = useAuthGuard('client');
  const [cars, setCars] = useState<Car[]>([]);
  const [recentRepairs, setRecentRepairs] = useState<Repair[]>([]);
  const [upcomingAppointments, setUpcomingAppointments] = useState<Appointment[]>([]);
  const [error, setError] = useState<string | null>(null);
  
  // Só mostra loading, não faz redirecionamentos (AuthContext cuida disso)
  if (isLoading || !isAuthorized) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  const navConfig = getNavigationConfig('client');

  return (
    <DashboardLayout
      title="Dashboard"
      subtitle="Customer Dashboard"
      activeTab="dashboard"
      navigationItems={navConfig.items}
    >
      {/* Conteúdo do dashboard */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white p-6 rounded-lg shadow">
          <h3 className="text-lg font-semibold mb-2">My Cars</h3>
          <p className="text-3xl font-bold text-blue-600">3</p>
        </div>
        {/* Mais conteúdo... */}
        {error && (
          <div className={styles.errorAlert}>
            <span>{error}</span>
          </div>
        )}

        {/* Stats Grid */}
        <div className={styles.statsGrid}>
          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-2m-2 0H7m5 0v-5a2 2 0 012-2h2a2 2 0 012 2v5" />
              </svg>
            </div>
            <div>
              <h3>My Cars</h3>
              <p className={styles.statNumber}>{cars.length}</p>
            </div>
          </div>

          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div>
              <h3>Active Repairs</h3>
              <p className={styles.statNumber}>
                {recentRepairs.filter(r => r.status === 'in_progress').length}
              </p>
            </div>
          </div>

          <div className={styles.statCard}>
            <div className={styles.statIcon}>
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h3>Upcoming Appointments</h3>
              <p className={styles.statNumber}>{upcomingAppointments.length}</p>
            </div>
          </div>
        </div>

        {/* Content Grid */}
        <div className={styles.contentGrid}>
          {/* Recent Cars */}
          <div className={styles.card}>
            <div className={styles.cardHeader}>
              <h3>My Cars</h3>
              <button 
                onClick={() => router.push('/cars')}
                className={styles.linkButton}
              >
                View All
              </button>
            </div>
            <div className={styles.cardBody}>
              {cars.length === 0 ? (
                <div className={styles.emptyState}>
                  <p>No cars registered yet</p>
                  <button 
                    onClick={() => router.push('/cars')}
                    className={styles.primaryButton}
                  >
                    Add Your First Car
                  </button>
                </div>
              ) : (
                <div className={styles.carsList}>
                  {cars.slice(0, 3).map((car) => (
                    <div key={car.id} className={styles.carItem}>
                      <div className={styles.carIcon}>
                        🚗
                      </div>
                      <div className={styles.carInfo}>
                        <h4>{car.year} {car.make} {car.model}</h4>
                        <p>{car.license_plate}</p>
                      </div>
                      <button 
                        onClick={() => router.push(`/cars/${car.id}`)}
                        className={styles.viewButton}
                      >
                        View
                      </button>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Recent Repairs */}
          <div className={styles.card}>
            <div className={styles.cardHeader}>
              <h3>Recent Repairs</h3>
              <button 
                onClick={() => router.push('/cars')}
                className={styles.linkButton}
              >
                View All
              </button>
            </div>
            <div className={styles.cardBody}>
              {recentRepairs.length === 0 ? (
                <div className={styles.emptyState}>
                  <p>No repairs yet</p>
                </div>
              ) : (
                <div className={styles.repairsList}>
                  {recentRepairs.map((repair) => (
                    <div key={repair.id} className={styles.repairItem}>
                      <div className={styles.repairStatus}>
                        <span className={`${styles.statusBadge} ${styles[repair.status]}`}>
                          {repair.status.replace('_', ' ')}
                        </span>
                      </div>
                      <div className={styles.repairInfo}>
                        <h4>{repair.description}</h4>
                        <p>${repair.cost.toFixed(2)}</p>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

    </DashboardLayout>
  );
}
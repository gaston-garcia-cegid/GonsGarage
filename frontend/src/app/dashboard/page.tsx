'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient, Car, Repair, Appointment } from '@/lib/api';
import styles from './dashboard.module.css';

export default function DashboardPage() {
  const [cars, setCars] = useState<Car[]>([]);
  const [recentRepairs, setRecentRepairs] = useState<Repair[]>([]);
  const [upcomingAppointments, setUpcomingAppointments] = useState<Appointment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const { user, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    fetchDashboardData();
  }, [user, router]);

  const fetchDashboardData = async () => {
    try {
      setLoading(true);
      setError(null);

      const [carsResponse, repairsResponse, appointmentsResponse] = await Promise.all([
        apiClient.getCars(),
        apiClient.getRepairs(),
        apiClient.getAppointments(),
      ]);

      if (carsResponse.data && !carsResponse.error) {
        setCars(carsResponse.data);
      }

      if (repairsResponse.data && !repairsResponse.error) {
        setRecentRepairs(repairsResponse.data.slice(0, 5));
      }

      if (appointmentsResponse.data && !appointmentsResponse.error) {
        setUpcomingAppointments(appointmentsResponse.data.slice(0, 3));
      }
    } catch (err) {
      setError('Failed to load dashboard data');
      console.error('Dashboard error:', err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <span>Loading dashboard...</span>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logoSection}>
            <div className={styles.logoIcon}>
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-2m-2 0H7m5 0v-5a2 2 0 012-2h2a2 2 0 012 2v5" />
              </svg>
            </div>
            <div>
              <h1>GonsGarage</h1>
              <p>Customer Dashboard</p>
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

      {/* Navigation */}
      <nav className={styles.navigation}>
        <button 
          onClick={() => router.push('/dashboard')}
          className={`${styles.navButton} ${styles.active}`}
        >
          Dashboard
        </button>
        <button 
          onClick={() => router.push('/cars')}
          className={styles.navButton}
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
      </main>
    </div>
  );
}
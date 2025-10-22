'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient, Car, Repair } from '@/lib/api';
import styles from './car-details.module.css';

export default function CarDetailsPage() {
  const [car, setCar] = useState<Car | null>(null);
  const [repairs, setRepairs] = useState<Repair[]>([]);
  const [loading, setLoading] = useState(true);
  const [repairsLoading, setRepairsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const { user, logout } = useAuth();
  const router = useRouter();
  const params = useParams();
  const carId = params.id as string;

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    if (carId) {
      fetchCarDetails();
      fetchCarRepairs();
    }
  }, [user, router, carId]);

  const fetchCarDetails = async () => {
    try {
      setLoading(true);
      setError(null);

      const { data, error: apiError } = await apiClient.getCar(carId);
      
      if (data && !apiError) {
        setCar(data);
      } else {
        setError(apiError?.message || 'Failed to fetch car details');
      }
    } catch (err) {
      setError('Network error occurred');
    } finally {
      setLoading(false);
    }
  };

  const fetchCarRepairs = async () => {
    try {
      setRepairsLoading(true);

      const { data, error: apiError } = await apiClient.getRepairs(carId);
      
      if (data && !apiError) {
        setRepairs(data);
      } else {
        console.error('Failed to fetch repairs:', apiError);
      }
    } catch (err) {
      console.error('Network error fetching repairs:', err);
    } finally {
      setRepairsLoading(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
    }).format(amount);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending': return 'warning';
      case 'in_progress': return 'info';
      case 'completed': return 'success';
      case 'cancelled': return 'error';
      default: return 'default';
    }
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <span>Loading car details...</span>
      </div>
    );
  }

  if (error || !car) {
    return (
      <div className={styles.errorContainer}>
        <div className={styles.errorContent}>
          <div className={styles.errorIcon}>⚠️</div>
          <h2>Car Not Found</h2>
          <p>{error || 'The requested car could not be found.'}</p>
          <button 
            onClick={() => router.push('/cars')}
            className={styles.backButton}
          >
            Back to Cars
          </button>
        </div>
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
              <p>Car Details</p>
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

      {/* Navigation */}
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

      {/* Breadcrumbs */}
      <div className={styles.breadcrumbs}>
        <button onClick={() => router.push('/cars')} className={styles.breadcrumbLink}>
          My Cars
        </button>
        <span className={styles.breadcrumbSeparator}>›</span>
        <span className={styles.breadcrumbCurrent}>
          {car.year} {car.make} {car.model}
        </span>
      </div>

      {/* Main Content */}
      <main className={styles.main}>
        {/* Car Info Card */}
        <div className={styles.carInfoCard}>
          <div className={styles.carInfoHeader}>
            <div className={styles.carIconLarge}>🚗</div>
            <div className={styles.carTitle}>
              <h2>{car.year} {car.make} {car.model}</h2>
              <p className={styles.licensePlate}>{car.license_plate}</p>
            </div>
            <div className={styles.carActions}>
              <button 
                onClick={() => router.push(`/appointments/new?carId=${car.id}`)}
                className={styles.scheduleServiceButton}
              >
                <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Schedule Service
              </button>
            </div>
          </div>

          <div className={styles.carInfoGrid}>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Color</span>
              <span className={styles.infoValue}>{car.color}</span>
            </div>
            {car.vin && (
              <div className={styles.infoItem}>
                <span className={styles.infoLabel}>VIN</span>
                <span className={styles.infoValue}>{car.vin}</span>
              </div>
            )}
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Added</span>
              <span className={styles.infoValue}>{formatDate(car.created_at)}</span>
            </div>
            <div className={styles.infoItem}>
              <span className={styles.infoLabel}>Last Updated</span>
              <span className={styles.infoValue}>{formatDate(car.updated_at)}</span>
            </div>
          </div>
        </div>

        {/* Repairs Section */}
        <div className={styles.repairsSection}>
          <div className={styles.sectionHeader}>
            <h3>Service History</h3>
            <div className={styles.repairsStats}>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'completed').length}
                </span>
                <span className={styles.statLabel}>Completed</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'in_progress').length}
                </span>
                <span className={styles.statLabel}>In Progress</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {repairs.filter(r => r.status === 'pending').length}
                </span>
                <span className={styles.statLabel}>Pending</span>
              </div>
              <div className={styles.statItem}>
                <span className={styles.statValue}>
                  {formatCurrency(repairs.reduce((sum, r) => sum + r.cost, 0))}
                </span>
                <span className={styles.statLabel}>Total Cost</span>
              </div>
            </div>
          </div>

          {repairsLoading ? (
            <div className={styles.repairsLoading}>
              <div className={styles.spinner}></div>
              <span>Loading repairs...</span>
            </div>
          ) : repairs.length === 0 ? (
            <div className={styles.emptyRepairs}>
              <div className={styles.emptyIcon}>🔧</div>
              <h4>No Service History</h4>
              <p>This car hasn&apos;t had any services yet.</p>
              <button 
                onClick={() => router.push(`/appointments/new?carId=${car.id}`)}
                className={styles.scheduleFirstServiceButton}
              >
                Schedule First Service
              </button>
            </div>
          ) : (
            <div className={styles.repairsList}>
              {repairs.map((repair) => (
                <div key={repair.id} className={styles.repairCard}>
                  <div className={styles.repairHeader}>
                    <div className={styles.repairStatus}>
                      <span className={`${styles.statusBadge} ${styles[getStatusColor(repair.status)]}`}>
                        {repair.status.replace('_', ' ').toUpperCase()}
                      </span>
                    </div>
                    <div className={styles.repairCost}>
                      {formatCurrency(repair.cost)}
                    </div>
                  </div>

                  <div className={styles.repairContent}>
                    <h4 className={styles.repairDescription}>{repair.description}</h4>
                    
                    <div className={styles.repairMeta}>
                      <div className={styles.metaItem}>
                        <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>Created: {formatDate(repair.created_at)}</span>
                      </div>
                      
                      {repair.started_at && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>Started: {formatDate(repair.started_at)}</span>
                        </div>
                      )}
                      
                      {repair.completed_at && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>Completed: {formatDate(repair.completed_at)}</span>
                        </div>
                      )}

                      {repair.technician && (
                        <div className={styles.metaItem}>
                          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                          </svg>
                          <span>Technician: {repair.technician.first_name} {repair.technician.last_name}</span>
                        </div>
                      )}
                    </div>
                  </div>

                  <div className={styles.repairProgress}>
                    <div className={styles.progressBar}>
                      <div 
                        className={`${styles.progressFill} ${styles[getStatusColor(repair.status)]}`}
                        style={{
                          width: repair.status === 'completed' ? '100%' : 
                                repair.status === 'in_progress' ? '60%' : 
                                repair.status === 'pending' ? '20%' : '0%'
                        }}
                      ></div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </main>
    </div>
  );
}
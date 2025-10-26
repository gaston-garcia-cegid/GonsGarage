'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth, useAppointments, useCars } from '@/stores';
import styles from './appointments.module.css';

export default function AppointmentsPage() {
  const [selectedStatus, setSelectedStatus] = useState<string>('all');

  const { user, logout } = useAuth();
  const { appointments, isLoading: appointmentsLoading, error: appointmentsError, fetchAppointments } = useAppointments();
  const { cars, isLoading: carsLoading, error: carsError, fetchCars } = useCars();
  const router = useRouter();

  // Combined loading and error states
  const loading = appointmentsLoading || carsLoading;
  const error = appointmentsError || carsError;

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    // Fetch data using stores
    fetchAppointments();
    fetchCars();
  }, [user, router, fetchAppointments, fetchCars]);

  const filteredAppointments = appointments.filter(appointment => 
    selectedStatus === 'all' || appointment.status === selectedStatus
  );

  const getCarInfo = (carId: string) => {
    return cars.find(car => car.id === carId);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <span>Loading appointments...</span>
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
              <p>Appointments</p>
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
          className={styles.navButton}
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
          className={`${styles.navButton} ${styles.active}`}
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

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.controlsLeft}>
            <h2>My Appointments ({filteredAppointments.length})</h2>
            <p>Manage your service appointments</p>
          </div>
          <div className={styles.controlsRight}>
            <select 
              value={selectedStatus}
              onChange={(e) => setSelectedStatus(e.target.value)}
              className={styles.statusFilter}
            >
              <option value="all">All Status</option>
              <option value="scheduled">Scheduled</option>
              <option value="confirmed">Confirmed</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
            <button 
              onClick={() => router.push('/appointments/new')}
              className={styles.addButton}
            >
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
              Schedule Service
            </button>
          </div>
        </div>

        {/* Appointments List */}
        {filteredAppointments.length === 0 ? (
          <div className={styles.emptyState}>
            <div className={styles.emptyIcon}>📅</div>
            <h3>No appointments found</h3>
            <p>
              {selectedStatus === 'all' 
                ? 'Schedule your first service appointment'
                : `No ${selectedStatus} appointments found`
              }
            </p>
            <button 
              onClick={() => router.push('/appointments/new')}
              className={styles.primaryButton}
            >
              Schedule Service
            </button>
          </div>
        ) : (
          <div className={styles.appointmentsList}>
            {filteredAppointments.map((appointment) => {
              const car = getCarInfo(appointment.carId);
              return (
                <div key={appointment.id} className={styles.appointmentCard}>
                  <div className={styles.appointmentHeader}>
                    <div className={styles.appointmentDate}>
                      <span className={styles.dateLabel}>Scheduled</span>
                      <span className={styles.dateValue}>
                        {formatDate(appointment.scheduledAt)}
                      </span>
                    </div>
                    <span className={`${styles.statusBadge} ${styles[appointment.status]}`}>
                      {appointment.status.replace('_', ' ')}
                    </span>
                  </div>
                  
                  <div className={styles.appointmentBody}>
                    <div className={styles.serviceInfo}>
                      <h3>{appointment.serviceType}</h3>
                      {appointment.notes && (
                        <p className={styles.notes}>{appointment.notes}</p>
                      )}
                    </div>
                    
                    {car && (
                      <div className={styles.carInfo}>
                        <div className={styles.carIcon}>🚗</div>
                        <div>
                          <h4>{car.year} {car.make} {car.model}</h4>
                          <p>{car.licensePlate}</p>
                        </div>
                      </div>
                    )}
                  </div>
                  
                  <div className={styles.appointmentFooter}>
                    <span className={styles.appointmentId}>
                      ID: {appointment.id.slice(0, 8)}...
                    </span>
                    <div className={styles.appointmentActions}>
                      {appointment.status === 'scheduled' && (
                        <>
                          <button className={styles.editButton}>
                            Edit
                          </button>
                          <button className={styles.cancelButton}>
                            Cancel
                          </button>
                        </>
                      )}
                      {appointment.status === 'completed' && (
                        <button 
                          onClick={() => car && router.push(`/cars/${car.id}`)}
                          className={styles.viewButton}
                        >
                          View Car
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}
      </main>
    </div>
  );
}
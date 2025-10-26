'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuthStore } from '@/stores/auth.store';
import { useAppointments } from '@/stores/appointment.store';
import { useCarStore } from '@/stores/car.store';
import { Appointment } from '@/shared/types';
import styles from './appointments.module.css';

export default function AppointmentsPage() {
  const [selectedStatus, setSelectedStatus] = useState<string>('all');

  const { user } = useAuthStore();
  const { 
    appointments, 
    isLoading: appointmentsLoading, 
    error: appointmentsError, 
    fetchAppointments,
    cancelAppointment,
    confirmAppointment,
    completeAppointment
  } = useAppointments();
  const { cars, isLoading: carsLoading, fetchCars } = useCarStore();
  const router = useRouter();

  const loading = appointmentsLoading || carsLoading;
  const error = appointmentsError;

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
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

  const handleStatusChange = async (appointmentId: string, action: 'cancel' | 'confirm' | 'complete') => {
    try {
      let success = false;
      
      switch (action) {
        case 'cancel':
          success = await cancelAppointment(appointmentId);
          break;
        case 'confirm':
          success = await confirmAppointment(appointmentId);
          break;
        case 'complete':
          success = await completeAppointment(appointmentId);
          break;
      }

      if (success) {
        // Refresh appointments
        fetchAppointments();
      }
    } catch (error) {
      console.error(`Failed to ${action} appointment:`, error);
    }
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
            <div className={styles.logoIcon}>🔧</div>
            <div>
              <h1>GonsGarage</h1>
              <p>Appointments</p>
            </div>
          </div>
          <div className={styles.userSection}>
            <span>Welcome, {user?.firstName} {user?.lastName}</span>
            <button onClick={() => router.push('/auth/login')} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className={styles.navigation}>
        <button 
          onClick={() => router.push('/client')}
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
              <option value="in-progress">In Progress</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
            <button 
              onClick={() => router.push('/appointments/new')}
              className={styles.addButton}
            >
              ➕ Schedule Service
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
            {filteredAppointments.map((appointment: Appointment) => {
              const car = getCarInfo(appointment.carId);  // ✅ camelCase per Agent.md
              return (
                <div key={appointment.id} className={styles.appointmentCard}>
                  <div className={styles.appointmentHeader}>
                    <div className={styles.appointmentDate}>
                      <span className={styles.dateLabel}>Scheduled</span>
                      <span className={styles.dateValue}>
                        {formatDate(appointment.scheduledAt)}  {/* ✅ camelCase per Agent.md */}
                      </span>
                    </div>
                    <span className={`${styles.statusBadge} ${styles[appointment.status]}`}>
                      {appointment.status.replace('_', ' ').replace('-', ' ')}
                    </span>
                  </div>
                  
                  <div className={styles.appointmentBody}>
                    <div className={styles.serviceInfo}>
                      <h3>{appointment.serviceType}</h3>  {/* ✅ camelCase per Agent.md */}
                      {appointment.notes && (
                        <p className={styles.notes}>{appointment.notes}</p>
                      )}
                    </div>
                    
                    {car && (
                      <div className={styles.carInfo}>
                        <div className={styles.carIcon}>🚗</div>
                        <div>
                          <h4>{car.year} {car.make} {car.model}</h4>
                          <p>{car.licensePlate}</p>  {/* ✅ camelCase per Agent.md */}
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
                          <button 
                            onClick={() => handleStatusChange(appointment.id, 'confirm')}
                            className={styles.confirmButton}
                          >
                            Confirm
                          </button>
                          <button 
                            onClick={() => router.push(`/appointments/${appointment.id}/edit`)}
                            className={styles.editButton}
                          >
                            Edit
                          </button>
                          <button 
                            onClick={() => handleStatusChange(appointment.id, 'cancel')}
                            className={styles.cancelButton}
                          >
                            Cancel
                          </button>
                        </>
                      )}
                      {appointment.status === 'confirmed' && (
                        <button 
                          onClick={() => handleStatusChange(appointment.id, 'cancel')}
                          className={styles.cancelButton}
                        >
                          Cancel
                        </button>
                      )}
                      {appointment.status === 'completed' && car && (
                        <button 
                          onClick={() => router.push(`/cars/${car.id}`)}
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
'use client';

import React, { useEffect } from 'react';
import { useAppointmentStore } from '@/stores/appointment.store';
import { useCarStore } from '@/stores/car.store';
import EmptyAppointmentState from '@/components/empty-states/EmptyAppointmentState';
import AppointmentCard from '@/components/AppointmentCard';
import Link from 'next/link';
import styles from './appointments.module.css'; 

export default function AppointmentsPage() {
  const { 
    appointments, 
    isLoading: appointmentsLoading, 
    error: appointmentsError, 
    fetchAppointments 
  } = useAppointmentStore();
  
  const { 
    cars, 
    isLoading: carsLoading, 
    fetchCars 
  } = useCarStore();

  useEffect(() => {
    fetchAppointments();
    fetchCars();
  }, [fetchAppointments, fetchCars]);

  const isLoading = appointmentsLoading || carsLoading;

  if (isLoading) {
    return (
      <div className="loading-container">
        <div className="loading-spinner">Loading appointments...</div>
      </div>
    );
  }

  if (appointmentsError) {
    return (
      <div className="error-container">
        <p>Error loading appointments: {appointmentsError}</p>
        <button onClick={() => fetchAppointments()}>Retry</button>
      </div>
    );
  }

  // ✅ Show empty state based on data availability
  if (!isLoading && appointments.length === 0) {
    return <EmptyAppointmentState hasNoCars={cars.length === 0} />;
  }

  return (
    <div className="appointments-page">
      <div className="page-header">
        <h1>My Appointments ({appointments.length})</h1>
        {cars.length > 0 && (
          <Link href="/appointments/new" className="btn-primary">
            Schedule New Appointment
          </Link>
        )}
      </div>

      <div className="appointments-list">
        {appointments.map(appointment => (
          <AppointmentCard key={appointment.id} appointment={appointment} />
        ))}
      </div>
    </div>
  );
}
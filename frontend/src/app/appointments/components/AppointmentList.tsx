import React from 'react';
import { Appointment } from '@/types/appointment';
import AppointmentCard from './AppointmentCard';
import styles from '../appointments.module.css';

interface AppointmentListProps {
  appointments: Appointment[];
  onEdit: (appointment: Appointment) => void;
  onDelete: (id: string) => void;
  onViewDetails: (id: string) => void;
  onScheduleService: (id: string) => void;
}

// Appointment list component following Agent.md component conventions
export default function AppointmentList({
  appointments,
  onEdit,
  onDelete,
  onViewDetails,
  onScheduleService
}: AppointmentListProps) {
  if (appointments.length === 0) {
    return (
      <div className={styles.emptyState}>
        <div className={styles.emptyIcon}>📅</div>
        <h3>Ainda sem marcações</h3>
        <p>Marque a primeira visita para começar a usar os nossos serviços</p>
      </div>
    );
  }

  return (
    <div className={styles.appointmentsGrid}>
      {appointments.map((appointment) => (
        <AppointmentCard
          key={appointment.id}
          appointment={appointment}
          onEdit={() => onEdit(appointment)}
          onDelete={() => onDelete(appointment.id)}
          onViewDetails={() => onViewDetails(appointment.id)}
          onScheduleService={() => onScheduleService(appointment.id)}
        />
      ))}
    </div>
  );
}
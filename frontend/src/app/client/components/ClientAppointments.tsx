// src/components/client/ClientAppointments.tsx
'use client';

import { Appointment } from '@/types/appointment';
import styles from '../client.module.css';
import AppointmentsContainer from '@/app/appointments/components/AppointmentContainer';

interface ClientAppointmentsProps {
  appointments: Appointment[];
  onAddAppointment?: (appointment: Appointment) => void;
  onUpdateAppointment?: (appointments: Appointment[]) => void;
  showAddButton?: boolean;
  maxAppointments?: number;
  onScheduleService: (id: string) => void;
}

export default function ClientAppointments({ 
  appointments,
  onAddAppointment, 
  onUpdateAppointment, 
  showAddButton = true, 
  maxAppointments,
  onScheduleService
}: ClientAppointmentsProps) {
  return (
    <AppointmentsContainer
      onAddAppointment={onAddAppointment}
      onUpdateAppointment={onUpdateAppointment}
      maxAppointments={maxAppointments}
      headerTitle="Your Appointments"
      addButtonText="Add New Appointment"
    />
  );
}
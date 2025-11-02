// src/components/client/ClientAppointments.tsx
'use client';

import { Appointment } from '@/types/appointment';
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
  onAddAppointment, 
  onUpdateAppointment,
  maxAppointments,
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
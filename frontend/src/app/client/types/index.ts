// src/app/client/types/index.ts

import { Car, Repair, Appointment } from '@/shared/types';

// Types específicos do domínio cliente
export interface ClientDashboardProps {
  error: string | null;
  cars: Car[];
  recentRepairs: Repair[];
  upcomingAppointments: Appointment[];
  onNavigate: (tab: string) => void;
}

export interface ClientCarsProps {
  cars: Car[];
  onAddCar: () => void;
}

export interface ClientAppointmentsProps {
  appointments: Appointment[];
  cars: Car[];
  onScheduleService: () => void;
}
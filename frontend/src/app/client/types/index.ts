// src/app/client/types/index.ts

import {Repair } from '@/shared/types';
import { Appointment } from '@/types/appointment';
import { Car } from '@/types/car';

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
// src/app/client/hooks/useClientData.ts
'use client';

import { useEffect } from 'react';
import { useCars, useAppointments } from '@/stores';
import { Repair } from '@/shared/types';

export function useClientData() {
  const repairs: Repair[] = [];

  const { cars, isLoading: carsLoading, error: carsError, fetchCars } = useCars();
  const { appointments, isLoading: appointmentsLoading, error: appointmentsError, fetchAppointments } = useAppointments();

  // Combined loading and error states
  const loading = carsLoading || appointmentsLoading;
  const error = carsError || appointmentsError;

  useEffect(() => {
    fetchCars();
    fetchAppointments();
  }, [fetchCars, fetchAppointments]);

  return { cars, repairs, appointments, loading, error };
}
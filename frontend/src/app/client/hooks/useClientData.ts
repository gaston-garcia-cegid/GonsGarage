// src/app/client/hooks/useClientData.ts
'use client';

import { useState, useEffect } from 'react';
import { useCars, useAppointments } from '@/stores';
import { Repair } from '@/shared/types';

export function useClientData() {
  const [repairs, setRepairs] = useState<Repair[]>([]);
  
  const { cars, isLoading: carsLoading, error: carsError, fetchCars } = useCars();
  const { appointments, isLoading: appointmentsLoading, error: appointmentsError, fetchAppointments } = useAppointments();

  // Combined loading and error states
  const loading = carsLoading || appointmentsLoading;
  const error = carsError || appointmentsError;

  useEffect(() => {
    // Fetch data from stores
    fetchCars();
    fetchAppointments();
    
    // TODO: Load repairs data when repair store is available
    // For now keeping repairs empty
    setRepairs([]);
  }, [fetchCars, fetchAppointments]);

  return { cars, repairs, appointments, loading, error };
}
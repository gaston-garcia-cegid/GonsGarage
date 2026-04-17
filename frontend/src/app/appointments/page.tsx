'use client';

import React, { Suspense, useCallback, useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAppointmentStore } from '@/stores/appointment.store';
import { useCarStore } from '@/stores/car.store';
import EmptyAppointmentState from '@/components/empty-states/EmptyAppointmentState';
import AppointmentCard from '@/components/appointments/AppointmentCard';
import NewAppointmentModal from '@/components/appointments/NewAppointmentModal';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import styles from './appointments.module.css';
import emptyStyles from '@/components/empty-states/EmptyState.module.css';

function AppointmentsPageContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const scheduleQueryHandledRef = useRef(false);
  const authHydrated = useAuthHydrationReady();
  const { user, logout, isAuthenticated } = useAuth();

  const [scheduleOpen, setScheduleOpen] = useState(false);
  const [scheduleCarId, setScheduleCarId] = useState<string | undefined>();

  const {
    appointments,
    isLoading: appointmentsLoading,
    error: appointmentsError,
    fetchAppointments,
  } = useAppointmentStore();

  const { cars, isLoading: carsLoading, fetchCars } = useCarStore();

  useEffect(() => {
    if (!authHydrated) return;
    if (!user || !isAuthenticated) {
      router.replace('/auth/login');
    }
  }, [authHydrated, user, isAuthenticated, router]);

  useEffect(() => {
    if (!authHydrated || !user || !isAuthenticated) return;
    fetchAppointments();
    fetchCars();
  }, [authHydrated, user, isAuthenticated, fetchAppointments, fetchCars]);

  const scheduleFlag = searchParams.get('schedule');
  const carIdFromQuery = searchParams.get('carId');

  useEffect(() => {
    if (scheduleFlag !== '1') {
      scheduleQueryHandledRef.current = false;
      return;
    }
    if (scheduleQueryHandledRef.current) return;
    scheduleQueryHandledRef.current = true;
    setScheduleOpen(true);
    if (carIdFromQuery) setScheduleCarId(carIdFromQuery);
    router.replace('/appointments', { scroll: false });
  }, [scheduleFlag, carIdFromQuery, router]);

  useEffect(() => {
    if (!authHydrated || !user || !isAuthenticated) return;
    if (carsLoading) return;
    if (cars.length > 0) return;
    const t = setTimeout(() => {
      router.replace('/cars?addCar=1');
    }, 3200);
    return () => clearTimeout(t);
  }, [authHydrated, user, isAuthenticated, carsLoading, cars.length, router]);

  const openSchedule = useCallback((carId?: string) => {
    if (cars.length === 0) {
      router.replace('/cars?addCar=1');
      return;
    }
    setScheduleCarId(carId);
    setScheduleOpen(true);
  }, [cars.length, router]);

  const closeSchedule = useCallback(() => {
    setScheduleOpen(false);
    setScheduleCarId(undefined);
  }, []);

  if (!authHydrated || !user || !isAuthenticated) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
      </div>
    );
  }

  // Only block the full page while the shell has no data yet — not when cars refresh in the background
  const isBootstrapping =
    (appointmentsLoading && appointments.length === 0 && !appointmentsError) ||
    (carsLoading && cars.length === 0);

  if (isBootstrapping) {
    return (
      <AppShell user={user} subtitle="Appointments" activeNav="appointments" onLogout={logout}>
        <div className="flex flex-col items-center justify-center gap-3 py-16 text-gray-600">
          <div className="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600" />
          <span>Loading appointments...</span>
        </div>
      </AppShell>
    );
  }

  if (appointmentsError) {
    return (
      <AppShell user={user} subtitle="Appointments" activeNav="appointments" onLogout={logout}>
        <div className="rounded-md border border-red-200 bg-red-50 p-4 text-red-700">
          <p>Error loading appointments: {appointmentsError}</p>
          <button type="button" className="btn-primary mt-3" onClick={() => fetchAppointments()}>
            Retry
          </button>
        </div>
      </AppShell>
    );
  }

  if (cars.length === 0) {
    const toolbarNoCars = <h1>Mis turnos</h1>;
    return (
      <AppShell
        user={user}
        subtitle="Appointments"
        activeNav="appointments"
        onLogout={logout}
        toolbar={toolbarNoCars}
      >
        <div className={emptyStyles.emptyState}>
          <div className={emptyStyles.emptyStateIcon} aria-hidden>
            🚗
          </div>
          <h3 className={emptyStyles.emptyStateTitle}>Aún no tenés vehículos cargados</h3>
          <p className={emptyStyles.emptyStateDescription}>
            Para agendar un turno primero registrá al menos un auto en <strong>Mis autos</strong>.
          </p>
          <p className={`${emptyStyles.emptyStateDescription} ${styles.noCarsRedirectHint}`}>
            Te llevamos a Mis autos en unos segundos…
          </p>
          <button
            type="button"
            className={emptyStyles.emptyStateButton}
            onClick={() => router.replace('/cars?addCar=1')}
          >
            Ir a Mis autos ahora
          </button>
        </div>
      </AppShell>
    );
  }

  const toolbar = (
    <>
      <h1>Mis turnos ({appointments.length})</h1>
      <button type="button" className="btn-primary" onClick={() => openSchedule()}>
        Nuevo turno
      </button>
    </>
  );

  return (
    <>
      <AppShell
        user={user}
        subtitle="Appointments"
        activeNav="appointments"
        onLogout={logout}
        toolbar={toolbar}
      >
        {appointments.length === 0 ? (
          <EmptyAppointmentState onSchedule={() => openSchedule()} />
        ) : (
          <div className={styles.list}>
            {appointments.map((appointment) => (
              <AppointmentCard
                key={appointment.id}
                appointment={appointment}
                onReschedule={(carId) => openSchedule(carId)}
              />
            ))}
          </div>
        )}
      </AppShell>

      <NewAppointmentModal
        isOpen={scheduleOpen}
        onClose={closeSchedule}
        initialCarId={scheduleCarId}
        onCreated={() => {
          fetchAppointments().catch(() => {});
        }}
      />
    </>
  );
}

function AppointmentsFallback() {
  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
    </div>
  );
}

export default function AppointmentsPage() {
  return (
    <Suspense fallback={<AppointmentsFallback />}>
      <AppointmentsPageContent />
    </Suspense>
  );
}

'use client';

import React, { useCallback, useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/stores';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import AppShell from '@/components/layout/AppShell';
import { Button } from '@/components/ui/button';
import { apiClient, type Car, type ServiceJob } from '@/lib/api';
import styles from './workshop.module.css';

function statusLabel(s: ServiceJob['status']): string {
  switch (s) {
    case 'open':
      return 'Aberta';
    case 'in_progress':
      return 'Em curso';
    case 'closed':
      return 'Fechada';
    case 'cancelled':
      return 'Cancelada';
    default:
      return s;
  }
}

export default function WorkshopListPage() {
  const { user, logout } = useAuth();
  const authHydrated = useAuthHydrationReady();
  const router = useRouter();
  const [cars, setCars] = useState<Car[]>([]);
  const [carId, setCarId] = useState<string>('');
  const [jobs, setJobs] = useState<ServiceJob[]>([]);
  const [todayJobs, setTodayJobs] = useState<ServiceJob[]>([]);
  const [err, setErr] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const loadCars = useCallback(async () => {
    if (!user) return;
    setErr(null);
    setLoading(true);
    const { data, error } = await apiClient.getCars();
    setLoading(false);
    if (error) {
      setErr(error.message);
      return;
    }
    setCars(data ?? []);
  }, [user]);

  const loadTodayUTC = useCallback(async () => {
    if (!user) return;
    const d = new Date();
    const openedOn = `${d.getUTCFullYear()}-${String(d.getUTCMonth() + 1).padStart(2, '0')}-${String(d.getUTCDate()).padStart(2, '0')}`;
    const { data, error } = await apiClient.listServiceJobsByOpenedOn(openedOn);
    if (error) {
      setTodayJobs([]);
      return;
    }
    setTodayJobs(data ?? []);
  }, [user]);

  useEffect(() => {
    if (authHydrated && user) {
      void loadCars();
      void loadTodayUTC();
    }
  }, [authHydrated, user, loadCars, loadTodayUTC]);

  useEffect(() => {
    if (cars.length > 0 && !carId) {
      setCarId(cars[0].id);
    }
  }, [cars, carId]);

  const loadJobs = useCallback(async () => {
    if (!carId) {
      setJobs([]);
      return;
    }
    setErr(null);
    setLoading(true);
    const { data, error } = await apiClient.listServiceJobsByCar(carId);
    setLoading(false);
    if (error) {
      setErr(error.message);
      setJobs([]);
      return;
    }
    setJobs(data ?? []);
  }, [carId]);

  useEffect(() => {
    if (authHydrated && carId) {
      void loadJobs();
    }
  }, [authHydrated, carId, loadJobs]);

  const onNewVisit = async () => {
    if (!carId) return;
    setErr(null);
    setLoading(true);
    const { data, error } = await apiClient.createServiceJob(carId);
    setLoading(false);
    if (error) {
      setErr(error.message);
      return;
    }
    if (data?.id) {
      router.push(`/workshop/${data.id}`);
    }
  };

  if (!authHydrated || !user) return null;

  return (
    <AppShell
      user={user}
      subtitle="Taller"
      activeNav="workshop"
      onLogout={logout}
      logoVariant="branded"
      toolbar={
        <>
          <h1>Visitas (service jobs)</h1>
          <Button type="button" variant="outline" asChild>
            <Link href="/workshop/recepcion">Receção rápida</Link>
          </Button>
          <Button type="button" disabled={!carId || loading} onClick={() => void onNewVisit()}>
            Nova visita
          </Button>
        </>
      }
    >
      <section className={styles.todaySection} aria-labelledby="today-heading">
        <h2 id="today-heading" className={styles.sectionTitle}>
          Hoje (UTC)
        </h2>
        <p className={styles.hint}>Visitas abertas neste dia civil (UTC), parâmetro opened_on.</p>
        {todayJobs.length === 0 ? (
          <p className={styles.muted}>Nenhuma visita neste dia.</p>
        ) : (
          <table className={styles.table}>
            <thead>
              <tr>
                <th>Estado</th>
                <th>Aberta (UTC)</th>
                <th>ID</th>
                <th></th>
              </tr>
            </thead>
            <tbody>
              {todayJobs.map(j => (
                <tr key={j.id}>
                  <td>{statusLabel(j.status)}</td>
                  <td>{j.opened_at ? new Date(j.opened_at).toLocaleString('en-GB', { timeZone: 'UTC' }) : '—'}</td>
                  <td style={{ fontSize: '0.85rem', fontFamily: 'monospace' }}>{j.id.slice(0, 8)}…</td>
                  <td>
                    <Button type="button" variant="outline" size="sm" asChild>
                      <Link href={`/workshop/${j.id}`}>Detalhe</Link>
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </section>
      <p className={styles.hint}>Seleccione um veículo para listar e abrir visitas de oficina.</p>
      <div className={styles.select}>
        <label htmlFor="workshop-car">Viatura</label>
        <select
          id="workshop-car"
          value={carId}
          onChange={e => {
            setCarId(e.target.value);
          }}
        >
          <option value="">—</option>
          {cars.map(c => (
            <option key={c.id} value={c.id}>
              {c.license_plate} — {c.make} {c.model}
            </option>
          ))}
        </select>
      </div>
      {err ? <p className={styles.err}>{err}</p> : null}
      {loading ? <p>A carregar…</p> : null}
      <table className={styles.table}>
        <thead>
          <tr>
            <th>Estado</th>
            <th>Aberta</th>
            <th>ID</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {jobs.map(j => (
            <tr key={j.id}>
              <td>{statusLabel(j.status)}</td>
              <td>{j.opened_at ? new Date(j.opened_at).toLocaleString('pt-PT') : '—'}</td>
              <td style={{ fontSize: '0.85rem', fontFamily: 'monospace' }}>{j.id.slice(0, 8)}…</td>
              <td>
                <Button type="button" variant="outline" size="sm" asChild>
                  <Link href={`/workshop/${j.id}`}>Detalhe</Link>
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </AppShell>
  );
}

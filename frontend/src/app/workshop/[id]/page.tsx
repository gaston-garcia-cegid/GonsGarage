'use client';

import React, { useCallback, useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '@/stores';
import { useAuthHydrationReady } from '@/hooks/useAuthHydrationReady';
import AppShell from '@/components/layout/AppShell';
import { Button } from '@/components/ui/button';
import { apiClient, type ServiceJobDetail } from '@/lib/api';
import styles from '../workshop.module.css';

export default function WorkshopDetailPage() {
  const { id: jobId } = useParams<{ id: string }>();
  const { user, logout } = useAuth();
  const authHydrated = useAuthHydrationReady();
  const [detail, setDetail] = useState<ServiceJobDetail | null>(null);
  const [err, setErr] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [rKm, setRKm] = useState('');
  const [rNotes, setRNotes] = useState('');
  const [hKm, setHKm] = useState('');
  const [hNotes, setHNotes] = useState('');

  const load = useCallback(async () => {
    if (!jobId) return;
    setErr(null);
    const { data, error } = await apiClient.getServiceJob(jobId);
    if (error) {
      setErr(error.message);
      return;
    }
    setDetail(data ?? null);
  }, [jobId]);

  useEffect(() => {
    if (authHydrated && user && jobId) {
      void load();
    }
  }, [authHydrated, user, jobId, load]);

  const onSaveReception = async () => {
    if (!jobId) return;
    const km = parseInt(rKm, 10);
    if (Number.isNaN(km) || km < 0) {
      setErr('Indique quilometragem (km) válida.');
      return;
    }
    setErr(null);
    setSaving(true);
    const { error } = await apiClient.putServiceJobReception(jobId, {
      odometer_km: km,
      general_notes: rNotes.trim() || undefined,
    });
    setSaving(false);
    if (error) {
      setErr(error.message);
      return;
    }
    setRKm('');
    setRNotes('');
    await load();
  };

  const onSaveHandover = async () => {
    if (!jobId) return;
    const km = parseInt(hKm, 10);
    if (Number.isNaN(km) || km < 0) {
      setErr('Indique quilometragem (km) válida.');
      return;
    }
    setErr(null);
    setSaving(true);
    const { error } = await apiClient.putServiceJobHandover(jobId, {
      odometer_km: km,
      general_notes: hNotes.trim() || undefined,
    });
    setSaving(false);
    if (error) {
      setErr(error.message);
      return;
    }
    setHKm('');
    setHNotes('');
    await load();
  };

  if (!authHydrated || !user) return null;
  return (
    <AppShell
      user={user}
      subtitle="Taller — visita"
      activeNav="workshop"
      onLogout={logout}
      logoVariant="branded"
      toolbar={
        <>
          <h1>Visita</h1>
          <Button type="button" variant="outline" asChild>
            <Link href="/workshop">Voltar à lista</Link>
          </Button>
        </>
      }
    >
      {err ? <p className={styles.err}>{err}</p> : null}
      {detail ? (
        <>
          <p>
            <strong>Estado:</strong> {detail.job.status}
          </p>
          {detail.reception ? (
            <p className={styles.hint}>
              Recepção: {detail.reception.odometer_km} km — {detail.reception.general_notes || '—'}
            </p>
          ) : null}
          {detail.handover ? (
            <p className={styles.hint}>
              Entrega: {detail.handover.odometer_km} km — {detail.handover.general_notes || '—'}
            </p>
          ) : null}
          {detail.job.status !== 'closed' && detail.job.status !== 'cancelled' ? (
            <>
              <h2>Recepção</h2>
              <div className={styles.formGrid}>
                <label>
                  Km
                  <input
                    type="number"
                    min={0}
                    value={rKm}
                    onChange={e => setRKm(e.target.value)}
                    placeholder="ex. 120000"
                  />
                </label>
                <label>
                  Notas
                  <textarea
                    value={rNotes}
                    onChange={e => setRNotes(e.target.value)}
                    rows={3}
                    placeholder="Níveis, pneus, etc."
                  />
                </label>
                <div className={styles.actions}>
                  <Button type="button" disabled={saving} onClick={() => void onSaveReception()}>
                    Gravar recepção
                  </Button>
                </div>
              </div>
              {detail.reception ? (
                <>
                  <h2>Entrega</h2>
                  <div className={styles.formGrid}>
                    <label>
                      Km
                      <input
                        type="number"
                        min={0}
                        value={hKm}
                        onChange={e => setHKm(e.target.value)}
                        placeholder="ex. 120050"
                      />
                    </label>
                    <label>
                      Notas
                      <textarea
                        value={hNotes}
                        onChange={e => setHNotes(e.target.value)}
                        rows={3}
                      />
                    </label>
                    <div className={styles.actions}>
                      <Button type="button" disabled={saving} onClick={() => void onSaveHandover()}>
                        Gravar entrega e fechar
                      </Button>
                    </div>
                  </div>
                </>
              ) : null}
            </>
          ) : (
            <p className={styles.hint}>Visita concluída.</p>
          )}
        </>
      ) : !err ? (
        <p>A carregar…</p>
      ) : null}
    </AppShell>
  );
}

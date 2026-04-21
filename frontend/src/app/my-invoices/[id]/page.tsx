'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useParams } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import type { IssuedInvoice } from '@/types/accounting';
import styles from '../../accounting/accounting.module.css';
import { AppLoading } from '@/components/ui/AppLoading';

export default function MyInvoiceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { user, logout } = useAuth();
  const [row, setRow] = useState<IssuedInvoice | null>(null);
  const [notes, setNotes] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const load = useCallback(async () => {
    if (!id) return;
    setError(null);
    const res = await issuedInvoiceService.get(id);
    if (res.success && res.data) {
      setRow(res.data);
      setNotes(res.data.notes ?? '');
    } else {
      setError(res.error?.message ?? 'Fatura não encontrada.');
    }
  }, [id]);

  useEffect(() => {
    void load();
  }, [load]);

  if (!user) return null;

  async function onSaveNotes(e: React.FormEvent) {
    e.preventDefault();
    if (!id) return;
    setSaving(true);
    setError(null);
    const res = await issuedInvoiceService.patchNotes(id, notes);
    setSaving(false);
    if (res.success && res.data) {
      setRow(res.data);
      return;
    }
    setError(res.error?.message ?? 'Erro ao guardar as notas.');
  }

  return (
    <AppShell
      user={user}
      subtitle="Detalhe da fatura"
      activeNav="my_invoices"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        <Link href="/my-invoices" className={styles.mutedLink}>
          ← As minhas faturas
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      {row ? (
        <>
          <div className={styles.readonlyBlock}>
            <div>
              <strong>Valor:</strong> {row.amount.toFixed(2)} €
            </div>
            <div>
              <strong>Estado:</strong> {row.status}
            </div>
            <div>
              <strong>Criada:</strong> {row.createdAt?.slice(0, 19).replace('T', ' ') ?? '—'}
            </div>
          </div>
          <form className={styles.form} onSubmit={onSaveNotes}>
            <div className={styles.field}>
              <label htmlFor="notes">As suas notas (editável)</label>
              <textarea id="notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
            </div>
            <div className={styles.rowActions}>
              <button type="submit" className={styles.submitButton} disabled={saving}>
                {saving ? 'A guardar…' : 'Guardar notas'}
              </button>
            </div>
          </form>
        </>
      ) : !error ? (
        <div className="loadingStack" aria-busy="true">
          <AppLoading size="md" aria-busy={false} />
          <span>A carregar…</span>
        </div>
      ) : null}
    </AppShell>
  );
}

'use client';

import React, { useCallback, useEffect, useOptimistic, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import type { IssuedInvoice } from '@/types/accounting';
import styles from '../../accounting/accounting.module.css';
import { AppLoading } from '@/components/ui/AppLoading';

export type MyInvoiceDetailClientProps = Readonly<{
  invoiceId: string;
  /** When set (cookie-based server fetch), skip the first client GET for the same id. */
  initialRow: IssuedInvoice | null;
}>;

export default function MyInvoiceDetailClient({ invoiceId, initialRow }: MyInvoiceDetailClientProps) {
  const { user, logout } = useAuth();
  const [row, setRow] = useState<IssuedInvoice | null>(initialRow);
  const [notes, setNotes] = useState(initialRow?.notes ?? '');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const [optimisticRow, mergeOptimisticNotes] = useOptimistic(
    row,
    (current, nextNotes: string) => {
      if (!current) return current;
      return { ...current, notes: nextNotes };
    },
  );

  const shownRow = optimisticRow;

  const load = useCallback(async () => {
    if (!invoiceId) return;
    setError(null);
    const res = await issuedInvoiceService.get(invoiceId);
    if (res.success && res.data) {
      setRow(res.data);
      setNotes(res.data.notes ?? '');
    } else {
      setError(res.error?.message ?? 'Fatura não encontrada.');
    }
  }, [invoiceId]);

  useEffect(() => {
    if (initialRow) {
      setRow(initialRow);
      setNotes(initialRow.notes ?? '');
      return;
    }
    void load();
  }, [initialRow, load]);

  if (!user) return null;

  async function saveNotesAction(formData: FormData) {
    if (!invoiceId) return;
    const snapshot = row;
    if (!snapshot) return;
    const notesToSave = String(formData.get('notes') ?? notes);
    setError(null);
    mergeOptimisticNotes(notesToSave);
    setSaving(true);
    try {
      const res = await issuedInvoiceService.patchNotes(invoiceId, notesToSave);
      if (res.success && res.data) {
        setRow(res.data);
        setNotes(res.data.notes ?? '');
        return;
      }
      setError(res.error?.message ?? 'Erro ao guardar as notas.');
      setRow(snapshot);
    } finally {
      setSaving(false);
    }
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
      {shownRow ? (
        <>
          <span className="sr-only" data-testid="invoice-notes-optimistic">
            {shownRow.notes}
          </span>
          <div className={styles.readonlyBlock}>
            <div>
              <strong>Valor:</strong> {shownRow.amount.toFixed(2)} €
            </div>
            <div>
              <strong>Estado:</strong> {shownRow.status}
            </div>
            <div>
              <strong>Criada:</strong> {shownRow.createdAt?.slice(0, 19).replace('T', ' ') ?? '—'}
            </div>
          </div>
          <form className={styles.form} action={saveNotesAction}>
            <div className={styles.field}>
              <label htmlFor="notes">As suas notas (editável)</label>
              <textarea id="notes" name="notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
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

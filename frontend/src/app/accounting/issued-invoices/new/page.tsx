'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import styles from '../../accounting.module.css';
import { Button } from '@/components/ui/button';

export default function NewIssuedInvoicePage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [customerId, setCustomerId] = useState('');
  const [amount, setAmount] = useState('');
  const [status, setStatus] = useState('open');
  const [notes, setNotes] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  if (!user) return null;

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    const amt = Number.parseFloat(amount.replace(',', '.'));
    if (Number.isNaN(amt)) {
      setError('Indique um valor numérico válido.');
      return;
    }
    setSaving(true);
    setError(null);
    const res = await issuedInvoiceService.createStaff({
      customerId: customerId.trim(),
      amount: amt,
      status: status.trim() || 'open',
      notes: notes.trim(),
    });
    setSaving(false);
    if (res.success && res.data?.id) {
      router.replace(`/accounting/issued-invoices/${res.data.id}`);
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <AppShell
      user={user}
      subtitle="Nova fatura emitida"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        <Link href="/accounting/issued-invoices" className={styles.mutedLink}>
          ← Faturas emitidas
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="cid">ID do cliente (UUID)</label>
          <input id="cid" value={customerId} onChange={(ev) => setCustomerId(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="amt">Valor</label>
          <input id="amt" value={amount} onChange={(ev) => setAmount(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="st">Estado</label>
          <input id="st" value={status} onChange={(ev) => setStatus(ev.target.value)} placeholder="open" />
        </div>
        <div className={styles.field}>
          <label htmlFor="notes">Notas</label>
          <textarea id="notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
        </div>
        <div className={styles.rowActions}>
          <Button type="submit" className={styles.submitButton} disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </Button>
        </div>
      </form>
    </AppShell>
  );
}

'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { billingDocumentService } from '@/lib/services/billing-document.service';
import type { BillingDocumentKind } from '@/types/accounting';
import styles from '../../accounting.module.css';
import { Button } from '@/components/ui/button';

const KINDS: { value: BillingDocumentKind; label: string }[] = [
  { value: 'client_invoice', label: 'Fatura cliente' },
  { value: 'payroll', label: 'Salários' },
  { value: 'irs', label: 'IRS' },
  { value: 'other', label: 'Outro' },
];

export default function NewBillingDocumentPage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [kind, setKind] = useState<BillingDocumentKind>('other');
  const [title, setTitle] = useState('');
  const [amount, setAmount] = useState('');
  const [customerId, setCustomerId] = useState('');
  const [reference, setReference] = useState('');
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
    const cid = customerId.trim();
    const res = await billingDocumentService.create({
      kind,
      title: title.trim(),
      amount: amt,
      customerId: cid ? cid : undefined,
      reference: reference.trim(),
      notes: notes.trim(),
    });
    setSaving(false);
    if (res.success && res.data?.id) {
      router.replace(`/accounting/billing-documents/${res.data.id}`);
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <AppShell
      user={user}
      subtitle="Novo documento"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        <Link href="/accounting/billing-documents" className={styles.mutedLink}>
          ← Documentos
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="kind">Tipo</label>
          <select id="kind" value={kind} onChange={(ev) => setKind(ev.target.value as BillingDocumentKind)}>
            {KINDS.map((k) => (
              <option key={k.value} value={k.value}>
                {k.label}
              </option>
            ))}
          </select>
        </div>
        <div className={styles.field}>
          <label htmlFor="title">Título</label>
          <input id="title" value={title} onChange={(ev) => setTitle(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="amt">Valor</label>
          <input id="amt" value={amount} onChange={(ev) => setAmount(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="cid">ID cliente (opcional, UUID)</label>
          <input id="cid" value={customerId} onChange={(ev) => setCustomerId(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="ref">Referência</label>
          <input id="ref" value={reference} onChange={(ev) => setReference(ev.target.value)} required />
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

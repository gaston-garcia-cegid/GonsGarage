'use client';

import React, { useState } from 'react';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import styles from '../accounting.module.css';
import { Button } from '@/components/ui/button';

export interface IssuedInvoiceCreateFormProps {
  onSuccess: () => void;
  onCancel: () => void;
}

export function IssuedInvoiceCreateForm({ onSuccess, onCancel }: Readonly<IssuedInvoiceCreateFormProps>) {
  const [customerId, setCustomerId] = useState('');
  const [amount, setAmount] = useState('');
  const [status, setStatus] = useState('open');
  const [notes, setNotes] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

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
      onSuccess();
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <div>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="ii-cid">ID do cliente (UUID)</label>
          <input id="ii-cid" value={customerId} onChange={(ev) => setCustomerId(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ii-amt">Valor</label>
          <input id="ii-amt" value={amount} onChange={(ev) => setAmount(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ii-st">Estado</label>
          <input id="ii-st" value={status} onChange={(ev) => setStatus(ev.target.value)} placeholder="open" />
        </div>
        <div className={styles.field}>
          <label htmlFor="ii-notes">Notas</label>
          <textarea id="ii-notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
        </div>
        <div className={styles.rowActions}>
          <Button type="button" variant="outline" onClick={onCancel} disabled={saving}>
            Cancelar
          </Button>
          <Button type="submit" className={styles.submitButton} disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </Button>
        </div>
      </form>
    </div>
  );
}

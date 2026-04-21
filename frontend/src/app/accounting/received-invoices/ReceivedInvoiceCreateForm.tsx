'use client';

import React, { useState } from 'react';
import { receivedInvoiceService } from '@/lib/services/received-invoice.service';
import styles from '../accounting.module.css';
import { Button } from '@/components/ui/button';

export interface ReceivedInvoiceCreateFormProps {
  onSuccess: () => void;
  onCancel: () => void;
}

export function ReceivedInvoiceCreateForm({ onSuccess, onCancel }: Readonly<ReceivedInvoiceCreateFormProps>) {
  const [supplierId, setSupplierId] = useState('');
  const [vendorName, setVendorName] = useState('');
  const [category, setCategory] = useState('');
  const [amount, setAmount] = useState('');
  const [invoiceDate, setInvoiceDate] = useState('');
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
    const sid = supplierId.trim();
    const res = await receivedInvoiceService.create({
      supplierId: sid ? sid : undefined,
      vendorName: vendorName.trim(),
      category: category.trim(),
      amount: amt,
      invoiceDate: invoiceDate.trim(),
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
          <label htmlFor="ri-sid">ID fornecedor (opcional, UUID)</label>
          <input id="ri-sid" value={supplierId} onChange={(ev) => setSupplierId(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="ri-vendor">Nome do fornecedor</label>
          <input id="ri-vendor" value={vendorName} onChange={(ev) => setVendorName(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ri-cat">Categoria</label>
          <input id="ri-cat" value={category} onChange={(ev) => setCategory(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ri-amt">Valor</label>
          <input id="ri-amt" value={amount} onChange={(ev) => setAmount(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ri-d">Data da fatura</label>
          <input id="ri-d" type="date" value={invoiceDate} onChange={(ev) => setInvoiceDate(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="ri-notes">Notas</label>
          <textarea id="ri-notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
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

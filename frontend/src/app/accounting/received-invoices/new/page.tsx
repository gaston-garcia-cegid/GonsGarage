'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { receivedInvoiceService } from '@/lib/services/received-invoice.service';
import styles from '../../accounting.module.css';

export default function NewReceivedInvoicePage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [supplierId, setSupplierId] = useState('');
  const [vendorName, setVendorName] = useState('');
  const [category, setCategory] = useState('');
  const [amount, setAmount] = useState('');
  const [invoiceDate, setInvoiceDate] = useState('');
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
      router.replace(`/accounting/received-invoices/${res.data.id}`);
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <AppShell
      user={user}
      subtitle="Nova fatura recebida"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        <Link href="/accounting/received-invoices" className={styles.mutedLink}>
          ← Faturas recebidas
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="sid">ID fornecedor (opcional, UUID)</label>
          <input id="sid" value={supplierId} onChange={(ev) => setSupplierId(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="vendor">Nome do fornecedor</label>
          <input id="vendor" value={vendorName} onChange={(ev) => setVendorName(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="cat">Categoria</label>
          <input id="cat" value={category} onChange={(ev) => setCategory(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="amt">Valor</label>
          <input id="amt" value={amount} onChange={(ev) => setAmount(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="d">Data da fatura</label>
          <input id="d" type="date" value={invoiceDate} onChange={(ev) => setInvoiceDate(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="notes">Notas</label>
          <textarea id="notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
        </div>
        <div className={styles.rowActions}>
          <button type="submit" className={styles.submitButton} disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </button>
        </div>
      </form>
    </AppShell>
  );
}

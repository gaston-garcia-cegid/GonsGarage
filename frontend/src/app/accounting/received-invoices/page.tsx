'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { receivedInvoiceService } from '@/lib/services/received-invoice.service';
import type { ReceivedInvoice } from '@/types/accounting';
import styles from '../accounting.module.css';

export default function ReceivedInvoicesListPage() {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<ReceivedInvoice[]>([]);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const res = await receivedInvoiceService.list();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar as faturas.');
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  if (!user) return null;

  return (
    <AppShell
      user={user}
      subtitle="Contabilidade — Faturas recebidas"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <div className={styles.toolbar}>
        <h1>Faturas recebidas</h1>
        <Link href="/accounting/received-invoices/new" className={styles.primaryLink}>
          Nova fatura
        </Link>
      </div>
      <p className={styles.intro}>
        <Link href="/accounting" className={styles.mutedLink}>
          ← Contabilidade
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Fornecedor</th>
              <th>Categoria</th>
              <th>Valor</th>
              <th>Data</th>
            </tr>
          </thead>
          <tbody>
            {items.map((r) => (
              <tr key={r.id}>
                <td>
                  <Link href={`/accounting/received-invoices/${r.id}`}>{r.vendorName}</Link>
                </td>
                <td>{r.category}</td>
                <td>{r.amount.toFixed(2)}</td>
                <td>{r.invoiceDate?.slice(0, 10) ?? '—'}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={4}>
                  Sem registos. <Link href="/accounting/received-invoices/new">Criar o primeiro</Link>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import type { IssuedInvoice } from '@/types/accounting';
import styles from '../accounting.module.css';

export default function IssuedInvoicesStaffListPage() {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<IssuedInvoice[]>([]);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const res = await issuedInvoiceService.listStaff();
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
      subtitle="Contabilidade — Faturas emitidas"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <div className={styles.toolbar}>
        <h1>Faturas emitidas (clientes)</h1>
        <Link href="/accounting/issued-invoices/new" className={styles.primaryLink}>
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
              <th>Cliente (ID)</th>
              <th>Valor</th>
              <th>Estado</th>
            </tr>
          </thead>
          <tbody>
            {items.map((inv) => (
              <tr key={inv.id}>
                <td>
                  <Link href={`/accounting/issued-invoices/${inv.id}`}>{inv.customerId}</Link>
                </td>
                <td>{inv.amount.toFixed(2)}</td>
                <td>{inv.status}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={3}>
                  Sem faturas. <Link href="/accounting/issued-invoices/new">Criar a primeira</Link>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

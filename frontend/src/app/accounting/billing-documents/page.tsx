'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { billingDocumentService } from '@/lib/services/billing-document.service';
import type { BillingDocument } from '@/types/accounting';
import styles from '../accounting.module.css';

const kindLabel: Record<string, string> = {
  client_invoice: 'Fatura cliente',
  payroll: 'Salários',
  irs: 'IRS',
  other: 'Outro',
};

export default function BillingDocumentsListPage() {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<BillingDocument[]>([]);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const res = await billingDocumentService.list();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar os documentos.');
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  if (!user) return null;

  return (
    <AppShell
      user={user}
      subtitle="Contabilidade — Documentos"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <div className={styles.toolbar}>
        <h1>Documentos de faturação</h1>
        <Link href="/accounting/billing-documents/new" className={styles.primaryLink}>
          Novo documento
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
              <th>Tipo</th>
              <th>Título</th>
              <th>Valor</th>
            </tr>
          </thead>
          <tbody>
            {items.map((d) => (
              <tr key={d.id}>
                <td>{kindLabel[d.kind] ?? d.kind}</td>
                <td>
                  <Link href={`/accounting/billing-documents/${d.id}`}>{d.title}</Link>
                </td>
                <td>{d.amount.toFixed(2)}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={3}>
                  Sem registos. <Link href="/accounting/billing-documents/new">Criar o primeiro</Link>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

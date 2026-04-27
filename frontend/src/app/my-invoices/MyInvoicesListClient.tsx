'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import type { IssuedInvoice } from '@/types/accounting';
import styles from '../accounting/accounting.module.css';

export type MyInvoicesListClientProps = Readonly<{
  /** When non-empty (e.g. future cookie session + server fetch), skip the initial client list fetch. */
  initialItems: IssuedInvoice[];
}>;

export default function MyInvoicesListClient({ initialItems }: MyInvoicesListClientProps) {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<IssuedInvoice[]>(initialItems);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const res = await issuedInvoiceService.listMine();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar as suas faturas.');
  }, []);

  useEffect(() => {
    if (initialItems.length > 0) {
      setItems(initialItems);
      return;
    }
    void load();
  }, [initialItems, load]);

  if (!user) return null;

  return (
    <AppShell
      user={user}
      subtitle="As minhas faturas"
      activeNav="my_invoices"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>Faturas emitidas pela oficina para si. Pode atualizar as notas no detalhe.</p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Valor</th>
              <th>Estado</th>
              <th>Data</th>
            </tr>
          </thead>
          <tbody>
            {items.map((inv) => (
              <tr key={inv.id}>
                <td>
                  <Link href={`/my-invoices/${inv.id}`}>{inv.amount.toFixed(2)} €</Link>
                </td>
                <td>{inv.status}</td>
                <td>{inv.createdAt?.slice(0, 10) ?? '—'}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={3}>Ainda não tem faturas registadas.</td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

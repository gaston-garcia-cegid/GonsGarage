'use client';

import React, { Suspense, useCallback, useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { receivedInvoiceService } from '@/lib/services/received-invoice.service';
import type { ReceivedInvoice } from '@/types/accounting';
import styles from '../accounting.module.css';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { ReceivedInvoiceCreateForm } from './ReceivedInvoiceCreateForm';

function ReceivedInvoicesListContent() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [items, setItems] = useState<ReceivedInvoice[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const openedFromCreateQuery = useRef(false);

  const load = useCallback(async () => {
    setError(null);
    const res = await receivedInvoiceService.list();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar as faturas.');
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    if (searchParams.get('create') !== '1' || openedFromCreateQuery.current) return;
    openedFromCreateQuery.current = true;
    setCreateOpen(true);
    router.replace('/accounting/received-invoices');
  }, [searchParams, router]);

  if (!user) return null;

  const handleCreated = () => {
    setCreateOpen(false);
    void load();
  };

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
        <Button type="button" onClick={() => setCreateOpen(true)}>
          Nova fatura
        </Button>
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
                  Sem registos.{' '}
                  <button type="button" className={styles.inlineTextButton} onClick={() => setCreateOpen(true)}>
                    Criar o primeiro
                  </button>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-h-[90vh] overflow-hidden sm:max-w-lg" aria-describedby={undefined}>
          <DialogHeader>
            <DialogTitle>Nova fatura recebida</DialogTitle>
          </DialogHeader>
          <div className={styles.dialogFormBody}>
            <ReceivedInvoiceCreateForm onSuccess={handleCreated} onCancel={() => setCreateOpen(false)} />
          </div>
        </DialogContent>
      </Dialog>
    </AppShell>
  );
}

export default function ReceivedInvoicesListPage() {
  return (
    <Suspense fallback={null}>
      <ReceivedInvoicesListContent />
    </Suspense>
  );
}

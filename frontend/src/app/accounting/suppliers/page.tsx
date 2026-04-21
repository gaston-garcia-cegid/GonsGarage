'use client';

import React, { Suspense, useCallback, useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { supplierService } from '@/lib/services/supplier.service';
import type { Supplier } from '@/types/accounting';
import styles from '../accounting.module.css';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import { SupplierCreateForm } from './SupplierCreateForm';

function SuppliersListContent() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const [items, setItems] = useState<Supplier[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [createOpen, setCreateOpen] = useState(false);
  const openedFromCreateQuery = useRef(false);

  const load = useCallback(async () => {
    setError(null);
    const res = await supplierService.list();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar os fornecedores.');
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  useEffect(() => {
    if (searchParams.get('create') !== '1' || openedFromCreateQuery.current) return;
    openedFromCreateQuery.current = true;
    setCreateOpen(true);
    router.replace('/accounting/suppliers');
  }, [searchParams, router]);

  if (!user) return null;

  const handleCreated = () => {
    setCreateOpen(false);
    void load();
  };

  return (
    <AppShell
      user={user}
      subtitle="Contabilidade — Fornecedores"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <div className={styles.toolbar}>
        <h1>Fornecedores</h1>
        <Button type="button" onClick={() => setCreateOpen(true)}>
          Novo fornecedor
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
              <th>Nome</th>
              <th>Email</th>
              <th>Ativo</th>
            </tr>
          </thead>
          <tbody>
            {items.map((s) => (
              <tr key={s.id}>
                <td>
                  <Link href={`/accounting/suppliers/${s.id}`}>{s.name}</Link>
                </td>
                <td>{s.contactEmail || '—'}</td>
                <td>{s.isActive ? 'Sim' : 'Não'}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={3}>
                  Sem fornecedores.{' '}
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
            <DialogTitle>Novo fornecedor</DialogTitle>
          </DialogHeader>
          <div className={styles.dialogFormBody}>
            <SupplierCreateForm onSuccess={handleCreated} onCancel={() => setCreateOpen(false)} />
          </div>
        </DialogContent>
      </Dialog>
    </AppShell>
  );
}

export default function SuppliersListPage() {
  return (
    <Suspense fallback={null}>
      <SuppliersListContent />
    </Suspense>
  );
}

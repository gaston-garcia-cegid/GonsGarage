'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { supplierService } from '@/lib/services/supplier.service';
import type { Supplier } from '@/types/accounting';
import styles from '../accounting.module.css';

export default function SuppliersListPage() {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<Supplier[]>([]);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    setError(null);
    const res = await supplierService.list();
    if (res.success && res.data) setItems(res.data.items);
    else setError(res.error?.message ?? 'Não foi possível carregar os fornecedores.');
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  if (!user) return null;

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
        <Link href="/accounting/suppliers/new" className={styles.primaryLink}>
          Novo fornecedor
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
                  <Link href="/accounting/suppliers/new">Criar o primeiro</Link>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

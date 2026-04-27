'use client';

import React, { useCallback, useEffect, useState } from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { apiClient } from '@/lib/api-client';
import type { PartItem } from '@/types/parts';
import styles from './admin-parts.module.css';
import { Button } from '@/components/ui/button';

export default function AdminPartsPage() {
  const { user, logout } = useAuth();
  const [items, setItems] = useState<PartItem[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [barcode, setBarcode] = useState('');
  const [searchText, setSearchText] = useState('');

  const load = useCallback(async (opts?: { barcode?: string; search?: string }) => {
    setError(null);
    const res = await apiClient.listParts({
      barcode: opts?.barcode,
      search: opts?.search,
      limit: 200,
      offset: 0,
    });
    if (res.success && res.data) {
      setItems(res.data.items);
    } else {
      setError(res.error?.message ?? 'Não foi possível carregar as peças.');
    }
  }, []);

  useEffect(() => {
    void load();
  }, [load]);

  if (!user) return null;

  function onFilterSubmit(e: React.FormEvent) {
    e.preventDefault();
    const bc = barcode.trim();
    const st = searchText.trim();
    void load({
      barcode: bc || undefined,
      search: st || undefined,
    });
  }

  function onShowAll() {
    setBarcode('');
    setSearchText('');
    void load();
  }

  const toolbar = (
    <>
      <h1>Peças (stock)</h1>
      <Button type="button" asChild>
        <Link href="/admin/parts/new">Nova peça</Link>
      </Button>
    </>
  );

  return (
    <AppShell
      user={user}
      subtitle="Inventário"
      activeNav="admin_parts"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
      toolbar={toolbar}
    >
      <p className={styles.hint}>
        Procure por código de barras (Enter ou botão) ou por texto em referência / marca / nome. Apenas gestores e
        administradores acedem a esta área.
      </p>
      <form className={styles.searchRow} onSubmit={onFilterSubmit}>
        <input
          name="barcode"
          value={barcode}
          onChange={(ev) => setBarcode(ev.target.value)}
          placeholder="Código de barras"
          aria-label="Código de barras"
          autoComplete="off"
        />
        <input
          name="search"
          value={searchText}
          onChange={(ev) => setSearchText(ev.target.value)}
          placeholder="Texto (ref., marca, nome)"
          aria-label="Pesquisa por texto"
          autoComplete="off"
        />
        <Button type="submit">Aplicar</Button>
        <Button type="button" variant="outline" onClick={() => void onShowAll()}>
          Mostrar todas
        </Button>
      </form>
      {error ? <div className={styles.error}>{error}</div> : null}
      <div className={styles.tableWrap}>
        <table className={styles.table}>
          <thead>
            <tr>
              <th>Referência</th>
              <th>Marca</th>
              <th>Nome</th>
              <th>Cód. barras</th>
              <th>Qtd</th>
              <th>UoM</th>
            </tr>
          </thead>
          <tbody>
            {items.map((p) => (
              <tr key={p.id}>
                <td>
                  <Link href={`/admin/parts/${p.id}`}>{p.reference}</Link>
                </td>
                <td>{p.brand}</td>
                <td>{p.name}</td>
                <td>{p.barcode || '—'}</td>
                <td>{p.quantity}</td>
                <td>{p.uom}</td>
              </tr>
            ))}
            {items.length === 0 && !error ? (
              <tr>
                <td colSpan={6}>
                  Sem peças nesta vista.{' '}
                  <Link href="/admin/parts/new" className={styles.mutedLink}>
                    Criar a primeira
                  </Link>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>
    </AppShell>
  );
}

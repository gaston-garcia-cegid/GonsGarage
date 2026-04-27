'use client';

import React, { Suspense, useCallback, useEffect, useRef, useState } from 'react';
import Link from 'next/link';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { apiClient } from '@/lib/api-client';
import type { PartItem } from '@/types/parts';
import styles from './admin-parts.module.css';
import { Button } from '@/components/ui/button';
import { PartCreateModal } from './components/PartCreateModal';

function AdminPartsPageContent() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();
  const openedFromCreateQuery = useRef(false);
  const [createOpen, setCreateOpen] = useState(false);
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
    let cancelled = false;
    queueMicrotask(() => {
      if (!cancelled) void load();
    });
    return () => {
      cancelled = true;
    };
  }, [load]);

  useEffect(() => {
    if (searchParams.get('create') !== '1' || openedFromCreateQuery.current) return;
    openedFromCreateQuery.current = true;
    setCreateOpen(true);
    router.replace('/admin/parts', { scroll: false });
  }, [searchParams, router]);

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

  const handlePartCreated = (partId: string) => {
    setCreateOpen(false);
    void load();
    router.push(`/admin/parts/${partId}`);
  };

  const toolbar = (
    <>
      <h1>Peças (stock)</h1>
      <Button type="button" onClick={() => setCreateOpen(true)}>
        Nova peça
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
          onChange={ev => setBarcode(ev.target.value)}
          placeholder="Código de barras"
          aria-label="Código de barras"
          autoComplete="off"
        />
        <input
          name="search"
          value={searchText}
          onChange={ev => setSearchText(ev.target.value)}
          placeholder="Texto (ref., marca, nome)"
          aria-label="Pesquisa por texto"
          autoComplete="off"
        />
        <Button type="submit">Aplicar</Button>
        <Button type="button" variant="outline" onClick={() => onShowAll()}>
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
            {items.map(p => (
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
                  <button type="button" className={styles.inlineTextButton} onClick={() => setCreateOpen(true)}>
                    Criar a primeira
                  </button>
                </td>
              </tr>
            ) : null}
          </tbody>
        </table>
      </div>

      <PartCreateModal open={createOpen} onOpenChange={setCreateOpen} onSuccess={handlePartCreated} />
    </AppShell>
  );
}

export default function AdminPartsPage() {
  return (
    <Suspense fallback={null}>
      <AdminPartsPageContent />
    </Suspense>
  );
}

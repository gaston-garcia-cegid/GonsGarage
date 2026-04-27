'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { apiClient } from '@/lib/api-client';
import type { PartItemWriteBody, PartUOM } from '@/types/parts';
import styles from '../admin-parts.module.css';
import { Button } from '@/components/ui/button';

const UOM_OPTIONS: { value: PartUOM; label: string }[] = [
  { value: 'unit', label: 'Unidade (unit)' },
  { value: 'liter', label: 'Litro (liter)' },
];

export default function AdminPartsNewPage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [reference, setReference] = useState('');
  const [brand, setBrand] = useState('');
  const [name, setName] = useState('');
  const [barcode, setBarcode] = useState('');
  const [quantity, setQuantity] = useState('0');
  const [uom, setUom] = useState<PartUOM>('unit');
  const [minimumQuantity, setMinimumQuantity] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  if (!user) return null;

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError(null);
    const qty = Number(quantity);
    const body: PartItemWriteBody = {
      reference: reference.trim(),
      brand: brand.trim(),
      name: name.trim(),
      barcode: barcode.trim(),
      quantity: Number.isFinite(qty) ? qty : 0,
      uom,
    };
    const min = minimumQuantity.trim();
    if (min !== '') {
      const m = parseFloat(min);
      body.minimumQuantity = Number.isFinite(m) ? m : null;
    }
    const res = await apiClient.createPart(body);
    setSaving(false);
    if (res.success && res.data) {
      router.replace(`/admin/parts/${res.data.id}`);
      return;
    }
    setError(res.error?.message ?? 'Não foi possível criar a peça.');
  }

  const toolbar = (
    <>
      <h1>Nova peça</h1>
      <Button type="button" variant="outline" asChild>
        <Link href="/admin/parts">Voltar à lista</Link>
      </Button>
    </>
  );

  return (
    <AppShell
      user={user}
      subtitle="Inventário — Nova peça"
      activeNav="admin_parts"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
      toolbar={toolbar}
    >
      <p className={styles.intro}>
        <Link href="/admin/parts" className={styles.mutedLink}>
          ← Peças (stock)
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={(ev) => void onSubmit(ev)}>
        <div className={styles.field}>
          <label htmlFor="reference">Referência</label>
          <input id="reference" value={reference} onChange={(ev) => setReference(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="brand">Marca</label>
          <input id="brand" value={brand} onChange={(ev) => setBrand(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="name">Nome</label>
          <input id="name" value={name} onChange={(ev) => setName(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="barcode">Código de barras</label>
          <input id="barcode" value={barcode} onChange={(ev) => setBarcode(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="quantity">Quantidade</label>
          <input
            id="quantity"
            type="number"
            min={0}
            step="any"
            value={quantity}
            onChange={(ev) => setQuantity(ev.target.value)}
            required
          />
        </div>
        <div className={styles.field}>
          <label htmlFor="uom">Unidade de medida</label>
          <select id="uom" value={uom} onChange={(ev) => setUom(ev.target.value as PartUOM)} required>
            {UOM_OPTIONS.map((o) => (
              <option key={o.value} value={o.value}>
                {o.label}
              </option>
            ))}
          </select>
        </div>
        <div className={styles.field}>
          <label htmlFor="minimumQuantity">Quantidade mínima (opcional)</label>
          <input
            id="minimumQuantity"
            type="number"
            min={0}
            step="any"
            value={minimumQuantity}
            onChange={(ev) => setMinimumQuantity(ev.target.value)}
          />
        </div>
        <div className={styles.rowActions}>
          <Button type="submit" disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </Button>
        </div>
      </form>
    </AppShell>
  );
}

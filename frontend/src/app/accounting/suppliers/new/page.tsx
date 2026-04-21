'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import { supplierService } from '@/lib/services/supplier.service';
import styles from '../../accounting.module.css';
import { Button } from '@/components/ui/button';

export default function NewSupplierPage() {
  const { user, logout } = useAuth();
  const router = useRouter();
  const [name, setName] = useState('');
  const [contactEmail, setContactEmail] = useState('');
  const [contactPhone, setContactPhone] = useState('');
  const [taxId, setTaxId] = useState('');
  const [notes, setNotes] = useState('');
  const [isActive, setIsActive] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  if (!user) return null;

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setSaving(true);
    setError(null);
    const res = await supplierService.create({
      name: name.trim(),
      contactEmail: contactEmail.trim(),
      contactPhone: contactPhone.trim(),
      taxId: taxId.trim(),
      notes: notes.trim(),
      isActive,
    });
    setSaving(false);
    if (res.success && res.data?.id) {
      router.replace(`/accounting/suppliers/${res.data.id}`);
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <AppShell
      user={user}
      subtitle="Novo fornecedor"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        <Link href="/accounting/suppliers" className={styles.mutedLink}>
          ← Fornecedores
        </Link>
      </p>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="name">Nome</label>
          <input id="name" value={name} onChange={(ev) => setName(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="email">Email de contacto</label>
          <input id="email" type="email" value={contactEmail} onChange={(ev) => setContactEmail(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="phone">Telefone</label>
          <input id="phone" value={contactPhone} onChange={(ev) => setContactPhone(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="tax">NIF / Tax ID</label>
          <input id="tax" value={taxId} onChange={(ev) => setTaxId(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="notes">Notas</label>
          <textarea id="notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label>
            <input type="checkbox" checked={isActive} onChange={(ev) => setIsActive(ev.target.checked)} /> Ativo
          </label>
        </div>
        <div className={styles.rowActions}>
          <Button type="submit" className={styles.submitButton} disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </Button>
        </div>
      </form>
    </AppShell>
  );
}

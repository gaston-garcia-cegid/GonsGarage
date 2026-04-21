'use client';

import React, { useState } from 'react';
import { supplierService } from '@/lib/services/supplier.service';
import styles from '../accounting.module.css';
import { Button } from '@/components/ui/button';

export interface SupplierCreateFormProps {
  onSuccess: () => void;
  onCancel: () => void;
}

export function SupplierCreateForm({ onSuccess, onCancel }: Readonly<SupplierCreateFormProps>) {
  const [name, setName] = useState('');
  const [contactEmail, setContactEmail] = useState('');
  const [contactPhone, setContactPhone] = useState('');
  const [taxId, setTaxId] = useState('');
  const [notes, setNotes] = useState('');
  const [isActive, setIsActive] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

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
      onSuccess();
      return;
    }
    setError(res.error?.message ?? 'Erro ao criar.');
  }

  return (
    <div>
      {error ? <div className={styles.error}>{error}</div> : null}
      <form className={styles.form} onSubmit={onSubmit}>
        <div className={styles.field}>
          <label htmlFor="supplier-name">Nome</label>
          <input id="supplier-name" value={name} onChange={(ev) => setName(ev.target.value)} required />
        </div>
        <div className={styles.field}>
          <label htmlFor="supplier-email">Email de contacto</label>
          <input id="supplier-email" type="email" value={contactEmail} onChange={(ev) => setContactEmail(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="supplier-phone">Telefone</label>
          <input id="supplier-phone" value={contactPhone} onChange={(ev) => setContactPhone(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="supplier-tax">NIF / Tax ID</label>
          <input id="supplier-tax" value={taxId} onChange={(ev) => setTaxId(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label htmlFor="supplier-notes">Notas</label>
          <textarea id="supplier-notes" value={notes} onChange={(ev) => setNotes(ev.target.value)} />
        </div>
        <div className={styles.field}>
          <label>
            <input type="checkbox" checked={isActive} onChange={(ev) => setIsActive(ev.target.checked)} /> Ativo
          </label>
        </div>
        <div className={styles.rowActions}>
          <Button type="button" variant="outline" onClick={onCancel} disabled={saving}>
            Cancelar
          </Button>
          <Button type="submit" className={styles.submitButton} disabled={saving}>
            {saving ? 'A guardar…' : 'Criar'}
          </Button>
        </div>
      </form>
    </div>
  );
}

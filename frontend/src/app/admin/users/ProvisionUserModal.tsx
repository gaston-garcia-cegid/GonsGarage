'use client';

import React, { useEffect, useState } from 'react';
import { apiClient } from '@/lib/api-client';
import { UserRole } from '@/types';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog';
import styles from './admin-users.module.css';

export type ProvisionRole = 'manager' | 'employee' | 'client';

export interface RoleOption {
  value: ProvisionRole;
  label: string;
}

export interface ProvisionUserModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  roleOptions: RoleOption[];
  defaultRole: ProvisionRole;
  callerRole: UserRole;
  onProvisioned?: (info: { email: string; role: string }) => void;
}

export default function ProvisionUserModal({
  open,
  onOpenChange,
  roleOptions,
  defaultRole,
  callerRole,
  onProvisioned,
}: Readonly<ProvisionUserModalProps>) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [role, setRole] = useState<ProvisionRole>(defaultRole);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!open) return;
    let cancelled = false;
    queueMicrotask(() => {
      if (cancelled) return;
      setEmail('');
      setPassword('');
      setFirstName('');
      setLastName('');
      setRole(defaultRole);
      setError(null);
    });
    return () => {
      cancelled = true;
    };
  }, [open, defaultRole]);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    setSubmitting(true);
    try {
      const res = await apiClient.provisionUser({
        email: email.trim(),
        password,
        firstName: firstName.trim(),
        lastName: lastName.trim(),
        role,
      });
      if (!res.success || !res.data?.user) {
        setError(res.error?.message ?? 'Pedido falhou');
        return;
      }
      onProvisioned?.({ email: res.data.user.email, role: res.data.user.role });
      onOpenChange(false);
      setRole(callerRole === UserRole.ADMIN ? 'client' : 'employee');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] overflow-hidden sm:max-w-lg" aria-describedby={undefined}>
        <DialogHeader>
          <DialogTitle>Novo utilizador</DialogTitle>
        </DialogHeader>
        <form className={styles.modalForm} onSubmit={onSubmit}>
          <p className={styles.modalHint}>
            O servidor recusa combinações inválidas (por exemplo, gestor a criar outro gestor).
          </p>
          <label className={styles.field}>
            <span>E-mail</span>
            <Input type="email" value={email} onChange={(ev) => setEmail(ev.target.value)} required autoComplete="off" />
          </label>
          <label className={styles.field}>
            <span>Palavra-passe inicial</span>
            <Input
              type="password"
              value={password}
              onChange={(ev) => setPassword(ev.target.value)}
              required
              minLength={6}
              autoComplete="new-password"
            />
          </label>
          <label className={styles.field}>
            <span>Nome</span>
            <Input value={firstName} onChange={(ev) => setFirstName(ev.target.value)} required />
          </label>
          <label className={styles.field}>
            <span>Apelido</span>
            <Input value={lastName} onChange={(ev) => setLastName(ev.target.value)} required />
          </label>
          <label className={styles.field}>
            <span>Papel</span>
            <select
              className={styles.select}
              value={role}
              onChange={(ev) => setRole(ev.target.value as ProvisionRole)}
              required
            >
              {roleOptions.map((o) => (
                <option key={o.value} value={o.value}>
                  {o.label}
                </option>
              ))}
            </select>
          </label>
          {error ? <p className={styles.error}>{error}</p> : null}
          <div className={styles.modalActions}>
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              Cancelar
            </Button>
            <Button type="submit" disabled={submitting}>
              {submitting ? 'A criar…' : 'Criar utilizador'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}

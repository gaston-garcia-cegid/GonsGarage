'use client';

import React, { useMemo, useState } from 'react';
import { useAuth } from '@/stores';
import { UserRole } from '@/types';
import { canManageUsers } from '@/types/user';
import { apiClient } from '@/lib/api-client';
import AppShell from '@/components/layout/AppShell';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import styles from './admin-users.module.css';

type ProvisionRole = 'manager' | 'employee' | 'client';

export default function AdminUsersPage() {
  const { user, logout } = useAuth();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [role, setRole] = useState<ProvisionRole>('client');
  const [submitting, setSubmitting] = useState(false);
  const [message, setMessage] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const roleOptions = useMemo(() => {
    if (!user || !canManageUsers(user)) return [] as { value: ProvisionRole; label: string }[];
    if (user.role === UserRole.ADMIN) {
      return [
        { value: 'manager' as const, label: 'Gestor (manager)' },
        { value: 'employee' as const, label: 'Funcionário (employee)' },
        { value: 'client' as const, label: 'Cliente (client)' },
      ];
    }
    return [
      { value: 'employee' as const, label: 'Funcionário (employee)' },
      { value: 'client' as const, label: 'Cliente (client)' },
    ];
  }, [user]);

  if (!user) return null;

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!user) return;
    setMessage(null);
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
      setMessage(`Utilizador criado: ${res.data.user.email} (${res.data.user.role})`);
      setEmail('');
      setPassword('');
      setFirstName('');
      setLastName('');
      setRole(user.role === UserRole.ADMIN ? 'client' : 'employee');
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <AppShell
      user={user}
      subtitle="Utilizadores"
      activeNav="admin_users"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <div className={styles.wrap}>
        <h2 className={styles.title}>Criar utilizador</h2>
        <p className={styles.hint}>
          Apenas administradores e gestores podem usar este formulário. O servidor recusa combinações inválidas (por
          exemplo, gestor a criar outro gestor).
        </p>
        <form className={styles.form} onSubmit={onSubmit}>
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
          {message ? <p className={styles.ok}>{message}</p> : null}
          <Button type="submit" disabled={submitting}>
            {submitting ? 'A criar…' : 'Criar utilizador'}
          </Button>
        </form>
      </div>
    </AppShell>
  );
}

'use client';

import React, { useMemo, useState } from 'react';
import { useAuth } from '@/stores';
import { UserRole } from '@/types';
import { canManageUsers } from '@/types/user';
import AppShell from '@/components/layout/AppShell';
import { Button } from '@/components/ui/button';
import styles from './admin-users.module.css';
import ProvisionUserModal, { type ProvisionRole, type RoleOption } from './ProvisionUserModal';

export default function AdminUsersPage() {
  const { user, logout } = useAuth();
  const [createOpen, setCreateOpen] = useState(false);
  const [message, setMessage] = useState<string | null>(null);

  const roleOptions = useMemo((): RoleOption[] => {
    if (!user || !canManageUsers(user)) return [];
    if (user.role === UserRole.ADMIN) {
      return [
        { value: 'manager', label: 'Gestor (manager)' },
        { value: 'employee', label: 'Funcionário (employee)' },
        { value: 'client', label: 'Cliente (client)' },
      ];
    }
    return [
      { value: 'employee', label: 'Funcionário (employee)' },
      { value: 'client', label: 'Cliente (client)' },
    ];
  }, [user]);

  const defaultRole: ProvisionRole = user?.role === UserRole.ADMIN ? 'client' : 'employee';

  if (!user) return null;

  const toolbar = (
    <>
      <h1>Utilizadores</h1>
      <Button type="button" onClick={() => setCreateOpen(true)}>
        Novo utilizador
      </Button>
    </>
  );

  return (
    <AppShell
      user={user}
      subtitle="Utilizadores"
      activeNav="admin_users"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
      toolbar={toolbar}
    >
      <p className={styles.hint}>
        Apenas administradores e gestores podem criar utilizadores. Utilize o botão acima para abrir o formulário.
      </p>
      {message ? <p className={styles.ok}>{message}</p> : null}

      <ProvisionUserModal
        open={createOpen}
        onOpenChange={setCreateOpen}
        roleOptions={roleOptions}
        defaultRole={defaultRole}
        callerRole={user.role}
        onProvisioned={({ email, role }) => {
          setMessage(`Utilizador criado: ${email} (${role})`);
        }}
      />
    </AppShell>
  );
}

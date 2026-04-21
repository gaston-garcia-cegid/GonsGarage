'use client';

import React from 'react';
import { Button } from '@/components/ui/button';
import styles from './dashboard.module.css';

export interface StaffUsersDashboardCtaProps {
  onManageUsers: () => void;
}

/** Secondary entry to `/admin/users` for admin/manager (parent applies `canManageUsers`). */
export function StaffUsersDashboardCta({ onManageUsers }: Readonly<StaffUsersDashboardCtaProps>) {
  return (
    <div className={styles.staffUsersCta}>
      <p className={styles.staffUsersCtaText}>Criar ou atualizar contas de staff e clientes.</p>
      <Button type="button" variant="outline" onClick={onManageUsers}>
        Gestão de utilizadores
      </Button>
    </div>
  );
}

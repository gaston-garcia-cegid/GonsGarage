'use client';

import React from 'react';
import Link from 'next/link';
import { useAuth } from '@/stores';
import AppShell from '@/components/layout/AppShell';
import styles from './accounting.module.css';

export default function AccountingHubPage() {
  const { user, logout } = useAuth();
  if (!user) return null;

  return (
    <AppShell
      user={user}
      subtitle="Contabilidade"
      activeNav="accounting"
      carsNavLabel="Viaturas"
      onLogout={logout}
      logoVariant="branded"
    >
      <p className={styles.intro}>
        Gestão de fornecedores, faturas recebidas, documentos de faturação e faturas emitidas a clientes.
      </p>
      <div className={styles.cardGrid}>
        <Link href="/accounting/suppliers" className={styles.card}>
          Fornecedores
          <span>CRUD de fornecedores</span>
        </Link>
        <Link href="/accounting/received-invoices" className={styles.card}>
          Faturas recebidas
          <span>Despesas e compras</span>
        </Link>
        <Link href="/accounting/billing-documents" className={styles.card}>
          Documentos de faturação
          <span>IRS, salários, outros</span>
        </Link>
        <Link href="/accounting/issued-invoices" className={styles.card}>
          Faturas emitidas
          <span>Faturas a clientes</span>
        </Link>
      </div>
    </AppShell>
  );
}

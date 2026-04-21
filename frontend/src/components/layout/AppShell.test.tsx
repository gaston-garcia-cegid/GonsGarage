import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AppShell from './AppShell';
import { UserRole } from '@/types';
import type { User } from '@/types';

const mockPush = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: mockPush,
  }),
}));

function buildUser(overrides: Partial<User> & Pick<User, 'role'>): User {
  return {
    id: '11111111-1111-1111-1111-111111111111',
    email: 'u@test.com',
    firstName: 'Test',
    lastName: 'User',
    createdAt: '2020-01-01T00:00:00.000Z',
    updatedAt: '2020-01-01T00:00:00.000Z',
    ...overrides,
  };
}

describe('AppShell staff user management nav', () => {
  beforeEach(() => {
    mockPush.mockClear();
  });

  it('shows Utilizadores and navigates to /admin/users for manager', async () => {
    const user = userEvent.setup();
    const manager = buildUser({ role: UserRole.MANAGER });
    render(
      <AppShell
        user={manager}
        subtitle="Teste"
        activeNav="dashboard"
        onLogout={vi.fn()}
      >
        <p>Conteúdo</p>
      </AppShell>
    );

    const navUsers = screen.getByRole('button', { name: 'Utilizadores' });
    expect(navUsers).toBeInTheDocument();
    await user.click(navUsers);
    expect(mockPush).toHaveBeenCalledWith('/admin/users');
  });

  it('shows Utilizadores for admin', () => {
    const admin = buildUser({ role: UserRole.ADMIN });
    render(
      <AppShell user={admin} subtitle="Teste" activeNav="dashboard" onLogout={vi.fn()}>
        <p>X</p>
      </AppShell>
    );
    expect(screen.getByRole('button', { name: 'Utilizadores' })).toBeInTheDocument();
  });

  it('does not show Utilizadores for client', () => {
    const client = buildUser({ role: UserRole.CLIENT });
    render(
      <AppShell user={client} subtitle="Teste" activeNav="dashboard" onLogout={vi.fn()}>
        <p>X</p>
      </AppShell>
    );
    expect(screen.queryByRole('button', { name: 'Utilizadores' })).not.toBeInTheDocument();
  });

  it('does not show Utilizadores for employee', () => {
    const employee = buildUser({ role: UserRole.EMPLOYEE });
    render(
      <AppShell user={employee} subtitle="Teste" activeNav="accounting" onLogout={vi.fn()}>
        <p>X</p>
      </AppShell>
    );
    expect(screen.queryByRole('button', { name: 'Utilizadores' })).not.toBeInTheDocument();
  });

  it('marks Utilizadores active when activeNav is admin_users', () => {
    const manager = buildUser({ role: UserRole.MANAGER });
    render(
      <AppShell user={manager} subtitle="Utilizadores" activeNav="admin_users" onLogout={vi.fn()}>
        <p>Form</p>
      </AppShell>
    );
    const btn = screen.getByRole('button', { name: 'Utilizadores' });
    expect(btn.className).toMatch(/active/);
  });
});

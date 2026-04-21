import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DashboardLayout from './DashboardLayout';
import { UserRole } from '@/types';

vi.mock('next/image', () => ({
  default: (props: { src: string; alt: string }) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={props.src} alt={props.alt} data-testid="layout-logo" />
  ),
}));

const push = vi.fn();
vi.mock('next/navigation', () => ({
  useRouter: () => ({ push }),
}));

const logout = vi.fn();
vi.mock('@/stores', () => ({
  useAuth: () => ({
    user: {
      id: 'u-1',
      email: 'client@test.com',
      firstName: 'Maria',
      lastName: 'Costa',
      role: UserRole.CLIENT,
      createdAt: '2026-01-01T00:00:00.000Z',
      updatedAt: '2026-01-01T00:00:00.000Z',
    },
    logout,
  }),
}));

describe('DashboardLayout', () => {
  const navigationItems = [
    { key: 'dashboard', label: 'Painel', href: '#' },
    { key: 'cars', label: 'Os meus automóveis', href: '#' },
  ];

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('greets the signed-in user from the store', () => {
    render(
      <DashboardLayout
        title="Painel"
        subtitle="Área do cliente"
        activeTab="dashboard"
        navigationItems={navigationItems}
      >
        <p>Conteúdo</p>
      </DashboardLayout>,
    );

    expect(screen.getByText(/Welcome, Maria Costa/)).toBeInTheDocument();
  });

  it('invokes logout when the user clicks Logout', async () => {
    const user = userEvent.setup();
    render(
      <DashboardLayout
        title="Painel"
        subtitle="Área do cliente"
        activeTab="dashboard"
        navigationItems={navigationItems}
      >
        <p>Conteúdo</p>
      </DashboardLayout>,
    );

    await user.click(screen.getByRole('button', { name: 'Logout' }));
    expect(logout).toHaveBeenCalledTimes(1);
  });

  it('uses onNavClick for SPA tabs when provided', async () => {
    const user = userEvent.setup();
    const onNavClick = vi.fn();

    render(
      <DashboardLayout
        title="Automóveis"
        subtitle="Área do cliente"
        activeTab="cars"
        navigationItems={navigationItems}
        onNavClick={onNavClick}
      >
        <p>Lista</p>
      </DashboardLayout>,
    );

    await user.click(screen.getByRole('button', { name: 'Painel' }));
    expect(onNavClick).toHaveBeenCalledWith('dashboard');
    expect(push).not.toHaveBeenCalled();
  });
});

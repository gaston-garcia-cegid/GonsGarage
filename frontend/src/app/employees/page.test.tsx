import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import EmployeesPage from './page';

const mockLogout = vi.fn();
const mockGetEmployees = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn() }),
}));

vi.mock('@/contexts/AuthContext', () => ({
  useAuth: () => ({
    user: {
      id: 'emp-admin',
      email: 'admin@test.com',
      first_name: 'Pat',
      last_name: 'Lee',
      role: 'admin',
      is_active: true,
      created_at: '',
      updated_at: '',
    },
    logout: mockLogout,
  }),
}));

vi.mock('@/lib/api', async (importOriginal) => {
  const mod = await importOriginal<typeof import('@/lib/api')>();
  return {
    ...mod,
    apiClient: {
      ...mod.apiClient,
      getEmployees: (...args: unknown[]) => mockGetEmployees(...args),
    },
  };
});

describe('EmployeesPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockGetEmployees.mockResolvedValue({
      data: { employees: [], total: 0 },
      error: null,
    });
  });

  it('opens the create-employee modal from the toolbar CTA', async () => {
    const user = userEvent.setup();
    render(<EmployeesPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Adicionar colaborador/i })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: /Adicionar colaborador/i }));

    expect(screen.getByRole('heading', { name: 'Novo colaborador' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancelar' })).toBeInTheDocument();
  });

  it('terminates the session when the user confirms logout', async () => {
    const user = userEvent.setup();
    render(<EmployeesPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: /Terminar sessão/i })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: /Terminar sessão/i }));
    expect(mockLogout).toHaveBeenCalledTimes(1);
  });
});

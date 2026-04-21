import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import SuppliersListPage from './page';
import { UserRole } from '@/types';

const mockReplace = vi.fn();
let searchParamsString = '';

const { listMock, createMock } = vi.hoisted(() => ({
  listMock: vi.fn(),
  createMock: vi.fn(),
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    replace: mockReplace,
    push: vi.fn(),
  }),
  useSearchParams: () => new URLSearchParams(searchParamsString),
}));

vi.mock('@/stores', () => ({
  useAuth: () => ({
    user: {
      id: '11111111-1111-1111-1111-111111111111',
      email: 'mgr@test.com',
      firstName: 'Maria',
      lastName: 'Gestora',
      role: UserRole.MANAGER,
      createdAt: '2020-01-01T00:00:00.000Z',
      updatedAt: '2020-01-01T00:00:00.000Z',
    },
    logout: vi.fn(),
  }),
}));

vi.mock('@/lib/services/supplier.service', () => ({
  supplierService: {
    list: (...args: unknown[]) => listMock(...args),
    create: (...args: unknown[]) => createMock(...args),
  },
}));

const emptyList = { success: true, data: { items: [], total: 0 } };

const supplierRow = {
  id: 'new-sup-1',
  name: 'Acme Lda',
  contactEmail: '',
  contactPhone: '',
  taxId: '',
  notes: '',
  isActive: true,
  createdAt: '2020-01-01T00:00:00.000Z',
  updatedAt: '2020-01-01T00:00:00.000Z',
};

describe('SuppliersListPage create modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    searchParamsString = '';
    listMock.mockResolvedValue(emptyList);
    createMock.mockResolvedValue({ success: true, data: supplierRow });
  });

  it('opens the create dialog from the toolbar button', async () => {
    const user = userEvent.setup();
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Novo fornecedor' })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Novo fornecedor' }));

    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Novo fornecedor' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Cancelar' })).toBeInTheDocument();
  });

  it('opens the create dialog when URL has create=1 and clears the query', async () => {
    searchParamsString = 'create=1';
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Novo fornecedor' })).toBeInTheDocument();
    });

    expect(mockReplace).toHaveBeenCalledWith('/accounting/suppliers');
  });

  it('opens the create dialog from the empty-state CTA', async () => {
    const user = userEvent.setup();
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Criar o primeiro' })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Criar o primeiro' }));

    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Novo fornecedor' })).toBeInTheDocument();
  });

  it('closes the dialog on Cancel without calling create', async () => {
    const user = userEvent.setup();
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Novo fornecedor' })).toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Novo fornecedor' }));
    await user.click(screen.getByRole('button', { name: 'Cancelar' }));

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
    expect(createMock).not.toHaveBeenCalled();
  });

  it('closes the dialog on Escape without calling create', async () => {
    const user = userEvent.setup();
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Novo fornecedor' })).toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Novo fornecedor' }));
    await user.keyboard('{Escape}');

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
    expect(createMock).not.toHaveBeenCalled();
  });

  it('submits create, closes the dialog, and reloads the list', async () => {
    const user = userEvent.setup();
    render(<SuppliersListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Novo fornecedor' })).toBeInTheDocument();
    });
    expect(listMock).toHaveBeenCalledTimes(1);

    await user.click(screen.getByRole('button', { name: 'Novo fornecedor' }));
    await user.type(screen.getByLabelText('Nome'), 'Acme Lda');
    await user.click(screen.getByRole('button', { name: 'Criar' }));

    await waitFor(() => {
      expect(createMock).toHaveBeenCalledTimes(1);
    });
    await waitFor(() => {
      expect(listMock).toHaveBeenCalledTimes(2);
    });
    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
  });
});

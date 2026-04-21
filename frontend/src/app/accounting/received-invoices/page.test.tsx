import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ReceivedInvoicesListPage from './page';
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

vi.mock('@/lib/services/received-invoice.service', () => ({
  receivedInvoiceService: {
    list: (...args: unknown[]) => listMock(...args),
    create: (...args: unknown[]) => createMock(...args),
  },
}));

const emptyList = { success: true, data: { items: [], total: 0 } };

const invoiceRow = {
  id: 'ri-new-1',
  vendorName: 'Forn A',
  category: 'Mat',
  amount: 42.5,
  invoiceDate: '2024-03-15',
  notes: '',
  createdAt: '2020-01-01T00:00:00.000Z',
  updatedAt: '2020-01-01T00:00:00.000Z',
};

describe('ReceivedInvoicesListPage create modal', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    searchParamsString = '';
    listMock.mockResolvedValue(emptyList);
    createMock.mockResolvedValue({ success: true, data: invoiceRow });
  });

  it('opens the create dialog from the toolbar button', async () => {
    const user = userEvent.setup();
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Nova fatura' })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Nova fatura' }));

    expect(screen.getByRole('dialog')).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Nova fatura recebida' })).toBeInTheDocument();
  });

  it('opens the create dialog when URL has create=1', async () => {
    searchParamsString = 'create=1';
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('heading', { name: 'Nova fatura recebida' })).toBeInTheDocument();
    });

    expect(mockReplace).toHaveBeenCalledWith('/accounting/received-invoices');
  });

  it('opens the create dialog from the empty-state CTA', async () => {
    const user = userEvent.setup();
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Criar o primeiro' })).toBeInTheDocument();
    });

    await user.click(screen.getByRole('button', { name: 'Criar o primeiro' }));

    expect(screen.getByRole('dialog')).toBeInTheDocument();
  });

  it('closes the dialog on Cancel without calling create', async () => {
    const user = userEvent.setup();
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Nova fatura' })).toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Nova fatura' }));
    await user.click(screen.getByRole('button', { name: 'Cancelar' }));

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
    expect(createMock).not.toHaveBeenCalled();
  });

  it('closes the dialog on Escape without calling create', async () => {
    const user = userEvent.setup();
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Nova fatura' })).toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Nova fatura' }));
    await user.keyboard('{Escape}');

    await waitFor(() => {
      expect(screen.queryByRole('dialog')).not.toBeInTheDocument();
    });
    expect(createMock).not.toHaveBeenCalled();
  });

  it('submits create, closes the dialog, and reloads the list', async () => {
    const user = userEvent.setup();
    render(<ReceivedInvoicesListPage />);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Nova fatura' })).toBeInTheDocument();
    });
    expect(listMock).toHaveBeenCalledTimes(1);

    await user.click(screen.getByRole('button', { name: 'Nova fatura' }));
    await user.type(screen.getByLabelText('Nome do fornecedor'), 'Forn A');
    await user.type(screen.getByLabelText('Categoria'), 'Mat');
    await user.type(screen.getByLabelText('Valor'), '42,5');
    await user.type(screen.getByLabelText('Data da fatura'), '2024-03-15');
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

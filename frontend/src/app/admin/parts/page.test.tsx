import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import AdminPartsPage from './page';
import { UserRole } from '@/types';

const listPartsMock = vi.fn();
const mockReplace = vi.fn();
const { searchParamsRef } = vi.hoisted(() => ({
  searchParamsRef: { current: '' },
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: vi.fn(),
    replace: mockReplace,
  }),
  useSearchParams: () => new URLSearchParams(searchParamsRef.current),
}));

vi.mock('@/stores', () => ({
  useAuth: () => ({
    user: {
      id: '22222222-2222-2222-2222-222222222222',
      email: 'mgr-parts@test.com',
      firstName: 'Mário',
      lastName: 'Gestor',
      role: UserRole.MANAGER,
      createdAt: '2020-01-01T00:00:00.000Z',
      updatedAt: '2020-01-01T00:00:00.000Z',
    },
    logout: vi.fn(),
  }),
}));

vi.mock('@/lib/api-client', () => ({
  apiClient: {
    listParts: (...args: unknown[]) => listPartsMock(...args),
  },
}));

const emptyList = { success: true as const, data: { items: [], total: 0 } };

describe('AdminPartsPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    searchParamsRef.current = '';
    listPartsMock.mockResolvedValue(emptyList);
  });

  it('loads list and shows empty state without error', async () => {
    render(<AdminPartsPage />);
    const main = screen.getByRole('main');

    await waitFor(() => {
      expect(listPartsMock).toHaveBeenCalledWith(
        expect.objectContaining({ limit: 200, offset: 0 })
      );
    });

    expect(within(main).getByRole('heading', { name: 'Peças (stock)', level: 1 })).toBeInTheDocument();
    expect(await within(main).findByText(/Sem peças nesta vista/i)).toBeInTheDocument();
    expect(within(main).getByRole('button', { name: /Criar a primeira/i })).toBeInTheDocument();
  });

  it('renders rows when API returns items', async () => {
    listPartsMock.mockResolvedValue({
      success: true,
      data: {
        items: [
          {
            id: 'p1111111-1111-1111-1111-111111111111',
            reference: 'REF-1',
            brand: 'ACME',
            name: 'Filtro',
            barcode: '5900000000001',
            quantity: 3,
            uom: 'unit',
            createdAt: '2020-01-01T00:00:00.000Z',
            updatedAt: '2020-01-01T00:00:00.000Z',
          },
        ],
        total: 1,
      },
    });

    render(<AdminPartsPage />);
    const main = screen.getByRole('main');

    expect(await within(main).findByRole('link', { name: 'REF-1' })).toBeInTheDocument();
    expect(within(main).getByRole('cell', { name: 'ACME' })).toBeInTheDocument();
    expect(within(main).getByRole('cell', { name: 'Filtro' })).toBeInTheDocument();
  });

  it('submits filters with barcode and search params', async () => {
    const user = userEvent.setup();
    render(<AdminPartsPage />);

    await waitFor(() => expect(listPartsMock).toHaveBeenCalledTimes(1));

    await user.type(screen.getByLabelText('Código de barras'), '5900');
    await user.type(screen.getByLabelText('Pesquisa por texto'), 'filtro');
    await user.click(screen.getByRole('button', { name: 'Aplicar' }));

    await waitFor(() => {
      expect(listPartsMock).toHaveBeenLastCalledWith(
        expect.objectContaining({
          barcode: '5900',
          search: 'filtro',
          limit: 200,
          offset: 0,
        })
      );
    });
  });

  it('opens create dialog when URL has create=1 and clears the query (modal parity)', async () => {
    searchParamsRef.current = 'create=1';
    render(<AdminPartsPage />);

    await waitFor(() => {
      expect(screen.getByRole('dialog')).toBeInTheDocument();
    });
    expect(screen.getByRole('heading', { name: /Nova peça/i })).toBeInTheDocument();
    expect(mockReplace).toHaveBeenCalledWith('/admin/parts', expect.objectContaining({ scroll: false }));
  });

  it('does not open create dialog when create query is absent', async () => {
    searchParamsRef.current = '';
    render(<AdminPartsPage />);
    const main = screen.getByRole('main');
    await waitFor(() => expect(listPartsMock).toHaveBeenCalled());
    expect(within(main).queryByRole('dialog')).not.toBeInTheDocument();
  });
});

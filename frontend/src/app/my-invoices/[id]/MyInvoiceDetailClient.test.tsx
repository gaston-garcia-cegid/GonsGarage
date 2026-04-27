import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import MyInvoiceDetailClient from './MyInvoiceDetailClient';
import { issuedInvoiceService } from '@/lib/services/issued-invoice.service';
import { UserRole } from '@/types';
import type { IssuedInvoice } from '@/types/accounting';

vi.mock('@/stores', () => ({
  useAuth: () => ({
    user: {
      id: 'u-1',
      email: 'client@test.com',
      firstName: 'C',
      lastName: 'L',
      role: UserRole.CLIENT,
      createdAt: '2026-01-01T00:00:00.000Z',
      updatedAt: '2026-01-01T00:00:00.000Z',
    },
    logout: vi.fn(),
  }),
}));

vi.mock('@/components/layout/AppShell', () => ({
  default: ({ children }: { children: React.ReactNode }) => <div data-testid="shell">{children}</div>,
}));

const baseRow: IssuedInvoice = {
  id: 'inv-1',
  customerId: 'c1',
  amount: 42,
  status: 'issued',
  notes: 'versión servidor',
  createdAt: '2026-01-02T00:00:00.000Z',
  updatedAt: '2026-01-02T00:00:00.000Z',
};

describe('MyInvoiceDetailClient — useOptimistic notes save', () => {
  beforeEach(() => {
    vi.spyOn(issuedInvoiceService, 'get').mockResolvedValue({
      success: true,
      data: baseRow,
    });
  });

  it('shows optimistic server notes snapshot while saving, then reverts on API error', async () => {
    const user = userEvent.setup();
    const patchSpy = vi.spyOn(issuedInvoiceService, 'patchNotes').mockImplementation(
      () =>
        new Promise((resolve) => {
          setTimeout(
            () =>
              resolve({
                success: false,
                error: { message: 'fallo', status: 500 },
              }),
            40,
          );
        }),
    );

    render(<MyInvoiceDetailClient invoiceId="inv-1" initialRow={baseRow} />);

    const textarea = await screen.findByRole('textbox', { name: /notas/i });
    await user.clear(textarea);
    await user.type(textarea, 'borrador cliente');
    expect(textarea).toHaveValue('borrador cliente');

    await user.click(screen.getByRole('button', { name: /Guardar notas/i }));

    await waitFor(() => expect(patchSpy).toHaveBeenCalled());
    expect(patchSpy.mock.calls[0]?.[1]).toBe('borrador cliente');

    await waitFor(() => {
      expect(screen.getByTestId('invoice-notes-optimistic')).toHaveTextContent('borrador cliente');
    });

    await waitFor(() => {
      expect(screen.getByTestId('invoice-notes-optimistic')).toHaveTextContent('versión servidor');
    });

    expect(screen.getByText(/fallo/i)).toBeInTheDocument();
  });

  it('reconciles optimistic state with server row after successful patch', async () => {
    const user = userEvent.setup();
    const updated: IssuedInvoice = { ...baseRow, notes: 'gardado no API' };
    vi.spyOn(issuedInvoiceService, 'patchNotes').mockResolvedValue({
      success: true,
      data: updated,
    });

    render(<MyInvoiceDetailClient invoiceId="inv-1" initialRow={baseRow} />);

    const textarea = await screen.findByRole('textbox', { name: /notas/i });
    await user.clear(textarea);
    await user.type(textarea, 'texto novo');
    await user.click(screen.getByRole('button', { name: /Guardar notas/i }));

    await waitFor(() => {
      expect(screen.getByTestId('invoice-notes-optimistic')).toHaveTextContent('gardado no API');
    });
  });
});

import React from 'react';
import { describe, it, expect, vi, beforeEach, afterEach, type MockInstance } from 'vitest';
import { render, waitFor } from '@testing-library/react';
import MyInvoicesListClient from './MyInvoicesListClient';
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

const sampleInvoices: IssuedInvoice[] = [
  {
    id: 'inv-1',
    customerId: 'cust-1',
    amount: 10,
    status: 'issued',
    notes: '',
    createdAt: '2026-01-02T00:00:00.000Z',
    updatedAt: '2026-01-02T00:00:00.000Z',
  },
  {
    id: 'inv-2',
    customerId: 'cust-1',
    amount: 20,
    status: 'issued',
    notes: '',
    createdAt: '2026-01-03T00:00:00.000Z',
    updatedAt: '2026-01-03T00:00:00.000Z',
  },
];

describe('MyInvoicesListClient', () => {
  let listMineSpy: MockInstance;

  beforeEach(() => {
    listMineSpy = vi.spyOn(issuedInvoiceService, 'listMine').mockResolvedValue({
      success: true,
      data: { items: [], total: 0 },
    });
  });

  afterEach(() => {
    listMineSpy.mockRestore();
  });

  it('does not call listMine on mount when server passed initial rows', async () => {
    render(<MyInvoicesListClient initialItems={sampleInvoices} />);
    await waitFor(() => {
      expect(listMineSpy).not.toHaveBeenCalled();
    });
  });

  it('calls listMine once when initial list is empty (client JWT path)', async () => {
    render(<MyInvoicesListClient initialItems={[]} />);
    await waitFor(() => {
      expect(listMineSpy).toHaveBeenCalledTimes(1);
    });
  });
});

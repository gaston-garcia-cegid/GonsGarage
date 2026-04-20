import { beforeEach, describe, expect, it, vi } from 'vitest';

vi.mock('../api-client', () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}));

import { apiClient } from '../api-client';
import { billingDocumentService } from './billing-document.service';
import { issuedInvoiceService } from './issued-invoice.service';
import { receivedInvoiceService } from './received-invoice.service';
import { supplierService } from './supplier.service';

const mocked = vi.mocked(apiClient);

describe('accounting API services (contract)', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('supplierService.list uses GET /suppliers with limit and offset', async () => {
    mocked.get.mockResolvedValueOnce({ success: true, data: { items: [], total: 0 } });
    await supplierService.list(25, 10);
    expect(mocked.get).toHaveBeenCalledWith('/suppliers?limit=25&offset=10');
  });

  it('receivedInvoiceService.create posts body', async () => {
    mocked.post.mockResolvedValueOnce({ success: true, data: { id: 'x' } });
    await receivedInvoiceService.create({
      vendorName: 'ACME',
      category: 'parts',
      amount: 12.5,
      invoiceDate: '2026-04-01',
      notes: 'n',
    });
    expect(mocked.post).toHaveBeenCalledWith('/received-invoices', {
      vendorName: 'ACME',
      category: 'parts',
      amount: 12.5,
      invoiceDate: '2026-04-01',
      notes: 'n',
    });
  });

  it('billingDocumentService.update puts to /billing-documents/:id', async () => {
    mocked.put.mockResolvedValueOnce({ success: true, data: { id: 'b' } });
    await billingDocumentService.update('uuid-1', {
      kind: 'irs',
      title: 'T',
      amount: 1,
      reference: 'R',
      notes: '',
    });
    expect(mocked.put).toHaveBeenCalledWith('/billing-documents/uuid-1', {
      kind: 'irs',
      title: 'T',
      amount: 1,
      reference: 'R',
      notes: '',
    });
  });

  it('issuedInvoiceService.listMine uses GET /invoices/me', async () => {
    mocked.get.mockResolvedValueOnce({ success: true, data: { items: [], total: 0 } });
    await issuedInvoiceService.listMine(15, 3);
    expect(mocked.get).toHaveBeenCalledWith('/invoices/me?limit=15&offset=3');
  });

  it('issuedInvoiceService.patchIssuedInvoice uses PATCH', async () => {
    mocked.patch.mockResolvedValueOnce({ success: true, data: { id: 'i' } });
    await issuedInvoiceService.patchIssuedInvoice('inv-1', { notes: 'hello', status: 'paid' });
    expect(mocked.patch).toHaveBeenCalledWith('/invoices/inv-1', { notes: 'hello', status: 'paid' });
  });

  it('issuedInvoiceService.removeStaff uses DELETE', async () => {
    mocked.delete.mockResolvedValueOnce({ success: true, data: {} });
    await issuedInvoiceService.removeStaff('inv-2');
    expect(mocked.delete).toHaveBeenCalledWith('/invoices/inv-2');
  });
});

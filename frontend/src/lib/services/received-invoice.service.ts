import { apiClient, type ApiResponse } from '../api-client';
import type { ItemsTotal, ReceivedInvoice } from '@/types/accounting';

export type { ReceivedInvoice };

export class ReceivedInvoiceService {
  private static instance: ReceivedInvoiceService;
  static getInstance(): ReceivedInvoiceService {
    if (!ReceivedInvoiceService.instance) ReceivedInvoiceService.instance = new ReceivedInvoiceService();
    return ReceivedInvoiceService.instance;
  }
  private constructor() {}

  async list(limit = 50, offset = 0): Promise<ApiResponse<ItemsTotal<ReceivedInvoice>>> {
    return apiClient.get<ItemsTotal<ReceivedInvoice>>(
      `/received-invoices?limit=${limit}&offset=${offset}`,
    );
  }

  async get(id: string): Promise<ApiResponse<ReceivedInvoice>> {
    return apiClient.get<ReceivedInvoice>(`/received-invoices/${id}`);
  }

  async create(
    body: Pick<ReceivedInvoice, 'vendorName' | 'category' | 'amount' | 'invoiceDate' | 'notes'> & {
      supplierId?: string;
    },
  ): Promise<ApiResponse<ReceivedInvoice>> {
    return apiClient.post<ReceivedInvoice>('/received-invoices', body);
  }

  async update(
    id: string,
    body: Pick<ReceivedInvoice, 'vendorName' | 'category' | 'amount' | 'invoiceDate' | 'notes'> & {
      supplierId?: string;
    },
  ): Promise<ApiResponse<ReceivedInvoice>> {
    return apiClient.put<ReceivedInvoice>(`/received-invoices/${id}`, body);
  }

  async remove(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.delete(`/received-invoices/${id}`);
  }
}

export const receivedInvoiceService = ReceivedInvoiceService.getInstance();

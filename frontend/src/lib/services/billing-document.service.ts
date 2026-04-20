import { apiClient, type ApiResponse } from '../api-client';
import type { BillingDocument, BillingDocumentKind, ItemsTotal } from '@/types/accounting';

export type { BillingDocument, BillingDocumentKind };

export class BillingDocumentService {
  private static instance: BillingDocumentService;
  static getInstance(): BillingDocumentService {
    if (!BillingDocumentService.instance) BillingDocumentService.instance = new BillingDocumentService();
    return BillingDocumentService.instance;
  }
  private constructor() {}

  async list(limit = 50, offset = 0): Promise<ApiResponse<ItemsTotal<BillingDocument>>> {
    return apiClient.get<ItemsTotal<BillingDocument>>(
      `/billing-documents?limit=${limit}&offset=${offset}`,
    );
  }

  async get(id: string): Promise<ApiResponse<BillingDocument>> {
    return apiClient.get<BillingDocument>(`/billing-documents/${id}`);
  }

  async create(
    body: Pick<BillingDocument, 'title' | 'amount' | 'reference' | 'notes'> & {
      kind: BillingDocumentKind;
      customerId?: string;
    },
  ): Promise<ApiResponse<BillingDocument>> {
    return apiClient.post<BillingDocument>('/billing-documents', body);
  }

  async update(
    id: string,
    body: Pick<BillingDocument, 'title' | 'amount' | 'reference' | 'notes'> & {
      kind: BillingDocumentKind;
      customerId?: string;
    },
  ): Promise<ApiResponse<BillingDocument>> {
    return apiClient.put<BillingDocument>(`/billing-documents/${id}`, body);
  }

  async remove(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.delete(`/billing-documents/${id}`);
  }
}

export const billingDocumentService = BillingDocumentService.getInstance();

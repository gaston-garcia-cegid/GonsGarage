import { apiClient, type ApiResponse } from '../api-client';
import type { ItemsTotal, Supplier } from '@/types/accounting';

export type { Supplier };

export class SupplierService {
  private static instance: SupplierService;
  static getInstance(): SupplierService {
    if (!SupplierService.instance) SupplierService.instance = new SupplierService();
    return SupplierService.instance;
  }
  private constructor() {}

  async list(limit = 50, offset = 0): Promise<ApiResponse<ItemsTotal<Supplier>>> {
    return apiClient.get<ItemsTotal<Supplier>>(`/suppliers?limit=${limit}&offset=${offset}`);
  }

  async get(id: string): Promise<ApiResponse<Supplier>> {
    return apiClient.get<Supplier>(`/suppliers/${id}`);
  }

  async create(body: Partial<Supplier> & { name: string }): Promise<ApiResponse<Supplier>> {
    return apiClient.post<Supplier>('/suppliers', body);
  }

  async update(id: string, body: Partial<Supplier> & { name: string }): Promise<ApiResponse<Supplier>> {
    return apiClient.put<Supplier>(`/suppliers/${id}`, body);
  }

  async remove(id: string): Promise<ApiResponse<unknown>> {
    return apiClient.delete(`/suppliers/${id}`);
  }
}

export const supplierService = SupplierService.getInstance();

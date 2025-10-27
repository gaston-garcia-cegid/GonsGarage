// ✅ Vehicle service using centralized API client
// Provides vehicle management methods following Agent.md patterns

import { apiClient, type ApiResponse } from '../api-client';

// ✅ Vehicle interfaces following Agent.md camelCase conventions
export interface Vehicle {
  id: string;
  make: string;
  model: string;
  year: number;
  licensePlate: string;
  vin?: string;
  color?: string;
  mileage?: number;
  ownerId: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateVehicleRequest {
  make: string;
  model: string;
  year: number;
  licensePlate: string;
  vin?: string;
  color?: string;
  mileage?: number;
}

export interface UpdateVehicleRequest extends Partial<CreateVehicleRequest> {
  id: string;
}

export interface VehicleFilters {
  ownerId?: string;
  make?: string;
  model?: string;
  year?: number;
  limit?: number;
  offset?: number;
}

export interface VehicleListResponse {
  vehicles: Vehicle[];
  total: number;
  limit: number;
  offset: number;
}

// ✅ Vehicle service class following Agent.md clean architecture
export class VehicleService {
  private static instance: VehicleService;

  // ✅ Singleton pattern for service consistency
  static getInstance(): VehicleService {
    if (!VehicleService.instance) {
      VehicleService.instance = new VehicleService();
    }
    return VehicleService.instance;
  }

  private constructor() {
    // Private constructor for singleton
  }

  // ✅ Get all vehicles with optional filters
  async getVehicles(filters: VehicleFilters = {}): Promise<ApiResponse<VehicleListResponse>> {
    try {
      const queryParams = new URLSearchParams();
      
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null) {
          queryParams.append(key, String(value));
        }
      });

      const queryString = queryParams.toString();
      const endpoint = `/vehicles${queryString ? `?${queryString}` : ''}`;

      const response = await apiClient.get<VehicleListResponse>(endpoint);
      if (response.success) {
        console.log('Vehicles:', response.data);
      } else {
        console.error('Error:', response.error?.message);
      }
      
      return response;
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to fetch vehicles',
          status: 0,
          code: 'FETCH_VEHICLES_ERROR'
        }
      };
    }
  }

  // ✅ Get single vehicle by ID
  async getVehicle(id: string): Promise<ApiResponse<Vehicle>> {
    try {
      return await apiClient.get<Vehicle>(`/vehicles/${id}`);
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to fetch vehicle',
          status: 0,
          code: 'FETCH_VEHICLE_ERROR'
        }
      };
    }
  }

  // ✅ Create new vehicle
  async createVehicle(vehicleData: CreateVehicleRequest): Promise<ApiResponse<Vehicle>> {
    try {
      return await apiClient.post<Vehicle>('/vehicles', vehicleData);
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to create vehicle',
          status: 0,
          code: 'CREATE_VEHICLE_ERROR'
        }
      };
    }
  }

  // ✅ Update existing vehicle
  async updateVehicle(id: string, updates: Partial<CreateVehicleRequest>): Promise<ApiResponse<Vehicle>> {
    try {
      return await apiClient.put<Vehicle>(`/vehicles/${id}`, updates);
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to update vehicle',
          status: 0,
          code: 'UPDATE_VEHICLE_ERROR'
        }
      };
    }
  }

  // ✅ Delete vehicle
  async deleteVehicle(id: string): Promise<ApiResponse<{ message: string }>> {
    try {
      return await apiClient.delete<{ message: string }>(`/vehicles/${id}`);
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to delete vehicle',
          status: 0,
          code: 'DELETE_VEHICLE_ERROR'
        }
      };
    }
  }

  // ✅ Get vehicles by owner ID (for client dashboard)
  async getVehiclesByOwner(ownerId: string): Promise<ApiResponse<Vehicle[]>> {
    try {
      return await apiClient.get<Vehicle[]>(`/vehicles/owner/${ownerId}`);
    } catch {
      return {
        success: false,
        error: {
          message: 'Failed to fetch owner vehicles',
          status: 0,
          code: 'FETCH_OWNER_VEHICLES_ERROR'
        }
      };
    }
  }
}

// ✅ Export singleton instance for application use
export const vehicleService = VehicleService.getInstance();

// ✅ Default export for convenience
export default vehicleService;
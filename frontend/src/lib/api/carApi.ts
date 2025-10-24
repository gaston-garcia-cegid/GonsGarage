import { Car, CreateCarRequest, UpdateCarRequest } from '@/types/car';

export interface ApiResponse<T> {
  data: T | null;
  error: {
    message: string;
    status?: number;
  } | null;
}

class CarApiClient {
  private baseUrl: string;

  constructor(baseUrl: string = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080') {
    this.baseUrl = baseUrl;
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    try {
      const token = localStorage.getItem('authToken');
      const url = `${this.baseUrl}${endpoint}`;
      
      const response = await fetch(url, {
        ...options,
        headers: {
          'Content-Type': 'application/json',
          ...(token && { Authorization: `Bearer ${token}` }),
          ...options.headers,
        },
      });

      if (!response.ok) {
        const errorText = await response.text();
        return {
          data: null,
          error: {
            message: errorText || `HTTP ${response.status}: ${response.statusText}`,
            status: response.status
          }
        };
      }

      const data = await response.json();
      return { data, error: null };
    } catch (error) {
      return {
        data: null,
        error: {
          message: error instanceof Error ? error.message : 'Network error occurred'
        }
      };
    }
  }

  // Create car - following Agent.md error handling
  async createCar(carData: CreateCarRequest): Promise<ApiResponse<Car>> {
    return this.request<Car>('/api/v1/cars', {
      method: 'POST',
      body: JSON.stringify(carData),
    });
  }

  // Read cars - following Agent.md naming conventions
  async getCars(): Promise<ApiResponse<Car[]>> {
    return this.request<Car[]>('/api/v1/cars');
  }

  // Read single car
  async getCar(id: string): Promise<ApiResponse<Car>> {
    return this.request<Car>(`/api/v1/cars/${id}`);
  }

  // Update car - following Agent.md clean architecture
  async updateCar(id: string, carData: Partial<CreateCarRequest>): Promise<ApiResponse<Car>> {
    return this.request<Car>(`/api/v1/cars/${id}`, {
      method: 'PUT',
      body: JSON.stringify(carData),
    });
  }

  // Delete car - following Agent.md proper error handling
  async deleteCar(id: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/api/v1/cars/${id}`, {
      method: 'DELETE',
    });
  }

  // Get cars with repairs (complex read operation)
  async getCarWithRepairs(id: string): Promise<ApiResponse<Car>> {
    return this.request<Car>(`/api/v1/cars/${id}/repairs`);
  }
}

export const carApi = new CarApiClient();
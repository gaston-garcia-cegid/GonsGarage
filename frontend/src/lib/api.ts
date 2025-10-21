interface ApiError {
  message: string;
  status: number;
  code?: string;
}

interface ApiResponse<T> {
  data?: T;
  error?: ApiError;
}

// Types for API responses
export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  role: string;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface Employee {
  id: string;
  first_name: string;
  last_name: string;
  email: string;
  phone?: string;
  department: string;
  position: string;
  hire_date: string;
  salary: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  token: string;
  message: string;
}

export interface CreateEmployeeRequest {
  first_name: string;
  last_name: string;
  email: string;
  phone?: string;
  department: string;
  position: string;
  hire_date: string;
  salary: number;
}

export interface EmployeeListResponse {
  employees: Employee[];
  total: number;
  limit: number;
  offset: number;
}

//const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

class ApiClient {
  private baseURL: string;
  private token: string | null = null;

  constructor() {
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';
    if (typeof window !== 'undefined') {
      this.token = localStorage.getItem('token');
    }
  }

  setToken(token: string) {
    this.token = token;
    if (typeof window !== 'undefined') {
      localStorage.setItem('token', token);
    }
  }

  clearToken() {
    this.token = null;
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }

  private async request<T>(
    endpoint: string,
    options: RequestInit = {}
  ): Promise<ApiResponse<T>> {
    const url = `${this.baseURL}${endpoint}`;
    
    const config: RequestInit = {
      headers: {
        'Content-Type': 'application/json',
        ...(this.token && { Authorization: `Bearer ${this.token}` }),
        ...options.headers,
      },
      ...options,
    };

    try {
      const response = await fetch(url, config);
      
      if (response.status === 204) {
        return { data: {} as T };
      }

      const data = await response.json();

      if (!response.ok) {
        return {
          error: {
            message: data.error || data.message || 'An error occurred',
            status: response.status,
            code: data.code,
          },
        };
      }

      return { data };
    } catch (error) {
      console.error('API request failed:', error);
      return {
        error: {
          message: 'Network error occurred, Line 134',
          status: 0,
        },
      };
    }
  }

  // Authentication endpoints
  async login(loginData: LoginRequest): Promise<ApiResponse<LoginResponse>> {
    return this.request<LoginResponse>('/auth/login', {
      method: 'POST',
      body: JSON.stringify(loginData),
    });
  }

  async logout(): Promise<ApiResponse<{ message: string }>> {
    const response = await this.request<{ message: string }>('/auth/logout', {
      method: 'POST',
    });
    this.clearToken();
    return response;
  }

  // Employee endpoints
  async getEmployees(limit = 20, offset = 0): Promise<ApiResponse<EmployeeListResponse>> {
    const response = await this.request<EmployeeListResponse>(`/employees?limit=${limit}&offset=${offset}`);
    return response;
  }

  async getEmployee(id: string): Promise<ApiResponse<Employee>> {
    const response = await this.request<Employee>(`/employees/${id}`);
    return response;
  }

  async createEmployee(employeeData: CreateEmployeeRequest): Promise<ApiResponse<Employee>> {
    return this.request<Employee>('/employees', {
      method: 'POST',
      body: JSON.stringify(employeeData),
    });
  }

  async updateEmployee(id: string, employeeData: Partial<CreateEmployeeRequest>): Promise<ApiResponse<Employee>> {
    return this.request<Employee>(`/employees/${id}`, {
      method: 'PUT',
      body: JSON.stringify(employeeData),
    });
  }

  async deleteEmployee(id: string): Promise<ApiResponse<{ message: string }>> {
    return this.request<{ message: string }>(`/employees/${id}`, {
      method: 'DELETE',
    });
  }

  // Profile endpoint
  async getProfile(): Promise<ApiResponse<User>> {
    return this.request<User>('/auth/profile');
  }
}

export const apiClient = new ApiClient();
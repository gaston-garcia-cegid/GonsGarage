// src/shared/types/index.ts
export interface Car {
  id: string;
  make: string;
  model: string;
  year: number;
  licensePlate: string;
  vin?: string;
  createdAt: string;
  updatedAt: string;
}

export interface Repair {
  id: string;
  car_id: string;
  description: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  cost: number;
  created_at: string;
  completed_at?: string;
}

export interface Appointment {
  id: string;
  customerId: string;
  carId: string;
  serviceType: string;
  scheduledAt: string;
  status: 'scheduled' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled';
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateAppointmentRequest {
  carId: string;         // ✅ camelCase
  serviceType: string;   // ✅ camelCase
  scheduledAt: string;   // ✅ camelCase
  notes?: string;
}

export interface UpdateAppointmentRequest {
  serviceType?: string;  // ✅ camelCase
  scheduledAt?: string;  // ✅ camelCase
  status?: 'scheduled' | 'in-progress' | 'completed' | 'cancelled';
  notes?: string;
}
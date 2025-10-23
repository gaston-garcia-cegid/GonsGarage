// src/shared/types/index.ts
export interface Car {
  id: string;
  make: string;
  model: string;
  year: number;
  license_plate: string;
  vin?: string;
  created_at: string;
  updated_at: string;
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
  car_id: string;
  service_type: string;
  scheduled_at: string;
  status: 'scheduled' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled';
  notes?: string;
  created_at: string;
}
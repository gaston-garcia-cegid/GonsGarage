// src/shared/types/index.ts
// export interface Car {
//   id: string;
//   make: string;
//   model: string;
//   year: number;
//   licensePlate: string;
//   vin?: string;
//   createdAt: string;
//   updatedAt: string;
// }

export interface Repair {
  id: string;
  car_id: string;
  description: string;
  status: 'pending' | 'in_progress' | 'completed' | 'cancelled';
  cost: number;
  created_at: string;
  completed_at?: string;
}

// export interface Appointment {
//   id: string;
//   customerId: string;
//   carId: string;
//   serviceType: string;
//   scheduledAt: string;
//   status: 'scheduled' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled';
//   notes?: string;
//   createdAt: string;
//   updatedAt: string;
// }

// export interface CreateAppointmentRequest {
//   carId: string;         // ✅ camelCase
//   serviceType: string;   // ✅ camelCase
//   scheduledAt: string;   // ✅ camelCase
//   notes?: string;
// }

// export interface UpdateAppointmentRequest {
//   serviceType?: string;  // ✅ camelCase
//   scheduledAt?: string;  // ✅ camelCase
//   status?: 'scheduled' | 'in-progress' | 'completed' | 'cancelled';
//   notes?: string;
// }

export const SERVICE_TYPES = [
  { id: 'oil_change', name: 'Oil Change', description: 'Regular oil and filter change' },
  { id: 'brake_service', name: 'Brake Service', description: 'Brake pads, rotors, and fluid' },
  { id: 'tire_service', name: 'Tire Service', description: 'Tire rotation, alignment, replacement' },
  { id: 'engine_diagnostic', name: 'Engine Diagnostic', description: 'Check engine light and diagnostics' },
  { id: 'transmission_service', name: 'Transmission Service', description: 'Transmission fluid and inspection' },
  { id: 'air_conditioning', name: 'Air Conditioning', description: 'A/C repair and maintenance' },
  { id: 'battery_service', name: 'Battery Service', description: 'Battery testing and replacement' },
  { id: 'general_maintenance', name: 'General Maintenance', description: 'Multi-point inspection' },
  { id: 'other', name: 'Other', description: 'Custom service request' },
] as const;
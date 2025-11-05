// Appointment type definition

export interface Appointment {
  id: string;
  clientName: string;
  carId: string;
  service: string;
  date: string;
  time: string;
  status: 'scheduled' | 'completed' | 'cancelled' | 'in_progress' | 'confirmed';
  notes?: string;
  createdAt: string;
  updatedAt: string;
  deletedAt?: string;
}

export interface CreateAppointmentRequest {
  clientName: string;
  carId: string;
  service: string;
  date: string;
  time: string;
  status: 'scheduled' | 'completed' | 'cancelled' | 'in_progress' | 'confirmed'; // ✅ camelCase per Agent.md
  notes?: string;
  createdAt: string; // ✅ camelCase per Agent.md
  updatedAt: string; // ✅ camelCase per Agent.md
  deletedAt?: string; // ✅ camelCase per Agent.md

}

export interface UpdateAppointmentRequest {
  service?: string;
  date?: string;
  time?: string;
  status?: 'scheduled' | 'completed' | 'cancelled'; // ✅ camelCase per Agent.md
  createdAt?: string; // ✅ camelCase per Agent.md
  updatedAt?: string; // ✅ camelCase per Agent.md
  deletedAt?: string; // ✅ camelCase per Agent.md
  notes?: string;
}

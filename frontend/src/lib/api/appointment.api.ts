import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/shared/types';
import { apiClient } from '@/lib/api-client';

export const appointmentApi = {
  getAppointments: async (): Promise<Appointment[]> => {
    const response = await apiClient.get<Appointment[]>('/appointments');
    return response.data ?? [];
  },

  getAppointment: async (id: string): Promise<Appointment | null> => {
    const response = await apiClient.get<Appointment | null>(`/appointments/${id}`);
    return response.data ?? null;
  },

  createAppointment: async (data: CreateAppointmentRequest): Promise<Appointment | null> => {
    const response = await apiClient.post<Appointment | null>('/appointments', {
      carId: data.carId,         // ✅ camelCase per Agent.md
      serviceType: data.serviceType, // ✅ camelCase per Agent.md
      scheduledAt: data.scheduledAt, // ✅ camelCase per Agent.md
      notes: data.notes,
    });
    return response.data ?? null;
  },

  updateAppointment: async (id: string, data: UpdateAppointmentRequest): Promise<Appointment | null> => {
    const response = await apiClient.put<Appointment | null>(`/appointments/${id}`, {
      serviceType: data.serviceType, // ✅ camelCase per Agent.md
      scheduledAt: data.scheduledAt, // ✅ camelCase per Agent.md
      status: data.status,
      notes: data.notes,
    });
    return response.data ?? null;
  },

  deleteAppointment: async (id: string): Promise<void> => {
    await apiClient.delete(`/appointments/${id}`);
  },

  cancelAppointment: async (id: string): Promise<Appointment | null> => {
    const response = await apiClient.patch<Appointment | null>(`/appointments/${id}/cancel`);
    return response.data ?? null;
  },

  confirmAppointment: async (id: string): Promise<Appointment | null> => {
    const response = await apiClient.patch<Appointment | null>(`/appointments/${id}/confirm`);
    return response.data ?? null;
  },

  completeAppointment: async (id: string): Promise<Appointment | null> => {
    const response = await apiClient.patch<Appointment | null>(`/appointments/${id}/complete`);
    return response.data ?? null;
  },
};
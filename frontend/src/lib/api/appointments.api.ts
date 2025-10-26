import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/shared/types';
import { apiClient } from './client';

export const appointmentApi = {
  getAppointments: async (): Promise<Appointment[]> => {
    const response = await apiClient.get('/appointments');
    return response.data;
  },

  getAppointment: async (id: string): Promise<Appointment> => {
    const response = await apiClient.get(`/appointments/${id}`);
    return response.data;
  },

  createAppointment: async (data: CreateAppointmentRequest): Promise<Appointment> => {
    const response = await apiClient.post('/appointments', {
      carId: data.carId,         // ✅ camelCase
      serviceType: data.serviceType, // ✅ camelCase
      scheduledAt: data.scheduledAt, // ✅ camelCase
      notes: data.notes,
    });
    return response.data;
  },

  updateAppointment: async (id: string, data: UpdateAppointmentRequest): Promise<Appointment> => {
    const response = await apiClient.put(`/appointments/${id}`, {
      serviceType: data.serviceType, // ✅ camelCase
      scheduledAt: data.scheduledAt, // ✅ camelCase
      status: data.status,
      notes: data.notes,
    });
    return response.data;
  },

  deleteAppointment: async (id: string): Promise<void> => {
    await apiClient.delete(`/appointments/${id}`);
  },
};
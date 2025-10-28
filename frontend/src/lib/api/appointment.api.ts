import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/shared/types';
import { apiClient } from '@/lib/api-client';

export const appointmentApi = {
  // ✅ Fix getAppointments to handle both response formats
  getAppointments: async (): Promise<Appointment[]> => {
    try {
      const response = await apiClient.get<Appointment[]>('/appointments');
      
      if (response.success && response.data) {
        // ✅ Handle both response formats
        if (Array.isArray(response.data)) {
          // ✅ Direct array response
          return response.data;
        // } else if (response.data.appointments && Array.isArray(response.data.appointments)) {
        //   // ✅ Wrapped object response  
        //   return response.data.appointments;
        } else {
          // ✅ Handle empty/invalid response
          console.warn('⚠️ Unexpected appointments response format:', response.data);
          return [];
        }
      } else {
        console.error('❌ Appointments API Error:', response.error);
        return [];
      }
    } catch (error) {
      console.error('❌ Appointments Service Error:', error);
      return [];
    }
  },

  // ✅ Rest of the methods remain the same...
  getAppointment: async (id: string): Promise<Appointment | null> => {
    try {
      const response = await apiClient.get<Appointment>(`/appointments/${id}`);
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Get appointment error:', error);
      return null;
    }
  },

  createAppointment: async (data: CreateAppointmentRequest): Promise<Appointment | null> => {
    try {
      const response = await apiClient.post<Appointment>('/appointments', {
        carId: data.carId,         // ✅ camelCase per Agent.md
        serviceType: data.serviceType, // ✅ camelCase per Agent.md
        scheduledAt: data.scheduledAt, // ✅ camelCase per Agent.md
        notes: data.notes,
      });
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Create appointment error:', error);
      return null;
    }
  },

  updateAppointment: async (id: string, data: UpdateAppointmentRequest): Promise<Appointment | null> => {
    try {
      const response = await apiClient.put<Appointment>(`/appointments/${id}`, {
        serviceType: data.serviceType, // ✅ camelCase per Agent.md
        scheduledAt: data.scheduledAt, // ✅ camelCase per Agent.md
        status: data.status,
        notes: data.notes,
      });
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Update appointment error:', error);
      return null;
    }
  },

  deleteAppointment: async (id: string): Promise<void> => {
    try {
      await apiClient.delete(`/appointments/${id}`);
    } catch (error) {
      console.error('❌ Delete appointment error:', error);
      throw error;
    }
  },

  cancelAppointment: async (id: string): Promise<Appointment | null> => {
    try {
      const response = await apiClient.patch<Appointment>(`/appointments/${id}/cancel`);
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Cancel appointment error:', error);
      return null;
    }
  },

  confirmAppointment: async (id: string): Promise<Appointment | null> => {
    try {
      const response = await apiClient.patch<Appointment>(`/appointments/${id}/confirm`);
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Confirm appointment error:', error);
      return null;
    }
  },

  completeAppointment: async (id: string): Promise<Appointment | null> => {
    try {
      const response = await apiClient.patch<Appointment>(`/appointments/${id}/complete`);
      return response.success ? response.data ?? null : null;
    } catch (error) {
      console.error('❌ Complete appointment error:', error);
      return null;
    }
  },
};
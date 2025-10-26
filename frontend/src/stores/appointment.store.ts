// ✅ Appointment Store - Zustand + Immer + TypeScript following Agent.md patterns
// Provides centralized appointment state management with CRUD operations

import { create } from 'zustand';
import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/shared/types';
import { appointmentApi } from '@/lib/api';

// ✅ Appointment interfaces following Agent.md camelCase conventions
export interface Appointment {
  id: string;
  customerId: string; // ✅ camelCase per Agent.md
  carId: string; // ✅ camelCase per Agent.md
  serviceType: string; // ✅ camelCase per Agent.md
  status: 'scheduled' | 'confirmed' | 'completed' | 'cancelled';
  scheduledAt: string; // ✅ camelCase per Agent.md
  notes?: string;
  createdAt: string; // ✅ camelCase per Agent.md
  updatedAt: string; // ✅ camelCase per Agent.md
}

export interface CreateAppointmentRequest {
  carId: string; // ✅ camelCase per Agent.md
  serviceType: string; // ✅ camelCase per Agent.md
  scheduledAt: string; // ✅ camelCase per Agent.md
  notes?: string;
}

export interface UpdateAppointmentRequest {
  serviceType?: string; // ✅ camelCase per Agent.md
  scheduledAt?: string; // ✅ camelCase per Agent.md
  notes?: string;
  status?: Appointment['status'];
}

// ✅ Service types enum following Agent.md conventions
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

// ✅ Appointment state interface following Agent.md conventions
interface AppointmentState {
  // Data state
  appointments: Appointment[];
  selectedAppointment: Appointment | null;
  
  // Loading and error states
  isLoading: boolean;
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;
  error: string | null;
  
  // Filtering and pagination
  totalAppointments: number;
  currentPage: number;
  pageSize: number;
  filters: {
    status?: Appointment['status'] | 'all';
    serviceType?: string;
    customerId?: string;
    carId?: string;
    dateFrom?: string;
    dateTo?: string;
  };
  
  // Actions
  fetchAppointments: (filters?: AppointmentState['filters']) => Promise<void>;
  fetchAppointmentById: (id: string) => Promise<void>;
  createAppointment: (appointmentData: CreateAppointmentRequest) => Promise<boolean>;
  updateAppointment: (id: string, appointmentData: UpdateAppointmentRequest) => Promise<boolean>;
  cancelAppointment: (id: string) => Promise<boolean>;
  confirmAppointment: (id: string) => Promise<boolean>;
  completeAppointment: (id: string) => Promise<boolean>;
  deleteAppointment: (id: string) => Promise<boolean>;
  selectAppointment: (appointment: Appointment | null) => void;
  setFilters: (filters: Partial<AppointmentState['filters']>) => void;
  setPage: (page: number) => void;
  clearError: () => void;
  reset: () => void;
}

// ✅ Initial state following Agent.md conventions
const initialState = {
  appointments: [],
  selectedAppointment: null,
  isLoading: false,
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  error: null,
  totalAppointments: 0,
  currentPage: 1,
  pageSize: 10,
  filters: { status: 'all' as const },
};

// ✅ Appointment store implementation with Zustand + Immer
export const useAppointmentStore = create<AppointmentState>()(
  immer((set, get) => ({
    // Initial state
    ...initialState,
    
    // ✅ Fetch appointments with optional filters
    fetchAppointments: async (filters) => {
      set((state) => {
        state.isLoading = true;
        state.error = null;
        if (filters) {
          state.filters = { ...state.filters, ...filters };
        }
      });
      
      try {
        const response = await appointmentApi.getAppointments();
        
        if (response.data) {
          set((state) => {
            state.appointments = response.data!.appointments;
            state.totalAppointments = response.data!.total;
            state.isLoading = false;
          });
        } else {
          throw new Error(response.error?.message || 'Failed to fetch appointments');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to fetch appointments';
          state.isLoading = false;
        });
      }
    },
    
    // ✅ Fetch single appointment by ID
    fetchAppointmentById: async (id: string) => {
      set((state) => {
        state.isLoading = true;
        state.error = null;
      });
      
      try {
        const response = await appointmentApi.getAppointment(id);
        
        if (response.data) {
          set((state) => {
            state.selectedAppointment = response.data!;
            state.isLoading = false;
          });
        } else {
          throw new Error(response.error?.message || 'Appointment not found');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to fetch appointment';
          state.isLoading = false;
        });
      }
    },
    
    // ✅ Create new appointment
    createAppointment: async (appointmentData: CreateAppointmentRequest): Promise<boolean> => {
      set((state) => {
        state.isCreating = true;
        state.error = null;
      });
      
      try {
        const response = await appointmentApi.createAppointment(appointmentData);
        
        if (response.data) {
          set((state) => {
            state.appointments.unshift(response.data!); // Add to beginning of list
            state.totalAppointments += 1;
            state.isCreating = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to create appointment');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to create appointment';
          state.isCreating = false;
        });
        return false;
      }
    },
    
    // ✅ Update existing appointment
    updateAppointment: async (id: string, appointmentData: UpdateAppointmentRequest): Promise<boolean> => {
      set((state) => {
        state.isUpdating = true;
        state.error = null;
      });
      
      try {
        const response = await appointmentApi.updateAppointment(id, appointmentData);
        
        if (response.data) {
          set((state) => {
            const index = state.appointments.findIndex(appointment => appointment.id === id);
            if (index !== -1) {
              state.appointments[index] = response.data!;
            }
            if (state.selectedAppointment?.id === id) {
              state.selectedAppointment = response.data!;
            }
            state.isUpdating = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to update appointment');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to update appointment';
          state.isUpdating = false;
        });
        return false;
      }
    },
    
    // ✅ Cancel appointment (update status)
    cancelAppointment: async (id: string): Promise<boolean> => {
      return get().updateAppointment(id, { status: 'cancelled' });
    },
    
    // ✅ Confirm appointment (update status)
    confirmAppointment: async (id: string): Promise<boolean> => {
      return get().updateAppointment(id, { status: 'confirmed' });
    },
    
    // ✅ Complete appointment (update status)
    completeAppointment: async (id: string): Promise<boolean> => {
      return get().updateAppointment(id, { status: 'completed' });
    },
    
    // ✅ Delete appointment
    deleteAppointment: async (id: string): Promise<boolean> => {
      set((state) => {
        state.isDeleting = true;
        state.error = null;
      });
      
      try {
        const response = await appointmentApi.deleteAppointment(id);
        
        if (response.success) {
          set((state) => {
            state.appointments = state.appointments.filter(appointment => appointment.id !== id);
            state.totalAppointments -= 1;
            if (state.selectedAppointment?.id === id) {
              state.selectedAppointment = null;
            }
            state.isDeleting = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to delete appointment');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to delete appointment';
          state.isDeleting = false;
        });
        return false;
      }
    },
    
    // ✅ Select appointment for detailed view
    selectAppointment: (appointment: Appointment | null) => {
      set((state) => {
        state.selectedAppointment = appointment;
      });
    },
    
    // ✅ Update filters
    setFilters: (filters: Partial<AppointmentState['filters']>) => {
      set((state) => {
        state.filters = { ...state.filters, ...filters };
        state.currentPage = 1; // Reset to first page when filters change
      });
    },
    
    // ✅ Set current page
    setPage: (page: number) => {
      set((state) => {
        state.currentPage = page;
      });
    },
    
    // ✅ Clear error state
    clearError: () => {
      set((state) => {
        state.error = null;
      });
    },
    
    // ✅ Reset store to initial state
    reset: () => {
      set(() => initialState);
    },
  }))
);

// ✅ Convenience hooks following AuthStore pattern
export const useAppointments = () => {
  const store = useAppointmentStore();
  return {
    // Data
    appointments: store.appointments,
    selectedAppointment: store.selectedAppointment,
    totalAppointments: store.totalAppointments,
    currentPage: store.currentPage,
    pageSize: store.pageSize,
    filters: store.filters,
    
    // Loading states
    isLoading: store.isLoading,
    isCreating: store.isCreating,
    isUpdating: store.isUpdating,
    isDeleting: store.isDeleting,
    error: store.error,
    
    // Actions
    fetchAppointments: store.fetchAppointments,
    fetchAppointmentById: store.fetchAppointmentById,
    createAppointment: store.createAppointment,
    updateAppointment: store.updateAppointment,
    cancelAppointment: store.cancelAppointment,
    confirmAppointment: store.confirmAppointment,
    completeAppointment: store.completeAppointment,
    deleteAppointment: store.deleteAppointment,
    selectAppointment: store.selectAppointment,
    setFilters: store.setFilters,
    setPage: store.setPage,
    clearError: store.clearError,
    reset: store.reset,
  };
};

// ✅ Filtered appointments hook
export const useFilteredAppointments = () => {
  const { appointments, filters } = useAppointments();
  
  const filteredAppointments = appointments.filter((appointment) => {
    // Status filter
    if (filters.status && filters.status !== 'all' && appointment.status !== filters.status) {
      return false;
    }
    
    // Service type filter
    if (filters.serviceType && appointment.serviceType !== filters.serviceType) {
      return false;
    }
    
    // Customer filter
    if (filters.customerId && appointment.customerId !== filters.customerId) {
      return false;
    }
    
    // Car filter
    if (filters.carId && appointment.carId !== filters.carId) {
      return false;
    }
    
    // Date range filter
    if (filters.dateFrom || filters.dateTo) {
      const appointmentDate = new Date(appointment.scheduledAt);
      
      if (filters.dateFrom && appointmentDate < new Date(filters.dateFrom)) {
        return false;
      }
      
      if (filters.dateTo && appointmentDate > new Date(filters.dateTo)) {
        return false;
      }
    }
    
    return true;
  });
  
  return filteredAppointments;
};

// ✅ Specific hooks for common use cases
export const useAppointmentList = () => {
  const { appointments, isLoading, error, fetchAppointments, setFilters, setPage } = useAppointments();
  const filteredAppointments = useFilteredAppointments();
  
  return { 
    appointments: filteredAppointments,
    allAppointments: appointments,
    isLoading, 
    error, 
    fetchAppointments, 
    setFilters, 
    setPage 
  };
};

export const useAppointmentDetails = () => {
  const { selectedAppointment, isLoading, error, fetchAppointmentById, selectAppointment } = useAppointments();
  return { selectedAppointment, isLoading, error, fetchAppointmentById, selectAppointment };
};

export const useAppointmentMutations = () => {
  const { 
    isCreating, 
    isUpdating, 
    isDeleting, 
    createAppointment, 
    updateAppointment,
    cancelAppointment,
    confirmAppointment,
    completeAppointment,
    deleteAppointment,
    error,
    clearError
  } = useAppointments();
  
  return { 
    isCreating, 
    isUpdating, 
    isDeleting, 
    createAppointment, 
    updateAppointment,
    cancelAppointment,
    confirmAppointment,
    completeAppointment,
    deleteAppointment,
    error,
    clearError
  };
};
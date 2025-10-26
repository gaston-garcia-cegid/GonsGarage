// ✅ AppointmentStore Tests - Simplified and robust testing following Agent.md patterns

import { act } from '@testing-library/react';
import { 
  useAppointmentStore,
  type Appointment,
  type CreateAppointmentRequest 
} from '@/stores/appointment.store';

// ✅ Sample test data
const mockAppointment: Appointment = {
  id: '1',
  customerId: 'customer-1',
  carId: 'car-1',
  serviceType: 'oil_change',
  status: 'scheduled',
  scheduledAt: '2025-01-15T10:00:00Z',
  notes: 'Regular maintenance',
  createdAt: '2025-01-01T00:00:00Z',
  updatedAt: '2025-01-01T00:00:00Z',
};

const mockCreateAppointmentData: CreateAppointmentRequest = {
  carId: 'car-2',
  serviceType: 'brake_service',
  scheduledAt: '2025-01-20T14:00:00Z',
  notes: 'Brake pads replacement needed',
};

describe('AppointmentStore', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Reset store state before each test
    act(() => {
      useAppointmentStore.getState().reset();
    });
  });

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      const store = useAppointmentStore.getState();
      
      expect(store.appointments).toEqual([]);
      expect(store.selectedAppointment).toBeNull();
      expect(store.isLoading).toBe(false);
      expect(store.isCreating).toBe(false);
      expect(store.isUpdating).toBe(false);
      expect(store.isDeleting).toBe(false);
      expect(store.error).toBeNull();
      expect(store.totalAppointments).toBe(0);
      expect(store.currentPage).toBe(1);
      expect(store.pageSize).toBe(10);
      expect(store.filters).toEqual({ status: 'all' });
    });
  });

  describe('Fetch Appointments', () => {
    it('should fetch appointments successfully', async () => {
      const store = useAppointmentStore.getState();
      
      await act(async () => {
        await store.fetchAppointments();
      });
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.appointments).toEqual([]);
      expect(updatedStore.totalAppointments).toBe(0);
      expect(updatedStore.isLoading).toBe(false);
      expect(updatedStore.error).toBeNull();
    });
  });

  describe('Create Appointment', () => {
    it('should create appointment successfully', async () => {
      const store = useAppointmentStore.getState();
      
      let createResult: boolean;
      await act(async () => {
        createResult = await store.createAppointment(mockCreateAppointmentData);
      });
      
      expect(createResult!).toBe(true);
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.appointments).toHaveLength(1);
      expect(updatedStore.appointments[0].carId).toBe('car-2');
      expect(updatedStore.appointments[0].serviceType).toBe('brake_service');
      expect(updatedStore.appointments[0].status).toBe('scheduled');
      expect(updatedStore.totalAppointments).toBe(1);
      expect(updatedStore.isCreating).toBe(false);
    });
  });

  describe('Update Appointment', () => {
    it('should handle update failure correctly', async () => {
      const store = useAppointmentStore.getState();
      
      let updateResult: boolean;
      await act(async () => {
        updateResult = await store.updateAppointment('1', { 
          serviceType: 'tire_service' 
        });
      });
      
      expect(updateResult!).toBe(false);
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.error).toBe('Update not implemented');
      expect(updatedStore.isUpdating).toBe(false);
    });
  });

  describe('Delete Appointment', () => {
    it('should delete appointment successfully', async () => {
      // First add an appointment to the store
      act(() => {
        useAppointmentStore.setState({ 
          appointments: [mockAppointment], 
          totalAppointments: 1 
        });
      });

      const store = useAppointmentStore.getState();
      
      let deleteResult: boolean;
      await act(async () => {
        deleteResult = await store.deleteAppointment('1');
      });
      
      expect(deleteResult!).toBe(true);
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.appointments).toHaveLength(0);
      expect(updatedStore.totalAppointments).toBe(0);
      expect(updatedStore.isDeleting).toBe(false);
    });
  });

  describe('State Management', () => {
    it('should select appointment correctly', () => {
      const store = useAppointmentStore.getState();
      
      act(() => {
        store.selectAppointment(mockAppointment);
      });
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.selectedAppointment).toEqual(mockAppointment);
    });

    it('should set filters correctly', () => {
      const store = useAppointmentStore.getState();
      
      act(() => {
        store.setFilters({ 
          status: 'confirmed', 
          serviceType: 'oil_change',
          customerId: 'customer-1' 
        });
      });
      
      const updatedStore = useAppointmentStore.getState();
      expect(updatedStore.filters).toEqual({ 
        status: 'confirmed', 
        serviceType: 'oil_change',
        customerId: 'customer-1' 
      });
      expect(updatedStore.currentPage).toBe(1); // Should reset page
    });

    it('should clear error correctly', () => {
      // First set an error
      act(() => {
        useAppointmentStore.setState({ error: 'Test error' });
      });
      
      const storeWithError = useAppointmentStore.getState();
      expect(storeWithError.error).toBe('Test error');
      
      // Then clear it
      act(() => {
        storeWithError.clearError();
      });
      
      const clearedStore = useAppointmentStore.getState();
      expect(clearedStore.error).toBeNull();
    });

    it('should reset store correctly', () => {
      // First modify the store
      act(() => {
        useAppointmentStore.setState({
          appointments: [mockAppointment],
          selectedAppointment: mockAppointment,
          isLoading: true,
          error: 'Test error',
          totalAppointments: 1,
          currentPage: 2,
          filters: { status: 'confirmed' },
        });
      });
      
      const modifiedStore = useAppointmentStore.getState();
      expect(modifiedStore.appointments).toHaveLength(1);
      
      // Then reset it
      act(() => {
        modifiedStore.reset();
      });
      
      const resetStore = useAppointmentStore.getState();
      expect(resetStore.appointments).toEqual([]);
      expect(resetStore.selectedAppointment).toBeNull();
      expect(resetStore.isLoading).toBe(false);
      expect(resetStore.error).toBeNull();
      expect(resetStore.totalAppointments).toBe(0);
      expect(resetStore.currentPage).toBe(1);
      expect(resetStore.filters).toEqual({ status: 'all' });
    });
  });
});
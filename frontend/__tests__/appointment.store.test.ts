import { renderHook, act } from '@testing-library/react';
import { useAppointmentStore, useAppointments } from '@/stores/appointment.store';
import { appointmentApi } from '@/lib/api/appointment.api';
import { Appointment, CreateAppointmentRequest } from '@/shared/types';

// Mock API
jest.mock('@/lib/api/appointment.api');
const mockedAppointmentApi = appointmentApi as jest.Mocked<typeof appointmentApi>;

describe('AppointmentStore', () => {
  beforeEach(() => {
    // Reset store state
    useAppointmentStore.getState().reset();
    jest.clearAllMocks();
  });

  describe('fetchAppointments', () => {
    it('should fetch appointments successfully', async () => {
      // Arrange
      const mockAppointments: Appointment[] = [
        {
          id: '1',
          customerId: 'user1',
          carId: 'car1',
          serviceType: 'oil_change',
          scheduledAt: '2024-12-15T10:00:00Z',
          status: 'scheduled',
          createdAt: '2024-12-01T10:00:00Z',
          updatedAt: '2024-12-01T10:00:00Z',
        },
      ];

      mockedAppointmentApi.getAppointments.mockResolvedValueOnce(mockAppointments);

      const { result } = renderHook(() => useAppointments());

      // Act
      await act(async () => {
        await result.current.fetchAppointments();
      });

      // Assert
      expect(result.current.appointments).toEqual(mockAppointments);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(mockedAppointmentApi.getAppointments).toHaveBeenCalledTimes(1);
    });

    it('should handle fetch appointments error', async () => {
      // Arrange
      const errorMessage = 'Failed to fetch appointments';
      mockedAppointmentApi.getAppointments.mockRejectedValueOnce(new Error(errorMessage));

      const { result } = renderHook(() => useAppointments());

      // Act
      await act(async () => {
        await result.current.fetchAppointments();
      });

      // Assert
      expect(result.current.appointments).toEqual([]);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBe(errorMessage);
    });
  });

  describe('createAppointment', () => {
    it('should create appointment successfully', async () => {
      // Arrange
      const newAppointmentData: CreateAppointmentRequest = {
        carId: 'car1',
        serviceType: 'oil_change',
        scheduledAt: '2024-12-15T10:00:00Z',
        notes: 'Regular maintenance',
      };

      const createdAppointment: Appointment = {
        id: '1',
        customerId: 'user1',
        ...newAppointmentData,
        status: 'scheduled',
        createdAt: '2024-12-01T10:00:00Z',
        updatedAt: '2024-12-01T10:00:00Z',
      };

      mockedAppointmentApi.createAppointment.mockResolvedValueOnce(createdAppointment);

      const { result } = renderHook(() => useAppointments());

      // Act
      let success: boolean = false;
      await act(async () => {
        success = await result.current.createAppointment(newAppointmentData);
      });

      // Assert
      expect(success).toBe(true);
      expect(result.current.appointments).toContain(createdAppointment);
      expect(result.current.isCreating).toBe(false);
      expect(result.current.error).toBeNull();
      expect(mockedAppointmentApi.createAppointment).toHaveBeenCalledWith(newAppointmentData);
    });

    it('should handle create appointment error', async () => {
      // Arrange
      const newAppointmentData: CreateAppointmentRequest = {
        carId: 'car1',
        serviceType: 'oil_change',
        scheduledAt: '2024-12-15T10:00:00Z',
      };

      const errorMessage = 'Failed to create appointment';
      mockedAppointmentApi.createAppointment.mockRejectedValueOnce(new Error(errorMessage));

      const { result } = renderHook(() => useAppointments());

      // Act
      let success: boolean = true;
      await act(async () => {
        success = await result.current.createAppointment(newAppointmentData);
      });

      // Assert
      expect(success).toBe(false);
      expect(result.current.appointments).toEqual([]);
      expect(result.current.isCreating).toBe(false);
      expect(result.current.error).toBe(errorMessage);
    });
  });

  describe('cancelAppointment', () => {
    it('should cancel appointment successfully', async () => {
      // Arrange
      const appointment: Appointment = {
        id: '1',
        customerId: 'user1',
        carId: 'car1',
        serviceType: 'oil_change',
        scheduledAt: '2024-12-15T10:00:00Z',
        status: 'scheduled',
        createdAt: '2024-12-01T10:00:00Z',
        updatedAt: '2024-12-01T10:00:00Z',
      };

      const cancelledAppointment: Appointment = {
        ...appointment,
        status: 'cancelled',
        updatedAt: '2024-12-02T10:00:00Z',
      };

      mockedAppointmentApi.cancelAppointment.mockResolvedValueOnce(cancelledAppointment);

      const { result } = renderHook(() => useAppointments());

      // Set initial state
      act(() => {
        useAppointmentStore.setState({ appointments: [appointment] });
      });

      // Act
      let success: boolean = false;
      await act(async () => {
        success = await result.current.cancelAppointment('1');
      });

      // Assert
      expect(success).toBe(true);
      expect(result.current.appointments[0].status).toBe('cancelled');
      expect(mockedAppointmentApi.cancelAppointment).toHaveBeenCalledWith('1');
    });
  });

  describe('setFilters', () => {
    it('should update filters correctly', () => {
      // Arrange
      const { result } = renderHook(() => useAppointments());

      // Act
      act(() => {
        result.current.setFilters({ status: 'completed', carId: 'car1' });
      });

      // Assert
      expect(result.current.filters).toEqual({
        status: 'completed',
        carId: 'car1',
      });
    });
  });

  describe('clearError', () => {
    it('should clear error state', () => {
      // Arrange
      const { result } = renderHook(() => useAppointments());

      // Set error state
      act(() => {
        useAppointmentStore.setState({ error: 'Some error' });
      });

      expect(result.current.error).toBe('Some error');

      // Act
      act(() => {
        result.current.clearError();
      });

      // Assert
      expect(result.current.error).toBeNull();
    });
  });

  describe('reset', () => {
    it('should reset store to initial state', () => {
      // Arrange
      const { result } = renderHook(() => useAppointments());

      // Set some state
      act(() => {
        useAppointmentStore.setState({
          appointments: [
            {
              id: '1',
              customerId: 'user1',
              carId: 'car1',
              serviceType: 'oil_change',
              scheduledAt: '2024-12-15T10:00:00Z',
              status: 'scheduled',
              createdAt: '2024-12-01T10:00:00Z',
              updatedAt: '2024-12-01T10:00:00Z',
            },
          ],
          error: 'Some error',
          isLoading: true,
        });
      });

      // Act
      act(() => {
        result.current.reset();
      });

      // Assert
      expect(result.current.appointments).toEqual([]);
      expect(result.current.error).toBeNull();
      expect(result.current.isLoading).toBe(false);
      expect(result.current.filters).toEqual({ status: 'all' });
    });
  });
});
import { appointmentApi } from '@/lib/api/appointment.api';
import { apiClient } from '@/lib/api/client';
import { Appointment, CreateAppointmentRequest } from '@/shared/types';

// Mock API client
jest.mock('@/lib/api/client');
const mockedApiClient = apiClient as jest.Mocked<typeof apiClient>;

describe('appointmentApi', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('getAppointments', () => {
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

      mockedApiClient.get.mockResolvedValueOnce({ data: mockAppointments });

      // Act
      const result = await appointmentApi.getAppointments();

      // Assert
      expect(result).toEqual(mockAppointments);
      expect(mockedApiClient.get).toHaveBeenCalledWith('/appointments');
    });
  });

  describe('createAppointment', () => {
    it('should create appointment with camelCase properties', async () => {
      // Arrange
      const appointmentData: CreateAppointmentRequest = {
        carId: 'car1',           // ✅ camelCase per Agent.md
        serviceType: 'oil_change', // ✅ camelCase per Agent.md
        scheduledAt: '2024-12-15T10:00:00Z', // ✅ camelCase per Agent.md
        notes: 'Regular maintenance',
      };

      const createdAppointment: Appointment = {
        id: '1',
        customerId: 'user1',
        ...appointmentData,
        status: 'scheduled',
        createdAt: '2024-12-01T10:00:00Z',
        updatedAt: '2024-12-01T10:00:00Z',
      };

      mockedApiClient.post.mockResolvedValueOnce({ data: createdAppointment });

      // Act
      const result = await appointmentApi.createAppointment(appointmentData);

      // Assert
      expect(result).toEqual(createdAppointment);
      expect(mockedApiClient.post).toHaveBeenCalledWith('/appointments', {
        carId: 'car1',           // ✅ camelCase per Agent.md
        serviceType: 'oil_change', // ✅ camelCase per Agent.md
        scheduledAt: '2024-12-15T10:00:00Z', // ✅ camelCase per Agent.md
        notes: 'Regular maintenance',
      });
    });
  });

  describe('updateAppointment', () => {
    it('should update appointment with camelCase properties', async () => {
      // Arrange
      const updateData = {
        serviceType: 'brake_service', // ✅ camelCase per Agent.md
        scheduledAt: '2024-12-16T10:00:00Z', // ✅ camelCase per Agent.md
        notes: 'Updated notes',
      };

      const updatedAppointment: Appointment = {
        id: '1',
        customerId: 'user1',
        carId: 'car1',
        status: 'scheduled',
        createdAt: '2024-12-01T10:00:00Z',
        updatedAt: '2024-12-02T10:00:00Z',
        ...updateData,
      };

      mockedApiClient.put.mockResolvedValueOnce({ data: updatedAppointment });

      // Act
      const result = await appointmentApi.updateAppointment('1', updateData);

      // Assert
      expect(result).toEqual(updatedAppointment);
      expect(mockedApiClient.put).toHaveBeenCalledWith('/appointments/1', {
        serviceType: 'brake_service', // ✅ camelCase per Agent.md
        scheduledAt: '2024-12-16T10:00:00Z', // ✅ camelCase per Agent.md
        notes: 'Updated notes',
      });
    });
  });

  describe('cancelAppointment', () => {
    it('should cancel appointment successfully', async () => {
      // Arrange
      const cancelledAppointment: Appointment = {
        id: '1',
        customerId: 'user1',
        carId: 'car1',
        serviceType: 'oil_change',
        scheduledAt: '2024-12-15T10:00:00Z',
        status: 'cancelled',
        createdAt: '2024-12-01T10:00:00Z',
        updatedAt: '2024-12-02T10:00:00Z',
      };

      mockedApiClient.patch.mockResolvedValueOnce({ data: cancelledAppointment });

      // Act
      const result = await appointmentApi.cancelAppointment('1');

      // Assert
      expect(result).toEqual(cancelledAppointment);
      expect(mockedApiClient.patch).toHaveBeenCalledWith('/appointments/1/cancel');
    });
  });

  describe('deleteAppointment', () => {
    it('should delete appointment successfully', async () => {
      // Arrange
      mockedApiClient.delete.mockResolvedValueOnce({});

      // Act
      await appointmentApi.deleteAppointment('1');

      // Assert
      expect(mockedApiClient.delete).toHaveBeenCalledWith('/appointments/1');
    });
  });
});
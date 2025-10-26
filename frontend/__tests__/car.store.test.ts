// ✅ CarStore Tests - Comprehensive testing following Agent.md patterns

import { act, renderHook } from '@testing-library/react';
import { useCars, useCarStore } from '@/stores/car.store';
import { vehicleService } from '@/lib/services/vehicle.service';
import type { Car, CreateCarRequest } from '@/types/car';

// ✅ Mock vehicle service
jest.mock('@/lib/services/vehicle.service', () => ({
  vehicleService: {
    getVehicles: jest.fn(),
    getVehicle: jest.fn(),
    createVehicle: jest.fn(),
    updateVehicle: jest.fn(),
    deleteVehicle: jest.fn(),
  }
}));

const mockVehicleService = vehicleService as jest.Mocked<typeof vehicleService>;

// ✅ Sample test data
const mockCar: Car = {
  id: '1',
  make: 'Toyota',
  model: 'Camry',
  year: 2023,
  licensePlate: 'ABC123',
  vin: '1234567890',
  color: 'Blue',
  mileage: 15000,
  ownerId: 'owner-1',
  createdAt: '2025-01-01T00:00:00Z',
  updatedAt: '2025-01-01T00:00:00Z',
};

const mockCreateCarData: CreateCarRequest = {
  make: 'Honda',
  model: 'Civic',
  year: 2024,
  licensePlate: 'XYZ789',
  vin: '0987654321',
  color: 'Red',
  mileage: 5000,
};

describe('CarStore', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    // Reset store state before each test
    useCarStore.getState().reset();
  });

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      const { result } = renderHook(() => useCars());
      
      expect(result.current.cars).toEqual([]);
      expect(result.current.selectedCar).toBeNull();
      expect(result.current.isLoading).toBe(false);
      expect(result.current.isCreating).toBe(false);
      expect(result.current.isUpdating).toBe(false);
      expect(result.current.isDeleting).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.totalCars).toBe(0);
      expect(result.current.currentPage).toBe(1);
      expect(result.current.pageSize).toBe(10);
      expect(result.current.filters).toEqual({});
    });
  });

  describe('Fetch Cars', () => {
    it('should fetch cars successfully', async () => {
      const mockResponse = {
        success: true,
        data: {
          vehicles: [
            {
              id: '1',
              make: 'Toyota',
              model: 'Camry',
              year: 2023,
              licensePlate: 'ABC123',
              vin: '1234567890',
              color: 'Blue',
              mileage: 15000,
              ownerId: 'owner-1',
              createdAt: '2025-01-01T00:00:00Z',
              updatedAt: '2025-01-01T00:00:00Z',
            }
          ],
          total: 1,
          limit: 10,
          offset: 0,
        }
      };
      
      mockVehicleService.getVehicles.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      await act(async () => {
        await result.current.fetchCars();
      });
      
      expect(result.current.cars).toHaveLength(1);
      expect(result.current.cars[0].make).toBe('Toyota');
      expect(result.current.totalCars).toBe(1);
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });

    it('should handle fetch cars failure', async () => {
      const mockResponse = {
        success: false,
        error: {
          message: 'Network error',
          status: 500,
          code: 'FETCH_ERROR'
        }
      };
      
      mockVehicleService.getVehicles.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      await act(async () => {
        await result.current.fetchCars();
      });
      
      expect(result.current.cars).toEqual([]);
      expect(result.current.error).toBe('Network error');
      expect(result.current.isLoading).toBe(false);
    });
  });

  describe('Fetch Car By ID', () => {
    it('should fetch single car successfully', async () => {
      const mockResponse = {
        success: true,
        data: {
          id: '1',
          make: 'Toyota',
          model: 'Camry',
          year: 2023,
          licensePlate: 'ABC123',
          vin: '1234567890',
          color: 'Blue',
          mileage: 15000,
          ownerId: 'owner-1',
          createdAt: '2025-01-01T00:00:00Z',
          updatedAt: '2025-01-01T00:00:00Z',
        }
      };
      
      mockVehicleService.getVehicle.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      await act(async () => {
        await result.current.fetchCarById('1');
      });
      
      expect(result.current.selectedCar).toBeDefined();
      expect(result.current.selectedCar?.make).toBe('Toyota');
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
    });
  });

  describe('Create Car', () => {
    it('should create car successfully', async () => {
      const mockResponse = {
        success: true,
        data: {
          id: '2',
          ...mockCreateCarData,
          ownerId: 'owner-1',
          createdAt: '2025-01-01T00:00:00Z',
          updatedAt: '2025-01-01T00:00:00Z',
        }
      };
      
      mockVehicleService.createVehicle.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      let createResult: boolean;
      await act(async () => {
        createResult = await result.current.createCar(mockCreateCarData);
      });
      
      expect(createResult!).toBe(true);
      expect(result.current.cars).toHaveLength(1);
      expect(result.current.cars[0].make).toBe('Honda');
      expect(result.current.totalCars).toBe(1);
      expect(result.current.isCreating).toBe(false);
    });

    it('should handle create car failure', async () => {
      const mockResponse = {
        success: false,
        error: {
          message: 'Validation error',
          status: 400,
          code: 'VALIDATION_ERROR'
        }
      };
      
      mockVehicleService.createVehicle.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      let createResult: boolean;
      await act(async () => {
        createResult = await result.current.createCar(mockCreateCarData);
      });
      
      expect(createResult!).toBe(false);
      expect(result.current.error).toBe('Validation error');
      expect(result.current.isCreating).toBe(false);
    });
  });

  describe('Update Car', () => {
    it('should update car successfully', async () => {
      // First add a car to the store
      const initialState = useCarStore.getState();
      act(() => {
        initialState.selectCar(mockCar);
        useCarStore.setState({ cars: [mockCar] });
      });

      const mockResponse = {
        success: true,
        data: {
          ...mockCar,
          make: 'Honda',
          updatedAt: '2025-01-02T00:00:00Z',
        }
      };
      
      mockVehicleService.updateVehicle.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      let updateResult: boolean;
      await act(async () => {
        updateResult = await result.current.updateCar('1', { make: 'Honda' });
      });
      
      expect(updateResult!).toBe(true);
      expect(result.current.cars[0].make).toBe('Honda');
      expect(result.current.isUpdating).toBe(false);
    });
  });

  describe('Delete Car', () => {
    it('should delete car successfully', async () => {
      // First add a car to the store
      act(() => {
        useCarStore.setState({ cars: [mockCar], totalCars: 1 });
      });

      const mockResponse = {
        success: true,
        data: { message: 'Car deleted successfully' }
      };
      
      mockVehicleService.deleteVehicle.mockResolvedValue(mockResponse);
      
      const { result } = renderHook(() => useCars());
      
      let deleteResult: boolean;
      await act(async () => {
        deleteResult = await result.current.deleteCar('1');
      });
      
      expect(deleteResult!).toBe(true);
      expect(result.current.cars).toHaveLength(0);
      expect(result.current.totalCars).toBe(0);
      expect(result.current.isDeleting).toBe(false);
    });
  });

  describe('State Management', () => {
    it('should select car correctly', () => {
      const { result } = renderHook(() => useCars());
      
      act(() => {
        result.current.selectCar(mockCar);
      });
      
      expect(result.current.selectedCar).toEqual(mockCar);
    });

    it('should set filters correctly', () => {
      const { result } = renderHook(() => useCars());
      
      act(() => {
        result.current.setFilters({ make: 'Toyota', year: 2023 });
      });
      
      expect(result.current.filters).toEqual({ make: 'Toyota', year: 2023 });
      expect(result.current.currentPage).toBe(1); // Should reset page
    });

    it('should set page correctly', () => {
      const { result } = renderHook(() => useCars());
      
      act(() => {
        result.current.setPage(3);
      });
      
      expect(result.current.currentPage).toBe(3);
    });

    it('should clear error correctly', () => {
      const { result } = renderHook(() => useCars());
      
      // First set an error
      act(() => {
        useCarStore.setState({ error: 'Test error' });
      });
      
      expect(result.current.error).toBe('Test error');
      
      // Then clear it
      act(() => {
        result.current.clearError();
      });
      
      expect(result.current.error).toBeNull();
    });

    it('should reset store correctly', () => {
      const { result } = renderHook(() => useCars());
      
      // First modify the store
      act(() => {
        useCarStore.setState({
          cars: [mockCar],
          selectedCar: mockCar,
          isLoading: true,
          error: 'Test error',
          totalCars: 1,
          currentPage: 2,
          filters: { make: 'Toyota' },
        });
      });
      
      expect(result.current.cars).toHaveLength(1);
      
      // Then reset it
      act(() => {
        result.current.reset();
      });
      
      expect(result.current.cars).toEqual([]);
      expect(result.current.selectedCar).toBeNull();
      expect(result.current.isLoading).toBe(false);
      expect(result.current.error).toBeNull();
      expect(result.current.totalCars).toBe(0);
      expect(result.current.currentPage).toBe(1);
      expect(result.current.filters).toEqual({});
    });
  });
});
// ✅ Car Store - Zustand + Immer + TypeScript following Agent.md patterns
// Provides centralized car state management with CRUD operations

import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import type { Car, CreateCarRequest } from '@/types/car';
import { carService, type Car as CarServiceType } from '@/lib/services/car.service';

// ✅ Helper function to convert CarServiceType to Car (same interface, just ensuring compatibility)
const serviceCarToCar = (car: CarServiceType): Car => ({
  id: car.id,
  make: car.make,
  model: car.model,
  year: car.year,
  licensePlate: car.licensePlate,
  vin: car.vin || '',
  color: car.color || '',
  mileage: car.mileage || 0,
  ownerId: car.ownerId,
  createdAt: car.createdAt,
  updatedAt: car.updatedAt,
});

// ✅ Car state interface following Agent.md conventions
interface CarState {
  // Data state
  cars: Car[];
  selectedCar: Car | null;
  
  // Loading and error states
  isLoading: boolean;
  isCreating: boolean;
  isUpdating: boolean;
  isDeleting: boolean;
  error: string | null;
  
  // Pagination and filters
  totalCars: number;
  currentPage: number;
  pageSize: number;
  filters: {
    make?: string;
    model?: string;
    year?: number;
    ownerId?: string;
  };
  
  // Actions
  fetchCars: (filters?: CarState['filters']) => Promise<void>;
  fetchCarById: (id: string) => Promise<void>;
  createCar: (carData: CreateCarRequest) => Promise<boolean>;
  updateCar: (id: string, carData: Partial<CreateCarRequest>) => Promise<boolean>;
  deleteCar: (id: string) => Promise<boolean>;
  selectCar: (car: Car | null) => void;
  setFilters: (filters: Partial<CarState['filters']>) => void;
  setPage: (page: number) => void;
  clearError: () => void;
  reset: () => void;
}

// ✅ Initial state following Agent.md conventions
const initialState = {
  cars: [],
  selectedCar: null,
  isLoading: false,
  isCreating: false,
  isUpdating: false,
  isDeleting: false,
  error: null,
  totalCars: 0,
  currentPage: 1,
  pageSize: 10,
  filters: {},
};

// ✅ Car store implementation with Zustand + Immer
export const useCarStore = create<CarState>()(
  immer((set, get) => ({
    // Initial state
    ...initialState,
    
    // ✅ Fetch cars with optional filters
    fetchCars: async (filters) => {
      set((state) => {
        state.isLoading = true;
        state.error = null;
        if (filters) {
          state.filters = { ...state.filters, ...filters };
        }
      });
      
      try {
        const response = await carService.getCars({
          ...get().filters,
          limit: get().pageSize,
          offset: (get().currentPage - 1) * get().pageSize,
        });
        
        if (response.data) {
          set((state) => {
            state.cars = response.data!.cars.map(serviceCarToCar);
            state.totalCars = response.data!.total;
            state.isLoading = false;
          });
        } else {
          throw new Error(response.error?.message || 'Failed to fetch cars');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to fetch cars';
          state.isLoading = false;
        });
      }
    },
    
    // ✅ Fetch single car by ID
    fetchCarById: async (id: string) => {
      set((state) => {
        state.isLoading = true;
        state.error = null;
      });
      
      try {
        const response = await carService.getCar(id);
        
        if (response.data) {
          set((state) => {
            state.selectedCar = serviceCarToCar(response.data!);
            state.isLoading = false;
          });
        } else {
          throw new Error(response.error?.message || 'Car not found');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to fetch car';
          state.isLoading = false;
        });
      }
    },
    
    // ✅ Create new car
    createCar: async (carData: CreateCarRequest): Promise<boolean> => {
      set((state) => {
        state.isCreating = true;
        state.error = null;
      });
      
      try {
        const response = await carService.createCar(carData);
        
        if (response.data) {
          set((state) => {
            state.cars.unshift(serviceCarToCar(response.data!)); // Add to beginning of list
            state.totalCars += 1;
            state.isCreating = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to create car');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to create car';
          state.isCreating = false;
        });
        return false;
      }
    },
    
    // ✅ Update existing car
    updateCar: async (id: string, carData: Partial<CreateCarRequest>): Promise<boolean> => {
      set((state) => {
        state.isUpdating = true;
        state.error = null;
      });
      
      try {
        const response = await carService.updateCar(id, carData);
        
        if (response.data) {
          set((state) => {
            const index = state.cars.findIndex(car => car.id === id);
            if (index !== -1) {
              state.cars[index] = serviceCarToCar(response.data!);
            }
            if (state.selectedCar?.id === id) {
              state.selectedCar = serviceCarToCar(response.data!);
            }
            state.isUpdating = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to update car');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to update car';
          state.isUpdating = false;
        });
        return false;
      }
    },
    
    // ✅ Delete car
    deleteCar: async (id: string): Promise<boolean> => {
      set((state) => {
        state.isDeleting = true;
        state.error = null;
      });
      
      try {
        const response = await carService.deleteCar(id);
        
        if (response.success) {
          set((state) => {
            state.cars = state.cars.filter(car => car.id !== id);
            state.totalCars -= 1;
            if (state.selectedCar?.id === id) {
              state.selectedCar = null;
            }
            state.isDeleting = false;
          });
          return true;
        } else {
          throw new Error(response.error?.message || 'Failed to delete car');
        }
      } catch (error) {
        set((state) => {
          state.error = error instanceof Error ? error.message : 'Failed to delete car';
          state.isDeleting = false;
        });
        return false;
      }
    },
    
    // ✅ Select car for detailed view
    selectCar: (car: Car | null) => {
      set((state) => {
        state.selectedCar = car;
      });
    },
    
    // ✅ Update filters
    setFilters: (filters: Partial<CarState['filters']>) => {
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
export const useCars = () => {
  const store = useCarStore();
  return {
    // Data
    cars: store.cars,
    selectedCar: store.selectedCar,
    totalCars: store.totalCars,
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
    fetchCars: store.fetchCars,
    fetchCarById: store.fetchCarById,
    createCar: store.createCar,
    updateCar: store.updateCar,
    deleteCar: store.deleteCar,
    selectCar: store.selectCar,
    setFilters: store.setFilters,
    setPage: store.setPage,
    clearError: store.clearError,
    reset: store.reset,
  };
};

// ✅ Specific hooks for common use cases
export const useCarList = () => {
  const { cars, isLoading, error, fetchCars, setFilters, setPage } = useCars();
  return { cars, isLoading, error, fetchCars, setFilters, setPage };
};

export const useCarDetails = () => {
  const { selectedCar, isLoading, error, fetchCarById, selectCar } = useCars();
  return { selectedCar, isLoading, error, fetchCarById, selectCar };
};

export const useCarMutations = () => {
  const { 
    isCreating, 
    isUpdating, 
    isDeleting, 
    createCar, 
    updateCar, 
    deleteCar,
    error,
    clearError
  } = useCars();
  return { 
    isCreating, 
    isUpdating, 
    isDeleting, 
    createCar, 
    updateCar, 
    deleteCar,
    error,
    clearError
  };
};
import { useState, useEffect, useCallback } from 'react';
import { Car, CreateCarRequest } from '@/types/car';
import { carApi } from '@/lib/api/carApi';

export interface UseCarsReturn {
  cars: Car[];
  loading: boolean;
  error: string | null;
  createCar: (carData: CreateCarRequest) => Promise<boolean>;
  updateCar: (id: string, carData: Partial<CreateCarRequest>) => Promise<boolean>;
  deleteCar: (id: string) => Promise<boolean>;
  refreshCars: () => Promise<void>;
}

// Custom hook following Agent.md clean architecture principles
export function useCars(): UseCarsReturn {
  const [cars, setCars] = useState<Car[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Fetch cars - following Agent.md error handling
  const fetchCars = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);

      const { data, error: apiError } = await carApi.getCars();
      
      if (apiError) {
        setError(apiError.message);
        return;
      }

      setCars(data || []);
    } catch (err) {
      setError('Failed to fetch cars');
    } finally {
      setLoading(false);
    }
  }, []);

  // Create car - following Agent.md TDD patterns
  const createCar = useCallback(async (carData: CreateCarRequest): Promise<boolean> => {
    try {
      const { data, error: apiError } = await carApi.createCar(carData);
      
      if (apiError) {
        setError(apiError.message);
        return false;
      }

      if (data) {
        setCars(prevCars => [...prevCars, data]);
        return true;
      }

      return false;
    } catch (err) {
      setError('Failed to create car');
      return false;
    }
  }, []);

  // Update car - following Agent.md naming conventions
  const updateCar = useCallback(async (
    id: string, 
    carData: Partial<CreateCarRequest>
  ): Promise<boolean> => {
    try {
      const { data, error: apiError } = await carApi.updateCar(id, carData);
      
      if (apiError) {
        setError(apiError.message);
        return false;
      }

      if (data) {
        setCars(prevCars => 
          prevCars.map(car => car.id === id ? data : car)
        );
        return true;
      }

      return false;
    } catch (err) {
      setError('Failed to update car');
      return false;
    }
  }, []);

  // Delete car - following Agent.md clean architecture
  const deleteCar = useCallback(async (id: string): Promise<boolean> => {
    try {
      const { error: apiError } = await carApi.deleteCar(id);
      
      if (apiError) {
        setError(apiError.message);
        return false;
      }

      setCars(prevCars => prevCars.filter(car => car.id !== id));
      return true;
    } catch (err) {
      setError('Failed to delete car');
      return false;
    }
  }, []);

  // Refresh cars - utility function
  const refreshCars = useCallback(async () => {
    await fetchCars();
  }, [fetchCars]);

  useEffect(() => {
    fetchCars();
  }, [fetchCars]);

  return {
    cars,
    loading,
    error,
    createCar,
    updateCar,
    deleteCar,
    refreshCars,
  };
}
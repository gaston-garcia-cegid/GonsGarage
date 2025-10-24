// src/components/client/ClientCars.tsx
'use client';

import React, { useState, useCallback } from 'react';
import { Car, CreateCarRequest } from '@/types/car';
import { useCars } from '@/hooks/useCars';
import CarList from '@/app/cars/components/CarList';
import CarModal from '@/app/cars/components/CarModal';
import LoadingSpinner from '@/components/ui/Loading/LoadingSpinner';
import ErrorAlert from '@/components/ui/Error/ErrorAlert';
import { useRouter } from 'next/navigation';
import styles from '../client.module.css';

interface ClientCarsProps {
  onAddCar: (car: Car) => void;
  onUpdateCar: (cars: Car[]) => void;
  showAddButton: boolean;
  maxCars?: number;
}

// Client cars component following Agent.md component patterns
export default function ClientCars({ 
  onAddCar, 
  onUpdateCar, 
  showAddButton = true, 
  maxCars }: ClientCarsProps) {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCar, setEditingCar] = useState<Car | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  
  const router = useRouter();
  const { cars, loading, error, createCar, updateCar, deleteCar, refreshCars } = useCars();

  // Enhanced create car handler - following Agent.md error handling
  const handleCreateCar = useCallback(async (carData: CreateCarRequest): Promise<boolean> => {
    // Check max cars limit if specified
    if (maxCars && cars.length >= maxCars) {
      alert(`You can only register up to ${maxCars} cars.`);
      return false;
    }

    setIsCreating(true);

    try {
      const success = await createCar(carData);
      
      if (success) {
        // Refresh cars to get the latest data with the new car
        await refreshCars();
        
        // Get the newly created car (assuming it's the last one after refresh)
        const updatedCars = await refreshCars();
        const newCar = cars.find(car => 
          car.make === carData.make && 
          car.model === carData.model && 
          car.licensePlate === carData.licensePlate
        );

        // Notify parent component about the new car - following Agent.md callback patterns
        if (onAddCar && newCar) {
          onAddCar(newCar);
        }

        // Notify about cars list update
        if (onUpdateCar) {
          onUpdateCar(cars);
        }

        setShowCreateModal(false);
        return true;
      }

      return false;
    } catch (error) {
      console.error('Failed to create car:', error);
      return false;
    } finally {
      setIsCreating(false);
    }
  }, [createCar, refreshCars, cars, maxCars, onAddCar, onUpdateCar]);

  // Enhanced update car handler - following Agent.md consistency
  const handleUpdateCar = useCallback(async (id: string, carData: Partial<CreateCarRequest>): Promise<boolean> => {
    try {
      const success = await updateCar(id, carData);
      
      if (success) {
        // Refresh cars to get updated data
        await refreshCars();
        
        // Notify about cars list update
        if (onUpdateCar) {
          onUpdateCar(cars);
        }

        setEditingCar(null);
        return true;
      }

      return false;
    } catch (error) {
      console.error('Failed to update car:', error);
      return false;
    }
  }, [updateCar, refreshCars, cars, onUpdateCar]);

  // Handle car operations following Agent.md clean patterns
  const handleDeleteCar = useCallback(async (id: string) => {
    if (!confirm('Are you sure you want to delete this car?')) {
      return;
    }

    try {
      const success = await deleteCar(id);
      
      if (success) {
        // Notify about cars list update
        if (onUpdateCar) {
          onUpdateCar(cars.filter(car => car.id !== id));
        }
      } else {
        alert('Failed to delete car. Please try again.');
      }
    } catch (error) {
      console.error('Failed to delete car:', error);
      alert('An error occurred while deleting the car.');
    }
  }, [deleteCar, cars, onUpdateCar]);

  const handleEditCar = (car: Car) => {
    setEditingCar(car);
  };

  const handleModalClose = () => {
    setShowCreateModal(false);
    setEditingCar(null);
  };

  // Enhanced add car button handler - following Agent.md UX patterns
  const handleAddCarClick = useCallback(() => {
    // Check if user has reached max cars limit
    if (maxCars && cars.length >= maxCars) {
      alert(`You have reached the maximum limit of ${maxCars} cars.`);
      return;
    }

    setShowCreateModal(true);
  }, [cars.length, maxCars]);

  const handleAddCar = () => {
    setShowCreateModal(true);
    // Removed onAddCar() call since car is not available here
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <span>Loading cars...</span>
      </div>
    );
  }

  // Check if user can add more cars
  const canAddCars = !maxCars || cars.length < maxCars;

  return (
    <div className={styles.clientCarsContainer}>
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <h2>Your Vehicles</h2>
          <p>
            {maxCars 
              ? `Manage your registered cars (${cars.length}/${maxCars})`
              : 'Manage your registered cars'
            }
          </p>
        </div>
        <button 
          onClick={handleAddCar}
          className={styles.primaryButton}
          disabled={isCreating}
        >
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          {isCreating ? 'Adding...' : 'Add New Car'}
        </button>

        {/* Max cars reached message */}
        {maxCars && cars.length >= maxCars && (
          <div className={styles.maxCarsMessage}>
            <span>Maximum cars limit reached</span>
          </div>
        )}
      </div>

      {error && <ErrorAlert message={error} />}

      <CarList
        cars={cars}
        onEdit={handleEditCar}
        onDelete={handleDeleteCar}
        onViewDetails={(id) => router.push(`/cars/${id}`)}
        onScheduleService={(id) => router.push(`/appointments/new?carId=${id}`)}
      />

      {/* Car Modal */}
      {(showCreateModal || editingCar) && (
        <CarModal
          car={editingCar}
          onClose={handleModalClose}
          onCreate={handleCreateCar}
          onUpdate={handleUpdateCar}
        />
      )}
    </div>
  );
}
'use client';

import React, { useState, useCallback } from 'react';
import { Car, CreateCarRequest } from '@/types/car';
import { useCars } from '@/stores';
import CarList from '@/app/cars/components/CarList';
import CarModal from '@/app/cars/components/CarModal';
import LoadingSpinner from '@/components/ui/Loading/LoadingSpinner';
import ErrorAlert from '@/components/ui/Error/ErrorAlert';
import ConfirmModal from '@/components/ui/Modal/ConfirmModal';
import EmptyCarState from '@/components/empty-states/EmptyCarState';
import { useRouter } from 'next/navigation';
import styles from './CarsContainer.module.css'; // ✅ Import styles

interface CarsContainerProps {
  // Optional callbacks for parent components
  onAddCar?: (car: Car) => void;
  onUpdateCar?: (cars: Car[]) => void;
  onDeleteCar?: (carId: string) => void;
  
  // UI customization
  maxCars?: number;
  showHeader?: boolean;
  headerTitle?: string;
  headerSubtitle?: string;
  addButtonText?: string;
  className?: string;
  
  // Layout control
  renderHeader?: () => React.ReactNode;
  renderEmptyState?: () => React.ReactNode;
}

export default function CarsContainer({
  onAddCar,
  onUpdateCar,
  onDeleteCar,
  maxCars,
  showHeader = true,
  headerTitle = 'Your Cars',
  headerSubtitle,
  addButtonText = 'Add New Car',
  className = '',
  renderHeader,
  renderEmptyState,
}: CarsContainerProps) {
  const router = useRouter();
  const { cars, isLoading, error, createCar, updateCar, deleteCar, fetchCars } = useCars();

  // Modal states
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCar, setEditingCar] = useState<Car | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  
  // Delete confirmation state
  const [deleteConfirmation, setDeleteConfirmation] = useState<{
    isOpen: boolean;
    carId: string | null;
    carName: string;
  }>({
    isOpen: false,
    carId: null,
    carName: '',
  });

  // ✅ Create car handler
  const handleCreateCar = useCallback(async (carData: CreateCarRequest): Promise<boolean> => {
    if (maxCars && cars.length >= maxCars) {
      alert(`You can only register up to ${maxCars} cars.`);
      return false;
    }

    setIsCreating(true);

    try {
      const success = await createCar(carData);
      
      if (success) {
        await fetchCars();
        
        const newCar = cars.find(car => 
          car.make === carData.make && 
          car.model === carData.model && 
          car.licensePlate === carData.licensePlate
        );

        if (onAddCar && newCar) {
          onAddCar(newCar);
        }

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
  }, [createCar, fetchCars, cars, maxCars, onAddCar, onUpdateCar]);

  // ✅ Update car handler
  const handleUpdateCar = useCallback(async (id: string, carData: Partial<CreateCarRequest>): Promise<boolean> => {
    try {
      const success = await updateCar(id, carData);
      
      if (success) {
        await fetchCars();
        
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
  }, [updateCar, fetchCars, cars, onUpdateCar]);

  // ✅ Delete car handler - opens confirmation modal
  const handleDeleteCar = useCallback((id: string) => {
    const car = cars.find(c => c.id === id);
    if (!car) return;

    setDeleteConfirmation({
      isOpen: true,
      carId: id,
      carName: `${car.make} ${car.model} (${car.licensePlate})`,
    });
  }, [cars]);

  // ✅ Confirm deletion
  const confirmDelete = useCallback(async () => {
    const { carId } = deleteConfirmation;
    if (!carId) return;

    try {
      const success = await deleteCar(carId);
      
      if (success) {
        if (onDeleteCar) {
          onDeleteCar(carId);
        }

        if (onUpdateCar) {
          onUpdateCar(cars.filter(car => car.id !== carId));
        }
      } else {
        alert('Failed to delete car. Please try again.');
      }
    } catch (error) {
      console.error('Failed to delete car:', error);
      alert('An error occurred while deleting the car.');
    } finally {
      setDeleteConfirmation({ isOpen: false, carId: null, carName: '' });
    }
  }, [deleteConfirmation, deleteCar, cars, onDeleteCar, onUpdateCar]);

  // ✅ Cancel deletion
  const cancelDelete = useCallback(() => {
    setDeleteConfirmation({ isOpen: false, carId: null, carName: '' });
  }, []);

  // ✅ Edit car handler
  const handleEditCar = useCallback((car: Car) => {
    setEditingCar(car);
  }, []);

  // ✅ Close modal
  const handleModalClose = useCallback(() => {
    setShowCreateModal(false);
    setEditingCar(null);
  }, []);

  // ✅ Open create modal
  const handleAddCarClick = useCallback(() => {
    if (maxCars && cars.length >= maxCars) {
      alert(`You have reached the maximum limit of ${maxCars} cars.`);
      return;
    }
    setShowCreateModal(true);
  }, [cars.length, maxCars]);

  // ✅ Loading state
  if (isLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <span>Loading cars...</span>
      </div>
    );
  }

  // ✅ Empty state
  if (!isLoading && cars.length === 0) {
    return (
      <div className={styles.emptyStateWrapper}>
        {renderEmptyState ? renderEmptyState() : <EmptyCarState />}
      </div>
    );
  }

  const canAddCars = !maxCars || cars.length < maxCars;
  const computedSubtitle = headerSubtitle || (
    maxCars 
      ? `Manage your registered cars (${cars.length}/${maxCars})`
      : 'Manage your registered cars'
  );

  return (
    <div className={`${styles.container} ${className}`}>
      {/* ✅ Custom header or default */}
      {renderHeader ? (
        renderHeader()
      ) : showHeader ? (
        <div className={styles.header}>
          <div className={styles.headerContent}>
            <h2>{headerTitle}</h2>
            <p>{computedSubtitle}</p>
          </div>
          <button 
            onClick={handleAddCarClick}
            className={styles.addButton}
            disabled={isCreating || !canAddCars}
          >
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} 
                d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            {isCreating ? 'Adding...' : addButtonText}
          </button>
        </div>
      ) : null}

      {/* ✅ Max cars warning */}
      {maxCars && cars.length >= maxCars && (
        <div className={styles.maxCarsMessage}>
          <span>⚠️ Maximum cars limit reached ({cars.length}/{maxCars})</span>
        </div>
      )}

      {/* ✅ Error display */}
      {error && <ErrorAlert message={error} />}

      {/* ✅ Car list */}
      <CarList
        cars={cars}
        onEdit={handleEditCar}
        onDelete={handleDeleteCar}
        onViewDetails={(id) => router.push(`/cars/${id}`)}
        onScheduleService={(id) => router.push(`/appointments/new?carId=${id}`)}
      />

      {/* ✅ Car Modal */}
      {(showCreateModal || editingCar) && (
        <CarModal
          car={editingCar}
          onClose={handleModalClose}
          onCreate={handleCreateCar}
          onUpdate={handleUpdateCar}
        />
      )}

      {/* ✅ Confirmation Modal */}
      <ConfirmModal
        isOpen={deleteConfirmation.isOpen}
        title="Delete Car"
        message={`Are you sure you want to delete ${deleteConfirmation.carName}? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        variant="danger"
        onConfirm={confirmDelete}
        onCancel={cancelDelete}
      />
    </div>
  );
}
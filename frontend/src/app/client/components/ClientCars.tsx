// src/components/client/ClientCars.tsx
'use client';

import React, { useState } from 'react';
import { Car } from '@/types/car';
import { useCars } from '@/hooks/useCars';
import CarList from '@/app/cars/components/CarList';
import CarModal from '@/app/cars/components/CarModal';
import LoadingSpinner from '@/components/ui/Loading/LoadingSpinner';
import ErrorAlert from '@/components/ui/Error/ErrorAlert';
import { useRouter } from 'next/navigation';
import styles from '../client.module.css';

interface ClientCarsProps {
  onAddCar: () => void;
}

// Client cars component following Agent.md component patterns
export default function ClientCars({ onAddCar }: ClientCarsProps) {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCar, setEditingCar] = useState<Car | null>(null);
  
  const router = useRouter();
  const { cars, loading, error, createCar, updateCar, deleteCar } = useCars();

  // Handle car operations following Agent.md clean patterns
  const handleDeleteCar = async (id: string) => {
    if (!confirm('Are you sure you want to delete this car?')) {
      return;
    }
    await deleteCar(id);
  };

  const handleEditCar = (car: Car) => {
    setEditingCar(car);
  };

  const handleModalClose = () => {
    setShowCreateModal(false);
    setEditingCar(null);
  };

  const handleAddCar = () => {
    setShowCreateModal(true);
    onAddCar(); // Notify parent component
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <span>Loading cars...</span>
      </div>
    );
  }

  return (
    <div className={styles.clientCarsContainer}>
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <h2>Your Vehicles</h2>
          <p>Manage your registered cars</p>
        </div>
        <button 
          onClick={handleAddCar}
          className={styles.primaryButton}
        >
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Add New Car
        </button>
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
          onCreate={createCar}
          onUpdate={updateCar}
        />
      )}
    </div>
  );
}
'use client';

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import { Car } from '@/types/car';
import { useCars } from '@/stores';
import CarList from './components/CarList';
import CarModal from './components/CarModal';
import LoadingSpinner from '@/components/ui/Loading/LoadingSpinner';
import ErrorAlert from '@/components/ui/Error/ErrorAlert';
import Image from 'next/image';
import styles from './cars.module.css';

// Main Cars page component following Agent.md component conventions
export default function CarsPage() {
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCar, setEditingCar] = useState<Car | null>(null);

  const { user, logout } = useAuth();
  const router = useRouter();
  const { cars, isLoading: loading, error, createCar, updateCar, deleteCar, fetchCars } = useCars();

  // Redirect if not authenticated and fetch cars - following Agent.md security practices
  React.useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    fetchCars();
  }, [user, router, fetchCars]);

  // Handle car deletion - following Agent.md user experience
  const handleDeleteCar = async (id: string) => {
    if (!confirm('Are you sure you want to delete this car?')) {
      return;
    }

    const success = await deleteCar(id);
    if (!success) {
      alert('Failed to delete car. Please try again.');
    }
  };

  // Handle edit car - following Agent.md clean patterns
  const handleEditCar = (car: Car) => {
    setEditingCar(car);
  };

  // Handle modal close - following Agent.md state management
  const handleModalClose = () => {
    setShowCreateModal(false);
    setEditingCar(null);
  };

  // Loading state - following Agent.md user experience
  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <span>Loading cars...</span>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header - following Agent.md UI conventions */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logoSection} onClick={() => router.push('/')}>
            <div className={styles.logoIcon}>
              <Image
                src="/images/LogoGonsGarage.jpg"
                alt="GonsGarage Logo"
                width={24}
                height={24}
                style={{ objectFit: 'contain' }}
              />
            </div>
            <div>
              <h1>GonsGarage</h1>
              <p>My Cars</p>
            </div>
          </div>
          <div className={styles.userSection}>
            <span>Welcome, {user?.firstName} {user?.lastName}</span>
            <button onClick={logout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation - following Agent.md navigation patterns */}
      <nav className={styles.navigation}>
        <button 
          onClick={() => router.push('/dashboard')}
          className={styles.navButton}
        >
          Dashboard
        </button>
        <button 
          onClick={() => router.push('/cars')}
          className={`${styles.navButton} ${styles.active}`}
        >
          My Cars
        </button>
        <button 
          onClick={() => router.push('/appointments')}
          className={styles.navButton}
        >
          Appointments
        </button>
      </nav>

      {/* Main Content */}
      <main className={styles.main}>
        {/* Error Alert - following Agent.md error handling */}
        {error && <ErrorAlert message={error} />}

        {/* Controls */}
        <div className={styles.controls}>
          <div className={styles.controlsLeft}>
            <h2>My Cars ({cars.length})</h2>
            <p>Manage your registered vehicles</p>
          </div>
          <button 
            onClick={() => setShowCreateModal(true)}
            className={styles.addButton}
          >
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            Add Car
          </button>
        </div>

        {/* Car List - following Agent.md component composition */}
        <CarList
          cars={cars}
          onEdit={handleEditCar}
          onDelete={handleDeleteCar}
          onViewDetails={(id) => router.push(`/cars/${id}`)}
          onScheduleService={(id) => router.push(`/appointments/new?carId=${id}`)}
        />
      </main>

      {/* Car Modal - following Agent.md modal patterns */}
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
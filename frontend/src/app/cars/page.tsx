'use client';

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient, Car, CreateCarRequest } from '@/lib/api';
import Image from 'next/image';
import styles from './cars.module.css';

export default function CarsPage() {
  const [cars, setCars] = useState<Car[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingCar, setEditingCar] = useState<Car | null>(null);

  const { user, logout } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    fetchCars();
  }, [user, router]);

  const fetchCars = async () => {
    try {
      setLoading(true);
      setError(null);

      const { data, error: apiError } = await apiClient.getCars();
      
      if (data && !apiError) {
        setCars(data);
      } else {
        setError(apiError?.message || 'Failed to fetch cars');
      }
    } catch (err) {
      setError('Network error occurred');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this car?')) {
      return;
    }

    try {
      const { error } = await apiClient.deleteCar(id);
      if (!error) {
        setCars(cars.filter(car => car.id !== id));
      } else {
        alert('Failed to delete car: ' + error.message);
      }
    } catch (err) {
      alert('Network error occurred');
    }
  };

  if (loading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <span>Loading cars...</span>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
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
            <span>Welcome, {user?.first_name} {user?.last_name}</span>
            <button onClick={logout} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Navigation */}
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
        {error && (
          <div className={styles.errorAlert}>
            <span>{error}</span>
          </div>
        )}

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

        {/* Cars Grid */}
        {cars.length === 0 ? (
          <div className={styles.emptyState}>
            <div className={styles.emptyIcon}>🚗</div>
            <h3>No cars registered yet</h3>
            <p>Add your first car to get started with our services</p>
            <button 
              onClick={() => setShowCreateModal(true)}
              className={styles.primaryButton}
            >
              Add Your First Car
            </button>
          </div>
        ) : (
          <div className={styles.carsGrid}>
            {cars.map((car) => (
              <div key={car.id} className={styles.carCard}>
                <div className={styles.carHeader}>
                  <div className={styles.carIcon}>🚗</div>
                  <div className={styles.carActions}>
                    <button 
                      onClick={() => setEditingCar(car)}
                      className={styles.editButton}
                      title="Edit car"
                    >
                      <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                      </svg>
                    </button>
                    <button 
                      onClick={() => handleDelete(car.id)}
                      className={styles.deleteButton}
                      title="Delete car"
                    >
                      <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                </div>
                <div className={styles.carContent}>
                  <h3>{car.year} {car.make} {car.model}</h3>
                  <div className={styles.carDetails}>
                    <div className={styles.detail}>
                      <span className={styles.label}>License Plate:</span>
                      <span className={styles.value}>{car.license_plate}</span>
                    </div>
                    <div className={styles.detail}>
                      <span className={styles.label}>Color:</span>
                      <span className={styles.value}>{car.color}</span>
                    </div>
                    {car.vin && (
                      <div className={styles.detail}>
                        <span className={styles.label}>VIN:</span>
                        <span className={styles.value}>{car.vin}</span>
                      </div>
                    )}
                  </div>
                </div>
                <div className={styles.carFooter}>
                  <button 
                    onClick={() => router.push(`/cars/${car.id}`)}
                    className={styles.viewDetailsButton}
                  >
                    View Details & Repairs
                  </button>
                  <button 
                    onClick={() => router.push(`/appointments/new?carId=${car.id}`)}
                    className={styles.scheduleButton}
                  >
                    Schedule Service
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>

      {/* Car Modal */}
      {(showCreateModal || editingCar) && (
        <CarModal
          car={editingCar}
          onClose={() => {
            setShowCreateModal(false);
            setEditingCar(null);
          }}
          onSuccess={() => {
            fetchCars();
            setShowCreateModal(false);
            setEditingCar(null);
          }}
        />
      )}
    </div>
  );
}

// Car Modal Component
interface CarModalProps {
  car?: Car | null;
  onClose: () => void;
  onSuccess: () => void;
}

function CarModal({ car, onClose, onSuccess }: CarModalProps) {
  const [formData, setFormData] = useState<CreateCarRequest>({
    make: car?.make || '',
    model: car?.model || '',
    year: car?.year || new Date().getFullYear(),
    license_plate: car?.license_plate || '',
    vin: car?.vin || '',
    color: car?.color || '',
  });
  const [errors, setErrors] = useState<{[key: string]: string}>({});
  const [isLoading, setIsLoading] = useState(false);

  const validateForm = () => {
    const newErrors: {[key: string]: string} = {};

    if (!formData.make.trim()) newErrors.make = 'Make is required';
    if (!formData.model.trim()) newErrors.model = 'Model is required';
    if (!formData.year || formData.year < 1900 || formData.year > new Date().getFullYear() + 2) {
      newErrors.year = 'Valid year is required';
    }
    if (!formData.license_plate.trim()) newErrors.license_plate = 'License plate is required';
    if (!formData.color.trim()) newErrors.color = 'Color is required';

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsLoading(true);

    try {
      if (car) {
        const { error } = await apiClient.updateCar(car.id, formData);
        if (error) {
          alert('Failed to update car: ' + error.message);
          return;
        }
      } else {
        const { error } = await apiClient.createCar(formData);
        if (error) {
          alert('Failed to create car: ' + error.message);
          return;
        }
      }

      onSuccess();
    } catch (err) {
      alert('Network error occurred');
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: name === 'year' ? parseInt(value) || 0 : value
    }));

    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  return (
    <div className={styles.modalOverlay}>
      <div className={styles.modal}>
        <div className={styles.modalHeader}>
          <h3>{car ? 'Edit Car' : 'Add New Car'}</h3>
          <button onClick={onClose} className={styles.closeButton}>
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.modalForm}>
          <div className={styles.formGrid}>
            <div className={styles.formGroup}>
              <label htmlFor="make">Make</label>
              <input
                id="make"
                name="make"
                type="text"
                value={formData.make}
                onChange={handleChange}
                placeholder="e.g., Toyota"
                className={errors.make ? styles.inputError : ''}
              />
              {errors.make && <span className={styles.errorText}>{errors.make}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="model">Model</label>
              <input
                id="model"
                name="model"
                type="text"
                value={formData.model}
                onChange={handleChange}
                placeholder="e.g., Camry"
                className={errors.model ? styles.inputError : ''}
              />
              {errors.model && <span className={styles.errorText}>{errors.model}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="year">Year</label>
              <input
                id="year"
                name="year"
                type="number"
                value={formData.year}
                onChange={handleChange}
                min="1900"
                max={new Date().getFullYear() + 2}
                className={errors.year ? styles.inputError : ''}
              />
              {errors.year && <span className={styles.errorText}>{errors.year}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="color">Color</label>
              <input
                id="color"
                name="color"
                type="text"
                value={formData.color}
                onChange={handleChange}
                placeholder="e.g., Blue"
                className={errors.color ? styles.inputError : ''}
              />
              {errors.color && <span className={styles.errorText}>{errors.color}</span>}
            </div>
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="license_plate">License Plate</label>
            <input
              id="license_plate"
              name="license_plate"
              type="text"
              value={formData.license_plate}
              onChange={handleChange}
              placeholder="e.g., ABC-1234"
              className={errors.license_plate ? styles.inputError : ''}
            />
            {errors.license_plate && <span className={styles.errorText}>{errors.license_plate}</span>}
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="vin">VIN (Optional)</label>
            <input
              id="vin"
              name="vin"
              type="text"
              value={formData.vin}
              onChange={handleChange}
              placeholder="17-character VIN"
              maxLength={17}
            />
          </div>

          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={onClose}
              className={styles.cancelButton}
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={isLoading}
              className={styles.submitButton}
            >
              {isLoading ? 'Saving...' : (car ? 'Update Car' : 'Add Car')}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
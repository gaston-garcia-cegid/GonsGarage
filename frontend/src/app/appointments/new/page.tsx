'use client';

import React, { useState, useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuthStore } from '@/stores/auth.store';
import { useCarStore } from '@/stores/car.store';
import { useAppointments } from '@/stores/appointment.store';
import { CreateAppointmentRequest, SERVICE_TYPES } from '@/shared/types';
import styles from './new-appointment.module.css';

function NewAppointmentForm() {
  const [formData, setFormData] = useState<CreateAppointmentRequest>({
    carId: '',           // ✅ camelCase per Agent.md
    serviceType: '',     // ✅ camelCase per Agent.md
    scheduledAt: '',     // ✅ camelCase per Agent.md
    notes: '',
  });
  const [customServiceType, setCustomServiceType] = useState('');
  const [errors, setErrors] = useState<{[key: string]: string}>({});

  const { user } = useAuthStore();
  const { cars, isLoading: carsLoading, fetchCars } = useCarStore();
  const { createAppointment, isCreating, error } = useAppointments();
  const router = useRouter();
  const searchParams = useSearchParams();
  const preselectedCarId = searchParams.get('carId');

  useEffect(() => {
    if (!user) {
      router.push('/auth/login');
      return;
    }
    fetchCars();
  }, [user, router, fetchCars]);

  useEffect(() => {
    if (preselectedCarId && cars.length > 0) {
      setFormData(prev => ({
        ...prev,
        carId: preselectedCarId    // ✅ camelCase per Agent.md
      }));
    }
  }, [preselectedCarId, cars.length]);

  const validateForm = (): boolean => {
    const newErrors: {[key: string]: string} = {};

    if (!formData.carId) {
      newErrors.carId = 'Please select a car';
    }

    if (!formData.serviceType) {
      newErrors.serviceType = 'Please select a service type';
    }

    if (formData.serviceType === 'other' && !customServiceType.trim()) {
      newErrors.customServiceType = 'Please specify the service type';
    }

    if (!formData.scheduledAt) {
      newErrors.scheduledAt = 'Please select a date and time';
    } else {
      const selectedDate = new Date(formData.scheduledAt);
      const now = new Date();
      if (selectedDate <= now) {
        newErrors.scheduledAt = 'Please select a future date and time';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    try {
      // ✅ Data already in camelCase format per Agent.md
      const appointmentData: CreateAppointmentRequest = {
        carId: formData.carId,
        serviceType: formData.serviceType === 'other' ? customServiceType : formData.serviceType,
        scheduledAt: formData.scheduledAt,
        notes: formData.notes,
      };

      const success = await createAppointment(appointmentData);
      
      if (success) {
        router.push('/appointments?success=true');
      }
    } catch (error) {
      console.error('Failed to create appointment:', error);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: value
    }));

    // Clear error when user starts typing
    if (errors[name]) {
      setErrors(prev => ({
        ...prev,
        [name]: ''
      }));
    }
  };

  const getMinDateTime = () => {
    const now = new Date();
    now.setHours(now.getHours() + 1);
    return now.toISOString().slice(0, 16);
  };

  const selectedCar = cars.find(car => car.id === formData.carId);
  const selectedService = SERVICE_TYPES.find(service => service.id === formData.serviceType);

  if (carsLoading) {
    return (
      <div className={styles.loadingContainer}>
        <div className={styles.spinner}></div>
        <span>Loading...</span>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.logoSection}>
            <div className={styles.logoIcon}>🔧</div>
            <div>
              <h1>GonsGarage</h1>
              <p>Schedule Service</p>
            </div>
          </div>
          <div className={styles.userSection}>
            <span>Welcome, {user?.firstName} {user?.lastName}</span>
            <button onClick={() => router.push('/auth/login')} className={styles.logoutButton}>
              Logout
            </button>
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className={styles.main}>
        <div className={styles.formContainer}>
          <div className={styles.formHeader}>
            <button 
              onClick={() => router.back()}
              className={styles.backButton}
            >
              ← Back
            </button>
            <h2>Schedule a Service Appointment</h2>
            <p>Book an appointment for your vehicle maintenance</p>
          </div>

          {error && (
            <div className={styles.errorAlert}>
              {error}
            </div>
          )}

          <form onSubmit={handleSubmit} className={styles.appointmentForm}>
            {/* Car Selection */}
            <div className={styles.section}>
              <h3>Select Your Vehicle *</h3>
              <div className={styles.formGroup}>
                <select
                  id="carId"
                  name="carId"           // ✅ camelCase per Agent.md
                  value={formData.carId}
                  onChange={handleChange}
                  className={errors.carId ? styles.inputError : ''}
                >
                  <option value="">Choose a vehicle...</option>
                  {cars.map(car => (
                    <option key={car.id} value={car.id}>
                      {car.year} {car.make} {car.model} - {car.licensePlate}
                    </option>
                  ))}
                </select>
                {errors.carId && <span className={styles.errorText}>{errors.carId}</span>}
              </div>

              {selectedCar && (
                <div className={styles.selectedCarInfo}>
                  <h4>Selected Vehicle:</h4>
                  <p>{selectedCar.year} {selectedCar.make} {selectedCar.model}</p>
                  <p>License Plate: {selectedCar.licensePlate}</p>
                </div>
              )}
            </div>

            {/* Service Type Selection */}
            <div className={styles.section}>
              <h3>Service Type *</h3>
              <div className={styles.serviceGrid}>
                {SERVICE_TYPES.map(service => (
                  <label key={service.id} className={styles.serviceOption}>
                    <input
                      type="radio"
                      name="serviceType"   // ✅ camelCase per Agent.md
                      value={service.id}
                      checked={formData.serviceType === service.id}
                      onChange={handleChange}
                    />
                    <div className={styles.serviceContent}>
                      <h4>{service.name}</h4>
                      <p>{service.description}</p>
                    </div>
                  </label>
                ))}
              </div>
              {errors.serviceType && <span className={styles.errorText}>{errors.serviceType}</span>}

              {formData.serviceType === 'other' && (
                <div className={styles.formGroup}>
                  <label htmlFor="customServiceType">Specify Service Type *</label>
                  <input
                    id="customServiceType"
                    type="text"
                    value={customServiceType}
                    onChange={(e) => setCustomServiceType(e.target.value)}
                    placeholder="Describe the service you need..."
                    className={errors.customServiceType ? styles.inputError : ''}
                  />
                  {errors.customServiceType && <span className={styles.errorText}>{errors.customServiceType}</span>}
                </div>
              )}
            </div>

            {/* Date and Time */}
            <div className={styles.section}>
              <h3>Preferred Date & Time *</h3>
              <div className={styles.formGroup}>
                <label htmlFor="scheduledAt">Date and Time</label>
                <input
                  id="scheduledAt"
                  name="scheduledAt"     // ✅ camelCase per Agent.md
                  type="datetime-local"
                  value={formData.scheduledAt}
                  onChange={handleChange}
                  min={getMinDateTime()}
                  className={errors.scheduledAt ? styles.inputError : ''}
                />
                {errors.scheduledAt && <span className={styles.errorText}>{errors.scheduledAt}</span>}
              </div>
            </div>

            {/* Additional Notes */}
            <div className={styles.section}>
              <h3>Additional Information</h3>
              <div className={styles.formGroup}>
                <label htmlFor="notes">Notes</label>
                <textarea
                  id="notes"
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  placeholder="Any specific concerns, parts needed, or special instructions..."
                  rows={4}
                  className={styles.textarea}
                />
              </div>
            </div>

            {/* Summary */}
            {selectedCar && selectedService && formData.scheduledAt && (
              <div className={styles.appointmentSummary}>
                <h3>Appointment Summary</h3>
                <div className={styles.summaryContent}>
                  <p><strong>Vehicle:</strong> {selectedCar.year} {selectedCar.make} {selectedCar.model}</p>
                  <p><strong>Service:</strong> {selectedService.name}</p>
                  <p><strong>Date:</strong> {new Date(formData.scheduledAt).toLocaleDateString()}</p>
                  <p><strong>Time:</strong> {new Date(formData.scheduledAt).toLocaleTimeString()}</p>
                </div>
              </div>
            )}

            {/* Form Actions */}
            <div className={styles.formActions}>
              <button 
                type="button" 
                onClick={() => router.back()}
                className={styles.cancelButton}
              >
                Cancel
              </button>
              <button 
                type="submit" 
                disabled={isCreating}
                className={styles.submitButton}
              >
                {isCreating ? 'Scheduling...' : 'Schedule Appointment'}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}

function LoadingFallback() {
  return (
    <div className={styles.loadingFallback}>
      <div className={styles.spinner}></div>
      <span>Loading...</span>
    </div>
  );
}

export default function NewAppointmentPage() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <NewAppointmentForm />
    </Suspense>
  );
}
'use client';

import React, { useState, useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth, useCars, useAppointments } from '@/stores';
import { CreateAppointmentRequest } from '@/lib/api';
import styles from '../appointments.module.css';

const SERVICE_TYPES = [
  { id: 'oil_change', name: 'Oil Change', description: 'Regular oil and filter change' },
  { id: 'brake_service', name: 'Brake Service', description: 'Brake pads, rotors, and fluid' },
  { id: 'tire_service', name: 'Tire Service', description: 'Tire rotation, alignment, replacement' },
  { id: 'engine_diagnostic', name: 'Engine Diagnostic', description: 'Check engine light and diagnostics' },
  { id: 'transmission_service', name: 'Transmission Service', description: 'Transmission fluid and inspection' },
  { id: 'air_conditioning', name: 'Air Conditioning', description: 'A/C repair and maintenance' },
  { id: 'battery_service', name: 'Battery Service', description: 'Battery testing and replacement' },
  { id: 'general_maintenance', name: 'General Maintenance', description: 'Multi-point inspection' },
  { id: 'other', name: 'Other', description: 'Custom service request' },
];

// Separate component that uses useSearchParams
function NewAppointmentForm() {
  const [formData, setFormData] = useState<CreateAppointmentRequest>({
    car_id: '',
    service_type: '',
    scheduled_at: '',
    notes: '',
  });
  const [customServiceType, setCustomServiceType] = useState('');
  const [errors, setErrors] = useState<{[key: string]: string}>({});
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { user, logout } = useAuth();
  const { cars, isLoading: carsLoading, fetchCars } = useCars();
  const { createAppointment, isCreating } = useAppointments();
  const router = useRouter();
  const searchParams = useSearchParams();
  const preselectedCarId = searchParams.get('carId');
  
  const loading = carsLoading;

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
        car_id: preselectedCarId
      }));
    }
  }, [preselectedCarId, cars.length]);

  const validateForm = (): boolean => {
    const newErrors: {[key: string]: string} = {};

    if (!formData.car_id) {
      newErrors.car_id = 'Please select a car';
    }

    if (!formData.service_type) {
      newErrors.service_type = 'Please select a service type';
    }

    if (formData.service_type === 'other' && !customServiceType.trim()) {
      newErrors.customServiceType = 'Please specify the service type';
    }

    if (!formData.scheduled_at) {
      newErrors.scheduled_at = 'Please select a date and time';
    } else {
      const selectedDate = new Date(formData.scheduled_at);
      const now = new Date();
      if (selectedDate <= now) {
        newErrors.scheduled_at = 'Please select a future date and time';
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) return;

    setIsSubmitting(true);

    try {
      // Convert to camelCase for store
      const appointmentData = {
        carId: formData.car_id,
        serviceType: formData.service_type === 'other' ? customServiceType : formData.service_type,
        scheduledAt: formData.scheduled_at,
        notes: formData.notes,
      };

      const success = await createAppointment(appointmentData);
      
      if (!success) {
        alert('Failed to create appointment');
        return;
      }

      // Success - redirect to appointments list
      router.push('/appointments?success=true');
    } catch {
      alert('Network error occurred');
    } finally {
      setIsSubmitting(false);
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

  const handleCustomServiceTypeChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setCustomServiceType(e.target.value);
    if (errors.customServiceType) {
      setErrors(prev => ({
        ...prev,
        customServiceType: ''
      }));
    }
  };

  const getMinDateTime = () => {
    const now = new Date();
    now.setHours(now.getHours() + 1); // At least 1 hour from now
    return now.toISOString().slice(0, 16);
  };

  const selectedCar = cars.find(car => car.id === formData.car_id);
  const selectedService = SERVICE_TYPES.find(service => service.id === formData.service_type);

  if (loading) {
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
            <div className={styles.logoIcon}>
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-2m-2 0H7m5 0v-5a2 2 0 012-2h2a2 2 0 012 2v5" />
              </svg>
            </div>
            <div>
              <h1>GonsGarage</h1>
              <p>Schedule Service</p>
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
          className={styles.navButton}
        >
          My Cars
        </button>
        <button 
          onClick={() => router.push('/appointments')}
          className={`${styles.navButton} ${styles.active}`}
        >
          Appointments
        </button>
      </nav>

      {/* Main Content */}
      <main className={styles.main}>
        <div className={styles.formContainer}>
          <div className={styles.formHeader}>
            <button 
              onClick={() => router.back()}
              className={styles.backButton}
            >
              <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
              </svg>
              Back
            </button>
            <h2>Schedule a Service Appointment</h2>
            <p>Book an appointment for your vehicle maintenance</p>
          </div>

          <form onSubmit={handleSubmit} className={styles.appointmentForm}>
            {/* Car Selection */}
            <div className={styles.section}>
              <h3>Select Your Vehicle</h3>
              <div className={styles.formGroup}>
                <label htmlFor="car_id">Choose a car *</label>
                <select
                  id="car_id"
                  name="car_id"
                  value={formData.car_id}
                  onChange={handleChange}
                  className={errors.car_id ? styles.inputError : ''}
                >
                  <option value="">Select a car...</option>
                  {cars.map((car) => (
                    <option key={car.id} value={car.id}>
                      {car.year} {car.make} {car.model} - {car.licensePlate}
                    </option>
                  ))}
                </select>
                {errors.car_id && <span className={styles.errorText}>{errors.car_id}</span>}
              </div>

              {selectedCar && (
                <div className={styles.selectedCarInfo}>
                  <div className={styles.carIcon}>🚗</div>
                  <div className={styles.carDetails}>
                    <h4>{selectedCar.year} {selectedCar.make} {selectedCar.model}</h4>
                    <p>{selectedCar.licensePlate} • {selectedCar.color}</p>
                  </div>
                </div>
              )}
            </div>

            {/* Service Type Selection */}
            <div className={styles.section}>
              <h3>Service Type *</h3>
              <div className={styles.serviceGrid}>
                {SERVICE_TYPES.map((service) => (
                  <label key={service.id} className={styles.serviceOption}>
                    <input
                      type="radio"
                      name="service_type"
                      value={service.id}
                      checked={formData.service_type === service.id}
                      onChange={handleChange}
                    />
                    <div className={styles.serviceCard}>
                      <h4>{service.name}</h4>
                      <p>{service.description}</p>
                    </div>
                  </label>
                ))}
              </div>
              {errors.service_type && <span className={styles.errorText}>{errors.service_type}</span>}

              {formData.service_type === 'other' && (
                <div className={styles.formGroup}>
                  <label htmlFor="customServiceType">Specify service type *</label>
                  <input
                    id="customServiceType"
                    type="text"
                    value={customServiceType}
                    onChange={handleCustomServiceTypeChange}
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
                <label htmlFor="scheduled_at">Select date and time</label>
                <input
                  id="scheduled_at"
                  name="scheduled_at"
                  type="datetime-local"
                  value={formData.scheduled_at}
                  onChange={handleChange}
                  min={getMinDateTime()}
                  className={errors.scheduled_at ? styles.inputError : ''}
                />
                {errors.scheduled_at && <span className={styles.errorText}>{errors.scheduled_at}</span>}
                <small className={styles.helpText}>
                  Please select a date and time at least 1 hour from now
                </small>
              </div>
            </div>

            {/* Additional Notes */}
            <div className={styles.section}>
              <h3>Additional Information</h3>
              <div className={styles.formGroup}>
                <label htmlFor="notes">Special requests or notes (optional)</label>
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
            {selectedCar && selectedService && formData.scheduled_at && (
              <div className={styles.appointmentSummary}>
                <h3>Appointment Summary</h3>
                <div className={styles.summaryGrid}>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Vehicle:</span>
                    <span className={styles.summaryValue}>
                      {selectedCar.year} {selectedCar.make} {selectedCar.model} ({selectedCar.licensePlate})
                    </span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Service:</span>
                    <span className={styles.summaryValue}>
                      {formData.service_type === 'other' ? customServiceType : selectedService.name}
                    </span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Date & Time:</span>
                    <span className={styles.summaryValue}>
                      {new Date(formData.scheduled_at).toLocaleString()}
                    </span>
                  </div>
                </div>
              </div>
            )}

            {/* Form Actions */}
            <div className={styles.formActions}>
              <button
                type="button"
                onClick={() => router.push('/appointments')}
                className={styles.cancelButton}
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={isSubmitting}
                className={styles.submitButton}
              >
                {isSubmitting ? 'Scheduling...' : 'Schedule Appointment'}
              </button>
            </div>
          </form>
        </div>
      </main>
    </div>
  );
}

// Loading component
function LoadingFallback() {
  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      flexDirection: 'column',
      alignItems: 'center',
      justifyContent: 'center',
      gap: '0.75rem'
    }}>
      <div style={{
        width: '20px',
        height: '20px',
        border: '2px solid #e5e7eb',
        borderTop: '2px solid #2563eb',
        borderRadius: '50%',
        animation: 'spin 1s linear infinite'
      }}></div>
      <span>Loading...</span>
    </div>
  );
}

// Main component with Suspense boundary
export default function NewAppointmentPage() {
  return (
    <Suspense fallback={<LoadingFallback />}>
      <NewAppointmentForm />
    </Suspense>
  );
}
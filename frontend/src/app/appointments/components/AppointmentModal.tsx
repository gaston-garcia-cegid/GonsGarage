import React, { useState } from 'react';
import { Car, CreateCarRequest, CarFormData } from '@/types/car';
import { useCarValidation } from '@/hooks/useCarValidation';
import styles from '../cars.module.css';

interface CarModalProps {
  car: Car | null;
  onClose: () => void;
  onCreate: (carData: CreateCarRequest) => Promise<boolean>;
  onUpdate: (id: string, carData: Partial<CreateCarRequest>) => Promise<boolean>;
}

// Car modal component following Agent.md modal patterns
export default function CarModal({ 
  car, 
  onClose, 
  onCreate, 
  onUpdate 
}: CarModalProps) {
  const [formData, setFormData] = useState<CarFormData>({
    make: car?.make || '',
    model: car?.model || '',
    year: car?.year || new Date().getFullYear(),
    licensePlate: car?.licensePlate || '',
    vin: car?.vin || '',
    color: car?.color || '',
    mileage: car?.mileage || 0,
  });
  const [isLoading, setIsLoading] = useState(false);

  const { errors, validateCar, clearFieldError } = useCarValidation();

  // Handle form field changes - following Agent.md form handling
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    
    setFormData(prev => ({
      ...prev,
      [name]: name === 'year' || name === 'mileage' ? parseInt(value) || 0 : value
    }));

    // Clear field error when user starts typing
    if (errors[name as keyof typeof errors]) {
      clearFieldError(name as keyof typeof errors);
    }
  };

  // Handle form submission - following Agent.md error handling
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateCar(formData)) {
      return;
    }

    setIsLoading(true);

    try {
      const carData: CreateCarRequest = {
        make: formData.make,
        model: formData.model,
        year: formData.year,
        licensePlate: formData.licensePlate,
        vin: formData.vin || undefined,
        color: formData.color,
        mileage: formData.mileage || undefined,
      };

      let success = false;

      if (car) {
        // Update existing car
        success = await onUpdate(car.id, carData);
      } else {
        // Create new car
        success = await onCreate(carData);
      }

      if (success) {
        onClose();
      } else {
        errors.general = 'Failed to save car. Please try again.';
        //alert('Failed to save car. Please try again.');
      }
    } catch (err) {
      console.log(err);
      alert('An error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  // Handle modal backdrop click
  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div className={styles.modalOverlay} onClick={handleBackdropClick}>
      <div className={styles.modal}>
        <div className={styles.modalHeader}>
          <h3>{car ? 'Edit Car' : 'Add New Car'}</h3>
          <button 
            onClick={onClose} 
            className={styles.closeButton}
            aria-label="Close modal"
          >
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form onSubmit={handleSubmit} className={styles.modalForm}>
            {errors.general && (
              <div style={{
                backgroundColor: '#fef2f2',
                border: '1px solid #fecaca',
                color: '#dc2626',
                padding: 'var(--space-3)',
                borderRadius: 'var(--radius)',
                marginBottom: 'var(--space-4)',
                fontSize: '0.875rem',
              }}>
                {errors.general}
              </div>
            )}
          <div className={styles.formGrid}>
            {/* Make Field */}
            <div className={styles.formGroup}>
              <label htmlFor="make">Make *</label>
              <input
                id="make"
                name="make"
                type="text"
                value={formData.make}
                onChange={handleChange}
                placeholder="e.g., Toyota"
                className={errors.make ? styles.inputError : ''}
                disabled={isLoading}
              />
              {errors.make && <span className={styles.errorText}>{errors.make}</span>}
            </div>

            {/* Model Field */}
            <div className={styles.formGroup}>
              <label htmlFor="model">Model *</label>
              <input
                id="model"
                name="model"
                type="text"
                value={formData.model}
                onChange={handleChange}
                placeholder="e.g., Camry"
                className={errors.model ? styles.inputError : ''}
                disabled={isLoading}
              />
              {errors.model && <span className={styles.errorText}>{errors.model}</span>}
            </div>

            {/* Year Field */}
            <div className={styles.formGroup}>
              <label htmlFor="year">Year *</label>
              <input
                id="year"
                name="year"
                type="number"
                value={formData.year}
                onChange={handleChange}
                min="1900"
                max={new Date().getFullYear() + 2}
                className={errors.year ? styles.inputError : ''}
                disabled={isLoading}
              />
              {errors.year && <span className={styles.errorText}>{errors.year}</span>}
            </div>

            {/* Color Field */}
            <div className={styles.formGroup}>
              <label htmlFor="color">Color *</label>
              <input
                id="color"
                name="color"
                type="text"
                value={formData.color}
                onChange={handleChange}
                placeholder="e.g., Blue"
                className={errors.color ? styles.inputError : ''}
                disabled={isLoading}
              />
              {errors.color && <span className={styles.errorText}>{errors.color}</span>}
            </div>
          </div>

          {/* License Plate Field */}
          <div className={styles.formGroup}>
            <label htmlFor="licensePlate">License Plate *</label>
            <input
              id="licensePlate"
              name="licensePlate"
              type="text"
              value={formData.licensePlate}
              onChange={handleChange}
              placeholder="e.g., ABC-1234"
              className={errors.licensePlate ? styles.inputError : ''}
              disabled={isLoading}
            />
            {errors.licensePlate && <span className={styles.errorText}>{errors.licensePlate}</span>}
          </div>

          {/* VIN Field */}
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
              disabled={isLoading}
            />
          </div>

          {/* Mileage Field */}
          <div className={styles.formGroup}>
            <label htmlFor="mileage">Mileage (Optional)</label>
            <input
              id="mileage"
              name="mileage"
              type="number"
              value={formData.mileage}
              onChange={handleChange}
              min="0"
              placeholder="Current mileage"
              className={errors.mileage ? styles.inputError : ''}
              disabled={isLoading}
            />
            {errors.mileage && <span className={styles.errorText}>{errors.mileage}</span>}
          </div>

          {/* Modal Actions */}
          <div className={styles.modalActions}>
            <button
              type="button"
              onClick={onClose}
              className={styles.cancelButton}
              disabled={isLoading}
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
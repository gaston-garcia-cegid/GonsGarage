'use client';

import React, { useState } from 'react';
import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/types/appointment';
import { SERVICE_TYPES } from '@/shared/types';
import { useCarStore } from '@/stores/car.store';
import styles from '../appointments.module.css';

interface AppointmentModalProps {
  appointment?: Appointment | null;
  onClose: () => void;
  onCreate: (data: CreateAppointmentRequest) => Promise<boolean>;
  onUpdate?: (id: string, data: Partial<UpdateAppointmentRequest>) => Promise<boolean>;
  preSelectedCarId?: string;
}

interface FormData {
  carId: string;
  service: string;
  date: string;
  time: string;
  notes: string;
}

export default function AppointmentModal({
  appointment,
  onClose,
  onCreate,
  onUpdate,
  preSelectedCarId,
}: AppointmentModalProps) {
  const { cars } = useCarStore();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string>('');

  const [formData, setFormData] = useState<FormData>({
    carId: appointment?.carId || preSelectedCarId || '',
    service: appointment?.service || '',
    date: appointment?.date ? appointment.date.split('T')[0] : '',
    time: appointment?.date ? new Date(appointment.date).toTimeString().slice(0, 5) : '',
    notes: appointment?.notes || '',
  });

  const selectedCar = cars.find(c => c.id === formData.carId);

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
  ) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
    setError('');
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    // Validation
    if (!formData.carId) {
      setError('Please select a car');
      return;
    }

    if (!formData.service) {
      setError('Please select a service');
      return;
    }

    if (!formData.date || !formData.time) {
      setError('Please select date and time');
      return;
    }

    setIsLoading(true);

    try {
      const appointmentData: CreateAppointmentRequest = {
        clientName: selectedCar ? `${selectedCar.make} ${selectedCar.model}` : 'Unknown',
        carId: formData.carId,
        service: formData.service,
        date: `${formData.date}T${formData.time}:00Z`,
        notes: formData.notes || undefined,
        status: 'scheduled',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        time: ''
      };

      let success = false;

      if (appointment && onUpdate) {
        const appointmentUpdateData: UpdateAppointmentRequest = {
              carId: formData.carId,
              service: formData.service,
              date: `${formData.date}T${formData.time}:00Z`,
              notes: formData.notes || undefined,
              status: 'scheduled',
            };

        success = await onUpdate(appointment.id, appointmentUpdateData);
      } else {
        success = await onCreate(appointmentData);
      }

      if (success) {
        onClose();
      } else {
        setError('Failed to save appointment. Please try again.');
      }
    } catch (err) {
      console.error('Error saving appointment:', err);
      setError('An error occurred. Please try again.');
    } finally {
      setIsLoading(false);
    }
  };

  const handleBackdropClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  // Get minimum date (today)
  const today = new Date().toISOString().split('T')[0];

  return (
    <div className={styles.modalOverlay} onClick={handleBackdropClick}>
      <div className={styles.modal}>
        {/* Header */}
        <div className={styles.modalHeader}>
          <h3>{appointment ? 'Edit Appointment' : 'Schedule Appointment'}</h3>
          <button
            onClick={onClose}
            className={styles.closeButton}
            aria-label="Close modal"
          >
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Body */}
        <form onSubmit={handleSubmit}>
          <div className={styles.modalBody}>
            {error && (
              <div className={styles.errorAlert}>
                {error}
              </div>
            )}

            {/* Car Selection */}
            <div className={styles.section}>
              <h3>Select Vehicle</h3>
              {selectedCar ? (
                <div className={styles.selectedCarInfo}>
                  <div className={styles.carIcon}>🚗</div>
                  <div className={styles.carDetails}>
                    <h4>{selectedCar.year} {selectedCar.make} {selectedCar.model}</h4>
                    <p>{selectedCar.licensePlate} • {selectedCar.color}</p>
                  </div>
                  {!appointment && (
                    <button
                      type="button"
                      onClick={() => setFormData(prev => ({ ...prev, carId: '' }))}
                      className={styles.cancelButton}
                      style={{ padding: '0.5rem 1rem' }}
                    >
                      Change
                    </button>
                  )}
                </div>
              ) : (
                <div className={styles.formGroup}>
                  <select
                    name="carId"
                    value={formData.carId}
                    onChange={handleChange}
                    required
                    disabled={isLoading}
                  >
                    <option value="">Select a car...</option>
                    {cars.map(car => (
                      <option key={car.id} value={car.id}>
                        {car.year} {car.make} {car.model} - {car.licensePlate}
                      </option>
                    ))}
                  </select>
                </div>
              )}
            </div>

            {/* Service Selection */}
            <div className={styles.section}>
              <h3>Select Service</h3>
              <div className={styles.serviceGrid}>
                {SERVICE_TYPES.map(service => (
                  <label key={service.id} className={styles.serviceOption}>
                    <input
                      type="radio"
                      name="service"
                      value={service.id}
                      checked={formData.service === service.id}
                      onChange={handleChange}
                      disabled={isLoading}
                    />
                    <div className={styles.serviceCard}>
                      <h4>{service.name}</h4>
                      <p>{service.description}</p>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Date & Time */}
            <div className={styles.section}>
              <h3>Select Date & Time</h3>
              <div className={styles.formGrid}>
                <div className={styles.formGroup}>
                  <label htmlFor="date">Date *</label>
                  <input
                    id="date"
                    name="date"
                    type="date"
                    value={formData.date}
                    onChange={handleChange}
                    min={today}
                    required
                    disabled={isLoading}
                  />
                </div>
                <div className={styles.formGroup}>
                  <label htmlFor="time">Time *</label>
                  <input
                    id="time"
                    name="time"
                    type="time"
                    value={formData.time}
                    onChange={handleChange}
                    min="08:00"
                    max="18:00"
                    required
                    disabled={isLoading}
                  />
                </div>
              </div>
            </div>

            {/* Notes */}
            <div className={styles.section}>
              <h3>Additional Notes (Optional)</h3>
              <div className={styles.formGroup}>
                <textarea
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  placeholder="Any specific concerns or requests..."
                  rows={3}
                  disabled={isLoading}
                />
              </div>
            </div>

            {/* Summary */}
            {formData.carId && formData.service && formData.date && (
              <div className={styles.appointmentSummary}>
                <h3>Appointment Summary</h3>
                <div className={styles.summaryGrid}>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Vehicle:</span>
                    <span className={styles.summaryValue}>
                      {selectedCar && `${selectedCar.year} ${selectedCar.make} ${selectedCar.model}`}
                    </span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Service:</span>
                    <span className={styles.summaryValue}>
                      {SERVICE_TYPES.find(s => s.id === formData.service)?.name}
                    </span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Date:</span>
                    <span className={styles.summaryValue}>
                      {new Date(formData.date).toLocaleDateString('en-US', {
                        weekday: 'short',
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric',
                      })}
                    </span>
                  </div>
                  <div className={styles.summaryItem}>
                    <span className={styles.summaryLabel}>Time:</span>
                    <span className={styles.summaryValue}>{formData.time}</span>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Footer Actions */}
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
              className={styles.submitButton}
              disabled={isLoading}
            >
              {isLoading
                ? 'Saving...'
                : appointment
                ? 'Update Appointment'
                : 'Schedule Appointment'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
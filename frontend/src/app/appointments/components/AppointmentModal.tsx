'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/types/appointment';
import { SERVICE_TYPES } from '@/shared/types';
import {
  combineLocalDateAndSlot,
  formatSlotLabel24h,
  getBookableSlotTimesForDay,
  isWeekdayLocalYmd,
  isWithinWorkshopHours,
  localYmdFromIso,
  todayLocalYmd,
  weekdayErrorMessage,
  workshopHoursErrorMessage,
  workshopSlotFromIso,
} from '@/lib/workshopAppointmentRules';
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
    date: appointment?.date ? localYmdFromIso(appointment.date) : '',
    time: appointment?.date ? workshopSlotFromIso(appointment.date) : '',
    notes: appointment?.notes || '',
  });

  const bookableSlots = useMemo(() => getBookableSlotTimesForDay(formData.date), [formData.date]);

  useEffect(() => {
    if (!formData.time) return;
    if (!bookableSlots.includes(formData.time)) {
      setFormData((prev) => ({ ...prev, time: '' }));
    }
  }, [formData.date, bookableSlots, formData.time]);

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
    if (!isWeekdayLocalYmd(formData.date)) {
      setError(weekdayErrorMessage());
      return;
    }
    if (bookableSlots.length === 0) {
      setError('No hay horarios disponibles para ese día.');
      return;
    }
    const scheduledLocal = combineLocalDateAndSlot(formData.date, formData.time);
    const scheduledDate = new Date(scheduledLocal);
    if (Number.isNaN(scheduledDate.getTime()) || scheduledDate <= new Date()) {
      setError('El turno tiene que ser en el futuro');
      return;
    }
    if (!isWithinWorkshopHours(scheduledDate)) {
      setError(workshopHoursErrorMessage());
      return;
    }

    setIsLoading(true);

    try {
      const appointmentData: CreateAppointmentRequest = {
        clientName: selectedCar ? `${selectedCar.make} ${selectedCar.model}` : 'Unknown',
        carId: formData.carId,
        service: formData.service,
        date: scheduledLocal,
        notes: formData.notes || undefined,
        status: 'scheduled',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        time: formData.time,
      };

      let success = false;

      if (appointment && onUpdate) {
        const appointmentUpdateData: UpdateAppointmentRequest = {
              carId: formData.carId,
              service: formData.service,
              date: scheduledLocal,
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

  return (
    <div className={styles.modalOverlay} onClick={handleBackdropClick}>
      <div
        className={styles.modal}
        role="dialog"
        aria-modal="true"
        aria-labelledby="appointment-modal-title"
      >
        {/* Header */}
        <div className={styles.modalHeader}>
          <h3 id="appointment-modal-title">{appointment ? 'Edit Appointment' : 'Schedule Appointment'}</h3>
          <button
            type="button"
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
                      className={styles.changeCarButton}
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
              <p className={styles.helpText}>
                Mon–Fri, every 30 minutes: 9:30–12:30 and 14:00–17:30 (24-hour display).
              </p>
              <div className={styles.formGrid}>
                <div className={styles.formGroup}>
                  <label htmlFor="date">Date *</label>
                  <input
                    id="date"
                    name="date"
                    type="date"
                    value={formData.date}
                    onChange={handleChange}
                    min={todayLocalYmd()}
                    required
                    disabled={isLoading}
                  />
                </div>
                <div className={styles.formGroup}>
                  <label htmlFor="time">Time *</label>
                  <select
                    id="time"
                    name="time"
                    value={formData.time}
                    onChange={handleChange}
                    required
                    disabled={isLoading || !formData.date || bookableSlots.length === 0}
                  >
                    <option value="">
                      {!formData.date
                        ? 'Pick a date first…'
                        : bookableSlots.length === 0
                          ? 'No slots that day'
                          : 'Select time…'}
                    </option>
                    {bookableSlots.map((hhmm) => (
                      <option key={hhmm} value={hhmm}>
                        {formatSlotLabel24h(hhmm)}
                      </option>
                    ))}
                  </select>
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
            {formData.carId && formData.service && formData.date && formData.time && (
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
                    <span className={styles.summaryValue}>{formatSlotLabel24h(formData.time)}</span>
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
              className={styles.modalCancelButton}
              disabled={isLoading}
            >
              Cancel
            </button>
            <button
              type="submit"
              className={styles.modalSubmitButton}
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
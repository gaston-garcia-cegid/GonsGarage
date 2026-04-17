'use client';

import React, { useEffect, useMemo, useState } from 'react';
import { useCarStore } from '@/stores/car.store';
import { useAppointments, useAppointmentStore } from '@/stores';
import { CreateAppointmentRequest } from '@/types/appointment';
import { SERVICE_TYPES } from '@/shared/types';
import {
  combineLocalDateAndSlot,
  formatAppointmentSummaryLocal,
  formatSlotLabel24h,
  getBookableSlotTimesForDay,
  isWeekdayLocalYmd,
  isWithinWorkshopHours,
  todayLocalYmd,
  weekdayErrorMessage,
  workshopHoursErrorMessage,
} from '@/lib/workshopAppointmentRules';
import formStyles from '@/app/appointments/new/new-appointment.module.css';
import modalStyles from './NewAppointmentModal.module.css';

const emptyForm = (): CreateAppointmentRequest => ({
  carId: '',
  clientName: '',
  service: '',
  date: '',
  time: '',
  status: 'scheduled',
  notes: '',
  createdAt: new Date().toISOString(),
  updatedAt: new Date().toISOString(),
  deletedAt: undefined,
});

export interface NewAppointmentModalProps {
  readonly isOpen: boolean;
  readonly onClose: () => void;
  readonly initialCarId?: string | null;
  readonly onCreated?: () => void;
}

export default function NewAppointmentModal({
  isOpen,
  onClose,
  initialCarId,
  onCreated,
}: NewAppointmentModalProps) {
  const [formData, setFormData] = useState<CreateAppointmentRequest>(emptyForm);
  const [appointmentDay, setAppointmentDay] = useState('');
  const [appointmentSlot, setAppointmentSlot] = useState('');
  const [customServiceType, setCustomServiceType] = useState('');
  const [errors, setErrors] = useState<Record<string, string>>({});

  const bookableSlots = useMemo(
    () => getBookableSlotTimesForDay(appointmentDay),
    [appointmentDay],
  );

  const { cars } = useCarStore();
  const { createAppointment, isCreating, error } = useAppointments();

  // Do not call fetchCars() here: it toggles the global car store `isLoading`, and the
  // appointments page treats that as a full-page load — it unmounts this modal and the
  // effect runs again → an infinite loop. Cars are already loaded on /appointments.
  useEffect(() => {
    if (!isOpen) return;
    setFormData(emptyForm());
    setAppointmentDay('');
    setAppointmentSlot('');
    setCustomServiceType('');
    setErrors({});
    useAppointmentStore.getState().clearError();
  }, [isOpen]);

  useEffect(() => {
    if (!appointmentSlot) return;
    if (!bookableSlots.includes(appointmentSlot)) {
      setAppointmentSlot('');
    }
  }, [appointmentDay, bookableSlots, appointmentSlot]);

  useEffect(() => {
    if (!isOpen) return;
    if (initialCarId && cars.some((c) => c.id === initialCarId)) {
      setFormData((prev) => ({ ...prev, carId: initialCarId }));
      return;
    }
    if (cars.length === 1) {
      setFormData((prev) => ({ ...prev, carId: cars[0].id }));
    }
  }, [isOpen, initialCarId, cars]);

  const validateForm = (): boolean => {
    const newErrors: Record<string, string> = {};

    if (!formData.carId) {
      newErrors.carId = 'Please select a car';
    }
    if (!formData.service) {
      newErrors.service = 'Please select a service type';
    }
    if (formData.service === 'other' && !customServiceType.trim()) {
      newErrors.customServiceType = 'Please specify the service type';
    }
    if (!appointmentDay) {
      newErrors.date = 'Elegí un día';
    } else if (!isWeekdayLocalYmd(appointmentDay)) {
      newErrors.date = weekdayErrorMessage();
    } else if (bookableSlots.length === 0) {
      newErrors.date = 'No hay horarios disponibles para ese día. Elegí otro día.';
    }
    if (!appointmentSlot) {
      newErrors.time = 'Elegí un horario';
    }
    if (!newErrors.date && !newErrors.time && appointmentDay && appointmentSlot) {
      const combined = combineLocalDateAndSlot(appointmentDay, appointmentSlot);
      const selectedDate = new Date(combined);
      const now = new Date();
      if (Number.isNaN(selectedDate.getTime())) {
        newErrors.date = 'Fecha u hora no válida';
      } else if (selectedDate <= now) {
        newErrors.date = 'El turno tiene que ser en el futuro';
      } else if (!isWithinWorkshopHours(selectedDate)) {
        newErrors.date = workshopHoursErrorMessage();
      }
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!validateForm()) return;

    const scheduledAt = combineLocalDateAndSlot(appointmentDay, appointmentSlot);
    const appointmentData: CreateAppointmentRequest = {
      carId: formData.carId,
      clientName: formData.clientName,
      service: formData.service === 'other' ? customServiceType : formData.service,
      time: appointmentSlot,
      date: scheduledAt,
      status: formData.status,
      notes: formData.notes,
      createdAt: formData.createdAt,
      updatedAt: formData.updatedAt,
      deletedAt: formData.deletedAt,
    };

    const success = await createAppointment(appointmentData);
    if (success) {
      onCreated?.();
      onClose();
      setFormData(emptyForm());
      setCustomServiceType('');
      setErrors({});
    }
  };

  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>,
  ) => {
    const { name, value } = e.target;
    setFormData((prev) => ({ ...prev, [name]: value }));
    if (errors[name]) {
      setErrors((prev) => ({ ...prev, [name]: '' }));
    }
  };

  if (!isOpen) return null;

  const selectedCar = cars.find((car) => car.id === formData.carId);
  const selectedService = SERVICE_TYPES.find((service) => service.id === formData.service);

  return (
    <div
      className={modalStyles.overlay}
      role="presentation"
      onMouseDown={(e) => {
        if (e.target === e.currentTarget) onClose();
      }}
    >
      <div className={modalStyles.dialog} role="dialog" aria-modal="true" aria-labelledby="new-apt-title">
        <div className={formStyles.formContainer}>
          <div className={modalStyles.header}>
            <div>
              <h2 id="new-apt-title">Schedule a service</h2>
              <p>Book an appointment for your car</p>
            </div>
            <button type="button" className={modalStyles.closeButton} onClick={onClose} aria-label="Close">
              ×
            </button>
          </div>

          {error ? <div className={formStyles.errorAlert}>{error}</div> : null}

          <form onSubmit={handleSubmit} className={formStyles.appointmentForm}>
            <div className={formStyles.section}>
              <h3>Select your car *</h3>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-carId">Vehicle</label>
                <select
                  id="modal-carId"
                  name="carId"
                  value={formData.carId}
                  onChange={handleChange}
                  className={errors.carId ? formStyles.inputError : ''}
                >
                  <option value="">Choose a car…</option>
                  {cars.map((car) => (
                    <option key={car.id} value={car.id}>
                      {car.year} {car.make} {car.model} — {car.licensePlate}
                    </option>
                  ))}
                </select>
                {errors.carId ? <span className={formStyles.errorText}>{errors.carId}</span> : null}
              </div>

              {selectedCar ? (
                <div className={formStyles.selectedCarInfo}>
                  <h4>Selected car</h4>
                  <p>
                    {selectedCar.year} {selectedCar.make} {selectedCar.model} — {selectedCar.licensePlate}
                  </p>
                </div>
              ) : null}
            </div>

            <div className={formStyles.section}>
              <h3>Service type *</h3>
              <div className={formStyles.serviceGrid}>
                {SERVICE_TYPES.map((service) => (
                  <label key={service.id} className={formStyles.serviceOption}>
                    <input
                      type="radio"
                      name="service"
                      value={service.id}
                      checked={formData.service === service.id}
                      onChange={handleChange}
                    />
                    <div className={formStyles.serviceContent}>
                      <h4>{service.name}</h4>
                      <p>{service.description}</p>
                    </div>
                  </label>
                ))}
              </div>
              {errors.service ? <span className={formStyles.errorText}>{errors.service}</span> : null}

              {formData.service === 'other' ? (
                <div className={formStyles.formGroup}>
                  <label htmlFor="modal-customService">Specify service *</label>
                  <input
                    id="modal-customService"
                    type="text"
                    value={customServiceType}
                    onChange={(e) => setCustomServiceType(e.target.value)}
                    placeholder="Describe the service"
                    className={errors.customServiceType ? formStyles.inputError : ''}
                  />
                  {errors.customServiceType ? (
                    <span className={formStyles.errorText}>{errors.customServiceType}</span>
                  ) : null}
                </div>
              ) : null}
            </div>

            <div className={formStyles.section}>
              <h3>Fecha y horario *</h3>
              <p className={formStyles.helpText}>
                Lunes a viernes, cada 30 minutos: 9:30 a 12:30 y 14:00 a 17:30 (horario en formato 24 h).
              </p>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-appointment-day">Día</label>
                <input
                  id="modal-appointment-day"
                  name="appointmentDay"
                  type="date"
                  value={appointmentDay}
                  min={todayLocalYmd()}
                  onChange={(e) => {
                    const v = e.target.value;
                    setAppointmentDay(v);
                    setErrors((prev) => ({ ...prev, date: '', time: '' }));
                  }}
                  className={errors.date ? formStyles.inputError : ''}
                />
                {errors.date ? <span className={formStyles.errorText}>{errors.date}</span> : null}
              </div>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-appointment-slot">Horario</label>
                <select
                  id="modal-appointment-slot"
                  value={appointmentSlot}
                  onChange={(e) => {
                    setAppointmentSlot(e.target.value);
                    setErrors((prev) => ({ ...prev, time: '', date: prev.date }));
                  }}
                  disabled={!appointmentDay || bookableSlots.length === 0}
                  className={errors.time ? formStyles.inputError : ''}
                >
                  <option value="">
                    {!appointmentDay
                      ? 'Primero elegí un día…'
                      : bookableSlots.length === 0
                        ? 'No hay horarios para ese día'
                        : 'Elegí un horario…'}
                  </option>
                  {bookableSlots.map((hhmm) => (
                    <option key={hhmm} value={hhmm}>
                      {formatSlotLabel24h(hhmm)}
                    </option>
                  ))}
                </select>
                {errors.time ? <span className={formStyles.errorText}>{errors.time}</span> : null}
              </div>
            </div>

            <div className={formStyles.section}>
              <h3>Additional information</h3>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-notes">Notes</label>
                <textarea
                  id="modal-notes"
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  placeholder="Concerns, parts, or instructions…"
                  rows={4}
                  className={formStyles.textarea}
                />
              </div>
            </div>

            {selectedCar && selectedService && appointmentDay && appointmentSlot ? (
              <div className={formStyles.appointmentSummary}>
                <h3>Summary</h3>
                <div className={formStyles.summaryContent}>
                  <p>
                    <strong>Car:</strong> {selectedCar.year} {selectedCar.make} {selectedCar.model}
                  </p>
                  <p>
                    <strong>Service:</strong> {selectedService.name}
                  </p>
                  <p>
                    <strong>Cuándo:</strong> {formatAppointmentSummaryLocal(appointmentDay, appointmentSlot)}
                  </p>
                </div>
              </div>
            ) : null}

            <div className={formStyles.formActions}>
              <button type="button" onClick={onClose} className={formStyles.cancelButton}>
                Cancel
              </button>
              <button type="submit" disabled={isCreating} className={formStyles.submitButton}>
                {isCreating ? 'Scheduling…' : 'Schedule appointment'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

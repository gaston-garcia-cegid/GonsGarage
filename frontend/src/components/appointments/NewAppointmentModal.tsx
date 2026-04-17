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
      newErrors.carId = 'Selecione um automóvel';
    }
    if (!formData.service) {
      newErrors.service = 'Selecione o tipo de serviço';
    }
    if (formData.service === 'other' && !customServiceType.trim()) {
      newErrors.customServiceType = 'Indique o tipo de serviço';
    }
    if (!appointmentDay) {
      newErrors.date = 'Escolha um dia';
    } else if (!isWeekdayLocalYmd(appointmentDay)) {
      newErrors.date = weekdayErrorMessage();
    } else if (bookableSlots.length === 0) {
      newErrors.date = 'Não há horários disponíveis para esse dia. Escolha outro dia.';
    }
    if (!appointmentSlot) {
      newErrors.time = 'Escolha um horário';
    }
    if (!newErrors.date && !newErrors.time && appointmentDay && appointmentSlot) {
      const combined = combineLocalDateAndSlot(appointmentDay, appointmentSlot);
      const selectedDate = new Date(combined);
      const now = new Date();
      if (Number.isNaN(selectedDate.getTime())) {
        newErrors.date = 'Data ou hora inválida';
      } else if (selectedDate <= now) {
        newErrors.date = 'A marcação tem de ser no futuro';
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
              <h2 id="new-apt-title">Marcar serviço</h2>
              <p>Reserve uma visita para o seu automóvel</p>
            </div>
            <button type="button" className={modalStyles.closeButton} onClick={onClose} aria-label="Fechar">
              ×
            </button>
          </div>

          {error ? <div className={formStyles.errorAlert}>{error}</div> : null}

          <form onSubmit={handleSubmit} className={formStyles.appointmentForm}>
            <div className={formStyles.section}>
              <h3>Selecione o automóvel *</h3>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-carId">Automóvel</label>
                <select
                  id="modal-carId"
                  name="carId"
                  value={formData.carId}
                  onChange={handleChange}
                  className={errors.carId ? formStyles.inputError : ''}
                >
                  <option value="">Escolha um automóvel…</option>
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
                  <h4>Automóvel selecionado</h4>
                  <p>
                    {selectedCar.year} {selectedCar.make} {selectedCar.model} — {selectedCar.licensePlate}
                  </p>
                </div>
              ) : null}
            </div>

            <div className={formStyles.section}>
              <h3>Tipo de serviço *</h3>
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
                  <label htmlFor="modal-customService">Especificar serviço *</label>
                  <input
                    id="modal-customService"
                    type="text"
                    value={customServiceType}
                    onChange={(e) => setCustomServiceType(e.target.value)}
                    placeholder="Descreva o serviço"
                    className={errors.customServiceType ? formStyles.inputError : ''}
                  />
                  {errors.customServiceType ? (
                    <span className={formStyles.errorText}>{errors.customServiceType}</span>
                  ) : null}
                </div>
              ) : null}
            </div>

            <div className={formStyles.section}>
              <h3>Data e horário *</h3>
              <p className={formStyles.helpText}>
                Segunda a sexta, de 30 em 30 minutos: 9:30–12:30 e 14:00–17:30 (formato 24 h).
              </p>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-appointment-day">Dia</label>
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
                <label htmlFor="modal-appointment-slot">Horário</label>
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
                      ? 'Primeiro escolha um dia…'
                      : bookableSlots.length === 0
                        ? 'Sem horários para esse dia'
                        : 'Escolha um horário…'}
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
              <h3>Informação adicional</h3>
              <div className={formStyles.formGroup}>
                <label htmlFor="modal-notes">Notas</label>
                <textarea
                  id="modal-notes"
                  name="notes"
                  value={formData.notes}
                  onChange={handleChange}
                  placeholder="Observações, peças ou instruções…"
                  rows={4}
                  className={formStyles.textarea}
                />
              </div>
            </div>

            {selectedCar && selectedService && appointmentDay && appointmentSlot ? (
              <div className={formStyles.appointmentSummary}>
                <h3>Resumo</h3>
                <div className={formStyles.summaryContent}>
                  <p>
                    <strong>Automóvel:</strong> {selectedCar.year} {selectedCar.make} {selectedCar.model}
                  </p>
                  <p>
                    <strong>Serviço:</strong> {selectedService.name}
                  </p>
                  <p>
                    <strong>Quando:</strong> {formatAppointmentSummaryLocal(appointmentDay, appointmentSlot)}
                  </p>
                </div>
              </div>
            ) : null}

            <div className={formStyles.formActions}>
              <button type="button" onClick={onClose} className={formStyles.cancelButton}>
                Cancelar
              </button>
              <button type="submit" disabled={isCreating} className={formStyles.submitButton}>
                {isCreating ? 'A marcar…' : 'Confirmar marcação'}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  );
}

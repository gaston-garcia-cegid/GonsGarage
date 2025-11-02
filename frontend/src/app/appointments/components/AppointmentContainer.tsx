'use client';

import React, { useState, useCallback } from 'react';
import { Appointment, CreateAppointmentRequest, UpdateAppointmentRequest } from '@/types/appointment';
import AppointmentList from '@/app/appointments/components/AppointmentList';
import AppointmentModal from '@/app/appointments/components/AppointmentModal';
import LoadingSpinner from '@/components/ui/Loading/LoadingSpinner';
import ErrorAlert from '@/components/ui/Error/ErrorAlert';
import ConfirmModal from '@/components/ui/Modal/ConfirmModal';
import EmptyCarState from '@/components/empty-states/EmptyCarState';
import { useRouter } from 'next/navigation';
import styles from './AppointmentContainer.module.css';
import { useAppointments } from '@/stores';
import EmptyAppointmentState from '@/components/empty-states/EmptyAppointmentState';

interface AppointmentsContainerProps {
  // Optional callbacks for parent components
  onAddAppointment?: (appointment: Appointment) => void;
  onUpdateAppointment?: (appointments: Appointment[]) => void;
  onDeleteAppointment?: (appointmentId: string) => void;
  
  // UI customization
  maxAppointments?: number;
  showHeader?: boolean;
  headerTitle?: string;
  headerSubtitle?: string;
  addButtonText?: string;
  className?: string;
  
  // Layout control
  renderHeader?: () => React.ReactNode;
  renderEmptyState?: () => React.ReactNode;
}

export default function AppointmentsContainer({
  onAddAppointment,
  onUpdateAppointment,
  onDeleteAppointment,
  maxAppointments,
  showHeader = true,
  headerTitle = 'Your Appointments',
  headerSubtitle,
  addButtonText = 'Add New Appointment',
  className = '',
  renderHeader,
  renderEmptyState,
}: AppointmentsContainerProps) {
  const router = useRouter();
  const { appointments, isLoading, error, createAppointment, updateAppointment, deleteAppointment, fetchAppointments } = useAppointments();

  // Modal states
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [editingAppointment, setEditingAppointment] = useState<Appointment | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  
  // Delete confirmation state
  const [deleteConfirmation, setDeleteConfirmation] = useState<{
    appointmentId: string | null;
    clientName: string;
    isOpen: boolean;
    carId: string | null;
    carName: string;
  }>({
    appointmentId: null,
    clientName: '',
    isOpen: false,
    carId: null,
    carName: '',
  });

  // ✅ Create appointment handler
  const handleCreateAppointment = useCallback(async (appointmentData: CreateAppointmentRequest): Promise<boolean> => {
    if (maxAppointments && appointments.length >= maxAppointments) {
      alert(`You can only register up to ${maxAppointments} appointments.`);
      return false;
    }

    setIsCreating(true);

    try {
      const success = await createAppointment(appointmentData);
      
      if (success) {
        await fetchAppointments();
        
        const newAppointment = appointments.find(appointment => 
          appointment.date === appointmentData.date &&
          appointment.time === appointmentData.time &&
          appointment.carId === appointmentData.carId
        );

        if (onAddAppointment && newAppointment) {
          onAddAppointment(newAppointment);
        }

        if (onUpdateAppointment) {
          onUpdateAppointment(appointments);
        }

        setShowCreateModal(false);
        return true;
      }

      return false;
    } catch (error) {
      console.error('Failed to create car:', error);
      return false;
    } finally {
      setIsCreating(false);
    }
  }, [createAppointment, fetchAppointments, appointments, maxAppointments, onAddAppointment, onUpdateAppointment]);

  // ✅ Update appointment handler
  const handleUpdateAppointment = useCallback(async (id: string, appointmentData: Partial<UpdateAppointmentRequest>): Promise<boolean> => {
    try {
      const success = await updateAppointment(id, appointmentData);
      
      if (success) {
        await fetchAppointments();
        
        if (onUpdateAppointment) {
          onUpdateAppointment(appointments);
        }

        setEditingAppointment(null);
        return true;
      }

      return false;
    } catch (error) {
      console.error('Failed to update car:', error);
      return false;
    }
  }, [updateAppointment, fetchAppointments, appointments, onUpdateAppointment]);

  // ✅ Delete appointment handler - opens confirmation modal
  const handleDeleteAppointment = useCallback((id: string) => {
    const appointment = appointments.find(a => a.id === id);
    if (!appointment) return;

    setDeleteConfirmation({
      appointmentId: id,
      clientName: appointment.clientName,
      isOpen: true,
      carId: id,
      carName: `${appointment.clientName} (${appointment.carId})`,
    });
  }, [appointments]);

  // ✅ Confirm deletion
  const confirmDelete = useCallback(async () => {
    const { appointmentId } = deleteConfirmation;
    if (!appointmentId) return;

    try {
      const success = await deleteAppointment(appointmentId);
      
      if (success) {
        if (onDeleteAppointment) {
          onDeleteAppointment(appointmentId);
        }

        if (onUpdateAppointment) {
          onUpdateAppointment(appointments.filter(appointment => appointment.id !== appointmentId));
        }
      } else {
        alert('Failed to delete appointment. Please try again.');
      }
    } catch (error) {
      console.error('Failed to delete car:', error);
      alert('An error occurred while deleting the appointment.');
    } finally {
      setDeleteConfirmation({ isOpen: false, appointmentId: null, clientName: '', carId: null, carName: '' });
    }
  }, [deleteConfirmation, deleteAppointment, appointments, onDeleteAppointment, onUpdateAppointment]);

  // ✅ Cancel deletion
  const cancelDelete = useCallback(() => {
    setDeleteConfirmation({ isOpen: false, appointmentId: null, clientName: '', carId: null, carName: ''});
  }, []);

  // ✅ Edit appointment handler
  const handleEditAppointment = useCallback((appointment: Appointment) => {
    setEditingAppointment(appointment);
  }, []);

  // ✅ Close modal
  const handleModalClose = useCallback(() => {
    setShowCreateModal(false);
    setEditingAppointment(null);
  }, []);

  // ✅ Open create modal
  const handleAddAppointmentClick = useCallback(() => {
    if (maxAppointments && appointments.length >= maxAppointments) {
      alert(`You have reached the maximum limit of ${maxAppointments} appointments.`);
      return;
    }
    setShowCreateModal(true);
  }, [appointments.length, maxAppointments]);

  // Loading state
  if (isLoading) {
    return (
      <div className="loading-container">
        <LoadingSpinner />
        <span>Loading cars...</span>
      </div>
    );
  }

  // Empty state
  if (!isLoading && appointments.length === 0) {
    return renderEmptyState ? renderEmptyState() : <EmptyAppointmentState />;
  }

  const canAddAppointments = !maxAppointments || appointments.length < maxAppointments;
  const computedSubtitle = headerSubtitle || (
    maxAppointments 
      ? `Manage your registered appointments (${appointments.length}/${maxAppointments})`
      : 'Manage your registered appointments'
  );

  return (
    <div className={className}>
      {/* Custom header or default */}
      {renderHeader ? (
        renderHeader()
      ) : showHeader ? (
        <div className={`${styles.carsContainer} ${className}`}>
          <div className={styles.header}>
            <h2>{headerTitle}</h2>
            <p>{computedSubtitle}</p>
          </div>
          <button 
            onClick={handleAddAppointmentClick}
            className="primary-button"
            disabled={isCreating || !canAddAppointments}
          >
            <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} 
                d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            {isCreating ? 'Adding...' : addButtonText}
          </button>

          {maxAppointments && appointments.length >= maxAppointments && (
            <div className="max-cars-message">
              <span>Maximum cars limit reached</span>
            </div>
          )}
        </div>
      ) : null}

      {/* Error display */}
      {error && <ErrorAlert message={error} />}

      {/* Appointment list */}
      <AppointmentList
        appointments={appointments}
        onEdit={handleEditAppointment}
        onDelete={handleDeleteAppointment}
        onViewDetails={(id) => router.push(`/appointments/${id}`)}
        onScheduleService={(id) => router.push(`/appointments/new?appointmentId=${id}`)}
      />

      {/* Car Modal */}
      {/* {(showCreateModal || editingCar) && (
        <CarModal
          car={editingCar}
          onClose={handleModalClose}
          onCreate={handleCreateCar}
          onUpdate={handleUpdateCar}
        />
      )} */}

      {/* Confirmation Modal */}
      <ConfirmModal
        isOpen={deleteConfirmation.isOpen}
        title="Delete Car"
        message={`Are you sure you want to delete ${deleteConfirmation.carName}? This action cannot be undone.`}
        confirmText="Delete"
        cancelText="Cancel"
        variant="danger"
        onConfirm={confirmDelete}
        onCancel={cancelDelete}
      />
    </div>
  );
}
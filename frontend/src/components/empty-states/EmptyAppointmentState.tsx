'use client';

import React from 'react';
import Link from 'next/link';

interface EmptyAppointmentStateProps {
  hasNoCars?: boolean;
  onScheduleAppointment?: () => void;
}

export default function EmptyAppointmentState({ 
  hasNoCars = false, 
  onScheduleAppointment 
}: EmptyAppointmentStateProps) {
  if (hasNoCars) {
    return (
      <div className="empty-state">
        <div className="empty-state-content">
          <div className="empty-state-icon">
            📅
          </div>
          <h3>No Cars Available</h3>
          <p>
            You need to add a car before you can schedule appointments.
          </p>
          <div className="empty-state-actions">
            <Link href="/cars/new" className="btn-primary">
              Add a Car First
            </Link>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="empty-state">
      <div className="empty-state-content">
        <div className="empty-state-icon">
          📅
        </div>
        <h3>No Appointments Scheduled</h3>
        <p>
          You don&apos;t have any appointments scheduled yet. 
          Book your first service appointment.
        </p>
        <div className="empty-state-actions">
          {onScheduleAppointment ? (
            <button 
              onClick={onScheduleAppointment}
              className="btn-primary"
            >
              Schedule Appointment
            </button>
          ) : (
            <Link href="/appointments/new" className="btn-primary">
              Schedule Appointment
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
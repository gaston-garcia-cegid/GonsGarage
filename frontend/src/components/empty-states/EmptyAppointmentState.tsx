'use client';

import React from 'react';
import styles from './EmptyState.module.css';

interface EmptyAppointmentStateProps {
  onSchedule?: () => void;
}

export default function EmptyAppointmentState({ onSchedule }: EmptyAppointmentStateProps) {
  return (
    <div className={styles.emptyState}>
      <div className={styles.emptyStateIcon}>
        📅
      </div>
      <h3 className={styles.emptyStateTitle}>No Appointments Yet</h3>
      <p className={styles.emptyStateDescription}>
        You haven&apos;t scheduled any service appointments. 
        Book your first appointment to keep your vehicle in top condition.
      </p>
      {onSchedule && (
        <button onClick={onSchedule} className={styles.emptyStateButton}>
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} 
              d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Schedule First Appointment
        </button>
      )}
    </div>
  );
}
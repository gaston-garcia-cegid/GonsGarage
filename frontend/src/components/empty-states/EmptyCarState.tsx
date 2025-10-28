'use client';

import React from 'react';
import Link from 'next/link';

interface EmptyCarStateProps {
  onAddCar?: () => void;
}

export default function EmptyCarState({ onAddCar }: EmptyCarStateProps) {
  return (
    <div className="empty-state">
      <div className="empty-state-content">
        <div className="empty-state-icon">
          🚗
        </div>
        <h3>No Cars Registered</h3>
        <p>
          You haven&apos;t added any cars to your account yet. 
          Add your first car to start booking appointments.
        </p>
        <div className="empty-state-actions">
          {onAddCar ? (
            <button 
              onClick={onAddCar}
              className="btn-primary"
            >
              Add Your First Car
            </button>
          ) : (
            <Link href="/cars/new" className="btn-primary">
              Add Your First Car
            </Link>
          )}
        </div>
      </div>
    </div>
  );
}
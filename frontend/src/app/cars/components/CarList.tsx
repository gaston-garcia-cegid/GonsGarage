import React from 'react';
import { Car } from '@/types/car';
import CarCard from './CarCard';
import styles from '../cars.module.css';

interface CarListProps {
  cars: Car[];
  onEdit: (car: Car) => void;
  onDelete: (id: string) => void;
  onViewDetails: (id: string) => void;
  onScheduleService: (id: string) => void;
}

// Car list component following Agent.md component conventions
export default function CarList({
  cars,
  onEdit,
  onDelete,
  onViewDetails,
  onScheduleService
}: CarListProps) {
  // Empty state - following Agent.md user experience
  if (cars.length === 0) {
    return (
      <div className={styles.emptyState}>
        <div className={styles.emptyIcon}>🚗</div>
        <h3>No cars registered yet</h3>
        <p>Add your first car to get started with our services</p>
      </div>
    );
  }

  return (
    <div className={styles.carsGrid}>
      {cars.map((car) => (
        <CarCard
          key={car.id}
          car={car}
          onEdit={() => onEdit(car)}
          onDelete={() => onDelete(car.id)}
          onViewDetails={() => onViewDetails(car.id)}
          onScheduleService={() => onScheduleService(car.id)}
        />
      ))}
    </div>
  );
}
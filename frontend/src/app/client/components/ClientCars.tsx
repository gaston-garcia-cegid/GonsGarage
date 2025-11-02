// src/components/client/ClientCars.tsx
'use client';

import React from 'react';
import { Car } from '@/types/car';
import CarsContainer from '@/app/cars/components/CarsContainer';
import styles from '../client.module.css';

interface ClientCarsProps {
  onAddCar?: (car: Car) => void;
  onUpdateCar?: (cars: Car[]) => void;
  showAddButton?: boolean;
  maxCars?: number;
}

export default function ClientCars({ 
  onAddCar, 
  onUpdateCar, 
  showAddButton = true, 
  maxCars 
}: ClientCarsProps) {
  return (
    <CarsContainer
      onAddCar={onAddCar}
      onUpdateCar={onUpdateCar}
      maxCars={maxCars}
      headerTitle="Your Cars"
      addButtonText="Add New Car"
      className={styles.clientCarsContainer}
      showHeader={showAddButton}
    />
  );
}
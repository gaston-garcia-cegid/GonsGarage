// src/components/client/ClientCars.tsx
'use client';

import styles from '../client.module.css';

export default function ClientCars(
  { onAddCar }: { onAddCar: () => void }) {
  return (
    <div className="bg-white rounded-lg shadow p-6">
      <div className="flex justify-between items-center mb-6">
        <h3 className="text-lg font-semibold">Your Vehicles</h3>
        <button onClick={onAddCar}
          className={styles.primaryButton}>
          Add New Car
        </button>
      </div>
      <p>Your cars will be listed here...</p>
    </div>
  );
}
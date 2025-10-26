// src/app/client/page.tsx
'use client';

import React, { useState, useCallback } from 'react';
import { useAuthGuard } from '@hooks/useAuthGuard';
import DashboardLayout from '@/components/layouts/DashboardLayout';
import ClientCars from './components/ClientCars';
import { useClientData } from './hooks/useClientData';
import { Car } from '@/shared/types';
import ClientDashboard from './components/ClientDashboard';
import ClientAppointments from './components/ClientAppointments';

type ActiveTab = 'dashboard' | 'cars' | 'appointments';

export default function ClientPage() {
  const { isLoading: authLoading, isAuthorized } = useAuthGuard('client');
  const { cars, repairs, appointments, loading, error } = useClientData();
  const [activeTab, setActiveTab] = useState<ActiveTab>('dashboard');

  // Move all hooks above any conditional returns
  // Add missing userCars state for callbacks to work
  const [userCars, setUserCars] = useState<Car[]>([]);

  // Handle when a car is added - following Agent.md callback patterns
  const handleAddCar = useCallback((newCar: Car) => {
  console.log('New car added:', newCar);
  setUserCars(prevCars => [...prevCars, newCar]);
  alert(`${newCar.year} ${newCar.make} ${newCar.model} has been added successfully!`);
}, []);

  // Handle when cars list is updated - following Agent.md state management
  const handleUpdateCar = useCallback((updatedCars: Car[]) => {
    setUserCars(updatedCars);

    // Optional: Sync with parent state or context
    // updateUserCarsContext(updatedCars);
  }, []);

  if (authLoading || loading || !isAuthorized) {
    return <div>Loading...</div>;
  }

  const navigationItems = [
    { key: 'dashboard', label: 'Dashboard', href: '#' },
    { key: 'cars', label: 'My Cars', href: '#' },
    { key: 'appointments', label: 'Appointments', href: '#' },
  ];

  const renderContent = () => {
    switch (activeTab) {
      case 'dashboard':
        return (
          <ClientDashboard 
            error={error}
            cars={cars}  // ✅ Direct use, no mapping needed
            recentRepairs={repairs}
            upcomingAppointments={appointments
              .filter(a => new Date(a.scheduledAt) > new Date())
            }
            onNavigate={(tab: string) => setActiveTab(tab as ActiveTab)}
          />
        );
      case 'cars':
        return <ClientCars 
        onAddCar={handleAddCar}
        onUpdateCar={handleUpdateCar}
        showAddButton={true}
        maxCars={5} />;
      case 'appointments':
        return <ClientAppointments onScheduleService={() => {}} />;
    }
  };

  return (
    <DashboardLayout
      title={activeTab.charAt(0).toUpperCase() + activeTab.slice(1)}
      subtitle="Customer Dashboard"
      activeTab={activeTab}
      navigationItems={navigationItems}
      onNavClick={(tab: string) => setActiveTab(tab as ActiveTab)}
    >
      {renderContent()}
    </DashboardLayout>
  );
}
// src/app/client/page.tsx
'use client';

import React, { useState } from 'react';
import { useAuthGuard } from '@/hooks/useAuthGuard';
import DashboardLayout from '@/components/layouts/DashboardLayout';
import ClientDashboard from './components/ClientDashboard';
import ClientCars from './components/ClientCars';
import ClientAppointments from './components/ClientAppointments';
import { useClientData } from './hooks/useClientData';

type ActiveTab = 'dashboard' | 'cars' | 'appointments';

export default function ClientPage() {
  const { isLoading: authLoading, isAuthorized } = useAuthGuard('client');
  const { cars, repairs, appointments, loading, error } = useClientData();
  const [activeTab, setActiveTab] = useState<ActiveTab>('dashboard');
  
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
            cars={cars}
            recentRepairs={repairs}
            upcomingAppointments={appointments.filter(a => new Date(a.scheduled_at) > new Date())}
            onNavigate={(tab: string) => setActiveTab(tab as ActiveTab)}
          />
        );
      case 'cars':
        return <ClientCars onAddCar={() => {}} />;
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
import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import ClientDashboard from './ClientDashboard';
import type { Car } from '@/types/car';
import type { Repair } from '@/shared/types';
import type { Appointment } from '@/types/appointment';

const baseCar = (overrides: Partial<Car> = {}): Car => ({
  id: 'car-1',
  make: 'Toyota',
  model: 'Camry',
  year: 2010,
  licensePlate: 'AA-00-BB',
  color: 'silver',
  ownerId: 'owner-1',
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
  ...overrides,
});

const baseRepair = (overrides: Partial<Repair> = {}): Repair => ({
  id: 'rep-1',
  car_id: 'car-1',
  description: 'Oil change',
  status: 'in_progress',
  cost: 120.5,
  created_at: '2026-01-02T00:00:00.000Z',
  ...overrides,
});

const baseAppointment = (overrides: Partial<Appointment> = {}): Appointment => ({
  id: 'appt-1',
  clientName: 'João',
  carId: 'car-1',
  service: 'Inspection',
  date: '2026-12-31',
  time: '10:00',
  status: 'scheduled',
  createdAt: '2026-01-01T00:00:00.000Z',
  updatedAt: '2026-01-01T00:00:00.000Z',
  ...overrides,
});

describe('ClientDashboard', () => {
  const onNavigate = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('shows empty cars CTA and navigates to cars tab when adding the first car', async () => {
    const user = userEvent.setup();
    render(
      <ClientDashboard
        error={null}
        cars={[]}
        recentRepairs={[]}
        upcomingAppointments={[]}
        onNavigate={onNavigate}
      />,
    );

    expect(screen.getByText('No cars registered yet')).toBeInTheDocument();

    await user.click(screen.getByRole('button', { name: 'Add Your First Car' }));
    expect(onNavigate).toHaveBeenCalledTimes(1);
    expect(onNavigate).toHaveBeenCalledWith('cars');
  });

  it('lists a car and navigates to cars from list actions', async () => {
    const user = userEvent.setup();
    const cars = [baseCar({ id: 'c-a', year: 2012, make: 'Honda', model: 'Civic', licensePlate: 'ZZ-99-YY' })];

    render(
      <ClientDashboard
        error={null}
        cars={cars}
        recentRepairs={[]}
        upcomingAppointments={[]}
        onNavigate={onNavigate}
      />,
    );

    expect(screen.getByText('2012 Honda Civic')).toBeInTheDocument();
    expect(screen.getByText('ZZ-99-YY')).toBeInTheDocument();

    const viewAllButtons = screen.getAllByRole('button', { name: 'View All' });
    expect(viewAllButtons).toHaveLength(2);

    await user.click(viewAllButtons[0]);
    expect(onNavigate).toHaveBeenCalledWith('cars');

    await user.click(screen.getByRole('button', { name: 'View' }));
    expect(onNavigate).toHaveBeenCalledWith('cars');
  });

  it('navigates to appointments from recent repairs header', async () => {
    const user = userEvent.setup();
    const repairs: Repair[] = [baseRepair()];

    render(
      <ClientDashboard
        error={null}
        cars={[baseCar()]}
        recentRepairs={repairs}
        upcomingAppointments={[baseAppointment()]}
        onNavigate={onNavigate}
      />,
    );

    const viewAllButtons = screen.getAllByRole('button', { name: 'View All' });
    await user.click(viewAllButtons[1]);
    expect(onNavigate).toHaveBeenCalledWith('appointments');
  });
});

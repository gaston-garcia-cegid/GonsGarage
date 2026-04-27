import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import WorkshopListPage from './page';
import { UserRole } from '@/types';

const sampleCar = {
  id: 'cccccccc-cccc-cccc-cccc-cccccccccccc',
  make: 'Toyota',
  model: 'Corolla',
  year: 2020,
  license_plate: 'AA-00-BB',
  color: 'blue',
  owner_id: '11111111-1111-1111-1111-111111111111',
  created_at: '2020-01-01T00:00:00.000Z',
  updated_at: '2020-01-01T00:00:00.000Z',
};

const { getCarsMock, listOpenedMock, listByCarMock, createJobMock, pushMock } = vi.hoisted(() => ({
  getCarsMock: vi.fn(),
  listOpenedMock: vi.fn(),
  listByCarMock: vi.fn(),
  createJobMock: vi.fn(),
  pushMock: vi.fn(),
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({
    push: pushMock,
    replace: vi.fn(),
  }),
}));

vi.mock('@/hooks/useAuthHydrationReady', () => ({
  useAuthHydrationReady: () => true,
}));

vi.mock('@/stores', () => ({
  useAuth: () => ({
    user: {
      id: '22222222-2222-2222-2222-222222222222',
      email: 'mech@test.com',
      firstName: 'Mec',
      lastName: 'Técnico',
      role: UserRole.EMPLOYEE,
      createdAt: '2020-01-01T00:00:00.000Z',
      updatedAt: '2020-01-01T00:00:00.000Z',
    },
    logout: vi.fn(),
  }),
}));

vi.mock('@/lib/api', () => ({
  apiClient: {
    getCars: (...args: unknown[]) => getCarsMock(...args),
    listServiceJobsByOpenedOn: (...args: unknown[]) => listOpenedMock(...args),
    listServiceJobsByCar: (...args: unknown[]) => listByCarMock(...args),
    createServiceJob: (...args: unknown[]) => createJobMock(...args),
  },
}));

describe('WorkshopListPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    getCarsMock.mockResolvedValue({ data: [sampleCar], error: undefined });
    listOpenedMock.mockResolvedValue({ data: [], error: undefined });
    listByCarMock.mockResolvedValue({ data: [], error: undefined });
    createJobMock.mockResolvedValue({
      data: {
        id: 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb',
        car_id: sampleCar.id,
        status: 'open' as const,
        opened_by_user_id: '22222222-2222-2222-2222-222222222222',
        opened_at: '2020-01-02T10:00:00.000Z',
        created_at: '2020-01-02T10:00:00.000Z',
        updated_at: '2020-01-02T10:00:00.000Z',
      },
      error: undefined,
    });
  });

  it('opens Nova visita dialog with selected vehicle summary', async () => {
    const user = userEvent.setup();
    render(<WorkshopListPage />);

    const novaBtn = await waitFor(() => {
      const b = screen.getByRole('button', { name: 'Nova visita' });
      expect(b).not.toBeDisabled();
      return b;
    });

    await user.click(novaBtn);

    const dialog = await waitFor(() => screen.getByRole('dialog'));
    expect(within(dialog).getByText(/Confirma criar uma visita para/i)).toBeInTheDocument();
    expect(within(dialog).getByText(/AA-00-BB/)).toBeInTheDocument();
    expect(within(dialog).getByText(/Toyota/)).toBeInTheDocument();
    expect(within(dialog).getByRole('button', { name: 'Cancelar' })).toBeInTheDocument();
  });

  it('calls createServiceJob and navigates to new visit after confirm', async () => {
    const user = userEvent.setup();
    render(<WorkshopListPage />);

    const novaBtn = await waitFor(() => {
      const b = screen.getByRole('button', { name: 'Nova visita' });
      expect(b).not.toBeDisabled();
      return b;
    });
    await user.click(novaBtn);

    await waitFor(() => {
      expect(screen.getByRole('button', { name: 'Criar visita' })).toBeInTheDocument();
    });
    await user.click(screen.getByRole('button', { name: 'Criar visita' }));

    await waitFor(() => {
      expect(createJobMock).toHaveBeenCalledWith(sampleCar.id);
    });
    expect(pushMock).toHaveBeenCalledWith('/workshop/bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb');
  });
});

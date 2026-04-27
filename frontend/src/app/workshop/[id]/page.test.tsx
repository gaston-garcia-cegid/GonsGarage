import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import WorkshopDetailPage from './page';
import { UserRole } from '@/types';

const { getServiceJobMock } = vi.hoisted(() => ({
  getServiceJobMock: vi.fn(),
}));

vi.mock('next/navigation', () => ({
  useParams: () => ({ id: 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa' }),
  useRouter: () => ({
    push: vi.fn(),
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
      lastName: 'Ânico',
      role: UserRole.EMPLOYEE,
      createdAt: '2020-01-01T00:00:00.000Z',
      updatedAt: '2020-01-01T00:00:00.000Z',
    },
    logout: vi.fn(),
  }),
}));

vi.mock('@/lib/api', () => ({
  apiClient: {
    getServiceJob: (...args: unknown[]) => getServiceJobMock(...args),
  },
}));

describe('WorkshopDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it('does not leave only “A carregar…” when GET succeeds without a usable job (invalid body)', async () => {
    getServiceJobMock.mockResolvedValue({ data: null });

    render(<WorkshopDetailPage />);

    // Single waitFor: Strict Mode can remount after an error pass; require stable final UI (no loading + error visible).
    await waitFor(() => {
      expect(screen.getByText(/dados da visita em falta/i)).toBeInTheDocument();
      expect(screen.queryByText(/^A carregar/i)).not.toBeInTheDocument();
    });
  });

  it('shows raw job status when API returns a valid detail', async () => {
    getServiceJobMock.mockResolvedValue({
      data: {
        job: {
          id: 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa',
          car_id: 'cccccccc-cccc-cccc-cccc-cccccccccccc',
          status: 'open',
          opened_by_user_id: '22222222-2222-2222-2222-222222222222',
          opened_at: '2020-01-01T12:00:00.000Z',
          created_at: '2020-01-01T12:00:00.000Z',
          updated_at: '2020-01-01T12:00:00.000Z',
        },
        repair_ids: [],
      },
    });

    render(<WorkshopDetailPage />);

    await waitFor(() => {
      expect(screen.getByText(/Estado:/i)).toBeInTheDocument();
      expect(screen.getByText('open')).toBeInTheDocument();
    });
  });
});

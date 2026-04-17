import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useAuthStore } from './auth.store';
import { UserRole } from '@/types';

describe('auth store — client role after login', () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.persist.clearStorage();
    useAuthStore.setState({
      user: null,
      token: null,
      isLoading: false,
      error: null,
      isAuthenticated: false,
    });
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('maps auth/me user with role client', async () => {
    vi.stubGlobal(
      'fetch',
      vi.fn(async (input: RequestInfo | URL) => {
        const url = String(input);
        if (url.includes('/api/v1/auth/login')) {
          return {
            ok: true,
            json: async () => ({ token: 'fake-jwt' }),
          };
        }
        if (url.includes('/api/v1/auth/me')) {
          return {
            ok: true,
            json: async () => ({
              user: {
                id: '11111111-1111-1111-1111-111111111111',
                email: 'client@test.com',
                firstName: 'Test',
                lastName: 'Cliente',
                role: 'client',
              },
            }),
          };
        }
        return { ok: false, status: 404, text: async () => 'not found' };
      }),
    );

    const res = await useAuthStore.getState().login('client@test.com', 'password12');
    expect(res.success).toBe(true);

    const u = useAuthStore.getState().user;
    expect(u?.role).toBe(UserRole.CLIENT);
    expect(u?.email).toBe('client@test.com');
    expect(useAuthStore.getState().token).toBe('fake-jwt');
  });
});

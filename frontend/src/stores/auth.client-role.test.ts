import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { useAuthStore } from './auth.store';
import { UserRole } from '@/types';
import { apiClient } from '@/lib/api-client';
import { apiClient as legacyApiClient } from '@/lib/api';

describe('auth store — client role after login', () => {
  beforeEach(() => {
    localStorage.clear();
    useAuthStore.persist.clearStorage();
    apiClient.clearToken();
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
    vi.restoreAllMocks();
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
    expect(apiClient.getToken()).toBe('fake-jwt');
  });

  it('forwards legacy @/lib/api token helpers to the centralized api-client', () => {
    legacyApiClient.setToken('sync-check');
    expect(apiClient.getToken()).toBe('sync-check');
    legacyApiClient.clearToken();
    expect(apiClient.getToken()).toBeNull();
  });

  it('clears apiClient when the store logs out', async () => {
    const hrefSpy = vi.fn();
    vi.stubGlobal(
      'fetch',
      vi.fn(async (input: RequestInfo | URL) => {
        const url = String(input);
        if (url.includes('/api/v1/auth/login')) {
          return {
            ok: true,
            json: async () => ({ token: 'logout-test-jwt' }),
          };
        }
        if (url.includes('/api/v1/auth/me')) {
          return {
            ok: true,
            json: async () => ({
              user: {
                id: '33333333-3333-3333-3333-333333333333',
                email: 'logout@test.com',
                firstName: 'Out',
                lastName: 'User',
                role: 'client',
              },
            }),
          };
        }
        return { ok: false, status: 404, text: async () => 'not found' };
      }),
    );

    vi.stubGlobal('location', {
      ...window.location,
      assign: vi.fn(),
      replace: vi.fn(),
      reload: vi.fn(),
      get href() {
        return '';
      },
      set href(_v: string) {
        hrefSpy();
      },
    } as unknown as Location);

    await useAuthStore.getState().login('logout@test.com', 'password12');
    expect(apiClient.getToken()).toBe('logout-test-jwt');

    useAuthStore.getState().logout();
    expect(apiClient.getToken()).toBeNull();
    expect(hrefSpy).toHaveBeenCalled();
  });
});

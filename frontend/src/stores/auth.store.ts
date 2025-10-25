// AuthStore using Zustand following Agent.md patterns

import { create } from 'zustand';
import { immer } from 'zustand/middleware/immer';
import { persist } from 'zustand/middleware';
import { User, UserRole, RegisterRequest } from '@/types';

// ✅ Store state interface following Agent.md
interface AuthState {
  // Authentication state
  user: User | null;
  token: string | null;
  isLoading: boolean;
  error: string | null;
  
  // Derived state
  isAuthenticated: boolean;
}

// ✅ Store actions interface following Agent.md  
interface AuthActions {
  // Authentication methods (maintain same API as AuthContext)
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (data: RegisterRequest) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  
  // State management
  setUser: (user: User | null) => void;
  setToken: (token: string | null) => void;
  setLoading: (loading: boolean) => void;
  setError: (error: string | null) => void;
  clearError: () => void;
  
  // Initialization
  checkAuthStatus: () => Promise<void>;
}

// ✅ Complete store type
type AuthStore = AuthState & AuthActions;

// ✅ API configuration per Agent.md
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// ✅ Helper function for role-based redirects (maintain existing behavior)
const redirectBasedOnRole = (userData: User) => {
  console.log('Redirecting user with role:', userData.role);
  
  if (typeof window === 'undefined') return;
  
  switch (userData.role) {
    case UserRole.ADMIN:
      console.log('Redirecting to admin dashboard');
      window.location.href = '/admin/dashboard';
      break;
    case UserRole.EMPLOYEE:
      console.log('Redirecting to employee dashboard');  
      window.location.href = '/employee/dashboard';
      break;
    case UserRole.CLIENT:
      console.log('Redirecting to client dashboard');
      window.location.href = '/client/';
      break;
    default:
      console.warn('Unknown role:', userData.role, 'redirecting to default dashboard');
      window.location.href = '/';
  }
};

// ✅ Helper function for localStorage operations
const storage = {
  setAuthData: (token: string, user: User) => {
    if (typeof window !== 'undefined') {
      localStorage.setItem('auth_token', token);
      localStorage.setItem('auth_user', JSON.stringify(user));
    }
  },
  
  getAuthData: (): { token: string | null; user: User | null } => {
    if (typeof window === 'undefined') {
      return { token: null, user: null };
    }
    
    try {
      const token = localStorage.getItem('auth_token');
      const userData = localStorage.getItem('auth_user');
      
      if (token && userData) {
        const user = JSON.parse(userData);
        return { token, user };
      }
    } catch (error) {
      console.error('Failed to parse stored auth data:', error);
    }
    
    return { token: null, user: null };
  },
  
  clearAuthData: () => {
    if (typeof window !== 'undefined') {
      localStorage.removeItem('auth_token');
      localStorage.removeItem('auth_user');
    }
  }
};

// ✅ API functions following Agent.md
const authAPI = {
  login: async (email: string, password: string) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ email, password }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      let errorMessage = `Login failed (${response.status})`;
      
      try {
        const errorData = JSON.parse(errorText);
        errorMessage = errorData.error || errorMessage;
      } catch {
        // Keep the default errorMessage
      }
      
      throw new Error(errorMessage);
    }

    return await response.json();
  },

  register: async (data: RegisterRequest) => {
    const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email: data.email,
        password: data.password,
        firstName: data.firstName, // ✅ camelCase per Agent.md
        lastName: data.lastName,   // ✅ camelCase per Agent.md
        role: data.role || UserRole.CLIENT,
      }),
    });

    if (!response.ok) {
      const errorText = await response.text();
      let errorMessage = `Registration failed (${response.status})`;
      
      try {
        const errorData = JSON.parse(errorText);
        errorMessage = errorData.error || errorData.message || errorMessage;
      } catch {
        // Keep the default errorMessage
      }
      
      throw new Error(errorMessage);
    }

    return await response.json();
  }
};

// ✅ Zustand store with Immer and Persist middleware
export const useAuthStore = create<AuthStore>()(
  persist(
    immer((set) => ({
      // ✅ Initial state
      user: null,
      token: null,
      isLoading: false,
      error: null,
      isAuthenticated: false,

      // ✅ Authentication methods (maintain AuthContext API)
      login: async (email: string, password: string) => {
        try {
          set((state) => {
            state.isLoading = true;
            state.error = null;
          });

          const data = await authAPI.login(email, password);

          if (data.token) {
            // ✅ Create unified User object per Agent.md
            const userData: User = {
              id: data?.user?.id || '1',
              email: data?.user?.email || email,
              firstName: data?.user?.firstName || data?.user?.first_name || email.split('@')[0], // ✅ Handle both formats
              lastName: data?.user?.lastName || data?.user?.last_name || 'User',  // ✅ Handle both formats
              role: data?.user?.role || UserRole.CLIENT,
              createdAt: data?.user?.createdAt || data?.user?.created_at || new Date().toISOString(), // ✅ Handle both formats
              updatedAt: data?.user?.updatedAt || data?.user?.updated_at || new Date().toISOString(), // ✅ Handle both formats
            };

            console.log('Setting user data:', userData);

            set((state) => {
              state.token = data.token;
              state.user = userData;
              state.isAuthenticated = true;
              state.isLoading = false;
            });

            // Store in localStorage
            storage.setAuthData(data.token, userData);

            // Redirect based on role
            setTimeout(() => {
              redirectBasedOnRole(userData);
            }, 100);

            return { success: true };
          } else {
            set((state) => {
              state.error = data.error || 'Login failed';
              state.isLoading = false;
            });
            return { success: false, error: data.error || 'Login failed' };
          }
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Network error. Please try again.';
          set((state) => {
            state.error = errorMessage;
            state.isLoading = false;
          });
          return { success: false, error: errorMessage };
        }
      },

      register: async (data: RegisterRequest) => {
        try {
          set((state) => {
            state.isLoading = true;
            state.error = null;
          });

          await authAPI.register(data);

          set((state) => {
            state.isLoading = false;
          });

          return { success: true };
        } catch (error) {
          const errorMessage = error instanceof Error ? error.message : 'Network error. Please try again.';
          set((state) => {
            state.error = errorMessage;
            state.isLoading = false;
          });
          return { success: false, error: errorMessage };
        }
      },

      logout: () => {
        set((state) => {
          state.user = null;
          state.token = null;
          state.isAuthenticated = false;
          state.error = null;
        });

        // Clear localStorage
        storage.clearAuthData();

        // Redirect to landing page
        if (typeof window !== 'undefined') {
          window.location.href = '/';
        }
      },

      // ✅ State management methods
      setUser: (user: User | null) => {
        set((state) => {
          state.user = user;
          state.isAuthenticated = !!user && !!state.token;
        });
      },

      setToken: (token: string | null) => {
        set((state) => {
          state.token = token;
          state.isAuthenticated = !!state.user && !!token;
        });
      },

      setLoading: (loading: boolean) => {
        set((state) => {
          state.isLoading = loading;
        });
      },

      setError: (error: string | null) => {
        set((state) => {
          state.error = error;
        });
      },

      clearError: () => {
        set((state) => {
          state.error = null;
        });
      },

      // ✅ Initialize auth status from localStorage
      checkAuthStatus: async () => {
        try {
          set((state) => {
            state.isLoading = true;
          });

          const { token, user } = storage.getAuthData();

          if (token && user) {
            set((state) => {
              state.token = token;
              state.user = user;
              state.isAuthenticated = true;
            });
          } else {
            set((state) => {
              state.isAuthenticated = false;
            });
          }
        } catch (error) {
          console.error('Auth check failed:', error);
          storage.clearAuthData();
          set((state) => {
            state.user = null;
            state.token = null;
            state.isAuthenticated = false;
            state.error = 'Authentication check failed';
          });
        } finally {
          set((state) => {
            state.isLoading = false;
          });
        }
      },
    })),
    {
      name: 'gons-garage-auth', // ✅ Persistent storage key
      partialize: (state) => ({
        user: state.user,
        token: state.token,
        // Don't persist loading states or errors
      }),
    }
  )
);

// ✅ Convenience hooks following Agent.md patterns
export const useAuth = () => {
  const store = useAuthStore();
  return {
    // State
    user: store.user,
    token: store.token,
    isLoading: store.isLoading,
    error: store.error,
    isAuthenticated: store.isAuthenticated,
    
    // Actions (maintain AuthContext API)
    login: store.login,
    register: store.register,
    logout: store.logout,
    clearError: store.clearError,
  };
};

// ✅ Additional convenience hooks
export const useUser = () => useAuthStore((state) => state.user);
export const useAuthToken = () => useAuthStore((state) => state.token);
export const useAuthLoading = () => useAuthStore((state) => state.isLoading);
export const useAuthError = () => useAuthStore((state) => state.error);
export const useIsAuthenticated = () => useAuthStore((state) => state.isAuthenticated);
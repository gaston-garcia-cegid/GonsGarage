// Stores index - Central exports following Agent.md

// ✅ Authentication store
export {
  useAuthStore,
  useAuth,
  useUser,
  useAuthToken,
  useAuthLoading,
  useAuthError,
  useIsAuthenticated
} from './auth.store';

// ✅ Store types available for advanced usage
// AuthStore type is internal to the store implementation

// ✅ Re-export types for convenience
export type {
  User,
  RegisterRequest,
  LoginRequest,
  LoginResponse
} from '@/types';

export { UserRole } from '@/types';
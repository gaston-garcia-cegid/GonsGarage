import { useContext } from 'react';
import { AuthContext } from '@/contexts/AuthContext';
import { carApi } from '@/lib/api/carApi';
import React from 'react';

// ✅ Hook to integrate API client with authentication context
export function useApiClient() {
  const authContext = useContext(AuthContext) as AuthContextType;

  // ✅ Sync token with API client when auth context changes
  React.useEffect(() => {
    if (authContext?.token) {
      carApi.setAuthToken(authContext.token);
    } else {
      carApi.clearAuthToken();
    }
  }, [authContext?.token]);

  return carApi;
}
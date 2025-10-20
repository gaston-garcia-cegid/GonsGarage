'use client';

import { createContext, useContext, useEffect, useState } from 'react';
import { apiClient, type User } from '@/src/lib/api';

interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (data: RegisterData) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  loading: boolean;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('user');
    
    if (token && userData) {
      try {
        const parsedUser = JSON.parse(userData);
        setUser(parsedUser);
        apiClient.setToken(token);
      } catch (error) {
        console.error('Failed to parse user data:', error);
        localStorage.removeItem('token');
        localStorage.removeItem('user');
      }
    }
    
    setLoading(false);
  }, []);

  const login = async (email: string, password: string) => {
    const loginRequest = { email, password };
    const response = await apiClient.login(loginRequest);
    
    if (response.error) {
      // Ensure error is always a string
      const errorMsg = typeof response.error === 'string'
        ? response.error
        : response.error?.message || 'An error occurred';
      return { success: false, error: errorMsg };
    }

    if (response.data) {
      setUser(response.data.user);
      apiClient.setToken(response.data.token);
      localStorage.setItem('user', JSON.stringify(response.data.user));
      localStorage.setItem('token', response.data.token);
      return { success: true };
    }

    return { success: false, error: 'Unknown error occurred' };
  };

  const register = async (data: RegisterData): Promise<{ success: boolean; error?: string }> => {
    try {
      const requestUrl = `${process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'}/api/v1/auth/register`;
      const requestBody = {
        email: data.email,
        password: data.password,
        first_name: data.firstName,
        last_name: data.lastName,
        role: data.role,
      };

      console.log('Attempting registration to:', requestUrl);
      console.log('Request body:', requestBody);

      const response = await fetch(requestUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(requestBody),
      });

      console.log('Register response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Registration failed:', response.status, errorText);
        
        try {
          const errorData = JSON.parse(errorText);
          return { 
            success: false, 
            error: errorData.error || errorData.message || `Registration failed (${response.status})` 
          };
        } catch {
          return { success: false, error: `Registration failed (${response.status})` };
        }
      }

      const responseData = await response.json();
      console.log('Registration successful:', responseData);

      return { success: true };
    } catch (error) {
      console.error('Registration error:', error);
      return { 
        success: false, 
        error: error instanceof Error ? error.message : 'Network error. Please try again.' 
      };
    }
  };

  const logout = () => {
    setUser(null);
    //apiClient.removeToken();
    localStorage.removeItem('user');
  };

  return (
    <AuthContext.Provider value={{ 
      user,
      login,
      register, 
      logout, 
      loading,
      isAuthenticated: !!user }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
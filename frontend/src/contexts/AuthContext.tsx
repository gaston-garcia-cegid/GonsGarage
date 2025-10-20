'use client';

import React, { createContext, useContext, useEffect, useState, ReactNode } from 'react';

interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: string;
  isActive: boolean;
}

interface RegisterData {
  email: string;
  password: string;
  firstName: string;
  lastName: string;
  role: string;
}

interface AuthContextType {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<{ success: boolean; error?: string }>;
  register: (data: RegisterData) => Promise<{ success: boolean; error?: string }>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

interface AuthProviderProps {
  children: ReactNode;
}

export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    checkAuthStatus();
  }, []);

  const checkAuthStatus = async () => {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        setIsLoading(false);
        return;
      }

      apiClient.setToken(token);
      const { data, error } = await apiClient.getProfile();
      
      if (data && !error) {
        setUser(data);
      } else {
        localStorage.removeItem('token');
        apiClient.clearToken();
      }
    } catch (error) {
      console.error('Auth check failed:', error);
      localStorage.removeItem('token');
      apiClient.clearToken();
    } finally {
      setIsLoading(false);
    }
  };

  const login = async (email: string, password: string) : Promise<{ success: boolean; error?: string }> => {
    try {
      setIsLoading(true);
      const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email, password }),
      });

      console.log('Login response status:', response.status);

      if (!response.ok) {
        const errorText = await response.text();
        console.error('Login failed:', response.status, errorText);
        
        try {
          const errorData = JSON.parse(errorText);
          return { success: false, error: errorData.error || `Login failed (${response.status})` };
        } catch {
          return { success: false, error: `Login failed (${response.status})` };
        }
      }

      const data = await response.json();

      if (response.ok && data.token) {
        setToken(data.token);
        
        // Mock user data based on email for now
        const userData: User = {
          id: '1',
          email: email,
          firstName: email.split('@')[0],
          lastName: 'User',
          role: 'admin',
          isActive: true,
        };
        
        setUser(userData);
        
        // Store in localStorage (with window check)
        if (typeof window !== 'undefined') {
          localStorage.setItem('auth_token', data.token);
          localStorage.setItem('auth_user', JSON.stringify(userData));
        }
        
        return { success: true };
      } else {
        console.error('Login error:', data.error);
        return { success: false, error: data.error || 'Login failed' };
      }
    } catch (error) {
      console.log('Login error:', error);
      return { success: false, error: 'Network error occurred' };
    } finally {
      setIsLoading(false);
    }
  };

  const register = async (data: RegisterData): Promise<{ success: boolean; error?: string }> => {
    try {
      const requestUrl = `${API_BASE_URL}/api/v1/auth/register`;
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
      return { success: false, error: 'Network error. Please try again.' };
    }
  };

  const logout = async () => {
    try {
      await apiClient.logout();
    } catch (error) {
      console.error('Logout error:', error);
    } finally {
      setUser(null);
      apiClient.clearToken();
    }
  };

  const value: AuthContextType = {
    user,
    isLoading,
    login,
    logout,
    register,
    isAuthenticated: !!user,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}
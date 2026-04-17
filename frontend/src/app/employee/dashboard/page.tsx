'use client';

import React, { useEffect } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/stores';
import { UserRole } from '@/types';
import { getPostLoginPath } from '@/lib/post-login-paths';

export default function EmployeeDashboardPage() {
  const { user, isLoading, isAuthenticated } = useAuth();
  const router = useRouter();

  useEffect(() => {
    if (isLoading) return;
    if (!isAuthenticated || !user) {
      router.replace('/auth/login');
      return;
    }
    if (user.role !== UserRole.EMPLOYEE) {
      router.replace(getPostLoginPath(user.role));
    }
  }, [isLoading, isAuthenticated, user, router]);

  if (isLoading || !isAuthenticated || user?.role !== UserRole.EMPLOYEE) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-zinc-950 text-zinc-200">
        <p className="text-sm text-zinc-400">Loading…</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-zinc-950 text-zinc-100 p-8">
      <header className="max-w-2xl mx-auto mb-10">
        <h1 className="text-2xl font-semibold tracking-tight">Staff dashboard</h1>
        <p className="text-zinc-400 mt-1">
          Welcome, {user.firstName} {user.lastName}
        </p>
      </header>
      <nav className="max-w-2xl mx-auto flex flex-col gap-3">
        <Link
          href="/appointments"
          className="rounded-lg border border-zinc-800 bg-zinc-900 px-4 py-3 text-sm font-medium hover:bg-zinc-800 transition-colors"
        >
          Appointments
        </Link>
        <Link
          href="/cars"
          className="rounded-lg border border-zinc-800 bg-zinc-900 px-4 py-3 text-sm font-medium hover:bg-zinc-800 transition-colors"
        >
          Cars
        </Link>
      </nav>
    </div>
  );
}

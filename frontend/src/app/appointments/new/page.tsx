'use client';

import { Suspense, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';

/**
 * Legacy URL: /appointments/new?carId=…
 * Scheduling now lives in a modal on /appointments.
 */
function RedirectInner() {
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const carId = searchParams.get('carId');
    const q = carId
      ? `?schedule=1&carId=${encodeURIComponent(carId)}`
      : '?schedule=1';
    router.replace(`/appointments${q}`);
  }, [router, searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center">
      <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
    </div>
  );
}

export default function LegacyNewAppointmentPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600" />
        </div>
      }
    >
      <RedirectInner />
    </Suspense>
  );
}

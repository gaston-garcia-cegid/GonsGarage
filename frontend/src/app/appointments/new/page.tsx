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
    <div className="loadingScreen">
      <div className="spinnerLg" aria-hidden />
    </div>
  );
}

export default function LegacyNewAppointmentPage() {
  return (
    <Suspense
      fallback={
        <div className="loadingScreen">
          <div className="spinnerLg" aria-hidden />
        </div>
      }
    >
      <RedirectInner />
    </Suspense>
  );
}

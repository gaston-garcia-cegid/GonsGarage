'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

/** Legacy bookmark: opens create modal on the list via `?create=1`. */
export default function AdminPartsNewRedirectPage() {
  const router = useRouter();
  useEffect(() => {
    router.replace('/admin/parts?create=1');
  }, [router]);
  return null;
}

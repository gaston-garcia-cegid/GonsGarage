'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function NewBillingDocumentRedirectPage() {
  const router = useRouter();
  useEffect(() => {
    router.replace('/accounting/billing-documents?create=1');
  }, [router]);
  return null;
}

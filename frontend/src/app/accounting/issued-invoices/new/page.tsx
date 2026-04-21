'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function NewIssuedInvoiceRedirectPage() {
  const router = useRouter();
  useEffect(() => {
    router.replace('/accounting/issued-invoices?create=1');
  }, [router]);
  return null;
}

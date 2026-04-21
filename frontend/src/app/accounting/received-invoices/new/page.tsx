'use client';

import { useEffect } from 'react';
import { useRouter } from 'next/navigation';

export default function NewReceivedInvoiceRedirectPage() {
  const router = useRouter();
  useEffect(() => {
    router.replace('/accounting/received-invoices?create=1');
  }, [router]);
  return null;
}

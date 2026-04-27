import { fetchMyInvoicesInitialAuthenticated } from '@/lib/server/my-invoices-initial';
import MyInvoicesListClient from './MyInvoicesListClient';

export default async function MyInvoicesPage() {
  const initialItems = await fetchMyInvoicesInitialAuthenticated();
  return <MyInvoicesListClient initialItems={initialItems} />;
}

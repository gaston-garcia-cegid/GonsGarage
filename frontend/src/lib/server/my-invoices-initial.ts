import { cookies } from 'next/headers';
import { getPublicApiOrigin } from '@/lib/api-public-origin';
import type { IssuedInvoice } from '@/types/accounting';

function apiOriginBase(): string {
  return getPublicApiOrigin().replace(/\/+$/, '');
}

/**
 * Optional server-side initial list when `auth_token` is present as an **http** cookie
 * on the same request (e.g. future BFF / cookie-based session). Today login persists only
 * `localStorage`, so this usually returns `[]` and the client island performs `listMine`.
 */
export async function fetchMyInvoicesInitialAuthenticated(): Promise<IssuedInvoice[]> {
  const jar = await cookies();
  const token = jar.get('auth_token')?.value;
  if (!token) return [];

  const res = await fetch(`${apiOriginBase()}/api/v1/invoices/me?limit=50&offset=0`, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: 'application/json',
    },
    cache: 'no-store',
  });
  if (!res.ok) return [];

  const json = (await res.json()) as { items?: IssuedInvoice[] };
  return Array.isArray(json.items) ? json.items : [];
}

/**
 * Same contract as `issuedInvoiceService.get(id)` when a bearer cookie exists; otherwise `null`.
 */
export async function fetchMyInvoiceDetailAuthenticated(id: string): Promise<IssuedInvoice | null> {
  const jar = await cookies();
  const token = jar.get('auth_token')?.value;
  if (!token || !id) return null;

  const res = await fetch(`${apiOriginBase()}/api/v1/invoices/${encodeURIComponent(id)}`, {
    headers: {
      Authorization: `Bearer ${token}`,
      Accept: 'application/json',
    },
    cache: 'no-store',
  });
  if (!res.ok) return null;

  return (await res.json()) as IssuedInvoice;
}

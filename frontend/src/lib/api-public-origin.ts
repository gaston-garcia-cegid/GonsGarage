/**
 * Evita bloqueios de **Private Network Access** (Chrome): se a página abre-se desde
 * um hostname "público" (ex.: DuckDNS) mas `NEXT_PUBLIC_API_URL` aponta a um IP
 * privado (192.168.x, 10.x, …), o browser recusa `fetch` para esse IP.
 * Neste caso usamos a mesma origem da página (`/api/v1` via nginx).
 *
 * @see https://developer.chrome.com/blog/private-network-access-update
 */

export function isPrivateOrLocalhostHostname(hostname: string): boolean {
  const h = hostname.toLowerCase();
  if (h === 'localhost' || h.endsWith('.localhost')) return true;

  if (h.includes(':')) {
    if (h === '::1') return true;
    const he = h.replace(/^[\[]+|[\]]+$/g, '');
    if (he.startsWith('fe80:')) return true;
    if (he.startsWith('fc') || he.startsWith('fd')) return true;
    return false;
  }

  const m = /^(\d{1,3})\.(\d{1,3})\.(\d{1,3})\.(\d{1,3})$/.exec(h);
  if (!m) return false;
  const [a, b] = [Number(m[1]), Number(m[2])];
  if (a === 10) return true;
  if (a === 172 && b >= 16 && b <= 31) return true;
  if (a === 192 && b === 168) return true;
  if (a === 127) return true;
  if (a === 169 && b === 254) return true;
  if (a === 100 && b >= 64 && b <= 127) return true;
  return false;
}

/**
 * Origem (protocolo + host + porta) para onde o cliente deve enviar pedidos `/api/v1`.
 * No browser, substitui IPs privados pela origem da página quando necessário.
 */
export function getPublicApiOrigin(): string {
  const raw = (process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080').trim();
  const withoutApiPath = raw.replace(/\/+$/, '').replace(/\/api\/v1$/, '');
  let origin: string;
  try {
    origin = new URL(withoutApiPath).origin;
  } catch {
    origin = 'http://localhost:8080';
  }

  if (typeof window === 'undefined') {
    return origin;
  }

  try {
    const apiHost = new URL(origin).hostname;
    const pageHost = window.location.hostname;
    if (isPrivateOrLocalhostHostname(apiHost) && !isPrivateOrLocalhostHostname(pageHost)) {
      return window.location.origin;
    }
  } catch {
    /* manter origin */
  }

  return origin;
}

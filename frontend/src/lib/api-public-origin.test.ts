import { describe, it, expect, vi, afterEach } from 'vitest';
import { isPrivateOrLocalhostHostname, getPublicApiOrigin } from './api-public-origin';

describe('isPrivateOrLocalhostHostname', () => {
  it('treats RFC1918 and loopback as private', () => {
    expect(isPrivateOrLocalhostHostname('192.168.1.100')).toBe(true);
    expect(isPrivateOrLocalhostHostname('10.0.0.1')).toBe(true);
    expect(isPrivateOrLocalhostHostname('172.20.1.1')).toBe(true);
    expect(isPrivateOrLocalhostHostname('127.0.0.1')).toBe(true);
    expect(isPrivateOrLocalhostHostname('localhost')).toBe(true);
  });

  it('treats public hostnames as non-private', () => {
    expect(isPrivateOrLocalhostHostname('gonsgarage.duckdns.org')).toBe(false);
    expect(isPrivateOrLocalhostHostname('api.example.com')).toBe(false);
  });
});

describe('getPublicApiOrigin', () => {
  const env = process.env.NEXT_PUBLIC_API_URL;
  const originalLocation = window.location;

  afterEach(() => {
    process.env.NEXT_PUBLIC_API_URL = env;
    Object.defineProperty(window, 'location', {
      value: originalLocation,
      writable: true,
      configurable: true,
    });
  });

  it('replaces private API origin when page host is public (PNA-safe)', () => {
    process.env.NEXT_PUBLIC_API_URL = 'http://192.168.1.100:8102';
    Object.defineProperty(window, 'location', {
      value: {
        origin: 'http://gonsgarage.duckdns.org',
        hostname: 'gonsgarage.duckdns.org',
      },
      configurable: true,
    });
    expect(getPublicApiOrigin()).toBe('http://gonsgarage.duckdns.org');
  });

  it('keeps private API origin when page host is also private', () => {
    process.env.NEXT_PUBLIC_API_URL = 'http://192.168.1.100:8102';
    Object.defineProperty(window, 'location', {
      value: {
        origin: 'http://192.168.1.100:8102',
        hostname: '192.168.1.100',
      },
      configurable: true,
    });
    expect(getPublicApiOrigin()).toBe('http://192.168.1.100:8102');
  });
});

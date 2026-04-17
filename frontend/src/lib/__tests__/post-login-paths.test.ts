import fs from 'fs';
import path from 'path';
import { getPostLoginPath } from '@/lib/post-login-paths';
import { UserRole } from '@/types';

const appDir = path.join(__dirname, '../../app');

function pageFileForRoute(route: string): string {
  const normalized = route.replace(/\/+$/, '') || '/';
  if (normalized === '/') {
    return path.join(appDir, 'page.tsx');
  }
  const segments = normalized.split('/').filter(Boolean);
  return path.join(appDir, ...segments, 'page.tsx');
}

const POST_LOGIN_ROLE_CASES = [
  ...Object.values(UserRole).map((role) => ({ role, path: getPostLoginPath(role) })),
  { role: 'technician', path: getPostLoginPath('technician') },
  { role: 'unknown-role', path: getPostLoginPath('unknown-role') },
];

describe('getPostLoginPath', () => {
  it.each(POST_LOGIN_ROLE_CASES)('role $role -> $path', ({ role, path: expected }) => {
    expect(getPostLoginPath(role)).toBe(expected);
  });

  it('every mapped post-login path resolves to a real App Router page file', () => {
    const uniquePaths = [...new Set(POST_LOGIN_ROLE_CASES.map((c) => c.path))];
    for (const route of uniquePaths) {
      const file = pageFileForRoute(route);
      expect(fs.existsSync(file)).toBe(true);
    }
  });

  it('includes every UserRole in the contract table', () => {
    const rolesInTable = new Set(
      POST_LOGIN_ROLE_CASES.map((c) => c.role).filter((r): r is UserRole =>
        Object.values(UserRole).includes(r as UserRole),
      ),
    );
    for (const role of Object.values(UserRole)) {
      expect(rolesInTable.has(role)).toBe(true);
    }
  });
});

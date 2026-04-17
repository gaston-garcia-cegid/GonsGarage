import fs from 'fs';
import path from 'path';

/**
 * Documents `router.push` targets that still have no matching `page.tsx`.
 * When you implement a route, delete its entry here (or switch the assertion).
 * This prevents "silent" 404s like the old `/employee/dashboard` gap.
 */
const appDir = path.join(__dirname, '../../app');

describe('Internal navigation vs App Router files', () => {
  const routesStillWithoutPage = [
    {
      id: 'appointment-detail',
      relativeDir: path.join('appointments', '[id]'),
    },
    {
      id: 'appointment-edit',
      relativeDir: path.join('appointments', '[id]', 'edit'),
    },
  ];

  it.each(routesStillWithoutPage)(
    '$id: no page.tsx yet under $relativeDir',
    ({ relativeDir }) => {
      const pagePath = path.join(appDir, relativeDir, 'page.tsx');
      expect(fs.existsSync(pagePath)).toBe(false);
    },
  );
});

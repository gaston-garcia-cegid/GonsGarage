// src/lib/navigation.ts
export const getNavigationConfig = (role: string) => {
  const configs = {
    client: {
      basePath: '/client',
      items: [
        { key: 'dashboard', label: 'Dashboard', href: '/client' },        // ✅ Sem /dashboard
        { key: 'cars', label: 'My Cars', href: '/client/cars' },
        { key: 'appointments', label: 'Appointments', href: '/client/appointments' },
      ]
    },
    admin: {
      basePath: '/admin',
      items: [
        { key: 'dashboard', label: 'Dashboard', href: '/admin' },         // ✅ Sem /dashboard
        { key: 'users', label: 'Users', href: '/admin/users' },
        { key: 'reports', label: 'Reports', href: '/admin/reports' },
      ]
    },
    technician: {
      basePath: '/technician',
      items: [
        { key: 'dashboard', label: 'Dashboard', href: '/technician' },    // ✅ Sem /dashboard
        { key: 'tasks', label: 'My Tasks', href: '/technician/tasks' },
      ]
    }
  };

  return configs[role as keyof typeof configs] || configs.client;
};
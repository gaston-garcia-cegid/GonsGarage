import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { StaffUsersDashboardCta } from './StaffUsersDashboardCta';

describe('StaffUsersDashboardCta', () => {
  it('invokes onManageUsers when Gestão de utilizadores is clicked', async () => {
    const onManageUsers = vi.fn();
    const user = userEvent.setup();
    render(<StaffUsersDashboardCta onManageUsers={onManageUsers} />);

    await user.click(screen.getByRole('button', { name: 'Gestão de utilizadores' }));
    expect(onManageUsers).toHaveBeenCalledTimes(1);
  });

  it('shows supporting copy', () => {
    render(<StaffUsersDashboardCta onManageUsers={vi.fn()} />);
    expect(screen.getByText(/Criar ou atualizar contas/i)).toBeInTheDocument();
  });
});

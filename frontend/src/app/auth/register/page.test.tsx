import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import RegisterPage from './page';

vi.mock('next/image', () => ({
  default: (props: { src: string; alt: string }) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={props.src} alt={props.alt} data-testid="brand-logo" />
  ),
}));

vi.mock('next/navigation', () => ({
  useRouter: () => ({ push: vi.fn(), replace: vi.fn() }),
}));

vi.mock('@/stores', () => ({
  useAuth: () => ({
    register: vi.fn().mockResolvedValue({ success: false }),
    isAuthenticated: false,
  }),
}));

describe('RegisterPage', () => {
  it('uses the shared auth shell title and subtitle', () => {
    render(<RegisterPage />);
    expect(screen.getByRole('heading', { level: 1, name: 'Criar conta' })).toBeInTheDocument();
    expect(screen.getByText('Junte-se à equipa GonsGarage')).toBeInTheDocument();
    expect(screen.getByTestId('brand-logo')).toHaveAttribute('src', '/images/LogoGonsGarage.jpg');
  });

  it('shows the cross-link to login as a button', () => {
    render(<RegisterPage />);
    expect(screen.getByRole('button', { name: 'Iniciar sessão' })).toBeInTheDocument();
  });
});

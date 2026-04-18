import React from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
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

const { mockRegister } = vi.hoisted(() => ({
  mockRegister: vi.fn().mockResolvedValue({ success: false }),
}));

vi.mock('@/stores', () => ({
  useAuth: () => ({
    register: mockRegister,
    isAuthenticated: false,
  }),
}));

beforeEach(() => {
  mockRegister.mockReset();
  mockRegister.mockResolvedValue({ success: false });
});

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

  it('shows email field label and placeholder in Portuguese aligned with login', () => {
    render(<RegisterPage />);
    expect(screen.getByLabelText('E-mail')).toBeInTheDocument();
    expect(screen.getByPlaceholderText('O seu e-mail')).toBeInTheDocument();
  });

  it('shows confirm password label in Portuguese', () => {
    render(<RegisterPage />);
    expect(screen.getByLabelText('Confirmar palavra-passe')).toBeInTheDocument();
  });

  it('does not prefix unexpected errors with English "Error:"', async () => {
    const user = userEvent.setup();
    mockRegister.mockRejectedValueOnce(new Error('Falha de rede'));
    render(<RegisterPage />);

    await user.type(screen.getByLabelText('Nome'), 'Maria');
    await user.type(screen.getByLabelText('Apelido'), 'Silva');
    await user.type(screen.getByLabelText('E-mail'), 'maria@example.com');
    await user.type(screen.getByLabelText('Palavra-passe'), 'secret12');
    await user.type(screen.getByLabelText('Confirmar palavra-passe'), 'secret12');
    await user.click(screen.getByRole('button', { name: 'Criar conta' }));

    const alert = await screen.findByRole('alert');
    expect(alert).toHaveTextContent('Falha de rede');
    expect(alert.textContent).not.toMatch(/^Error:/);
    await waitFor(() => expect(mockRegister).toHaveBeenCalled());
  });
});

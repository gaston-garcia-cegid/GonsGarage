import React, { Suspense } from 'react';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { act, render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import LoginForm from './LoginForm';

vi.mock('next/image', () => ({
  default: (props: { src: string; alt: string }) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={props.src} alt={props.alt} data-testid="brand-logo" />
  ),
}));

const replace = vi.fn();
const push = vi.fn();

vi.mock('next/navigation', () => ({
  useRouter: () => ({ replace, push }),
  useSearchParams: () => new URLSearchParams('message=Registo%20ok'),
}));

const loginMock = vi.fn();

vi.mock('@/stores', () => ({
  useAuth: () => ({
    login: loginMock,
    isAuthenticated: false,
  }),
}));

function renderLogin() {
  return render(
    <Suspense fallback={null}>
      <LoginForm />
    </Suspense>
  );
}

describe('LoginForm', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    loginMock.mockResolvedValue({ success: true });
  });

  it('shows auth heading and subtitle from the shared shell', async () => {
    renderLogin();
    await act(async () => {
      await Promise.resolve();
    });
    expect(screen.getByRole('heading', { level: 1, name: 'Iniciar sessão' })).toBeInTheDocument();
    expect(screen.getByText('Aceda com a sua conta GonsGarage')).toBeInTheDocument();
  });

  it('shows registration success message from query in a status region', async () => {
    renderLogin();
    await waitFor(() => {
      expect(screen.getByRole('status')).toHaveTextContent('Registo ok');
    });
  });

  it('submits email and password to login and navigates to dashboard on success', async () => {
    const user = userEvent.setup();
    renderLogin();
    await act(async () => {
      await Promise.resolve();
    });

    await user.type(screen.getByRole('textbox', { name: /e-mail/i }), 'a@b.com');
    await user.type(screen.getByPlaceholderText('A sua palavra-passe'), 'secret12');
    await user.click(screen.getByRole('button', { name: /iniciar sessão/i }));

    expect(loginMock).toHaveBeenCalledWith('a@b.com', 'secret12');
    expect(replace).toHaveBeenCalledWith('/dashboard');
  });
});

import React from 'react';
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '@testing-library/react';
import { AuthShell, AuthShellFooter } from './AuthShell';

vi.mock('next/image', () => ({
  default: (props: { src: string; alt: string }) => (
    // eslint-disable-next-line @next/next/no-img-element
    <img src={props.src} alt={props.alt} data-testid="brand-logo" />
  ),
}));

describe('AuthShell', () => {
  it('renders the page title as the main heading', () => {
    render(
      <AuthShell title="Iniciar sessão" subtitle="Subtítulo">
        <p>Conteúdo</p>
      </AuthShell>
    );
    expect(screen.getByRole('heading', { level: 1, name: 'Iniciar sessão' })).toBeInTheDocument();
    expect(screen.getByText('Subtítulo')).toBeInTheDocument();
    expect(screen.getByText('Conteúdo')).toBeInTheDocument();
  });

  it('renders success banner with message in a status region', () => {
    render(
      <AuthShell
        title="Login"
        banner={{ variant: 'success', message: 'Registo concluído.' }}
      >
        <span>Form</span>
      </AuthShell>
    );
    const region = screen.getByRole('status');
    expect(region).toHaveTextContent('Registo concluído.');
  });

  it('renders error banner with message in an alert region', () => {
    render(
      <AuthShell
        title="Login"
        banner={{ variant: 'error', message: 'Credenciais inválidas' }}
      >
        <span>Form</span>
      </AuthShell>
    );
    const alert = screen.getByRole('alert');
    expect(alert).toHaveTextContent('Credenciais inválidas');
  });

  it('exposes brand logo with correct path for accessibility', () => {
    render(
      <AuthShell title="Título">
        <span>X</span>
      </AuthShell>
    );
    const img = screen.getByTestId('brand-logo');
    expect(img).toHaveAttribute('src', '/images/LogoGonsGarage.jpg');
    expect(img).toHaveAttribute('alt', 'Logótipo GonsGarage');
  });
});

describe('AuthShellFooter', () => {
  it('renders children inside the footer wrapper', () => {
    render(
      <AuthShellFooter>
        <button type="button">Iniciar sessão</button>
      </AuthShellFooter>
    );
    expect(screen.getByRole('button', { name: 'Iniciar sessão' })).toBeInTheDocument();
  });
});

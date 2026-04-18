'use client';

import React, { useState, useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/stores';
import { AuthShell, AuthShellFooter } from '@/components/auth/AuthShell';
import authShellStyles from '@/components/auth/AuthShell.module.css';
import { Input } from '@/components/ui/Input/Input';
import { Button } from '@/components/ui/Button/Button';
import styles from './login.module.css';

export default function LoginForm() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [errors, setErrors] = useState<{ [key: string]: string }>({});
  const [isLoading, setIsLoading] = useState(false);
  const [showPassword, setShowPassword] = useState(false);
  const [successMessage, setSuccessMessage] = useState('');

  const { login } = useAuth();
  const router = useRouter();
  const searchParams = useSearchParams();

  useEffect(() => {
    const message = searchParams.get('message');
    if (message) {
      setSuccessMessage(message);
    }
  }, [searchParams]);

  const validateForm = () => {
    const newErrors: { [key: string]: string } = {};

    if (!formData.email) {
      newErrors.email = 'O e-mail é obrigatório';
    } else if (!/\S+@\S+\.\S+/.test(formData.email)) {
      newErrors.email = 'E-mail inválido';
    }

    if (!formData.password) {
      newErrors.password = 'A palavra-passe é obrigatória';
    } else if (formData.password.length < 6) {
      newErrors.password = 'A palavra-passe deve ter pelo menos 6 caracteres';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validateForm()) {
      return;
    }

    setIsLoading(true);
    setErrors({});

    try {
      const result = await login(formData.email, formData.password);

      if (result.success) {
        router.replace('/dashboard');
      } else {
        setErrors({ general: result.error || 'Falha no início de sessão' });
      }
    } catch {
      setErrors({ general: 'Ocorreu um erro inesperado' });
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const { name, value } = e.target;
    setFormData((prev) => ({
      ...prev,
      [name]: value,
    }));

    if (errors[name]) {
      setErrors((prev) => ({
        ...prev,
        [name]: '',
      }));
    }
  };

  let banner: { variant: 'success' | 'error'; message: string } | null = null;
  if (errors.general) {
    banner = { variant: 'error', message: errors.general };
  } else if (successMessage) {
    banner = { variant: 'success', message: successMessage };
  }

  return (
    <AuthShell
      title="Iniciar sessão"
      subtitle="Aceda com a sua conta GonsGarage"
      banner={banner}
    >
      <form className={styles.form} onSubmit={handleSubmit}>
        <div className={styles.fieldStack}>
          <Input
            id="email"
            name="email"
            type="email"
            label="E-mail"
            autoComplete="email"
            required
            value={formData.email}
            onChange={handleChange}
            error={errors.email}
            placeholder="O seu e-mail"
          />

          <div className={styles.passwordField}>
            <label htmlFor="password" className={styles.passwordLabel}>
              Palavra-passe
            </label>
            <div className={styles.passwordInputWrap}>
              <input
                id="password"
                name="password"
                type={showPassword ? 'text' : 'password'}
                autoComplete="current-password"
                required
                value={formData.password}
                onChange={handleChange}
                className={`${styles.passwordInput} ${errors.password ? styles.passwordInputError : ''}`}
                placeholder="A sua palavra-passe"
                aria-invalid={Boolean(errors.password)}
                aria-describedby={errors.password ? 'password-error' : undefined}
              />
              <button
                type="button"
                className={styles.passwordToggle}
                onClick={() => setShowPassword(!showPassword)}
                aria-label={showPassword ? 'Ocultar palavra-passe' : 'Mostrar palavra-passe'}
              >
                <svg className={styles.eyeIcon} fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden>
                  {showPassword ? (
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L7.05 7.05m2.828 2.828L12 12m2.122 2.122L17.95 17.95"
                    />
                  ) : (
                    <>
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"
                      />
                    </>
                  )}
                </svg>
              </button>
            </div>
            {errors.password ? (
              <p id="password-error" className={styles.fieldError}>
                {errors.password}
              </p>
            ) : null}
          </div>
        </div>

        <Button type="submit" variant="primary" size="lg" loading={isLoading} className={styles.submitWide}>
          Iniciar sessão
        </Button>

        <AuthShellFooter>
          <p className={authShellStyles.footerText}>
            Ainda não tem conta?{' '}
            <button
              type="button"
              className={authShellStyles.footerLink}
              onClick={() => router.push('/auth/register')}
            >
              Criar conta
            </button>
          </p>
          <p className={styles.demoBlock}>
            Conta de demonstração: admin@gonsgarage.com / admin123
          </p>
        </AuthShellFooter>
      </form>
    </AuthShell>
  );
}

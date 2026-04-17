'use client';

import { useCallback, useEffect, useState } from 'react';
import { applyTheme, getStoredPreference, type ThemePreference } from '@/lib/theme';
import styles from './ThemeSwitcher.module.css';

export default function ThemeSwitcher() {
  const [pref, setPref] = useState<ThemePreference>('system');
  const [mounted, setMounted] = useState(false);

  useEffect(() => {
    setMounted(true);
    setPref(getStoredPreference());
  }, []);

  const syncFromStorage = useCallback(() => {
    setPref(getStoredPreference());
  }, []);

  useEffect(() => {
    if (!mounted) return;
    applyTheme(pref);
    if (pref !== 'system') return undefined;
    const mq = globalThis.matchMedia('(prefers-color-scheme: dark)');
    const onChange = () => applyTheme('system');
    mq.addEventListener('change', onChange);
    globalThis.addEventListener('storage', syncFromStorage);
    return () => {
      mq.removeEventListener('change', onChange);
      globalThis.removeEventListener('storage', syncFromStorage);
    };
  }, [pref, mounted, syncFromStorage]);

  if (!mounted) {
    return null;
  }

  return (
    <div className={styles.root}>
      <label className={styles.label} htmlFor="gg-theme-select">
        Tema
      </label>
      <select
        id="gg-theme-select"
        className={styles.select}
        value={pref}
        onChange={(e) => setPref(e.target.value as ThemePreference)}
        aria-label="Tema da aplicação"
      >
        <option value="system">Sistema</option>
        <option value="light">Claro</option>
        <option value="dark">Escuro</option>
      </select>
    </div>
  );
}

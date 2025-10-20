import React, { forwardRef } from 'react';
import styles from './Input.module.css';

interface InputProps extends React.InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  helperText?: string;
  variant?: 'default' | 'error';
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, helperText, variant = 'default', className, ...props }, ref) => {
    const inputVariant = error ? 'error' : variant;
    
    return (
      <div className={styles.inputWrapper}>
        {label && (
          <label htmlFor={props.id} className={styles.label}>
            {label}
          </label>
        )}
        <input
          ref={ref}
          className={`${styles.input} ${styles[inputVariant]} ${className || ''}`}
          {...props}
        />
        {error && (
          <span className={styles.errorMessage}>{error}</span>
        )}
        {helperText && !error && (
          <span className={styles.helperText}>{helperText}</span>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';
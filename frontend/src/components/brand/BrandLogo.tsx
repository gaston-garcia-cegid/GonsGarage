'use client';

import React, { useCallback, useState } from 'react';
import type { CSSProperties } from 'react';
import Image from 'next/image';

/** Preferido quando o asset existe em `public/images/`. */
export const BRAND_LOGO_JPEG_SRC = '/images/LogoGonsGarage.jpg';
/** Fallback versionado no repo (sempre disponível). */
export const BRAND_LOGO_FALLBACK_SRC = '/images/logo-gonsgarage-mark.svg';

export type BrandLogoProps = {
  width: number;
  height: number;
  alt: string;
  className?: string;
  style?: CSSProperties;
  priority?: boolean;
  /** Para testes ou telemetria; omitir em páginas com vários logótipos. */
  dataTestId?: string;
};

/**
 * Logótipo de marca: tenta o JPG da oficina; se falhar (404, clone sem asset), usa SVG em repo.
 * `unoptimized`: evita `/_next/image` — o loader por defeito pode 404 no JPG e tratar SVG local
 * de forma pouco fiável; `<img>` directo garante `onError` e fallback visível.
 */
export function BrandLogo({
  width,
  height,
  alt,
  className,
  style,
  priority,
  dataTestId,
}: Readonly<BrandLogoProps>) {
  const [src, setSrc] = useState(BRAND_LOGO_JPEG_SRC);

  const handleError = useCallback(() => {
    setSrc((current) => (current === BRAND_LOGO_FALLBACK_SRC ? current : BRAND_LOGO_FALLBACK_SRC));
  }, []);

  return (
    <Image
      src={src}
      alt={alt}
      width={width}
      height={height}
      className={className}
      style={style}
      priority={priority}
      unoptimized
      onError={handleError}
      data-testid={dataTestId}
    />
  );
}

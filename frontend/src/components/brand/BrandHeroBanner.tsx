'use client';

import React, { useCallback, useState } from 'react';
import Image from 'next/image';

const BANNER_JPEG_SRC = '/images/BannerGonsGarage.jpg';
const BANNER_FALLBACK_SRC = '/images/banner-gonsgarage-fallback.svg';

export type BrandHeroBannerProps = {
  alt: string;
  priority?: boolean;
};

/**
 * Hero full-bleed: tenta foto da oficina; se falhar, usa SVG de fallback no repo.
 * Ver `BrandLogo` — `unoptimized` para carregar ficheiros em `public/` sem pipeline `/_next/image`.
 */
export function BrandHeroBanner({ alt, priority }: Readonly<BrandHeroBannerProps>) {
  const [src, setSrc] = useState(BANNER_JPEG_SRC);

  const handleError = useCallback(() => {
    setSrc((current) => (current === BANNER_FALLBACK_SRC ? current : BANNER_FALLBACK_SRC));
  }, []);

  return (
    <Image
      src={src}
      alt={alt}
      fill
      style={{ objectFit: 'cover' }}
      priority={priority}
      unoptimized
      onError={handleError}
    />
  );
}

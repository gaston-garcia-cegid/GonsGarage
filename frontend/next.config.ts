import type { NextConfig } from "next";
import path from "node:path";
import { fileURLToPath } from "node:url";

// Turbopack (Next 16+) can infer the wrong workspace root (e.g. `src/app`) and then
// fail to resolve `next/package.json`. Pin the root to this app directory.
const projectRoot = path.dirname(fileURLToPath(import.meta.url));

const nextConfig: NextConfig = {
  turbopack: {
    root: projectRoot,
  },
};

export default nextConfig;

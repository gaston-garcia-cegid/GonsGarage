import type { NextConfig } from "next";

// Standalone solo en CI/Linux o build Docker (NEXT_STANDALONE=true); en Windows sin permisos de symlink `pnpm build` falla.
const nextConfig: NextConfig = {
  ...(process.env.NEXT_STANDALONE === "true" ? { output: "standalone" as const } : {}),
};

export default nextConfig;

import { readFileSync, existsSync } from "node:fs";
import path from "node:path";

const repoRoot = path.resolve(__dirname, "../..", "..");
const adrPath = path.join(
  repoRoot,
  "docs",
  "adr",
  "0003-nextjs16-react19-migration-main.md",
);
const packageJsonPath = path.join(repoRoot, "frontend", "package.json");
const postcssConfigPath = path.join(repoRoot, "frontend", "postcss.config.mjs");
const globalsCssPath = path.join(repoRoot, "frontend", "src", "app", "globals.css");
const tailwindConfigPath = path.join(repoRoot, "frontend", "tailwind.config.ts");
const authLoginPagePath = path.join(
  repoRoot,
  "frontend",
  "src",
  "app",
  "auth",
  "login",
  "page.tsx",
);
const myInvoicesListPagePath = path.join(
  repoRoot,
  "frontend",
  "src",
  "app",
  "my-invoices",
  "page.tsx",
);
const myInvoicesLayoutPath = path.join(
  repoRoot,
  "frontend",
  "src",
  "app",
  "my-invoices",
  "layout.tsx",
);

function readPackageJson(): {
  dependencies?: Record<string, string>;
  devDependencies?: Record<string, string>;
} {
  return JSON.parse(readFileSync(packageJsonPath, "utf8")) as {
    dependencies?: Record<string, string>;
    devDependencies?: Record<string, string>;
  };
}

/** First major version digit in a dep range (e.g. `^4.1.0` → 4, `16.2.4` → 16). */
function depMajor(versionSpec: string): number {
  const m = versionSpec.match(/(\d+)/);
  if (!m) return 0;
  return Number.parseInt(m[1]!, 10);
}

describe("SDD change nextjs-16-react19-migration — ADR on main", () => {
  it("documents scope, GO/NO-GO, and quality commands", () => {
    expect(existsSync(adrPath)).toBe(true);
    const body = readFileSync(adrPath, "utf8");
    expect(body).toContain("## Scope");
    expect(body).toMatch(/GO\s*\/\s*NO-GO|NO-GO|GO/);
    expect(body).toContain("pnpm lint");
    expect(body).toContain("pnpm typecheck");
    expect(body).toContain("pnpm build");
  });

  it("links the openspec change folder for traceability", () => {
    const body = readFileSync(adrPath, "utf8");
    expect(body).toContain("nextjs-16-react19-migration");
  });
});

describe("SDD change nextjs-16-react19-migration — Next 16 on main", () => {
  it("declares Next.js 16.x in frontend package.json", () => {
    const pkg = readPackageJson();
    const next = pkg.dependencies?.next;
    expect(next).toBeDefined();
    expect(depMajor(next!)).toBeGreaterThanOrEqual(16);
  });

  it("aligns eslint-config-next major with Next", () => {
    const pkg = readPackageJson();
    const eslintNext = pkg.devDependencies?.["eslint-config-next"];
    expect(eslintNext).toBeDefined();
    expect(depMajor(eslintNext!)).toBeGreaterThanOrEqual(16);
  });
});

describe("SDD change nextjs-16-react19-migration — Tailwind CSS v4", () => {
  it("uses @tailwindcss/postcss in PostCSS config", () => {
    const body = readFileSync(postcssConfigPath, "utf8");
    expect(body).toContain("@tailwindcss/postcss");
    expect(body).not.toContain("tailwindcss: {}");
  });

  it("loads design tokens then Tailwind v4 then utilities in globals.css", () => {
    const body = readFileSync(globalsCssPath, "utf8");
    expect(body).toContain('@import "tailwindcss"');
    const iTok = body.indexOf("@import '../styles/tokens.css'");
    const iShad = body.indexOf("@import '../styles/shadcn-theme.css'");
    const iTw = body.indexOf('@import "tailwindcss"');
    const iAnim = body.indexOf('@import "tw-animate-css"');
    const iUtil = body.indexOf("@import '../styles/utilities.css'");
    expect(iTok).toBeGreaterThanOrEqual(0);
    expect(iShad).toBeGreaterThan(iTok);
    expect(iTw).toBeGreaterThan(iShad);
    expect(iAnim).toBeGreaterThan(iTw);
    expect(iUtil).toBeGreaterThan(iAnim);
  });

  it("declares tailwindcss 4.x and tw-animate-css (tailwindcss-animate removed)", () => {
    const pkg = readPackageJson();
    const tw = pkg.devDependencies?.tailwindcss;
    expect(tw).toBeDefined();
    expect(depMajor(tw!)).toBeGreaterThanOrEqual(4);
    expect(pkg.dependencies?.["tailwindcss-animate"]).toBeUndefined();
    expect(pkg.devDependencies?.["tw-animate-css"]).toBeDefined();
  });

  it("pins @tailwindcss/postcss for the v4 CSS pipeline", () => {
    const pkg = readPackageJson();
    expect(pkg.devDependencies?.["@tailwindcss/postcss"]).toBeDefined();
    expect(depMajor(pkg.devDependencies!["@tailwindcss/postcss"]!)).toBeGreaterThanOrEqual(4);
  });
});

describe("SDD change nextjs-16-react19-migration — tailwind.config (dark + HSL parity)", () => {
  it("uses selector darkMode aligned with html[data-theme=dark]", () => {
    const body = readFileSync(tailwindConfigPath, "utf8");
    expect(body).toContain('darkMode: ["selector", \'[data-theme="dark"]\']');
  });

  it("maps shadcn semantic colors to hsl(var(--token)) for Tailwind utilities", () => {
    const body = readFileSync(tailwindConfigPath, "utf8");
    expect(body).toContain("hsl(var(--primary))");
    expect(body).toContain("hsl(var(--background))");
    expect(body).toContain('"./src/**/*.{ts,tsx}"');
  });
});

describe("SDD change nextjs-16-react19-migration — Phase 4 RSC shells", () => {
  it("keeps auth login page as Server Component (LoginForm is the client island)", () => {
    const raw = readFileSync(authLoginPagePath, "utf8").trimStart();
    expect(raw.startsWith("'use client'")).toBe(false);
    expect(raw).toContain("LoginForm");
  });

  it("uses server entry for my-invoices list with optional cookie prefetch helper", () => {
    const raw = readFileSync(myInvoicesListPagePath, "utf8").trimStart();
    expect(raw.startsWith("'use client'")).toBe(false);
    expect(raw).toContain("fetchMyInvoicesInitialAuthenticated");
    expect(raw).toContain("MyInvoicesListClient");
  });

  it("keeps my-invoices layout as Server Component wrapping a client auth gate", () => {
    const raw = readFileSync(myInvoicesLayoutPath, "utf8").trimStart();
    expect(raw.startsWith("'use client'")).toBe(false);
    expect(raw).toContain("MyInvoicesAuthGate");
  });
});

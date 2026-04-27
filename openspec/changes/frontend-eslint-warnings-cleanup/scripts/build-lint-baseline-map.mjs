import fs from 'node:fs';

const jsonPath = new URL('../lint-baseline.json', import.meta.url);
const outPath = new URL('../lint-baseline-map.md', import.meta.url);

const j = JSON.parse(fs.readFileSync(jsonPath, 'utf8'));
const rows = [];
for (const f of j) {
  for (const m of f.messages ?? []) {
    rows.push({
      file: f.filePath.replace(/\\/g, '/'),
      line: m.line,
      rule: m.ruleId ?? '(unknown)',
      sev: m.severity,
    });
  }
}

const warn = rows.filter((r) => r.sev === 1);
const err = rows.filter((r) => r.sev === 2);
const byFile = new Map();
for (const r of warn) {
  if (!byFile.has(r.file)) byFile.set(r.file, new Set());
  byFile.get(r.file).add(`${r.rule} @${r.line}`);
}

let out = `# Lint baseline (Phase 1)\n\n`;
out += `**When:** ${new Date().toISOString()}\n`;
out += `**Command:** \`pnpm exec eslint . --format json\` (cwd: \`frontend/\`)\n\n`;
out += `**Totals:** ${warn.length} warnings, ${err.length} errors\n\n`;
out += `## Map: file → rule (line)\n\n`;

for (const [file, rules] of [...byFile.entries()].sort((a, b) => a[0].localeCompare(b[0]))) {
  const rel = file.includes('/frontend/')
    ? 'frontend/' + file.split('/frontend/').pop()
    : file;
  out += `### \`${rel}\`\n`;
  for (const x of [...rules].sort()) out += `- ${x}\n`;
  out += '\n';
}

fs.writeFileSync(outPath, out);
console.log('Wrote', outPath, 'warnings', warn.length, 'errors', err.length);

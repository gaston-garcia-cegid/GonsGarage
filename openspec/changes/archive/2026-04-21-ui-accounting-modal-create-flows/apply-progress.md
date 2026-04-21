# Apply progress — ui-accounting-modal-create-flows

**Mode**: Strict TDD  
**Batch**: 2026-04-21 — implementação fases 1–6 + **actualização 2026-04-21** (RTL extra: empty CTA, Cancel, ESC, submit + `list` 2×)

## TDD cycle evidence

| Task | Test file(s) | Layer | Safety net | RED | GREEN | TRIANGULATE | REFACTOR |
|------|--------------|-------|------------|-----|-------|---------------|----------|
| 1.1 | — | CSS | ✅ baseline suite | N/A estrutural | ✅ Passed (`accounting.module.css`) | ➖ Single surface | ➖ None |
| 1.2 | — | N/A | — | — | ➖ **Skipped**: shell inline (sem `AccountingCreateDialogShell.tsx`) | — | — |
| 2.1–2.4 | `suppliers/page.test.tsx` | RTL | ✅ baseline | ✅ Written | ✅ Passed | ✅ **6** casos | ➖ None |
| 3.1–3.3 | `received-invoices/page.test.tsx` | RTL | ✅ | ✅ Written | ✅ Passed | ✅ **6** casos | ➖ None |
| 4.1–4.2 | `billing-documents/page.test.tsx` | RTL | ✅ | ✅ Written | ✅ Passed | ✅ **6** casos | ➖ None |
| 5.1–5.2 | `issued-invoices/page.test.tsx` | RTL | ✅ | ✅ Written | ✅ Passed | ✅ **6** casos | ➖ None |
| 6.1 | suite + typecheck + build | — | ✅ | — | ✅ Passed (**73** Vitest; `tsc`; `next build`) | — | — |
| 6.2 | manual + RTL ESC/Cancel | — | — | — | ✅ Passed (Cancel/ESC cobertos em RTL; checklist manual opcional) | — | — |

Casos por ficheiro `page.test.tsx` (cada domínio): toolbar → modal; `?create=1` + `replace`; CTA estado vazio; Cancel sem `create`; ESC sem `create`; submit com `create` + segunda chamada `list`/`listStaff`.

## Test summary

| Métrica | Valor |
|---------|--------|
| Testes RTL accounting (`*accounting/**/page.test.tsx`) | **24** (6 × 4 ficheiros) |
| Testes totais Vitest (`pnpm test`, Abr 2026) | **73** (18 ficheiros de teste) |
| Camada | Integração RTL (`render`, `userEvent`, mocks `@/stores` + serviços) |

## Deviations vs `design.md`

- **Pós-criação**: design preferia só `load()` na lista; **implementado** igual (sem `router.replace` para detalhe após criar no modal).  
- **1.2**: não foi criado `AccountingCreateDialogShell.tsx` — duplicação de header Dialog ainda aceitável em 4 ficheiros; REFACTOR futuro.

## Issues

- Radix avisa `DialogDescription` em falta se `aria-describedby` não for suprimido — usámos `aria-describedby={undefined}` em `DialogContent` nos quatro modais.

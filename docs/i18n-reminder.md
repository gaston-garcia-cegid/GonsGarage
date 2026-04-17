# Lembrete i18n (quando implementarmos traduções)

**Contexto (2026-04):** a UI está em **português de Portugal (pt_PT)** com strings **inline** no código (sem biblioteca i18n).

**Ao adicionar `es_ES` e `en_GB`:**

1. Extrair todas as cadeias visíveis ao utilizador para um único sistema (ex.: **next-intl**, **react-i18next**, ou dicionários por locale).
2. Manter **pt_PT** (ou `pt`) como **locale por defeito** em `layout`, `Intl.*`, e mensagens de validação.
3. Rever `Intl.DateTimeFormat` / `NumberFormat` para usar o locale ativo (ou `pt-PT` / `es-ES` / `en-GB` explicitamente).
4. Alinhar `lang` no `<html>` com o locale escolhido (ex.: `lang="pt-PT"` vs convenção do projeto).
5. Os ficheiros `docs/mvp-minimum-phases.md` e `README` podem referenciar esta nota para onboarding.

*Este ficheiro existe para o equipa não perder o contexto quando abrir a tarefa de i18n.*

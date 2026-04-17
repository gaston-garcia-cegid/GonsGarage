# Relatório de situação do MVP (GonsGarage)

**Data de referência:** 16 de abril de 2026.

## Onde estamos

| Fase (ver `mvp-minimum-phases.md`) | Estado resumido |
|-----------------------------------|-------------------|
| **A — Base técnica** | Concluída: Compose, CI, testes base, documentação. |
| **B — Identidade e dados** | Concluída: JWT, roles, coches/citas com permissões, `/auth/me`, UI cliente alinhada. |
| **C — Reparações** | **Parcial:** API `GET /api/v1/repairs/car/:carId` ativa; histórico por automóvel na página de detalhe; **painel** passa a agregar reparações recentes dos automóveis do utilizador (leitura). Falta CRUD de reparações na API/UI para staff, se for obrigatório para o MVP. |
| **D — Produção** | Por fazer: segredos, deploy real, workflow de deploy operacional. |

## Área do cliente

**Funcionalidades principais cobertas:** registo/login, automóveis (CRUD), marcações (lista + modal), painel com resumo, detalhe do automóvel com histórico de serviços quando existem dados na API, idioma UI **pt_PT** por defeito.

**Gaps ainda relevantes para “MVP cliente fechado”:** notificações por e-mail, recuperação de password, e (opcional) uma única chamada API para “todas as minhas reparações” em vez de N chamadas por automóvel no painel — melhoria de performance, não bloqueante.

## Melhor sugestão para continuar

1. **Se o MVP inclui oficina/staff:** expor `POST/PATCH` de reparações na API Gin (hoje o serviço existe; faltam rotas e UI mínima para técnico/admin) e documentar no Swagger — fecha a **Fase C** de forma completa.
2. **Se o foco é “cliente primeiro”:** avançar para **Fase D** mínima — `JWT_SECRET` obrigatório fora de defaults, `deploy.yml` com ambiente de staging, e checklist de CORS em `release`.
3. **Paralelo recomendado:** item da roadmap “sincronizar Swagger / tipos frontend” para reduzir divergências (`es_ES` / `en_GB` ficam para a tarefa i18n; ver `docs/i18n-reminder.md`).

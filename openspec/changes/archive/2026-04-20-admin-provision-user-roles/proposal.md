# Proposal: Admin (taller) crea utilizadores manager, employee e client

## Intent

Falta um fluxo **JWT + autorização** para um **admin do taller** criar `User` com roles **`manager`**, **`employee`**, **`client`** (hoje: registo público, seeds, e `POST /employees` no domínio colaborador). **“Cliente admin”** = role **`admin`** do MVP. **Client** como administrador multi-tenant → fora de âmbito.

## Scope

### In Scope

- `POST` protegido (ex. `/api/v1/admin/users`) cria `User` com `role` ∈ {manager, employee, client}; **proibir** `role=admin` neste fluxo.
- **`admin` MUST** invocar; **`manager` MAY** só {employee, client} — ou só `admin` se tasks fecharem assim.
- Validação alinhada a `Register`; testes Go (403/201/roles inválidas).
- UI pt_PT mínima **ou** API-only documentado em tasks.

### Out of Scope

- Tenant “client admin”; convites e-mail; SCIM; mudanças grandes em repairs/faturação.

## Capabilities

### New Capabilities

- `staff-user-provisioning`: Quem pode criar utilizador e que `role` pode atribuir; conjunto de roles fechado.

### Modified Capabilities

- `mvp-role-access`: Matriz + testes CI para o novo endpoint.

## Approach

Handler + serviço reutilizando repo/hashing de `Register`; gate de role no serviço; UI com `api-client` se entregar front na mesma change.

## Affected Areas

| Area | Impact |
|------|--------|
| `backend/cmd/api/main.go` | Rota protegida |
| `backend/internal/handler/`, `service/` | Handler + regras |
| `mvp_role_access_test.go` (ou equivalente) | Novos asserts |
| `openspec/specs/mvp-role-access/spec.md` | Delta pós-spec |
| `frontend/src/app/` | UI opcional |

## Risks

| Risk | L | Mitigation |
|------|---|------------|
| Escalação a `admin` | M | Testes negam `role=admin`; review |
| Solapamento com `POST /employees` | M | Design documenta ou unifica |

## Rollback Plan

Revert commits; quitar ruta en `main.go`. Sin migraciones destructivas.

## Dependencies

JWT y modelo `User` actuales.

## Success Criteria

- [ ] `admin` → 201 para roles permitidos; `client`/`employee` → 403 sin permiso.
- [ ] `role=admin` en body → 4xx.
- [ ] CI verde; matriz `mvp-role-access` alineada.

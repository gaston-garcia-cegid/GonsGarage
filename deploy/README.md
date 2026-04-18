# Despliegue LAN / Docker (plantilla Arnela)

## Orden recomendado (desde tu PC con `deploy.ps1`)

1. **`git pull`** en el clon local (`D:\Repos\GonsGarage` o el que uses), para que lo que subís sea la última versión.
2. **`.\deploy.ps1`** en **PowerShell** desde la raíz del repo (no copia el `.git`; solo `backend/`, `frontend/`, `nginx/`, compose y example de env).

Si trabajás **solo en el servidor** con una carpeta que rellenás a mano o con `git clone` allí, entonces el equivalente es **`git pull`** en `/DATA/AppData/gonsgarage` (si usás git en el servidor) y después `docker compose … up -d --build` — no hace falta `deploy.ps1` en ese flujo.

---

Archivos en la **raíz del repo**:

| Archivo | Uso |
|---------|-----|
| [`docker-compose.prod.yml`](../docker-compose.prod.yml) | Redis + API + Next (standalone) + nginx en **8102**. |
| [`.env.prod.example`](../.env.prod.example) | Plantilla; copiar a **`.env.prod`** en el servidor (no git). |
| [`deploy.ps1`](../deploy.ps1) | `scp` + `docker compose` remoto. Por defecto `$REMOTE_DIR` = `/DATA/AppData/gonsgarage` (cambiar en el script si hace falta). |
| [`docker-setup.prod.ps1`](../docker-setup.prod.ps1) | Mismo compose en local para probar antes de subir. |
| [`nginx/default.conf`](../nginx/default.conf) | `/` → front, `/api/` y `/swagger/` → API. |
| [`backend/Dockerfile`](../backend/Dockerfile) | Binario `gonsgarage-api`. |
| [`frontend/Dockerfile`](../frontend/Dockerfile) | Activa `DOCKER_BUILD=1` para `output: "standalone"` (solo en build Linux/Docker; `pnpm build` local en Windows sigue sin standalone). |

## Postgres en el host (Arnela / otro)

1. Crear base `gonsgarage` y usuario con contraseña fuerte en tu Postgres del servidor.
2. En `.env.prod`, `DATABASE_URL` debe apuntar a **`host.docker.internal`** (Linux: ya está en `extra_hosts` del compose) si Postgres escucha en el mismo host que Docker.

## CORS y `GIN_MODE=release`

El API lee **`CORS_ORIGINS`** (coma-separado). Debe coincidir con la URL exacta del navegador (p. ej. `http://192.168.1.100:8102`, sin barra final).

## Secretos (checklist 4.2)

- En el servidor solo **`/DATA/AppData/gonsgarage/.env.prod`** (o la ruta que uses); **no** subir `.env.prod` al git (está en [`.gitignore`](../.gitignore)).
- Obligatorios: **`JWT_SECRET`** (largo, aleatorio), **`DATABASE_URL`** (usuario dedicado, no el superusuario de Postgres).
- **`REDIS_URL`** en el compose prod apunta al contenedor Redis interno; no hace falta secret extra para Redis salvo que cambies el diseño.

## Verificación tras el primer `up` (4.3 / 4.4)

Desde cualquier máquina en la LAN (sustituí host/puerto si cambiaron):

```bash
curl -sS "http://192.168.1.100:8102/health"
curl -sS "http://192.168.1.100:8102/ready"
```

- **4.3:** respuestas JSON OK del API vía nginx.
- **4.4:** abrir `http://192.168.1.100:8102` en el navegador, **login** contra el mismo origen (el front ya fue construido con `NEXT_PUBLIC_API_URL` igual a esa URL pública).

Si `/ready` falla, Postgres no es alcanzable desde el contenedor del API (revisá `DATABASE_URL`, firewall y `host.docker.internal`).

## Login → **502 Bad Gateway** (nginx)

El front carga pero **`POST /api/v1/auth/login`** devuelve 502: nginx **no** recibe respuesta OK del contenedor **`gonsgarage-api`** (caído, reiniciando o no escucha en `8080`).

En el **servidor**:

```bash
cd /DATA/AppData/gonsgarage
docker compose -f docker-compose.prod.yml --env-file .env.prod ps
docker logs gonsgarage-api --tail 120
```

- Si el API está **Exited** o en bucle de reinicio, casi siempre es **`DATABASE_URL`**: Postgres en el host debe escuchar en una interfaz alcanzable desde Docker (no solo `127.0.0.1`). Revisá `postgresql.conf` → `listen_addresses` y `pg_hba.conf` para permitir el bridge Docker (p. ej. `172.17.0.0/16`) o la IP del host. Alternativa: en `.env.prod`, probá el host real en lugar de `host.docker.internal` si Postgres escucha en `192.168.1.100:5432`.
- Comprobar red interna (desde el contenedor nginx al API):

```bash
docker exec gonsgarage-nginx wget -qO- -S http://gonsgarage-api:8080/health 2>&1 | head -20
```

Si esto falla con *connection refused*, el API no está arriba → prioridad: logs de `gonsgarage-api`.

### Logs: `connection refused` a `172.17.0.1:5432` (`host.docker.internal`)

Suelen convivir **dos** cosas:

1. **Contraseña / usuario de ejemplo**  
   Si en el log aparece `CHANGE_ME_strong_password_here`, tu `.env.prod` **no** está listo: editá **`DATABASE_URL`** con la contraseña **real** del rol `gonsgarage_app` (o el usuario que hayas creado en Postgres).

2. **Postgres solo escucha en `127.0.0.1`**  
   Desde el contenedor, `host.docker.internal` apunta al host (a menudo `172.17.0.1`). Si Postgres está configurado con `listen_addresses = 'localhost'` (solo `127.0.0.1`), **no** atiende en esa IP del bridge → *connection refused*.

**En el servidor (host Linux, fuera de Docker):**

```bash
ss -tlnp | grep 5432
```

- Si ves solo `127.0.0.1:5432`, tenés que ampliar escucha y permisos, por ejemplo:
  - En `postgresql.conf`: `listen_addresses = '*'` (o al menos incluir la IP del host en la LAN si preferís algo más acotado).
  - En `pg_hba.conf` una línea para la red Docker, p. ej.:  
    `host  gonsgarage  gonsgarage_app  172.17.0.0/16  scram-sha-256`  
    (ajustá método si usás `md5`. Reiniciá Postgres después: `systemctl restart postgresql` o el servicio que uses.)

**Alternativa rápida:** si tras `listen_addresses = '*'` Postgres queda en `0.0.0.0:5432`, probá en `.env.prod` sustituir el host por la **IP LAN del servidor** (p. ej. `192.168.1.100`) en lugar de `host.docker.internal`, siempre que no haya firewall bloqueando el contenedor hacia esa IP.

Después de cambiar `.env.prod`:

```bash
cd /DATA/AppData/gonsgarage
docker compose -f docker-compose.prod.yml --env-file .env.prod up -d --build
docker logs gonsgarage-api --tail 30
```

## pgAdmin desde tu PC hacia Postgres del servidor

**Recomendado (sin abrir el puerto 5432 a la LAN): túnel SSH**

1. En **PowerShell** o terminal de tu PC (con `ssh` instalado), dejá abierta una sesión que reenvíe el puerto (elegí un puerto local libre, p. ej. **5433**):

```powershell
ssh -L 5433:127.0.0.1:5432 root@192.168.1.100
```

   Eso asume que **Postgres en el servidor escucha en `127.0.0.1:5432`**. Si tu Postgres solo escucha en otra interfaz, ajustá el segundo tramo (consultá con `ss -tlnp | grep 5432` en el servidor).

2. En **pgAdmin** → *Register* → *Server* → pestaña **General**: nombre libre. Pestaña **Connection**:
   - **Host:** `127.0.0.1` o `localhost`
   - **Port:** `5433` (el local del túnel)
   - **Username / Password:** los de tu base `gonsgarage` (los mismos que usaste en `DATABASE_URL` del `.env.prod`, usuario dedicado, no hace falta ser `postgres`).

3. Mantené la ventana del **ssh** abierta mientras usás pgAdmin.

**Alternativa (directo por red):** solo si en el servidor Postgres tiene `listen_addresses` que incluya la LAN y firewall/pg_hba permiten tu IP. Entonces en pgAdmin: **Host** `192.168.1.100`, **Port** `5432`. No lo recomendamos en WiFi compartido sin TLS ni restricción por IP.

## HTTPS vs HTTP en LAN (4.5)

Para **solo red local** (`192.168.x.x`) es aceptable **HTTP** en el checklist: sin certificado, sin Let’s Encrypt. Si exponés el servidor a **internet**, añadí TLS (nginx + certbot, o TLS en el proxy) y actualizá **`CORS_ORIGINS`** / **`NEXT_PUBLIC_API_URL`** a `https://…`.

## Rollback (4.6)

1. **Antes de cada deploy:** copia del `.env.prod` vigente y, si aplica, volcado rápido de BD (`pg_dump -Fc gonsgarage > backup.dump`).
2. **Volver atrás en contenedores:** en el directorio del deploy, `docker compose -f docker-compose.prod.yml --env-file .env.prod down`, restaurá la carpeta/código anterior (o `git checkout` + rebuild) y `… up -d --build`.
3. **Solo API (binario):** sustituí el binario por la versión anterior y reiniciá el proceso; la base suele seguir compatible salvo que hayas dependido de migraciones destructivas (evitar `RESET_DATABASE` en servidor).

## Checklist MVP

Cerrar tareas **4.x** en [`docs/mvp-solo-checklist.md`](../docs/mvp-solo-checklist.md) cuando hayas hecho la **verificación en servidor** (curl + login). Este documento cubre criterios **4.2–4.6** a nivel de runbook; **4.1** son las URLs ya anotadas en el checklist.

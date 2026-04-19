#!/usr/bin/env bash
#
# Arnela — actualiza el código en ESTE clone y reconstruye contenedores de producción.
# (Nombre explícito para servidores con varias apps en /DATA/AppData/…)
#
# IMPORTANTE: en producción se usa SOLO docker-compose.prod.yml. NO mezclar con
# docker-compose.yml: en el repo base, go-api tiene profile "app" y al fusionar
# Compose deja frontend dependiendo de un servicio go-api "inexistente".
#
# Uso típico vía SSH:
#   ssh -i "$HOME/.ssh/tu_clave" user@host 'bash -s' < scripts/update-server-arnela.sh
#
# Requisitos: git, Docker Compose v2, .env.prod en ARNELA_DIR.

set -euo pipefail

# --- configuración (override con export VAR=... antes de ejecutar) ---
: "${ARNELA_DIR:=/DATA/AppData/gonsgarage}"
: "${GIT_REF:=main}"
: "${COMPOSE_FILE:=docker-compose.prod.yml}"
: "${ENV_FILE:=.env.prod}"

cd "$ARNELA_DIR"

echo "==> Directorio: $(pwd)"
echo "==> Rama/ref: $GIT_REF"
echo "==> Compose: $COMPOSE_FILE"

if [[ ! -f "$COMPOSE_FILE" ]]; then
  echo "Error: no se encuentra $COMPOSE_FILE en $GONSGARAGE_DIR" >&2
  exit 1
fi

if [[ ! -f "$ENV_FILE" ]]; then
  echo "Error: falta $ENV_FILE (variables de producción)." >&2
  exit 1
fi

echo "==> git fetch"
git fetch --all --prune

echo "==> git checkout $GIT_REF && pull"
git checkout "$GIT_REF"
git pull --ff-only origin "$GIT_REF"

echo "==> docker compose up -d --build"
docker compose -f "$COMPOSE_FILE" --env-file "$ENV_FILE" up -d --build

echo "==> Estado de contenedores"
docker compose -f "$COMPOSE_FILE" ps

echo "==> Health (backend en el host)"
curl -sS -o /dev/null -w "GET /health -> %{http_code}\n" http://127.0.0.1:8080/health || true
curl -sS -o /dev/null -w "GET /readiness -> %{http_code}\n" http://127.0.0.1:8080/readiness || true

echo "==> Listo (GonsGarage)."

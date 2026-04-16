/**
 * Resuelve la URL base pública del API v1 a partir de NEXT_PUBLIC_API_URL.
 *
 * - Origen: `http://localhost:8080` → `http://localhost:8080/api/v1`
 * - Ya con sufijo: `http://localhost:8080/api/v1` → igual
 * - Corrige sufijos repetidos: `.../api/v1/api/v1` → `.../api/v1`
 */
export function getPublicApiV1BaseUrl(): string {
  let raw = (process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080").trim().replace(/\/+$/, "");
  raw = raw.replace(/(\/api\/v1)+$/i, "/api/v1");
  if (/\/api\/v1$/i.test(raw)) {
    return raw;
  }
  return `${raw}/api/v1`;
}

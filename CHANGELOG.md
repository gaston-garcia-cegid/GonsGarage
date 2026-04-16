# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html) for the **public HTTP contract** (`apiVersion` in `GET /health`). See [docs/api/versioning.md](docs/api/versioning.md).

## [Unreleased]

### Added

- `GET /ready` readiness probe (PostgreSQL ping).
- `apiVersion` field on `GET /health` liveness response.
- GitHub issue templates (bug / feature).
- API versioning policy: `docs/api/versioning.md`.

### Changed

- (none)

## [1.0.0] - 2026-04-16

### Added

- Initial changelog cut aligned with MVP phases A–D: auth, cars, appointments, repairs, Docker deploy smoke, CORS whitelist in `release`, JWT hardening.

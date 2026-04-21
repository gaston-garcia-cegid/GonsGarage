# staff-user-provisioning Specification

> **Canonical:** promoted from `openspec/changes/archive/2026-04-20-admin-provision-user-roles/specs/staff-user-provisioning/spec.md` (2026-04-20).

## Purpose

Define **who MAY create `User` accounts** with which **discrete `role` values** via an authenticated API, without exposing self-service **admin** creation.

In GonsGarage, the HTTP surface is **`POST /api/v1/admin/users`** (JWT required; only `admin` and `manager` pass the staff-manager gate; role matrix below applies in the auth service).

## Normative matrix (summary)

| Caller `role` | MAY assign `manager` | MAY assign `employee` | MAY assign `client` | SHALL NOT assign `admin` via this flow |
|---------------|------------------------|------------------------|----------------------|----------------------------------------|
| `admin`       | MUST                   | MUST                   | MUST                 | MUST                                   |
| `manager`     | MUST NOT               | MUST                   | MUST                 | MUST                                   |
| `employee`    | MUST NOT               | MUST NOT               | MUST NOT             | MUST NOT                               |
| `client`      | MUST NOT               | MUST NOT               | MUST NOT             | MUST NOT                               |

> **SHALL NOT** `admin` in this table means: the provisioning endpoint **MUST NOT** accept a request whose intended role is `admin` (including body tampering).

## Requirements

### Requirement: Authenticated provisioning only

The provisioning endpoint **SHALL** require a valid JWT. Unauthenticated callers **MUST NOT** receive 2xx success.

#### Scenario: Missing token rejected

- GIVEN no `Authorization` header
- WHEN `POST` the provisioning endpoint with an otherwise valid body
- THEN the response **MUST NOT** be 2xx success

### Requirement: Caller role gates assignment

The system **SHALL** enforce the matrix: only permitted caller-role and target-role pairs **MAY** result in user creation.

#### Scenario: Admin creates client

- GIVEN JWT role `admin` and a valid payload with `role` = `client`
- WHEN `POST` the provisioning endpoint
- THEN the response **SHALL** be success (2xx) if all other domain validations pass

#### Scenario: Manager cannot create manager

- GIVEN JWT role `manager` and payload with `role` = `manager`
- WHEN `POST` the provisioning endpoint
- THEN the response **MUST NOT** be 2xx success

#### Scenario: Client cannot provision

- GIVEN JWT role `client`
- WHEN `POST` the provisioning endpoint with any allowed-role payload
- THEN the response **MUST NOT** be 2xx success

### Requirement: No admin escalation via body

The endpoint **MUST NOT** create a user whose stored role is `admin` through this flow.

#### Scenario: Reject admin role in payload

- GIVEN JWT role `admin` and payload with `role` = `admin`
- WHEN `POST` the provisioning endpoint
- THEN the response **MUST NOT** be 2xx success

### Requirement: Valid target roles only

Requested `role` **SHALL** be exactly one of `manager`, `employee`, or `client` (after normalisation). Any other value **MUST NOT** yield 2xx success.

#### Scenario: Unknown role rejected

- GIVEN JWT role `admin` and payload with `role` = `superuser`
- WHEN `POST` the provisioning endpoint
- THEN the response **MUST NOT** be 2xx success

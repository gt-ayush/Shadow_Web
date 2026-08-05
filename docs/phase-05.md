# Phase 5 – Registry Infrastructure Design Document

**Project:** Private Internet Infrastructure

**Phase:** 5

**Component:** Domain Registry

**Status:** Design Specification

---

# 1. Purpose

The Registry is the official operator of one or more Top-Level Domains (TLDs) within the private Internet.

Examples:

* `.x`
* `.web`
* `.shop`
* `.mail`
* `.cloud`

The Registry is the authoritative source for domain registration information under its managed TLDs.

It maintains ownership records, registration metadata, nameserver delegations, and generates the DNS delegation information consumed by the TLD DNS servers.

The Registry **does not** host websites or manage customer DNS records.

---

# 2. Objectives

This phase aims to build a production-style registry service that:

* Operates one or more TLDs.
* Stores all registered domains.
* Validates registration requests.
* Publishes DNS delegations automatically.
* Provides an API for Registrars.
* Supports future scaling across multiple servers.
* Is fully containerized using Podman.

---

# 3. Registry Responsibilities

The Registry is responsible for:

* Managing TLD ownership.
* Maintaining the official domain database.
* Storing domain metadata.
* Storing delegated nameservers.
* Validating registrations.
* Handling renewals.
* Handling transfers.
* Handling deletions.
* Handling expiration.
* Publishing delegation records.
* Maintaining audit logs.
* Providing an internal API for Registrars.

---

# 4. Registry Does NOT

The Registry does **not**:

* Sell domains directly to users.
* Manage websites.
* Store website files.
* Manage customer email.
* Issue SSL certificates.
* Act as an Authoritative DNS server.
* Act as a Recursive Resolver.
* Replace the Registrar.

---

# 5. High-Level Architecture

```text
                           User
                             │
                             ▼
                        Registrar
                             │
                             ▼
                      Registry API
                             │
          ┌──────────────────┼──────────────────┐
          │                  │                  │
          ▼                  ▼                  ▼
     PostgreSQL         Validation        Redis Cache
          │
          ▼
    Zone Generator
          │
          ▼
    DNS Publisher
          │
          ▼
      TLD DNS Servers
```

The Registry is the central authority for all domains under its TLDs.

---

# 6. Podman Architecture

One logical Registry server is represented by one Podman Pod.

```
Registry Pod

├── Registry API
├── PostgreSQL
├── Redis
├── Validation Worker
├── Zone Generator
├── Publisher
├── Metrics Exporter
└── Log Agent
```

Each container has one responsibility.

---

# 7. Container Responsibilities

## Registry API

Purpose:

Provides REST endpoints for Registrars.

Responsibilities:

* Register domain
* Renew domain
* Delete domain
* Transfer domain
* Query domain
* Manage nameservers

---

## PostgreSQL

Purpose:

Permanent storage.

Stores:

* Domains
* Owners
* Registrars
* Transactions
* TLD information

---

## Redis

Purpose:

Temporary data.

Used for:

* Cache
* Queue
* Locking
* Rate limiting

---

## Validation Worker

Checks:

* Domain syntax
* Reserved words
* Existing registrations
* TLD ownership
* Nameserver validation
* Business rules

---

## Zone Generator

Reads PostgreSQL.

Produces:

Delegation records for TLD DNS servers.

Example:

```
example.x

NS ns1.example.x

NS ns2.example.x
```

---

## Publisher

Updates:

Primary TLD DNS

↓

Secondary TLD DNS

↓

Reloads DNS

↓

Domain becomes active

---

## Metrics Exporter

Publishes:

* API latency
* Database metrics
* Registration count
* Publishing status
* Queue size

---

## Log Agent

Collects:

* API logs
* Database logs
* Worker logs
* Publishing logs

---

# 8. Directory Layout

```
internet/

registry/

├── api/
│
├── database/
│
├── migrations/
│
├── workers/
│
├── zone-generator/
│
├── publisher/
│
├── redis/
│
├── monitoring/
│
├── compose/
│
├── ansible/
│
├── documentation/
│
└── backups/
```

---

# 9. Database Design

## Table: TLD

Stores every managed TLD.

Fields:

* ID
* Name
* Status
* Created Date

Example:

```
.x

.web

.shop
```

---

## Table: Domain

Stores every registered domain.

Fields:

* ID
* Domain
* TLD
* Owner
* Status
* Created
* Updated
* Expiration

Example:

```
example.x

ACTIVE
```

---

## Table: Owner

Stores domain owner information.

Fields:

* ID
* Organization
* Contact Email
* Status

---

## Table: Nameserver

Stores delegated nameservers.

Example:

```
ns1.example.x

ns2.example.x
```

---

## Table: Registrar

Stores authorized registrars.

Fields:

* ID
* Name
* API Key
* Status

---

## Table: Transaction

Stores every operation.

Examples:

* Register
* Renew
* Delete
* Transfer

---

# 10. Registration Workflow

```
User

↓

Registrar

↓

Registry API

↓

Validation

↓

Database

↓

Zone Generator

↓

Publisher

↓

TLD DNS

↓

Registration Complete
```

The Registry never communicates directly with end users.

---

# 11. Domain Lifecycle

```
Available

↓

Pending Validation

↓

Registered

↓

Active

↓

Renewed

↓

Expired

↓

Suspended (optional)

↓

Deleted

↓

Available Again
```

---

# 12. Zone Generation Workflow

The Registry stores:

```
example.x

Owner

Nameservers
```

↓

Zone Generator

↓

Produces:

```
example.x

NS ns1.example.x

NS ns2.example.x
```

↓

Publisher

↓

TLD DNS Updated

---

# 13. Registry API

Suggested Endpoints

```
POST   /domains/register

POST   /domains/renew

POST   /domains/transfer

DELETE /domains/delete

GET    /domains/{domain}

GET    /domains

GET    /tlds

GET    /registrars
```

Future:

```
PATCH /domains

POST /domains/restore
```

---

# 14. Validation Rules

Every registration request should verify:

* Domain format
* Valid characters
* Domain length
* Reserved names
* Duplicate registrations
* Existing owner
* Valid nameservers
* Valid registrar credentials
* TLD existence

Only valid registrations proceed.

---

# 15. Publishing Process

Whenever a registration changes:

```
Database Updated

↓

Worker Triggered

↓

Zone Generated

↓

Publisher

↓

Primary TLD DNS

↓

Secondary TLD DNS

↓

DNS Reload

↓

Domain Live
```

No manual intervention is required.

---

# 16. Monitoring

Monitor:

* Total domains
* Registrations/hour
* Renewals/day
* Expirations
* Failed registrations
* Queue size
* API response time
* Database latency
* DNS publishing duration

---

# 17. Logging

Record:

* API requests
* Authentication
* Validation failures
* Registrations
* Renewals
* Transfers
* Zone generation
* Publishing
* System errors

Logs should be centralized for analysis.

---

# 18. Security

Recommendations:

* Private PostgreSQL network.
* API authentication required.
* Encrypt API traffic.
* Rotate API keys.
* Encrypted backups.
* Role-based access.
* Audit every write operation.
* Restrict administrative access.

---

# 19. Scaling Strategy

## Initial Deployment

```
Oracle VPS

Registry Pod
```

Everything runs inside one Podman Pod.

---

## Medium Deployment

```
VPS1

Registry API
```

```
VPS2

PostgreSQL
```

```
VPS3

Publisher
```

```
VPS4

Validation Workers
```

---

## Large Deployment

```
Load Balancer

↓

Registry API Cluster

↓

Redis Cluster

↓

PostgreSQL Cluster

↓

Worker Cluster

↓

Zone Generator Cluster

↓

Publisher Cluster
```

Each component scales independently.

---

# 20. Backup Strategy

Daily:

* PostgreSQL dump
* Redis snapshot (if required)
* Configuration backup
* Registry secrets
* Zone generation history

Weekly:

* Full system backup

Monthly:

* Disaster recovery test

---

# 21. Future Features

* DNSSEC support.
* EPP server for registrar integration.
* WHOIS/RDAP service.
* Billing integration.
* Domain redemption period.
* Registry analytics dashboard.
* Multi-region deployment.
* Multi-TLD management.
* Automated failover.
* High-availability PostgreSQL.
* API rate limiting per registrar.

---

# 22. Success Criteria

Phase 5 is complete when:

* The Registry API is operational.
* PostgreSQL stores all registry data.
* Domain validation is implemented.
* Domain lifecycle management is functional.
* Delegation records are generated automatically.
* TLD DNS updates automatically after approved registrations.
* Registry operations are fully logged.
* Metrics are exposed for monitoring.
* The Registry can scale by moving individual containers or the entire Pod to another VPS with minimal changes.

---

# 23. Deliverables

At the completion of this phase, the project should include:

* Registry Pod definition.
* Registry API service.
* PostgreSQL database schema.
* Redis configuration.
* Validation worker.
* Zone generator.
* DNS publisher.
* Monitoring configuration.
* Logging configuration.
* Backup procedures.
* API documentation.
* Deployment documentation.
* Scaling guide.
* Security guidelines.

This completes the Registry layer of the private Internet architecture. The next phase (Phase 6) will introduce the **Registrar**, which provides the user-facing interface and communicates with the Registry to register, renew, transfer, and manage domains while keeping the Registry isolated from direct public access.

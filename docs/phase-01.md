# Phase 1 – Foundation

## Purpose

Phase 1 establishes the base infrastructure for the entire Private Internet project. It does **not** provide any Internet services such as DNS, web hosting, or domain registration. Instead, it prepares a secure, standardized, and reusable environment on which every later component will be deployed.

The objective is to ensure that future services can be installed, updated, moved, or replaced without requiring changes to the underlying operating system.

---

# Goals

The primary goals of Phase 1 are to:

* Prepare the operating system.
* Secure the server.
* Install the container platform.
* Create a standard project structure.
* Configure automation tools.
* Establish networking for future services.
* Prepare the environment for scaling.

By the end of this phase, the server should be capable of hosting any component of the Private Internet.

---

# Infrastructure Philosophy

The project separates **infrastructure** from **services**.

The Virtual Private Server (VPS) is considered only a resource provider.

It provides:

* CPU
* Memory
* Storage
* Network connectivity

The VPS itself should contain as little project-specific configuration as possible.

Instead, every Internet component will later be deployed as containers grouped into Podman Pods.

This separation provides portability, repeatability, and easier maintenance.

---

# Operating System Preparation

A minimal and stable Linux distribution is installed.

During this step:

* The operating system is updated.
* Security patches are applied.
* Required utilities are installed.
* System time is synchronized.
* Hostname is configured.

The operating system becomes the stable foundation for every future phase.

---

# Security Configuration

Before any services are deployed, the server is secured.

Security tasks include:

* Creating a dedicated administrator account.
* Disabling direct root login.
* Using SSH key authentication.
* Enabling the firewall.
* Removing unnecessary services.
* Applying basic system hardening.

The goal is to minimize the attack surface before exposing any public services.

---

# Container Platform

Instead of installing applications directly onto the operating system, the project uses **Podman**.

Podman allows each Internet component to run inside isolated containers.

Examples of future containers include:

* Root DNS
* Recursive Resolver
* Registry
* Registrar
* Certificate Authority
* Database
* Reverse Proxy
* Monitoring
* Mail Server

Running applications inside containers provides:

* Isolation
* Portability
* Easier upgrades
* Simpler backups
* Reproducible deployments

---

# Pod-Based Architecture

The project is designed around the concept of logical servers.

Each logical server is represented by a Podman Pod.

Examples:

* Root DNS Server
* TLD Server
* Recursive Resolver
* Registry
* Registrar
* Hosting Node

Each Pod may contain multiple containers that work together.

Example:

```text
Pod: Registry

├── Registry API
├── PostgreSQL
├── Redis
└── Metrics Exporter
```

The Pod behaves as a single logical server even though multiple containers are running inside it.

---

# Networking Preparation

During Phase 1, dedicated container networks are created.

Examples include:

* dns-net
* registry-net
* hosting-net
* monitoring-net
* backend-net

These networks isolate services from one another while allowing controlled communication between components.

This design makes future deployments cleaner and easier to secure.

---

# Automation

Automation is introduced from the beginning.

Git is used for:

* Version control
* Infrastructure history
* Configuration management

Ansible is used for:

* Server provisioning
* Package installation
* Configuration deployment
* Future infrastructure automation

Every configuration change should eventually become automated instead of being performed manually.

---

# Project Directory Structure

A standardized directory hierarchy is created.

Each major component receives its own directory.

Example:

* DNS
* Registry
* Registrar
* Certificate Authority
* Hosting
* Monitoring
* Scripts
* Documentation
* Backups

Keeping the project organized from the beginning makes future expansion significantly easier.

---

# Monitoring Preparation

Basic monitoring utilities are installed to observe:

* CPU usage
* Memory usage
* Disk usage
* Network activity

These tools help verify that the infrastructure remains healthy while additional services are added in later phases.

More advanced monitoring platforms such as Prometheus and Grafana will be introduced in a later phase.

---

# Backup Preparation

Directories and procedures for backups are established.

Although no production data exists yet, backup practices are introduced early so they become part of the normal operational workflow.

Future phases will expand this to include:

* Database backups
* Certificate backups
* DNS configuration backups
* Application backups

---

# Documentation

Documentation begins in Phase 1.

The documentation should include:

* Server specifications
* Operating system version
* Installed software
* Network layout
* Directory structure
* Security configuration
* Deployment procedures

Maintaining accurate documentation throughout the project reduces future maintenance effort and simplifies troubleshooting.

---

# Expected Outcome

After completing Phase 1, the environment should provide:

* A secure Linux server.
* A fully functional Podman installation.
* Standardized networking.
* Automation tools.
* Version control.
* Organized project directories.
* Basic monitoring.
* Backup preparation.

No Internet-facing services are deployed during this phase.

Instead, the infrastructure is fully prepared to host them.

---

# Role in the Overall Project

Phase 1 serves as the foundation upon which every remaining phase depends.

Future phases will build on this environment without requiring significant changes to the underlying server.

The progression is as follows:

**Phase 1:** Foundation → Prepare the infrastructure.

**Phase 2:** Root DNS Infrastructure → Create the private DNS root.

**Phase 3:** Recursive Resolver → Enable client DNS resolution.

**Phase 4:** TLD Infrastructure → Introduce custom Top-Level Domains.

**Phase 5 and Beyond:** Build the remaining Internet services, including registries, registrars, authoritative DNS servers, certificate authorities, hosting platforms, monitoring, and user-facing services.

A properly completed Phase 1 ensures that all future components can be deployed consistently, migrated between VPSs when needed, and managed using the same standardized infrastructure.

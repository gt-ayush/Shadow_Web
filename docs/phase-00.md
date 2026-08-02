# Phase 0 – Foundation

## Purpose

Phase 0 establishes the foundation of the entire private Internet project. Before building DNS infrastructure, registries, hosting platforms, or web services, a stable and consistent operating environment is required.

The objective of this phase is **not** to create Internet services. Instead, it prepares the infrastructure that every later phase will depend on. A well-prepared foundation reduces future maintenance, simplifies deployment, and makes scaling much easier.

Think of this phase as constructing the land, roads, power, and utilities before building a city.

---

# Objectives

The primary goals of Phase 0 are to:

* Prepare the virtual infrastructure.
* Install and configure the operating system.
* Secure remote administration.
* Install the container runtime.
* Install automation tools.
* Create a standardized project structure.
* Configure basic system security.
* Establish documentation and version control.

When Phase 0 is complete, the server should be ready to host any Internet component without requiring additional operating system configuration.

---

# Infrastructure Philosophy

The project follows a layered architecture.

```text
Physical Hardware
        │
        ▼
Oracle Cloud Virtual Machine
        │
        ▼
Linux Operating System
        │
        ▼
Podman
        │
        ▼
Pods
        │
        ▼
Containers
        │
        ▼
Internet Services
```

Each layer has a single responsibility.

* The Oracle VM provides compute resources.
* Linux provides the operating system.
* Podman manages containerized workloads.
* Pods represent logical servers.
* Containers run the individual applications.

Keeping these responsibilities separate makes the platform portable and easier to maintain.

---

# Why Oracle Cloud?

Oracle Cloud Free Tier provides sufficient resources to host many lightweight services on a single virtual machine. Instead of creating dozens of small VPS instances, the project uses one or more larger VMs that run multiple isolated Podman pods.

This approach provides:

* Better resource utilization.
* Lower operational complexity.
* Easier backup and monitoring.
* Simple migration of workloads to additional VMs later.

The virtual machine is viewed only as infrastructure. Internet services are never tightly coupled to a specific VM.

---

# Operating System Preparation

The operating system forms the base of every service.

During this phase the server is:

* Updated to the latest packages.
* Configured with correct time synchronization.
* Assigned a permanent hostname.
* Prepared with required system utilities.
* Cleaned of unnecessary packages.

Keeping the operating system minimal reduces security risks and simplifies future maintenance.

---

# System Hardening

Because every future service depends on the operating system, security begins here.

The server is hardened by:

* Disabling direct root login.
* Using SSH key authentication.
* Limiting remote access.
* Applying security updates.
* Enabling a firewall.
* Removing unnecessary services.

These measures create a secure baseline before any public-facing applications are deployed.

---

# Memory Management

Even though the Oracle VM provides a generous amount of RAM, emergency memory management is configured.

A small swap area is created to prevent applications from terminating unexpectedly if memory becomes exhausted.

The Linux swappiness value is adjusted to encourage the system to use physical memory first while reserving swap for exceptional situations.

This configuration improves system stability without relying on swap for normal operation.

---

# Container Platform

Podman is selected as the container engine for the project.

Unlike traditional virtualization, containers share the host operating system while remaining isolated from one another.

Each Internet component will eventually run inside one or more containers.

Examples include:

* Root DNS servers
* Recursive resolvers
* Registry
* Registrar
* Certificate Authority
* Web hosting
* Mail services
* Monitoring

Podman provides a lightweight and secure method for deploying these services.

---

# Why Pods?

A Pod represents a logical server.

Rather than placing every service into one large container, related containers are grouped into a Pod.

For example, a Registry Pod may contain:

* Registry API
* PostgreSQL
* Redis
* Metrics exporter

These containers work together as if they were running on a dedicated server.

This design makes each Internet component modular and portable.

---

# Automation

Manual configuration does not scale well.

Therefore, automation tools are installed during Phase 0.

Ansible will later automate:

* Software installation.
* Configuration management.
* Pod deployment.
* Updates.
* Maintenance.

Using Infrastructure as Code ensures the environment can be recreated consistently.

---

# Version Control

Every configuration file is stored in Git.

This includes:

* Pod definitions.
* DNS configuration.
* Automation scripts.
* Documentation.
* Certificates (excluding private keys).
* Deployment files.

Version control provides:

* Change history.
* Rollback capability.
* Collaboration.
* Configuration tracking.

No manual configuration should exist outside the repository unless it contains sensitive secrets.

---

# Project Directory Structure

A standardized directory layout is created before any services are deployed.

Organizing the project from the beginning makes it easier to locate configurations, automate deployments, and maintain documentation as the infrastructure grows.

Every future phase will build upon this directory structure.

---

# Networking Preparation

Although networking is implemented in the next phase, the project prepares for it by defining a consistent naming strategy.

Dedicated networks will later separate:

* DNS infrastructure.
* Registry services.
* Hosting platform.
* Monitoring.
* Storage.
* Mail services.

Separating traffic improves security and simplifies troubleshooting.

---

# Documentation

Documentation begins with Phase 0.

Every architectural decision, configuration change, deployment procedure, and troubleshooting step should be recorded.

Good documentation ensures the project remains understandable as it becomes more complex and allows the entire environment to be recreated in the future.

---

# Expected Outcome

When Phase 0 is complete, the server should not yet function as a private Internet.

Instead, it provides a stable platform that is ready to host all future components.

The operating system is secured, automation tools are installed, the container runtime is operational, version control is configured, and the project structure is in place.

From this point onward, every new Internet component can be deployed as a Podman pod without requiring changes to the underlying operating system.

This foundation allows the project to grow gradually—from a single Oracle Cloud VM running several pods to a distributed environment spanning multiple VMs—while maintaining the same architecture and deployment process.

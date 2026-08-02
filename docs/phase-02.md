# Phase 2 – Root DNS Infrastructure

## Purpose

The Root DNS Infrastructure is the foundation of the entire private Internet. Every DNS lookup begins here.

Just as the public Internet has a single logical root zone (`.`) that delegates requests to Top-Level Domain (TLD) servers such as `.com`, `.org`, and `.net`, this project creates its own private root that delegates requests to custom TLDs such as `.x`, `.web`, `.shop`, `.mail`, and `.cloud`.

The Root DNS servers do **not** know the IP address of every website. Their only responsibility is to direct clients to the correct TLD server.

---

# Role in the DNS Hierarchy

The Root DNS is the first layer of the DNS resolution process.

```text
User

↓

Recursive Resolver

↓

Root DNS

↓

TLD Server

↓

Authoritative DNS

↓

Website
```

When a client requests:

```text
blog.example.x
```

the Root DNS does **not** answer with the website's IP address.

Instead, it replies with:

> "I don't host `blog.example.x`, but the `.x` TLD is managed by these name servers."

The client (through its recursive resolver) then contacts the `.x` TLD server to continue the lookup.

This delegation process mirrors how the public Internet functions.

---

# Responsibilities

The Root DNS is responsible for:

* Hosting the private root zone (`.`).
* Maintaining the list of available custom TLDs.
* Delegating each TLD to its authoritative TLD servers.
* Providing referrals to recursive resolvers.
* Remaining highly available and consistent across all root instances.

The Root DNS is **not** responsible for:

* Hosting website records.
* Managing customer domains.
* Registering domains.
* Issuing certificates.
* Serving web content.

Those responsibilities belong to later phases.

---

# Root Zone

The Root Zone is the highest level of the DNS hierarchy.

Example structure:

```text
.

├── x
├── web
├── shop
├── cloud
├── mail
└── dev
```

Each entry represents a custom Top-Level Domain.

For each TLD, the Root Zone stores a delegation to the TLD's name servers.

Example concept:

```text
Root Zone

↓

.x

↓

ns1.x
ns2.x
```

The Root DNS simply tells resolvers where the `.x` TLD is hosted.

---

# Root Server Architecture

The project uses multiple Root DNS servers to simulate the public Internet.

Example:

```text
Root-A

Root-B

Root-C

Root-D
```

Each Root Server is identical.

They all contain:

* The same root zone.
* The same delegations.
* The same configuration.

Initially, these servers may all run on the same Oracle Cloud VM for simplicity.

As the infrastructure grows, each Root Server Pod can be moved to a separate VPS without changing the overall architecture.

This provides:

* Better scalability.
* High availability.
* Easier maintenance.
* Fault isolation.

---

# Pod Design

Each Root Server is implemented as a dedicated Podman Pod.

Example:

```text
Pod: Root-A

├── CoreDNS
├── Metrics Exporter
└── Log Collector
```

### CoreDNS

Hosts the Root Zone and responds to DNS queries.

### Metrics Exporter

Collects operational metrics such as:

* Queries per second.
* Response latency.
* Request counts.
* Resource usage.

These metrics are consumed by the monitoring system.

### Log Collector

Collects DNS logs and forwards them to the centralized logging platform.

This separation keeps each component focused on a single responsibility while allowing them to share the same network namespace.

---

# Why Use Pods?

A Pod represents one logical server.

Although multiple containers exist inside the Pod, they behave as one machine.

Advantages include:

* Shared IP address.
* Shared ports.
* Easier communication between containers.
* Simplified deployment.
* Better modularity.

Each Root Server Pod can later be migrated independently to another VPS if additional capacity or redundancy is required.

---

# Networking

All Root Servers are connected to the dedicated DNS network.

Example:

```text
dns-net

├── Root-A
├── Root-B
├── Root-C
└── Root-D
```

Only DNS services are exposed externally.

Monitoring and management traffic remain on internal networks wherever possible.

---

# Synchronization

Every Root Server must always contain the same Root Zone.

Several synchronization methods can be used:

* Git-based deployments.
* Configuration management tools (such as Ansible).
* DNS zone transfers between primary and secondary servers.

Regardless of the mechanism, consistency is critical because all Root Servers must return identical delegation information.

---

# High Availability

The initial deployment may place all Root Server Pods on one VPS.

Example:

```text
Oracle VM

├── Root-A
├── Root-B
├── Root-C
└── Root-D
```

This configuration is acceptable for development and testing.

A production-like deployment distributes Root Servers across multiple VPSs.

Example:

```text
Oracle VM 1

Root-A
```

```text
Oracle VM 2

Root-B
```

```text
Oracle VM 3

Root-C
```

```text
Oracle VM 4

Root-D
```

If one VPS becomes unavailable, the remaining Root Servers continue serving the Root Zone.

---

# Scalability

One of the design goals is portability.

Each Root Server Pod should be completely self-contained.

It includes:

* Configuration.
* DNS data.
* Volumes.
* Networking.
* Monitoring.
* Logging.

If a VPS becomes overloaded, the entire Pod can be moved to another VPS without redesigning the DNS architecture.

Only the deployment location changes.

This makes the infrastructure flexible and easy to expand over time.

---

# Monitoring

Each Root Server exports operational metrics.

Typical metrics include:

* DNS query rate.
* Successful responses.
* Referral responses.
* Error responses.
* CPU usage.
* Memory usage.
* Network throughput.
* Container health.

Logs are forwarded to the centralized monitoring platform for troubleshooting and long-term analysis.

---

# Security

The Root DNS is one of the most critical components in the infrastructure.

Recommended security practices include:

* Restrict administrative access to trusted hosts.
* Allow only DNS traffic on TCP and UDP port 53.
* Use SSH keys instead of password authentication.
* Keep the Root Zone read-only during normal operation.
* Maintain regular backups of the Root Zone.
* Monitor for unexpected configuration changes.

---

# Expected Outcome

At the completion of Phase 2, the project will have a fully functional private Root DNS infrastructure.

The Root DNS will:

* Host the private Root Zone (`.`).
* Delegate custom TLDs.
* Respond to recursive resolver requests.
* Operate through multiple synchronized Root Server Pods.
* Be portable between VPSs.
* Be monitored and logged.
* Form the foundation for all later DNS components.

Although websites and customer domains do not yet exist, the DNS hierarchy is now established.

Subsequent phases will build upon this foundation by introducing Recursive Resolvers, TLD Servers, Registries, Registrars, and Authoritative DNS servers, completing the private Internet architecture.

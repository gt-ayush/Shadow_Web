# Phase 4 – Top-Level Domain (TLD) Infrastructure

## Purpose

The purpose of Phase 4 is to build the **Top-Level Domain (TLD) layer** of the private Internet.

A TLD server acts as the bridge between the Root DNS servers and the domain owner's authoritative DNS servers.

In the public Internet, TLD servers are responsible for domains such as:

* `.com`
* `.org`
* `.net`
* `.io`

In this project, custom TLDs will be created instead, such as:

* `.x`
* `.web`
* `.shop`
* `.cloud`
* `.mail`
* `.dev`

These TLDs exist only within the private Internet and are controlled entirely by the platform.

---

# Position in the DNS Hierarchy

The TLD layer sits between the Root DNS servers and the Authoritative DNS servers.

```text
Client

↓

Recursive Resolver

↓

Root DNS (.)

↓

TLD Server (.x)

↓

Authoritative DNS (example.x)

↓

Website
```

The Root DNS knows where every TLD is located.

The TLD server knows which authoritative DNS servers manage each registered domain.

The authoritative DNS server knows the actual DNS records for that domain.

---

# Responsibilities of the TLD Server

The TLD server has only one responsibility:

**Delegate registered domains to their authoritative name servers.**

For example, if a user owns:

```text
example.x
```

The TLD server stores only the delegation information.

Example:

```text
example.x

↓

ns1.example.x

ns2.example.x
```

When a resolver asks:

> "Who manages example.x?"

The TLD replies:

> "Ask ns1.example.x or ns2.example.x."

It does **not** answer questions about website IP addresses, mail servers, or TXT records.

---

# What the TLD Server Stores

The TLD server maintains a list of registered domains and the authoritative name servers responsible for them.

Typical information includes:

* Domain name
* Authoritative name servers (NS records)
* Glue records (when required)
* Domain status
* Expiration information (provided later by the Registry)

The TLD server does **not** maintain website configuration.

---

# What the TLD Server Does NOT Store

The following records belong to the domain owner's authoritative DNS server and are **not** stored by the TLD server:

* A records
* AAAA records
* MX records
* TXT records
* CNAME records
* SRV records
* HTTPS records

Example:

The TLD server never stores:

```text
blog.example.x

↓

10.20.5.15
```

Instead, it simply replies:

```text
Ask ns1.example.x
```

---

# Example DNS Resolution

Suppose a client wants to visit:

```text
blog.example.x
```

The DNS lookup occurs as follows:

1. The client sends the request to the Recursive Resolver.
2. The Recursive Resolver asks the Root DNS.
3. The Root DNS replies:

```text
The .x TLD server is responsible.
```

4. The Resolver contacts the `.x` TLD server.
5. The `.x` TLD server replies:

```text
example.x

↓

ns1.example.x
ns2.example.x
```

6. The Resolver contacts the authoritative DNS server.
7. The authoritative DNS server returns:

```text
blog.example.x

↓

10.20.5.15
```

Only the authoritative server knows the final IP address.

---

# Glue Records

Sometimes the authoritative DNS server is inside the same domain.

Example:

```text
example.x

↓

ns1.example.x
```

This creates a dependency problem.

To contact `ns1.example.x`, the resolver first needs its IP address.

However, the IP address is inside the same domain that has not yet been resolved.

To solve this, the TLD provides a **Glue Record**.

Example:

```text
example.x

NS

ns1.example.x

Glue

ns1.example.x

A

10.0.50.20
```

The resolver can now contact the authoritative server without entering a circular dependency.

Glue records are only provided when necessary.

---

# TLD Server Deployment

Each TLD is deployed independently.

Example:

```text
Pod

↓

TLD-.x
```

Inside the Pod:

* DNS Server (CoreDNS or BIND)
* Metrics Exporter
* Log Collector

Every TLD Pod behaves as one logical DNS server.

This design allows each TLD to be moved to another VPS without changing the overall architecture.

---

# Supported TLDs

Initially, the following TLDs will be created:

```text
.x

.web

.shop

.cloud

.mail

.dev
```

Each TLD operates independently and has its own DNS configuration.

Future TLDs can be added without modifying the Root DNS architecture.

---

# Registration Workflow (Current Phase)

During Phase 4, domain registration is manual.

The administrator performs the following steps:

1. Choose a domain name.
2. Add the domain delegation to the TLD configuration.
3. Define the authoritative name servers.
4. Reload the DNS service.
5. Verify that the delegation is working.

This manual process will later be automated.

---

# Future Automation

In Phase 5 (Registry), a central database becomes the source of truth for registered domains.

In Phase 6 (Registrar), users register domains through a web interface or API.

The automated workflow becomes:

```text
User

↓

Registrar

↓

Registry Database

↓

TLD Zone Generator

↓

DNS Reload

↓

Domain Active
```

The TLD servers no longer require manual editing because they consume data generated by the Registry.

---

# High Availability

Every TLD should eventually have:

* One Primary Server
* One or more Secondary Servers

Example:

```text
Primary

↓

Secondary A

↓

Secondary B
```

Secondary servers synchronize their zones from the primary using DNS zone transfers (AXFR/IXFR) or another replication mechanism.

This improves availability and resilience.

---

# Monitoring

Every TLD server should expose operational metrics, including:

* DNS query count
* Response time
* Delegation count
* Zone reload status
* CPU usage
* Memory usage
* Network traffic

Centralized logging should capture DNS queries, errors, and operational events for troubleshooting and auditing.

---

# Security

The TLD infrastructure should follow these principles:

* Restrict management access to administrators.
* Allow only DNS traffic on TCP/UDP port 53.
* Keep zone data under version control.
* Regularly back up zone files and configuration.
* Use immutable container images for consistent deployments.
* Run containers with the minimum required privileges.

---

# Success Criteria

Phase 4 is complete when:

* The Root DNS successfully delegates requests to each custom TLD.
* Each TLD server responds with the correct NS delegations.
* Glue records are returned when required.
* Registered domains can be delegated to authoritative DNS servers.
* Multiple TLDs operate independently.
* TLD services can be relocated between VPSs without changing the logical architecture.

At this point, the DNS hierarchy consists of:

```text
Client

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

The next phase introduces the **Registry**, which becomes the authoritative source of all domain registration data and automates the generation of TLD delegations.

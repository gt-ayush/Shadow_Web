# Phase 3 – Recursive DNS Resolver

## Purpose

The **Recursive DNS Resolver** is the first DNS server that users interact with. Unlike a Root DNS server or a TLD server, it does **not** own any DNS records. Its responsibility is to find the answer on behalf of the client by walking through the DNS hierarchy.

In this private Internet, every client device (browser, operating system, or application) will be configured to use the recursive resolver as its primary DNS server. Clients will never contact the Root DNS or TLD servers directly.

---

# Role in the DNS Hierarchy

The recursive resolver sits between clients and the rest of the DNS infrastructure.

```text
                    Client
                       │
                       ▼
             Recursive Resolver
                       │
         ┌─────────────┴─────────────┐
         ▼                           ▼
      Root DNS                  Cached Results
         │
         ▼
      TLD Server
         │
         ▼
 Authoritative DNS
         │
         ▼
     Final DNS Record
```

The resolver performs all communication with the DNS hierarchy while presenting a simple interface to clients.

---

# Why Is a Recursive Resolver Needed?

Without a recursive resolver, every application would need to:

1. Know where the Root DNS servers are.
2. Query the Root DNS.
3. Query the correct TLD server.
4. Query the authoritative DNS server.
5. Cache the results.

This would be inefficient and would duplicate DNS logic across every application.

Instead, the recursive resolver performs these tasks once and shares the cached results with all clients.

---

# Responsibilities

The recursive resolver is responsible for:

* Accepting DNS requests from clients.
* Searching its cache for existing answers.
* Querying the Root DNS servers when necessary.
* Following referrals to TLD servers.
* Following referrals to authoritative DNS servers.
* Returning the final DNS records to the client.
* Caching results to improve performance.
* Logging DNS activity.
* Collecting operational metrics.

It is **not** responsible for hosting DNS records or domains.

---

# How DNS Resolution Works

Consider the example:

```
blog.example.x
```

The resolver processes the request as follows.

### Step 1 – Client Request

A client sends a request to the recursive resolver asking for the IP address of:

```
blog.example.x
```

The client expects a single response and does not know how DNS is structured internally.

---

### Step 2 – Cache Lookup

The resolver first checks its cache.

If a valid cached record already exists:

```
Cache

blog.example.x

↓

192.168.100.20
```

The answer is returned immediately.

No additional DNS servers are contacted.

---

### Step 3 – Query the Root DNS

If the record is not cached, the resolver contacts one of the configured Root DNS servers.

The Root DNS does **not** know the address of `blog.example.x`.

Instead, it responds with a referral.

Example:

```
"I do not know blog.example.x.

The .x TLD is served by:

ns1.x

ns2.x"
```

---

### Step 4 – Query the TLD Server

The resolver then contacts the `.x` TLD server.

The TLD server also does not know the website's IP address.

Instead, it knows which authoritative DNS servers manage the domain.

Example:

```
example.x

↓

ns1.example.x

ns2.example.x
```

The resolver receives another referral.

---

### Step 5 – Query the Authoritative DNS

The resolver contacts the authoritative DNS server for `example.x`.

This server owns the DNS records for the domain.

Example:

```
blog.example.x

↓

192.168.100.20
```

The authoritative server returns the final answer.

---

### Step 6 – Return the Result

The resolver:

* stores the result in its cache,
* returns the IP address to the client.

The client can now connect to the web server.

---

# Cache

Caching is one of the most important responsibilities of the recursive resolver.

Without caching:

```
100 Users

↓

100 Root Queries

↓

100 TLD Queries

↓

100 Authoritative Queries
```

Every lookup would repeat the entire resolution process.

With caching:

```
First User

↓

Full DNS Lookup

↓

Cache
```

```
Next 99 Users

↓

Resolver Cache

↓

Immediate Response
```

This significantly reduces latency and lowers the load on the DNS infrastructure.

---

# Root Hints

Unlike public DNS resolvers, this private resolver does not use the Internet's root servers.

Instead, it is configured with a custom list of private root servers.

Example concept:

```
Root-A

10.10.0.10
```

```
Root-B

10.10.0.11
```

These addresses form the starting point for every lookup.

---

# Networking

The resolver communicates with three groups of systems.

### Clients

Clients send DNS requests to the resolver.

### Root DNS

The resolver queries the Root DNS to locate TLD servers.

### TLD and Authoritative DNS

The resolver continues following referrals until it receives the final DNS records.

---

# Security

The recursive resolver should only accept DNS requests from trusted clients.

Public access should be disabled unless the resolver is intentionally being operated as a public recursive resolver.

Security measures include:

* restricting access to internal networks,
* limiting request rates,
* logging all queries,
* monitoring unusual traffic,
* running with minimal privileges.

---

# High Availability

Initially, a single resolver is sufficient.

As the environment grows, multiple resolver instances can be deployed.

Example:

```
Clients

↓

Load Balancer

↓

Resolver-01

Resolver-02

Resolver-03
```

If one resolver fails, another can continue processing requests without interrupting client access.

---

# Pod Design

Each resolver is deployed as a dedicated Podman pod.

Example:

```
Resolver Pod

├── Unbound
├── Metrics Exporter
└── Log Collector
```

The main DNS service performs recursive resolution, while supporting containers provide monitoring and centralized logging.

This separation allows the resolver to remain lightweight while integrating with the platform's monitoring and observability systems.

---

# Phase Outcome

At the completion of Phase 3, the environment includes:

* one or more recursive DNS resolver pods,
* a functional cache,
* communication with the private Root DNS servers,
* referral handling for custom TLDs,
* centralized logging,
* operational metrics,
* secure client access.

Although authoritative DNS servers are not yet implemented, the resolver is fully prepared to integrate with them in later phases. Once the TLD and authoritative DNS infrastructure is added, clients will be able to resolve domains across the complete private DNS hierarchy exactly as they would on the public Internet.


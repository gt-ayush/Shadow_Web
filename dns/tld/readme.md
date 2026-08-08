# Shadow_Web - Private TLD DNS Infrastructure

This component establishes the authoritative Top-Level Domain (TLD) zone layer for the **Shadow_Web** private hierarchical DNS infrastructure. Built using **CoreDNS**, it manages multiple custom TLDs (`.x`, `.web`, `.shop`, `.cloud`, `.mail`, `.dev`) and provides downstream delegation to Second-Level Domains (SLDs).

---

## 🌐 Network & Addressing

* **Network**: `dns-net` (Bridge)


* **Subnet**: `10.89.10.0/24`
* **TLD Containers & Static IPs**:
* `tld-x-01`: `10.89.10.20`

* `tld-web-01`: `10.89.10.30`

* `tld-shop-01`: `10.89.10.40`

* `tld-cloud-01`: `10.89.10.50`

* `tld-mail-01`: `10.89.10.60`

* `tld-dev-01`: `10.89.10.70`




---

## 📁 Directory Structure

```bash
~/Shadow_Web/dns/tld$ tree
.
├── readme.md
├── tld-cloud
│   ├── Corefile
│   └── zone
│       └── db.cloud
├── tld-dev
│   ├── Corefile
│   └── zone
│       └── db.dev
├── tld-mail
│   ├── Corefile
│   └── zone
│       └── db.mail
├── tld-shop
│   ├── Corefile
│   └── zone
│       └── db.shop
├── tld-web
│   ├── Corefile
│   └── zone
│       └── db.web
└── tld-x
    ├── Corefile
    └── zone
        └── db.x
```

---

## ⚙️ Configuration Patterns

### 1. Corefile Example (e.g., `.web`)

Instructs CoreDNS to serve the specific TLD zone file with logging and errors enabled.

```corefile
web. {
    file /etc/coredns/zone/db.web
    log
    errors
}

```

### 2. Zone File Example (e.g., `db.web`)

Defines the Start of Authority (SOA), nameservers, glue records, and sample SLD delegations.

```zone
$TTL 86400
@   IN  SOA ns1.dns.web. admin.dns.web. (
            2026080301 ; serial
            3600       ; refresh
            1800       ; retry
            604800     ; expire
            86400 )    ; minimum

@   IN  NS  ns1.dns.web.
@   IN  NS  ns2.dns.web.

; Glue records for .web TLD NS itself
ns1.dns.web.   IN  A   10.89.10.30
ns2.dns.web.   IN  A   10.89.10.31

; --- Subdomain / SLD Delegations ---
shop.web.   IN  NS  ns1.shop.web.
ns1.shop.web.  IN  A   10.89.30.10

```

---

## 🚀 Deployment

### Automated Deployment via Ansible

The entire TLD tier can be provisioned using the provided Ansible playbook:

```bash
ansible-playbook -i inventory.ini playbooks/deploy_dns_tlds.yml

```

### Manual Deployment via Podman (Example for `.web`)

```bash
podman run --name tld-web-01 \
  --network dns-net \
  --ip 10.89.10.30 \
  -v $(pwd)/tld-web/Corefile:/etc/coredns/Corefile:ro \
  -v $(pwd)/tld-web/zone:/etc/coredns/zone:ro \
  -d docker.io/coredns/coredns:latest -conf /etc/coredns/Corefile

```

---

## 🔍 Verification

Test each TLD server directly using an ephemeral container on the `dns-net` network:

```bash
# Test .x TLD
podman run --rm --network dns-net docker.io/library/alpine sh -c \
  "apk add --no-cache bind-tools >/dev/null && dig @10.89.10.20 example.x NS"

# Test .web TLD
podman run --rm --network dns-net docker.io/library/alpine sh -c \
  "apk add --no-cache bind-tools >/dev/null && dig @10.89.10.30 shop.web NS"

```

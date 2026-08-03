# Shadow_Web: Phase 5 - Second-Level Domain (SLD) Authoritative Servers

This directory contains the configurations, Corefiles, zone files, and deployment playbooks for the Second-Level Domain (SLD) authoritative name servers within the private `Shadow_Web` DNS infrastructure.

## Architecture & IP Allocations

The SLD tier resolves specific subdomains under custom Top-Level Domains (TLDs). Each SLD runs an isolated CoreDNS container connected to the standard `dns-net` bridge network (`10.89.0.0/16`).

| Domain | SLD Container Name | Container IP | Host / Service Records | Target IP |
| :--- | :--- | :--- | :--- | :--- |
| **`example.x`** | `ns-example-x` | `10.89.20.10` | `@`, `blog`<br>`api` | `10.89.20.100`<br>`10.89.20.101` |
| **`shop.web`** | `ns-shop-web` | `10.89.30.10` | `@`, `portal`<br>`checkout` | `10.89.30.100`<br>`10.89.30.101` |
| **`store.shop`** | `ns-store-shop` | `10.89.40.10` | `api` | `10.89.40.101` |
| **`cluster.cloud`** | `ns-cluster-cloud` | `10.89.50.10` | `@`, `node1` | `10.89.50.100`<br>`10.89.50.101` |
| **`secure.mail`** | `ns-secure-mail` | `10.89.60.10` | `@`, `smtp` | `10.89.60.100`<br>`10.89.60.101` |
| **`lab.dev`** | `ns-lab-dev` | `10.89.70.10` | `@`, `git` | `10.89.70.100`<br>`10.89.70.101` |

---

## Directory Structure

```text
dns/sld/
├── example-x/
│   ├── Corefile
│   └── zone/
│       └── db.example.x
├── shop-web/
│   ├── Corefile
│   └── zone/
│       └── db.shop.web
├── shop-store/
│   ├── Corefile
│   └── zone/
│       └── db.store.shop
├── cloud-cluster/
│   ├── Corefile
│   └── zone/
│       └── db.cluster.cloud
├── mail-secure/
│   ├── Corefile
│   └── zone/
│       └── db.secure.mail
└── dev-lab/
    ├── Corefile
    └── zone/
        └── db.lab.dev
```

---

## Deployment

All SLD authoritative name servers are orchestrated via Ansible. Run the following playbook from your control machine to deploy or update the SLD containers:

```bash
ansible-playbook -i ~/Shadow_Web/inventory.ini ~/Shadow_Web/playbooks/deploy_dns_slds.yml
```

---

## Verification & Testing

You can test each SLD authoritative name server directly using an ephemeral Alpine container attached to the `dns-net` network:

### Direct Authoritative Queries
```bash
# Test example.x
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.20.10 blog.example.x +short"

# Test shop.web
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.30.10 portal.shop.web +short"

# Test store.shop
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.40.10 api.store.shop +short"

# Test cluster.cloud
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.50.10 node1.cluster.cloud +short"

# Test secure.mail
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.60.10 smtp.secure.mail +short"

# Test lab.dev
podman run --rm --net dns-net docker.io/library/alpine:latest sh -c "apk add --no-cache bind-tools > /dev/null && dig @10.89.70.10 git.lab.dev +short"
```

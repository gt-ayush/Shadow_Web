# Shadow_Web - Private Recursive Resolver (`resolver-01`)

This component provides the local recursive DNS resolution layer for the **Shadow_Web** private hierarchical infrastructure. It runs on **CoreDNS** to ensure complete cross-architecture compatibility (including native ARM64 support) and forwards unresolved recursive queries to the private root name servers (`root-a` and `root-b`).

---

## 🌐 Network & Addressing

* **Container Name**: `resolver-01`
* **Network**: `dns-net` (Bridge)
* **Subnet**: `10.89.10.0/24`
* **Static IP**: `10.89.10.10`
* **Upstream Roots**: `10.89.10.2` (Root-A), `10.89.10.4` (Root-B)

---

## 📁 Directory Structure

```bash
~/Shadow_Web/dns/resolver$ tree
.
├── issue.md
├── readme.md
└── resolver-01
    ├── Corefile
    ├── Dockerfile
    ├── compose.yaml
    └── config
        ├── root.hints
        └── unbound.conf

```

---

## ⚙️ Configuration Files

### 1. `Corefile`

Configures CoreDNS to act as a recursive forwarder with caching and error logging enabled.

```corefile
. {
    forward . 10.89.10.2 10.89.10.4 {
        policy random
    }
    cache 3600
    log
    errors
}

```

### 2. `compose.yaml`

Defines the Podman/Docker Compose specification for standalone deployment.

```yaml
services:
  resolver-01:
    image: docker.io/coredns/coredns:latest
    container_name: resolver-01
    restart: always
    command: -conf /etc/coredns/Corefile
    networks:
      dns-net:
        ipv4_address: 10.89.10.10
    ports:
      - "53:53/udp"
      - "53:53/tcp"
    volumes:
      - ./Corefile:/etc/coredns/Corefile:ro

networks:
  dns-net:
    external: true

```

---

## 🚀 Deployment

### Manual Deployment via Podman

Ensure the `dns-net` external network is created, then spin up the container:

```bash
podman network create dns-net --subnet 10.89.10.0/24

podman run --name resolver-01 \
  --network dns-net \
  --ip 10.89.10.10 \
  -v $(pwd)/Corefile:/etc/coredns/Corefile:ro \
  -d docker.io/coredns/coredns:latest -conf /etc/coredns/Corefile

```

### Automated Deployment via Ansible

Run the provided playbook from your project root:

```bash
ansible-playbook -i inventory.ini playbooks/deploy_dns_resolver.yml

```

---

## 🔍 Verification

Test the resolver directly from an ephemeral container on the `dns-net` network:

```bash
podman run --rm --network dns-net docker.io/library/alpine sh -c \
  "apk add --no-cache bind-tools >/dev/null && dig @10.89.10.10 a.root-servers.internal"

```

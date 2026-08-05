# Shadow_Web - Private Root DNS Infrastructure (`root-a` & `root-b`)

This component establishes the redundant, authoritative root zone (`.`) layer for the **Shadow_Web** private hierarchical DNS infrastructure. Built using **CoreDNS**, it guarantees cross-architecture compatibility (including native ARM64 execution) and serves foundational root records to local recursive resolvers.

---

## 🌐 Network & Addressing

* **Network**: `dns-net` (Bridge)
* **Subnet**: `10.89.10.0/24`
* **Root-A Container**: `root-a` (`10.89.10.2`)
* **Root-B Container**: `root-b` (`10.89.10.4`)

---

## 📁 Directory Structure

```bash
~/Shadow_Web/dns/root$ tree
.
├── readme.md
├── root-a
│   ├── Corefile
│   ├── compose.yaml
│   └── zone
│       └── db.root
└── root-b
    ├── Corefile
    ├── compose.yaml
    └── zone
        └── db.root
```

---

## ⚙️ Configuration Files

### 1. Root Zone File (`db.root`)

Defines the authoritative root zone records, delegating nameservers and binding their respective static IP addresses.

```zone
$ORIGIN .
@   IN  SOA     a.root-servers.internal. admin.shadow.web. (
                2026080301 ; serial
                7200       ; refresh
                3600       ; retry
                1209600    ; expire
                3600       ; minimum TTL
            )

    IN  NS      a.root-servers.internal.
    IN  NS      b.root-servers.internal.

a.root-servers.internal. IN A 10.89.10.2
b.root-servers.internal. IN A 10.89.10.4

```

### 2. Corefile (`Corefile`)

Instructs CoreDNS to load the authoritative root zone file.

```corefile
. {
    file /etc/coredns/db.root
    log
    errors
}

```

### 3. Compose Specification (`compose.yaml` - Example for Root-A)

```yaml
services:
  root-a:
    image: docker.io/coredns/coredns:latest
    container_name: root-a
    restart: always
    command: -conf /etc/coredns/Corefile
    networks:
      dns-net:
        ipv4_address: 10.89.10.2
    volumes:
      - ./Corefile:/etc/coredns/Corefile:ro
      - ./db.root:/etc/coredns/db.root:ro

networks:
  dns-net:
    external: true

```

---

## 🚀 Deployment

### Manual Deployment via Podman

Create the external bridge network if it doesn't already exist, then deploy both root nodes:

```bash
podman network create dns-net --subnet 10.89.10.0/24

# Deploy Root-A
podman run --name root-a --network dns-net --ip 10.89.10.2 \
  -v $(pwd)/root-a/Corefile:/etc/coredns/Corefile:ro \
  -v $(pwd)/root-a/db.root:/etc/coredns/db.root:ro \
  -d docker.io/coredns/coredns:latest -conf /etc/coredns/Corefile

# Deploy Root-B
podman run --name root-b --network dns-net --ip 10.89.10.4 \
  -v $(pwd)/root-b/Corefile:/etc/coredns/Corefile:ro \
  -v $(pwd)/root-b/db.root:/etc/coredns/db.root:ro \
  -d docker.io/coredns/coredns:latest -conf /etc/coredns/Corefile

```

---

## 🔍 Verification

Test authoritative responses directly from an ephemeral container on the network:

```bash
# Query Root-A
podman run --rm --network dns-net docker.io/library/alpine sh -c \
  "apk add --no-cache bind-tools >/dev/null && dig @10.89.10.2 . NS"

# Query Root-B
podman run --rm --network dns-net docker.io/library/alpine sh -c \
  "apk add --no-cache bind-tools >/dev/null && dig @10.89.10.4 . NS"

```

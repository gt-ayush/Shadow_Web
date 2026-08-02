**Naming Convention Policy**

We will document this naming specification inside our `docs/` folder to maintain consistency across all Ansible playbooks, Podman containers, network bridges, and persistent volumes.

---

## Standardized Naming Blueprint

### **1. Container Pods & Service Identifiers**

| Name | Type / Description | Target Network |
| --- | --- | --- |
| **`root-a`** | Primary Root DNS Server Instance | `dns-net` |
| **`root-b`** | Secondary/Redundant Root DNS Instance | `dns-net` |
| **`resolver-01`** | Recursive DNS Resolver Instance | `dns-net`, `hosting-net` |
| **`registry`** | Container Image Registry Service | `registry-net` |
| **`registrar`** | Domain Registration / Internal Management API | `hosting-net` |

---

### **2. Podman Isolated Bridge Networks**

| Network Name | Subnet Purpose | Connected Components |
| --- | --- | --- |
| **`dns-net`** | Isolated DNS infrastructure layer | `root-a`, `root-b`, `resolver-01` |
| **`hosting-net`** | Application backends, databases & APIs | `resolver-01`, `registrar` |
| **`registry-net`** | Image distribution & deployment layer | `registry` |

---

### **3. Persistent Named Volumes**

| Volume Name | Mounted Path Purpose | Target Pod / Service |
| --- | --- | --- |
| **`dns-data`** | Zone files, DNSSEC keys, and resolver state | `root-a`, `root-b`, `resolver-01` |
| **`registry-db`** | Storage bucket for raw container image blobs | `registry` |
| **`grafana-data`** | Monitoring dashboards, user preferences, SQLite DB | `grafana` *(Monitor Stack)* |


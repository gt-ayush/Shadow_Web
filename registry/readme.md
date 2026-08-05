# Shadow_Web — Registry & DNS Infrastructure

Shadow_Web is a containerized, registry-backed internet infrastructure platform designed to establish a fully automated registration-to-publishing pipeline for custom Top-Level Domains (TLDs).

---

## Architecture Stack

* **Core API**: Go-based microservice handling registrar authentication, business logic validation, and database updates.
* **Database**: PostgreSQL maintaining relational tables for domains, nameservers, owners, registrars, TLDs, and transaction audit trails.
* **Cache & Queue**: Redis for operational caching and task queuing.
* **Publisher Worker**: Automated daemon responsible for compiling and reloading custom TLD DNS zone files.
* **Orchestration**: Podman Compose (`podman-compose`) for container lifecycle management and Ansible for automated deployment.
* **Networking**: Isolated bridge network (`dns-net`, subnet `10.89.0.0/16`).

---

## Directory Structure

```text
Shadow_Web/
├── playbooks/
│   └── deploy_registry.yml       # Ansible deployment playbook
└── registry/
    ├── ansible/
    │   └── inventory.ini         # Ansible target definitions
    ├── api/                      # Go API source code & generated zones
    ├── compose/
    │   └── docker-compose.yml    # Podman Compose service definitions
    └── database/
        └── migrations/           # SQL schema initialization files

```

---

## Automated Deployment

The entire registry infrastructure can be deployed and verified automatically using the provided Ansible playbook.

### 1. Prerequisites

* Podman & Podman Compose installed locally.
* Ansible installed on the control node.

### 2. Run the Deployment Playbook

Navigate to your playbooks directory and execute the deployment playbook against your local environment:

```bash
cd ~/Shadow_Web/playbooks/
ansible-playbook -i ../registry/ansible/inventory.ini deploy_registry.yml

```

This playbook automates:

1. Creation of the standardized `dns-net` bridge network.
2. Verification and provisioning of mandatory directory structures.
3. Building and launching the container stack (`registry-db`, `registry-redis`, `registry-api`) via Podman Compose.
4. Health-check validation against the Registry API endpoint (`http://localhost:8081/health`).

---

## API Usage & Verification

To register a custom domain under a managed TLD (e.g., `.x`), send a JSON payload to the Registry API:

```bash
curl -X POST http://localhost:8081/domains/register \
  -H "Content-Type: application/json" \
  -d '{
    "domain": "example.x",
    "tld": "x",
    "owner_organization": "CyberCorp",
    "owner_email": "admin@cybercorp.x",
    "registrar": "ShadowRegistrarPrimary",
    "api_key": "sec_key_alpha_99887766",
    "nameservers": [
      {
        "hostname": "ns1.example.x",
        "ipv4": "10.89.20.10"
      }
    ]
  }'

```

### Inspecting Database Audit Logs

Verify transaction logs and recorded domain state directly inside the containerized PostgreSQL instance:

```bash
podman exec -it registry-db psql -U registry_admin -d tld_registry -c "SELECT * FROM transactions;"

```

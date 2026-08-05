# Shadow_Web DNS - Troubleshooting & Issue Ledger

This document logs the architectural issues, compatibility errors, and configuration hurdles encountered during the setup and automation of the `Shadow_Web` sovereign DNS infrastructure, along with their implemented resolutions.

---

## 1. Iterative Resolution Loop via Unresolvable TLD (`.internal`)
* **Problem:** Recursive resolution failed due to an infinite loop caused by root servers incorrectly identifying themselves using an unresolvable `.internal` TLD suffix (`a.root-servers.internal`).
* **Root Cause:** CoreDNS/Unbound interaction stalled because the nameservers pointed to non-existent internal domain paths that could not be resolved from root hints.
* **Fix:** Updated all zone files, root hints, and NS records to use clean, top-level root nomenclature (`a.root.` and `b.root.`) pointing directly to static internal IPs (`10.89.10.2` and `10.89.10.4`).

---

## 2. `dig +trace` Failure on Private Resolvers
* **Problem:** Running `dig @resolver-01 example.x +trace` failed with `couldn't get address for 'b.root': not found`.
* **Root Cause:** `+trace` forces `dig` to manually step through delegation trees. Because Unbound did not bundle glue records precisely how `dig +trace` expects in the immediate referral packet, `dig` tried a separate independent lookup for `a.root` / `b.root`.
* **Fix:** Verified standard recursive queries directly (`dig @resolver-01 example.x A`), confirming that normal DNS resolution works perfectly fine without manual trace loops.

---

## 3. Ansible Playbook Mismatch (CoreDNS vs. Unbound)
* **Problem:** The original deployment playbook targeted CoreDNS (`docker.io/coredns/coredns:latest`) for `resolver-01` instead of the migrated Unbound service.
* **Root Cause:** Infrastructure configuration drift between manual container migration and automation scripts.
* **Fix:** Rewrote `deploy_dns_resolver.yml` to support Unbound configuration files (`unbound.conf` and `root.hints`).

---

## 4. Ansible Inventory Parsing Error
* **Problem:** Running ansible-playbook resulted in warnings: `Unable to parse inventory.ini as an inventory source`, skipping target hosts (`dns_servers`).
* **Root Cause:** Missing or incorrectly formatted host definitions in the initialization inventory.
* **Fix:** Configured `inventory.ini` to explicitly target local execution for the `dns_servers` group:
  ```ini
  [dns_servers]
  localhost ansible_connection=local

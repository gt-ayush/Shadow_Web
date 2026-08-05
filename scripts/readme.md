# Shadow_Web: Backup & Automation System

This document outlines the architecture, configuration, and management of the automated backup and archival system for the `Shadow_Web` project.

---

## Overview

The backup system is designed to maintain operational continuity and data integrity by capturing a complete snapshot of the project repository, configuration files, and database states on a scheduled daily basis. It runs entirely under user privileges using **systemd user timers**, avoiding the need for root access while safely interacting with containerized services via Podman.

---

## System Components

* **Backup Script**: `scripts/local_backup.sh`
* **Archive Destination**: `archives/` (Stores compressed `.tar.gz` bundles)
* **Automation Engine**: Systemd User Services (`shadow-web-backup.service` and `shadow-web-backup.timer`)
* **Retention Policy**: Automatic pruning of backups older than **7 days**

---

## Core Script Workflow (`local_backup.sh`)

1. **Environment Setup**: Defines paths, generates a unique timestamp (`YYYY-MM-DD_HH-MM-SS`), and establishes retention limits.
2. **Database Dumper**: Executes `pg_dumpall` inside the active PostgreSQL container (`my_postgres_db`) to safely perform a logical backup of database states.
3. **Compression**: Combines the main project directory (`$HOME/Shadow_Web`) and the temporary database SQL dump into a single compressed tarball (`server_backup_[TIMESTAMP].tar.gz`).
4. **Cleanup**: Removes the raw SQL dump file and purges any compressed archives exceeding the 7-day retention threshold.

---

## Manual Execution

To trigger an immediate backup manually outside of the automated schedule:

```bash
bash ~/Shadow_Web/scripts/local_backup.sh

```

---

## Automation Management (Systemd User Timer)

The backup routine is automated using a systemd user timer, which runs persistently and handles missed schedules if the system was powered off.
### Create
```bash
mkdir -p ~/.config/systemd/user

cat << 'EOF' > ~/.config/systemd/user/shadow-web-backup.service
[Unit]
Description=Shadow_Web Automated Local Backup Service
After=network.target

[Service]
Type=oneshot
ExecStart=/home/min/Shadow_Web/scripts/local_backup.sh
EOF

cat << 'EOF' > ~/.config/systemd/user/shadow-web-backup.timer
[Unit]
Description=Run Shadow_Web Backup Daily

[Timer]
OnCalendar=daily
Persistent=true

[Install]
WantedBy=timers.target
EOF

systemctl --user daemon-reload
systemctl --user enable --now shadow-web-backup.timer
```
### Check Timer Status

To verify if the timer is active and check the next scheduled run:

```bash
systemctl --user list-timers

```

### View Service Logs

To inspect the execution history and logs of recent backups:

```bash
journalctl --user -u shadow-web-backup.service -b

```

### Stop or Disable Automation

To temporarily halt or permanently disable the automated daily backup:

* **Stop current timer instance:**
```bash
systemctl --user stop shadow-web-backup.timer

```


* **Disable completely:**
```bash
systemctl --user disable shadow-web-backup.timer

```

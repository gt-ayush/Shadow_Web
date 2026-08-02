#!/bin/bash

# 1. Configuration
BACKUP_DIR="$HOME/Shadow_Web/archives"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
ARCHIVE_NAME="server_backup_$TIMESTAMP.tar.gz"
RETENTION_DAYS=7

# Directories to back up (Config, Certs, DNS, Registry)
# Adjust these paths to where your actual volume data lives
SOURCE_DIRS="/path/to/certs /path/to/dns /path/to/registry $HOME/Shadow_Web"

echo "Starting backup at $TIMESTAMP..."

# 2. Database Dumps (Logical Backups)
# Safely extract data from running containers without stopping them
echo "Dumping databases..."
podman exec my_postgres_db pg_dumpall -U postgres > "$BACKUP_DIR/db_dump_$TIMESTAMP.sql"
# Example for MariaDB/MySQL:
# podman exec my_mariadb mariadb-dump -u root -p'PASSWORD' --all-databases > "$BACKUP_DIR/mysql_dump_$TIMESTAMP.sql"

# 3. Compress Data & Dumps into a single Tarball
echo "Compressing files..."
tar -czf "$BACKUP_DIR/$ARCHIVE_NAME" $SOURCE_DIRS "$BACKUP_DIR/db_dump_$TIMESTAMP.sql"

# 4. Clean up the raw SQL dumps (they are now inside the tarball)
rm "$BACKUP_DIR/db_dump_$TIMESTAMP.sql"

# 5. Prune Old Backups
echo "Removing backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "server_backup_*.tar.gz" -type f -mtime +$RETENTION_DAYS -delete

echo "Backup complete: $BACKUP_DIR/$ARCHIVE_NAME"

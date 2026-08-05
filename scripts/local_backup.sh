#!/bin/bash

# 1. Configuration
BACKUP_DIR="$HOME/Shadow_Web/archives"
TIMESTAMP=$(date +"%Y-%m-%d_%H-%M-%S")
ARCHIVE_NAME="server_backup_$TIMESTAMP.tar.gz"
RETENTION_DAYS=7

# Target the main project repository (covers dns, registry, configs, etc.)
SOURCE_DIRS="$HOME/Shadow_Web"

echo "Starting backup at $TIMESTAMP..."

# 2. Database Dumps (Logical Backups)
echo "Dumping databases..."
podman exec my_postgres_db pg_dumpall -U postgres > "$BACKUP_DIR/db_dump_$TIMESTAMP.sql" 2>/dev/null || echo "Warning: Database container not running or failed to dump."

# 3. Compress Data & Dumps into a single Tarball
echo "Compressing files..."
if [ -f "$BACKUP_DIR/db_dump_$TIMESTAMP.sql" ]; then
    tar -czf "$BACKUP_DIR/$ARCHIVE_NAME" $SOURCE_DIRS "$BACKUP_DIR/db_dump_$TIMESTAMP.sql"
    # Clean up the raw SQL dump (it is now inside the tarball)
    rm "$BACKUP_DIR/db_dump_$TIMESTAMP.sql"
else
    tar -czf "$BACKUP_DIR/$ARCHIVE_NAME" $SOURCE_DIRS
fi

# 4. Prune Old Backups
echo "Removing backups older than $RETENTION_DAYS days..."
find "$BACKUP_DIR" -name "server_backup_*.tar.gz" -type f -mtime +$RETENTION_DAYS -delete

echo "Backup complete: $BACKUP_DIR/$ARCHIVE_NAME"

#!/bin/bash
set -e

TYPO3_CLI="/var/www/html/typo3/sysext/core/bin/typo3"
SETTINGS_FILE="/var/www/html/typo3conf/system/settings.php"
FIRST_INSTALL_FILE="/var/www/html/FIRST_INSTALL"
DB_PORT="${TYPO3_DB_PORT:-3306}"
DB_HOST="${TYPO3_DB_HOST:-localhost}"
MAX_RETRIES=30
RETRY_INTERVAL=2

if [ -n "$TYPO3_SETUP_ADMIN_PASSWORD" ] && [ ! -f "$SETTINGS_FILE" ]; then
    echo "TYPO3 auto-setup: waiting for database at ${DB_HOST}:${DB_PORT}..."

    retries=0
    until php -r "@\$sock = fsockopen('${DB_HOST}', (int)'${DB_PORT}', \$errno, \$errstr, 1); if (!\$sock) exit(1); fclose(\$sock);" 2>/dev/null; do
        retries=$((retries + 1))
        if [ "$retries" -ge "$MAX_RETRIES" ]; then
            echo "TYPO3 auto-setup: database not reachable after $((MAX_RETRIES * RETRY_INTERVAL))s, giving up."
            break
        fi
        sleep "$RETRY_INTERVAL"
    done

    if [ "$retries" -lt "$MAX_RETRIES" ]; then
        echo "TYPO3 auto-setup: database is ready, running setup..."
        runuser -u www-data -- "$TYPO3_CLI" setup --force --no-interaction --server-type=apache
        rm -f "$FIRST_INSTALL_FILE"
        echo "TYPO3 auto-setup: complete."
    fi
fi

exec "$@"

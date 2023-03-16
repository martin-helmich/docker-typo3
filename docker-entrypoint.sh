#!/usr/bin/env bash
for dir in fileadmin typo3conf typo3temp uploads;
do
  chown www-data:www-data \${DOCUMENT_ROOT:-/var/www/html}/\${dir}
done
[ -d /docker-entrypoint.d/ ] && [ $(ls -1 /docker-entrypoint.d/*.sh 2> /dev/null) ] && {
  chmod +x /docker-entrypoint.d/*.sh
  run-parts --report --regex='\.sh' /docker-entrypoint.d
}

sh /usr/local/bin/docker-php-entrypoint "$@"

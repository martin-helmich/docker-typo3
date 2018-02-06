#!/bin/bash
set -euo pipefail

if [[ "$1" == apache2* ]] || [ "$1" == php-fpm ]; then
	if ! [ -e index.php -a -e typo3/index.php ]; then
		echo >&2 "TYPO3 not found in $PWD - copying now..."
		if [ "$(ls -A)" ]; then
			echo >&2 "WARNING: $PWD is not empty - press Ctrl+C now if this is an error!"
			( set -x; ls -A; sleep 10 )
		fi
		tar cf - --one-file-system -C /usr/src/typo3 . | tar xf -
		ln -s typo3_src-* typo3_src
		ln -s typo3_src/index.php
		ln -s typo3_src/typo3
		ln -s typo3_src/_.htaccess .htaccess
		touch FIRST_INSTALL
		chown -R www-data:www-data FIRST_INSTALL
		echo >&2 "Complete! TYPO3 has been successfully copied to $PWD"
	fi
fi

exec "$@"

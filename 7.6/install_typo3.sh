
if [ ! -d /var/www/html/typo3/sysext ]
    then
	
	cd /var/www/html && \
	wget -O - https://get.typo3.org/$1 | tar -xzf - && \
    ln -s typo3_src-* typo3_src && \
    ln -s typo3_src/index.php && \
    ln -s typo3_src/typo3 && \
    ln -s typo3_src/_.htaccess .htaccess && \
    touch FIRST_INSTALL && \
    chown -R www-data /var/www/html
    
fi

apache2-foreground
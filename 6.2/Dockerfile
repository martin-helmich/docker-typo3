FROM php:5.6-apache
MAINTAINER Martin Helmich <typo3@martin-helmich.de>

# Install TYPO3
RUN apt-get update &&\
    apt-get install -y wget

# Configure PHP
RUN apt-get install -y \
        libxml2-dev libfreetype6-dev \
        libjpeg62-turbo-dev \
        libmcrypt-dev \
        libpng12-dev \
        zlib1g-dev
RUN docker-php-ext-configure gd --with-freetype-dir=/usr/include/ --with-jpeg-dir=/usr/include/ && \
    docker-php-ext-install -j$(nproc) mysqli soap gd zip
RUN echo 'always_populate_raw_post_data = -1\nmax_execution_time = 240\nmax_input_vars = 1500\nupload_max_filesize = 32M\npost_max_size = 32M' > /usr/local/etc/php/conf.d/typo3.ini

# Install required 3rd party tools
RUN apt-get install -y graphicsmagick

# Configure Apache as needed
RUN a2enmod rewrite

# Configure volumes
VOLUME /var/www/html


ADD install_typo3.sh /install_typo3.sh
RUN chmod 755 /install_typo3.sh

CMD /install_typo3.sh 6.2
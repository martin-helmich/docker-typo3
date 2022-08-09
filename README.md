TYPO3 Quickstart in Docker
==========================

[![Docker Automated build](https://img.shields.io/docker/automated/martinhelmich/typo3.svg)](https://hub.docker.com/r/martinhelmich/typo3/)
[![Docker build status](https://img.shields.io/docker/build/martinhelmich/typo3.svg)](https://hub.docker.com/r/martinhelmich/typo3/)

This repository contains build instructions for a simple TYPO3 Docker image.

**Note** that this image is not intended for production usage (yet). It's goal is to provide users an easy quickstart for working with TYPO3.

Usage
-----

This container does not ship a database management system; which means you'll have to create your own database container. The upside of this is that you're not bound to any specific version of MySQL or MariaDB and can even use a PostgreSQL database if you like.

1. So, the first step should be to create a database container:

        $ docker run -d --name typo3-db \
           -e MYSQL_ROOT_PASSWORD=yoursupersecretpassword \
           -e MYSQL_USER=typo3 \
           -e MYSQL_PASSWORD=yourothersupersecretpassword \
           -e MYSQL_DATABASE=typo3 \
           mariadb:latest \
           --character-set-server=utf8mb4 \
           --collation-server=utf8mb4_unicode_ci

2. Next, use this image to create your TYPO3 container and link it with the database container:

        $ docker run -d --name typo3-web \
            --link typo3-db:db \
            -p 80:80 \
            martinhelmich/typo3:11

3. After that, simply open `http://localhost/` in your browser to start the TYPO3 install tool. **Note**: If you're using Docker Machine to run Docker on Windows or MacOS, you'll need the Docker VM's IP instead (which you can find out using the `docker-machine ip default` command).

4. Complete the install tool. When prompted for database credentials, use the environment variables that you've passed to the database container in step 1. If you've linked the containers using the `--link` flag as shown in step 2, use `db` as database host name.

 ![](doc/database-setup.png)

Special use cases
-----------------

### Using PostgreSQL instead of MySQL

This image comes with full support for PostgreSQL as database driver. In this case, simply start a PostgreSQL database instead of MySQL:

    $ docker run -d --name typo3-db \
        -e POSTGRES_PASSWORD=yoursupersecretpassword \
        -e POSTGRES_USER=typo3 \
        -e POSTGRES_DATABASE=typo3 \
        postgres:latest

Then, proceed as before.

Available tags
--------------

This repository offers the following image tags:

- `latest` maps to the latest available LTS version (currently, latest `11.5.*`)
- `11.5` and `11` for the latest available version from the `11.*` respectively `11.5.*` branch.
- `10.4` and `10` for the latest available version from the `10.*` respectively `10.4.*` branch.
- `9.5` and `9` for the latest available version from the `9.*` respectively `9.5.*` branch.
- `8.7` and `8` for the latest available version from the `8.*` respectively `8.7.*` branch.
- `7.6` and `7` for the latest available version from the `7.*` respectively `7.6.*` branch.
- `6.2` and `6` for the latest available version from the `6.*` respectively `6.2.*` branch.

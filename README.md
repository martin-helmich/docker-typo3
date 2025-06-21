TYPO3 Quickstart in Docker
==========================

[![Docker Pulls](https://img.shields.io/docker/pulls/martinhelmich/typo3)](https://hub.docker.com/r/martinhelmich/typo3)
[![Docker Image Version (tag latest semver)](https://img.shields.io/docker/v/martinhelmich/typo3/latest?logo=docker)](https://hub.docker.com/r/martinhelmich/typo3)
[![Build and push images](https://github.com/martin-helmich/docker-typo3/actions/workflows/build.yml/badge.svg)](https://github.com/martin-helmich/docker-typo3/actions/workflows/build.yml)

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
            martinhelmich/typo3:13

3. After that, simply open `http://localhost/` in your browser to start the TYPO3 install tool. **Note**: If you're using Docker Machine to run Docker on Windows or MacOS, you'll need the Docker VM's IP instead (which you can find out using the `docker-machine ip default` command).

4. Complete the install tool. When prompted for database credentials, use the environment variables that you've passed to the database container in step 1. If you've linked the containers using the `--link` flag as shown in step 2, use `db` as database host name.

 ![](doc/database-setup.png)
 
Volumes to be Mounted
~~~~~~~~~~~~~~~~~~~~~

The Docker images expose the following volumes, which **should** be mounted to ensure proper functionality and data persistence:

- `/var/www/html/fileadmin` — Contains files and images uploaded by editors.
- `/var/www/html/typo3conf` — Contains uploaded extensions, configuration, and language files.
- `/var/www/html/typo3temp` — Contains temporary files such as logs that may be useful for debugging.
- `/var/www/html/uploads` — Deprecated since TYPO3 6.2 and unused after TYPO3 10.0. Can be ignored for newer installations.

If you are using the supplied `docker-compose.yml` file, these volumes are mounted automatically.

If you are using a custom configuration strategy, you must ensure these volumes are mounted manually.

Exposed Ports and Reverse Proxies
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The Docker image exposes **port 80**. If you are using the provided `docker-compose.yml`, this is automatically mapped to port 80 on your host system, allowing you to access the site at [http://localhost](http://localhost).

If you are using a different setup, you may need to manually map port 80 to a suitable host port.

If you plan to use HTTPS/TLS, you must additionally configure your reverse proxy. For more details, refer to the TYPO3 documentation:  
[Reverse Proxy Container Configuration](https://docs.typo3.org/permalink/t3coreapi:reverse-proxy-container)


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

- `latest` maps to the latest available LTS version (currently, latest `13.4.*`)
- `13.4` and `13` for the latest available version from the `13.*` respectively `13.4.*` branch.
- `12.4` and `12` for the latest available version from the `12.*` respectively `12.4.*` branch.

The following tags are still available, but not updated any longer:

- `11.5` and `11` for the latest available version from the `11.*` respectively `11.5.*` branch.
- `10.4` and `10` for the latest available version from the `10.*` respectively `10.4.*` branch.
- `9.5` and `9` for the latest available version from the `9.*` respectively `9.5.*` branch.
- `8.7` and `8` for the latest available version from the `8.*` respectively `8.7.*` branch.
- `7.6` and `7` for the latest available version from the `7.*` respectively `7.6.*` branch.
- `6.2` and `6` for the latest available version from the `6.*` respectively `6.2.*` branch.

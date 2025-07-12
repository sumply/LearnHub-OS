# Dockerfile (Single image for PHP app + Pre-initialized PostgreSQL)

# --- Stage 1: Build the pre-initialized PostgreSQL database ---
FROM postgres:17 AS init_db
ENV POSTGRES_DB main
ENV POSTGRES_USER pdadmin
ENV POSTGRES_PASSWORD pgadmin
COPY ./init.sql /docker-entrypoint-initdb.d/
CMD ["postgres"] 

# --- Stage 2: PHP Application Builder ---
FROM ubuntu:plucky AS php_builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        software-properties-common ca-certificates curl git unzip
RUN apt-get install -y --no-install-recommends \
    php8.4-cli php8.4-fpm php8.4-common php8.4-mbstring php8.4-xml \
    php8.4-gd php8.4-intl php8.4-pgsql php8.4-redis composer && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
WORKDIR /app
COPY composer.json ./
RUN composer install --no-interaction --no-dev --prefer-dist --optimize-autoloader
COPY . . 

# --- Final Stage: Combine PHP App and Pre-initialized DB ---
FROM ubuntu:plucky AS final_app_with_db

# 1. Install necessary system dependencies for PostgreSQL client and PHP runtime
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        postgresql-client \
        libpq-dev \
        php8.4-fpm \
        php8.4-cli \
        php8.4-common \
        php8.4-mbstring \
        php8.4-xml \
        php8.4-gd \
        php8.4-intl \
        php8.4-pgsql \
        php8.4-redis && \
    apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*


# 2. Copy the initialized DB data from the init_db stage
#    This makes the DB part of THIS image
COPY --from=init_db /var/lib/postgresql/data /var/lib/postgresql/data

# 3. Copy the PHP application and its dependencies from the php_builder stage
COPY --from=php_builder /app /app

# 4. Set up necessary environment variables for PostgreSQL
#    These are for the *internal* DB instance
ENV PGDATA=/var/lib/postgresql/data
ENV POSTGRES_DB main
ENV POSTGRES_USER pdadmin
ENV POSTGRES_PASSWORD pgadmin

# 5. Expose the port for the database and for PHP-FPM
EXPOSE 5432
EXPOSE 9000

# 6. Set working directory for the application
WORKDIR /app

# 7. Define the entrypoint script to start both PostgreSQL and PHP-FPM
#    This is a bit more complex, as you need a wrapper script.

CMD ["php8.4", "-S", "0.0.0.0:9000", "-t",".", "-c","php.ini", "index.php"]
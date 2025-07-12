# --- Builder Stage ---
# This stage has all the tools needed to install PHP dependencies with Composer.
# Using Ubuntu 22.04 LTS for stability.
FROM ubuntu:plucky AS builder

# 1. Set a non-interactive frontend for apt to prevent it from prompting for input.
ENV DEBIAN_FRONTEND=noninteractive

# 2. Install PHP 8.3, Composer, and buil4 tools.
# We need to add the ondrej/php PPA to get the latest PHP versions.
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        software-properties-common \
        ca-certificates \
        curl \
        php8.4-common \
        php8.4-cli \
        php8.4-fpm \
        # --- Add required PHP extensions here (uncomment what you need) ---
        # Most applications need mbstring and xml.
        php8.4-mbstring \
        php8.4-xml \
        php8.4-gd \
        php8.4-intl \
        php8.4-pgsql \
        # php8.4-mysql \
        php8.4-redis \
        apt-transport-https && \
    apt-get update && \
    # Install PHP, Composer, and other necessary tools for building
    apt-get install -y --no-install-recommends \
        php8.4-cli \
        php8.4-fpm \
        composer \
        git \
        unzip
# 3. Set up the application directory
ENV DEBIAN_FRONTEND=noninteractive
RUN install -d /usr/share/postgresql-common/pgdg
RUN curl -o /usr/share/postgresql-common/pgdg/apt.postgresql.org.asc --fail https://www.postgresql.org/media/keys/ACCC4CF8.asc
RUN sh -c "echo 'deb [signed-by=/usr/share/postgresql-common/pgdg/apt.postgresql.org.asc] https://apt.postgresql.org/pub/repos/apt $VERSION_CODENAME-pgdg main' > /etc/apt/sources.list.d/pgdg.list"
RUN . /etc/os-release
RUN apt install -y --no-install-recommends postgresql-17
WORKDIR /app

# 4. Copy only dependency files and install them.
# This leverages Docker's layer cache. `composer install` only re-runs if these files change.
COPY composer.json ./
RUN composer install --no-interaction --no-dev --prefer-dist --optimize-autoloader
# 5. Copy the rest of the application source code
COPY . .

# --- Final Application Stage ---
# This stage creates the final, lean image for running the application.
FROM ubuntu:plucky AS app

# 1. Install only the required PHP runtime and common extensions.
RUN apt-get update && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 2. Set up the application directory
WORKDIR /app

# 3. Copy the application code and installed dependencies from the builder stage
COPY --from=builder /app .

# 4. Expose port if it's a web service running on php-fpm
EXPOSE 9000

# 5. Define the entrypoint
# Choose one of the following based on your application type.

# --- OPTION A: For a command-line service (like your original C++ service) ---
# ENTRYPOINT ["php"]
CMD ["php8.4", "-S","0.0.0.0:9000","-t",".","index.php" ]

# --- OPTION B: For a web application served by PHP-FPM ---
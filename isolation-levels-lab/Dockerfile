FROM postgres:16-alpine

# Set environment variables
ENV POSTGRES_DB=isolation_db
ENV POSTGRES_USER=isolation_user
ENV POSTGRES_PASSWORD=isolation_pass

# Install PostgreSQL client tools
RUN apk add --no-cache postgresql-client

# Expose PostgreSQL port
EXPOSE 5432

# Add any custom initialization scripts if needed
COPY init.sql /docker-entrypoint-initdb.d/ 
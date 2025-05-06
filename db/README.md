# Database Schema for Blog Service

This directory contains the database schema and migration scripts for the Blog service defined in `protos/blog/v1/blog.proto`.

## Schema Overview

The database schema consists of two main tables:

1. **blogs** - Stores blog posts with the following columns:
   - `id` (UUID, primary key)
   - `title` (VARCHAR, max 100 chars)
   - `content` (TEXT, max 10000 chars)
   - `created_at` (TIMESTAMP WITH TIME ZONE)
   - `updated_at` (TIMESTAMP WITH TIME ZONE)

2. **comments** - Stores comments on blog posts with the following columns:
   - `id` (UUID, primary key)
   - `blog_id` (UUID, foreign key to blogs.id)
   - `content` (TEXT, max 1000 chars)
   - `author` (VARCHAR, max 50 chars)
   - `created_at` (TIMESTAMP WITH TIME ZONE)

## Flyway Migration

This project uses [Flyway](https://flywaydb.org/) for database migrations. The migration scripts are located in the `migrations` directory.

### Configuration

The Flyway configuration is stored in `flyway.conf`. You may need to update the database connection details:

```
flyway.url=jdbc:postgresql://localhost:5432/blog_db
flyway.user=postgres
flyway.password=postgres
```

### Running Migrations

The following commands are available in the Makefile:

- `make db-migrate` - Run all pending migrations
- `make db-clean` - Clean the database (drop all objects)
- `make db-info` - Show information about migrations
- `make db-validate` - Validate applied migrations
- `make db-repair` - Repair the schema history table

### Creating New Migrations

To create a new migration, add a new SQL file to the `migrations` directory following the Flyway naming convention:

```
V{version}__{description}.sql
```

For example:
- `V1__initial_schema.sql` - Initial schema creation
- `V2__add_tags_to_blogs.sql` - Adding tags to blogs

## Database Setup

Before running migrations, ensure you have a PostgreSQL database created:

```bash
# Create the database
createdb blog_db

# Or using psql
psql -U postgres -c "CREATE DATABASE blog_db;"
```

Then run the migrations:

```bash
make db-migrate
```
# Environment Variables

This document describes all environment variables used by the Iivineri API application.

## Setup

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your specific configuration values.

## Application Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `ENV` | Application environment | `development` | `development`, `production` |
| `PORT` | HTTP server port | `8080` | `8080` |
| `LOG_LEVEL` | Logging level | `info` | `debug`, `info`, `warn`, `error` |
| `PREFORK` | Enable Fiber prefork mode | `false` | `true`, `false` |

### Environment Details

- **`development`**: Enables debug features, stack traces, verbose logging
- **`production`**: Optimized for performance, minimal logging, rate limiting enabled

## Database Configuration

| Variable | Description | Default | Example |
|----------|-------------|---------|---------|
| `DB_HOST` | Database host | `localhost` | `localhost`, `db.example.com` |
| `DB_PORT` | Database port | `5432` | `5432` |
| `DB_NAME` | Database name | `iivineri_dev` | `iivineri_dev`, `iivineri_prod` |
| `DB_USERNAME` | Database username | `iivineri_user` | `iivineri_user` |
| `DB_PASSWORD` | Database password | - | `your_secure_password` |

### Connection Pool Settings

| Variable | Description | Default | Recommended |
|----------|-------------|---------|-------------|
| `DB_MIN_CONNS` | Minimum connections | `5` | `5-10` |
| `DB_MAX_CONNS` | Maximum connections | `25` | `25-50` |
| `DB_MIN_IDLE_CONNS` | Minimum idle connections | `5` | `5-10` |

### Database Logging

| Variable | Description | Default | Options |
|----------|-------------|---------|---------|
| `DB_LOG_LEVEL` | Database query logging | `error` | `trace`, `debug`, `info`, `warn`, `error` |

- **`trace`**: Logs all SQL queries (development only)
- **`error`**: Only logs database errors (recommended for production)

## Future Configuration

These variables are prepared for future features:

### Cache (Redis)
```env
CACHE_HOST=localhost
CACHE_PORT=6379
CACHE_PASSWORD=
CACHE_DB=0
```

### Thumbor (Image Processing)
```env
THUMBOR_URL=http://localhost:8888
THUMBOR_SECURITY_KEY=your_thumbor_key_here
```

## Production Recommendations

For production deployment:

```env
ENV=production
LOG_LEVEL=warn
PREFORK=true
DB_LOG_LEVEL=error
DB_MAX_CONNS=50
```

### Security Notes

- Never commit `.env` files to version control
- Use strong, unique passwords for database connections
- Consider using environment-specific database users with minimal permissions
- In production, consider using secrets management systems instead of `.env` files

## Docker Compose Example

```yaml
version: '3.8'
services:
  api:
    build: .
    environment:
      - ENV=production
      - DB_HOST=postgres
      - DB_NAME=iivineri
      - DB_USERNAME=iivineri_user
      - DB_PASSWORD=${DB_PASSWORD}
    depends_on:
      - postgres
  
  postgres:
    image: postgres:15
    environment:
      - POSTGRES_DB=iivineri
      - POSTGRES_USER=iivineri_user
      - POSTGRES_PASSWORD=${DB_PASSWORD}
```
# Iivineri API - Metrics Documentation

## Overview

The Iivineri API implements comprehensive monitoring and observability using **Prometheus metrics**. This document describes all metrics exposed by the application, their purpose, and how to use them for monitoring and alerting.

## Metrics Endpoint

All metrics are exposed at: `GET /metrics`

This endpoint returns metrics in Prometheus format, ready for scraping by Prometheus server.

## Metrics Categories

### 1. HTTP Metrics

#### `http_requests_total`
- **Type**: Counter
- **Description**: Total number of HTTP requests processed
- **Labels**:
  - `method`: HTTP method (GET, POST, PUT, DELETE, etc.)
  - `path`: Route pattern (e.g., `/api/v1`, `/health`, `/metrics`)
  - `status_code`: HTTP response status code (200, 404, 500, etc.)

**Example**:
```prometheus
http_requests_total{method="GET",path="/health",status_code="200"} 42
http_requests_total{method="GET",path="/api/v1",status_code="200"} 15
http_requests_total{method="POST",path="/api/v1/tracks",status_code="201"} 8
```

#### `http_request_duration_seconds`
- **Type**: Histogram
- **Description**: HTTP request duration in seconds
- **Labels**:
  - `method`: HTTP method
  - `path`: Route pattern
- **Buckets**: Default Prometheus buckets (.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10, +Inf)

**Example**:
```prometheus
http_request_duration_seconds_bucket{method="GET",path="/api/v1",le="0.005"} 10
http_request_duration_seconds_bucket{method="GET",path="/api/v1",le="0.01"} 15
http_request_duration_seconds_sum{method="GET",path="/api/v1"} 0.123
http_request_duration_seconds_count{method="GET",path="/api/v1"} 15
```

#### `http_requests_in_flight`
- **Type**: Gauge
- **Description**: Number of HTTP requests currently being processed
- **Labels**: None

**Example**:
```prometheus
http_requests_in_flight 3
```

---

### 2. Database Metrics

#### `db_query_duration_seconds`
- **Type**: Histogram
- **Description**: Database query execution time in seconds
- **Labels**:
  - `operation`: Type of database operation (health_check, select, insert, update, delete)
- **Buckets**: [0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10, +Inf]

**Example**:
```prometheus
db_query_duration_seconds_bucket{operation="health_check",le="0.001"} 45
db_query_duration_seconds_bucket{operation="select",le="0.01"} 120
db_query_duration_seconds_sum{operation="health_check"} 0.856
db_query_duration_seconds_count{operation="health_check"} 50
```

#### `db_queries_total`
- **Type**: Counter
- **Description**: Total number of database queries executed
- **Labels**:
  - `operation`: Type of database operation
  - `status`: Query result status (success, error)

**Example**:
```prometheus
db_queries_total{operation="health_check",status="success"} 48
db_queries_total{operation="health_check",status="error"} 2
db_queries_total{operation="select",status="success"} 156
```

#### `db_connections_active`
- **Type**: Gauge
- **Description**: Number of active database connections in the pool
- **Labels**: None

**Example**:
```prometheus
db_connections_active 12
```

---

### 3. Business Metrics

#### `tracks_played_total`
- **Type**: Counter
- **Description**: Total number of tracks played
- **Labels**:
  - `track_id`: Unique identifier of the track

**Example**:
```prometheus
tracks_played_total{track_id="track_001"} 25
tracks_played_total{track_id="track_042"} 18
tracks_played_total{track_id="track_156"} 7
```

#### `user_actions_total`
- **Type**: Counter
- **Description**: Total number of user actions performed
- **Labels**:
  - `action`: Type of action performed

**Common actions**:
- `api_request`: General API endpoint access
- `health_check`: Health endpoint access
- `track_play`: Track playback initiated
- `track_pause`: Track playback paused
- `track_skip`: Track skipped
- `playlist_create`: Playlist created
- `user_login`: User authentication
- `user_logout`: User session ended

**Example**:
```prometheus
user_actions_total{action="api_request"} 89
user_actions_total{action="health_check"} 156
user_actions_total{action="track_play"} 42
user_actions_total{action="user_login"} 12
```

#### `cache_operations_total`
- **Type**: Counter
- **Description**: Total number of cache operations
- **Labels**:
  - `cache_type`: Type of cache (redis, memory, file)
  - `result`: Operation result (hit, miss)

**Example**:
```prometheus
cache_operations_total{cache_type="redis",result="hit"} 234
cache_operations_total{cache_type="redis",result="miss"} 45
cache_operations_total{cache_type="memory",result="hit"} 567
```

---

### 4. System Metrics

#### `memory_usage_bytes`
- **Type**: Gauge
- **Description**: Current memory usage in bytes (allocated heap memory)
- **Labels**: None

**Example**:
```prometheus
memory_usage_bytes 45728768
```

#### `goroutines_active`
- **Type**: Gauge
- **Description**: Number of active goroutines
- **Labels**: None

**Example**:
```prometheus
goroutines_active 25
```

---

## Metrics Collection

### Automatic Collection

The following metrics are automatically collected:

1. **HTTP Metrics**: Collected by Fiber middleware on every request
2. **System Metrics**: Collected every 30 seconds by background collector
3. **Database Health**: Checked every 30 seconds with timing metrics

### Manual Collection

Business metrics need to be manually recorded in your application code:

```go
// Record track play
container.Metrics.IncrementTrackPlayed("track_123")

// Record user action
container.Metrics.IncrementUserAction("playlist_create")

// Record cache operation
container.Metrics.RecordCacheHit("redis", true) // hit
container.Metrics.RecordCacheHit("redis", false) // miss

// Record database query with timing
start := time.Now()
// ... execute query ...
duration := time.Since(start)
container.Metrics.RecordDBQuery("select", duration, err == nil)
```

---

## Prometheus Configuration

### Scrape Configuration

Add this job to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'iivineri-api'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/metrics'
    scrape_interval: 15s
    scrape_timeout: 10s
```

### Recording Rules

Useful recording rules for dashboards:

```yaml
groups:
  - name: iivineri_api
    rules:
      # Request rate (per second)
      - record: iivineri:http_requests:rate5m
        expr: rate(http_requests_total[5m])
      
      # Error rate percentage
      - record: iivineri:http_errors:rate5m
        expr: rate(http_requests_total{status_code=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100
      
      # Average response time
      - record: iivineri:http_duration:avg5m
        expr: rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])
      
      # Database query rate
      - record: iivineri:db_queries:rate5m
        expr: rate(db_queries_total[5m])
```

---

## Alerting Rules

### Critical Alerts

```yaml
groups:
  - name: iivineri_critical
    rules:
      # High error rate
      - alert: HighErrorRate
        expr: iivineri:http_errors:rate5m > 5
        for: 2m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }}% over the last 5 minutes"
      
      # High response time
      - alert: HighResponseTime
        expr: iivineri:http_duration:avg5m > 1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High response time detected"
          description: "Average response time is {{ $value }}s over the last 5 minutes"
      
      # Database connection issues
      - alert: DatabaseConnectionLow
        expr: db_connections_active < 3
        for: 1m
        labels:
          severity: warning
        annotations:
          summary: "Low database connections"
          description: "Only {{ $value }} database connections active"
      
      # High memory usage
      - alert: HighMemoryUsage
        expr: memory_usage_bytes > 100 * 1024 * 1024  # 100MB
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High memory usage"
          description: "Memory usage is {{ $value | humanizeBytes }}"
```

---

## Grafana Dashboard

### Key Panels

1. **HTTP Request Rate**: `rate(http_requests_total[5m])`
2. **HTTP Error Rate**: `rate(http_requests_total{status_code=~"5.."}[5m])`
3. **Response Time (95th percentile)**: `histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))`
4. **Database Query Duration**: `rate(db_query_duration_seconds_sum[5m]) / rate(db_query_duration_seconds_count[5m])`
5. **Active Connections**: `db_connections_active`
6. **Memory Usage**: `memory_usage_bytes`
7. **Goroutines**: `goroutines_active`
8. **Top Played Tracks**: `topk(10, rate(tracks_played_total[1h]))`

### Sample Query Examples

```promql
# Request rate by endpoint
sum(rate(http_requests_total[5m])) by (path)

# Error rate by status code
sum(rate(http_requests_total{status_code!~"2.."}[5m])) by (status_code)

# Database performance
histogram_quantile(0.95, rate(db_query_duration_seconds_bucket[5m]))

# Most popular tracks in last hour
topk(10, increase(tracks_played_total[1h]))

# User activity heatmap
sum(rate(user_actions_total[5m])) by (action)
```

---

## Troubleshooting

### Common Issues

1. **No metrics data**: Check if `/metrics` endpoint is accessible
2. **Missing labels**: Ensure middleware is properly configured
3. **High cardinality**: Avoid using dynamic values as labels (use track_id sparingly)
4. **Stale metrics**: System collector runs every 30s, some metrics may have delays

### Health Checks

Monitor these key metrics for application health:

- `http_requests_total`: Should be increasing
- `db_connections_active`: Should be > 0
- `memory_usage_bytes`: Should be stable, not constantly growing
- `goroutines_active`: Should be reasonable (< 1000 for typical load)

### Performance Impact

Metrics collection has minimal performance impact:
- HTTP middleware: ~0.1ms overhead per request
- System collection: Runs every 30s in background
- Memory overhead: ~1-2MB for metrics storage

---

## Future Enhancements

Planned metrics additions:

1. **Cache metrics**: Redis/Memory cache hit ratios
2. **External API metrics**: Third-party service latencies
3. **Queue metrics**: Background job processing stats
4. **Custom business metrics**: Genre popularity, user engagement
5. **Resource metrics**: CPU usage, disk I/O

---

For questions about metrics implementation, see the source code in:
- `internal/metrics/metrics.go` - Core metrics definitions
- `internal/metrics/middleware.go` - HTTP middleware
- `internal/metrics/collector.go` - System metrics collector
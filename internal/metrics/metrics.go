package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsInterface interface {
	// HTTP Metrics
	RecordHTTPRequest(method, path string, statusCode int, duration time.Duration)
	IncrementHTTPRequests(method, path string, statusCode int)
	ObserveHTTPDuration(method, path string, duration time.Duration)
	
	// Database Metrics
	RecordDBQuery(operation string, duration time.Duration, success bool)
	IncrementDBConnections()
	DecrementDBConnections()
	
	// Business Metrics
	IncrementTrackPlayed(trackID string)
	IncrementUserAction(action string)
	RecordCacheHit(cacheType string, hit bool)
	
	// System Metrics
	RecordMemoryUsage(bytes uint64)
	RecordGoroutines(count int)
}

type Metrics struct {
	// HTTP Metrics
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	httpRequestsInFlight prometheus.Gauge
	
	// Database Metrics
	dbQueryDuration     *prometheus.HistogramVec
	dbQueriesTotal      *prometheus.CounterVec
	dbConnectionsActive prometheus.Gauge
	
	// Business Metrics
	tracksPlayedTotal   *prometheus.CounterVec
	userActionsTotal    *prometheus.CounterVec
	cacheHitsTotal      *prometheus.CounterVec
	
	// System Metrics
	memoryUsageBytes prometheus.Gauge
	goroutinesActive prometheus.Gauge
}

func NewMetrics() MetricsInterface {
	m := &Metrics{
		// HTTP Metrics
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status_code"},
		),
		
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		
		httpRequestsInFlight: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_requests_in_flight",
				Help: "Number of HTTP requests currently being processed",
			},
		),
		
		// Database Metrics
		dbQueryDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "db_query_duration_seconds",
				Help:    "Database query duration in seconds",
				Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
			},
			[]string{"operation"},
		),
		
		dbQueriesTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "db_queries_total",
				Help: "Total number of database queries",
			},
			[]string{"operation", "status"},
		),
		
		dbConnectionsActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "db_connections_active",
				Help: "Number of active database connections",
			},
		),
		
		// Business Metrics
		tracksPlayedTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "tracks_played_total",
				Help: "Total number of tracks played",
			},
			[]string{"track_id"},
		),
		
		userActionsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "user_actions_total",
				Help: "Total number of user actions",
			},
			[]string{"action"},
		),
		
		cacheHitsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_operations_total",
				Help: "Total number of cache operations",
			},
			[]string{"cache_type", "result"},
		),
		
		// System Metrics
		memoryUsageBytes: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "memory_usage_bytes",
				Help: "Current memory usage in bytes",
			},
		),
		
		goroutinesActive: promauto.NewGauge(
			prometheus.GaugeOpts{
				Name: "goroutines_active",
				Help: "Number of active goroutines",
			},
		),
	}
	
	return m
}

// HTTP Metrics
func (m *Metrics) RecordHTTPRequest(method, path string, statusCode int, duration time.Duration) {
	m.IncrementHTTPRequests(method, path, statusCode)
	m.ObserveHTTPDuration(method, path, duration)
}

func (m *Metrics) IncrementHTTPRequests(method, path string, statusCode int) {
	m.httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(statusCode)).Inc()
}

func (m *Metrics) ObserveHTTPDuration(method, path string, duration time.Duration) {
	m.httpRequestDuration.WithLabelValues(method, path).Observe(duration.Seconds())
}

// Database Metrics
func (m *Metrics) RecordDBQuery(operation string, duration time.Duration, success bool) {
	status := "success"
	if !success {
		status = "error"
	}
	
	m.dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
	m.dbQueriesTotal.WithLabelValues(operation, status).Inc()
}

func (m *Metrics) IncrementDBConnections() {
	m.dbConnectionsActive.Inc()
}

func (m *Metrics) DecrementDBConnections() {
	m.dbConnectionsActive.Dec()
}

// Business Metrics
func (m *Metrics) IncrementTrackPlayed(trackID string) {
	m.tracksPlayedTotal.WithLabelValues(trackID).Inc()
}

func (m *Metrics) IncrementUserAction(action string) {
	m.userActionsTotal.WithLabelValues(action).Inc()
}

func (m *Metrics) RecordCacheHit(cacheType string, hit bool) {
	result := "miss"
	if hit {
		result = "hit"
	}
	m.cacheHitsTotal.WithLabelValues(cacheType, result).Inc()
}

// System Metrics
func (m *Metrics) RecordMemoryUsage(bytes uint64) {
	m.memoryUsageBytes.Set(float64(bytes))
}

func (m *Metrics) RecordGoroutines(count int) {
	m.goroutinesActive.Set(float64(count))
}
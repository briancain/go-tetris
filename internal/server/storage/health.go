package storage

// HealthChecker provides health check capabilities for storage backends
type HealthChecker interface {
	HealthCheck() error
}

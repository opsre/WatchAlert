package provider

import (
	"fmt"
	"watchAlert/internal/models"
)

// NewTraceClient creates the trace client for a supported trace datasource.
func NewTraceClient(datasource models.AlertDataSource) (TracesFactoryProvider, error) {
	switch datasource.Type {
	case JaegerDsProviderName:
		return NewJaegerClient(datasource)
	case SkyWalkingDsProviderName:
		return NewSkyWalkingClient(datasource)
	default:
		return nil, fmt.Errorf("unsupported trace datasource type: %s", datasource.Type)
	}
}

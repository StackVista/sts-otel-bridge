package collector

import (
	"context"

	colmetricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	"google.golang.org/grpc"
)

var _ colmetricspb.MetricsServiceServer = (*MetricsCollector)(nil)
var _ OpenTelemetryCollector[colmetricspb.ExportMetricsServiceRequest, colmetricspb.ExportMetricsServiceResponse] = (*MetricsCollector)(nil)

type MetricsCollector struct {
	colmetricspb.UnimplementedMetricsServiceServer
}

func (m *MetricsCollector) Export(ctx context.Context, req *colmetricspb.ExportMetricsServiceRequest) (*colmetricspb.ExportMetricsServiceResponse, error) {
	// for _, rs := range req.GetResourceMetrics() {
	// }
	return &colmetricspb.ExportMetricsServiceResponse{}, nil
}

func (m *MetricsCollector) Register(s *grpc.Server) {
	colmetricspb.RegisterMetricsServiceServer(s, m)
}

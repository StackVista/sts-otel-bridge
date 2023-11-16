package collector

import (
	"context"

	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	"google.golang.org/grpc"
)

var _ collogspb.LogsServiceServer = (*LogsCollector)(nil)
var _ OpenTelemetryCollector[collogspb.ExportLogsServiceRequest, collogspb.ExportLogsServiceResponse] = (*LogsCollector)(nil)

type LogsCollector struct {
	collogspb.UnimplementedLogsServiceServer
}

func (l *LogsCollector) Export(ctx context.Context, req *collogspb.ExportLogsServiceRequest) (*collogspb.ExportLogsServiceResponse, error) {
	// for _, rs := range req.GetResourceLogs() {
	// }
	return &collogspb.ExportLogsServiceResponse{}, nil
}

func (l *LogsCollector) Register(s *grpc.Server) {
	collogspb.RegisterLogsServiceServer(s, l)
}

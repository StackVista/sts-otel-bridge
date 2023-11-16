package collector

import (
	"context"

	"github.com/stackvista/sts-otel-bridge/pkg/grpc"
	collogspb "go.opentelemetry.io/proto/otlp/collector/logs/v1"
	colmetricspb "go.opentelemetry.io/proto/otlp/collector/metrics/v1"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
)

type OpenTelemetrySourceReqData interface {
	collogspb.ExportLogsServiceRequest | coltracepb.ExportTraceServiceRequest | colmetricspb.ExportMetricsServiceRequest
}

type OpenTelemetrySourceRespData interface {
	collogspb.ExportLogsServiceResponse | coltracepb.ExportTraceServiceResponse | colmetricspb.ExportMetricsServiceResponse
}

type OpenTelemetryCollector[Req OpenTelemetrySourceReqData, Resp OpenTelemetrySourceRespData] interface {
	grpc.Listener
	Export(ctx context.Context, req *Req) (*Resp, error)
}

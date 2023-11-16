package collector

import (
	"context"
	"encoding/binary"
	"encoding/json"

	"github.com/stackvista/sts-otel-bridge/internal/config"
	"github.com/stackvista/sts-otel-bridge/internal/logging"
	"github.com/stackvista/sts-otel-bridge/pkg/identifier"
	ststracepb "github.com/stackvista/sts-otel-bridge/proto/sts/trace"
	coltracepb "go.opentelemetry.io/proto/otlp/collector/trace/v1"
	v1 "go.opentelemetry.io/proto/otlp/common/v1"
	"google.golang.org/grpc"
)

var _ coltracepb.TraceServiceServer = (*TraceCollector)(nil)
var _ OpenTelemetryCollector[coltracepb.ExportTraceServiceRequest, coltracepb.ExportTraceServiceResponse] = (*TraceCollector)(nil)

type TraceID uint64

type TraceCollector struct {
	coltracepb.UnimplementedTraceServiceServer
	Out    chan []*ststracepb.APITrace
	Active bool
}

func NewTraceCollector(cfg config.Trace) *TraceCollector {
	return &TraceCollector{
		Out:    make(chan []*ststracepb.APITrace),
		Active: true,
	}
}

func (t *TraceCollector) Export(ctx context.Context, req *coltracepb.ExportTraceServiceRequest) (*coltracepb.ExportTraceServiceResponse, error) {
	logger := logging.LoggerFor(ctx, "trace-collector")
	var traces map[TraceID]*ststracepb.APITrace = map[TraceID]*ststracepb.APITrace{}

	for _, rs := range req.GetResourceSpans() {
		// TODO - remove this once we have a better way to debug
		js, err := json.Marshal(rs.GetResource())
		if err != nil {
			logger.Warn().Err(err).Msg("Failed to marshal resource")
		}
		println(string(js))

		logger.Debug().Str("resource-attributes", rs.GetResource().String()).Msg("Processing ResourceSpan")
		ident, err := identifier.Identify(rs.GetResource())
		if err != nil {
			js, err := json.Marshal(rs.GetResource())
			if err != nil {
				logger.Warn().Err(err).Msg("Failed to marshal resource")
			}
			logger.Warn().Err(err).Str("resource", string(js)).Msg("Failed to identify resource")
		}
		logger.Info().Str("resource", string(ident)).Msg("Identified resource")

		for _, ss := range rs.GetScopeSpans() {
			for _, s := range ss.GetSpans() {
				// Group the Spans by their TraceIDs for StS
				traceID := TraceID(binary.BigEndian.Uint64(s.TraceId))
				if _, ok := traces[traceID]; !ok {
					traces[traceID] = &ststracepb.APITrace{
						TraceID: uint64(traceID),
						Spans:   []*ststracepb.Span{},
					}
				}
				stsTrace := traces[traceID]

				if s.ParentSpanId == nil || len(s.ParentSpanId) == 0 {
					// Root span, this contains the trace start and end times
					stsTrace.StartTime = int64(s.StartTimeUnixNano)
					stsTrace.EndTime = int64(s.EndTimeUnixNano)
				}

				stsSpan := &ststracepb.Span{
					Name:     s.Name,
					TraceID:  binary.BigEndian.Uint64(s.TraceId),
					SpanID:   binary.BigEndian.Uint64(s.SpanId),
					Start:    int64(s.StartTimeUnixNano),
					Duration: int64(s.EndTimeUnixNano - s.StartTimeUnixNano),
					Meta:     convertAttributes(s.Attributes),
				}

				if s.ParentSpanId != nil && len(s.ParentSpanId) > 0 {
					stsSpan.ParentID = binary.BigEndian.Uint64(s.ParentSpanId)
				}

				stsTrace.Spans = append(stsTrace.Spans, stsSpan)
			}
		}
	}

	tt := []*ststracepb.APITrace{}
	for _, stsTrace := range traces {
		logger.Debug().Uint64("trace_id", stsTrace.TraceID).Msg("Sending trace to StS")
		tt = append(tt, stsTrace)
	}

	t.Out <- tt

	return &coltracepb.ExportTraceServiceResponse{}, nil
}

func convertAttributes(attrs []*v1.KeyValue) map[string]string {
	m := map[string]string{}
	for _, attr := range attrs {
		m[attr.GetKey()] = attr.GetValue().GetStringValue()
	}

	return m
}

func (t *TraceCollector) Register(s *grpc.Server) {
	coltracepb.RegisterTraceServiceServer(s, t)
}

func (t *TraceCollector) Stop(ctx context.Context) error {
	t.Active = false
	close(t.Out)

	return nil
}

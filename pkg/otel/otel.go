package otel

import (
	"encoding/json"

	ststracepb "github.com/stackvista/sts-otel-bridge/proto/sts/trace"
)

type OpenTelemetryData interface {
	*ststracepb.APITrace
}

func PrintOTelData[D OpenTelemetryData](data D) {
	switch dd := any(data).(type) {
	case *ststracepb.APITrace:
		js, _ := json.Marshal(dd)
		println(js)
	}
}

package stackstate

type ReceiverMessage struct {
	ApiKey              string          `json:"apiKey"`
	CollectionTimestamp int64           `json:"collection_timestamp"`
	InternalHostname    string          `json:"internalHostname"`
	Events              Events          `json:"events"`
	Metrics             []Metrics       `json:"metrics"`
	ServiceChecks       []ServiceChecks `json:"service_checks"`
	Health              []Health        `json:"health"`
	Topologies          []Topology      `json:"topologies"`
}

type Events struct{}

type Metrics struct{}

type ServiceChecks struct{}

type Health struct{}

type Topology struct {
	StartSnapshot bool        `json:"start_snapshot"`
	StopSnapshot  bool        `json:"stop_snapshot"`
	Instance      Instance    `json:"instance"`
	Components    []Component `json:"components"`
	Relations     []Relation  `json:"relations"`
}

type Instance struct {
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Component struct {
	ExternalId       string                 `json:"externalId"`
	Type             ComponentType          `json:"type"`
	Data             map[string]interface{} `json:"data"`
	SourceProperties map[string]interface{} `json:"sourceProperties"`
}

type ComponentType struct {
	Name string `json:"name"`
}

type Relation struct {
	ExternalId string                 `json:"externalId"`
	Type       RelationType           `json:"type"`
	SourceId   string                 `json:"sourceId"`
	TargetId   string                 `json:"targetId"`
	Data       map[string]interface{} `json:"data"`
}

type RelationType struct {
	Name string `json:"name"`
}

func NewTopologyMessage() *ReceiverMessage {
	return &ReceiverMessage{
		Events:        Events{},
		Metrics:       []Metrics{},
		ServiceChecks: []ServiceChecks{},
		Health:        []Health{},
		Topologies:    []Topology{},
	}
}

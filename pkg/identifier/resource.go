package identifier

import (
	"fmt"

	v1 "go.opentelemetry.io/proto/otlp/resource/v1"
)

type Identifier string

var Mapping = map[string][]string{
	"urn:kubernetes:/{k8s.cluster.name}:{k8s.namespace.name}:pod/{k8s.pod.name}": []string{"k8s.cluster.name", "k8s.namespace.name", "k8s.pod.name"},
}

func Identify(resource *v1.Resource) (Identifier, error) {
	kvs := toMap(resource)
	for pattern, keys := range Mapping {
		for _, key := range keys {
			if _, ok := kvs[key]; !ok {
				// not this pattern
				break
			}
		}

		// all keys matched
		return CreateIdentifier(pattern, kvs)
	}

	return "", fmt.Errorf("No identifier matched resource")
}

func toMap(resource *v1.Resource) map[string]string {
	m := make(map[string]string)
	for _, kv := range resource.Attributes {
		m[kv.Key] = kv.Value.String()
	}

	return m
}

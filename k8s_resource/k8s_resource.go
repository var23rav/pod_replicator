package resource

import (
	"fmt"
	"strings"

	prclient "github.com/var23rav/pod_replicator/client"
)

type Resource interface {
	Name() string
	Migrate(bool) error
}

type ResourceType string

func (resourceType ResourceType) String() string {
	return strings.ToLower(string(resourceType))
}

func (resourceType ResourceType) Build(replicationClient *prclient.ReplicationClient) (Resource, error) {
	switch resourceType.String() {
	case "deployment":
		return Deployment{replicationClient}, nil
	default:
		return nil, fmt.Errorf("Unknown or Unsupported ResourceType %s", resourceType)
	}
}

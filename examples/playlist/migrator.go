package playlist

import (
	"context"
	"fmt"

	"github.com/grafana/kindsys"
)

type MigrationLookupHooks interface {
	GetUIDFromID(ctx context.Context, id int64) (string, error)
	GetTitleAndIDFromUID(ctx context.Context, uid string) (int64, string, error)
}

func newMigrator(hooks MigrationLookupHooks) kindsys.ResourceMigrator {
	return func(ctx context.Context, obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
		return nil, fmt.Errorf("TODO... actually migrate")
	}
}

var _ MigrationLookupHooks = &dummyLookupHooks{}

type dummyLookupHooks struct{}

func (k *dummyLookupHooks) GetUIDFromID(ctx context.Context, id int64) (string, error) {
	switch id {
	case 111:
		return "AAA", nil
	case 222:
		return "BBB", nil
	}
	return "", fmt.Errorf("unknown internal id")
}

func (k *dummyLookupHooks) GetTitleAndIDFromUID(ctx context.Context, uid string) (int64, string, error) {
	switch uid {
	case "AAA":
		return 111, "Title for ID(111)", nil
	case "BBB":
		return 222, "Title for ID(222)", nil
	}
	return 0, "", fmt.Errorf("unknown uid")
}

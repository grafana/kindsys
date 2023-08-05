package playlist

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/examples/playlist/v0x"
	"github.com/grafana/kindsys/examples/playlist/v1x"
)

// The real version will need access to a database for this to work
type MigrationLookupHooks interface {
	GetUIDFromID(ctx context.Context, id int64) (string, error)
	GetTitleAndIDFromUID(ctx context.Context, uid string) (int64, string, error)
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

func newMigrator(hooks MigrationLookupHooks) kindsys.ResourceMigrator {
	return func(ctx context.Context, obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
		srcVersion := obj.StaticMetadata().Version
		if srcVersion == targetVersion {
			return obj, nil
		}

		srcMajor := -1
		srcMinor := -1
		targetMajor := -1
		targetMinor := -1
		n, err := fmt.Sscanf(srcVersion, "v%d-%d", &srcMajor, &srcMinor)
		if err != nil || n != 2 {
			return nil, fmt.Errorf("error reading source version")
		}
		n, err = fmt.Sscanf(targetVersion, "v%d-%d", &targetMajor, &targetMinor)
		if err != nil || n != 2 {
			return nil, fmt.Errorf("error reading target version")
		}

		data, err := json.Marshal(obj.SpecObject())
		if err != nil {
			return nil, err
		}

		switch srcMajor {
		case 0:
			spec := &v0x.Spec{}
			err = json.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}

			switch targetMajor {
			case 0:

			case 1:
				// from 0 to 1
				targetSpec := v1x.Spec{}
				// TODO!!

				return &ResourceV1{
					StaticMeta: obj.StaticMetadata(),
					CommonMeta: obj.CommonMetadata(),
					Spec:       targetSpec,
				}, nil

			default:
				return nil, fmt.Errorf("invalid target")
			}

			return nil, fmt.Errorf("TODO... actually migrate v0")

		case 1:
			spec := &v1x.Spec{}
			err = json.Unmarshal(data, spec)
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("TODO... actually migrate v0")
		}

		return nil, fmt.Errorf("invalid version")
	}
}

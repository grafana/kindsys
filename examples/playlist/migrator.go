package playlist

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

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

// EEEEP... this is awful, but at least it works
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

			static := obj.StaticMetadata()
			static.Version = targetVersion
			if static.Name == "" {
				static.Name = spec.Uid
			}
			switch targetMajor {
			case 0:
				switch targetMinor {
				case 0:
					return nil, fmt.Errorf("TODO... migrate down")
				case 1:
					// Migrate minor up (id to uid if possible)
					for i, item := range spec.Items {
						if item.Type == v0x.ItemTypeDashboardById {
							id, err := strconv.ParseInt(item.Value, 10, 64)
							if err != nil {
								return nil, err
							}
							uid, err := hooks.GetUIDFromID(ctx, id)
							if err == nil {
								spec.Items[i] = v0x.Item{
									Type:  v0x.ItemTypeDashboardByUid,
									Value: uid,
								}
							}
						}
					}
				}
				return &ResourceV0{
					StaticMeta: static,
					CommonMeta: obj.CommonMetadata(),
					Spec:       *spec,
				}, nil

			case 1:
				// from 0 to 1
				targetSpec := v1x.Spec{
					Interval: spec.Interval,
					Name:     spec.Name,
					Items:    make([]v1x.Item, len(spec.Items)),
					Xxx:      "just here for the change detection version bypass",
				}
				for i, item := range spec.Items {
					dst, err := migrateItemV0ToV1(ctx, item, hooks)
					if err != nil {
						return nil, err
					}
					targetSpec.Items[i] = dst
				}
				return &ResourceV1{
					StaticMeta: static,
					CommonMeta: obj.CommonMetadata(),
					Spec:       targetSpec,
					//CustomMeta: obj.CustomMetadata(),
				}, nil

			default:
				return nil, fmt.Errorf("invalid target")
			}

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

func migrateItemV0ToV1(ctx context.Context, src v0x.Item, hooks MigrationLookupHooks) (v1x.Item, error) {
	dst := v1x.Item{
		Type:  v1x.ItemTypeDashboardByUid,
		Value: src.Value,
	}
	switch src.Type {
	case v0x.ItemTypeDashboardById:
		id, err := strconv.ParseInt(src.Value, 10, 64)
		if err != nil {
			return dst, err
		}
		uid, err := hooks.GetUIDFromID(ctx, id)
		if err != nil {
			return dst, err
		}
		dst.Value = uid
	case v0x.ItemTypeDashboardByTag:
		dst.Type = v1x.ItemTypeDashboardByTag
	case v0x.ItemTypeDashboardByUid:
		dst.Type = v1x.ItemTypeDashboardByUid
	default:
		return dst, fmt.Errorf("invalid src type")
	}
	return dst, nil
}

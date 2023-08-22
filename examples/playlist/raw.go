package playlist

import (
	"context"
	"fmt"
	"io"

	"github.com/grafana/kindsys"
	"github.com/grafana/kindsys/examples/playlist/v0x"
	"github.com/grafana/kindsys/examples/playlist/v1x"
	jsoniter "github.com/json-iterator/go"
	"k8s.io/kube-openapi/pkg/common"
)

var _ kindsys.ResourceKind = &rawPlaylistKind{}

// This implements a playlist directly in golang
type rawPlaylistKind struct {
	migrator kindsys.ResourceMigrator
}

func NewRawPlaylistKind(hooks MigrationLookupHooks) kindsys.ResourceKind {
	return &rawPlaylistKind{
		migrator: newMigrator(hooks),
	}
}

func (k *rawPlaylistKind) GetMachineNames() kindsys.MachineNames {
	return kindsys.MachineNames{
		Plural:   "playlists",
		Singular: "playlist",
	}
}

func (k *rawPlaylistKind) GetKindInfo() kindsys.KindInfo {
	return kindsys.KindInfo{
		Group:       "playlists.ext.grafana.com",
		Kind:        "Playlist",
		Description: "A set of dashboards that will be displayed in a loop (dummy for testing)",
		Maturity:    kindsys.MaturityMerged,
	}
}

func (k *rawPlaylistKind) CurrentVersion() string {
	return "v1-0" // would be nice pick a middle one
}

func (k *rawPlaylistKind) GetVersions() []kindsys.VersionInfo {
	return []kindsys.VersionInfo{
		{
			Version:         "v0-0",
			SoftwareVersion: "v6.0", // when playlists were introduced
		},
		{
			Version:         "v0-1",
			SoftwareVersion: "v9.1", // when we added uid support
			Changelog: []string{
				"adding the dashboard_by_uid type",
				"deprecating the dashboard_by_id type",
				"deprecating the PlaylistItem.title property (now optional and unused)",
				"TODO: verify that k8s name and spec.uid match",
			},
		},
		{
			Version:         "v1-0",
			SoftwareVersion: "v10.5", // when we remove internal id support
			Changelog: []string{
				"removed the dashboard_by_id enumeration type",
				"removed the PlaylistItem.title property",
				"remove the spec.uid property",
				"TODO! added xxx so thema will detect a breaking version change",
			},
		},
	}
}

// NOTE: this files are not used to do validation, but can be used generically to describe the kind
func (k *rawPlaylistKind) GetOpenAPIDefinition(version string, ref common.ReferenceCallback) (common.OpenAPIDefinition, error) {
	api := common.OpenAPIDefinition{}
	s, err := packageFS.ReadFile("schemas/" + version + ".json")
	if err != nil {
		return api, fmt.Errorf("unknown schema version")
	}
	return kindsys.LoadOpenAPIDefinition(s)
}

type ResourceV0 = kindsys.GenericResource[v0x.Spec, kindsys.SimpleCustomMetadata, any]
type ResourceV1 = kindsys.GenericResource[v1x.Spec, kindsys.SimpleCustomMetadata, any]

func (k *rawPlaylistKind) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	major := -1
	minor := -1

	var rV0 *ResourceV0
	var rV1 *ResourceV1
	var obj kindsys.Resource

	var static kindsys.StaticMetadata
	var common kindsys.CommonMetadata

	err := kindsys.ReadResourceJSON(reader, kindsys.JSONResourceBuilder{
		SetGroupVersionKind: func(group, version, kind string) error {
			if group != k.GetKindInfo().Group {
				return fmt.Errorf("invalid group")
			}
			if kind != k.GetKindInfo().Kind {
				return fmt.Errorf("invalid kind")
			}
			n, err := fmt.Sscanf(version,
				"v%d-%d", &major, &minor)
			if err != nil || n != 2 {
				return fmt.Errorf("unable to read version")
			}
			return nil
		},
		SetMetadata: func(s kindsys.StaticMetadata, c kindsys.CommonMetadata) {
			static = s
			common = c
		},
		ReadSpec: func(iter *jsoniter.Iterator) error {
			if major == 1 {
				rV1 = &ResourceV1{StaticMeta: static, CommonMeta: common}
				iter.ReadVal(&rV1.Spec)
				if iter.Error == nil && strict {
					err := rV1.Spec.Validate(major, minor)
					if err != nil {
						return err
					}
				}
				obj = rV1
			} else if major == 0 {
				rV0 = &ResourceV0{StaticMeta: static, CommonMeta: common}
				iter.ReadVal(&rV0.Spec)
				if iter.Error == nil && strict {
					err := rV0.Spec.Validate(major, minor)
					if err != nil {
						return err
					}

					if minor > 0 && rV0.Spec.Uid != rV0.StaticMeta.Name {
						return fmt.Errorf("the spec.uid must match metadata.name")
					}
				}
				obj = rV0
			} else {
				return fmt.Errorf("unknown major version")
			}
			return nil
		},
		SetAnnotation: func(key, val string) {
			fmt.Printf("??? unknown")
		},
	})
	if err != nil {
		return nil, err
	}
	if obj == nil {
		return nil, fmt.Errorf("missing spec")
	}
	return obj, err
}

func (k *rawPlaylistKind) Migrate(ctx context.Context, obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	return k.migrator(ctx, obj, targetVersion)
}

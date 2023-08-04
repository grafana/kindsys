package santhoshsys

import (
	"fmt"
	"io"
	"io/fs"
	"strings"

	"github.com/grafana/kindsys"
	jsoniter "github.com/json-iterator/go"
)

var _ kindsys.ResourceKind = &resourceKindFromManifest{}

type resourceKindFromManifest struct {
	kindFromManifest // the base properties

	//
	names kindsys.MachineNames
}

// Load a jsonschema based kind from a file system
// the file system will have a manifest that exists
func NewResourceKind(sfs fs.FS) (kindsys.ResourceKind, error) {
	m := &resourceKindFromManifest{}
	info, err := m.init(sfs)
	if err != nil {
		return m, err
	}
	if info.ComposableType != "" || len(info.ComposableSlots) > 0 {
		return nil, fmt.Errorf("invalid info in the manifest (should not have composable types)")
	}

	if info.MachineNames != nil {
		m.names = *info.MachineNames
	}
	if m.names.Singular == "" {
		m.names.Singular = strings.ToLower(info.Kind)
	}
	if m.names.Plural == "" {
		m.names.Plural = m.names.Singular + "s"
	}
	return m, nil
}

func (m *resourceKindFromManifest) GetMachineNames() kindsys.MachineNames {
	return m.names
}

func (m *resourceKindFromManifest) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	obj := &kindsys.UnstructuredResource{}
	err := kindsys.ReadResourceJSON(reader, kindsys.JSONResourceBuilder{
		SetStaticMetadata: func(v kindsys.StaticMetadata) { obj.StaticMeta = v },
		SetCommonMetadata: func(v kindsys.CommonMetadata) { obj.CommonMeta = v },
		ReadSpec: func(iter *jsoniter.Iterator) error {
			obj.Spec = make(map[string]any)
			iter.ReadVal(&obj.Spec)
			return iter.Error
		},
		SetAnnotation: func(key, val string) {
			fmt.Printf("??? unknown")
		},
		ReadStatus: func(iter *jsoniter.Iterator) error {
			obj.Status = make(map[string]any)
			iter.ReadVal(&obj.Status)
			return iter.Error
		},
		ReadSub: func(name string, iter *jsoniter.Iterator) error {
			return fmt.Errorf("unsupported sub resource")
		},
	})
	if err != nil {
		return obj, err
	}

	if strict {
		meta := obj.StaticMetadata()
		if meta.Group != m.info.Group {
			return obj, fmt.Errorf("wrong group")
		}
		if meta.Kind != m.info.Kind {
			return obj, fmt.Errorf("wrong kind")
		}

		schema, ok := m.parsed[meta.Version]
		if !ok || schema == nil {
			return obj, fmt.Errorf("unknown version")
		}

		// TODO!!! schema is right now just on the spec!!!
		doc := obj.SpecObject()
		// TODO: need to make sure the doc+resource are ones that we can parse ()
		err = schema.ValidateInterface(doc)
	}
	return obj, err
}

func (m *resourceKindFromManifest) Migrate(obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	return nil, fmt.Errorf("TODO")
}

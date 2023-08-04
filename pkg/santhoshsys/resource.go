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

func (k *resourceKindFromManifest) GetMachineNames() kindsys.MachineNames {
	return k.names
}

func (k *resourceKindFromManifest) Read(reader io.Reader, strict bool) (kindsys.Resource, error) {
	obj := &kindsys.UnstructuredResource{}
	err := kindsys.ReadResourceJSON(reader, kindsys.JSONResourceBuilder{
		SetGroupVersionKind: func(group, version, kind string) error {
			if group != k.info.Group {
				return fmt.Errorf("invalid group")
			}
			if kind != k.info.Kind {
				return fmt.Errorf("invalid kind")
			}
			return nil
		},
		SetMetadata: func(s kindsys.StaticMetadata, c kindsys.CommonMetadata) {
			obj.StaticMeta = s
			obj.CommonMeta = c
		},
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
		schema, ok := k.parsed[obj.StaticMeta.Version]
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

func (k *resourceKindFromManifest) Migrate(obj kindsys.Resource, targetVersion string) (kindsys.Resource, error) {
	return nil, fmt.Errorf("TODO")
}

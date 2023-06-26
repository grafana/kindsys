package resource

import (
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// A Resource is a single instance of a Grafana kind.
//
// The relationship between Resource and kindsys.Kind is similar to the
// relationship between objects and classes in conventional object oriented
// design:
//
// - Objects are instantiated from classes. The name of the class is the type of object.
// - Resources are instantiated from kinds. The name of the kind is the type of resource.
//
// Resource provides helper methods for using the associated Metadata in accordance
// with the patterns specified by Grafana's kindsys. TODO link to docs once written
//
// Resource represents the shared parts of Grafana resources, common regardless
// of the underlying kind. It is generic over its Spec and Status fields. It is
// expected that the generic parameters for Resource are structs, generated from
// a CUE kind definition.
type Resource[Spec, Status any] struct {
	APIVersion string `json:"apiVersion"`
	// Kind is the name of the Grafana kind this type represents.
	//
	// NOTE This field name may be somewhat confusing; KindName would be better. But this
	// must be named Kind to match the field name in the Kubernetes API types.
	Kind string `json:"kind"`

	Metadata Metadata `json:"metadata"`
	Spec     *Spec    `json:"spec,omitempty"`
	Status   *Status  `json:"status,omitempty"`
}

// Metadata is standard k8s object metadata, but with added helper functions
// that assist with using it in accordance with the schema formats and other
// patterns expected by Grafana's kind system.
//
// Metadata.Annotations is not guarded by a lock in upstream k8s. As a result,
// the helper functions provided are NOT safe for use from multiple goroutines.
//
// TODO link to schema format docs once written
type Metadata v1.ObjectMeta

// Annotation keys
const annoKeyCreatedBy = "grafana.com/createdBy"
const annoKeyUpdatedTimestamp = "grafana.com/updatedTimestamp"
const annoKeyUpdatedBy = "grafana.com/updatedBy"

// The folder identifier
const annoKeyFolder = "grafana.com/folder"
const annoKeySlug = "grafana.com/slug"

// Identify where values came from
const annoKeyOriginName = "grafana.com/origin/name"
const annoKeyOriginPath = "grafana.com/origin/path"
const annoKeyOriginKey = "grafana.com/origin/key"
const annoKeyOriginTime = "grafana.com/origin/time"

// GetUpdatedTimestamp retrieves the timestamp when this Grafana resource was
// last updated from the kindsys-specified key in the standard Kubernetes
// annotations map.
func (m *Metadata) GetUpdatedTimestamp() *time.Time {
	v, ok := m.Annotations[annoKeyUpdatedTimestamp]
	if ok {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			return &t
		}
	}
	return nil
}

// SetUpdatedTimestamp sets the timestamp for when this Grafana resource was
// last updated at the expected key in Kubernetes annotations metadata. It is serialized
// to a string using the RFC3339 format.
func (m *Metadata) SetUpdatedTimestamp(v *time.Time) {
	if v == nil {
		delete(m.Annotations, annoKeyUpdatedTimestamp)
	} else {
		m.Annotations[annoKeyUpdatedTimestamp] = v.Format(time.RFC3339)
	}
}

// GetCreatedBy retrieves the string identifying the creator of this Grafana
// resource from Kubernetes annotations metadata.
func (m *Metadata) GetCreatedBy() string {
	return m.Annotations[annoKeyCreatedBy]
}

// SetCreatedBy sets the string identifying the creator of this Grafana
// resource at the expected key in Kubernetes annotations metadata.
func (m *Metadata) SetCreatedBy(user string) {
	m.Annotations[annoKeyCreatedBy] = user // user GRN
}

// GetUpdatedBy retrieves the string identifying the last updater of this
// Grafana resource from Kubernetes annotations metadata.
func (m *Metadata) GetUpdatedBy() string {
	return m.Annotations[annoKeyUpdatedBy]
}

// SetUpdatedBy sets the string identifying the last updater of this Grafana
// resource at the expected key in the annotations map.
func (m *Metadata) SetUpdatedBy(user string) {
	m.Annotations[annoKeyUpdatedBy] = user // user GRN
}

// GetFolder retrieves the folder identifier from Kubernetes annotations metadata.
func (m *Metadata) GetFolder() string {
	return m.Annotations[annoKeyFolder]
}

// SetFolder sets the folder identifier in Kubernetes annotations metadata.
func (m *Metadata) SetFolder(uid string) {
	m.Annotations[annoKeyFolder] = uid
}

// GetSlug retrieves the slug from Kubernetes annotations metadata.
func (m *Metadata) GetSlug() string {
	return m.Annotations[annoKeySlug]
}

// SetSlug sets the slug in Kubernetes annotations metadata.
func (m *Metadata) SetSlug(v string) {
	m.Annotations[annoKeySlug] = v
}

// OriginInfo is saved in annotations. When the canonical definition of some resource
// lives outside of a Grafana instance's storage, the values in this object identify
// where it came from.
//
// This object can model the same data as our existing provisioning table
// or a more general git sync.
//
// TODO generalize this further, paths may be a poor assumption
type OriginInfo struct {
	// Name of the origin/provisioning source
	Name string `json:"name,omitempty"`

	// The path within the named origin above (external_id in the existing dashboard provisioing)
	Path string `json:"path,omitempty"`

	// Verification/identification key (check_sum in existing dashboard provisioning)
	Key string `json:"key,omitempty"`

	// Origin modification timestamp when the resource was saved
	// This will be before the resource updated time
	Timestamp *time.Time `json:"time,omitempty"`
}

// SetOriginInfo sets the origin info for a resource in Kubernetes annotations metadata.
//
// Calling always clears any origin info already stored in the annotations metadata.
//
// Only non-empty fields on the provided OriginInfo are stored in metadata. If
// info is nil or info.Name is empty, nothing is stored in metadata, and any
// existing origin info is removed.
func (m *Metadata) SetOriginInfo(info *OriginInfo) {
	delete(m.Annotations, annoKeyOriginName)
	delete(m.Annotations, annoKeyOriginPath)
	delete(m.Annotations, annoKeyOriginKey)
	delete(m.Annotations, annoKeyOriginTime)
	if info != nil || info.Name != "" {
		m.Annotations[annoKeyOriginName] = info.Name
		if info.Path != "" {
			m.Annotations[annoKeyOriginPath] = info.Path
		}
		if info.Key != "" {
			m.Annotations[annoKeyOriginKey] = info.Key
		}
		if info.Timestamp != nil {
			m.Annotations[annoKeyOriginTime] = info.Timestamp.Format(time.RFC3339)
		}
	}
}

// GetOriginInfo returns the origin info stored in Kubernetes annotations
// metadata.
func (m *Metadata) GetOriginInfo() *OriginInfo {
	v, ok := m.Annotations[annoKeyOriginName]
	if !ok {
		return nil
	}
	info := &OriginInfo{
		Name: v,
		Path: m.Annotations[annoKeyOriginPath],
		Key:  m.Annotations[annoKeyOriginKey],
	}
	v, ok = m.Annotations[annoKeyOriginTime]
	if ok {
		t, err := time.Parse(time.RFC3339, v)
		if err != nil {
			info.Timestamp = &t
		}
	}
	return info
}

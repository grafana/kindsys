package kindsys

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// WireFormat enumerates values for possible message wire formats.
// Constants with these values are in this package with a `WireFormat` prefix.
type WireFormat int

const (
	// WireFormatUnknown is an unknown message wire format.
	WireFormatUnknown WireFormat = iota
	// WireFormatJSON is a JSON message wire format, which should be handle-able by the `json` package.
	// (messages which _contain_ JSON, but are not parsable by the go json package should not be
	// considered to be of the JSON wire format).
	WireFormatJSON
)

// UnmarshalConfig is the config used for unmarshaling Resources.
// It consists of fields that are descriptive of the underlying content, based on knowledge the caller has.
type UnmarshalConfig struct {
	// WireFormat is the wire format of the provided payload
	WireFormat WireFormat
	// VersionHint is what the client thinks the version is (if non-empty)
	VersionHint string
}

// A Resource is a single instance of a Grafana [Kind], either [Core] or [Custom].
//
// A Resource is broadly composed of metadata, spec, and additional "subresources."
// Metadata is split into three sub-components:
// - KubeMetadata is the metadata provided in the "metadata" component of a kubernetes resource
// - GrafanaMetadata is additional, standard Grafana Resource metadata
// - KindMetadata is metadata specific to the Kind
// Spec is the Resource's main payload, what could classicly be considered the "body" of the Resource.
// Subresources are additional components of the Resource which are not considered part of either the metadata or
// the spec ("body"). For properly-defined Grafana Resources, this will always include the "status" subresource,
// and may include others on a per-Kind basis.
//
// The relationship between Resource and [Kind] is similar to the
// relationship between objects and classes in conventional object oriented
// design:
//
// - Objects are instantiated from classes. The name of the class is the type of object.
// - Resources are instantiated from kinds. The name of the kind is the type of resource.
//
// Resource is an interface, rather than a concrete struct, for two reasons:
//
// - Some use cases need to operate generically over resources of any kind.
// - Go generics do not allow the ergonomic expression of certain needed constraints.
//
// The [Core] and [Custom] interfaces are intended for the generic operation
// use case, fulfilling [Resource] using [UnstructuredResource].
//
// For known, specific kinds, it is usually possible to rely on code generation
// to produce a struct that implements [Resource] for each kind. Such a struct
// can be used as the generic type parameter to create a [TypedCore] or [TypedCustom]
type Resource interface {
	// KubeMetadata returns the kubernetes standard resource metadata
	KubeMetadata() KubeMetadata

	// GrafanaMetadata returns the grafana metadata for the resource
	GrafanaMetadata() GrafanaMetadata

	// CustomMetadata returns metadata unique to this Resource's kind, as opposed to Common and Static metadata,
	// which are the same across all kinds. An object may have no kind-specific CustomMetadata.
	// CustomMetadata can only be read from this interface, for use with resource.Client implementations,
	// those who wish to set CustomMetadata should use the interface's underlying type.
	KindMetadata() CustomMetadata

	// SetCommonMetadata overwrites the CommonMetadata of the object.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with CommonMetadata() and set specific values.
	SetKubeMetadata(metadata KubeMetadata)

	SetGrafanaMetadata(metadata GrafanaMetadata)

	// StaticMetadata returns the Resource's StaticMetadata
	StaticMetadata() StaticMetadata

	// SetStaticMetadata overwrites the Resource's StaticMetadata with the provided StaticMetadata.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with StaticMetadata() and set specific values.
	// Note that StaticMetadata is only mutable in an object create context.
	SetStaticMetadata(metadata StaticMetadata)

	// SpecObject returns the actual "schema" object, which holds the main body of data
	SpecObject() any

	// Subresources returns a map of subresource name(s) to the object value for that subresource.
	// Spec is not considered a subresource, and should only be returned by SpecObject
	// No guarantees are made that mutations to objects in the map will affect the underlying resource.
	Subresources() map[string]any

	// Copy returns a full copy of the Resource with all its data
	Copy() Resource
}

// CustomMetadata is an interface describing a kindsys.Resource's kind-specific metadata
type CustomMetadata interface {
	// MapFields converts the custom metadata's fields into a map of field key to value.
	// This is used so Clients don't need to engage in reflection for marshaling metadata,
	// as various implementations may not store kind-specific metadata the same way.
	MapFields() map[string]any
}

// StaticMetadata consists of all non-mutable metadata for an object.
// It is set in the initial Create call for an Resource, then will always remain the same.
type StaticMetadata struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// Identifier creates an Identifier struct from the StaticMetadata
func (s StaticMetadata) Identifier() Identifier {
	return Identifier{
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

// FullIdentifier returns a FullIdentifier struct from the StaticMetadata.
// Plural cannot be inferred so is left empty.
func (s StaticMetadata) FullIdentifier() FullIdentifier {
	return FullIdentifier{
		Group:     s.Group,
		Version:   s.Version,
		Kind:      s.Kind,
		Namespace: s.Namespace,
		Name:      s.Name,
	}
}

type Identifier struct {
	Namespace string
	Name      string
}

// FullIdentifier is a globally-unique identifier, consisting of Schema identity information
// (Group, Version, Kind, Plural) and within-schema identity information (Namespace, Name)
type FullIdentifier struct {
	Namespace string
	Name      string
	Group     string
	Version   string
	Kind      string
	Plural    string
}

// ListResource represents a List of Resource-implementing objects with list metadata.
// The simplest way to use it is to use the implementation returned by a Client's List call.
type ListResource interface {
	ListMetadata() ListMetadata
	SetListMetadata(ListMetadata)
	ListItems() []Resource
	SetItems([]Resource)
}

// ListMetadata is metadata for a list of objects. This is typically only used in responses from the storage layer.
type ListMetadata struct {
	ResourceVersion string `json:"resourceVersion"`

	Continue string `json:"continue"`

	RemainingItemCount *int64 `json:"remainingItemCount"`

	// ExtraFields stores implementation-specific metadata.
	// Not all Client implementations are required to honor all ExtraFields keys.
	// Generally, this field should be shied away from unless you know the specific
	// Client implementation you're working with and wish to track or mutate extra information.
	ExtraFields map[string]any `json:"extraFields"`
}

// KubeMetadata is kubernetes object metadata. It does not directly import kubernetes go structs to avoid requiring kubernetes dependencies in projects
// which use kindsys.
// TODO: does this matter? Could we just import kubernetes' struct and attach methods we need to it?
type KubeMetadata struct {
	Name      string
	Namespace string
	// UID is the unique ID of the object. This can be used to uniquely identify objects,
	// but is not guaranteed to be usable for lookups.
	UID string `json:"uid"`
	// ResourceVersion is a version string used to identify any and all changes to the object.
	// Any time the object changes in storage, the ResourceVersion will be changed.
	// This can be used to block updates if a change has been made to the object between when the object was
	// retrieved, and when the update was applied.
	ResourceVersion string `json:"resourceVersion"`
	// Labels are string key/value pairs attached to the object. They can be used for filtering,
	// or as additional metadata.
	Labels map[string]string `json:"labels"`
	// Annotations
	Annotations map[string]string `json:"annotations"`
	// CreationTimestamp indicates when the resource has been created.
	CreationTimestamp time.Time `json:"creationTimestamp"`
	// DeletionTimestamp indicates that the resource is pending deletion as of the provided time if non-nil.
	// Depending on implementation, this field may always be nil, or it may be a "tombstone" indicator.
	// It may also indicate that the system is waiting on some task to finish before the object is fully removed.
	DeletionTimestamp *time.Time `json:"deletionTimestamp,omitempty"`
	// Finalizers are a list of identifiers of interested parties for delete events for this resource.
	// Once a resource with finalizers has been deleted, the object should remain in the store,
	// DeletionTimestamp is set to the time of the "delete," and the resource will continue to exist
	// until the finalizers list is cleared.
	Finalizers []string `json:"finalizers"`
}

func (k KubeMetadata) Copy() KubeMetadata {
	n := KubeMetadata{
		Name:              k.Name,
		Namespace:         k.Namespace,
		UID:               k.UID,
		ResourceVersion:   k.ResourceVersion,
		CreationTimestamp: k.CreationTimestamp,
	}
	if k.DeletionTimestamp != nil {
		*n.DeletionTimestamp = *(k.DeletionTimestamp)
	}
	copy(n.Finalizers, k.Finalizers)
	for key, val := range k.Annotations {
		n.Annotations[key] = val
	}
	for key, val := range k.Labels {
		k.Labels[key] = val
	}
	return n
}

// TODO
// On encoding to kubernetes, fields in GrafanaMetadata MUST be encoded into annotations with the name "grafana.com/X", where X is the JSON name of the field.
type GrafanaMetadata struct {
	// UpdateTimestamp is the timestamp of the last update to the resource.
	UpdateTimestamp time.Time `json:"updateTimestamp"`
	// CreatedBy is a string which indicates the user or process which created the resource.
	// Implementations may choose what this indicator should be.
	CreatedBy string `json:"createdBy"`
	// UpdatedBy is a string which indicates the user or process which last updated the resource.
	// Implementations may choose what this indicator should be.
	UpdatedBy string `json:"updatedBy"`
}

func (g GrafanaMetadata) Copy() GrafanaMetadata {
	return GrafanaMetadata{}
}

// TODO guard against skew, use indirection through an internal package
// var _ CommonMetadata = encoding.CommonMetadata{}

// SimpleCustomMetadata is an implementation of CustomMetadata
type SimpleCustomMetadata map[string]any

// MapFields returns a map of string->value for all CustomMetadata fields
func (s SimpleCustomMetadata) MapFields() map[string]any {
	return s
}

// BasicMetadataObject is a composable base struct to attach Metadata, and its associated functions, to another struct.
// BasicMetadataObject provides a Metadata field composed of StaticMetadata and ObjectMetadata, as well as the
// ObjectMetadata(),SetObjectMetadata(), StaticMetadata(), and SetStaticMetadata() receiver functions.
type BasicMetadataObject struct {
	Kind        string               `json:"kind"`
	APIVersion  string               `json:"apiVersion"`
	Metadata    KubeMetadata         `json:"metadata"`
	GrafanaMeta GrafanaMetadata      `json:"grafanaMetadata"`
	KindMeta    SimpleCustomMetadata `json:"kindMetadata"`
}

// StaticMetadata returns the object's StaticMetadata
func (b *BasicMetadataObject) StaticMetadata() StaticMetadata {
	gv := strings.Split(b.APIVersion, "/")
	g := ""
	v := ""
	if len(gv) > 1 {
		g = gv[0]
		v = gv[1]
	} else if len(gv) == 1 {
		// For kubernetes core resources, the group is empty and only the version appears in the APIVersion string
		v = gv[0]
	}
	return StaticMetadata{
		Kind:      b.Kind,
		Group:     g,
		Version:   v,
		Namespace: b.Metadata.Namespace,
		Name:      b.Metadata.Name,
	}
}

// SetStaticMetadata overwrites the StaticMetadata supplied by BasicMetadataObject.StaticMetadata()
// Note that in implementations, this may impact the KubeMetadata, as they have overlapping information.
func (b *BasicMetadataObject) SetStaticMetadata(m StaticMetadata) {
	b.Kind = m.Kind
	if m.Group != "" {
		b.APIVersion = fmt.Sprintf("%s/%s", m.Group, m.Version)
	} else {
		b.APIVersion = m.Version
	}
	b.Metadata.Namespace = m.Namespace
	b.Metadata.Name = m.Name
}

// KubeMetadata returns the object's KubeMetadata
func (b *BasicMetadataObject) KubeMetadata() KubeMetadata {
	return b.Metadata
}

// SetKubeMetadata overwrites the KubeMetadata supplied by BasicMetadataObject.KubeMetadata()
func (b *BasicMetadataObject) SetKubeMetadata(m KubeMetadata) {
	b.Metadata = m
}

// GrafanaMetadata returns the object's GrafanaMetadata
func (b *BasicMetadataObject) GrafanaMetadata() GrafanaMetadata {
	return b.GrafanaMeta
}

// SetGrafanaMetadata overwrites the GrafanaMetadata supplied by BasicMetadataObject.GrafanaMetadata()
func (b *BasicMetadataObject) SetGrafanaMetadata(m GrafanaMetadata) {
	b.GrafanaMeta = m
}

// KindMetadata returns the object's CustomMetadata (kind-specific metadata)
func (b *BasicMetadataObject) KindMetadata() CustomMetadata {
	return b.KindMeta
}

// TODO delete these?

// CopyResource is an implementation of the receiver method `Copy()` required for implementing Resource.
// It should be used in your own runtime.Resource implementations if you do not wish to implement custom behavior.
// Example:
//
//	func (c *CustomResource) Copy() kindsys.Resource {
//	    return resource.CopyResource(c)
//	}
func CopyResource(in any) Resource {
	val := reflect.ValueOf(in).Elem()

	cpy := reflect.New(val.Type())
	cpy.Elem().Set(val)

	// Using the <obj>, <ok> for the type conversion ensures that it doesn't panic if it can't be converted
	if obj, ok := cpy.Interface().(Resource); ok {
		return obj
	}

	// TODO: better return than nil?
	return nil
}

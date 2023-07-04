package kindsys

import "time"

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

// UnmarshalConfig is the config used for unmarshaling Objects.
// It consists of fields that are descriptive of the underlying content, based on knowledge the caller has.
type UnmarshalConfig struct {
	// WireFormat is the wire format of the provided payload
	WireFormat WireFormat
	// VersionHint is what the client thinks the version is (if non-empty)
	VersionHint string
}

// ResourceBytes is the collection of different Kubernetes-shape (see
// [K8sToGrafana]) Resource components as raw bytes.
//
// It is used for unmarshaling a Resource, and can be used for marshaling as well.
// Client implementations are required to process their own storage representation into
// a uniform representation in ResourceBytes.
type ResourceBytes struct {
	// Spec contains the marshaled SpecObject. It should be unmarshalable directly into the Object-implementation's
	// Spec object using an unmarshaler of the appropriate WireFormat type
	Spec []byte
	// Metadata includes object-specific metadata, and may include CommonMetadata depending on implementation.
	// Clients must call SetCommonMetadata on the object after an Unmarshal if CommonMetadata is not provided in the bytes.
	Metadata []byte
	// Subresources contains a map of all subresources that are both part of the underlying Object implementation,
	// AND are supported by the Client implementation. Each entry should be unmarshalable directly into the
	// Object-implementation's relevant subresource using an unmarshaler of the appropriate WireFormat type
	Subresources map[string][]byte
}

// A Resource is a single instance of a Grafana [Kind], either [Core] or [Custom].
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
	// CommonMetadata returns the portion of this Resource's metadata that is common to all kinds.
	CommonMetadata() CommonMetadata

	// SetCommonMetadata overwrites the CommonMetadata of the object.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with CommonMetadata() and set specific values.
	SetCommonMetadata(metadata CommonMetadata)

	// StaticMetadata returns the Object's StaticMetadata
	StaticMetadata() StaticMetadata

	// SetStaticMetadata overwrites the Object's StaticMetadata with the provided StaticMetadata.
	// Implementations should always overwrite, rather than attempt merges of the metadata.
	// Callers wishing to merge should get current metadata with StaticMetadata() and set specific values.
	// Note that StaticMetadata is only mutable in an object create context.
	SetStaticMetadata(metadata StaticMetadata)

	ToBytes(metadataEncoder func(commonMeta *CommonMetadata, customMeta any, format WireFormat) []byte, format WireFormat) (*ResourceBytes, error)

	// FromBytes unmarshals raw bytes from a Kubernetes-form object in the provided
	// WireFormat, the spec object and all provided subresources according to the
	// provided WireFormat. It returns an error if any part of the provided bytes
	// cannot be unmarshaled.
	FromBytes(byt ResourceBytes, config UnmarshalConfig) error

	// Copy returns a full copy of the Resource with all its data.
	Copy() Resource

	// BELOW HERE ARE METHODS WE LIKELY WANT BUT WILL IMPLEMENT AS NEEDED

	// CustomMetadata returns metadata unique to the Kind of this Resource. A
	// resource may have no kind-specific CustomMetadata.
	//
	// CustomMetadata is read-only through this generic interface. Use cases needing
	// to modify custom metadata must know and convert to the underlying type.
	// FIXME finish
	// CustomMetadata() CustomMetadata

	// SpecObject returns the actual "schema" object, which holds the main body of data
	// SpecObject() any

	// Subresources returns a map of subresource name(s) to the object value for that subresource.
	// Spec is not considered a subresource, and should only be returned by SpecObject
	// Subresources() map[string]any
}

// CustomMetadata is an interface describing a resource.Object's kind-specific metadata
// type CustomMetadata interface {
// 	// MapFields converts the custom metadata's fields into a map of field key to value.
// 	// This is used so Clients don't need to engage in reflection for marshaling metadata,
// 	// as various implementations may not store kind-specific metadata the same way.
// 	MapFields() map[string]any
// }

// StaticMetadata consists of all immutable metadata for a Resource.
// It is set in the initial Create call for an Object, then will always remain the same.
type StaticMetadata struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
}

// ListObject represents a List of Object-implementing objects with list metadata.
// The simplest way to use it is to use the implementation returned by a Client's List call.
type ListObject interface {
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

// CommonMetadata is the system-defined common metadata associated with a [Resource].
// It combines Kubernetes standard metadata with certain Grafana-specific additions.
//
// It is analogous to [k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta] in vanilla Kubernetes.
//
// TODO generate this from the CUE definition
// TODO review this for optionality
type CommonMetadata struct {
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
	// UpdateTimestamp is the timestamp of the last update to the resource.
	UpdateTimestamp time.Time `json:"updateTimestamp"`
	// CreatedBy is a string which indicates the user or process which created the resource.
	// Implementations may choose what this indicator should be.
	CreatedBy string `json:"createdBy"`
	// UpdatedBy is a string which indicates the user or process which last updated the resource.
	// Implementations may choose what this indicator should be.
	UpdatedBy string `json:"updatedBy"`

	// ExtraFields stores implementation-specific metadata.
	// Not all Client implementations are required to honor all ExtraFields keys.
	// Generally, this field should be shied away from unless you know the specific
	// Client implementation you're working with and wish to track or mutate extra information.
	ExtraFields map[string]any `json:"extraFields"`
}

// UnstructuredResource is an untyped representation of [Resource]. In the same
// way that map[string]any can represent any JSON []byte, UnstructuredResource
// can represent a [Resource] for any [Core] or [Custom] kind. But it is not
// strongly typed, and lacks any user-defined methods that may exist on a
// kind-specific struct that implements [Resource].
type UnstructuredResource struct {
	CommonMeta CommonMetadata `json:"metadata"`
	KindMeta   map[string]any `json:"customMetadata"`
	Spec       map[string]any `json:"spec,omitempty"`
	Status     map[string]any `json:"status,omitempty"`
}

var _ Resource = UnstructuredResource{}

// SimpleCustomMetadata is an implementation of CustomMetadata
type SimpleCustomMetadata map[string]any

// MapFields returns a map of string->value for all CustomMetadata fields
func (s SimpleCustomMetadata) MapFields() map[string]any {
	return s
}

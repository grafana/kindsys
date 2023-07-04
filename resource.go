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
// - Some use cases need to operate over resources of any kind.
// - Go generics do not allow the ergonomic expression of certain needed constraints.
//
// For cases needing to operate generically, [UnstructuredResource] is
// available. But most Go code works directly with a known, finite set of
// kinds. In these cases, prefer using generated implementations of [Resource] that
// are specific to each kind. (TODO link to some docs)
//
// FIXME finish this doc before merging PR
// Each implementation of Resource also bridges the gap between the byte
// representation of objects, and the native Go representation.
//
// The byte representation is the standard Kubernetes object format. It's
// commonly seen committed to git repositories in YAML. the strongly typed
// representation of the object in Go. The wire format is a JSON representation
// of the object
//
// Resource represents the shared parts of Grafana resources, common regardless
// of the underlying kind. It is generic over its Spec and Status fields. It is
// expected that the generic parameters for Resource are structs, generated from
// a CUE kind definition.
type Resource interface {
	// CommonMetadata returns the portion of this Resource's metadata that is common to all kinds.
	CommonMetadata() CommonMetadata

	// CustomMetadata returns metadata unique to the Kind of this Resource. A
	// resource may have no kind-specific CustomMetadata.
	//
	// CustomMetadata is read-only through this generic interface. Use cases needing
	// to modify custom metadata must know and convert to the underlying type.
	// FIXME finish
	// CustomMetadata() map[string]any

	ToBytes(metadataEncoder func(commonMeta *CommonMetadata, customMeta any, format WireFormat) []byte, format WireFormat) (*ResourceBytes, error)

	// FromBytes unmarshals raw bytes from a Kubernetes-form object in the provided
	// WireFormat, the spec object and all provided subresources according to the
	// provided WireFormat. It returns an error if any part of the provided bytes
	// cannot be unmarshaled.
	FromBytes(bytes []byte, config UnmarshalConfig) error

	// Copy returns a full copy of the Resource with all its data.
	Copy() Resource
}

// KindInfo consists of all non-mutable metadata for an object.
// It is set in the initial Create call for an Object, then will always remain the same.
type KindInfo struct {
	Group     string `json:"group"`
	Version   string `json:"version"`
	Kind      string `json:"kind"`
	Namespace string `json:"namespace"`
	Name      string `json:"name"`
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

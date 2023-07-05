package encoding

// KubernetesShapeBytes are a partially-encoded representation of a []byte that
// is in the standard Kubernetes shape.
type KubernetesShapeBytes struct {
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

// TODO
type GrafanaShapeBytes struct {
	// TODO
	Spec []byte
	// TODO
	Metadata []byte
	// TODO
	CustomMetadata []byte
	// TODO
	Subresources map[string][]byte
}

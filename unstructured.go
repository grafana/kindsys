package kindsys

var _ Resource = UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource]. In the same
// way that map[string]any can represent any JSON []byte, UnstructuredResource
// can represent a [Resource] for any [Core] or [Custom] kind. But it is not
// strongly typed, and lacks any user-defined methods that may exist on a
// kind-specific struct that implements [Resource].
type UnstructuredResource struct {
	Metadata       CommonMetadata `json:"metadata"`
	CustomMetadata map[string]any `json:"customMetadata"`
	Spec           map[string]any `json:"spec,omitempty"`
	Status         map[string]any `json:"status,omitempty"`
}

func (u UnstructuredResource) CommonMetadata() CommonMetadata {
	return u.Metadata
}

func (u UnstructuredResource) SetCommonMetadata(metadata CommonMetadata) {
	u.Metadata = metadata
}

func (u UnstructuredResource) StaticMetadata() StaticMetadata {
	// TODO implement me
	panic("implement me")
}

func (u UnstructuredResource) SetStaticMetadata(metadata StaticMetadata) {
	// TODO implement me
	panic("implement me")
}

func (u UnstructuredResource) Copy() Resource {
	// TODO implement me
	panic("implement me")
}

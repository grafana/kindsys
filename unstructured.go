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
	com := CommonMetadata{
		UID:               u.Metadata.UID,
		ResourceVersion:   u.Metadata.ResourceVersion,
		Labels:            u.Metadata.Labels,
		CreationTimestamp: u.Metadata.CreationTimestamp.UTC(),
		UpdateTimestamp:   u.Metadata.UpdateTimestamp.UTC(),
		CreatedBy:         u.Metadata.CreatedBy,
		UpdatedBy:         u.Metadata.UpdatedBy,
		ExtraFields:       nil,
	}

	copy(u.Metadata.Finalizers, com.Finalizers)
	if u.Metadata.DeletionTimestamp != nil {
		*com.DeletionTimestamp = *(u.Metadata.DeletionTimestamp)
	}

	for k, v := range u.Metadata.Labels {
		com.Labels[k] = v
	}

	return UnstructuredResource{
		Metadata:       com,
		CustomMetadata: mapcopy(u.CustomMetadata),
		Spec:           mapcopy(u.Spec),
		Status:         mapcopy(u.Status),
	}
}

func mapcopy(m map[string]any) map[string]any {
	cp := make(map[string]any)
	for k, v := range m {
		if vm, ok := v.(map[string]any); ok {
			cp[k] = mapcopy(vm)
		} else {
			cp[k] = v
		}
	}

	return cp
}

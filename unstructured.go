package kindsys

var _ Resource = &UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource]. In the same
// way that map[string]any can represent any JSON []byte, UnstructuredResource
// can represent a [Resource] for any [Core] or [Custom] kind. But it is not
// strongly typed, and lacks any user-defined methods that may exist on a
// kind-specific struct that implements [Resource].
type UnstructuredResource struct {
	BasicMetadataObject `json:",inline"`
	Spec                map[string]any `json:"spec,omitempty"`
	Status              map[string]any `json:"status,omitempty"`
	// TODO: is there value in storing other subresources in UnstructuredResource?
}

func (u *UnstructuredResource) SpecObject() any {
	return u.Spec
}

func (u *UnstructuredResource) Subresources() map[string]any {
	return map[string]any{
		"status": u.Status,
	}
}

func (u *UnstructuredResource) Copy() Resource {
	n := UnstructuredResource{
		BasicMetadataObject: BasicMetadataObject{
			Kind:        u.Kind,
			APIVersion:  u.APIVersion,
			Metadata:    u.Metadata.Copy(),
			GrafanaMeta: u.GrafanaMeta.Copy(),
			KindMeta:    SimpleCustomMetadata{},
		},
		Spec:   mapcopy(u.Spec),
		Status: mapcopy(u.Status),
	}
	for k, v := range u.KindMeta {
		n.KindMeta[k] = v
	}
	return &n
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

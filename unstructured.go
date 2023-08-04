package kindsys

import (
	"bytes"
	"fmt"

	jsoniter "github.com/json-iterator/go"
)

var _ Resource = &UnstructuredResource{}

// UnstructuredResource is an untyped representation of [Resource].
type UnstructuredResource struct {
	GenericResource[map[string]any, SimpleCustomMetadata, map[string]any]
}

// UnmarshalJSON allows creating a resource from json
func (u *UnstructuredResource) UnmarshalJSON(b []byte) error {
	return ReadResourceJSON(bytes.NewReader(b), JSONResourceBuilder{
		SetStaticMetadata: func(v StaticMetadata) { u.StaticMeta = v },
		SetCommonMetadata: func(v CommonMetadata) { u.CommonMeta = v },
		ReadSpec: func(iter *jsoniter.Iterator) error {
			u.Spec = make(map[string]any)
			iter.ReadVal(&u.Spec)
			return iter.Error
		},
		SetAnnotation: func(key, val string) {
			fmt.Printf("??? unknown")
		},
		ReadStatus: func(iter *jsoniter.Iterator) error {
			u.Status = make(map[string]any)
			iter.ReadVal(&u.Status)
			return iter.Error
		},
		ReadSub: func(name string, iter *jsoniter.Iterator) error {
			return fmt.Errorf("unsupported sub resource")
		},
	})
}

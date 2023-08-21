package kindsys

import (
	"bytes"
	"encoding/gob"

	jsoniter "github.com/json-iterator/go"
)

// Verify this implements a real resource
var _ Resource = &GenericResource[any, CustomMetadata, any]{}

// GenericResource
type GenericResource[Spec any, CustomMeta CustomMetadata, Status any] struct {
	StaticMeta StaticMetadata
	CommonMeta CommonMetadata
	CustomMeta CustomMeta
	Spec       Spec
	Status     Status
}

func (u *GenericResource[Spec, CustomMeta, Status]) SpecObject() any {
	return u.Spec
}

func (u *GenericResource[Spec, CustomMeta, Status]) Subresources() map[string]any {
	return map[string]any{
		"status": u.Status,
	}
}

func (u *GenericResource[Spec, CustomMeta, Status]) Copy() Resource {
	dst := &GenericResource[Spec, CustomMeta, Status]{}
	buf := bytes.Buffer{}
	err := gob.NewEncoder(&buf).Encode(u)
	if err != nil {
		return dst // error
	}
	_ = gob.NewDecoder(&buf).Decode(dst)
	return dst
}

// MarshalJSON marshals Frame to JSON.
// NOTE: we can not do generic unmarshal because we need functions to create the typed objects
func (u *GenericResource[Spec, CustomMeta, Status]) MarshalJSON() ([]byte, error) {
	cfg := jsoniter.ConfigCompatibleWithStandardLibrary
	stream := cfg.BorrowStream(nil)
	defer cfg.ReturnStream(stream)

	err := WriteResourceJSON(u, stream)
	if err != nil {
		return nil, err
	}
	if stream.Error != nil {
		return nil, err
	}
	return append([]byte(nil), stream.Buffer()...), nil
}

// CommonMetadata returns the object's CommonMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) CommonMetadata() CommonMetadata {
	return u.CommonMeta
}

// SetCommonMetadata overwrites the ObjectMetadata.Common() supplied by BasicMetadataObject.ObjectMetadata()
func (u *GenericResource[Spec, CustomMeta, Status]) SetCommonMetadata(m CommonMetadata) {
	u.CommonMeta = m
}

// StaticMetadata returns the object's StaticMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) StaticMetadata() StaticMetadata {
	return u.StaticMeta
}

// SetStaticMetadata overwrites the StaticMetadata supplied by BasicMetadataObject.StaticMetadata()
func (u *GenericResource[Spec, CustomMeta, Status]) SetStaticMetadata(m StaticMetadata) {
	u.StaticMeta = m
}

// CustomMetadata returns the object's CustomMetadata
func (u *GenericResource[Spec, CustomMeta, Status]) CustomMetadata() CustomMetadata {
	return u.CustomMeta
}

package kindsys

import (
	"encoding/json"

	"github.com/grafana/kindsys/encoding"
	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"
)

type gis struct {
	// TODO
	Spec json.RawMessage `json:"spec"`
	// TODO
	Metadata json.RawMessage `json:"metadata"`
	// TODO
	// CustomMetadata json.RawMessage `json:"customMetadata"`
	// TODO
	// Subresources map[string]json.RawMessage
}

// genericCore is a general representation of a parsed and validated Core kind.
type genericCore struct {
	def Def[CoreProperties]
	lin thema.Lineage
}

func (k genericCore) Validate(b []byte, codec Decoder) error {
	_, err := k.bytesToAnyInstance(b, codec)
	return err
}

func (k genericCore) bytesToAnyInstance(b []byte, codec Decoder) (*thema.Instance, error) {
	// Transform from k8s shape to intermediate grafana shape
	var gb encoding.GrafanaShapeBytes
	gb, err := codec.Decode(b)
	if err != nil {
		return nil, err
	}
	// if gb.Group != k.Group() || gb.Kind != k.Name() {
	// 	return nil, fmt.Errorf("resource is %s.%s, not of kind %s.%s", gb.Group, gb.Kind, k.Group(), k.Name())
	// }
	// TODO make the intermediate type already look like this so we don't have to re-encode/decode
	gj := gis{
		Spec:     gb.Spec,
		Metadata: gb.Metadata,
		// CustomMetadata: gb.CustomMetadata,
		// Subresources: make(map[string]json.RawMessage),
	}
	// for k, v := range gb.Subresources {
	// 	gj.Subresources[k] = v
	// }
	gjb, err := json.Marshal(gj)
	if err != nil {
		return nil, err
	}

	// reuse the cue context already attached to the underlying lineage
	ctx := k.lin.Runtime().Context()
	// decode JSON into a cue.Value
	cval, err := vmux.NewJSONCodec(k.MachineName()+".json").Decode(ctx, gjb)
	if err != nil {
		return nil, err
	}

	// Try the current one first, it's probably what we want
	sch, _ := k.lin.Schema(k.CurrentVersion()) // we verified at bind of this kind that this schema exists
	inst, curvererr := sch.Validate(cval)
	if curvererr != nil {
		for sch := k.lin.First(); sch != nil; sch = sch.Successor() {
			if sch.Version() == k.CurrentVersion() {
				continue
			}
			if inst, err = sch.Validate(cval); err == nil {
				curvererr = nil
				break
			}
		}
	}

	// TODO improve this once thema stacks all schema validation errors https://github.com/grafana/thema/issues/156
	return inst, curvererr
}

func (k genericCore) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

func (k genericCore) Group() string {
	return k.def.Properties.CRD.Group
}

func (k genericCore) New() UnstructuredResource {
	// TODO implement me
	panic("implement me")
}

func (k genericCore) FromBytes(b []byte, codec Decoder) (UnstructuredResource, error) {
	inst, err := k.bytesToAnyInstance(b, codec)
	if err != nil {
		return UnstructuredResource{}, err
	}
	// we have a valid instance! decode into unstructured
	// TODO implement me
	_ = inst
	panic("implement me")
}

func (k genericCore) ToBytes(UnstructuredResource, codec Encoder) ([]byte, error) {
	// TODO implement me
	panic("implement me")
}

var _ Core = genericCore{}

func (k genericCore) Props() SomeKindProperties {
	return k.def.Properties
}

func (k genericCore) Name() string {
	return k.def.Properties.Name
}

func (k genericCore) MachineName() string {
	return k.def.Properties.MachineName
}

func (k genericCore) Maturity() Maturity {
	return k.def.Properties.Maturity
}

func (k genericCore) Def() Def[CoreProperties] {
	return k.def
}

func (k genericCore) Lineage() thema.Lineage {
	return k.lin
}

// TODO docs
func BindCore(rt *thema.Runtime, def Def[CoreProperties], opts ...thema.BindOption) (Core, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	return genericCore{
		def: def,
		lin: lin,
	}, nil
}

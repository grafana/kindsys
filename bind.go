package kindsys

import (
	"encoding/json"

	"github.com/grafana/thema"
	"github.com/grafana/thema/vmux"

	"github.com/grafana/kindsys/encoding"
)

type gis struct {
	// TODO
	Spec json.RawMessage `json:"spec"`
	// TODO
	Metadata json.RawMessage `json:"metadata"`
}

type withLineage interface {
	Kind
	Lineage() thema.Lineage
}

func bytesToAnyInstance(k withLineage, b []byte, codec Decoder) (*thema.Instance, error) {
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
		// TODO Status?
	}
	gjb, err := json.Marshal(gj)
	if err != nil {
		return nil, err
	}

	lin := k.Lineage()
	// reuse the cue context already attached to the underlying lineage
	ctx := lin.Runtime().Context()
	// decode JSON into a cue.Value
	cval, err := vmux.NewJSONCodec(k.MachineName()+".json").Decode(ctx, gjb)
	if err != nil {
		return nil, err
	}

	// Try the current one first, it's probably what we want
	sch, _ := lin.Schema(k.CurrentVersion()) // we verified at bind of this kind that this schema exists
	inst, curvererr := sch.Validate(cval)
	if curvererr != nil {
		for sch := lin.First(); sch != nil; sch = sch.Successor() {
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

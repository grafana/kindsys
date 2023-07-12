package kindsys

import (
	"github.com/grafana/thema"
)

// genericCore is a general representation of a parsed and validated Core kind.
type genericCore struct {
	def Def[CoreProperties]
	lin thema.Lineage
}

func (k genericCore) Validate(b []byte, codec Decoder) error {
	_, err := bytesToAnyInstance(k, b, codec)
	return err
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
	inst, err := bytesToAnyInstance(k, b, codec)
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

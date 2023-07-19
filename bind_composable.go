package kindsys

import "github.com/grafana/thema"

// genericComposable is a general representation of a parsed and validated
// Composable kind.
type genericComposable struct {
	def Def[ComposableProperties]
	lin thema.Lineage
}

func (k genericComposable) Maturity() Maturity {
	//TODO implement me
	panic("implement me")
}

func (k genericComposable) Validate(b []byte, codec Decoder) error {
	// TODO implement me
	panic("implement me")
}

func (k genericComposable) CurrentVersion() thema.SyntacticVersion {
	return k.def.Properties.CurrentVersion
}

var _ Composable = genericComposable{}

func (k genericComposable) Props() SomeKindProperties {
	return k.def.Properties
}

func (k genericComposable) Name() string {
	return k.def.Properties.Name
}

func (k genericComposable) MachineName() string {
	return k.def.Properties.MachineName
}

func (k genericComposable) Def() Def[ComposableProperties] {
	return k.def
}

func (k genericComposable) Lineage() thema.Lineage {
	return k.lin
}

// TODO docs
func BindComposable(rt *thema.Runtime, def Def[ComposableProperties], opts ...thema.BindOption) (Composable, error) {
	lin, err := def.Some().BindKindLineage(rt, opts...)
	if err != nil {
		return nil, err
	}

	return genericComposable{
		def: def,
		lin: lin,
	}, nil
}
